package argocd

import (
	"fmt"
	"os"
	"strings"

	"github.com/zcubbs/rgo/pkg/config"
	"github.com/zcubbs/rgo/pkg/k8s"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var gvrSecret = schema.GroupVersionResource{Group: "", Version: "v1", Resource: "secrets"}

// ensureGitSuffix ensures that Git repository URLs end with .git
func ensureGitSuffix(url string) string {
	// Skip if it's not a Git repository or already has .git suffix
	if !strings.HasPrefix(url, "http") || strings.HasSuffix(url, ".git") {
		return url
	}

	// Remove trailing slash if present
	if strings.HasSuffix(url, "/") {
		url = url[:len(url)-1]
	}

	// Add .git suffix
	return url + ".git"
}

func BuildRepoSecrets(repos []config.Repository, ns string) []k8s.Object {
	out := make([]k8s.Object, 0, len(repos))
	for _, r := range repos {
		var name string
		if r.Name != "" {
			name = fmt.Sprintf("repo-%s", r.Name)
		} else {
			name = secretNameFromURL(r.URL)
		}

		// Create stringData map with URL
		stringData := map[string]interface{}{}

		// Handle URL based on repository type
		if r.Type == "helm" || r.Type == "oci" {
			// For Helm/OCI repositories, don't add .git suffix
			if r.Type == "oci" {
				stringData["url"] = r.URL
				stringData["type"] = "helm"
				stringData["enableOCI"] = "true"
				stringData["name"] = name // Add name field for OCI repositories
			} else {
				// Regular Helm repo
				stringData["url"] = r.URL
				stringData["type"] = "helm"
				stringData["name"] = name // Add name field for Helm repositories
			}
		} else {
			// Default to Git repository type
			stringData["url"] = ensureGitSuffix(r.URL)
			stringData["type"] = "git"
			stringData["name"] = name // Add name field for Git repositories
		}

		// Add username if provided
		if r.Username != "" {
			stringData["username"] = resolveEnvVar(r.Username)
		}

		// Add password if provided
		if r.Password != "" {
			stringData["password"] = resolveEnvVar(r.Password)
		}

		// Add SSH key if provided
		if r.SSHKey != "" {
			stringData["sshPrivateKey"] = r.SSHKey
		}

		obj := &unstructured.Unstructured{Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Secret",
			"metadata": map[string]interface{}{
				"name":      name,
				"namespace": ns,
				"labels": map[string]interface{}{
					"argocd.argoproj.io/secret-type": "repository",
					"managed-by":                     "rgo",
					"created-at":                     getTimestamp(),
				},
			},
			"stringData": stringData,
		}}
		out = append(out, k8s.Object{Obj: obj, GVR: gvrSecret, NS: ns})
	}
	return out
}

// resolveEnvVar resolves environment variables in the format ${VAR_NAME}
func resolveEnvVar(value string) string {
	// If the value doesn't contain ${, return it as is
	if !strings.Contains(value, "${") {
		return value
	}

	// Find all environment variables in the format ${VAR_NAME}
	result := value
	for {
		start := strings.Index(result, "${")
		if start == -1 {
			break
		}
		end := strings.Index(result[start:], "}")
		if end == -1 {
			break
		}
		end += start

		// Extract the variable name
		varName := result[start+2 : end]

		// Get the environment variable value
		varValue := os.Getenv(varName)

		// Replace the variable in the string
		result = result[:start] + varValue + result[end+1:]
	}

	return result
}

func BuildCredentialSecrets(creds []config.Credential, ns string) []k8s.Object {
	out := make([]k8s.Object, 0, len(creds))
	for _, c := range creds {
		var name string
		if c.Name != "" {
			name = fmt.Sprintf("repo-%s", c.Name)
		} else {
			name = secretNameFromURL(c.URL)
		}
		stringData := map[string]interface{}{"url": ensureGitSuffix(c.URL)}
		if c.Username != "" {
			// Resolve environment variables in username
			stringData["username"] = resolveEnvVar(c.Username)
		}
		if c.Password != "" {
			// Resolve environment variables in password
			stringData["password"] = resolveEnvVar(c.Password)
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
					"managed-by":                     "rgo",
					"created-at":                     getTimestamp(),
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
