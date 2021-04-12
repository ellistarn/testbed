package main

import (
	"os"

	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/jsii-runtime-go"
	"github.com/ellistarn/testbed/pkg/testbed"
)


func main() {
	app := awscdk.NewApp(nil)

	testbed.NewStack(app, "Testbed", &testbed.StackOptions{
		StackProps: awscdk.StackProps{Env: &awscdk.Environment{
			Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
			Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
		}},
		SyncRepositories: []string{
			"https://github.com/ellistarn/testbed",
		},
	})

	app.Synth(nil)
}
