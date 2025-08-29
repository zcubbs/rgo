package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/zcubbs/rgo/pkg/k8s"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [kind] [name]",
	Short: "Delete a single resource by kind and name (kind: app|project|secret)",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		kind := strings.ToLower(args[0])
		name := args[1]

		obj, err := k8s.ObjectForDelete(kind, name, namespace)
		if err != nil {
			return err
		}

		if dryRun {
			fmt.Printf("[dry-run] would delete %s/%s in namespace %s\n", kind, name, obj.NS)
			return nil
		}

		client, err := k8s.New()
		if err != nil {
			return err
		}
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		return client.Delete(ctx, obj)
	},
}
