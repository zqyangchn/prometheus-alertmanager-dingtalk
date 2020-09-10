package main

import (
	"net/http"
	"time"

	"go.uber.org/zap"

	"prometheus-alertmanager-dingtalk/config"
	"prometheus-alertmanager-dingtalk/dingtalk"
	"prometheus-alertmanager-dingtalk/zaplog"
)


func init() {
	config.SetupInit()
	zaplog.SetupInit()
	dingtalk.SetupInit()

	http.HandleFunc("/ready", dingtalk.HandlerReady)
	http.HandleFunc("/healthy", dingtalk.HandlerHealthy)
	http.HandleFunc("/dingtalk/alertmanager", dingtalk.HandlerAlertManager)
}

func main() {
	server := &http.Server{
		Addr:              config.GetListenUri(),
		ReadTimeout:       1 * time.Minute,
		ReadHeaderTimeout: 20 * time.Second,
		WriteTimeout:      2 * time.Minute,
	}

	zaplog.Logger.Debug("ConfigSetting",
		zap.String("uri", config.GetDingTalkUri()),
		zap.String("securitySettingsType", config.GetSecuritySettingsType()),
		zap.String("secretKey", config.GetSecretKey()),
		zap.String("templatePath", config.GetTemplatePath()),
		zap.Strings("allowLabels", config.GetAllowLables()),
	)

	zaplog.Logger.Info("Web Starting Completed !", zap.String("ListenUri", config.GetListenUri()))
	if err := server.ListenAndServe(); err != nil {
		panic(err.Error())
	}
}
