/*
Copyright 2025.

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
	"time"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	replicav1 "github.com/shieldx-bot/replica-master/api/v1"
)

// ReplicaMasterReconciler reconciles a ReplicaMaster object
type ReplicaMasterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=apps.shieldx.local,resources=replicamasters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps.shieldx.local,resources=replicamasters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps.shieldx.local,resources=replicamasters/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;update;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ReplicaMaster object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.22.4/pkg/reconcile
func (r *ReplicaMasterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	// 1. Lấy thông tin ReplicaMaster (Cái tờ lệnh của Sếp)
	var replicaMaster replicav1.ReplicaMaster
	if err := r.Get(ctx, req.NamespacedName, &replicaMaster); err != nil {
		// Nếu không tìm thấy (do bị xóa), thì thôi không làm gì cả
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// 2. Tìm Deployment tương ứng (Nhân viên cần quản lý)
	var deployment appsv1.Deployment
	deploymentName := replicaMaster.Spec.DeploymentName
	// Giả sử Deployment nằm cùng namespace với CRD
	if err := r.Get(ctx, client.ObjectKey{Namespace: req.Namespace, Name: deploymentName}, &deployment); err != nil {
		l.Error(err, "❌ Cannot find Deployment", "Deployment.Name", deploymentName)
		return ctrl.Result{}, err // Thử lại sau (Re-queue)
	}

	// 3. SO SÁNH: Thực tế (Deployment) vs Mong muốn (ReplicaMaster)
	desiredReplicas := replicaMaster.Spec.Replicas
	currentReplicas := *deployment.Spec.Replicas

	if currentReplicas != desiredReplicas {
		l.Info("⚠️  Detected drift!", "Current", currentReplicas, "Desired", desiredReplicas)

		// 4. HÀNH ĐỘNG: Cập nhật lại Deployment
		deployment.Spec.Replicas = &desiredReplicas
		if err := r.Update(ctx, &deployment); err != nil {
			l.Error(err, "❌ Failed to update Deployment")
			return ctrl.Result{}, err
		}
		l.Info("✅ Fix applied! Deployment scaled successfully.")
	}

	return ctrl.Result{RequeueAfter: 2 * time.Second}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ReplicaMasterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&replicav1.ReplicaMaster{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
