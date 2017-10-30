# Environment operator operational guide.

Environment operator runs as a deployment in the namespace it manages. Example deployment (with local keycloak setup) can be found in [example](https://github.com/pearsontechnology/environment-operator/tree/master/example) directory in code repository.

*N.B.* This deployment depends on Keycloak. To use the deployment with static keys, you can ignore `OIDC_ISSUER_URL` and `OIDC_ALLOWED_GROUPS` environment variables in deployment, and setup `AUTH_TOKEN_FILE` environment variable. Recommended way to provide `AUTH_TOKEN_FILE` is via secrets, similar to the method for `GIT_PRIVATE_KEY` (see below).

## Environment variable list

* `GIT_REMOTE_REPOSITORY` - specifies remote repository, where your manifest/`environments.bitesize` file is located.
* `GIT_BRANCH` - specifies what branch to checkout from the GIT_REMOTE_REPOSITORY. If ommitted this defaults to "master"
* `GIT_PRIVATE_KEY` - git private key, used to authenticate against `GIT_REMOTE_REPOSITORY`. Must allow read-only access.
* `BITESIZE_FILE` - usually `environments.bitesize`, but can be anything, to suit project's needs better (for example, you can have file per environment, or per kubernetes cluster).
* `ENVIRONMENT_NAME` - corresponds to the "name" field in the manifest/environments.bitesize file. This is the environment that operator manages.
* `DOCKER_REGISTRY` - registry to download application images from.
* `DOCKER_PULL_SECRETS` - A comma delimited list of k8s secret names in your applications k8s namespace that will be used to pull images from your private registy. See [private registry](https://github.com/pearsontechnology/environment-operator/blob/dev/docs/Private_Registry.md) documentation for how to use private registries.
* `PROJECT`  - used for metadata (e.g. tags for managed services). 
* `OIDC_ISSUER_URL` - issuer ID for OpenID Connect.
* `OIDC_ALLOWED_GROUPS` - comma separated list of Keycloak provided groups, that can perform HTTP actions against environment-operator.
* `DEBUG` - debug mode.
* `NAMESPACE` - namespace this environment-operator actions on. Usually self-referenced to local namespace.
* `AUTH_TOKEN_FILE` - path to a static auth token file. Usually injected into environment-operator via kubernetes secret.


## Using kubernetes secrets in environment operator

It is recommended that `GIT_PRIVATE_KEY` would be used as a reference to the secret. Create file named key with private key contents (e.g. cp ~/.ssh/id_rsa key) and create secret git-private-key from it:

```
$ kubectl create secret generic git-private-key --from-file=./key
```

Then you can refer to git-private-key in environment-operator deployment

```

        - name: GIT_PRIVATE_KEY
          valueFrom:
            secretKeyRef:
              name: git-private-key
              key: key

```

Similarly, you can create `deploy-auth-token-file` secret (if you are not using Keycloak) and use it in volume mounts:


```
$ kubectl create secret generic deploy-auth-token-file --from-file=./token
```

and use it as a volume:

```
     containers:
       name: environment-operator
       env:
         - name: AUTH_TOKEN_FILE
           value: /etc/auth/token
       ...
       volumeMounts:
         name: "auth-token"
         mountPath: "/etc/auth"
         readOnly: true
     ...
     volumes:
       name: "auth-token"
       secret:
         secretName: deploy-auth-token-file
```

## Private registry support

The environment operator allows Docker images to be deployed into a Kubernetes namespace from private registries like
DockerHub as well as Google Container Registry. This document details the process for configuring the environment operator
to use a private registry as well as how to establish Kubernetes secrets to allow the pods to pull from said registry.  For more information review the [Private Registry](https://github.com/pearsontechnology/environment-operator/blob/dev/docs/Private_Registry.md) documentation.

***************

## Deploy Sequence

![deploy-sequence](https://github.com/pearsontechnology/environment-operator/blob/bite-1788/docs/images/deploy-sequence.png)

***************


## Troubleshooting

This section describes common issues and troubleshooting steps when something does not work.

### Environment operator fails to clone git repository

Usual symptom is these lines in the system output log:

```
time="2017-07-27T11:38:34Z" level=error msg="Git clone error: Failed to authenticate SSH session: Waiting for USERAUTH response"
time="2017-07-27T11:38:34Z" level=error msg="Failed to resolve path '/tmp/repository': No such file or directory"
```

Second line indicates that repository clone has failed (could be wrong GIT_REMOTE_REPOSITORY or GIT_PRIVATE_KEY). However, first line indicates that there was user authentication error (most likely, the wrong GIT_PRIVATE_KEY set).

To verify this is the case, you will need to exec into container:

```
kubectl exec -ti $(kubectl get pods | awk '/environment-operator/{print $1}') /bin/bash
```

Then execute the following commands to verify git clone fails:

```
echo $GIT_PRIVATE_KEY > /tmp/key
chmod 0400 /tmp/key
export GIT_SSH_COMMAND="ssh -i /tmp/key"
git clone $GIT_REMOTE_POSITORY
```

This should give you the error environment-operator encounters.
