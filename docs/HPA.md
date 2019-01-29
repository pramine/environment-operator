# Using Horizontal Pod Autoscaling

Horizontal Pod Autoscaling is a native Kubernetes feature. Horizontal Pod Autoscaling (HPA), allows a developer to dynamically scale the number of pods running in an application depending on CPU utilization, memory or other custom metrics.
Environment operator currently supports scaling the number of pods based on CPU utilization. The ability to dynamically scale the number of pods based on memory or custom metrics is currently on the Environment Operator roadmap.

The following  bitesize file shows an example HPA configuration:


**Example environments.bitesize**

```
project: pidah-app
environments:
  - name: dev
    namespace: pidah-app
    deployment:
      method: rolling-upgrade
    services:
      - name: api
        port: 80
        hpa:
          min_replicas: 1
          max_replicas: 5
          target_cpu_utilization_percentage: 80
        env:
          - name: API_PORT
            value: 80
        requests:
           cpu: 100m
```

You specify the minimum and maximum number of pod replicas required by the application and the threshold – target_cpu_utilization_percentage which would trigger a scale up of the replica count. The target_cpu_utilization_percentage is a weighted average across the available number of replicas. Once utilization drops below 80% across the pods, the HPA controller would again dynamically scale down the number of pods in the application.

*NOTE*
It is important to specify the `requests` field above as HPA requires this configuration. Without it, HPA will not work correctly and the pods will not scale dynamically.


More details on HPA is available [here](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/)


