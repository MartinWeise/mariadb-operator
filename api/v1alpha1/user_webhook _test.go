/*
Copyright 2022.

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

package v1alpha1

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("User webhook", Ordered, func() {
	Context("When updating a User", func() {
		key := types.NamespacedName{
			Name:      "user-create-webhook",
			Namespace: testNamespace,
		}
		BeforeAll(func() {
			user := User{
				ObjectMeta: metav1.ObjectMeta{
					Name:      key.Name,
					Namespace: key.Namespace,
				},
				Spec: UserSpec{
					MariaDBRef: MariaDBRef{
						ObjectReference: corev1.ObjectReference{
							Name: "mariadb-webhook",
						},
						WaitForIt: true,
					},
					PasswordSecretKeyRef: corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: "user-mariadb-webhook-root",
						},
						Key: "password",
					},
					MaxUserConnections: 10,
				},
			}
			Expect(k8sClient.Create(testCtx, &user)).To(Succeed())
		})
		DescribeTable(
			"Should validate",
			func(patchFn func(u *User), wantErr bool) {
				var user User
				Expect(k8sClient.Get(testCtx, key, &user)).To(Succeed())

				patch := client.MergeFrom(user.DeepCopy())
				patchFn(&user)

				err := k8sClient.Patch(testCtx, &user, patch)
				if wantErr {
					Expect(err).To(HaveOccurred())
				} else {
					Expect(err).ToNot(HaveOccurred())
				}
			},
			Entry(
				"Updating MariaDBRef",
				func(umdb *User) {
					umdb.Spec.MariaDBRef.Name = "another-mariadb"
				},
				true,
			),
			Entry(
				"Updating PasswordSecretKeyRef",
				func(umdb *User) {
					umdb.Spec.PasswordSecretKeyRef.Name = "another-secret"
				},
				true,
			),
			Entry(
				"Updating MaxUserConnections",
				func(umdb *User) {
					umdb.Spec.MaxUserConnections = 20
				},
				true,
			),
		)
	})
})
