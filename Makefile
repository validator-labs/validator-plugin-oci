# -include will silently skip missing files, which allows us
# to load those files with a target in the Makefile. If only
# "include" was used, the make command would fail and refuse
# to run a target until the include commands succeeded.
-include build/makelib/common.mk

# Image URL to use all building/pushing image targets
IMG ?= quay.io/validator-labs/validator-plugin-oci:latest

# Helm vars
CHART_NAME=validator-plugin-oci

.PHONY: dev
dev: ## Run a controller via devspace
	devspace dev -n validator-plugin-oci-system
