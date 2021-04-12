package file

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/aws/aws-cdk-go/awscdk/awseks"
	"github.com/aws/jsii-runtime-go"
	"gopkg.in/yaml.v3"
)

// Relative returns an absolute path relative to the provided path
func RelativeTo(path string) string {
	_, file, _, _ := runtime.Caller(1)
	return filepath.Join(filepath.Dir(file), path)
}

// ApplyYAML applies the YAML at the provided path to the specified cluster.
// This function makes @ellistarn very sad.
func ApplyYAML(cluster awseks.Cluster, path string) {
	file, err := os.Open(path)
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
		current := cluster.AddManifest(jsii.String(fmt.Sprint(i)), &resource)
		// Order resource application to ensure dependencies are respected (e.g. namespace creation)
		if last != nil {
			current.Node().AddDependency(last)
		}
		last = current
	}
}
