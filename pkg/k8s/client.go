package k8s

import (
	"context"
	"encoding/json"
	"fmt"

	"sigs.k8s.io/yaml"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Object struct {
	Obj *unstructured.Unstructured
	GVR schema.GroupVersionResource
	NS  string
}

// New returns a dynamic client using in-cluster config or local kubeconfig fallback
func New() (*Client, error) {
	cfg, err := rest.InClusterConfig()
	if err != nil {
		loading := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			clientcmd.NewDefaultClientConfigLoadingRules(), &clientcmd.ConfigOverrides{},
		)
		cfg, err = loading.ClientConfig()
		if err != nil {
			return nil, err
		}
	}
	dc, err := dynamic.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	return &Client{dc: dc}, nil
}

type Client struct{ dc dynamic.Interface }

// Apply creates or updates an object
func (c *Client) Apply(ctx context.Context, o Object) error {
	res := c.resource(o)
	name := o.Obj.GetName()
	// try get
	existing, err := res.Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			_, err = res.Create(ctx, o.Obj, metav1.CreateOptions{})
			return err
		}
		return err
	}
	// update with resourceVersion
	o.Obj.SetResourceVersion(existing.GetResourceVersion())
	_, err = res.Update(ctx, o.Obj, metav1.UpdateOptions{})
	return err
}

// Delete removes object by name
func (c *Client) Delete(ctx context.Context, o Object) error {
	res := c.resource(o)
	return res.Delete(ctx, o.Obj.GetName(), metav1.DeleteOptions{})
}

func (c *Client) resource(o Object) dynamic.ResourceInterface {
	if o.NS == "" {
		return c.dc.Resource(o.GVR)
	}
	return c.dc.Resource(o.GVR).Namespace(o.NS)
}

// PrintObjects prints objects as yaml or json
func PrintObjects(list []Object, format string) error {
	for _, o := range list {
		if err := PrintObject(o.Obj.Object, format); err != nil {
			return err
		}
	}
	return nil
}

func PrintObject(obj interface{}, format string) error {
	switch format {
	case "json":
		b, err := json.MarshalIndent(obj, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(b))
	default:
		b, err := yaml.Marshal(obj)
		if err != nil {
			return err
		}
		fmt.Println(string(b))
	}
	return nil
}

// ObjectForDelete builds a minimal object to locate & delete
func ObjectForDelete(kind, name, ns string) (Object, error) {
	switch kind {
	case "app", "application":
		return Object{Obj: &unstructured.Unstructured{Object: map[string]interface{}{
			"apiVersion": "argoproj.io/v1alpha1", "kind": "Application", "metadata": map[string]interface{}{"name": name, "namespace": ns},
		}}, GVR: schema.GroupVersionResource{Group: "argoproj.io", Version: "v1alpha1", Resource: "applications"}, NS: ns}, nil
	case "project", "appproject":
		return Object{Obj: &unstructured.Unstructured{Object: map[string]interface{}{
			"apiVersion": "argoproj.io/v1alpha1", "kind": "AppProject", "metadata": map[string]interface{}{"name": name, "namespace": ns},
		}}, GVR: schema.GroupVersionResource{Group: "argoproj.io", Version: "v1alpha1", Resource: "appprojects"}, NS: ns}, nil
	case "secret":
		return Object{Obj: &unstructured.Unstructured{Object: map[string]interface{}{
			"apiVersion": "v1", "kind": "Secret", "metadata": map[string]interface{}{"name": name, "namespace": ns},
		}}, GVR: schema.GroupVersionResource{Group: "", Version: "v1", Resource: "secrets"}, NS: ns}, nil
	default:
		return Object{}, fmt.Errorf("unsupported kind: %s", kind)
	}
}
