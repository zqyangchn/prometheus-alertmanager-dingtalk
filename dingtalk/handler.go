package dingtalk

import (
	"net/http"
	"text/template"

	"go.uber.org/zap"

	"prometheus-alertmanager-dingtalk/config"
	"prometheus-alertmanager-dingtalk/zaplog"
)

var DefaultDingTalk *DingTalk

func SetupInit() {
	tmpl = template.Must(
		template.New("DingTalk").Funcs(funcMap).ParseFiles(config.GetTemplatePath()))

	DefaultDingTalk = NewDingTalk()
}

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