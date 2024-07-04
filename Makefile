include build/makelib/common.mk
include build/makelib/plugin.mk

# Image URL to use all building/pushing image targets
IMG ?= quay.io/validator-labs/validator-plugin-oci:latest

# Helm vars
CHART_NAME=validator-plugin-oci

.PHONY: dev
dev: ## Run a controller via devspace
	devspace dev -n validator-plugin-oci-system
