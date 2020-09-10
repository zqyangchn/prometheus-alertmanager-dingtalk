package config

import (
	"os"

	"github.com/pkg/errors"
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
func SetupInit() {
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
			return
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
	if config.TemplatePath == "" {
		return "/etc/prometheus-alertmanager-dingtalk/dingtalk.tmpl"
	}
	return config.TemplatePath
}

func GetAllowLables() []string {
	return config.AllowLables
}
