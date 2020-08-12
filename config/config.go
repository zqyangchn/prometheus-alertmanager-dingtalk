package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

var config Config

type Config struct {
	ListenUri string `yaml:"listenUri"`

	LogLevel      string `yaml:"logLevel"`
	LogOutput     string `yaml:"logOutput"`
	LogOutputFile string `yaml:"logOutputFile"`

	DingTalkUri          string   `yaml:"dingTalkUri"`
	SecuritySettingsType string   `yaml:"securitySettingsType"`
	SecretKey            string   `yaml:"secretKey"`
	TemplatePath         string   `yaml:"templatePath"`
	AllowLables          []string `yaml:"allowLables"`
}

func GetLogLevel() string {
	return config.LogLevel
}
func GetLogOutput() string {
	return config.LogOutput
}
func GetLogOutputFile() string {
	return config.LogOutputFile
}
func init() {
	if err := initConfig(); err != nil {
		panic(err)
	}
}

func initConfig() error {
	configPath := os.Getenv("DINGTALK_CONFIG_FILE")

	if configPath == "" {
		return errors.New("Environment variable DINGTALK_CONFIG_FILE is null")
	}
	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return err
	}

	return err
}

func GetListenUri() string {
	return config.ListenUri
}

func GetDingTalkUri() string {
	return config.DingTalkUri
}

func GetSecuritySettingsType() string {
	return config.SecuritySettingsType
}

func GetSecretKey() string {
	return config.SecretKey
}

func GetTemplatePath() string {
	return config.TemplatePath
}

func GetAllowLables() []string {
	return config.AllowLables
}

func GetTemplateText() string {
	return `
{{ define "__title__" -}}
[{{ .Status | ToUpper }}{{ if eq .Status "firing" }}:{{ .Alerts | len }}{{ end }}] {{ index .GroupLabels.AlertName }}
{{- end }}


{{ define "__alertmanagerURL" -}}
{{ .ExternalURL }}/#/alerts?receiver={{ .Receiver }}
{{- end }}


{{ define "__content__" -}}
#### **\[{{ .Status | ToUpper }}{{ if eq .Status "firing" }}:{{ .Alerts | len }}{{ end }}\]** **[{{ index .GroupLabels.AlertName }}]({{ template "__alertmanagerURL" . }})**
{{- end }}


{{ define "__alerts_common_labels__" }}
{{ range .SortedAllowPairs }}
> **{{ .Name | Title }}**: {{ .Value | html }}
{{ end }}
{{- end }}


{{ define "__alerts_instance_lists__" }}
{{ range $i, $v := . }}
##### **警报: {{ Increase $i }}**
{{ range .Annotations.SortedAllowPairs }}
> ##### **{{ .Name | Title }}**: {{ .Value | html }}
{{ end }}
> ###### Report From **[Prometheus]({{ .GeneratorURL }})** at **{{ FormatTime .StartsAt }}**
{{- end }}
![screenshot](https://prom.i-morefun.net/MaineCoon.jpg)
{{- end }}


{{ define "__text__" -}}
{{ template "__content__" . }}
{{ template "__alerts_common_labels__" .CommonLabels }}
{{ template "__alerts_instance_lists__" .Alerts }}
{{- end }}
`
}
