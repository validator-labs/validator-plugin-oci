package controller

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	vapi "github.com/validator-labs/validator/api/v1alpha1"
	"github.com/validator-labs/validator/pkg/util"
	vres "github.com/validator-labs/validator/pkg/validationresult"

	"github.com/validator-labs/validator-plugin-oci/api/v1alpha1"
	// +kubebuilder:scaffold:imports
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
					RuleName:       "basic auth and empty artifact list",
					Host:           "foo1.registry.io",
					ValidationType: v1alpha1.ValidationTypeNone,
					Auth: v1alpha1.Auth{
						Basic: &v1alpha1.BasicAuth{
							Username: "userName",
							Password: "pa$$w0rd",
						},
					},
					Artifacts: []v1alpha1.Artifact{},
					SignatureVerification: v1alpha1.SignatureVerification{
						Provider: "cosign",
					},
				},
				{
					RuleName:       "full layer validation enabled on rule level",
					Host:           "foo.registry.io",
					ValidationType: v1alpha1.ValidationTypeFull,
					Artifacts: []v1alpha1.Artifact{
						{
							Ref: "foo/bar:latest",
						},
					},
					SignatureVerification: v1alpha1.SignatureVerification{
						Provider: "cosign",
					},
				},
				{
					RuleName:       "fast layer validation enabled on artifact level",
					Host:           "foo.registry.io",
					ValidationType: v1alpha1.ValidationTypeNone,
					Artifacts: []v1alpha1.Artifact{
						{
							Ref:            "foo/bar:latest",
							ValidationType: util.Ptr(v1alpha1.ValidationTypeFast),
						},
					},
					SignatureVerification: v1alpha1.SignatureVerification{
						Provider: "cosign",
					},
				},
				{
					RuleName:       "secret auth and ca cert provided",
					Host:           "foo2.registry.io",
					ValidationType: v1alpha1.ValidationTypeNone,
					Auth: v1alpha1.Auth{
						SecretName: util.Ptr("auth-secret"),
					},
					CaCert: "dummy-ca-cert",
					Artifacts: []v1alpha1.Artifact{
						{
							Ref:            "foo/bar:latest",
							ValidationType: util.Ptr(v1alpha1.ValidationTypeFast),
						},
					},
					SignatureVerification: v1alpha1.SignatureVerification{
						Provider: "cosign",
					},
				},
				{
					RuleName:       "ecr auth and pubkeys secret provided and created",
					Host:           "foo3.registry.io",
					ValidationType: v1alpha1.ValidationTypeNone,
					Auth: v1alpha1.Auth{
						ECR: &v1alpha1.ECRAuth{
							AccessKeyID:     "accessKeyID",
							SecretAccessKey: "secretAccessKey",
							SessionToken:    "sessionToken",
						},
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
