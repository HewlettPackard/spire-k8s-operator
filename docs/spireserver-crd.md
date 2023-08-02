# SpireServer Custom Resource Definition

The SpireServer Custom Resource Definition (CRD) is a cluster-wide resource that represents a SPIRE server as a Kubernetes resource. 

When an instance of the CRD is created, the controller creates the necessary components for a user-configured SPIRE server to function inside a Kubernetes cluster.  

The definition can be found [here](../api/v1/spireserver_types.go).

## SpireServerSpec
| Field | Required | Description |
| ----- | -------- | ----------- |
| `trustDomain`         | REQUIRED | Trust domain associated with the SPIRE server |
| `port`                | REQUIRED | Port on which the SPIRE server listens to agents |
| `nodeAttestors`       | REQUIRED | Node attestor plugins the SPIRE server uses |
| `keyStorage` | REQUIRED | Indicates whether the generated keys are stored on disk or in memory |
| `replicas` | REQUIRED | Number of replicas for SPIRE server |
| `dataStore` | REQUIRED | Indicates how server data should be stored (`sqlite3`, `mysql`, `postgres`) |
| `connectionString` | REQUIRED | Connection string for the datastore |

## SpireServerStatus
 Field | Description |
| ----- | ----------- |
| `health` | Indicates whether the SPIRE server is in an error state (`ERROR`), initializing (`INIT`), live (`LIVE`), or ready (`READY`) |

## Examples
1. SPIRE Server from [SPIRE's Quickstart for Kubernetes](https://spiffe.io/docs/latest/try/getting-started-k8s/)

    ```yaml
    apiVersion: spire.hpe.com/v1
    kind: SpireServer
    metadata:
        name: spire-server-01
    spec:
        trustDomain: example.org
        port: 8081
        nodeAttestors: 
            - k8s_sat
        keyStorage: disk
        replicas: 1
    ```

## Note
Under the High Availability (HA) model, if your cluster has more than one replica of a SPIRE Server, it cannot use `sqlite3` as its datastore. The operator will reject and delete any SPIRE server instances with this configuration. 