package certmanager

import (
	"github.com/aws/aws-cdk-go/awscdk/awseks"
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
)

type ControllerOptions struct {
	Cluster awseks.Cluster
}

func NewController(scope constructs.Construct, id string, options ControllerOptions) {
	options.Cluster.AddHelmChart(jsii.String("CertManager"), &awseks.HelmChartOptions{
		Release:         jsii.String("cert-manager"),
		Chart:           jsii.String("cert-manager"),
		Repository:      jsii.String("https://charts.jetstack.io/"),
		Version:         jsii.String("v1.2.0"),
		CreateNamespace: jsii.Bool(true),
	})
}
