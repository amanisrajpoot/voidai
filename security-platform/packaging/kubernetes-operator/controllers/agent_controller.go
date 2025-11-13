package controllers

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	securityplatformv1 "github.com/security-platform/kubernetes-operator/api/v1"
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

	var agent securityplatformv1.Agent
	if err := r.Get(ctx, req.NamespacedName, &agent); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Check if Deployment exists
	deployment := &appsv1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{
		Name:      agent.Name,
		Namespace: agent.Namespace,
	}, deployment)

	if err != nil && errors.IsNotFound(err) {
		// Create Deployment
		if err := r.createDeployment(ctx, &agent); err != nil {
			logger.Error(err, "Failed to create Deployment")
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		return ctrl.Result{}, err
	}

	// Update Deployment if needed
	if r.needsUpdate(&agent, deployment) {
		if err := r.updateDeployment(ctx, &agent, deployment); err != nil {
			logger.Error(err, "Failed to update Deployment")
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// Update status
	if err := r.updateStatus(ctx, &agent, deployment); err != nil {
		logger.Error(err, "Failed to update status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
}

func (r *AgentReconciler) createDeployment(ctx context.Context, agent *securityplatformv1.Agent) error {
	deployment := r.buildDeployment(agent)
	if err := ctrl.SetControllerReference(agent, deployment, r.Scheme); err != nil {
		return err
	}
	return r.Create(ctx, deployment)
}

func (r *AgentReconciler) updateDeployment(ctx context.Context, agent *securityplatformv1.Agent, deployment *appsv1.Deployment) error {
	newDeployment := r.buildDeployment(agent)
	deployment.Spec = newDeployment.Spec
	return r.Update(ctx, deployment)
}

func (r *AgentReconciler) buildDeployment(agent *securityplatformv1.Agent) *appsv1.Deployment {
	replicas := agent.Spec.Replicas
	if replicas == 0 {
		replicas = 1
	}

	image := agent.Spec.Image
	if image == "" {
		image = "securityplatform/agent:latest"
	}

	imagePullPolicy := corev1.PullPolicy(agent.Spec.ImagePullPolicy)
	if imagePullPolicy == "" {
		imagePullPolicy = corev1.PullIfNotPresent
	}

	envVars := []corev1.EnvVar{
		{
			Name:  "SECURITY_PLATFORM_CONTROL_PLANE_URL",
			Value: agent.Spec.ControlPlaneURL,
		},
		{
			Name: "SECURITY_PLATFORM_AUTH_KEY",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: agent.Spec.AuthKeySecretRef.Name,
					},
					Key: agent.Spec.AuthKeySecretRef.Key,
				},
			},
		},
		{
			Name:  "SECURITY_PLATFORM_SERVICE_NAME",
			Value: agent.Spec.ServiceName,
		},
		{
			Name:  "SECURITY_PLATFORM_ENVIRONMENT",
			Value: agent.Spec.Environment,
		},
	}

	if agent.Spec.OtlpEndpoint != "" {
		envVars = append(envVars, corev1.EnvVar{
			Name:  "SECURITY_PLATFORM_OTLP_ENDPOINT",
			Value: agent.Spec.OtlpEndpoint,
		})
	}

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      agent.Name,
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
							Name:            "agent",
							Image:           image,
							ImagePullPolicy: imagePullPolicy,
							Env:             envVars,
						},
					},
				},
			},
		},
	}
}

func (r *AgentReconciler) needsUpdate(agent *securityplatformv1.Agent, deployment *appsv1.Deployment) bool {
	if *deployment.Spec.Replicas != agent.Spec.Replicas {
		return true
	}
	if len(deployment.Spec.Template.Spec.Containers) == 0 {
		return true
	}
	container := deployment.Spec.Template.Spec.Containers[0]
	if container.Image != agent.Spec.Image && agent.Spec.Image != "" {
		return true
	}
	return false
}

func (r *AgentReconciler) updateStatus(ctx context.Context, agent *securityplatformv1.Agent, deployment *appsv1.Deployment) error {
	agent.Status.Replicas = deployment.Status.Replicas
	agent.Status.ReadyReplicas = deployment.Status.ReadyReplicas

	condition := metav1.Condition{
		Type:               "Ready",
		Status:             metav1.ConditionTrue,
		Reason:             "DeploymentReady",
		Message:            fmt.Sprintf("Deployment has %d ready replicas", deployment.Status.ReadyReplicas),
		LastTransitionTime: metav1.Now(),
	}

	if deployment.Status.ReadyReplicas != deployment.Status.Replicas {
		condition.Status = metav1.ConditionFalse
		condition.Reason = "DeploymentNotReady"
		condition.Message = fmt.Sprintf("Deployment has %d/%d ready replicas", deployment.Status.ReadyReplicas, deployment.Status.Replicas)
	}

	agent.Status.Conditions = []metav1.Condition{condition}

	return r.Status().Update(ctx, agent)
}

// SetupWithManager sets up the controller with the Manager.
func (r *AgentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&securityplatformv1.Agent{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
