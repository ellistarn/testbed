package flux

import (
	"fmt"
	"strings"

	"github.com/aws/aws-cdk-go/awscdk/awseks"
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
	"github.com/ellistarn/testbed/pkg/utils/file"
)

type ControllerOptions struct {
	Cluster      awseks.Cluster
	Repositories []string
}

func NewController(scope constructs.Construct, id string, options ControllerOptions) {
	// Generated with `flux install --export > pkg/addons/flux/controller.yaml`
	file.ApplyYAML(options.Cluster, file.RelativeTo("./controller.yaml"))

	for _, repository := range options.Repositories {
		fmt.Print(sanitize(repository) + "\n\n")
		options.Cluster.AddManifest(jsii.String(repository), &map[string]interface{}{
			"apiVersion": "source.toolkit.fluxcd.io/v1beta1",
			"kind":       "GitRepository",
			"metadata": map[string]interface{}{
				"name": sanitize(repository),
			},
			"spec": map[string]interface{}{
				"interval": "30s",
				"url":      repository,
			},
		})
	}
}

func sanitize(s string) string {
	for old, new := range map[string]string{
		"http://":  "",
		"https://": "",
		"/":        "-",
		".":        "-",
	} {
		s = strings.ReplaceAll(s, old, new)
	}
	return s
}
