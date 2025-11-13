package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AgentSpec defines the desired state of Agent
type AgentSpec struct {
	ControlPlaneURL string       `json:"controlPlaneUrl"`
	AuthKey         string       `json:"authKey,omitempty"`
	AuthKeySecret   SecretRef    `json:"authKeySecret,omitempty"`
	ServiceName     string       `json:"serviceName"`
	Environment     string       `json:"environment"`
	Namespace       string       `json:"namespace,omitempty"`
	Telemetry       TelemetryConfig `json:"telemetry,omitempty"`
	Policy          PolicyConfig    `json:"policy,omitempty"`
	Resources       ResourceRequirements `json:"resources,omitempty"`
}

type SecretRef struct {
	Name      string `json:"name"`
	Key       string `json:"key"`
	Namespace string `json:"namespace,omitempty"`
}

type TelemetryConfig struct {
	BatchSize    int    `json:"batchSize,omitempty"`
	BatchTimeout string `json:"batchTimeout,omitempty"`
	ExportTimeout string `json:"exportTimeout,omitempty"`
	MaxQueueSize int    `json:"maxQueueSize,omitempty"`
}

type PolicyConfig struct {
	Mode               string       `json:"mode,omitempty"`
	AutoEnableBlocking bool         `json:"autoEnableBlocking,omitempty"`
	ObservePeriodHours int          `json:"observePeriodHours,omitempty"`
}

type ResourceRequirements struct {
	Limits   ResourceList `json:"limits,omitempty"`
	Requests ResourceList `json:"requests,omitempty"`
}

type ResourceList map[string]string

// AgentStatus defines the observed state of Agent
type AgentStatus struct {
	Phase         string      `json:"phase,omitempty"`
	Message       string      `json:"message,omitempty"`
	LastUpdateTime metav1.Time `json:"lastUpdateTime,omitempty"`
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
