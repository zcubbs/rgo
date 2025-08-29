package argocd

import (
	"strings"

	"github.com/zcubbs/rgo/pkg/config"
	"github.com/zcubbs/rgo/pkg/k8s"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var gvrApplication = schema.GroupVersionResource{Group: "argoproj.io", Version: "v1alpha1", Resource: "applications"}

func BuildApplications(apps []config.Application, ns string) []k8s.Object {
	out := make([]k8s.Object, 0, len(apps))
	for _, a := range apps {
		syncPolicy := map[string]interface{}{}
		if strings.EqualFold(a.SyncPolicy, "automated") {
			syncPolicy["automated"] = map[string]interface{}{"prune": true, "selfHeal": true}
		}
		obj := &unstructured.Unstructured{Object: map[string]interface{}{
			"apiVersion": "argoproj.io/v1alpha1",
			"kind":       "Application",
			"metadata": map[string]interface{}{
				"name":      a.Name,
				"namespace": ns,
			},
			"spec": map[string]interface{}{
				"project": a.Project,
				"destination": map[string]interface{}{
					"namespace": a.DestinationNamespace,
					"server":    a.DestinationServer,
				},
				"source": map[string]interface{}{
					"repoURL": a.SourceRepoURL,
					"path":    a.SourcePath,
				},
				"syncPolicy": syncPolicy,
			},
		}}
		out = append(out, k8s.Object{Obj: obj, GVR: gvrApplication, NS: ns})
	}
	return out
}
