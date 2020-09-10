package dingtalk

import (
	"bytes"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/pkg/errors"

	"prometheus-alertmanager-dingtalk/config"
)

var (
	funcMap = template.FuncMap{
		"ToUpper":    strings.ToUpper,
		"ToLower":    strings.ToLower,
		"Title":      strings.Title,
		"Increase":   increase,
		"FormatTime": formatTime,
	}
	dingTalkTemplate *template.Template
)

func SetupInit() {
	dingTalkTemplate = template.Must(template.New("DingTalk").Funcs(funcMap).ParseFiles(config.GetTemplatePath()))

	DefaultDingTalk = New()
}

type Pair struct {
	Name, Value string
}
type Pairs []Pair

type KV map[string]string

// 解析AlertManager json消息的结构体
type AlertManagerMessage struct {
	Receiver string `json:"receiver"`
	Status   string `json:"status"`
	Alerts   Alerts

	GroupLabels struct {
		AlertName string `json:"alertname"`
	} `json:"groupLabels"`

	CommonLabels      KV     `json:"commonLabels"`
	CommonAnnotations KV     `json:"commonAnnotations"`
	ExternalURL       string `json:"externalURL"`

	dingTalkTemplate *template.Template
}

type Alerts []Alert
type Alert struct {
	Status       string    `json:"status"`
	Labels       KV        `json:"labels"`
	Annotations  KV        `json:"annotations"`
	StartsAt     time.Time `json:"startsAt"`
	EndsAt       time.Time `json:"endsAt"`
	GeneratorURL string    `json:"generatorURL"`
}

func NewAlertManagerMessage() *AlertManagerMessage {
	return &AlertManagerMessage{
		dingTalkTemplate: dingTalkTemplate,
	}
}

// 对map类型转化成结构体, 允许特定的kv输出, 并排序
func (kv KV) SortedAllowPairs() Pairs {
	pairs := make([]Pair, 0, len(kv))
	keys := make([]string, 0, len(kv))

	// 仅允许特定key输出
	for k := range kv {
		if !keyExist(k, config.GetAllowLables()) {
			continue
		} else {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	for _, k := range keys {
		pairs = append(pairs, Pair{k, kv[k]})
	}

	return pairs
}

// 判断key是否存在于允许发送消息列表
func keyExist(key string, allow []string) bool {
	for _, k := range allow {
		if strings.ToLower(key) == k {
			return true
		}
	}

	return false
}

func (m *AlertManagerMessage) FilterFiringInformation() {
	firingAlerts := make(Alerts, 0, len(m.Alerts))
	for _, alert := range m.Alerts {
		if alert.Status == "firing" {
			firingAlerts = append(firingAlerts, alert)
		}
	}
	m.Alerts = firingAlerts
}

// 模板函数, 显示实例ID的时候增加1
func increase(i int) int {
	return i + 1
}

// alertManager 发送消息的时间UTC, 对于CST, 增加8小时, 并格式化时间的输出
func formatTime(t time.Time) string {
	return t.Add(8 * time.Hour).Format("2006-01-02 15:04:05")
}

func (m *AlertManagerMessage) ParseDingTalkTemplate() (string, string, error) {
	var (
		titleBuffer bytes.Buffer
		textBuffer  bytes.Buffer
	)

	// title
	if err := m.dingTalkTemplate.ExecuteTemplate(&titleBuffer, "__title__", m); err != nil {
		return "", "", err
	}
	title := titleBuffer.String()

	// text
	switch strings.ToLower(m.Status) {
	case "firing":
		if err := m.dingTalkTemplate.ExecuteTemplate(&textBuffer, "__text__", m); err != nil {
			return "", "", err
		}
	case "resolved":
		if err := m.dingTalkTemplate.ExecuteTemplate(&textBuffer, "__content__", m); err != nil {
			return "", "", err
		}
	default:
		return "", "", errors.Errorf("unknown type: %s", m.Status)
	}
	text := textBuffer.String()

	//logger.Debug("ParseDingTalkTemplate", zap.String("title", title), zap.String("text", text))
	return title, text, nil
}
