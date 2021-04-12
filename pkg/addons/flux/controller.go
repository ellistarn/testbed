package flux

import (
	"strings"

	"github.com/aws/aws-cdk-go/awscdk/awseks"
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
	"github.com/ellistarn/testbed/pkg/utils/file"
	"github.com/ellistarn/testbed/pkg/utils/kubectl"
)

type ControllerOptions struct {
	Cluster awseks.Cluster
	// SyncRepositories defines a set of github uris where YAML at /config is synced to the cluster
	SyncRepositories []string
}

func NewController(scope constructs.Construct, id string, options ControllerOptions) {
	// Generated with `flux install --export > pkg/addons/flux/controller.yaml`
	controller := kubectl.ApplyFile(scope, id, kubectl.ApplyOptions{
		Cluster: options.Cluster,
		FilePath:    file.RelativeTo("./controller.yaml"),
	})

	for _, repository := range options.SyncRepositories {
		name := sanitize(repository)
		awseks.NewKubernetesManifest(scope, jsii.String(name), &awseks.KubernetesManifestProps{
			Cluster:   options.Cluster,
			Overwrite: jsii.Bool(true),
			Manifest: &[]*map[string]interface{}{
				{
					"apiVersion": "source.toolkit.fluxcd.io/v1beta1",
					"kind":       "GitRepository",
					"metadata": map[string]interface{}{
						"name": name,
					},
					"spec": map[string]interface{}{
						"interval": "30s",
						"url":      repository,
						"ref": map[string]interface{}{
							"branch": "main",
						},
					},
				},
				{
					"apiVersion": "kustomize.toolkit.fluxcd.io/v1beta1",
					"kind":       "Kustomization",
					"metadata": map[string]interface{}{
						"name": name,
					},
					"spec": map[string]interface{}{
						"interval": "30s",
						"path":     "/config",
						"prune":    true,
						"sourceRef": map[string]interface{}{
							"kind": "GitRepository",
							"name": name,
						},
					},
				},
			},
		}).Node().AddDependency(controller)
	}
}

// sanitize removes characters that do not conform to kubernetes naming requirements
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
