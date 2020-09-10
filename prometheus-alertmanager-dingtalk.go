package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
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
	srv := &http.Server{
		Addr:              config.GetListenUri(),
		ReadTimeout:       5 * time.Minute,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Minute,
	}

	zaplog.Logger.Debug("ConfigSetting",
		zap.String("uri", config.GetDingTalkUri()),
		zap.String("securitySettingsType", config.GetSecuritySettingsType()),
		zap.String("secretKey", config.GetSecretKey()),
		zap.String("templatePath", config.GetTemplatePath()),
		zap.Strings("allowLabels", config.GetAllowLables()),
	)

	go func() {
		zaplog.Logger.Info("Web Starting Completed !", zap.String("ListenUri", config.GetListenUri()))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zaplog.Logger.Fatal("web Server start Failed", zap.Error(err))
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 15 秒的超时时间）
	osSignal := make(chan os.Signal)
	signal.Notify(osSignal, os.Interrupt)
	<-osSignal

	// 启动服务器关闭流程
	zaplog.Logger.Info("shutdown server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zaplog.Logger.Fatal("Server Shutdown:", zap.Error(err))
	}
	zaplog.Logger.Info("server shutdown completed !")
}
