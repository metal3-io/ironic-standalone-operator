# Metal3 Ironic Operator

Operator to maintain an Ironic deployment for Metal3.

## OpenShift notes

- Running a database requires the user to have `nonroot` or similar SCC.
- Running Ironic requires `nonroot` and `hostnetwork`.
