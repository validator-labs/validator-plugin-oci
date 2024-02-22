package controller

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/spectrocloud-labs/validator-plugin-oci/api/v1alpha1"
	vapi "github.com/spectrocloud-labs/validator/api/v1alpha1"
	//+kubebuilder:scaffold:imports
)

const ociValidatorName = "oci-validator"

var _ = Describe("OCIValidator controller", Ordered, func() {

	BeforeEach(func() {
		// toggle true/false to enable/disable the OCIValidator controller specs
		if false {
			Skip("skipping")
		}
	})

	val := &v1alpha1.OciValidator{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ociValidatorName,
			Namespace: validatorNamespace,
		},
		Spec: v1alpha1.OciValidatorSpec{
			OciRegistryRules: []v1alpha1.OciRegistryRule{
				{
					Host:      "foo1.registry.io",
					Artifacts: []v1alpha1.Artifact{},
				},
				{
					Host: "foo.registry.io",
					Artifacts: []v1alpha1.Artifact{
						{
							Ref:             "foo/bar:latest",
							LayerValidation: true,
						},
					},
				},
				{
					Host:   "foo2.registry.io",
					CaCert: "dummy-ca-cert",
					Artifacts: []v1alpha1.Artifact{
						{
							Ref:             "foo/bar:latest",
							LayerValidation: true,
						},
					},
				},
				{
					Host: "foo3.registry.io",
					Auth: v1alpha1.Auth{
						SecretName: "mySecret",
					},
					Artifacts: []v1alpha1.Artifact{
						{
							Ref: "foo/bar:latest",
						},
					},
				},
			},
		},
	}

	vr := &vapi.ValidationResult{}
	vrKey := types.NamespacedName{Name: validationResultName(val), Namespace: validatorNamespace}

	It("Should create a ValidationResult and update its Status with a failed condition", func() {
		By("By creating a new OCIValidator")
		ctx := context.Background()

		Expect(k8sClient.Create(ctx, val)).Should(Succeed())

		// Wait for the ValidationResult's Status to be updated
		Eventually(func() bool {
			if err := k8sClient.Get(ctx, vrKey, vr); err != nil {
				return false
			}

			stateOk := vr.Status.State == vapi.ValidationFailed
			return stateOk
		}, timeout, interval).Should(BeTrue(), "failed to create a ValidationResult")
	})
})
