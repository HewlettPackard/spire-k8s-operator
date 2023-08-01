# Getting Started with SPIRE Operator for Kubernetes

This guide will walk you through the process of setting up the SPIRE operator for Kubernetes and deploying a user-configured SPIRE server and agent using the operator. It will then guide you through configuring a registration entry for a workload and fetching an x509-SVID over the SPIFFE Workload API, as adapted from the [SPIRE Quickstart for Kubernetes](https://spiffe.io/docs/latest/try/getting-started-k8s/). 

## Prerequisites
Before you begin, you should have a Kubernetes cluster running and access to the `kubectl` command line tool to control the cluster. 

You can use tools like `kind` or `minikube` to set up a local Kubernetes cluster. With `kind`, you can use the `kind create cluster` command to start the cluster and `kind delete cluster` command to delete the cluster. 

## Installing the CRDs and Operator
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

## Deploying a SPIRE Server Instance
3. In a separate terminal window, deploy the sample server yaml. 
```bash
kubectl apply -f config/samples/spire_server_sample.yaml
```

4. You can test whether the SPIRE server is deployed correctly with the following command. Running this command will also show the current health of the SPIRE server based on the operator's health check system. A "READY" health state means that the operator can now also install and configure a SPIRE agent for this server. 
```bash
kubectl get spireservers
```

5. To check whether all related resources have been deployed correctly, you can run the following commands. 
```bash
kubectl get statefulset
kubectl get pods
kubectl get services
```

## Deploying a SPIRE Agent Instance
6. Deploy the sample agent yaml. 
```bash
kubectl apply -f config/samples/spire_agent_sample.yaml
```

7. You can test whether the SPIRE agent is deployed correctly with the following command. 
```bash
kubectl get spireagents
```

8. To check whether all related resources have been deployed correctly, you can run the following commands. 
```bash
kubectl get daemonset
kubectl get pods
```

## Configuring Registration Entries
9. Using the following command, create a new registration entry for the node. This also specifies the SPIFFE ID to allocate to the node. 
```bash
kubectl exec spire-server-0 -- \
  /opt/spire/bin/spire-server entry create \
  -spiffeID spiffe://example.org/ns/default/sa/spire-agent \
  -selector k8s_sat:cluster:demo-cluster \
  -selector k8s_sat:agent_ns:default \
  -selector k8s_sat:agent_sa:spire-agent \
  -node
```

10. Using the following command, create a new registration entry for the workload. This also specifies the SPIFFE ID to allocate to the workload.
```bash
kubectl exec spire-server-0 -- \
  /opt/spire/bin/spire-server entry create \
  -spiffeID spiffe://example.org/ns/default/sa/default \
  -parentID spiffe://example.org/ns/default/sa/spire-agent \
  -selector k8s:ns:default \
  -selector k8s:sa:default
``` 

## Fetching an x509-SVID
11. To access SPIRE (or the Workload API UNIX domain socket), a workload container must be configured. The deployment file in the following command configures a no-op container using the `spire-k8s` docker image for the server and agent. 
```bash
kubectl apply -f config/samples/client-deployment.yaml
```

12. To verify that the container can access the socket, use the following command. If the agent is running, you will see a list of SVIDs, and if it is not, you will see an error message (such as “no such file or directory” or “connection refused”). 
```bash
kubectl exec -it $(kubectl get pods -o=jsonpath='{.items[0].metadata.name}' \
  -l app=client)  -- /opt/spire/bin/spire-agent api fetch -socketPath /run/spire/sockets/agent.sock
```

## Note
In line with conventional practices for SPIRE, it is not recommended to run more than 1 SPIRE server instance in a cluster. This also conflicts with the controller's logic for health checks. Nonetheless, you can have replicas of a SPIRE server instance corresponding to the same trust domain under the HA model with this controller. 