# Getting Started with SPIRE Operator for Kubernetes

This guide will walk you through the process of setting up the SPIRE operator for Kubernetes and deploying a user-configured SPIRE Server using the operator. 

## Prerequisites
Before you begin, you should have a Kubernetes cluster running and access to the `kubectl` command line tool to control the cluster. 

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
kubectl apply -f spire-server-01.yaml
```

4. You can test whether the CRD is deployed correctly with the following command. 
```bash
kubectl get spireservers
```