## Elasticity

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata | Standard object’s metadata. More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#metadata | [metav1.ObjectMeta](https://kubernetes.io/docs/api-reference/v1.6/#objectmeta-v1-meta) | false |
| spec | Define all related information to controll the elasticity of a deployment | [ElasticitySpec](#elasticity-spec) | true |
| status | | *[ElasticityStatus](#prometheusstatus) | false |

## Elasticity Spec
The Elasticity Spec describes the three dimensions of the controller: the target object, the buffer which anticipates the future and the queue-based workload 

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| deployment | Which deployment to controll | [DeploymentSpec](#deployment-spec) | true |
| buffer | What is the buffer specification | [BufferSpec](#buffer-spec) | true |
| workload | Specify the queue where the workload is coming from | [WorkloadSpec](#workload-spec) | true |

## Deployment Spec
The deployment specifies the scalability bounds of a certain deployment, and the capacity of a single instance

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| name | The name of the deployment to controll | [string] | true |
| capacity | The processing capacity of a single instance | [int32] | true |
| minReplicas | The minimum number of replicas no matter the value of the reference signal | [int32] | true |
| maxReplicas | The maximum number of replicas - scalability bound | [int32] | true |

## Buffer Spec
Specifies the queue that deployment is bound to

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| initial | Initial size of the buffer | [int32] | true |
| threshold | Threshold to increase the buffer linearly | [int32] | true |

## Workload Spec
Specifies the queue that deployment is bound to

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| queue | The name of the queue the workload is coming | [string] | true |
