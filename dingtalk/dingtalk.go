package dingtalk

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"prometheus-alertmanager-dingtalk/config"
	"prometheus-alertmanager-dingtalk/zaplog"
)

type DingTalk struct {
	HC *http.Client

	Uri string

	// Dingtalk Security Settings Type
	// 1. CustomKeywords		最多可以设置10个关键词，消息中至少包含其中1个关键词才可以发送成功
	// 2. Endorsement			加签, 需要提供 SecretKey
	// 3. IPAddressSegment		只有来自IP地址范围内的请求才会被正常处理
	SecuritySettingsType string
	SecretKey            string
}

// send to dingtalk json body struct
type Notification struct {
	MessageType string                `json:"msgtype"`
	Markdown    *NotificationMarkdown `json:"markdown,omitempty"`
	At          *NotificationAt       `json:"at,omitempty"`
}

type NotificationAt struct {
	AtMobiles []string `json:"atMobiles,omitempty"`
	IsAtAll   bool     `json:"isAtAll,omitempty"`
}

type NotificationMarkdown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

func New() *DingTalk {
	d := &DingTalk{
		Uri: config.GetDingTalkUri(),
		HC: &http.Client{
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   5 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,

				MaxIdleConns:          5,
				IdleConnTimeout:       30 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 5 * time.Second,
				ResponseHeaderTimeout: 5 * time.Second,
			},
			Timeout: 2 * time.Second,
		},
	}

	d.SecuritySettingsType = config.GetSecuritySettingsType()
	d.SecretKey = config.GetSecretKey()

	return d
}

func (d *DingTalk) MakeTimestamp() string {
	return strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
}

func (d *DingTalk) HmacSha256Base64Encode(timestamp string) string {
	stringToSign := strings.Join([]string{timestamp, d.SecretKey}, "\n")

	sig := hmac.New(sha256.New, []byte(d.SecretKey))
	sig.Write([]byte(stringToSign))

	return base64.StdEncoding.EncodeToString(sig.Sum(nil))
}

func (d *DingTalk) BuildAlertManagerMessagePayload(r *http.Request) (*bytes.Reader, error) {
	payload, err := ioutil.ReadAll(r.Body)
	defer func() {
		if err := r.Body.Close(); err != nil {
			return
		}
	}()
	if err != nil {
		return nil, err
	}

	alertManagerMessage := NewAlertManagerMessage()
	if err := json.Unmarshal(payload, &alertManagerMessage); err != nil {
		return nil, err
	}
	alertManagerMessage.FilterFiringInformation()
	title, text, err := alertManagerMessage.ParseDingTalkTemplate()
	if err != nil {
		return nil, err
	}

	notification, notificationAt, notificationMarkdown := new(Notification), new(NotificationAt), new(NotificationMarkdown)

	notificationAt.IsAtAll = true
	notificationMarkdown.Title, notificationMarkdown.Text = title, text

	notification.MessageType = "markdown"
	notification.At = notificationAt
	notification.Markdown = notificationMarkdown

	requestBody, err := json.Marshal(notification)
	if err != nil {
		return nil, err
	}

	zaplog.Logger.Debug("Build AlertManagerMessage Payload From AlertManager Message Completed !",
		zap.String("Status", alertManagerMessage.Status),
		zap.String("AlertName", alertManagerMessage.GroupLabels.AlertName),
	)
	return bytes.NewReader(requestBody), nil
}

func (d *DingTalk) SendAlertManagerMessage(r *http.Request) error {
	body, err := d.BuildAlertManagerMessagePayload(r)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", d.Uri, body)
	if err != nil {
		return err
	}

	if d.SecuritySettingsType == "Endorsement" {
		timestamp := d.MakeTimestamp()
		base64EncodeSign := d.HmacSha256Base64Encode(timestamp)

		q := req.URL.Query()            // Get a copy of the query values.
		q.Add("timestamp", timestamp)   // Add query timestamp
		q.Add("sign", base64EncodeSign) // Add query sign
		req.URL.RawQuery = q.Encode()   // Encode and assign back to the original query.
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := d.HC.Do(req)
	if resp != nil {
		defer func() {
			if err := resp.Body.Close(); err != nil {
				return
			}
		}()
	}
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.Errorf("Send Message Response StatusCode is not 200, StatusCode: %d", resp.StatusCode)
	}

	payload, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var response struct {
		ErrorCode    int    `json:"errcode"`
		ErrorMessage string `json:"errmsg"`
	}
	if err := json.Unmarshal(payload, &response); err != nil {
		return err
	}
	if response.ErrorCode != 0 {
		return errors.Errorf("Send Message Response ErrorCode is not zero, ErrorCode: %d, ErrorMessage: %s",
			response.ErrorCode, response.ErrorMessage)
	}

	zaplog.Logger.Debug("Send AlertManagerMessage Payload To DingTalk Completed !")
	return nil
}