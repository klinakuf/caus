# CAUS - A Custom Autoscaler for Containerized Microservices

A **prototype** elasticity controller that complements predictive approaches. It is implemented as an Kubernetes Controller and offers the Elasticity CRD extension

The controller processes definitions of kind Elasticity which cover three aspects of the controller. 
The controller can be depicted as follows: 

reference signal ---> | Elasticity Controller | ---> input ---> | Deployment.Scale | ---> output

# migration in progress

CAUS (Custom Autoscaler) is currently being migrated to work with Kubernetes CRDs

## Getting Started

First register the custom resource definition:

```
kubectl apply -f artifacts/elasticity-crd.yaml
```

Then add an example of the `Elasticity` kind:

```
kubectl apply -f artifacts/my-elastic-rules.yaml
```

Finally build and run the controller:

```
cd cmd/controller
go build
./controller -kubeconfig $HOME/.kube/config --logtostderr=1
```
The last line assumes that the kubernetes config file is located at the $HOME/.kube/ default directory.

