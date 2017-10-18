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

The statefulset, when deployed by environment operator, will spin up a specified number of mongo replicas into your namespace with each pod containing a MongoDB container. 

Prior to deploying (adding the service to your environments.bitesize file as shown below in the example) your mongo service, a secret needs to be created in your k8s namespace to allow bootstrapping 
of your mongo cluster with secure client communcation.  This will eventually will be automated, but the following must be executed before you add your mongo service to your bitesize file.

```
TMPFILE=$(mktemp)
/usr/bin/openssl rand -base64 741 | tr -d "=+/" > $TMPFILE
kubectl create secret generic shared-bootstrap-data --from-file=internal-auth-mongodb-keyfile=$TMPFILE --namespace=$YOURNAMESPACE
rm $TMPFILE

```

Once you have environment-operator deployed into a namespace in your cluster, you may specify a new Mongo service in your environments.bitesize file
to deploy the mongo cluster. 

**Example environments.bitesize**

```
project: somogyi-app
environments:
  - name: dev
    namespace: somogyi-app
    services:
      - name: mongodb
        database_type: mongo
        version: 3.4
        replicas: 3
        port: 27017
        volumes:
          - name: mongo-persistent-storage
            path: /data/db
            modes: ReadWriteOnce
            size: 10G
```

Unlike a normal deployment as shown in the [User Guide](https://github.com/pearsontechnology/environment-operator/blob/dev/docs/User_Guide.md), 
when a "database_type" (mongo is currently the only type supported)is specified, environment-operator will know this is 
a database service and that it should deploy it as a statefulset and also create a Headless service for that service. The
[headless service](https://kubernetes.io/docs/concepts/services-networking/service/#headless-services) will allow you to connect to the mongo database using a stable DNS name (more on that to follow).

Parameters in the environments.bitesize file are used as follows when specifying a mongo database service:

- **database_type**: As noted before, "mongo" is the only type supported. We may add other DB statefulsets in the future.
- **version**: This is the version of [mongo](https://hub.docker.com/_/mongo/) that will be deployed. 
- **replicas**: Number of mongo pods you would like. A replica value of 3 will provision a Primary and two Secondaries. This value can be changed to scale your mongo cluster up and down.
- **port**: The port your mongo service will accept requests on. Typically this is 27017 for a mongo database.
- **volumes**: This is the size of the EBS volume that will be dynamically provisioned and mounted into your Mongo containers. 
You'll most likely only want to configure the size and leave the other volume paramters (name, path, mode) as shown in the example.

Once your mongo service is deployed, you'll want to configure the initial replica set and add an administrator through the local host exception. This piece will eventually be 
automated as well. Utilize the following commands to initialize the cluster (note that you will need to change the namespace and secure hostnames depending on how you named your service in 
 the environments.bitesize file.  The example below assumes you named it "mongodb" and deployed it into the "somogyi-app" namespace with a value of "3" for replicas):

```
#Jump on a mongo pod
kubectl exec -it mongodb-0 --namespace=somogyi-app -c mongo bash

#Enter the mongo shell via localhost exception
mongo

#Initiate the replicaset
rs.initiate(
       {_id: "mongo", version: 1,members: [
       { _id: 0, host : "mongodb-0.mongodb.somogyi-app.svc.cluster.local:27017" },
       { _id: 1, host : "mongodb-1.mongodb.somogyi-app.svc.cluster.local:27017" },
       { _id: 2, host : "mongodb-2.mongodb.somogyi-app.svc.cluster.local:27017" },
       ],
       settings: {
          chainingAllowed : true, 
          heartbeatIntervalMillis : 2000,
          heartbeatTimeoutSecs: 10,
          electionTimeoutMillis : 10000,
          catchUpTimeoutMillis : 60000,
          getLastErrorModes : {},
          getLastErrorDefaults : {w: 'majority', wtimeout: 5000 }
       }}
);

#Ensure all members are in the replicaset
rs.status()

#Add a user to the admin database
db.getSiblingDB("admin").createUser({
      user : "admin",
      pwd  : "password",
      roles: [ { role: "root", db: "admin" } ]
 });

#Gain authorization as an administrator to add another user
db.getSiblingDB('admin').auth("admin", "password");
 
#Add an appication user to a database
db.createUser({
      user : "user",
      pwd  : "abc123",
      roles : [ { "role" : "readWrite", "db" : "nodeDB" } ]
  });

```

You may now connect to it and add some data. Below is a nodejs implementation that 
could be used built into a docker container, and deployed to your namespace alongside the mongo statefulset via environment-operator.
I'm leaving off those details as they are already covered in the User Guide. However, this code snippet will show you the 
proper syntax for a connection string to the mongo replicaset from a nodeJS app. Note that the connection url utilizes the stable DNS names 
for a 3 pod statefulset and that you have to specify a replicaset name, which will always be called "mongodb".  You do not need to specify all replicas
in the connection string.

```
var MongoClient = require('mongodb').MongoClient;
var assert = require('assert');
var randomstring = require('randomstring')
var url = 'mongodb://user:abc123@mongodb-0.mongodb,mongodb-1.mongodb,mongodb-2.mongodb:27017/nodeDB?replicaSet=mongodb'

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






