# environment-operator


More information on environment operator can be found in https://docs.google.com/document/d/1VYarE5SepyvVFjkXJfnpM3JAgybxkBT53TZzqsOv_pc

## Running Tests

Unit tests are stored in `test/unit`. To run full unit test suite, run:

```
% go test -v ./test/unit
```

End-to-end tests are located in `test/e2e` (currently none).


## Building new docker image

`hack/build/build.sh` script helps building docker image. It requires you to have
`IMAGE` environment variable set (e.g. `IMAGE=bitesize-registry.default.svc.cluster.local:5000/core/environment-operator`).
By default, image will be tagged with the HEAD commit tag. If you want to override
it and release versioned image, set `IMAGE_TAG` environment variable.

## Releasing new version

* Tag MASTER branch with release version (e.g. v0.0.1)
* Generate `CHANGELOG.md` by running `reportcopter`. Needs [reportcopter](https://github.com/3zcurdia/reportcopter) on local system.
* Build new docker image (see above).
