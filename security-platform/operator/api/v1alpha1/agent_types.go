package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AgentSpec defines the desired state of Agent
type AgentSpec struct {
	// ControlPlaneURL is the URL of the control plane
	ControlPlaneURL string `json:"controlPlaneURL"`

	// AuthKey is the authentication key (can be set via Secret)
	AuthKey string `json:"authKey,omitempty"`

	// AuthKeySecretRef references a secret containing the auth key
	AuthKeySecretRef *SecretRef `json:"authKeySecretRef,omitempty"`

	// ServiceName is the name of the service
	ServiceName string `json:"serviceName"`

	// Environment is the deployment environment
	Environment string `json:"environment,omitempty"`

	// OtlpEndpoint is the OTLP endpoint (optional)
	OtlpEndpoint string `json:"otlpEndpoint,omitempty"`

	// Replicas is the number of agent replicas
	Replicas *int32 `json:"replicas,omitempty"`

	// Image is the agent container image
	Image string `json:"image,omitempty"`

	// Resources are the resource requirements
	Resources *ResourceRequirements `json:"resources,omitempty"`
}

// SecretRef references a Kubernetes secret
type SecretRef struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

// ResourceRequirements defines resource requirements
type ResourceRequirements struct {
	Limits   *ResourceList `json:"limits,omitempty"`
	Requests *ResourceList `json:"requests,omitempty"`
}

// ResourceList defines resource quantities
type ResourceList struct {
	CPU    string `json:"cpu,omitempty"`
	Memory string `json:"memory,omitempty"`
}

// AgentStatus defines the observed state of Agent
type AgentStatus struct {
	// Conditions represent the latest available observations
	Conditions []Condition `json:"conditions,omitempty"`

	// Replicas is the number of running replicas
	Replicas int32 `json:"replicas,omitempty"`

	// ReadyReplicas is the number of ready replicas
	ReadyReplicas int32 `json:"readyReplicas,omitempty"`
}

// Condition represents a condition
type Condition struct {
	Type               string             `json:"type"`
	Status             metav1.ConditionStatus `json:"status"`
	LastTransitionTime metav1.Time        `json:"lastTransitionTime"`
	Reason             string             `json:"reason,omitempty"`
	Message            string             `json:"message,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Agent is the Schema for the agents API
type Agent struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AgentSpec   `json:"spec,omitempty"`
	Status AgentStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AgentList contains a list of Agent
type AgentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Agent `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Agent{}, &AgentList{})
}
