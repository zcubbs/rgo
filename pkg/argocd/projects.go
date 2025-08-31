package argocd

import (
	"github.com/zcubbs/rgo/pkg/config"
	"github.com/zcubbs/rgo/pkg/k8s"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var gvrAppProject = schema.GroupVersionResource{Group: "argoproj.io", Version: "v1alpha1", Resource: "appprojects"}

func BuildProjects(projects []config.Project, ns string) []k8s.Object {
	out := make([]k8s.Object, 0, len(projects))
	for _, p := range projects {
		obj := &unstructured.Unstructured{Object: map[string]interface{}{
			"apiVersion": "argoproj.io/v1alpha1",
			"kind":       "AppProject",
			"metadata": map[string]interface{}{
				"name":      p.Name,
				"namespace": ns,
				"labels": map[string]interface{}{
					"managed-by": "rgo",
					"created-at": getTimestamp(),
				},
			},
			"spec": map[string]interface{}{
				"description": p.Description,
				"sourceRepos": p.SourceRepos,
				"destinations": func() []interface{} {
					out := make([]interface{}, 0, len(p.Destinations))
					for _, d := range p.Destinations {
						out = append(out, map[string]interface{}{
							"server":    d.Server,
							"namespace": d.Namespace,
						})
					}
					return out
				}(),
			},
		}}
		out = append(out, k8s.Object{Obj: obj, GVR: gvrAppProject, NS: ns})
	}
	return out
}
