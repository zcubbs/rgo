package cmd

import (
	"context"
	"fmt"
	"time"

	"your/module/rgo/pkg/argocd"
	"your/module/rgo/pkg/config"
	"your/module/rgo/pkg/k8s"

	"github.com/spf13/cobra"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply all resources from config (projects, repos/creds, applications)",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		var objs []k8s.Object
		objs = append(objs, argocd.BuildProjects(cfg.Projects, namespace)...)
		objs = append(objs, argocd.BuildRepoSecrets(cfg.Repositories, namespace)...)
		objs = append(objs, argocd.BuildCredentialSecrets(cfg.Credentials, namespace)...)
		objs = append(objs, argocd.BuildApplications(cfg.Applications, namespace)...)

		if dryRun {
			return k8s.PrintObjects(objs, output)
		}

		client, err := k8s.New()
		if err != nil {
			return err
		}
		ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
		defer cancel()

		for _, obj := range objs {
			if err := client.Apply(ctx, obj); err != nil {
				return err
			}
		}
		fmt.Println("Applied successfully")
		return nil
	},
}
