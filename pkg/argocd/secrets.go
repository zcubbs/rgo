package argocd

import (
	"fmt"
	"strings"

	"github.com/zcubbs/rgo/pkg/config"
	"github.com/zcubbs/rgo/pkg/k8s"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var gvrSecret = schema.GroupVersionResource{Group: "", Version: "v1", Resource: "secrets"}

func BuildRepoSecrets(repos []config.Repository, ns string) []k8s.Object {
	out := make([]k8s.Object, 0, len(repos))
	for _, r := range repos {
		name := secretNameFromURL(r.URL)
		obj := &unstructured.Unstructured{Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Secret",
			"metadata": map[string]interface{}{
				"name":      name,
				"namespace": ns,
				"labels": map[string]interface{}{
					"argocd.argoproj.io/secret-type": "repository",
				},
			},
			"stringData": map[string]interface{}{
				"url": r.URL,
			},
		}}
		out = append(out, k8s.Object{Obj: obj, GVR: gvrSecret, NS: ns})
	}
	return out
}

func BuildCredentialSecrets(creds []config.Credential, ns string) []k8s.Object {
	out := make([]k8s.Object, 0, len(creds))
	for _, c := range creds {
		name := secretNameFromURL(c.URL)
		stringData := map[string]interface{}{"url": c.URL}
		if c.Username != "" {
			stringData["username"] = c.Username
		}
		if c.Password != "" {
			stringData["password"] = c.Password
		}
		if c.SSHKey != "" {
			stringData["sshPrivateKey"] = c.SSHKey
		}
		obj := &unstructured.Unstructured{Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Secret",
			"metadata": map[string]interface{}{
				"name":      name,
				"namespace": ns,
				"labels": map[string]interface{}{
					"argocd.argoproj.io/secret-type": "repository",
				},
			},
			"stringData": stringData,
		}}
		out = append(out, k8s.Object{Obj: obj, GVR: gvrSecret, NS: ns})
	}
	return out
}

func secretNameFromURL(url string) string {
	name := strings.ToLower(url)
	name = strings.TrimPrefix(name, "https://")
	name = strings.TrimPrefix(name, "http://")
	name = strings.TrimPrefix(name, "ssh://")
	name = strings.ReplaceAll(name, "/", "-")
	name = strings.ReplaceAll(name, ":", "-")
	name = strings.ReplaceAll(name, "@", "-")
	if len(name) > 50 {
		name = name[:50]
	}
	return fmt.Sprintf("repo-%s", name)
}
