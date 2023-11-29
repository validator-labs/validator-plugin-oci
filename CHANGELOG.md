# Changelog

## [0.0.2](https://github.com/spectrocloud-labs/validator-plugin-oci/compare/v0.0.1...v0.0.2) (2023-11-29)


### Features

* add Helm chart ([#25](https://github.com/spectrocloud-labs/validator-plugin-oci/issues/25)) ([f4295ae](https://github.com/spectrocloud-labs/validator-plugin-oci/commit/f4295ae9a509c52763c12ba01458d8d0150b0bae))
* allow initContainer image to be passed in via values.yaml ([#27](https://github.com/spectrocloud-labs/validator-plugin-oci/issues/27)) ([50c8647](https://github.com/spectrocloud-labs/validator-plugin-oci/commit/50c8647f76cc70453b1ec1a5f7e307fcda839235))
* implement OCI registry validation spec ([#6](https://github.com/spectrocloud-labs/validator-plugin-oci/issues/6)) ([f62c494](https://github.com/spectrocloud-labs/validator-plugin-oci/commit/f62c494d3a44bcf99c9d0bccecd1af2b8bc3ae78))
* support validating list of oci artifacts ([#16](https://github.com/spectrocloud-labs/validator-plugin-oci/issues/16)) ([d0cbecc](https://github.com/spectrocloud-labs/validator-plugin-oci/commit/d0cbecc24614a9a6ddf2a34e71e01ce23a313d8c))


### Bug Fixes

* **deps:** update kubernetes packages to v0.28.4 ([#17](https://github.com/spectrocloud-labs/validator-plugin-oci/issues/17)) ([f346f63](https://github.com/spectrocloud-labs/validator-plugin-oci/commit/f346f631c50d2fdc6236603055c792f116c554df))
* **deps:** update module github.com/aws/aws-sdk-go-v2/config to v1.25.6 ([#28](https://github.com/spectrocloud-labs/validator-plugin-oci/issues/28)) ([b12dabe](https://github.com/spectrocloud-labs/validator-plugin-oci/commit/b12dabe9730e9e12a48e979f796fde71dbd551a0))
* **deps:** update module github.com/aws/aws-sdk-go-v2/config to v1.25.8 ([#32](https://github.com/spectrocloud-labs/validator-plugin-oci/issues/32)) ([3eb0824](https://github.com/spectrocloud-labs/validator-plugin-oci/commit/3eb08241bd645d73cd50182fd562f941171b4a30))
* **deps:** update module github.com/aws/aws-sdk-go-v2/service/ecr to v1.23.2 ([#30](https://github.com/spectrocloud-labs/validator-plugin-oci/issues/30)) ([6375b2b](https://github.com/spectrocloud-labs/validator-plugin-oci/commit/6375b2bafbdcaa8649691eaba15ba52ec8eb80d9))
* **deps:** update module github.com/aws/aws-sdk-go-v2/service/ecr to v1.23.3 ([#34](https://github.com/spectrocloud-labs/validator-plugin-oci/issues/34)) ([a55ada3](https://github.com/spectrocloud-labs/validator-plugin-oci/commit/a55ada393e0ef05510388a152b4e0a03a573d3d4))
* **deps:** update module github.com/onsi/ginkgo/v2 to v2.13.1 ([#15](https://github.com/spectrocloud-labs/validator-plugin-oci/issues/15)) ([23673ac](https://github.com/spectrocloud-labs/validator-plugin-oci/commit/23673ac0092fac7eecc78f7d92c249b385537c39))
* **deps:** update module github.com/onsi/ginkgo/v2 to v2.13.2 ([#35](https://github.com/spectrocloud-labs/validator-plugin-oci/issues/35)) ([4f44a26](https://github.com/spectrocloud-labs/validator-plugin-oci/commit/4f44a26a67d141a8b953cd937d07c1c0482087eb))
* **deps:** update module github.com/spectrocloud-labs/validator to v0.0.18 ([#14](https://github.com/spectrocloud-labs/validator-plugin-oci/issues/14)) ([58c78f4](https://github.com/spectrocloud-labs/validator-plugin-oci/commit/58c78f43a7f21d8d22381042ea73fbe5f3b7f0d0))
* **deps:** update module github.com/spectrocloud-labs/validator to v0.0.21 ([#18](https://github.com/spectrocloud-labs/validator-plugin-oci/issues/18)) ([9373c1d](https://github.com/spectrocloud-labs/validator-plugin-oci/commit/9373c1d3541397948eca4a93df61c3a628661b56))
* **deps:** update module github.com/spectrocloud-labs/validator to v0.0.25 ([#21](https://github.com/spectrocloud-labs/validator-plugin-oci/issues/21)) ([76f1b24](https://github.com/spectrocloud-labs/validator-plugin-oci/commit/76f1b247a7bf69d990a0539e8ec73260cfe7ad5a))
* set owner references on validation result to ensure cleanup ([#19](https://github.com/spectrocloud-labs/validator-plugin-oci/issues/19)) ([9c7c28d](https://github.com/spectrocloud-labs/validator-plugin-oci/commit/9c7c28d1e69b9488263537e48415818826d96ebf))


### Other

* add license badge ([1eb5f1b](https://github.com/spectrocloud-labs/validator-plugin-oci/commit/1eb5f1b2ceafc7656816f42b4f51c11ad0057aba))
* **deps:** pin dependencies ([#9](https://github.com/spectrocloud-labs/validator-plugin-oci/issues/9)) ([9876cd7](https://github.com/spectrocloud-labs/validator-plugin-oci/commit/9876cd701be178016231d02661a78db1f2f48c85))
* **deps:** update actions/checkout action to v4 ([#10](https://github.com/spectrocloud-labs/validator-plugin-oci/issues/10)) ([cd110af](https://github.com/spectrocloud-labs/validator-plugin-oci/commit/cd110af99d4eed651d89dabe5565bcedcb3f4c35))
* **deps:** update anchore/sbom-action action to v0.15.0 ([#23](https://github.com/spectrocloud-labs/validator-plugin-oci/issues/23)) ([34253f0](https://github.com/spectrocloud-labs/validator-plugin-oci/commit/34253f03e491ebecc0ce8631d56558cc16bb4b82))
* **deps:** update docker/build-push-action digest to 4a13e50 ([#20](https://github.com/spectrocloud-labs/validator-plugin-oci/issues/20)) ([eace63e](https://github.com/spectrocloud-labs/validator-plugin-oci/commit/eace63e7d49fc14c8d1f8d0427bd11039bef140d))
* fix platform specification for manager image ([#13](https://github.com/spectrocloud-labs/validator-plugin-oci/issues/13)) ([539e8be](https://github.com/spectrocloud-labs/validator-plugin-oci/commit/539e8be372a623125d1ed04e602833c59acddd93))
* specify platform in Dockerfile and docker-build make target ([#12](https://github.com/spectrocloud-labs/validator-plugin-oci/issues/12)) ([a88d182](https://github.com/spectrocloud-labs/validator-plugin-oci/commit/a88d1820503bdfc6c2f99690db2a0bcd6befc5dc))
* switch back to public bulwark images ([010e7f8](https://github.com/spectrocloud-labs/validator-plugin-oci/commit/010e7f842a54cb0f9e0f572618007ad85009f766))
* update spectrocloud-labs/validator dependency to v0.0.15 ([f62c494](https://github.com/spectrocloud-labs/validator-plugin-oci/commit/f62c494d3a44bcf99c9d0bccecd1af2b8bc3ae78))


### Refactoring

* switch from oras to go-containerregistry ([#24](https://github.com/spectrocloud-labs/validator-plugin-oci/issues/24)) ([eef0013](https://github.com/spectrocloud-labs/validator-plugin-oci/commit/eef0013a7d1072f55bb3356304f287ad1cc61ff4))
* switch init container to image with ca-certificates pre installed ([#33](https://github.com/spectrocloud-labs/validator-plugin-oci/issues/33)) ([4550f4b](https://github.com/spectrocloud-labs/validator-plugin-oci/commit/4550f4bedb9807d8578fcc56d7fc4e3309cd6d8b))
