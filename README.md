# Environment operator [![dev](https://travis-ci.org/pearsontechnology/environment-operator.svg?branch=dev)](https://travis-ci.org/pearsontechnology/environment-operator/branches)

![environmentoperatoricon](https://github.com/pearsontechnology/environment-operator/blob/dev/docs/images/environmentoperatoricon.png)

The purpose of Environment Operator is to provide a seamless application deployment capability for a given environment within Kubernetes. It can easily hook into your existing CI/CD pipeline capabilities by installing our [Environment Operator Jenkins plugin](https://github.com/pearsontechnology/environment-operator-jenkins-plugin) to interface with environment operator and deploy your services.

Each environment (development, staging, production) has its own definition and a separate endpoint to perform deployments. Currently, environment operator supports Deployments, Services, Ingresses, MongoDB Statefulsets, and HorizonPodAutoscalers.

In order to begin deploying mircorservices through environment operator, you will need to start with the [Operations Guide](https://github.com/pearsontechnology/environment-operator/blob/dev/docs/Operatonal_Guide.md). The operations guide will provide the details required to get environment-operator itself deployed to a namespace and ready to manage your environment. Once environment operator is ready for use in your Kubernetes namespace, users of Environment Operator should start with our [User Guide](https://github.com/pearsontechnology/environment-operator/blob/dev/docs/User_Guide.md) to deploy their microservices.

Additionally, for those interested in developing against Environment Operator, check out our [Builder Guide](https://github.com/pearsontechnology/environment-operator/blob/dev/docs/Build.md).

*******************

### Just Show Me How To Run It...  

For those that would rather get an example running and then go back to read the docs on how to further configure environment operator, the [quick start guide](https://github.com/pearsontechnology/environment-operator/blob/dev/docs/Getting_Started.md) is for you...

*******************

![workflow](https://github.com/pearsontechnology/environment-operator/blob/dev/docs/images/workflow.png)

*******************

**Other documentation on Environment Operator:**  

* [Using Docker Registries (Dockerhub, Google Container Registry)](https://github.com/pearsontechnology/environment-operator/blob/dev/docs/Private_Registry.md)
* [Deploying a Mongo Statefulset](https://github.com/pearsontechnology/environment-operator/blob/dev/docs/Mongo.md)
