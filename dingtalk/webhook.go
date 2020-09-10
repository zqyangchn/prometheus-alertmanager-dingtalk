package dingtalk

import (
	"net/http"

	"go.uber.org/zap"

	"prometheus-alertmanager-dingtalk/zaplog"
)

var DefaultDingTalk *DingTalk

// Handler AlertManager WebHook Request, Send Message To DingTalk
func HandlerAlertManager(w http.ResponseWriter, r *http.Request) {
	if err := DefaultDingTalk.SendAlertManagerMessage(r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		zaplog.Logger.Warn("SendAlertManagerMessage Error", zap.Error(err))
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Handler Ready
func HandlerReady(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// Handler Healthy
func HandlerHealthy(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
