# Testbed
Testbed for testing Kubernetes on EC2.

## Prerequisites
```bash
# Add to your bashrc
export CDK_DEFAULT_ACCOUNT=$AWS_ACCOUNT_ID
export CDK_DEFAULT_REGION=$AWS_DEFAULT_REGION
```

## Deploy a testbed Cluster
```bash
make apply
```
## Connect to the testbed Cluster
```bash
$(aws cloudformation describe-stacks --stack-name Testbed | jq -r '.Stacks[].Outputs[0].OutputValue')
```
