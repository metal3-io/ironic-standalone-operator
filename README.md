# Metal3 Ironic Operator

Operator to maintain an Ironic deployment for Metal3.

## Usage

Let's assume we're creating an Ironic deployment called `ironic` in the
namespace `test`, and that [BMO][bmo] will access it via
`https://ironic.test.svc`.

Start with creating a TLS certificate and a secret for it. In this example,
I'll use a self-signed certificate, you may want to sign it with some
authority, e.g. using built-in Kubernetes facilities.

```bash
openssl req -x509 -new -subj "/CN=ironic.test.svc" -addext "subjectAltName = DNS:ironic.test.svc" \
    -newkey ec -pkeyopt ec_paramgen_curve:secp384r1 -nodes -keyout tls.key -out tls.crt
kubectl create secret tls ironic-tls -n test --key="tls.key" --cert="tls.crt"
```

Now create an API credentials secret and the Ironic configuration:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: ironic-credentials
  namespace: test
type: Opaque
data:
  username: <YOUR USERNAME HERE, BASE64>
  password: <YOUR PASSWORD HERE, BASE64>
---
apiVersion: metal3.io/v1alpha1
kind: Ironic
metadata:
  name: ironic
  namespace: test
spec:
  apiSecretName: ironic-credentials
  tlsSecretName: ironic-tls
  ramdiskSSHKey: "<YOUR SSH PUBLIC KEY HERE>"
```

```bash
kubectl create -f ironic.yaml
```

After some time, you can check the outcome:

```
$ kubectl describe ironic -n test ironic
...
Status:
  Conditions:
    Last Transition Time:  2023-08-25T13:05:35Z
    Message:               ironic is available
    Observed Generation:   2
    Reason:                DeploymentAvailable
    Status:                True
    Type:                  Available
    Last Transition Time:  2023-08-25T13:05:35Z
    Message:               ironic is available
    Observed Generation:   2
    Reason:                DeploymentAvailable
    Status:                False
    Type:                  Progressing
  Ironic Endpoint:
    https://10.96.197.213
```

The endpoint is the service IP corresponding to `ironic.test.svc`.

[bmo]: https://github.com/metal3-io/baremetal-operator

## OpenShift notes

- Running a database requires the user to have `nonroot` or similar SCC.
- Running Ironic requires `nonroot` and `hostnetwork`.
