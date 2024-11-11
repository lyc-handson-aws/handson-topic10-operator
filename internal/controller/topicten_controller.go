/*
Copyright 2024.

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

package controller

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appv1alpha1 "github.com/lyc-handson-aws/handson-topic10-operator/api/v1alpha1"
)

// TopicTenReconciler reconciles a TopicTen object
type TopicTenReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=app.lyc-handson-aws.com,resources=topictens,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=app.lyc-handson-aws.com,resources=topictens/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=app.lyc-handson-aws.com,resources=topictens/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the TopicTen object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.3/pkg/reconcile
func (r *TopicTenReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	log := log.FromContext(ctx)
	log.Info("starting reconciliation - TopicTen")

	topicTen := &appv1alpha1.TopicTen{}

	// Fetch the TopicTen instance
	if err := r.Get(ctx, req.NamespacedName, topicTen); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Try to fetch the existing Deployment
	foundDeployment := &appsv1.Deployment{}
	err := r.Get(ctx, req.NamespacedName, foundDeployment)
	if err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("Creating Deployment", "Namespace", req.Namespace, "Name", req.Name)
			foundDeployment := r.createDeployment(topicTen)
			if err := r.Create(ctx, foundDeployment); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{RequeueAfter: 2 * time.Minute}, nil
		}
		return ctrl.Result{}, err
	}

	// Reduce replica count every 2 minutes if greater than 1
	if *foundDeployment.Spec.Replicas > 1 {
		*foundDeployment.Spec.Replicas--
		if err := r.Update(ctx, foundDeployment); err != nil {
			return ctrl.Result{}, err
		}
		log.Info(fmt.Sprintf("Reduced replicas to %d", *foundDeployment.Spec.Replicas))
		return ctrl.Result{RequeueAfter: 2 * time.Minute}, nil
	}

	// Update status with current replica count
	topicTen.Status.CurrentReplicas = *foundDeployment.Spec.Replicas
	if err := r.Status().Update(ctx, topicTen); err != nil {
		log.Error(err, "Failed to update MyPod status")
		return ctrl.Result{}, err
	}

	log.Info("end of reconciliation - TopicTen")
	return ctrl.Result{}, nil
}

func (r *TopicTenReconciler) createDeployment(topicTen *appv1alpha1.TopicTen) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      topicTen.Name,
			Namespace: topicTen.Namespace,
			Labels:    map[string]string{"app": topicTen.Name},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &topicTen.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": topicTen.Name},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": topicTen.Name},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "topicten-container",
							Image: topicTen.Spec.Image,
							Env: []corev1.EnvVar{
								{
									Name:  "STORAGE_ARN",
									Value: topicTen.Spec.TargetArn,
								},
								{
									Name:  "AWS_KMS_ARN",
									Value: topicTen.Spec.KMSArn,
								},
								{
									Name:  "AWS_LOG_GROUP_ARN",
									Value: topicTen.Spec.CloudWatchArn,
								},
								{
									Name: "MY_POD_NAME",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "metadata.name",
										},
									},
								},
								{
									Name: "MY_POD_NAMESPACE",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "metadata.namespace",
										},
									},
								},
								{
									Name: "MY_POD_IP",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "status.podIP",
										},
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

// SetupWithManager sets up the controller with the Manager.
func (r *TopicTenReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appv1alpha1.TopicTen{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
