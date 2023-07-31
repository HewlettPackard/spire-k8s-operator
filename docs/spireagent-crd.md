# SpireAgent Custom Resource Definition

The SpireAgent Custom Resource Definition (CRD) is a cluster-wide resource that represents a SPIRE agent as a Kubernetes resource. 

When an instance of the CRD is created, the controller creates the necessary components for a user-configured SPIRE agent to function inside a Kubernetes cluster.  

The definition can be found [here](../api/v1/spireagent_types.go).

## SpireAgentSpec
| Field | Required | Description |
| ----- | -------- | ----------- |
| `trustDomain`         | REQUIRED | Trust domain that the SPIRE agent issues identities to |
| `nodeAttestor`       | REQUIRED | Node attestor plugin the SPIRE agent uses |
| `workloadAttestors` | REQUIRED | Workload attestor plugins the SPIRE agent uses |
| `keyStorage` | REQUIRED | Indicates whether the generated keys are stored on disk or in memory |
| `serverPort` | REQUIRED | Port on which the SPIRE server listens to agents |

## Examples
1. SPIRE Agent from [SPIRE's Quickstart for Kubernetes](https://spiffe.io/docs/latest/try/getting-started-k8s/)

    ```yaml
    apiVersion: spire.hpe.com/v1
    kind: SpireAgent
    metadata:
        name: spire-agent-01
    spec:
        trustDomain: example.org
        nodeAttestor: k8s_sat
        workloadAttestors: 
            - k8s
            - unix
        keyStorage: memory
        serverPort: 8081
    ```