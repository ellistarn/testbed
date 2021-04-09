package main

import (
	"os"

	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/awseks"
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
	"github.com/ellistarn/testbed/pkg/addons/certmanager"
	"github.com/ellistarn/testbed/pkg/addons/clusterapi"
)

type TestbedStackOptions struct {
	awscdk.StackProps
}

func NewTestbedStack(scope constructs.Construct, id string, props *TestbedStackOptions) awscdk.Stack {
	stack := awscdk.NewStack(scope, &id, &props.StackProps)

	vpc := awsec2.NewVpc(stack, jsii.String("VPC"), &awsec2.VpcProps{Cidr: jsii.String("10.0.0.0/16")})

	cluster := awseks.NewCluster(stack, jsii.String("Cluster"), &awseks.ClusterProps{
		Vpc:                     vpc,
		ClusterName:             jsii.String(id),
		Version:                 awseks.KubernetesVersion_V1_18(),
		DefaultCapacityInstance: awsec2.InstanceType_Of(awsec2.InstanceClass_BURSTABLE3, awsec2.InstanceSize_XLARGE2),
	})

	certmanager.NewController(stack, "CertManagerController", certmanager.ControllerOptions{Cluster: cluster})
	// karpenter.NewController(stack, "KarpenterController", karpenter.ControllerOptions{Cluster: cluster})
	clusterapi.NewController(stack, "ClusterapiController", clusterapi.ControllerOptions{Cluster: cluster})
	return stack
}

func main() {
	app := awscdk.NewApp(nil)

	NewTestbedStack(app, "Testbed", &TestbedStackOptions{
		awscdk.StackProps{Env: &awscdk.Environment{
			Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
			Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
		}},
	})

	app.Synth(nil)
}
