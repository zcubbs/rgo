package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Example YAML:
// ---
// projects:
//   - name: demo-proj
//     description: Demo project
//     sourceRepos: ["https://github.com/zcubbs/hotpot"]
////     destinations:
//       - namespace: default
//         server: https://kubernetes.default.svc
// applications:
//   - name: demo-app
//     project: demo-proj
//     sourceRepoURL: https://github.com/zcubbs/hotpot
//     sourcePath: manifests/app
//     destinationNamespace: default
//     destinationServer: https://kubernetes.default.svc
//     syncPolicy: automated
// repositories:
//   - url: https://github.com/zcubbs/go-k8s
//     type: git
//     name: demo-repo
// credentials:
//   - url: https://github.com
//     username: ${GIT_USERNAME}
//     password: ${GIT_PASSWORD}
//     name: demo-cred

type Config struct {
	Projects     []Project     `mapstructure:"projects"`
	Applications []Application `mapstructure:"applications"`
	Repositories []Repository  `mapstructure:"repositories"`
	Credentials  []Credential  `mapstructure:"credentials"`
}

type Project struct {
	Name         string        `mapstructure:"name"`
	Description  string        `mapstructure:"description"`
	SourceRepos  []string      `mapstructure:"sourceRepos"`
	Destinations []Destination `mapstructure:"destinations"`
}

type Destination struct {
	Namespace string `mapstructure:"namespace"`
	Server    string `mapstructure:"server"`
}

type Application struct {
	Name                 string `mapstructure:"name"`
	Project              string `mapstructure:"project"`
	DestinationNamespace string `mapstructure:"destinationNamespace"`
	DestinationServer    string `mapstructure:"destinationServer"`
	SourceRepoURL        string `mapstructure:"sourceRepoURL"`
	SourcePath           string `mapstructure:"sourcePath"`
	SyncPolicy           string `mapstructure:"syncPolicy"`
}

type Repository struct {
	URL  string `mapstructure:"url"`
	Type string `mapstructure:"type"`
	Name string `mapstructure:"name"`
}

type Credential struct {
	URL      string `mapstructure:"url"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	SSHKey   string `mapstructure:"sshKey"`
	Name     string `mapstructure:"name"`
}

func Load() (Config, error) {
	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		return c, fmt.Errorf("config unmarshal: %w", err)
	}
	return c, nil
}
