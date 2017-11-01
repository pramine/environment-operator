# Deploying a Mongo Statefulset 

Environment operator supports the ability to deploy a Mongo database as a Statefulset into a namespace of your choosing within your kubernetes cluster. This guide
will provide the details of now to add a mongo database to your environments.bitesize file. If you are not familiar with the environments.bitesize file,
please go back and start with our [User Guide](https://github.com/pearsontechnology/environment-operator/blob/dev/docs/User_Guide.md)


Kubernetes Prerequiste:

- To utilize this feature of environment-operator, you will need to be running k8s 1.5+ in your cluster to ensure Statefulsets are supported. The current release of environment operator is
tested on k8s 1.5.7
- You will also need to ensure your k8s deployment has EBS cloud provider support enabled.  The Mongo Statefulset that gets deployed needs to have a storage
 class defined that is used to dynamically provision volumes for the Mongo database containers.  More information on how to setup cloud-provider support and dynamic provisioning  
 may be found [here](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#aws)

The mongo statefulset, when deployed by environment operator, will spin up a specified number of mongo replicas into your namespace with each pod containing a "mongo" container. 
To deploy mongo, get environment-operator deployed into a namespace within your cluster, then specify a new Mongo service in your environments.bitesize file
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
when a "database_type" (mongo is currently the only type supported) is specified, environment-operator will know it is 
a database service and that it should  be deployed as a statefulset with an accompanying Headless service. The
[headless service](https://kubernetes.io/docs/concepts/services-networking/service/#headless-services) will allow you to 
connect to the mongo database using a stable DNS name (more on that to follow).

Parameters in the environments.bitesize file are used as follows when specifying a mongo database service:

- **database_type**: As noted before, "mongo" is the only type supported. We may add other DB statefulsets in the future.
- **version**: This is the version of [mongo](https://hub.docker.com/_/mongo/) that will be deployed. 
- **replicas**: Number of mongo pods you would like. A replica value of 3 will provision a Primary and two Secondaries. This value can be changed to scale your mongo cluster up and down.
- **port**: The port your mongo service will accept requests on. Typically this is 27017 for a mongo database.
- **volumes**: This is the size of the EBS volume that will be dynamically provisioned and mounted into your Mongo containers. 
You'll most likely only want to configure the size and leave the other volume paramters (name, path, mode) as shown in the example.

Once your mongo service is deployed, you'll want to configure the initial replica set and add an administrator through the local host exception. This piece will eventually be 
automated as well, but should be handled by your DBA (this is an example). Utilize the following commands to initialize the cluster (note that you will need to change the namespace 
and secure hostnames depending on how you named your service in the environments.bitesize file).  The example below assumes you named it "mongodb" and deployed it into the 
"somogyi-app" namespace with a value of "3" for replicas:

```
#Jump on a mongo pod within your cluster
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

#Create a DB
use nodeDB

#Add an appication user to a database
db.createUser({
      user : "user",
      pwd  : "abc123",
      roles : [ { "role" : "readWrite", "db" : "nodeDB" },
                { role: "dbAdmin", db: "nodeDB" } ]
  });
  
# If you scale your replicas by changing replicas in your environments.bitesize file, the new replicas will generate in the statefulset, but
#you still need to add them to the replicaset until automation is put in place  ex:
rs.add("mongodb-3.mongodb.somogyi-app.svc.cluster.local:27017"")

```

You may now connect to the DB and add some data. Below is a nodejs implementation that 
could be used as an example,  which could be deployed to your namespace alongside the mongo statefulset via environment-operator.
I'm leaving off the details with building an App as they are already covered in the User Guide. However, this code snippet will show you the 
proper syntax for a connection string to the mongo replicaset from a nodeJS app. Also, note that the connection url utilizes the stable DNS names 
for a 3 pod statefulset and that you have to specify a replicaset name, which will always be called "mongodb".  You do not need to specify all replicas
in the connection string.

```

var async = require("async");
var MongoClient = require('mongodb').MongoClient;
var assert = require('assert');
var randomstring = require('randomstring')

//Delay of 1 second
var delay = 1000
var uniqueNumber = 1;
var url = 'mongodb://user:abc123@mongodb-0.mongodb,mongodb-1.mongodb,mongodb-2.mongodb:27017/nodeDB?replicaSet=mongo'

async.forever( function(next) {
     MongoClient.connect(url, function (err, db) {
         assert.equal(null, err);
         db.collection('myNodeCollection').insertOne({

                "name": randomstring.generate(12),
                "address": randomstring.generate(25),
                "city": randomstring.generate(20),
                "state": randomstring.generate(2),
                "date": new Date(),
                "favoriteSports" : [ randomstring.generate(10), randomstring.generate(10)],
                "record" : uniqueNumber
         }, function (err, result) {
                assert.equal(err, null);
                console.log("Inserted document [ record : " + uniqueNumber + " ] into the myNodeCollection.");
                uniqueNumber++;
                setTimeout(function() {
                    next();
                }, delay)
         });
         db.close();
     });
},
function() {
});



```






