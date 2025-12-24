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

package v1

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	appsv1 "github.com/shieldx-bot/replica-master/api/v1"
)

// nolint:unused
// log is for logging in this package.
var replicamasterlog = logf.Log.WithName("replicamaster-resource")

// SetupReplicaMasterWebhookWithManager registers the webhook for ReplicaMaster in the manager.
func SetupReplicaMasterWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(&appsv1.ReplicaMaster{}).
		WithValidator(&ReplicaMasterCustomValidator{}).
		WithDefaulter(&ReplicaMasterCustomDefaulter{}).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-apps-shieldx-local-v1-replicamaster,mutating=true,failurePolicy=fail,sideEffects=None,groups=apps.shieldx.local,resources=replicamasters,verbs=create;update;delete,versions=v1,name=mreplicamaster-v1.kb.io,admissionReviewVersions=v1

// ReplicaMasterCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind ReplicaMaster when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type ReplicaMasterCustomDefaulter struct {
	// TODO(user): Add more fields as needed for defaulting
	client.Client
	Scheme *runtime.Scheme
}

var _ webhook.CustomDefaulter = &ReplicaMasterCustomDefaulter{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind ReplicaMaster.
func (d *ReplicaMasterCustomDefaulter) Default(_ context.Context, obj runtime.Object) error {
	replicamaster, ok := obj.(*appsv1.ReplicaMaster)

	if !ok {
		return fmt.Errorf("expected an ReplicaMaster object but got %T", obj)
	}
	replicamasterlog.Info("Defaulting for ReplicaMaster", "name", replicamaster.GetName())

	// TODO(user): fill in your defaulting logic.

	return nil
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: If you want to customise the 'path', use the flags '--defaulting-path' or '--validation-path'.
// +kubebuilder:webhook:path=/validate-apps-shieldx-local-v1-replicamaster,mutating=false,failurePolicy=fail,sideEffects=None,groups=apps.shieldx.local,resources=replicamasters,verbs=create;update;delete,versions=v1,name=vreplicamaster-v1.kb.io,admissionReviewVersions=v1

// ReplicaMasterCustomValidator struct is responsible for validating the ReplicaMaster resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type ReplicaMasterCustomValidator struct {
	// TODO(user): Add more fields as needed for validation
}

var _ webhook.CustomValidator = &ReplicaMasterCustomValidator{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type ReplicaMaster.
func (v *ReplicaMasterCustomValidator) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	replicamaster, ok := obj.(*appsv1.ReplicaMaster)
	if !ok {
		return nil, fmt.Errorf("expected a ReplicaMaster object but got %T", obj)
	}
	replicamasterlog.Info("Validation for ReplicaMaster upon creation", "name", replicamaster.GetName())

	// TODO(user): fill in your validation logic upon object creation.

	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type ReplicaMaster.
func (v *ReplicaMasterCustomValidator) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	replicamaster, ok := newObj.(*appsv1.ReplicaMaster)
	if !ok {
		return nil, fmt.Errorf("expected a ReplicaMaster object for the newObj but got %T", newObj)
	}
	replicamasterlog.Info("Validation for ReplicaMaster upon update", "name", replicamaster.GetName())

	// TODO(user): fill in your validation logic upon object update.

	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type ReplicaMaster.
func (v *ReplicaMasterCustomValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	replicamaster, ok := obj.(*appsv1.ReplicaMaster)
	if !ok {
		return nil, fmt.Errorf("expected a ReplicaMaster object but got %T", obj)
	}
	replicamasterlog.Info("Validation for ReplicaMaster upon deletion", "name", replicamaster.GetName())

	policy := replicamaster.GetLabels()["policy"]
	// Treat these policies as non-deletable.
	// NOTE: `critical-app.yaml` uses policy=critical.
	switch policy {
	case "protected", "critical":
		return nil, fmt.Errorf("không cho phép xóa: ReplicaMaster %s đang được bảo vệ bởi policy=%q", replicamaster.GetName(), policy)
	}

	return nil, nil
}
