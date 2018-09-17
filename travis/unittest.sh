#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

export RELEASE="dev"
export branch=${TRAVIS_BRANCH}
echo "export RELEASE=$RELEASE" >> ~/environment.sh
echo "export branch=$branch" >> ~/environment.sh

echo "**************************************************************"
echo "***************** Running Unit Tests *************************"
echo "**************************************************************"

BUILD_IMAGE="golang:1.10-alpine"

docker run --rm -v "$(pwd)":/go/src/github.com/pearsontechnology/environment-operator \
    -w /go/src/github.com/pearsontechnology/environment-operator \
    -e CGO_ENABLED=1 \
    ${BUILD_IMAGE} \
    /bin/sh -c "apk update && apk add git && go test -v ./pkg/bitesize ./pkg/cluster ./pkg/diff ./pkg/git ./pkg/reaper ./pkg/translator ./pkg/web ./pkg/util ./pkg/util/k8s"

