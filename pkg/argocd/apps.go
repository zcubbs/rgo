package argocd

import (
	"time"

	"github.com/zcubbs/rgo/pkg/config"
	"github.com/zcubbs/rgo/pkg/k8s"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var gvrApplication = schema.GroupVersionResource{Group: "argoproj.io", Version: "v1alpha1", Resource: "applications"}

func getTimestamp() string {
	// Format timestamp in a way that is compatible with Kubernetes label requirements
	// Replace colons and plus signs with dashes, remove any other invalid characters
	timestamp := time.Now().Format("2006-01-02T15-04-05")
	return timestamp
}

func BuildApplications(apps []config.Application, ns string) []k8s.Object {
	out := make([]k8s.Object, 0, len(apps))
	for _, a := range apps {
		// Set default target revision if not specified
		targetRevision := a.TargetRevision
		if targetRevision == "" {
			targetRevision = "HEAD"
		}

		// Handle OCI Helm charts differently
		if a.IsOCI && a.IsHelm {
			obj := &unstructured.Unstructured{Object: map[string]interface{}{
				"apiVersion": "argoproj.io/v1alpha1",
				"kind":       "Application",
				"metadata": map[string]interface{}{
					"name":      a.Name,
					"namespace": ns,
					"labels": map[string]interface{}{
						"managed-by": "rgo",
						"created-at": getTimestamp(),
					},
				},
				"spec": map[string]interface{}{
					"project": a.Project,
					"sources": []map[string]interface{}{
						{
							"repoURL":        a.OCIRepoURL,
							"targetRevision": a.OCIChartVersion,
							"chart":          a.OCIChartName,
							"helm": map[string]interface{}{
								"passCredentials": true,
								"valueFiles":      a.HelmValueFiles,
								"enableOCI":       true,
							},
						},
						{
							"repoURL":        a.SourceRepoURL,
							"targetRevision": targetRevision,
							"path":           a.SourcePath,
							"ref":            "values",
						},
					},
					"destination": map[string]interface{}{
						"server":    a.DestinationServer,
						"namespace": a.DestinationNamespace,
					},
					"syncPolicy": map[string]interface{}{
						"automated": map[string]interface{}{
							"prune":      true,
							"selfHeal":   true,
							"allowEmpty": false,
						},
					},
				},
			}}
			out = append(out, k8s.Object{Obj: obj, GVR: gvrApplication, NS: ns})
		} else if a.IsHelm {
			// Regular Helm chart
			obj := &unstructured.Unstructured{Object: map[string]interface{}{
				"apiVersion": "argoproj.io/v1alpha1",
				"kind":       "Application",
				"metadata": map[string]interface{}{
					"name":      a.Name,
					"namespace": ns,
					"labels": map[string]interface{}{
						"managed-by": "rgo",
						"created-at": getTimestamp(),
					},
				},
				"spec": map[string]interface{}{
					"project": a.Project,
					"source": map[string]interface{}{
						"repoURL":        a.SourceRepoURL,
						"targetRevision": targetRevision,
						"path":           a.SourcePath,
						"helm": map[string]interface{}{
							"passCredentials": true,
							"valueFiles":      a.HelmValueFiles,
						},
					},
					"destination": map[string]interface{}{
						"server":    a.DestinationServer,
						"namespace": a.DestinationNamespace,
					},
					"syncPolicy": map[string]interface{}{
						"automated": map[string]interface{}{
							"prune":      true,
							"selfHeal":   true,
							"allowEmpty": false,
						},
					},
				},
			}}
			out = append(out, k8s.Object{Obj: obj, GVR: gvrApplication, NS: ns})
		} else {
			// Regular Git application
			obj := &unstructured.Unstructured{Object: map[string]interface{}{
				"apiVersion": "argoproj.io/v1alpha1",
				"kind":       "Application",
				"metadata": map[string]interface{}{
					"name":      a.Name,
					"namespace": ns,
					"labels": map[string]interface{}{
						"managed-by": "rgo",
						"created-at": getTimestamp(),
					},
				},
				"spec": map[string]interface{}{
					"project": a.Project,
					"source": map[string]interface{}{
						"repoURL":        a.SourceRepoURL,
						"targetRevision": targetRevision,
						"path":           a.SourcePath,
					},
					"destination": map[string]interface{}{
						"server":    a.DestinationServer,
						"namespace": a.DestinationNamespace,
					},
					"syncPolicy": map[string]interface{}{
						"automated": map[string]interface{}{
							"prune":      true,
							"selfHeal":   true,
							"allowEmpty": false,
						},
					},
				},
			}}
			out = append(out, k8s.Object{Obj: obj, GVR: gvrApplication, NS: ns})
		}
	}
	return out
}
