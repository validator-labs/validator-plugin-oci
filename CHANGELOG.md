# Changelog

## [0.0.13](https://github.com/validator-labs/validator-plugin-oci/compare/v0.0.12...v0.0.13) (2024-08-05)


### Other

* add hook to install validator crds in devspace ([#237](https://github.com/validator-labs/validator-plugin-oci/issues/237)) ([4c51037](https://github.com/validator-labs/validator-plugin-oci/commit/4c51037c7527b5b207e7caedada476f775c2f359))


### Dependency Updates

* **deps:** update github.com/awslabs/amazon-ecr-credential-helper/ecr-login digest to a8d7d3c ([#233](https://github.com/validator-labs/validator-plugin-oci/issues/233)) ([c422454](https://github.com/validator-labs/validator-plugin-oci/commit/c422454310fab330b2ca024b4fc3e5139f0e2ea9))
* **deps:** update module github.com/onsi/gomega to v1.34.1 ([#232](https://github.com/validator-labs/validator-plugin-oci/issues/232)) ([253a725](https://github.com/validator-labs/validator-plugin-oci/commit/253a7250f84c74b5221fe19dfd8c138c185ecd6d))
* **deps:** update module github.com/validator-labs/validator to v0.0.50 ([#230](https://github.com/validator-labs/validator-plugin-oci/issues/230)) ([febe051](https://github.com/validator-labs/validator-plugin-oci/commit/febe051a38bfa162ee01cef58001da82fcc3a966))
* **deps:** update module github.com/validator-labs/validator to v0.0.51 ([#234](https://github.com/validator-labs/validator-plugin-oci/issues/234)) ([1f82f08](https://github.com/validator-labs/validator-plugin-oci/commit/1f82f08aabc579a303fa1224d7ef94dfe90a6934))
* **deps:** update module github.com/validator-labs/validator to v0.1.0 ([#236](https://github.com/validator-labs/validator-plugin-oci/issues/236)) ([4a997be](https://github.com/validator-labs/validator-plugin-oci/commit/4a997be722b6f8661be75fa109442accc0fdf573))


### Refactoring

* derive authn keychain from host ([#238](https://github.com/validator-labs/validator-plugin-oci/issues/238)) ([a9a3fda](https://github.com/validator-labs/validator-plugin-oci/commit/a9a3fda957ef2035a7d188a8474a333ece4e4979))
* support direct rule evaluation ([#240](https://github.com/validator-labs/validator-plugin-oci/issues/240)) ([9761384](https://github.com/validator-labs/validator-plugin-oci/commit/9761384be3bdef61060d28c1394bf8ab542ba69d))

## [0.0.12](https://github.com/validator-labs/validator-plugin-oci/compare/v0.0.11...v0.0.12) (2024-07-26)


### Dependency Updates

* **deps:** update github.com/validator-labs/validator digest to 81fd1cf ([#220](https://github.com/validator-labs/validator-plugin-oci/issues/220)) ([ab21a0c](https://github.com/validator-labs/validator-plugin-oci/commit/ab21a0cc9c0851321e8aaeed8b9730581c706c8e))
* **deps:** update module github.com/onsi/ginkgo/v2 to v2.19.1 ([#228](https://github.com/validator-labs/validator-plugin-oci/issues/228)) ([ee3e28b](https://github.com/validator-labs/validator-plugin-oci/commit/ee3e28b2e50e939480058101aed2590ecd5794e3))
* **deps:** update module github.com/onsi/gomega to v1.34.0 ([#227](https://github.com/validator-labs/validator-plugin-oci/issues/227)) ([90bb855](https://github.com/validator-labs/validator-plugin-oci/commit/90bb855ac2610db28544e90f4072e00d0051cb02))
* **deps:** update module github.com/sigstore/cosign/v2 to v2.3.0 ([#224](https://github.com/validator-labs/validator-plugin-oci/issues/224)) ([92ba260](https://github.com/validator-labs/validator-plugin-oci/commit/92ba260603b0f657e57460e83dc55783e4dc8c06))
* **deps:** update module github.com/validator-labs/validator to v0.0.47 ([#223](https://github.com/validator-labs/validator-plugin-oci/issues/223)) ([4e12c85](https://github.com/validator-labs/validator-plugin-oci/commit/4e12c85dfabada2c82649b651caf8d6ed2cf2628))
* **deps:** update module github.com/validator-labs/validator to v0.0.48 ([#225](https://github.com/validator-labs/validator-plugin-oci/issues/225)) ([0685ce0](https://github.com/validator-labs/validator-plugin-oci/commit/0685ce0bdc97cede2d4fc0ad5c7b72087c555639))
* **deps:** update module github.com/validator-labs/validator to v0.0.49 ([#226](https://github.com/validator-labs/validator-plugin-oci/issues/226)) ([d895c58](https://github.com/validator-labs/validator-plugin-oci/commit/d895c58c0fa7172cec2cc52894b81132b5e6a342))

## [0.0.11](https://github.com/validator-labs/validator-plugin-oci/compare/v0.0.10...v0.0.11) (2024-07-19)


### Features

* public OCI client with proxy from env support ([#216](https://github.com/validator-labs/validator-plugin-oci/issues/216)) ([a0ab6d6](https://github.com/validator-labs/validator-plugin-oci/commit/a0ab6d6003fdc93962d6f57d4e85fa439b977d11))


### Bug Fixes

* associate unexected errs w/ rules; always include validation result details ([#219](https://github.com/validator-labs/validator-plugin-oci/issues/219)) ([891a318](https://github.com/validator-labs/validator-plugin-oci/commit/891a318738ac5f552b2bf9325754add5e811a69a))
* **deps:** update aws-sdk-go-v2 monorepo ([#152](https://github.com/validator-labs/validator-plugin-oci/issues/152)) ([14d7e95](https://github.com/validator-labs/validator-plugin-oci/commit/14d7e959d24d1d98cc2f0feb0a202427f81d7ea6))
* **deps:** update aws-sdk-go-v2 monorepo ([#174](https://github.com/validator-labs/validator-plugin-oci/issues/174)) ([e166fff](https://github.com/validator-labs/validator-plugin-oci/commit/e166fffbe2ea5d0e6e3bd9a11a95108b916b75fa))
* **deps:** update aws-sdk-go-v2 monorepo ([#184](https://github.com/validator-labs/validator-plugin-oci/issues/184)) ([d4d6f31](https://github.com/validator-labs/validator-plugin-oci/commit/d4d6f319eb2f096bdd53b9e554ae2e1004c6c8ca))
* **deps:** update aws-sdk-go-v2 monorepo ([#191](https://github.com/validator-labs/validator-plugin-oci/issues/191)) ([b4941e9](https://github.com/validator-labs/validator-plugin-oci/commit/b4941e95d4fe9f69aa58a2879e5df2a8a54ef361))
* **deps:** update aws-sdk-go-v2 monorepo ([#195](https://github.com/validator-labs/validator-plugin-oci/issues/195)) ([b5e6e37](https://github.com/validator-labs/validator-plugin-oci/commit/b5e6e370249bb54edd5b2a88477e3882518fd119))
* **deps:** update aws-sdk-go-v2 monorepo ([#196](https://github.com/validator-labs/validator-plugin-oci/issues/196)) ([010cc4f](https://github.com/validator-labs/validator-plugin-oci/commit/010cc4f2283ae7477e60500cf8f02c8ebc711c7f))
* **deps:** update aws-sdk-go-v2 monorepo ([#197](https://github.com/validator-labs/validator-plugin-oci/issues/197)) ([5ce52c6](https://github.com/validator-labs/validator-plugin-oci/commit/5ce52c6c1436e13fa7019ea54cf42af5e3a1b833))
* **deps:** update aws-sdk-go-v2 monorepo ([#200](https://github.com/validator-labs/validator-plugin-oci/issues/200)) ([c482420](https://github.com/validator-labs/validator-plugin-oci/commit/c48242089d6c4d3cc8ac5d1b8696d7f055433049))
* **deps:** update kubernetes packages to v0.30.1 ([#165](https://github.com/validator-labs/validator-plugin-oci/issues/165)) ([d75bd41](https://github.com/validator-labs/validator-plugin-oci/commit/d75bd41663d70cfdf5a4709ff915dcf753eb75c4))
* **deps:** update kubernetes packages to v0.30.2 ([#193](https://github.com/validator-labs/validator-plugin-oci/issues/193)) ([9b10260](https://github.com/validator-labs/validator-plugin-oci/commit/9b10260e9dad9748ea18d4c845a60c938529ac23))
* **deps:** update module github.com/go-logr/logr to v1.4.2 ([#177](https://github.com/validator-labs/validator-plugin-oci/issues/177)) ([2ed9dba](https://github.com/validator-labs/validator-plugin-oci/commit/2ed9dba3ba66dd42a3ebdfa6b35e3b4aafd445ef))
* **deps:** update module github.com/google/go-containerregistry to v0.19.2 ([#194](https://github.com/validator-labs/validator-plugin-oci/issues/194)) ([65ecea1](https://github.com/validator-labs/validator-plugin-oci/commit/65ecea1c9998d62d18012e8997c9a734bd5ddc0c))
* **deps:** update module github.com/onsi/ginkgo/v2 to v2.19.0 ([#182](https://github.com/validator-labs/validator-plugin-oci/issues/182)) ([c0b10fc](https://github.com/validator-labs/validator-plugin-oci/commit/c0b10fce1e4557b1778fe4e68684e1aa0e3c3a06))
* **deps:** update module github.com/sigstore/cosign/v2 to v2.2.4 ([#162](https://github.com/validator-labs/validator-plugin-oci/issues/162)) ([2bf715f](https://github.com/validator-labs/validator-plugin-oci/commit/2bf715f0a0d70c2e6f46970dda3585124ac630d4))
* **deps:** update module github.com/sigstore/sigstore to v1.8.4 ([#178](https://github.com/validator-labs/validator-plugin-oci/issues/178)) ([67e2c8b](https://github.com/validator-labs/validator-plugin-oci/commit/67e2c8bef519c678f66f7f87a1456671b198ce69))
* **deps:** update module github.com/sigstore/sigstore to v1.8.5 ([#199](https://github.com/validator-labs/validator-plugin-oci/issues/199)) ([a454b94](https://github.com/validator-labs/validator-plugin-oci/commit/a454b94c5fffb71b7493f32aabfb333112409fa9))
* **deps:** update module github.com/validator-labs/validator to v0.0.41 ([#179](https://github.com/validator-labs/validator-plugin-oci/issues/179)) ([85f388f](https://github.com/validator-labs/validator-plugin-oci/commit/85f388f5c687e80e31d9e1c43ea0d933337e76a2))
* **deps:** update module github.com/validator-labs/validator to v0.0.42 ([#190](https://github.com/validator-labs/validator-plugin-oci/issues/190)) ([63f3dfd](https://github.com/validator-labs/validator-plugin-oci/commit/63f3dfd194b44c6bdd70632fd2aa1fb33b50ecac))
* **deps:** update module github.com/validator-labs/validator to v0.0.43 ([#198](https://github.com/validator-labs/validator-plugin-oci/issues/198)) ([3dc7de0](https://github.com/validator-labs/validator-plugin-oci/commit/3dc7de0decfdf2990dcdd2d18aff8c7b4301d671))
* **deps:** update module sigs.k8s.io/cluster-api to v1.7.2 ([#164](https://github.com/validator-labs/validator-plugin-oci/issues/164)) ([27a150c](https://github.com/validator-labs/validator-plugin-oci/commit/27a150c8e9d3fb9fe615531bb2e88590e846fbd6))
* **deps:** update module sigs.k8s.io/cluster-api to v1.7.3 ([#192](https://github.com/validator-labs/validator-plugin-oci/issues/192)) ([f9c2d5d](https://github.com/validator-labs/validator-plugin-oci/commit/f9c2d5dde679e42a2fbb2401c9b6f44ea2b391cc))
* **deps:** update module sigs.k8s.io/controller-runtime to v0.18.4 ([#188](https://github.com/validator-labs/validator-plugin-oci/issues/188)) ([8133d34](https://github.com/validator-labs/validator-plugin-oci/commit/8133d34578ffc841fc579db44e96f78746abeedd))


### Other

* **deps:** pin googleapis/release-please-action action to f3969c0 ([#171](https://github.com/validator-labs/validator-plugin-oci/issues/171)) ([8374a6f](https://github.com/validator-labs/validator-plugin-oci/commit/8374a6f7fb5201105646f72ab6c96b9ba929589c))
* **deps:** update actions/checkout digest to a5ac7e5 ([#172](https://github.com/validator-labs/validator-plugin-oci/issues/172)) ([c85b724](https://github.com/validator-labs/validator-plugin-oci/commit/c85b724e9ea5767191e6af656252b21190146bea))
* **deps:** update actions/setup-go digest to cdcb360 ([#175](https://github.com/validator-labs/validator-plugin-oci/issues/175)) ([133b586](https://github.com/validator-labs/validator-plugin-oci/commit/133b586e0cdb013fb887ebc412f8017af2f203b4))
* **deps:** update anchore/sbom-action action to v0.16.0 ([#180](https://github.com/validator-labs/validator-plugin-oci/issues/180)) ([7743e32](https://github.com/validator-labs/validator-plugin-oci/commit/7743e32713a6b4db09133f81e5624367bf56e823))
* **deps:** update azure/setup-helm digest to fe7b79c ([#163](https://github.com/validator-labs/validator-plugin-oci/issues/163)) ([5d50ba4](https://github.com/validator-labs/validator-plugin-oci/commit/5d50ba4391f58d354f9d895ad77a775034e0fa66))
* **deps:** update codecov/codecov-action digest to 125fc84 ([#173](https://github.com/validator-labs/validator-plugin-oci/issues/173)) ([0d023a6](https://github.com/validator-labs/validator-plugin-oci/commit/0d023a6fd103a1072124f2723bfc53616f6b50e5))
* **deps:** update codecov/codecov-action digest to 6d79887 ([#159](https://github.com/validator-labs/validator-plugin-oci/issues/159)) ([93abd02](https://github.com/validator-labs/validator-plugin-oci/commit/93abd02666d1e128e7a4af3cf7e2673b7af8a3c3))
* **deps:** update dependency go to v1.22.4 ([#185](https://github.com/validator-labs/validator-plugin-oci/issues/185)) ([e4288a5](https://github.com/validator-labs/validator-plugin-oci/commit/e4288a50f6216f8244f7a4636e025c9f80fc9ace))
* **deps:** update docker/login-action digest to 0d4c9c5 ([#176](https://github.com/validator-labs/validator-plugin-oci/issues/176)) ([2742f60](https://github.com/validator-labs/validator-plugin-oci/commit/2742f604cd1d6dc496419b9deb9c97a30b26b95f))
* **deps:** update docker/setup-buildx-action digest to d70bba7 ([#160](https://github.com/validator-labs/validator-plugin-oci/issues/160)) ([23e54a5](https://github.com/validator-labs/validator-plugin-oci/commit/23e54a50748ceb6acc7c07c689d67be7c58617fc))
* **deps:** update gcr.io/spectro-images-public/golang docker tag to v1.22 ([#105](https://github.com/validator-labs/validator-plugin-oci/issues/105)) ([c5edf12](https://github.com/validator-labs/validator-plugin-oci/commit/c5edf127bbd9d4ab44188b866250a36b795c4ffe))
* **deps:** update gcr.io/spectro-images-public/golang docker tag to v1.22.4 ([#186](https://github.com/validator-labs/validator-plugin-oci/issues/186)) ([3fe407d](https://github.com/validator-labs/validator-plugin-oci/commit/3fe407da4387ff4dbb57d37456ccd374d2e90ae4))
* **deps:** update helm/kind-action action to v1.10.0 ([#181](https://github.com/validator-labs/validator-plugin-oci/issues/181)) ([03458fd](https://github.com/validator-labs/validator-plugin-oci/commit/03458fd0f08cf58411831f81117b4330e17923d5))
* **deps:** update softprops/action-gh-release digest to 69320db ([#135](https://github.com/validator-labs/validator-plugin-oci/issues/135)) ([9989631](https://github.com/validator-labs/validator-plugin-oci/commit/9989631461b716ae9faa149b685a41530af129f3))


### Dependency Updates

* **deps:** update aws-sdk-go-v2 monorepo ([#204](https://github.com/validator-labs/validator-plugin-oci/issues/204)) ([0503bee](https://github.com/validator-labs/validator-plugin-oci/commit/0503bee4a591c982abf78f6568fd491fa2895b80))
* **deps:** update aws-sdk-go-v2 monorepo ([#205](https://github.com/validator-labs/validator-plugin-oci/issues/205)) ([196ef81](https://github.com/validator-labs/validator-plugin-oci/commit/196ef81516ce943c00cb758f3bb071383be7234f))
* **deps:** update aws-sdk-go-v2 monorepo ([#211](https://github.com/validator-labs/validator-plugin-oci/issues/211)) ([40350f0](https://github.com/validator-labs/validator-plugin-oci/commit/40350f029b7fb8d04890ed6c677b84744b42cd6f))
* **deps:** update dependency go to v1.22.5 ([#206](https://github.com/validator-labs/validator-plugin-oci/issues/206)) ([45ad3a8](https://github.com/validator-labs/validator-plugin-oci/commit/45ad3a8a48d543fdc51758a2447e3616fa1eefac))
* **deps:** update github.com/validator-labs/validator digest to de015d9 ([#218](https://github.com/validator-labs/validator-plugin-oci/issues/218)) ([a725d7f](https://github.com/validator-labs/validator-plugin-oci/commit/a725d7f1ee65a0e2146be3675bf840f59b8fbfe2))
* **deps:** update kubernetes packages to v0.30.3 ([#215](https://github.com/validator-labs/validator-plugin-oci/issues/215)) ([290ae5f](https://github.com/validator-labs/validator-plugin-oci/commit/290ae5f718ba6b0a5d46d643647b6b277101e2db))
* **deps:** update module github.com/google/go-containerregistry to v0.20.0 ([#207](https://github.com/validator-labs/validator-plugin-oci/issues/207)) ([aeeb24c](https://github.com/validator-labs/validator-plugin-oci/commit/aeeb24c276bdcc417e7313fcee5a61a96bd34090))
* **deps:** update module github.com/google/go-containerregistry to v0.20.1 ([#214](https://github.com/validator-labs/validator-plugin-oci/issues/214)) ([73525b1](https://github.com/validator-labs/validator-plugin-oci/commit/73525b1b6f31f4c51a3264461f3d4b72f2a4c8e5))
* **deps:** update module github.com/sigstore/sigstore to v1.8.6 ([#202](https://github.com/validator-labs/validator-plugin-oci/issues/202)) ([34d6274](https://github.com/validator-labs/validator-plugin-oci/commit/34d6274baaf41d45bb954fc5204d58db3b49a49f))
* **deps:** update module github.com/sigstore/sigstore to v1.8.7 ([#212](https://github.com/validator-labs/validator-plugin-oci/issues/212)) ([21a320a](https://github.com/validator-labs/validator-plugin-oci/commit/21a320a267a37cc540e1ca1404e8e2ecb597bb84))
* **deps:** update module github.com/validator-labs/validator to v0.0.44 ([#210](https://github.com/validator-labs/validator-plugin-oci/issues/210)) ([b7d8d5b](https://github.com/validator-labs/validator-plugin-oci/commit/b7d8d5b6022cad510e185976330e1331eb9ccc15))
* **deps:** update module github.com/validator-labs/validator to v0.0.46 ([#213](https://github.com/validator-labs/validator-plugin-oci/issues/213)) ([be1a840](https://github.com/validator-labs/validator-plugin-oci/commit/be1a84015077a2a23c775d6f888064da529bd872))
* **deps:** update module sigs.k8s.io/cluster-api to v1.7.4 ([#209](https://github.com/validator-labs/validator-plugin-oci/issues/209)) ([1e4bef0](https://github.com/validator-labs/validator-plugin-oci/commit/1e4bef0f3a58600d2a680c312dd537c13f77a204))


### Refactoring

* enable revive and address all lints ([#208](https://github.com/validator-labs/validator-plugin-oci/issues/208)) ([be2689d](https://github.com/validator-labs/validator-plugin-oci/commit/be2689d988d23ecf2e0b7a0ec61866ce9c80ab50))

## [0.0.10](https://github.com/validator-labs/validator-plugin-oci/compare/v0.0.9...v0.0.10) (2024-05-17)


### Other

* migrate from spectrocloud-labs to validator-labs ([#167](https://github.com/validator-labs/validator-plugin-oci/issues/167)) ([4a51c11](https://github.com/validator-labs/validator-plugin-oci/commit/4a51c117b26412f4d9e1358694c360f54446e649))

## [0.0.9](https://github.com/validator-labs/validator-plugin-oci/compare/v0.0.8...v0.0.9) (2024-05-15)


### Bug Fixes

* fix Ensure auth options are passed into the Cosign Verifier ([#166](https://github.com/validator-labs/validator-plugin-oci/issues/166)) ([26d4671](https://github.com/validator-labs/validator-plugin-oci/commit/26d46714d09c79d9dd4b9f3796ee0d2e3d9a5ad6))
* **deps:** update aws-sdk-go-v2 monorepo ([#151](https://github.com/validator-labs/validator-plugin-oci/issues/151)) ([4d6bfe5](https://github.com/validator-labs/validator-plugin-oci/commit/4d6bfe5759a340ca8380876fb3f99c1b31b20c72))
* **deps:** update kubernetes packages to v0.29.3 ([#148](https://github.com/validator-labs/validator-plugin-oci/issues/148)) ([cb599b9](https://github.com/validator-labs/validator-plugin-oci/commit/cb599b92152d36704b642eae7428877956a7f1e7))
* **deps:** update module github.com/google/go-containerregistry to v0.19.1 ([#147](https://github.com/validator-labs/validator-plugin-oci/issues/147)) ([413aa68](https://github.com/validator-labs/validator-plugin-oci/commit/413aa68a2822e954d4c16fe24ecd74c4da7531bf))
* **deps:** update module github.com/onsi/ginkgo/v2 to v2.17.0 ([#149](https://github.com/validator-labs/validator-plugin-oci/issues/149)) ([48dfcba](https://github.com/validator-labs/validator-plugin-oci/commit/48dfcbae5043da0591e11ef3fc65d6c85b398421))
* **deps:** update module github.com/onsi/ginkgo/v2 to v2.17.1 ([#153](https://github.com/validator-labs/validator-plugin-oci/issues/153)) ([06805c0](https://github.com/validator-labs/validator-plugin-oci/commit/06805c0f2445fa0a7b5006aec7da194874e5613d))
* **deps:** update module github.com/onsi/gomega to v1.32.0 ([#150](https://github.com/validator-labs/validator-plugin-oci/issues/150)) ([b04ef8d](https://github.com/validator-labs/validator-plugin-oci/commit/b04ef8db786c8ad0d804852606ab16041535863c))
* **deps:** update module github.com/sigstore/sigstore to v1.8.3 ([#157](https://github.com/validator-labs/validator-plugin-oci/issues/157)) ([8bbac88](https://github.com/validator-labs/validator-plugin-oci/commit/8bbac886c92f9decae14de301f78d32c35100769))
* **deps:** update module github.com/validator-labs/validator to v0.0.38 ([#95](https://github.com/validator-labs/validator-plugin-oci/issues/95)) ([5566f04](https://github.com/validator-labs/validator-plugin-oci/commit/5566f044527a811e862d333b7025868d62f78ef4))


### Other

* **deps:** update actions/setup-python digest to 82c7e63 ([#154](https://github.com/validator-labs/validator-plugin-oci/issues/154)) ([e3db8ed](https://github.com/validator-labs/validator-plugin-oci/commit/e3db8ed0b43f6cfb36f27a6c3eb2a478a259aa50))
* **deps:** update anchore/sbom-action action to v0.15.10 ([#156](https://github.com/validator-labs/validator-plugin-oci/issues/156)) ([c8c30aa](https://github.com/validator-labs/validator-plugin-oci/commit/c8c30aa74fea1c79eac1216d04385203f573ab14))
* **deps:** update codecov/codecov-action digest to c16abc2 ([#155](https://github.com/validator-labs/validator-plugin-oci/issues/155)) ([34c1c5c](https://github.com/validator-labs/validator-plugin-oci/commit/34c1c5c6fbafa73b161daf6ab8b632484c11b9d5))
* **deps:** update docker/build-push-action digest to 2cdde99 ([#143](https://github.com/validator-labs/validator-plugin-oci/issues/143)) ([c2f01b1](https://github.com/validator-labs/validator-plugin-oci/commit/c2f01b18d141c05907f3b4a191d3687930c388ed))
* **deps:** update docker/setup-buildx-action digest to 2b51285 ([#144](https://github.com/validator-labs/validator-plugin-oci/issues/144)) ([04a5e30](https://github.com/validator-labs/validator-plugin-oci/commit/04a5e30be8a2533839c0f97d58421413e8b3e421))
* **deps:** update gcr.io/kubebuilder/kube-rbac-proxy docker tag to v0.16.0 ([#158](https://github.com/validator-labs/validator-plugin-oci/issues/158)) ([58806a5](https://github.com/validator-labs/validator-plugin-oci/commit/58806a547be425ab3db76c5a30b08587ff570856))

## [0.0.8](https://github.com/validator-labs/validator-plugin-oci/compare/v0.0.7...v0.0.8) (2024-03-13)


### Features

* update helm chart to support signature verification and full layer validation ([#142](https://github.com/validator-labs/validator-plugin-oci/issues/142)) ([d7c9814](https://github.com/validator-labs/validator-plugin-oci/commit/d7c9814b8ab3c4e430b8a01a79d3f554191f3ca7))


### Bug Fixes

* **deps:** update github.com/validator-labs/validator digest to fc351f3 ([#137](https://github.com/validator-labs/validator-plugin-oci/issues/137)) ([c7c875f](https://github.com/validator-labs/validator-plugin-oci/commit/c7c875f2231264a09e7e132f92c9b422e3707759))
* **deps:** update module sigs.k8s.io/cluster-api to v1.6.3 ([#138](https://github.com/validator-labs/validator-plugin-oci/issues/138)) ([4f9947c](https://github.com/validator-labs/validator-plugin-oci/commit/4f9947cfbc99251f662513d07b3f9b8bb63b7e1c))
* ensure error is returned when signature verification is enabled but no public keys are provided ([#141](https://github.com/validator-labs/validator-plugin-oci/issues/141)) ([1e97bdf](https://github.com/validator-labs/validator-plugin-oci/commit/1e97bdf826faa05b26ba87eda9c61a8c872a4eb5))


### Other

* **deps:** update docker/login-action digest to e92390c ([#140](https://github.com/validator-labs/validator-plugin-oci/issues/140)) ([4343b61](https://github.com/validator-labs/validator-plugin-oci/commit/4343b618689b7f9318314b03fecae3524efcecee))

## [0.0.7](https://github.com/validator-labs/validator-plugin-oci/compare/v0.0.6...v0.0.7) (2024-03-12)


### Other

* **deps:** update google-github-actions/release-please-action digest to a37ac6e ([#131](https://github.com/validator-labs/validator-plugin-oci/issues/131)) ([fe6ad46](https://github.com/validator-labs/validator-plugin-oci/commit/fe6ad4697ca7b5792e5b4c8a703f7ec06c68e6ff))
* **deps:** update softprops/action-gh-release digest to 3198ee1 ([#130](https://github.com/validator-labs/validator-plugin-oci/issues/130)) ([857cedc](https://github.com/validator-labs/validator-plugin-oci/commit/857cedcb159b7aa66d17bc229b7adab65e8cd2c3))


### Refactoring

* use patch helpers ([#136](https://github.com/validator-labs/validator-plugin-oci/issues/136)) ([b25d6d0](https://github.com/validator-labs/validator-plugin-oci/commit/b25d6d06776c4f26bcdbf30c0c78158d1321c05b))

## [0.0.6](https://github.com/validator-labs/validator-plugin-oci/compare/v0.0.5...v0.0.6) (2024-03-11)


### Bug Fixes

* **deps:** update aws-sdk-go-v2 monorepo ([#123](https://github.com/validator-labs/validator-plugin-oci/issues/123)) ([aeddab7](https://github.com/validator-labs/validator-plugin-oci/commit/aeddab7a1175f8e7d2fb872bcbd8b3d9a6cde3fd))
* **deps:** update module github.com/onsi/ginkgo/v2 to v2.16.0 ([#122](https://github.com/validator-labs/validator-plugin-oci/issues/122)) ([22cd704](https://github.com/validator-labs/validator-plugin-oci/commit/22cd70405de123ef3b43010e13cd978bbe334f6d))
* **deps:** update module github.com/stretchr/testify to v1.9.0 ([#121](https://github.com/validator-labs/validator-plugin-oci/issues/121)) ([0cbc5bd](https://github.com/validator-labs/validator-plugin-oci/commit/0cbc5bdb104e40e0e003b9f89a66cb1a249c7a1c))


### Other

* **deps:** update anchore/sbom-action action to v0.15.9 ([#124](https://github.com/validator-labs/validator-plugin-oci/issues/124)) ([89ad98a](https://github.com/validator-labs/validator-plugin-oci/commit/89ad98a123e91b1e39f0c55d0efdc8e07e423f70))
* **deps:** update azure/setup-helm action to v4 ([#119](https://github.com/validator-labs/validator-plugin-oci/issues/119)) ([77978fa](https://github.com/validator-labs/validator-plugin-oci/commit/77978fa82d46696aa8cc649afacd1954816fe084))
* **deps:** update azure/setup-helm digest to b7246b1 ([#125](https://github.com/validator-labs/validator-plugin-oci/issues/125)) ([b5f0347](https://github.com/validator-labs/validator-plugin-oci/commit/b5f03474419840f3d000a8fc6e04a84be1ece00e))
* **deps:** update docker/build-push-action digest to af5a7ed ([#126](https://github.com/validator-labs/validator-plugin-oci/issues/126)) ([d3c13a7](https://github.com/validator-labs/validator-plugin-oci/commit/d3c13a7da4a05a7ea448e37cb9a96afaecd5709e))
* **deps:** update softprops/action-gh-release action to v2 ([#127](https://github.com/validator-labs/validator-plugin-oci/issues/127)) ([0494998](https://github.com/validator-labs/validator-plugin-oci/commit/0494998490147401a9881d99a72f93e9434f1cfc))
* **deps:** update softprops/action-gh-release digest to d99959e ([#129](https://github.com/validator-labs/validator-plugin-oci/issues/129)) ([00f28bf](https://github.com/validator-labs/validator-plugin-oci/commit/00f28bf4d93b7e1a569fcb145a79de5098d67e18))
* upgrade to validator v0.0.36 ([#128](https://github.com/validator-labs/validator-plugin-oci/issues/128)) ([50abf22](https://github.com/validator-labs/validator-plugin-oci/commit/50abf22f49637bdd2ee20baae8b7add0c4571f2c))

## [0.0.5](https://github.com/validator-labs/validator-plugin-oci/compare/v0.0.4...v0.0.5) (2024-02-28)


### Features

* add public key signature verification support ([#112](https://github.com/validator-labs/validator-plugin-oci/issues/112)) ([3a0b166](https://github.com/validator-labs/validator-plugin-oci/commit/3a0b166c554e717a9432ce95ab611e3827e728a6))
* introduce LayerValidation bool to expose go-containerregistry fast validation ([#110](https://github.com/validator-labs/validator-plugin-oci/issues/110)) ([68f5ae4](https://github.com/validator-labs/validator-plugin-oci/commit/68f5ae4f4168cc53904b401bdbbdcf9600508aa9))


### Bug Fixes

* **deps:** update aws-sdk-go-v2 monorepo ([#101](https://github.com/validator-labs/validator-plugin-oci/issues/101)) ([540842e](https://github.com/validator-labs/validator-plugin-oci/commit/540842e4910594e6ddb054650e6b671092279bf0))
* **deps:** update aws-sdk-go-v2 monorepo ([#107](https://github.com/validator-labs/validator-plugin-oci/issues/107)) ([ef46bbc](https://github.com/validator-labs/validator-plugin-oci/commit/ef46bbcd751d76bca2b3169107c91f091dac56cf))
* **deps:** update aws-sdk-go-v2 monorepo ([#109](https://github.com/validator-labs/validator-plugin-oci/issues/109)) ([15d5e1b](https://github.com/validator-labs/validator-plugin-oci/commit/15d5e1b5f45c55362fe553826785e28ff5944d64))
* **deps:** update aws-sdk-go-v2 monorepo ([#111](https://github.com/validator-labs/validator-plugin-oci/issues/111)) ([ef4cafe](https://github.com/validator-labs/validator-plugin-oci/commit/ef4cafef63e99af23a5887bd3fd618dba941667b))
* **deps:** update aws-sdk-go-v2 monorepo ([#114](https://github.com/validator-labs/validator-plugin-oci/issues/114)) ([b0830cd](https://github.com/validator-labs/validator-plugin-oci/commit/b0830cde070919921ec0ef4bfcd1e71d8cfe23d4))
* **deps:** update kubernetes packages to v0.29.2 ([#102](https://github.com/validator-labs/validator-plugin-oci/issues/102)) ([b6f58bc](https://github.com/validator-labs/validator-plugin-oci/commit/b6f58bcf497f5cc9bd35a4806ea03d2cb16e5fe5))
* **deps:** update module github.com/sigstore/sigstore to v1.8.2 ([#117](https://github.com/validator-labs/validator-plugin-oci/issues/117)) ([8f36e4b](https://github.com/validator-labs/validator-plugin-oci/commit/8f36e4b59434ba025da49b57726dd6af2d06bd01))
* **deps:** update module sigs.k8s.io/controller-runtime to v0.17.1 ([#75](https://github.com/validator-labs/validator-plugin-oci/issues/75)) ([adf2a72](https://github.com/validator-labs/validator-plugin-oci/commit/adf2a725447df8addc93665f8721e6ed4cc596ca))
* **deps:** update module sigs.k8s.io/controller-runtime to v0.17.2 ([#108](https://github.com/validator-labs/validator-plugin-oci/issues/108)) ([9d49ce1](https://github.com/validator-labs/validator-plugin-oci/commit/9d49ce1b25779e9b72509c28738da0a99ebec1a5))


### Other

* **deps:** update actions/upload-artifact digest to 5d5d22a ([#97](https://github.com/validator-labs/validator-plugin-oci/issues/97)) ([c21c3f4](https://github.com/validator-labs/validator-plugin-oci/commit/c21c3f4c697ed1dae9ca7497ad6c31c303edfcc9))
* **deps:** update codecov/codecov-action digest to 0cfda1d ([#113](https://github.com/validator-labs/validator-plugin-oci/issues/113)) ([77fc85a](https://github.com/validator-labs/validator-plugin-oci/commit/77fc85a7b03543003540d6506f759d07dd1070fb))
* **deps:** update codecov/codecov-action digest to 54bcd87 ([#115](https://github.com/validator-labs/validator-plugin-oci/issues/115)) ([50b800d](https://github.com/validator-labs/validator-plugin-oci/commit/50b800dd4fc5f72cf7bac704d188f5154d8f3468))
* **deps:** update docker/setup-buildx-action digest to 0d103c3 ([#116](https://github.com/validator-labs/validator-plugin-oci/issues/116)) ([f3a2b0d](https://github.com/validator-labs/validator-plugin-oci/commit/f3a2b0d9915668da3c495e20bda82f8506248d3c))
* **deps:** update helm/kind-action action to v1.9.0 ([#100](https://github.com/validator-labs/validator-plugin-oci/issues/100)) ([1ab8d01](https://github.com/validator-labs/validator-plugin-oci/commit/1ab8d012d00818c767b5023abe7dc3b272795eaa))
* fix broken build link in README ([#118](https://github.com/validator-labs/validator-plugin-oci/issues/118)) ([7aa2e93](https://github.com/validator-labs/validator-plugin-oci/commit/7aa2e93e7b9d795fd73a7d6f2377f187a9595cd7))

## [0.0.4](https://github.com/validator-labs/validator-plugin-oci/compare/v0.0.3...v0.0.4) (2024-02-06)


### Other

* update validator ([7717344](https://github.com/validator-labs/validator-plugin-oci/commit/7717344f5de4406b89caf5923a05370b8e0844a3))

## [0.0.3](https://github.com/validator-labs/validator-plugin-oci/compare/v0.0.2...v0.0.3) (2024-02-05)


### Bug Fixes

* CRD validation for rule host uniqueness ([#56](https://github.com/validator-labs/validator-plugin-oci/issues/56)) ([8dbdc15](https://github.com/validator-labs/validator-plugin-oci/commit/8dbdc15a2d23225e94630771d26eab26439c721c))
* **deps:** update aws-sdk-go-v2 monorepo ([#55](https://github.com/validator-labs/validator-plugin-oci/issues/55)) ([af7f8a4](https://github.com/validator-labs/validator-plugin-oci/commit/af7f8a47423f262b9d491b9bac7bda3ba8c21ac8))
* **deps:** update aws-sdk-go-v2 monorepo ([#61](https://github.com/validator-labs/validator-plugin-oci/issues/61)) ([b733807](https://github.com/validator-labs/validator-plugin-oci/commit/b7338076acaac8d86965a826cc2b27b1c626390b))
* **deps:** update aws-sdk-go-v2 monorepo ([#67](https://github.com/validator-labs/validator-plugin-oci/issues/67)) ([c1c5d0e](https://github.com/validator-labs/validator-plugin-oci/commit/c1c5d0e2543c53e4d36c8affb484c0847c6d6275))
* **deps:** update aws-sdk-go-v2 monorepo ([#76](https://github.com/validator-labs/validator-plugin-oci/issues/76)) ([55d84a8](https://github.com/validator-labs/validator-plugin-oci/commit/55d84a85bd44a435371068ea0329455332008bee))
* **deps:** update aws-sdk-go-v2 monorepo ([#81](https://github.com/validator-labs/validator-plugin-oci/issues/81)) ([1b4d64d](https://github.com/validator-labs/validator-plugin-oci/commit/1b4d64d305faaa6008e29d7d7723f3acf9188029))
* **deps:** update module github.com/aws/aws-sdk-go-v2/config to v1.25.11 ([#48](https://github.com/validator-labs/validator-plugin-oci/issues/48)) ([8567bef](https://github.com/validator-labs/validator-plugin-oci/commit/8567bef0cf026b79fb1fff133eb90da9d351d652))
* **deps:** update module github.com/aws/aws-sdk-go-v2/config to v1.26.6 ([#85](https://github.com/validator-labs/validator-plugin-oci/issues/85)) ([939b7cc](https://github.com/validator-labs/validator-plugin-oci/commit/939b7ccd1d5ec2d9be9649b9fb0c51f449e57b24))
* **deps:** update module github.com/aws/aws-sdk-go-v2/service/ecr to v1.24.2 ([#47](https://github.com/validator-labs/validator-plugin-oci/issues/47)) ([7275869](https://github.com/validator-labs/validator-plugin-oci/commit/727586939c110c34c2c60786279356ae40c3131b))
* **deps:** update module github.com/google/go-containerregistry to v0.18.0 ([#77](https://github.com/validator-labs/validator-plugin-oci/issues/77)) ([bfe4961](https://github.com/validator-labs/validator-plugin-oci/commit/bfe49617b3649e86f50a1062a5b8d923dade8f7c))
* **deps:** update module github.com/google/go-containerregistry to v0.19.0 ([#89](https://github.com/validator-labs/validator-plugin-oci/issues/89)) ([d07fd92](https://github.com/validator-labs/validator-plugin-oci/commit/d07fd92e2a2f08d6b054661a6132253b30240e6d))
* **deps:** update module github.com/onsi/ginkgo/v2 to v2.14.0 ([#73](https://github.com/validator-labs/validator-plugin-oci/issues/73)) ([f340000](https://github.com/validator-labs/validator-plugin-oci/commit/f34000094b81b712a86ced4b28b21aa2c7d295d9))
* **deps:** update module github.com/onsi/ginkgo/v2 to v2.15.0 ([#78](https://github.com/validator-labs/validator-plugin-oci/issues/78)) ([d0599fb](https://github.com/validator-labs/validator-plugin-oci/commit/d0599fbd0facae1048feedfabb037d564c37fc72))
* **deps:** update module github.com/onsi/gomega to v1.31.0 ([#79](https://github.com/validator-labs/validator-plugin-oci/issues/79)) ([d6d17a9](https://github.com/validator-labs/validator-plugin-oci/commit/d6d17a92de92d9167533b354f3f5bacac5ab033c))
* **deps:** update module github.com/onsi/gomega to v1.31.1 ([#83](https://github.com/validator-labs/validator-plugin-oci/issues/83)) ([731a624](https://github.com/validator-labs/validator-plugin-oci/commit/731a6241b95df5dea83a8bac65b93801ebf3c25c))
* **deps:** update module github.com/validator-labs/validator to v0.0.28 ([#52](https://github.com/validator-labs/validator-plugin-oci/issues/52)) ([4fb5e57](https://github.com/validator-labs/validator-plugin-oci/commit/4fb5e57a0065602ef699b03c6c72928d246dbfb0))
* **deps:** update module github.com/validator-labs/validator to v0.0.30 ([#66](https://github.com/validator-labs/validator-plugin-oci/issues/66)) ([dfc8fd7](https://github.com/validator-labs/validator-plugin-oci/commit/dfc8fd7eeddd644eec71ebc9402940f6bfbc8b98))
* **deps:** update module github.com/validator-labs/validator to v0.0.32 ([#69](https://github.com/validator-labs/validator-plugin-oci/issues/69)) ([105f4ce](https://github.com/validator-labs/validator-plugin-oci/commit/105f4ce027834c243998926dbd2c2df7646daf43))
* ensure codecov is run when code is pushed to main ([#59](https://github.com/validator-labs/validator-plugin-oci/issues/59)) ([22da463](https://github.com/validator-labs/validator-plugin-oci/commit/22da46306fc2d96f6aa10baaad0ff347a4ceb139))


### Other

* bump validator plugin version to support rule addition ([#96](https://github.com/validator-labs/validator-plugin-oci/issues/96)) ([4ca7680](https://github.com/validator-labs/validator-plugin-oci/commit/4ca7680c1107cefc613c04a78497167bb79b2a9e))
* **deps:** update actions/setup-go action to v5 ([#54](https://github.com/validator-labs/validator-plugin-oci/issues/54)) ([c2945d2](https://github.com/validator-labs/validator-plugin-oci/commit/c2945d29463a624530129707545e4b75a67bcc15))
* **deps:** update actions/setup-python action to v5 ([#53](https://github.com/validator-labs/validator-plugin-oci/issues/53)) ([2afc37b](https://github.com/validator-labs/validator-plugin-oci/commit/2afc37bda4e3f6e0e23ab57938f33af4ada2cd59))
* **deps:** update actions/upload-artifact action to v4 ([#63](https://github.com/validator-labs/validator-plugin-oci/issues/63)) ([4fbdafc](https://github.com/validator-labs/validator-plugin-oci/commit/4fbdafc6dbe48a19049667f02860815b469c83fe))
* **deps:** update actions/upload-artifact digest to 1eb3cb2 ([#74](https://github.com/validator-labs/validator-plugin-oci/issues/74)) ([84e1f0a](https://github.com/validator-labs/validator-plugin-oci/commit/84e1f0a1cc3a1b51785277c2333125c14ead1c2f))
* **deps:** update actions/upload-artifact digest to 26f96df ([#86](https://github.com/validator-labs/validator-plugin-oci/issues/86)) ([82e2806](https://github.com/validator-labs/validator-plugin-oci/commit/82e280628f1f5438162c26b50bdbf59d08370c43))
* **deps:** update actions/upload-artifact digest to 694cdab ([#82](https://github.com/validator-labs/validator-plugin-oci/issues/82)) ([7518f3f](https://github.com/validator-labs/validator-plugin-oci/commit/7518f3f94357116c9998423da3c5efaae4d8da5e))
* **deps:** update anchore/sbom-action action to v0.15.1 ([#51](https://github.com/validator-labs/validator-plugin-oci/issues/51)) ([ebf1d17](https://github.com/validator-labs/validator-plugin-oci/commit/ebf1d1770babd60e53a642041fe5d026f56c8838))
* **deps:** update anchore/sbom-action action to v0.15.2 ([#70](https://github.com/validator-labs/validator-plugin-oci/issues/70)) ([8e4c2f8](https://github.com/validator-labs/validator-plugin-oci/commit/8e4c2f82ffde501af505b08404a65927754b859d))
* **deps:** update anchore/sbom-action action to v0.15.3 ([#71](https://github.com/validator-labs/validator-plugin-oci/issues/71)) ([0e3dea8](https://github.com/validator-labs/validator-plugin-oci/commit/0e3dea878db1a962d81e11effd0737fb30d95eaf))
* **deps:** update anchore/sbom-action action to v0.15.4 ([#80](https://github.com/validator-labs/validator-plugin-oci/issues/80)) ([e1771fa](https://github.com/validator-labs/validator-plugin-oci/commit/e1771fa7d145004931095b3c68ff9ce424b388cc))
* **deps:** update anchore/sbom-action action to v0.15.5 ([#84](https://github.com/validator-labs/validator-plugin-oci/issues/84)) ([9576d15](https://github.com/validator-labs/validator-plugin-oci/commit/9576d15e3ac41592bb54d660e0fbdad655bb1439))
* **deps:** update anchore/sbom-action action to v0.15.7 ([#88](https://github.com/validator-labs/validator-plugin-oci/issues/88)) ([8367112](https://github.com/validator-labs/validator-plugin-oci/commit/83671123ed32ceabac34c46c1f14b78b354ec73e))
* **deps:** update anchore/sbom-action action to v0.15.8 ([#91](https://github.com/validator-labs/validator-plugin-oci/issues/91)) ([f0ee030](https://github.com/validator-labs/validator-plugin-oci/commit/f0ee030fb7383e5b46a333f1418bc27a0451457e))
* **deps:** update codecov/codecov-action digest to 4fe8c5f ([#87](https://github.com/validator-labs/validator-plugin-oci/issues/87)) ([1c8bf5a](https://github.com/validator-labs/validator-plugin-oci/commit/1c8bf5af72b3683cc085401bacdc79daa63376bc))
* **deps:** update codecov/codecov-action digest to ab904c4 ([#90](https://github.com/validator-labs/validator-plugin-oci/issues/90)) ([94a4bdf](https://github.com/validator-labs/validator-plugin-oci/commit/94a4bdf3b1972294c0aa76253b4e2043815ee976))
* **deps:** update codecov/codecov-action digest to e0b68c6 ([#94](https://github.com/validator-labs/validator-plugin-oci/issues/94)) ([31487d6](https://github.com/validator-labs/validator-plugin-oci/commit/31487d662d1b8f4dd741cb638c1fd25ab1b893eb))
* **deps:** update gcr.io/spectro-images-public/golang docker tag to v1.22 ([#72](https://github.com/validator-labs/validator-plugin-oci/issues/72)) ([ff3312c](https://github.com/validator-labs/validator-plugin-oci/commit/ff3312cd60da5cdb3aa0ee73ad1903ba5dff9a3c))
* **deps:** update google-github-actions/release-please-action action to v4 ([#50](https://github.com/validator-labs/validator-plugin-oci/issues/50)) ([78956e2](https://github.com/validator-labs/validator-plugin-oci/commit/78956e209a28d7cb4d756e360ecad096c3c48576))
* **deps:** update google-github-actions/release-please-action digest to a2d8d68 ([#58](https://github.com/validator-labs/validator-plugin-oci/issues/58)) ([23821e5](https://github.com/validator-labs/validator-plugin-oci/commit/23821e5a4388b5e519383c385ab24e75f9210893))
* **deps:** update google-github-actions/release-please-action digest to cc61a07 ([#64](https://github.com/validator-labs/validator-plugin-oci/issues/64)) ([b1dee9b](https://github.com/validator-labs/validator-plugin-oci/commit/b1dee9bf69fd3ed5114a8c330711c69eec0ce913))

## [0.0.2](https://github.com/validator-labs/validator-plugin-oci/compare/v0.0.1...v0.0.2) (2023-11-30)


### Bug Fixes

* **deps:** update module github.com/aws/aws-sdk-go-v2/config to v1.25.10 ([#44](https://github.com/validator-labs/validator-plugin-oci/issues/44)) ([4e221d7](https://github.com/validator-labs/validator-plugin-oci/commit/4e221d7cd8868c4677a84daaee5e1382c20f3ba5))
* **deps:** update module github.com/aws/aws-sdk-go-v2/config to v1.25.9 ([#37](https://github.com/validator-labs/validator-plugin-oci/issues/37)) ([74b6eae](https://github.com/validator-labs/validator-plugin-oci/commit/74b6eae6b16df62c9866898b76d0ceaf95edbe4f))
* **deps:** update module github.com/aws/aws-sdk-go-v2/service/ecr to v1.24.0 ([#39](https://github.com/validator-labs/validator-plugin-oci/issues/39)) ([d6d4314](https://github.com/validator-labs/validator-plugin-oci/commit/d6d4314c9e21a89c22542f099117bbb7ab20a5e4))
* **deps:** update module github.com/aws/aws-sdk-go-v2/service/ecr to v1.24.1 ([#45](https://github.com/validator-labs/validator-plugin-oci/issues/45)) ([9ba631c](https://github.com/validator-labs/validator-plugin-oci/commit/9ba631cc4b624cbfd8941d64b451a61b77392775))
* **deps:** update module github.com/google/go-containerregistry to v0.17.0 ([#42](https://github.com/validator-labs/validator-plugin-oci/issues/42)) ([a6c2c58](https://github.com/validator-labs/validator-plugin-oci/commit/a6c2c582ffa74a0c97a87f318f2b763029725037))
* **deps:** update module github.com/validator-labs/validator to v0.0.26 ([#36](https://github.com/validator-labs/validator-plugin-oci/issues/36)) ([2c18421](https://github.com/validator-labs/validator-plugin-oci/commit/2c184212663c4ac97048b50fb2399f163f519036))
* **deps:** update module github.com/validator-labs/validator to v0.0.27 ([#43](https://github.com/validator-labs/validator-plugin-oci/issues/43)) ([1113a9e](https://github.com/validator-labs/validator-plugin-oci/commit/1113a9ea320a34de327970093542681998959669))
* fix link to oci issues in readme ([#41](https://github.com/validator-labs/validator-plugin-oci/issues/41)) ([b3c1cea](https://github.com/validator-labs/validator-plugin-oci/commit/b3c1cea77e46bab05727be8bca15b85f2687c6e4))
* update leader election id ([#46](https://github.com/validator-labs/validator-plugin-oci/issues/46)) ([976487b](https://github.com/validator-labs/validator-plugin-oci/commit/976487bfe8edaebe638d3cf067f787d5ec2385b0))

## [0.0.1](https://github.com/validator-labs/validator-plugin-oci/compare/v0.0.1...v0.0.1) (2023-11-29)


### Features

* add Helm chart ([#25](https://github.com/validator-labs/validator-plugin-oci/issues/25)) ([f4295ae](https://github.com/validator-labs/validator-plugin-oci/commit/f4295ae9a509c52763c12ba01458d8d0150b0bae))
* allow initContainer image to be passed in via values.yaml ([#27](https://github.com/validator-labs/validator-plugin-oci/issues/27)) ([50c8647](https://github.com/validator-labs/validator-plugin-oci/commit/50c8647f76cc70453b1ec1a5f7e307fcda839235))
* implement OCI registry validation spec ([#6](https://github.com/validator-labs/validator-plugin-oci/issues/6)) ([f62c494](https://github.com/validator-labs/validator-plugin-oci/commit/f62c494d3a44bcf99c9d0bccecd1af2b8bc3ae78))
* support validating list of oci artifacts ([#16](https://github.com/validator-labs/validator-plugin-oci/issues/16)) ([d0cbecc](https://github.com/validator-labs/validator-plugin-oci/commit/d0cbecc24614a9a6ddf2a34e71e01ce23a313d8c))


### Bug Fixes

* **deps:** update kubernetes packages to v0.28.4 ([#17](https://github.com/validator-labs/validator-plugin-oci/issues/17)) ([f346f63](https://github.com/validator-labs/validator-plugin-oci/commit/f346f631c50d2fdc6236603055c792f116c554df))
* **deps:** update module github.com/aws/aws-sdk-go-v2/config to v1.25.6 ([#28](https://github.com/validator-labs/validator-plugin-oci/issues/28)) ([b12dabe](https://github.com/validator-labs/validator-plugin-oci/commit/b12dabe9730e9e12a48e979f796fde71dbd551a0))
* **deps:** update module github.com/aws/aws-sdk-go-v2/config to v1.25.8 ([#32](https://github.com/validator-labs/validator-plugin-oci/issues/32)) ([3eb0824](https://github.com/validator-labs/validator-plugin-oci/commit/3eb08241bd645d73cd50182fd562f941171b4a30))
* **deps:** update module github.com/aws/aws-sdk-go-v2/service/ecr to v1.23.2 ([#30](https://github.com/validator-labs/validator-plugin-oci/issues/30)) ([6375b2b](https://github.com/validator-labs/validator-plugin-oci/commit/6375b2bafbdcaa8649691eaba15ba52ec8eb80d9))
* **deps:** update module github.com/aws/aws-sdk-go-v2/service/ecr to v1.23.3 ([#34](https://github.com/validator-labs/validator-plugin-oci/issues/34)) ([a55ada3](https://github.com/validator-labs/validator-plugin-oci/commit/a55ada393e0ef05510388a152b4e0a03a573d3d4))
* **deps:** update module github.com/onsi/ginkgo/v2 to v2.13.1 ([#15](https://github.com/validator-labs/validator-plugin-oci/issues/15)) ([23673ac](https://github.com/validator-labs/validator-plugin-oci/commit/23673ac0092fac7eecc78f7d92c249b385537c39))
* **deps:** update module github.com/onsi/ginkgo/v2 to v2.13.2 ([#35](https://github.com/validator-labs/validator-plugin-oci/issues/35)) ([4f44a26](https://github.com/validator-labs/validator-plugin-oci/commit/4f44a26a67d141a8b953cd937d07c1c0482087eb))
* **deps:** update module github.com/validator-labs/validator to v0.0.18 ([#14](https://github.com/validator-labs/validator-plugin-oci/issues/14)) ([58c78f4](https://github.com/validator-labs/validator-plugin-oci/commit/58c78f43a7f21d8d22381042ea73fbe5f3b7f0d0))
* **deps:** update module github.com/validator-labs/validator to v0.0.21 ([#18](https://github.com/validator-labs/validator-plugin-oci/issues/18)) ([9373c1d](https://github.com/validator-labs/validator-plugin-oci/commit/9373c1d3541397948eca4a93df61c3a628661b56))
* **deps:** update module github.com/validator-labs/validator to v0.0.25 ([#21](https://github.com/validator-labs/validator-plugin-oci/issues/21)) ([76f1b24](https://github.com/validator-labs/validator-plugin-oci/commit/76f1b247a7bf69d990a0539e8ec73260cfe7ad5a))
* set owner references on validation result to ensure cleanup ([#19](https://github.com/validator-labs/validator-plugin-oci/issues/19)) ([9c7c28d](https://github.com/validator-labs/validator-plugin-oci/commit/9c7c28d1e69b9488263537e48415818826d96ebf))


### Other

* add license badge ([1eb5f1b](https://github.com/validator-labs/validator-plugin-oci/commit/1eb5f1b2ceafc7656816f42b4f51c11ad0057aba))
* **deps:** pin dependencies ([#9](https://github.com/validator-labs/validator-plugin-oci/issues/9)) ([9876cd7](https://github.com/validator-labs/validator-plugin-oci/commit/9876cd701be178016231d02661a78db1f2f48c85))
* **deps:** update actions/checkout action to v4 ([#10](https://github.com/validator-labs/validator-plugin-oci/issues/10)) ([cd110af](https://github.com/validator-labs/validator-plugin-oci/commit/cd110af99d4eed651d89dabe5565bcedcb3f4c35))
* **deps:** update anchore/sbom-action action to v0.15.0 ([#23](https://github.com/validator-labs/validator-plugin-oci/issues/23)) ([34253f0](https://github.com/validator-labs/validator-plugin-oci/commit/34253f03e491ebecc0ce8631d56558cc16bb4b82))
* **deps:** update docker/build-push-action digest to 4a13e50 ([#20](https://github.com/validator-labs/validator-plugin-oci/issues/20)) ([eace63e](https://github.com/validator-labs/validator-plugin-oci/commit/eace63e7d49fc14c8d1f8d0427bd11039bef140d))
* fix platform specification for manager image ([#13](https://github.com/validator-labs/validator-plugin-oci/issues/13)) ([539e8be](https://github.com/validator-labs/validator-plugin-oci/commit/539e8be372a623125d1ed04e602833c59acddd93))
* release 0.0.1 ([d4b32af](https://github.com/validator-labs/validator-plugin-oci/commit/d4b32afa1737b2ca4dc39e907bbb4bee871e15fc))
* specify platform in Dockerfile and docker-build make target ([#12](https://github.com/validator-labs/validator-plugin-oci/issues/12)) ([a88d182](https://github.com/validator-labs/validator-plugin-oci/commit/a88d1820503bdfc6c2f99690db2a0bcd6befc5dc))
* switch back to public bulwark images ([010e7f8](https://github.com/validator-labs/validator-plugin-oci/commit/010e7f842a54cb0f9e0f572618007ad85009f766))
* update validator-labs/validator dependency to v0.0.15 ([f62c494](https://github.com/validator-labs/validator-plugin-oci/commit/f62c494d3a44bcf99c9d0bccecd1af2b8bc3ae78))


### Refactoring

* switch from oras to go-containerregistry ([#24](https://github.com/validator-labs/validator-plugin-oci/issues/24)) ([eef0013](https://github.com/validator-labs/validator-plugin-oci/commit/eef0013a7d1072f55bb3356304f287ad1cc61ff4))
* switch init container to image with ca-certificates pre installed ([#33](https://github.com/validator-labs/validator-plugin-oci/issues/33)) ([4550f4b](https://github.com/validator-labs/validator-plugin-oci/commit/4550f4bedb9807d8578fcc56d7fc4e3309cd6d8b))
