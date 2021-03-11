/*
Copyright 2021 Tutkovics András.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ServicegraphSpec defines the desired state of Servicegraph
type ServicegraphSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:Required
	// Nodes contain the configs to service
	Nodes []*Node `json:"nodes,omitempty"`
}

// ServicegraphStatus defines the observed state of Servicegraph
type ServicegraphStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Nodes []string `json:"nodes"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Servicegraph is the Schema for the servicegraphs API
type Servicegraph struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ServicegraphSpec   `json:"spec,omitempty"`
	Status ServicegraphStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ServicegraphList contains a list of Servicegraph
type ServicegraphList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Servicegraph `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Servicegraph{}, &ServicegraphList{})
}

// Here come my structures

// Node struct. Contains specification of a node
type Node struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="[a-z-]*"
	// Name of the node
	Name string `json:"name"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=0
	// Number of replica to run
	Replicas uint `json:"replicas"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=0
	// Container port to open and listen
	ContainerPort uint `json:"port"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Minimum=0
	// If the service will listen on node port
	NodePort uint `json:"nodePort"`

	// +kubebuilder:validation:Required
	// Resource to consume
	Resources Resource `json:"resources"`

	// +kubebuilder:validation:Required
	// Setup and configure endpoints to node
	Endpoints []Endpoint `json:"endpoints"`

	// +kubebuilder:validation:Optional
	// HPA scaling configs
	HPA HPA `json:"hpa"`
}

// Contain HPA related configs
type HPA struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=0
	// CPU target utilization
	Utilization uint `json:"utilization"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=0
	// Minimum replicas to run
	MinReplicas uint `json:"min_replicas"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=0
	// Maximum replicas to run
	MaxReplicas uint `json:"max_replicas"`
}

// Resource structure
type Resource struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=0
	// Memory to use (kB)
	Memory uint `json:"memory"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=0
	// CPU to use (mCPU)
	CPU uint `json:"cpu"`
}

// Endpoint structure and behavior
type Endpoint struct {

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Type=string
	// Regex for pattern eg: /index
	// +kubebuilder:validation:Pattern="/[a-z]*"
	// Path to listen and answer
	Path string `json:"path"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Minimum=0
	// CPU usage until request
	CPULoad uint `json:"cpuLoad"`

	//TODO: Memory load

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=0
	// Delay / process time to "serve" request (ms)
	Delay uint `json:"delay"`

	// +kubebuilder:validation:Optional
	// Ask further services
	CallOuts []CallOut `json:"callouts,omitempty"`
}

// CallOut structure contains additional information to each call out
type CallOut struct {
	// Regex for pattern eg: db-user:890/read?from=table#site

	// +kubebuilder:validation:Required
	// Url to request
	URL string `json:"url,omitempty"`
	// todo: +kubebuilder:validation:Pattern="[a-z-]*:[0-9]*/[a-z?=#]*"
}
