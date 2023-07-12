# Kubernetes Operator for SPIRE
---

The [Kubernetes Operator](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/) for SPIRE configures, deploys, and helps ensure that a SPIRE server is up and running in a Kubernetes cluster based on a basic user-defined specification. 

This is a proof-of-concept project by the interns under the Identity and Access Management of Greenlake Cloud Platform (summer 2023). 

## How it Works

### Custom Resources

#### SPIRE Server

The [SPIRE Server](docs/spireserver-crd.md) resource is a CRD that represents a SPIRE Server as an individual Kubernetes resource. 

#### Configuring and Installing a SPIRE Server

The controller listens for the creation of a resource of type SPIRE Server for its reconciliation logic to be triggered. Based on the specifications in the user-inputted yaml file for a SPIRE Server instance, customized Kubernetes resources (such as `ConfigMap`, `StatefulSet`, `Service`, etc.) are generated and deployed in the Kubernetes cluster. The operator supports the High Availability model if more than 1 replica is specified in the user-inputted yaml file. 

#### Health Checks

Once all server-related components are deployed, the controller constantly runs a health check in the background by assessing the conditions of the SPIRE server pods deployed by the operator. The health status of the SPIRE Server is updated every 5 seconds and can be viewed by running `kubectl get spireservers`. 

#### Deploying the Operator
The operator is designed to be deployed in the same Kubernetes cluster as where the user intends on hosting the SPIRE servers. 

---

## Documentation

- [Getting Started Guide](docs/getting_started.md)
- [SPIRE Server CRD Configuration Reference](docs/spireserver-crd.md)
- [Design Document](https://docs.google.com/document/d/1F7h9khGMh2wz6tED40TXQH3wUlLYr-6FEt-Cukk3MnA/edit?usp=sharing)