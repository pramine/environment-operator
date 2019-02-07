# Demo Autoscaling features on the PAAS

This demo showcases the autoscaling features of the PAAS.

I will use environment-operator to deploy an application to the cluster

environment-operator is running in the sample-app namespace:

```
root@master-i-01d298add324432ca:~/environment-operator/example [0]# kubectl get pod -n sample-app
NAME                                   READY     STATUS    RESTARTS   AGE
environment-operator-69564d97f-t6ndw   1/1       Running   0          10h
```

Deploy an app using curl to the environment-operator:

```
root@master-i-01d298add324432ca:~/environment-operator/example/sample_app [0]# curl -k  -XPOST -H 'Content-Type: application/json' -d '{"application":"sample-app-back", "name":"back", "version":"latest"}'  environment-operator.sample-app.svc.cluster.local/deploy
{"status":"deploying"}
```


Bitesize file used in this demo:

```
project: pearsontechnology
environments:
  - name: sample-app-environment
    namespace: sample-app
    deployment:
      method: rolling-upgrade
    services:
      - name: front
        external_url: front.sample-app.domain
        ssl: false
        port: 80
        env:
          - name: APP_PORT
            value: 80
          - name: BACK_END
            value: back.sample-app.svc.cluster.local
        requests:
           cpu: 10m
           memory: 5000Mi
        limits:
           cpu: 10m
           memory: 5000Mi
      - name: back
        port: 80
        replicas: 1
        hpa:
          min_replicas: 1
          max_replicas: 5
          target_cpu_utilization_percentage: 75
        env:
          - name: APP_PORT
            value: 80
        requests:
           cpu: 10m
           memory: 5000Mi
        limits:
           cpu: 50m
           memory: 5000Mi
```

*NOTE*
It is important to specify the `requests` field above as HPA requires this configuration. Without it, HPA will not work correctly and the pods will not scale dynamically.


Add load to the deployed app:
```
 ab -k -c 1000 -n 2000000 http://back.sample-app.svc.cluster.local/
```

Observe HPA scaling stats:
```
 kubectl get  hpa -n sample-app

 kubectl top pod -n sample-app

 kubectl get pod -n sample-app -w

 kubectl get nodes
```

Add a large pod (if needed to speed up demo):

```

cat nginx-pod.yaml

apiVersion: v1
kind: Pod 
metadata:
  name: test-autoscaler
spec:
  containers:
  - name: test-autoscaler
    image: nginx
    resources:
      limits:
        cpu: "3000m"
      requests:
        cpu: "3000m"
```

Observe cluster node scaling events:
```
 kubectl get events --all-namespaces  --field-selector reason=TriggeredScaleUp

 kubectl get nodes

 kubectl top nodes
```

Take a look at #spam channel on slack

More details on HPA is available [here](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/)


