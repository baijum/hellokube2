package controllers

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	apixv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Guestbook Controller:", func() {
	const (
		timeout       = time.Second * 20
		interval      = time.Millisecond * 250
		testNamespace = "default"
	)

	Context("When creating Guestbook", func() {

		It("should be possible to create custom resource", func() {
			ctx := context.Background()

			backingServiceCRD := &apixv1.CustomResourceDefinition{
				ObjectMeta: metav1.ObjectMeta{
					Name: "backingservices.webapp.baiju.dev",
				},
				Spec: apixv1.CustomResourceDefinitionSpec{
					Group: "webapp.baiju.dev",
					Versions: []apixv1.CustomResourceDefinitionVersion{{
						Name:    "v1alpha1",
						Served:  true,
						Storage: true,
						Schema: &apixv1.CustomResourceValidation{
							OpenAPIV3Schema: &apixv1.JSONSchemaProps{
								Type: "object",
								Properties: map[string]apixv1.JSONSchemaProps{
									"status": {
										Type: "object",
										Properties: map[string]apixv1.JSONSchemaProps{
											"binding": {
												Type: "object",
												Properties: map[string]apixv1.JSONSchemaProps{
													"name": {
														Type: "string",
													},
												},
												Required: []string{"name"},
											},
										},
									},
								},
							},
						},
					},
					},
					Names: apixv1.CustomResourceDefinitionNames{
						Plural:   "backingservices",
						Singular: "backingservice",
						Kind:     "BackingService",
					},
					Scope: apixv1.ClusterScoped,
				}}
			Expect(k8sClient.Create(ctx, backingServiceCRD)).Should(Succeed())

			backingServiceCRDLookupKey := types.NamespacedName{Name: "backingservices.webapp.baiju.dev", Namespace: testNamespace}
			createdBackingServiceCRD := &apixv1.CustomResourceDefinition{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, backingServiceCRDLookupKey, createdBackingServiceCRD)
				if err != nil {
					return false
				}
				for _, condition := range createdBackingServiceCRD.Status.Conditions {
					if condition.Type == apixv1.Established &&
						condition.Status == apixv1.ConditionTrue {
						return true
					}
				}
				return false
			}, timeout, interval).Should(BeTrue())

			backingServiceCR := &unstructured.Unstructured{
				Object: map[string]interface{}{
					"kind":       "BackingService",
					"apiVersion": "webapp.baiju.dev/v1alpha1",
					"metadata": map[string]interface{}{
						"name": "back1",
					},
					"status": map[string]interface{}{
						"binding": map[string]interface{}{
							"name": "secret1",
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, backingServiceCR)).Should(Succeed())

		})
	})

})
