# Roadmap
* Support for SPIRE Agent: install and configure SPIRE Agents
* Support for SPIRE Controller Manager: install and configure the SPIRE Controller Manager (managing registration entries would still be out of scope)

## Recently Completed
* PoC goal has been completed
* Support Data Store Pluging SQLite (sqlite3)
    * Save Database in File
    * Save Database in Memory


## Near-Term and Medium-Term
* Spire Agent Support
* Support Spire Server Data Store Pluging...
    * MySQL
    * Postgres

## Long-Term
* Full customization of a SPIRE Server and Agent within the same trust domain.
* Establish communication between a SPIRE Server and Agent by viewing SVIDs 

### Initial Proof of Concept (PoC)

- **Status**: Completed
- **Goal**: The goal of this PoC aims to develop and test a Kubernetes operator and CRDs for deploying the SPIRE server in Kubernetes environments. The operator will include automation for common tasks such as installing, configuring, as well as making sure the SPIRE Server is up and running by running health checks.