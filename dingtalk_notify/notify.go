package dingtalk_notify

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

//
// NewRobot
// @Description: 生成一个新的机器人
// @param token		token用于校验
// @param secret	密钥
// @return *Robot
//
func NewRobot(token, secret string) *Robot {
	return &Robot{
		token:  token,
		secret: secret,
	}
}

//
// sign
// @Description: 验证钉钉服务器的sign
// @param t
// @param secret
// @return string
//
func sign(t int64, secret string) string {
	strToHash := fmt.Sprintf("%d\n%s", t, secret)
	hmac256 := hmac.New(sha256.New, []byte(secret))
	hmac256.Write([]byte(strToHash))
	data := hmac256.Sum(nil)
	return base64.StdEncoding.EncodeToString(data)
}

type Robot struct {
	token, secret string
}

//
// SendMessage
// @Description: 发送信息
// @receiver robot
// @param msg		返回消息
// @return error
//
func (robot *Robot) SendMessage(msg interface{}) error {
	body := bytes.NewBuffer(nil)
	err := json.NewEncoder(body).Encode(msg)
	if err != nil {
		return fmt.Errorf("msg json failed, msg: %v, err: %v", msg, err.Error())
	}
	value := url.Values{}
	value.Set("access_token", robot.token)
	if robot.secret != "" {
		t := time.Now().UnixNano() / 1e6
		value.Set("timestamp", fmt.Sprintf("%d", t))
		value.Set("sign", sign(t, robot.secret))
	}

	// 创建新请求
	request, err := http.NewRequest(http.MethodPost, "https://oapi.dingtalk.com/robot/send", body)
	if err != nil {
		return fmt.Errorf("error request: %v", err.Error())
	}
	request.URL.RawQuery = value.Encode()
	request.Header.Add("Content-Type", "application/json;charset=utf-8")
	res, err := (&http.Client{}).Do(request)
	if err != nil {
		return fmt.Errorf("send dingTalk message failed, error: %v", err.Error())
	}
	defer func() { _ = res.Body.Close() }()
	result, err := ioutil.ReadAll(res.Body)

	if res.StatusCode != 200 {
		return fmt.Errorf("send dingTalk message failed, %s", httpError(request, res, result, "http code is not 200"))
	}
	if err != nil {
		return fmt.Errorf("send dingTalk message failed, %s", httpError(request, res, result, err.Error()))
	}

	type response struct {
		ErrCode int `json:"errcode"`
	}
	var ret response

	if err := json.Unmarshal(result, &ret); err != nil {
		return fmt.Errorf("send dingTalk message failed, %s", httpError(request, res, result, err.Error()))
	}

	if ret.ErrCode != 0 {
		return fmt.Errorf("send dingTalk message failed, %s", httpError(request, res, result, "errcode is not 0"))
	}

	return nil
}
//
// httpError
// @Description: http异常
// @param request
// @param response
// @param body
// @param error
// @return string
//
func httpError(request *http.Request, response *http.Response, body []byte, error string) string {
	return fmt.Sprintf(
		"http request failure, error: %s, status code: %d, %s %s, body:\n%s",
		error,
		response.StatusCode,
		request.Method,
		request.URL.String(),
		string(body),
	)
}
//
// SendTextMessage
// @Description: 发送纯文本信息
// @receiver robot
// @param content	消息内容
// @param atUserIds		返回信息时艾特发送信息者
// @param isAtAll		是否艾特全体成员
// @return error
//
func (robot *Robot) SendTextMessage(content string, atUserIds []string, isAtAll bool) error {
	msg := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": content,
		},
		"at": map[string]interface{}{
			"atUserIds": atUserIds,
			"isAtAll":   isAtAll,
		},
	}

	return robot.SendMessage(msg)
}

//
// SendMarkdownMessage
// @Description: 发送Markdown类型信息
// @receiver robot
// @param title
// @param text
// @param atUserIds
// @param isAtAll
// @return error
//
func (robot *Robot) SendMarkdownMessage(title string, text string, atUserIds []string, isAtAll bool) error {
	msg := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"title": title,
			"text":  text,
		},
		"at": map[string]interface{}{
			"atUserIds": atUserIds,
			"isAtAll":   isAtAll,
		},
	}

	return robot.SendMessage(msg)
}

//
// SendActionCardMessage
// @Description: 发送actioncard类型信息
// @receiver robot
// @param title
// @param text
// @param singleTitle
// @param singleURL
// @return error
//
func (robot *Robot) SendActionCardMessage(title string, text string, singleTitle string, singleURL string) error {
	msg := map[string]interface{}{
		"msgtype": "actionCard",
		"actionCard": map[string]string{
			"title":       title,
			"text":        text,
			"singleTitle": singleTitle,
			"singleURL":   singleURL,
		},
	}

	return robot.SendMessage(msg)
}
