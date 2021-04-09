package clusterapi

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
	cloudformationinclude.NewCfnInclude(scope, jsii.String(id), &cloudformationinclude.CfnIncludeProps{
		TemplateFile: jsii.String(path.RelativeTo("./iam.cfn.yaml")),
	})
}
