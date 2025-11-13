package controllers

import (
	"context"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	securityplatformv1 "github.com/securityplatform/operator/api/v1"
)

// AgentReconciler reconciles a Agent object
type AgentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=securityplatform.io,resources=agents,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=securityplatform.io,resources=agents/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=securityplatform.io,resources=agents/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch

func (r *AgentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var agent securityplatformv1.Agent
	if err := r.Get(ctx, req.NamespacedName, &agent); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Create ConfigMap for agent configuration
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      agent.Name + "-config",
			Namespace: agent.Namespace,
		},
		Data: map[string]string{
			"config.yaml": r.generateConfig(&agent),
		},
	}

	if err := ctrl.SetControllerReference(&agent, configMap, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.CreateOrUpdate(ctx, configMap); err != nil {
		logger.Error(err, "failed to create/update ConfigMap")
		return ctrl.Result{}, err
	}

	// Create Deployment
	deployment := r.buildDeployment(&agent, configMap)
	if err := ctrl.SetControllerReference(&agent, deployment, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.CreateOrUpdate(ctx, deployment); err != nil {
		logger.Error(err, "failed to create/update Deployment")
		return ctrl.Result{}, err
	}

	// Update status
	agent.Status.Phase = "Running"
	agent.Status.Message = "Agent deployment created successfully"
	agent.Status.LastUpdateTime = metav1.Now()
	if err := r.Status().Update(ctx, &agent); err != nil {
		logger.Error(err, "failed to update status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *AgentReconciler) generateConfig(agent *securityplatformv1.Agent) string {
	// Generate YAML config from Agent spec
	// Simplified - in production use proper YAML marshaling
	return `control_plane_url: "` + agent.Spec.ControlPlaneURL + `"
service_name: "` + agent.Spec.ServiceName + `"
environment: "` + agent.Spec.Environment + `"
`
}

func (r *AgentReconciler) buildDeployment(agent *securityplatformv1.Agent, configMap *corev1.ConfigMap) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      agent.Name,
			Namespace: agent.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": agent.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": agent.Name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "agent",
							Image: "securityplatform/agent:latest",
							Args: []string{
								"-config=/etc/security-platform/config.yaml",
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "config",
									MountPath: "/etc/security-platform",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: configMap.Name,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *AgentReconciler) CreateOrUpdate(ctx context.Context, obj client.Object) error {
	existing := obj.DeepCopyObject().(client.Object)
	err := r.Get(ctx, client.ObjectKeyFromObject(obj), existing)
	if err != nil {
		if client.IgnoreNotFound(err) == nil {
			return r.Create(ctx, obj)
		}
		return err
	}
	return r.Update(ctx, obj)
}

func int32Ptr(i int32) *int32 { return &i }

func (r *AgentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&securityplatformv1.Agent{}).
		Complete(r)
}
