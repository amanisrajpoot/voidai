package controllers

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	securityplatformv1alpha1 "github.com/security-platform/operator/api/v1alpha1"
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
//+kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch

func (r *AgentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var agent securityplatformv1alpha1.Agent
	if err := r.Get(ctx, req.NamespacedName, &agent); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Get auth key from secret if specified
	authKey := agent.Spec.AuthKey
	if agent.Spec.AuthKeySecretRef != nil {
		var secret corev1.Secret
		if err := r.Get(ctx, client.ObjectKey{
			Namespace: req.Namespace,
			Name:      agent.Spec.AuthKeySecretRef.Name,
		}, &secret); err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to get secret: %w", err)
		}
		authKey = string(secret.Data[agent.Spec.AuthKeySecretRef.Key])
	}

	// Create or update Deployment
	deployment := r.buildDeployment(&agent, authKey)
	if err := ctrl.SetControllerReference(&agent, deployment, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	var existingDeployment appsv1.Deployment
	if err := r.Get(ctx, client.ObjectKeyFromObject(deployment), &existingDeployment); err != nil {
		if client.IgnoreNotFound(err) == nil {
			logger.Info("Creating Deployment")
			if err := r.Create(ctx, deployment); err != nil {
				return ctrl.Result{}, err
			}
		} else {
			return ctrl.Result{}, err
		}
	} else {
		logger.Info("Updating Deployment")
		if err := r.Update(ctx, deployment); err != nil {
			return ctrl.Result{}, err
		}
	}

	// Update status
	agent.Status.Replicas = *deployment.Spec.Replicas
	agent.Status.ReadyReplicas = existingDeployment.Status.ReadyReplicas
	if err := r.Status().Update(ctx, &agent); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *AgentReconciler) buildDeployment(agent *securityplatformv1alpha1.Agent, authKey string) *appsv1.Deployment {
	replicas := int32(1)
	if agent.Spec.Replicas != nil {
		replicas = *agent.Spec.Replicas
	}

	image := agent.Spec.Image
	if image == "" {
		image = "securityplatform/agent:latest"
	}

	env := []corev1.EnvVar{
		{Name: "SECURITY_PLATFORM_CONTROL_PLANE_URL", Value: agent.Spec.ControlPlaneURL},
		{Name: "SECURITY_PLATFORM_AUTH_KEY", Value: authKey},
		{Name: "SECURITY_PLATFORM_SERVICE_NAME", Value: agent.Spec.ServiceName},
		{Name: "SECURITY_PLATFORM_ENVIRONMENT", Value: agent.Spec.Environment},
	}

	if agent.Spec.OtlpEndpoint != "" {
		env = append(env, corev1.EnvVar{
			Name:  "SECURITY_PLATFORM_OTLP_ENDPOINT",
			Value: agent.Spec.OtlpEndpoint,
		})
	}

	resources := corev1.ResourceRequirements{}
	if agent.Spec.Resources != nil {
		if agent.Spec.Resources.Limits != nil {
			resources.Limits = corev1.ResourceList{}
			if agent.Spec.Resources.Limits.CPU != "" {
				resources.Limits[corev1.ResourceCPU] = resource.MustParse(agent.Spec.Resources.Limits.CPU)
			}
			if agent.Spec.Resources.Limits.Memory != "" {
				resources.Limits[corev1.ResourceMemory] = resource.MustParse(agent.Spec.Resources.Limits.Memory)
			}
		}
		if agent.Spec.Resources.Requests != nil {
			resources.Requests = corev1.ResourceList{}
			if agent.Spec.Resources.Requests.CPU != "" {
				resources.Requests[corev1.ResourceCPU] = resource.MustParse(agent.Spec.Resources.Requests.CPU)
			}
			if agent.Spec.Resources.Requests.Memory != "" {
				resources.Requests[corev1.ResourceMemory] = resource.MustParse(agent.Spec.Resources.Requests.Memory)
			}
		}
	}

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      agent.Name + "-agent",
			Namespace: agent.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
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
							Image: image,
							Env:   env,
							Resources: resources,
						},
					},
				},
			},
		},
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *AgentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&securityplatformv1alpha1.Agent{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
