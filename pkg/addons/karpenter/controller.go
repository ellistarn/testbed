package karpenter

import (
	"github.com/aws/aws-cdk-go/awscdk/awseks"
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
)

type ControllerOptions struct {
	Cluster awseks.Cluster
}

func NewController(scope constructs.Construct, id string, options ControllerOptions) {
	options.Cluster.AddHelmChart(jsii.String(id), &awseks.HelmChartOptions{
		Chart:           jsii.String("karpenter"),
		Namespace:       jsii.String("karpenter"),
		Release:         jsii.String("karpenter"),
		Repository:      jsii.String("https://awslabs.github.io/karpenter/charts/"),
		Version:         jsii.String("0.2.1"),
		CreateNamespace: jsii.Bool(true),
	})

	// cloudformationinclude.NewCfnInclude(scope, jsii.String(id), &cloudformationinclude.CfnIncludeProps{
	// 	TemplateFile: jsii.String(path.RelativeTo("./iam.cfn.yaml")),
	// 	Parameters: &map[string]interface{}{
	// 		"ClusterName":                   *options.Cluster.ClusterName(),
	// 		"OpenIDConnectIdentityProvider": *options.Cluster.OpenIdConnectProvider().OpenIdConnectProviderIssuer(),
	// 	},
	// })

	awseks.NewKubernetesManifest(scope, jsii.String("Provisioner"), &awseks.KubernetesManifestProps{
		Cluster:   options.Cluster,
		Overwrite: jsii.Bool(true),
		Manifest: &[]*map[string]interface{}{
			{
				"apiVersion": "provisioning.karpenter.sh/v1alpha1",
				"kind":       "Provisioner",
				"metadata": map[string]interface{}{
					"name": "default",
				},
				"spec": map[string]interface{}{
					"cluster": map[string]interface{}{
						"name":     *options.Cluster.ClusterName(),
						"caBundle": *options.Cluster.ClusterCertificateAuthorityData(),
						"endpoint": *options.Cluster.ClusterEndpoint(),
					},
				},
			},
		},
	})
}
