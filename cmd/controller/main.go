package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/golang/glog"
	v1 "github.com/klinakuf/crd-code-generation/pkg/apis/caus.rss.uni-stuttgart.de/v1"
	elasticityclientset "github.com/klinakuf/crd-code-generation/pkg/client/clientset/versioned"
	"github.com/klinakuf/crd-code-generation/pkg/client/informers/externalversions"
	lister "github.com/klinakuf/crd-code-generation/pkg/client/listers/caus/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"
)

var (
	kuberconfig = flag.String("kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	master      = flag.String("master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
)

func main() {
	flag.Parse()

	cfg, err := clientcmd.BuildConfigFromFlags(*master, *kuberconfig)
	if err != nil {
		glog.Fatalf("Error building kubeconfig: %v", err)
	}
	elasticityClient, err := elasticityclientset.NewForConfig(cfg)

	elSharedInformerFactory := externalversions.NewSharedInformerFactory(elasticityClient, 0)
	elInformer := elSharedInformerFactory.Caus().V1().Elasticities().Informer()
	elLister := elSharedInformerFactory.Caus().V1().Elasticities().Lister()
	controller := NewController(elInformer, elLister)

	stop := make(chan struct{})
	defer close(stop)
	fmt.Println("Controller will run")
	elSharedInformerFactory.Start(stop)
	controller.Run(1, stop)
}

type Controller struct {
	// pods gives cached access to pods.
	elasticities       lister.ElasticityLister
	elasticitiesSynced cache.InformerSynced
	informer           cache.SharedIndexInformer

	// queue is where incoming work is placed to de-dup and to allow "easy"
	// rate limited requeues on errors
	queue workqueue.RateLimitingInterface
}

func NewController(elasticityInformer cache.SharedIndexInformer, elasticityLister lister.ElasticityLister) *Controller {

	c := &Controller{
		elasticities:       elasticityLister,
		elasticitiesSynced: elasticityInformer.HasSynced,
		informer:           elasticityInformer,
		queue:              workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "controller-caus"),
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

	if exists {
		elasticity := obj.(*v1.Elasticity)
		glog.Infof("Elasticity with name %s processed", elasticity.Name)
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
		c.queue.Forget(key)
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
