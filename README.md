# caus-crd

CAUS (Custom Autoscaler) migrated to work with Kubernetes CRDs

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

