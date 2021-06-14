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
	"reflect"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/go-logr/logr"
	cfnv1alpha1 "github.com/jerry153fish/cloudformation-secrets/api/v1alpha1"
	utils "github.com/jerry153fish/cloudformation-secrets/utils"
	"github.com/patrickmn/go-cache"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	log := r.Log.WithValues("cfnSecrets", req.NamespacedName)

	cf = utils.GetCfnClient()

	secrets := &cfnv1alpha1.Secrets{}

	if err := r.Get(ctx, req.NamespacedName, secrets); err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("Secrets resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get Secrets")
		return ctrl.Result{}, err
	}
	// TODO: validate here
	// Check if the Secret already exists, if not create a new one
	found := &corev1.Secret{}
	err := r.Get(ctx, types.NamespacedName{Name: secrets.Name, Namespace: secrets.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		// Define a new Secret
		sec := r.SecretsCr2Secret(secrets)
		log.Info("Creating a new k8s Secret", "Secret.Namespace", sec.Namespace, "Secret.Name", sec.Name)
		err = r.Create(ctx, sec)
		if err != nil {
			log.Error(err, "Failed to create new K8s Secret", "Secret.Namespace", sec.Namespace, "Secret.Name", sec.Name)
			return ctrl.Result{}, err
		}
		// Secret created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Secret")
		return ctrl.Result{}, err
	}

	// TODO: update
	if shouldUpdate(secrets, found) {
		log.Info("Updating a k8s Secret", "Secret.Namespace", found.Namespace, "Secret.Name", found.Name)
		found.Data = getSecretData(secrets)
		err = r.Update(ctx, found)
		if err != nil {
			log.Error(err, "Failed to update K8s Secret", "Secret.Namespace", found.Namespace, "Secret.Name", found.Name)
			return ctrl.Result{}, err
		}
		// Secret created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SecretsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cfnv1alpha1.Secrets{}).
		Complete(r)
}

func (r *SecretsReconciler) SecretsCr2Secret(secrets *cfnv1alpha1.Secrets) *corev1.Secret {
	// TODO: verify
	// TODO: more metadata eg labels

	sec := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secrets.Name,
			Namespace: secrets.Namespace,
		},
		Data: getSecretData(secrets),
	}

	// Set Secrets CR as the owner and controller
	ctrl.SetControllerReference(secrets, sec, r.Scheme)
	return sec
}

func getSecretData(secrets *cfnv1alpha1.Secrets) map[string][]byte {
	// cfn := secrets.Spec.Cfn
	plainCreds := secrets.Spec.PlainCreds

	re := make(map[string][]byte)

	// if cfn != nil {

	// }

	for _, cred := range plainCreds {
		re[cred.KeyName] = []byte(cred.Value)
	}

	if len(re) > 0 {
		return re
	}

	return nil
}

func shouldUpdate(secrets *cfnv1alpha1.Secrets, k8sSec *corev1.Secret) bool {
	k8sSecData := k8sSec.Data

	secretsData := getSecretData(secrets)

	return !reflect.DeepEqual(secretsData, k8sSecData)
}
