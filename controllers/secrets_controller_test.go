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

	cfnv1alpha1 "github.com/jerry153fish/aws-secrets/api/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Secrets controller", func() {

	// Define utility constants for object names and testing timeouts/durations and intervals.
	const (
		SecretName      = "test-cfn-secrets"
		SecretNamespace = "default"

		timeout  = time.Second * 10
		duration = time.Second * 10
		interval = time.Millisecond * 250
	)

	Context("When updating Secrets Status", func() {
		It("Should increase Secrets Status.Active count when new Jobs are created", func() {
			By("By creating a new Secrets")
			ctx := context.Background()
			Secrets := &cfnv1alpha1.Secrets{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "cfn.jerry153fish.com/v1alpha1",
					Kind:       "Secrets",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      SecretName,
					Namespace: SecretNamespace,
				},
				Spec: cfnv1alpha1.SecretsSpec{
					SecretName: SecretName,
					PlainCreds: []cfnv1alpha1.PlainCred{
						cfnv1alpha1.PlainCred{
							KeyName: "testPlain1",
							Value:   "12345",
						},
						cfnv1alpha1.PlainCred{
							KeyName: "testPlain2",
							Value:   "12345",
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, Secrets)).Should(Succeed())

			/*
				After creating this CR Secrets, let's check that the Secrets's Spec fields match what we passed in.
			*/

			secretsLookupKey := types.NamespacedName{Name: SecretName, Namespace: SecretNamespace}
			createdSecrets := &cfnv1alpha1.Secrets{}

			// We'll need to retry getting this newly created CronJob, given that creation may not immediately happen.
			Eventually(func() bool {
				err := k8sClient.Get(ctx, secretsLookupKey, createdSecrets)
				return err == nil
			}, timeout, interval).Should(BeTrue())

			Expect(createdSecrets.Spec.SecretName).Should(Equal(SecretName))
			Expect(createdSecrets.Spec.PlainCreds).Should(HaveLen(2))
			Expect(createdSecrets.Spec.PlainCreds[0].Value).Should(Equal("12345"))
			Expect(createdSecrets.Spec.Cfn).To(BeNil())

			/*
				Next we need to check the created k8s secret
			*/
			createdSec := &corev1.Secret{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, secretsLookupKey, createdSec)
				return err == nil
			}, timeout, interval).Should(BeTrue())

			Expect(createdSec.ObjectMeta.Name).To(Equal(SecretName))
			Expect(createdSec.ObjectMeta.Namespace).To(Equal(SecretNamespace))
			Expect(createdSec.Data).NotTo(BeNil())
			Expect(createdSec.Data["testPlain1"]).To(Equal([]byte("12345")))
			// end to It
		})
	})

})
