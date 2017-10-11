# Deploying a Mongo Statefulset 

Environment operator supports the ability to deploy a Mongo database as a Statefulset into a namespace of your choosing within your kubernetes cluster. This guide
will provide the details of now to add a mongo database to your environments.bitesize file. If you are not familiar with the environments.bitesize file,
please go back and start with our [User Guide](https://github.com/pearsontechnology/environment-operator/blob/dev/docs/User_Guide.md)


Kubernetes Prerequiste:

- Additionally you will need to be running k8s 1.5+ in your cluster to ensure Statefulsets are supported. The current release of environment operator is
tested on k8s 1.5.7
- In order to make use of this feature, need to ensure your k8s deployment has EBS cloud provider support enabled as the Mongo Statefulset needs to have a storage
 class defined that is used to dynamically provision volumes for the Mongo database containers.  More information on how to do this 
 may be found [here](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#aws)

The statefulset, when deployed by environment operator, will spin up a specified number of mongo replicas into your namespace with each pod containing both a MongoDB container
as well as a [Mongo Sidecar](https://github.com/pearsontechnology/mongo-k8s-sidecar), which is used to configure the Mongo cluster and maintain quorum should loss of a node/pod occur. 


Once you have environment-operator deployed into a namespace in your cluster, you may specify a new Mongo service in your environments.bitesize file
to deploy the mongo cluster. 

**Example environments.bitesize**

```
project: somogyi-app
environments:
  - name: dev
    namespace: somogyi-app
    deployment:
      method: rolling-upgrade
    services:
      - name: mongodb
        database_type: mongo
        version: 3.4
        graceperiod: 10
        replicas: 3
        port: 27017
        volumes:
          - name: mongo-persistent-storage
            path: /data/db
            modes: ReadWriteOnce
            size: 10G
        env:
          - name: MONGO_PORT
            value: 27017
          - name: KUBE_NAMESPACE
            value: 'somogyi-app'
          - name: MONGO_SIDECAR_POD_LABELS
            value: "role=mongo,name=mongodb"
          - name: KUBERNETES_MONGO_SERVICE_NAME
            value: "mongodb"
          - name: MONGO_DB_USERNAME
            value: 'username'
          - name: MONGO_DB_PASSWORD
            value: 'password'

```

Unlike a normal deployment as shown in the [User Guide](https://github.com/pearsontechnology/environment-operator/blob/dev/docs/User_Guide.md), 
when a "database_type" (mongo is currently the only type supported)is specified, environment-operator will know this is 
a database service and that it should deploy it as a statefulset and also create a Headless service for that service. The
[headless service](https://kubernetes.io/docs/concepts/services-networking/service/#headless-services) will allow you to connect to the mongo database using a stable DNS name (more on that to follow).

Other parameters in the environments.bitesize file are used as follows when specifying a mongo database service:

- **database_type**: As noted before, "mongo" is the only type supported. We may add other DB statefulsets in the future.
- **version**: This is the version of [mongo](https://hub.docker.com/_/mongo/) that will be deployed. Note: the [mongo sidecar](https://github.com/pearsontechnology/mongo-k8s-sidecar)), which is a nodejs app
utilized to configure the mongo cluster has only been tested with version 3.2 and 3.4 of mongo.
- **graceperiod**: Typically this should be set to 10sec as shown in the example. This ensures mongo pods are terminated gracefully.
- **replicas**: Number of mongo pods you would like. A replica value of 3 will provision a Primary and two Secondaries.
- **port**: The port your mongo service will accept requests on. Typically this is 27017 for a mongo database.
- **volumes**: This is the size of the EBS volume that will be dynamically provisioned and mounted into your Mongo containers. 
You'll most likely only want to configure the size and leave the other volume paramters (name, path, mode) as shown in the example.
- **env**: These are variables that will be utilized by the mongo-sidecar. Information on other parameters that can be used to
configure your cluster are described in the sidecar docs, shown [here](https://github.com/pearsontechnology/mongo-k8s-sidecar)). 
Note: The MONGO_PORT should match the port specified for your service, MONGO_SIDECAR_POD_LABELS should contain the same "name" as the name
for the environments.bitesize service. Likewise, the KUBERNETES_MONGO_SERVICE_NAME should match the environments.bitesize service name. This is to ensure
the sidecar properly configures the mongo service and finds the pods that are part of the mongo statefulset.

Once your mongo service is deployed, you may now connect to it and add some data. Below is a nodejs implementation that 
could be used built into a docker container, and deployed to your namespace alongside the mongo statefulset via environment-operator.
I'm leaving off those details as they are already covered in the User Guide. However, this code snippet will show you the 
proper syntax for a connection string to the mongo replicaset. Note that the connection url utilizes the stable DNS names 
for a 3 pod statefulset and that you have to specify a replicaset name, which will always be called "mongodb".

```
var MongoClient = require('mongodb').MongoClient;
var assert = require('assert');
var randomstring = require('randomstring')
var url = 'mongodb://username:password@mongodb-0.mongodb,mongodb-1.mongodb,mongodb-2.mongodb:27017/nodeDB?replicaSet=mongodb'

MongoClient.connect(url, function (err, db) {
    assert.equal(null, err);

    console.log("Connected successfully to server!");

    insertDocuments(db, function () {
        db.close();
    });
});

var insertDocuments = function (db, callback) {

    for( var i = 0; i < 10000; i++) {

        db.collection('myNodeCollection2').insertOne({

            "name": randomstring.generate(12),
            "address": randomstring.generate(25),
            "city": randomstring.generate(20),
            "state": randomstring.generate(2),
            "date": new Date(),
            "favoriteSports" : [ randomstring.generate(10), randomstring.generate(10)],
            "record" : i
        }, function (err, result) {
            assert.equal(err, null);
            console.log("Inserted document into the myNodeCollection.");
            callback();


        });
    }
};



```






