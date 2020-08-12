package dingtalk

import (
	"net/http"

	"go.uber.org/zap"
)

var DefaultDingTalk *DingTalk

func init() {
	DefaultDingTalk = New()
}

// Handler AlertManager WebHook Request, Send Message To DingTalk
func HandlerAlertManager(w http.ResponseWriter, r *http.Request) {
	if err := DefaultDingTalk.SendAlertManagerMessage(r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Warn("SendAlertManagerMessage Error", zap.Error(err))
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
