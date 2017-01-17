#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

if ! which go > /dev/null; then
	echo "golang needs to be installed"
	exit 1
fi

if ! which docker > /dev/null; then
	echo "docker needs to be installed"
	exit 1
fi

: ${IMAGE:?"Need to set IMAGE, e.g. bitesize-registry.default.svc.cluster.local:5000/bitesize/environment-operator"}

bin_dir="_output/bin"
mkdir -p ${bin_dir} || true

GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build \
	-o ${bin_dir}/environment-operator ./cmd/operator/main.go


docker build --tag "${IMAGE}" -f hack/build/Dockerfile . 1>/dev/null
# For gcr users, do "gcloud docker -a" to have access.
docker push "${IMAGE}" 1>/dev/null
