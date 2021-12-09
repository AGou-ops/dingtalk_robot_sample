package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// 生成sign
func hmacSha256(stringToSign string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// 发送纯text
func sendPureText(s, url string) bool {
	content, data := make(map[string]string), make(map[string]interface{})
	content["content"] = s
	data["msgtype"] = "text"
	data["text"] = content
	b, _ := json.Marshal(data)

	resp, err := http.Post(url,
		"application/json",
		bytes.NewBuffer(b))
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	return true
}

// func sendMarkdown(s, url string) bool {

// }

//,func post<F8><F8><F8><F8>

// Getsign 获取sign
func Getsign(secret, webhook string) string {
	timestamp := time.Now().UnixNano() / 1e6
	stringToSign := fmt.Sprintf("%d\n%s", timestamp, secret)
	sign := hmacSha256(stringToSign, secret)
	url := fmt.Sprintf("%s&timestamp=%d&sign=%s", webhook, timestamp, sign)
	return url
}

func tmp() {
	// header中的timestamp + "\n" + 机器人的appSecret当做签名字符串，
	// 使用HmacSHA256算法计算签名，然后进行Base64 encode，得到最终的签名值。
	secret := "849051cc270053f33a0d683bff85d58712e38790dfe1ddfcf86a17ad9df895"
	webhook := "https://oapi.dingtalk.com/robot/send?access_token=37849051cc270053f33a0d683bff85d58712e38790dfe1ddfcf86a17ad9df895"
	url := Getsign(secret, webhook)
	// url = "https://oapi.dingtalk.com/robot/send?access_token=37849051cc270053f33a0d683bff85d58712e38790dfe1ddfcf86a17ad9df895"
	// text := os.Args[1]
	text := "hello_robot_dingtalk"
	sendPureText(text, url)

}
