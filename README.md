# Kubernetes Operator for SPIRE

The [Kubernetes Operator](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/) for SPIRE configures, deploys, and helps ensure that a SPIRE server and agents are up and running in a Kubernetes cluster based on basic user-defined specifications. 

This is a proof-of-concept project by the interns under the Identity and Access Management of GreenLake Platform (summer 2023). 

## How it Works

### Custom Resources

#### SPIRE Server

The [SPIRE Server](docs/spireserver-crd.md) resource is a CRD that represents a SPIRE server as an individual Kubernetes resource. 

#### SPIRE Agent

The [SPIRE Agent](docs/spireagent-crd.md) resource is a CRD that represents a SPIRE agent as an individual Kubernetes resource. 

### Configuring and Installing a SPIRE Server

The controller listens for the creation of a resource of type SPIRE Server for its reconciliation logic to be triggered. The user must create their own configuration for a SPIRE server in a yaml file for a resource of kind `SpireServer`. The user can run the command `kubectl apply -f <yaml-file-name>` to trigger the controller. Based on the specifications in the user-inputted yaml file for a SPIRE Server instance, customized Kubernetes resources (such as `ConfigMap`, `StatefulSet`, `Service`, etc.) are generated and deployed in the Kubernetes cluster. 

### Health Checks (SPIRE Server)

Once all server-related components are deployed, the controller constantly runs a health check in the background by assessing the conditions of the SPIRE server pods deployed by the operator. The health status of the SPIRE Server is updated every 5 seconds and can be viewed by running `kubectl get spireservers`. 

### Configuring and Installing a SPIRE Agent

Once the SPIRE server is in a "READY" health state, SPIRE agents can be deployed. The controller listens for the creation of a resource of type SPIRE Agent for its reconciliation logic to be triggered. The user must create their own configuration for a SPIRE agent in a yaml file for a resource of kind `SpireAgent`. The user can run the command `kubectl apply -f <yaml-file-name>` to trigger the controller. Based on the specifications in the user-inputted yaml file for a SPIRE Agent instance, customized Kubernetes resources (such as `ConfigMap`, `DaemonSet`, etc.) are generated and deployed in the Kubernetes cluster. 

### Running the Operator
The operator is designed to control/manage the same Kubernetes cluster where the SPIRE components will be deployed. 

---

## Documentation

- [Getting Started Guide](docs/getting-started.md)
- [SPIRE Server CRD Configuration Reference](docs/spireserver-crd.md)
- [SPIRE Agent CRD Configuration Reference](docs/spireagent-crd.md)
- [Design Document](https://docs.google.com/document/d/1F7h9khGMh2wz6tED40TXQH3wUlLYr-6FEt-Cukk3MnA/edit?usp=sharing)