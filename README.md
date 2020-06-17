# `kube-route-info`: kubectl tool

kube-route-info is a tool that provides the route information that has been configured on Ingresses or Services to access Pods.

## Usage

```sh
# View the route information of the service my-service
kubectl route-info service my-service

# View the route information of the ingress my-ingress in namespace my-namespace
kubectl route-info ingress my-ingress --namespace my-namespace
```

## Instalation

Work in Progress