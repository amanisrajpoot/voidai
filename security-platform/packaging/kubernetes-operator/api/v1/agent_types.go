package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AgentSpec defines the desired state of Agent
type AgentSpec struct {
	// ControlPlaneURL is the URL of the control plane
	ControlPlaneURL string `json:"controlPlaneUrl"`

	// AuthKey is the authentication key (stored in Secret)
	AuthKeySecretRef SecretRef `json:"authKeySecretRef"`

	// ServiceName is the name of the service
	ServiceName string `json:"serviceName"`

	// Environment is the deployment environment
	Environment string `json:"environment,omitempty"`

	// Namespace is the Kubernetes namespace
	Namespace string `json:"namespace,omitempty"`

	// OtlpEndpoint is the OTLP endpoint (optional)
	OtlpEndpoint string `json:"otlpEndpoint,omitempty"`

	// Telemetry configuration
	Telemetry TelemetryConfig `json:"telemetry,omitempty"`

	// Policy configuration
	Policy PolicyConfig `json:"policy,omitempty"`

	// Replicas is the number of agent replicas
	Replicas int32 `json:"replicas,omitempty"`

	// Image is the container image
	Image string `json:"image,omitempty"`

	// ImagePullPolicy
	ImagePullPolicy string `json:"imagePullPolicy,omitempty"`
}

type SecretRef struct {
	Name      string `json:"name"`
	Key       string `json:"key"`
	Namespace string `json:"namespace,omitempty"`
}

type TelemetryConfig struct {
	BatchSize     int    `json:"batchSize,omitempty"`
	BatchTimeout  string `json:"batchTimeout,omitempty"`
	ExportTimeout string `json:"exportTimeout,omitempty"`
	MaxQueueSize  int    `json:"maxQueueSize,omitempty"`
}

type PolicyConfig struct {
	Mode               string                   `json:"mode,omitempty"`
	AutoEnableBlocking bool                     `json:"autoEnableBlocking,omitempty"`
	ObservePeriodHours int                      `json:"observePeriodHours,omitempty"`
	Rules              []map[string]interface{} `json:"rules,omitempty"`
}

// AgentStatus defines the observed state of Agent
type AgentStatus struct {
	// Replicas is the number of running replicas
	Replicas int32 `json:"replicas"`

	// ReadyReplicas is the number of ready replicas
	ReadyReplicas int32 `json:"readyReplicas"`

	// Conditions represent the latest available observations
	Conditions []metav1.Condition `json:"conditions,omitempty"`
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
