package main

import (
	"fmt"
	"time"

	"github.com/golang/glog"
	v1 "github.com/klinakuf/caus/pkg/apis/caus.rss.uni-stuttgart.de/v1"
	elasticityclientset "github.com/klinakuf/caus/pkg/client/clientset/versioned"
	lister "github.com/klinakuf/caus/pkg/client/listers/caus/v1"
	extensions "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

//Controller represents the data the controller needs to operate.
type Controller struct {
	// pods gives cached access to pods.
	elasticities                 lister.ElasticityLister
	elasticitiesSynced           cache.InformerSynced
	informer                     cache.SharedIndexInformer
	kubernetes                   *kubernetes.Clientset
	monitor                      Monitor
	timeDurationForNextScaleDown time.Duration
	elasticityClient             *elasticityclientset.Clientset

	// queue is where incoming work is placed to de-dup and to allow "easy"
	// rate limited requeues on errors
	queue workqueue.RateLimitingInterface
}

//Monitor interface for gathering the rate from the queue
//error because it might fail
type Monitor interface {
	GetRate() (float64, error)
}

//NewController constructs an isntance of the controller
func NewController(elasticityInformer cache.SharedIndexInformer,
	elasticityLister lister.ElasticityLister,
	kubernetesClient *kubernetes.Clientset,
	monitor Monitor, elclient *elasticityclientset.Clientset) *Controller {

	var defaultTimeRateLimit, _ = time.ParseDuration(*TimeFrequency + "s")
	var timeDurationForNextScaleDown, _ = time.ParseDuration(*ScaleDown + "m")

	c := &Controller{
		elasticities:       elasticityLister,
		elasticitiesSynced: elasticityInformer.HasSynced,
		informer:           elasticityInformer,
		queue:              workqueue.NewNamedRateLimitingQueue(NewDefaultCAUSRateLimiter(defaultTimeRateLimit), "controller-caus"),
		kubernetes:         kubernetesClient,
		monitor:            monitor,
		timeDurationForNextScaleDown: timeDurationForNextScaleDown,
		elasticityClient:             elclient,
	}

	// register event handlers to fill the queue with elasticity creations, updates and deletions
	elasticityInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			glog.Info("Elasticity definition ADDED")
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				c.queue.Add(key)
			}
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			glog.Info("Elasticity definition UPTADET")

			key, err := cache.MetaNamespaceKeyFunc(new)
			if err == nil {
				c.queue.Add(key)
			}
		},
		DeleteFunc: func(obj interface{}) {
			glog.Info("Elasticity definition DELETED")

			// IndexerInformer uses a delta nodeQueue, therefore for deletes we have to use this
			// key function.
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err == nil {
				c.queue.Add(key)
			}
		},
	})

	return c
}

