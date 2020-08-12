package main

import (
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/zqyangchn/webhook-alertmanager-call-dingtalk/config"
	"github.com/zqyangchn/webhook-alertmanager-call-dingtalk/dingtalk"
	"github.com/zqyangchn/webhook-alertmanager-call-dingtalk/zaplog"
)

var logger = zaplog.Get()

func init() {
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

	logger.Debug("ConfigSetting", zap.String("uri", config.GetDingTalkUri()))
	logger.Debug("ConfigSetting", zap.String("securitySettingsType", config.GetSecuritySettingsType()))
	logger.Debug("ConfigSetting", zap.String("secretKey", config.GetSecretKey()))
	for _, label := range config.GetAllowLables() {
		logger.Debug("ConfigSetting", zap.String("allowLabel", label))
	}

	logger.Info("Web Starting Completed !", zap.String("ListenUri", config.GetListenUri()))
	if err := server.ListenAndServe(); err != nil {
		panic(err.Error())
	}
}
