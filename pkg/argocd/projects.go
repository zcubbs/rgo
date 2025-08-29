package argocd

import (
	"your/module/rgo/pkg/config"
	"your/module/rgo/pkg/k8s"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var gvrAppProject = schema.GroupVersionResource{Group: "argoproj.io", Version: "v1alpha1", Resource: "appprojects"}

func BuildProjects(ps []config.Project, ns string) []k8s.Object {
	out := make([]k8s.Object, 0, len(ps))
	for _, p := range ps {
		dests := make([]map[string]interface{}, 0, len(p.Destinations))
		for _, d := range p.Destinations {
			dests = append(dests, map[string]interface{}{"namespace": d.Namespace, "server": d.Server})
		}
		obj := &unstructured.Unstructured{Object: map[string]interface{}{
			"apiVersion": "argoproj.io/v1alpha1",
			"kind":       "AppProject",
			"metadata": map[string]interface{}{
				"name":      p.Name,
				"namespace": ns,
			},
			"spec": map[string]interface{}{
				"description":  p.Description,
				"destinations": dests,
				"sourceRepos":  anySlice(p.SourceRepos),
			},
		}}
		out = append(out, k8s.Object{Obj: obj, GVR: gvrAppProject, NS: ns})
	}
	return out
}

func anySlice(in []string) []interface{} {
	out := make([]interface{}, 0, len(in))
	for _, s := range in {
		out = append(out, s)
	}
	return out
}
