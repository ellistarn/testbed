package karpenter

import (
	"github.com/aws/aws-cdk-go/awscdk/awseks"
	"github.com/aws/aws-cdk-go/awscdk/cloudformationinclude"
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
	"github.com/ellistarn/testbed/pkg/utils/path"
)

type ControllerOptions struct {
	Cluster awseks.Cluster
}

func NewController(scope constructs.Construct, id string, options ControllerOptions) {
	// TODO wait for v0.2.1 to skip prometheus dependency
	// options.Cluster.AddHelmChart(jsii.String(id), &awseks.HelmChartOptions{
	// 	Chart:           jsii.String("karpenter"),
	// 	Namespace:       jsii.String("karpenter"),
	// 	Release:         jsii.String("karpenter"),
	// 	Repository:      jsii.String("https://awslabs.github.io/karpenter/charts/"),
	// 	Version:         jsii.String("0.2.0"),
	// 	CreateNamespace: jsii.Bool(true),
	// })

	cloudformationinclude.NewCfnInclude(scope, jsii.String(id), &cloudformationinclude.CfnIncludeProps{
		TemplateFile: jsii.String(path.RelativeTo("./iam.cfn.yaml")),
		Parameters: &map[string]interface{}{
			"ClusterName":                   *options.Cluster.ClusterName(),
			"OpenIDConnectIdentityProvider": *options.Cluster.OpenIdConnectProvider().OpenIdConnectProviderIssuer(),
		},
	})
}