//Run starts the controller
func (c *Controller) Run(threadiness int, stopCh chan struct{}) {
	// don't let panics crash the process
	defer utilruntime.HandleCrash()
	// make sure the work queue is shutdown which will trigger workers to end
	defer c.queue.ShutDown()

	glog.Infof("Starting CAUS controller")

	// wait for your secondary caches to fill before starting your work
	if !cache.WaitForCacheSync(stopCh, c.elasticitiesSynced) {
		return
	}

	// start up your worker threads based on threadiness.  Some controllers
	// have multiple kinds of workers
	for i := 0; i < threadiness; i++ {
		// runWorker will loop until "something bad" happens.  The .Until will
		// then rekick the worker after one second
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	// wait until we're told to stop
	<-stopCh
	glog.Infof("Shutting down CAUS controller")
}

func (c *Controller) runWorker() {
	// hot loop until we're told to stop.  processNextWorkItem will
	// automatically wait until there's work available, so we don't worry
	// about secondary waits
	for c.processNextWorkItem() {
	}
}

func (c *Controller) syncHandler(key string) error {
	glog.Infof("Processing key %s", key)
	obj, exists, err := c.informer.GetIndexer().GetByKey(key)

	if err != nil {
		return err
	}

	if !exists {
		glog.Infof("Object does not exist!")
		return nil
	}

	elasticity := obj.(*v1.Elasticity)
	elasticityCopy := elasticity.DeepCopy()
	glog.Infof("Elasticity with name %s processed", elasticity.Name)

	currentRate, err := c.monitor.GetRate()

	if err != nil {
		glog.Errorf("Failed to fetch data from the monitor: %v\n", err)
		return err
	}

	scale, err := c.kubernetes.
		ExtensionsV1beta1().
		Scales(elasticity.Namespace).
		Get("Deployment", elasticity.Spec.Deployment.Name)

	if err != nil {
		glog.Errorf("Obtaining the scale subresource for deployment: %v\n", err)
		return err
	}

	desiredReplicas, bufferedReplicas := CalcReplicas(elasticityCopy, currentRate, float64(scale.Status.Replicas))
	glog.Infof("Computed desired: %s and buffered: %s", desiredReplicas, bufferedReplicas)

	lastScalingDecision := elasticityCopy.Status.LastScaleTime

	glog.Infof("[caus] computed replicas: %v\n", desiredReplicas)
	glog.Infof("[caus] lastScalingDecision: %v\n", lastScalingDecision)

	// check if deployment is specified with 0 replicas in spec
	if scale.Spec.Replicas == 0 || desiredReplicas == scale.Spec.Replicas {
		glog.Infof("[<->] no need to perform scaling")
		//metrics.PushNumberOfReplicas(float64(scale.Spec.Replicas))
		return nil
	}

	// check scalability bound if it is exceeded
	maxReplicasAllowed := elasticityCopy.Spec.Deployment.MaxReplicas
	if scale.Status.Replicas > maxReplicasAllowed {
		glog.Infof("[<-] need to perform scale in without time dependency")
		desiredReplicas = elasticityCopy.Spec.Deployment.MaxReplicas
		return c.scaleTo(desiredReplicas, bufferedReplicas, scale, elasticityCopy)
	}

	// scaleDown decision
	if desiredReplicas < scale.Spec.Replicas && time.Now().After(lastScalingDecision.Add(c.timeDurationForNextScaleDown)) {
		glog.Infof("[<-] SCALING IN TO %v \n", desiredReplicas)
		return c.scaleTo(desiredReplicas, bufferedReplicas, scale, elasticityCopy)
	}

	//scaleUp decision
	if desiredReplicas > scale.Spec.Replicas {
		glog.Infof("[->] SCALING OUT TO %v \n", desiredReplicas)
		return c.scaleTo(desiredReplicas, bufferedReplicas, scale, elasticityCopy)

	}

	return nil
}

// processNextWorkItem deals with one key off the queue.  It returns false
// when it's time to quit.
func (c *Controller) processNextWorkItem() bool {
	// pull the next work item from queue.  It should be a key we use to lookup
	// something in a cache
	key, quit := c.queue.Get()
	if quit {
		return false
	}
	// you always have to indicate to the queue that you've completed a piece of
	// work
	defer c.queue.Done(key)

	// do your work on the key.  This method will contains your "do stuff" logic
	err := c.syncHandler(key.(string))
	if err == nil {
		// if you had no error, tell the queue to stop tracking history for your
		// key. This will reset things like failure counts for per-item rate
		// limiting
		c.queue.AddRateLimited(key)
		return true
	}

	// there was a failure so be sure to report it.  This method allows for
	// pluggable error handling which can be used for things like
	// cluster-monitoring
	utilruntime.HandleError(fmt.Errorf("%v failed with : %v", key, err))

	// since we failed, we should requeue the item to work on later.  This
	// method will add a backoff to avoid hotlooping on particular items
	// (they're probably still not going to work right away) and overall
	// controller protection (everything I've done is broken, this controller
	// needs to calm down or it can starve other useful work) cases.
	c.queue.AddRateLimited(key)

	return true
}

func (c *Controller) scaleTo(numberOfReplicas int32, bufferedReplicas int32, scale *extensions.Scale, elasticity *v1.Elasticity) error {
	//, scale *extensions.Scale
	now := metav1.NewTime(time.Now())
	// // update the deployment
	currentReplicas := scale.Spec.Replicas
	scale.Spec.Replicas = numberOfReplicas
	_, err := c.kubernetes.ExtensionsV1beta1().Scales("payroll").Update("Deployment", scale)

	if err != nil {
		glog.Errorf("updating the scale subresource: %v\n", err)
		return err
	}

	//metrics.PushNumberOfReplicas(float64(numberOfReplicas))

	//update elsasticity Status
	elasticity.Status = v1.ElasticityStatus{
		Message:          "Successfully scaled Deployment by controller",
		CurrentReplicas:  currentReplicas,
		BufferedReplicas: bufferedReplicas,
		DesiredReplicas:  numberOfReplicas,
		LastScaleTime:    &now,
	}

	elasticityUpdated, err := c.elasticityClient.CausV1().Elasticities(elasticity.ObjectMeta.Namespace).Update(elasticity)

	//TODO: check if somethign should be handled here.
	if err != nil {
		glog.Errorf("updating status: %v\n", err)
		return err
	}

	glog.Infof("Elasticity updated %s", elasticityUpdated)

	return nil
}
