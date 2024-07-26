include build/makelib/common.mk
include build/makelib/plugin.mk

# Image URL to use all building/pushing image targets
IMG ?= quay.io/validator-labs/validator-plugin-oci:latest

# Helm vars
CHART_NAME=validator-plugin-oci

.PHONY: dev
dev: ## Run a controller via devspace
	devspace dev -n validator-plugin-oci-system

# Static Analysis / CI

chartCrds = chart/validator-plugin-oci/crds/validation.spectrocloud.labs_ocivalidators.yaml

reviewable-ext:
	rm $(chartCrds)
	cp config/crd/bases/validation.spectrocloud.labs_ocivalidators.yaml $(chartCrds)
