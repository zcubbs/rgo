package cmd

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile   string
	namespace string
	dryRun    bool
	output    string // yaml|json
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, err := fmt.Fprintln(os.Stderr, err)
		if err != nil {
			// panic if we can't write to stderr
			panic(err)
		}
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "rgo",
	Short: "rgo — manage Argo CD apps, projects, repos & creds via Kubernetes CRDs",
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Path to config file (YAML)")
	rootCmd.PersistentFlags().StringVar(&namespace, "namespace", "argo-cd", "Argo CD namespace")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Preview resources instead of applying")
	rootCmd.PersistentFlags().StringVar(&output, "output", "yaml", "Output format for dry-run: yaml|json")

	rootCmd.AddCommand(applyCmd)
	rootCmd.AddCommand(deleteCmd)
}

func initConfig() {
	// 1. Load .env file if present
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			fmt.Println("warning: could not load .env:", err)
		}
	}

	// 2. Config file (YAML)
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
	}

	// 3. Env vars (including those loaded from .env)
	viper.SetEnvPrefix("RGO")
	viper.AutomaticEnv()

	// 4. Merge config file if found
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config:", viper.ConfigFileUsed())
	}
}
