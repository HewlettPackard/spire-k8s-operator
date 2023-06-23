# SpireServer Custom Resource Definition

The SpireServer Custom Resource Definition (CRD) is a cluster-wide resource that represents a SPIRE server as a Kubernetes resource. 

When an instance of the CRD is created, the controller creates the necessary components for a user-configured SPIRE server to function inside a Kubernetes cluster.  

The definition can be found [here](../api/v1/spireserver_types.go).

## SpireServerSpec
| Field | Required | Description |
| ----- | -------- | ----------- |
| `name`                | REQUIRED | The name of the SPIRE server |
| `trustDomain`         | REQUIRED | The trust domain associated with the SPIRE server |
| `port`                | REQUIRED | The port on which the SPIRE server Listens to agents |
| `nodeAttestors`       | REQUIRED | The node attestor plugins the SPIRE server uses |
| `keyStorage` | REQUIRED | Indicates whether the generated keys are stored on disk or in memory |

## SpireServerStatus
 Field | Description |
| ----- | ----------- |
| `health` | Indicates whether the SPIRE server is live or ready |

## Examples
1. SPIRE Server from [SPIRE's Quickstart for Kubernetes](https://spiffe.io/docs/latest/try/getting-started-k8s/)

    ```yaml
    apiVersion: spire.hpe.com/v1
    kind: SpireServer
    metadata:
        name: spire-server-01
    spec:
        name: spire-server-01
        trustDomain: example.org
        port: 8081
        nodeAttestors: 
            - k8s_sat
        keyStorage: disk
    ```