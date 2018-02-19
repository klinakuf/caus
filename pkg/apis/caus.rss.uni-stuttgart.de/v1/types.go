package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Elasticity struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              ElasticitySpec   `json:"spec"`
	Status            ElasticityStatus `json:"status,omitempty"`
}

type ElasticitySpec struct {
	Deployment DeploymentSpec `json:"deployment"`
	Buffer     BufferSpec     `json:"buffer"`
	Workload   WorkloadSpec   `json:"workload"`
}

type BufferSpec struct {
	Initial int32 `json:"initial"`
	Threshold int32 `json:"threshold"`
}

type DeploymentSpec struct {
	Name        string `json:"name"`
	Capacity    int32  `json:"capacity"`
	MinReplicas *int32 `json:"minReplicas"`
	MaxReplicas int32  `json:"maxReplicas"`
}

type WorkloadSpec struct {
	// Capacity of a single instance ... here we want a list of queue name, capacity, service
	// Classes:
	//   - name: "interactive"
	//     queue: simplePyRun
	//     capacity: 8
	//   - name: "batch"
	//     queue: simplePyRunBatch
	//     capacity: 4
	Queue string `json:"queue"`
}

type ElasticityStatus struct {
	Message          string       `json:"message,omitempty"`
	CurrentReplicas  int32        `json:"currentReplicas"`
	BufferedReplicas int32        `json:"bufferedReplicas"`
	DesiredReplicas  int32        `json:"desiredReplicas"`
	LastScaleTime    *metav1.Time `json:"lastScaleTime,omitempty" protobuf:"bytes,2,opt,name=lastScaleTime"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ElasticityList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Elasticity `json:"items"`
}