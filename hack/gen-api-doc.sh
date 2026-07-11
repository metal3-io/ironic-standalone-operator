#!/bin/bash

set -eu

CRDOC="${CRDOC:-crdoc}"

"${CRDOC}" --resources config/crd/bases/ --output docs/api.md

# FIXME: crdoc generates links to old k8s documentation, link checker complains
sed -i -e 's#/v1.20/#/v1.32/#' docs/api.md

