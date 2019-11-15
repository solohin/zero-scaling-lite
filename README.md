# Zero scaler for Kubernetes
Simple, straightforward, lightweight. Works on metrics from Nginx Ingress. Depends on Prometheus to retrieve the metrics. 

## How it works 

### Proxy component
On request received, holds it, asks scaler for scaling up to 1 replica.  

### Downscaler
Every minute:
Checks ingress with `zero-scaling/enabled = "true"` annotation every 60 seconds. If ingress has traffic after `zero-scaling/sleep-after` seconds, scales down deployment `zero-scaling\deployment_name` to zero and redirect all traffic to proxy. Original proxy value saved in 

### Upscaler
HTTP /wakeup/[domainName]
Rescales the deployment and sends traffic to `zero-scaling\service_name`.