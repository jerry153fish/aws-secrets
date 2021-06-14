/*
Copyright 2021.

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

package controllers

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/go-logr/logr"
	cfnv1alpha1 "github.com/jerry153fish/cloudformation-secrets/api/v1alpha1"
	utils "github.com/jerry153fish/cloudformation-secrets/utils"
	"github.com/patrickmn/go-cache"
)

var (
	cf *cloudformation.CloudFormation
	c  *cache.Cache = cache.New(5*time.Minute, 10*time.Minute)
)

// SecretsReconciler reconciles a Secrets object
type SecretsReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

//+kubebuilder:rbac:groups=cfn.jerry153fish.com,resources=secrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cfn.jerry153fish.com,resources=secrets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cfn.jerry153fish.com,resources=secrets/finalizers,verbs=update
//+kubebuilder:rbac:groups=,resources=secrets,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Secrets object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *SecretsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	slog := log.Log.WithValues("cfnSecrets", req.NamespacedName)

	slog.Info("Printing at INFO level")

	cf = utils.GetCfnClient()

	cfnSecret := &cfnv1alpha1.Secrets{}

	if err := r.Get(ctx, req.NamespacedName, cfnSecret); err != nil {
		slog.Error(err, "unable to fetch cfnSecrets")
		// TODO: we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	// your logic here

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SecretsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cfnv1alpha1.Secrets{}).
		Complete(r)
}
