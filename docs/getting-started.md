# Getting Started with SPIRE Operator for Kubernetes

This guide will walk you through the process of setting up the SPIRE operator for Kubernetes and deploying a user-configured SPIRE Server using the operator. 

## Prerequisites
Before you begin, you should have a Kubernetes cluster running and access to the `kubectl` command line tool to control the cluster. 

You can use tools like `kind` or `minikube` to set up a local Kubernetes cluster. With `kind`, you can use the `kind create cluster` command to start the cluster and `kind delete cluster` command to delete the cluster. 

## Installing the CRD and Operator
1. Clone this repository. 
2. Run the following commands to install the custom resource onto the cluster and start the controller. 
```bash
make manifests
```

```bash
make install
```


```bash
make run
```

3. Deploy the sample CRD yaml. 
```bash
kubectl apply -f config/samples/spire-server-01.yaml
```

4. You can test whether the CRD is deployed correctly with the following command. 
```bash
kubectl get spireservers
```

5. To check whether all related resources have been deployed correctly, you can run the following commands. 
```bash
kubectl get statefulset
kubectl get pods
kubectl get services
```

## Note
In line with conventional practices for SPIRE, it is not recommended to run more than 1 SPIRE server instance in a cluster. This also conflicts with the controller's logic for health checks. Nonetheless, you can have replicas of a SPIRE server instance corresponding to the same trust domain under the HA model with this controller. 