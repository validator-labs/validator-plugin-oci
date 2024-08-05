package controller

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/validator-labs/validator-plugin-oci/api/v1alpha1"
	vapi "github.com/validator-labs/validator/api/v1alpha1"
	vres "github.com/validator-labs/validator/pkg/validationresult"
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

	ociValidator := &v1alpha1.OciValidator{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ociValidatorName,
			Namespace: validatorNamespace,
		},
		Spec: v1alpha1.OciValidatorSpec{
			OciRegistryRules: []v1alpha1.OciRegistryRule{
				{
					RuleName:  "empty artifact list",
					Host:      "foo1.registry.io",
					Artifacts: []v1alpha1.Artifact{},
					SignatureVerification: v1alpha1.SignatureVerification{
						Provider: "cosign",
					},
				},
				{
					RuleName: "layer validation enabled",
					Host:     "foo.registry.io",
					Artifacts: []v1alpha1.Artifact{
						{
							Ref:             "foo/bar:latest",
							LayerValidation: true,
						},
					},
					SignatureVerification: v1alpha1.SignatureVerification{
						Provider: "cosign",
					},
				},
				{
					RuleName: "ca cert provided",
					Host:     "foo2.registry.io",
					CaCert:   "dummy-ca-cert",
					Artifacts: []v1alpha1.Artifact{
						{
							Ref:             "foo/bar:latest",
							LayerValidation: true,
						},
					},
					SignatureVerification: v1alpha1.SignatureVerification{
						Provider: "cosign",
					},
				},
				{
					RuleName: "auth secret and pubkeys secret provided and created",
					Host:     "foo3.registry.io",
					Auth: v1alpha1.Auth{
						SecretName: "auth-secret",
					},
					Artifacts: []v1alpha1.Artifact{
						{
							Ref: "foo/bar:latest",
						},
					},
					SignatureVerification: v1alpha1.SignatureVerification{
						Provider:   "cosign",
						SecretName: "pubkeys-secret",
					},
				},
			},
		},
	}

	vr := &vapi.ValidationResult{}
	vrKey := types.NamespacedName{Name: vres.Name(ociValidator), Namespace: validatorNamespace}

	authSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "auth-secret",
			Namespace: validatorNamespace,
		},
		StringData: map[string]string{
			"username": "userName",
			"password": "pa$$w0rd",
		},
		Type: corev1.SecretTypeOpaque,
	}

	pubKeysSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pubkeys-secret",
			Namespace: validatorNamespace,
		},
		Data: map[string][]byte{
			"key1.pub": []byte("LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUZrd0V3WUhLb1pJemowQ0FRWUlLb1pJemowREFRY0RRZ0FFS1B1Q285QW1KQ3BxR1doZWZqYmhrRmNyMUdBMwppTmE3NjVzZUUzallDM01HVWU1aDUyMzkzRGh5N0I1YlhHc2c2RWZQcE5ZYW1sQUVXanhDcEhGM0xnPT0KLS0tLS1FTkQgUFVCTElDIEtFWS0tLS0tCg=="),
		},
		Type: corev1.SecretTypeOpaque,
	}

	It("Should create a ValidationResult and update its Status with a failed condition", func() {
		By("By creating a new OCIValidator")
		ctx := context.Background()

		Expect(k8sClient.Create(ctx, authSecret)).Should(Succeed())
		Expect(k8sClient.Create(ctx, pubKeysSecret)).Should(Succeed())
		Expect(k8sClient.Create(ctx, ociValidator)).Should(Succeed())

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
