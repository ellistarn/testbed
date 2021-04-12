package kubectl

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-cdk-go/awscdk/awseks"
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
	"gopkg.in/yaml.v3"
)


type ApplyOptions struct {
	Cluster awseks.Cluster
	FilePath    string
}

// ApplyFile applies the YAML at the provided path to the specified cluster.
// This function makes @ellistarn very sad.
func ApplyFile(scope constructs.Construct, id string, options ApplyOptions) awseks.KubernetesManifest {
	file, err := os.Open(options.FilePath)
	if err != nil {
		panic(err)
	}
	decoder := yaml.NewDecoder(file)

	// Apply manifests in a loop rather than as a multidoc to avoid CFN size limits
	var last awseks.KubernetesManifest
	for i := 0; true; i++ {
		// Decode resource
		resource := map[string]interface{}{}
		err := decoder.Decode(resource)
		// End of file, exit
		if errors.Is(err, io.EOF) {
			break
		}
		// Failed to parse YAML
		if err != nil {
			panic(err)
		}
		// Skip if resource is empty
		if len(resource) == 0 {
			continue
		}
		// Apply to cluster
		current := awseks.NewKubernetesManifest(scope, jsii.String(fmt.Sprintf("Manifest-%d", i)), &awseks.KubernetesManifestProps{
			Cluster:   options.Cluster,
			Overwrite: jsii.Bool(true),
			Manifest:  &[]*map[string]interface{}{&resource},
		})
		// Order resource application to ensure dependencies are respected (e.g. namespace creation)
		if last != nil {
			current.Node().AddDependency(last)
		}
		last = current
	}
	return last
}
