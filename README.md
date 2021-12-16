[![Docker Repository on Quay](https://quay.io/repository/jerry153fish/aws-secrets/status "Docker Repository on Quay")](https://quay.io/repository/jerry153fish/aws-secrets)
![Github Action](https://github.com/jerry153fish/aws-secrets/actions/workflows/test.yaml/badge.svg)

## aws-secrets

A Kubernetes operator to convert cloudformation outputs to k8s secrets.

## How to use

1. apply the manifest file

```
kubeclt apply -f manifest.yaml
```

2. apply CRD file eg:

```
apiVersion: cfn.jerry153fish.com/v1alpha1
kind: Secrets
metadata:
  name: secrets-sample
spec:
  plainCreds:
    - key: test
      value: abcd
    - key: test1
      value: abc12
  cfn:
    - key: test2
      stackName: the-cfn-stack-name
      outputKey: rds-writer-endpoint

```
