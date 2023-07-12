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

3. Deploy the sample CRD. 
```bash
kubectl apply -f config/samples
```

4. You can test whether the CRD is deployed correctly with the following command. 
```bash
kubectl get spireservers
```

## Note
In line with conventional practices for SPIRE, it is not recommended to run more than 1 SPIRE server instance in a cluster. This also conflicts with the controller's logic for health checks. Nonetheless, you can have replicas of a SPIRE server instance corresponding to the same trust domain under the HA model with this controller. 