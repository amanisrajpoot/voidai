package main

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var (
	scheme   = runtime.NewScheme()
	groupVersion = schema.GroupVersion{Group: "securityplatform.io", Version: "v1"}
)

func init() {
	metav1.AddToGroupVersion(scheme, groupVersion)
}

// Agent represents the Security Platform Agent CRD
// +kubebuilder:object:root=true
type Agent struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              AgentSpec   `json:"spec,omitempty"`
	Status            AgentStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
type AgentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Agent `json:"items"`
}

type AgentSpec struct {
	ControlPlaneURL string `json:"controlPlaneUrl"`
	AuthKey         string `json:"authKey"`
	ServiceName     string `json:"serviceName"`
	Environment     string `json:"environment"`
	OtlpEndpoint     string `json:"otlpEndpoint"`
	LocalPolicy     string `json:"localPolicy"`
}

type AgentStatus struct {
	Phase       string    `json:"phase"`
	Message     string    `json:"message"`
	LastUpdated time.Time `json:"lastUpdated"`
}

// ReconcileAgent reconciles Agent resources
type ReconcileAgent struct {
	client client.Client
	scheme *runtime.Scheme
}

func (r *ReconcileAgent) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	log := log.FromContext(ctx)
	
	agent := &Agent{}
	if err := r.client.Get(ctx, req.NamespacedName, agent); err != nil {
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	// Update status
	agent.Status.Phase = "Running"
	agent.Status.Message = "Agent deployed successfully"
	agent.Status.LastUpdated = time.Now()

	if err := r.client.Status().Update(ctx, agent); err != nil {
		log.Error(err, "Failed to update agent status")
		return reconcile.Result{}, err
	}

	// Create or update DaemonSet
	if err := r.reconcileDaemonSet(ctx, agent); err != nil {
		log.Error(err, "Failed to reconcile DaemonSet")
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileAgent) reconcileDaemonSet(ctx context.Context, agent *Agent) error {
	clientset, err := kubernetes.NewForConfig(config.GetConfigOrDie())
	if err != nil {
		return err
	}

	daemonSet := &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-agent", agent.Name),
			Namespace: agent.Namespace,
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "security-platform-agent",
					"agent": agent.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "security-platform-agent",
						"agent": agent.Name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "agent",
							Image: "securityplatform/agent:latest",
							Env: []corev1.EnvVar{
								{Name: "SECURITY_PLATFORM_CONTROL_PLANE_URL", Value: agent.Spec.ControlPlaneURL},
								{Name: "SECURITY_PLATFORM_AUTH_KEY", Value: agent.Spec.AuthKey},
								{Name: "SECURITY_PLATFORM_SERVICE_NAME", Value: agent.Spec.ServiceName},
								{Name: "SECURITY_PLATFORM_ENVIRONMENT", Value: agent.Spec.Environment},
								{Name: "SECURITY_PLATFORM_OTLP_ENDPOINT", Value: agent.Spec.OtlpEndpoint},
							},
						},
					},
				},
			},
		},
	}

	// Create or update DaemonSet
	_, err = clientset.AppsV1().DaemonSets(agent.Namespace).Create(ctx, daemonSet, metav1.CreateOptions{})
	if err != nil {
		// Try update if exists
		_, err = clientset.AppsV1().DaemonSets(agent.Namespace).Update(ctx, daemonSet, metav1.UpdateOptions{})
		return err
	}

	return nil
}

func main() {
	log.SetLogger(zap.New())

	cfg, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	mgr, err := manager.New(cfg, manager.Options{Scheme: scheme})
	if err != nil {
		panic(err)
	}

	// Register CRD
	if err := scheme.AddKnownTypes(groupVersion, &Agent{}, &AgentList{}); err != nil {
		panic(err)
	}

	// Setup controller
	c, err := controller.New("agent-controller", mgr, controller.Options{
		Reconciler: &ReconcileAgent{client: mgr.GetClient(), scheme: scheme},
	})
	if err != nil {
		panic(err)
	}

	// Watch Agent resources
	if err := c.Watch(&source.Kind{Type: &Agent{}}, &handler.EnqueueRequestForObject{}); err != nil {
		panic(err)
	}

	// Start manager
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		panic(err)
	}
}
