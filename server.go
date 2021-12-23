package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/AGou-ops/dingtalk/dingtalk_notify"
)

const (
	dateFormat  = "2006-01-02 15:04:05"
	getIPAPI    = "https://ipv4.ipw.cn/api/ip/myip"
	accessToken = "456e2108405af31f4154d884c978641883f50fb5fe295c285e091c3f95ac1bae"
	version     = "0.12"
)

// Sender	  发送消息者
// Content	  消息内容
var (
	Sender  string
	Content string
)

// ReqBody	  请求消息体
// @Description: 请求body结构体
//
type ReqBody struct {
	SenderStaffID string
	Text          struct {
		Content string
	}
}

//
// HTTPServer
// @Description: 创建http服务
//
func HTTPServer() {
	http.HandleFunc("/", handleGetPost)

	log.Printf("Starting server for DingTalk robot HTTP POST...\n")
	log.Println("Server Host on: [::]:5432")
	if err := http.ListenAndServe(":5432", nil); err != nil {
		log.Fatal(err)
	}
}

//
// PrettyString
// @Description: 美化json,仅做日志收集用
// @param str	  未美化之前的json
// @return string	  美化之后的json
// @return error
//
func PrettyString(str string) (string, error) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(str), "", "    "); err != nil {
		return "", err
	}
	return prettyJSON.String(), nil
}

//
// handleGetPost
// @Description: 处理http请求,拒绝GET请求
// @param w
// @param r
//
func handleGetPost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		http.Error(w, "Method GET Denied.", http.StatusForbidden)
		return
	case "POST":
		// log.Fprintf(w,r.Header.Get("timestamp"))
		file, err := os.OpenFile("/var/log/dingtalk_robot.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		// file, err := os.OpenFile("/home/dmy/dingtalk_robot.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)
			os.Exit(-1)
		}

		log.SetOutput(file)
		log.SetPrefix("[info] ")

		log.Println(strings.Repeat("-", 100))
		log.Printf("client[%s] is connected...\n", r.RemoteAddr)
		log.Printf("remote addr is: %s \n", r.RequestURI)
		log.Println("\ntimestamp：", r.Header.Get("timestamp"))
		log.Println("sign：", r.Header.Get("sign"))
		log.Println(strings.Repeat("-", 100))
		// response body
		s, _ := ioutil.ReadAll(r.Body)
		reqBody, _ := PrettyString(string(s))
		log.Println("body: ", reqBody)
		log.Println(strings.Repeat("*", 100))
		resolveResp(string(s))
		log.Printf("Sender is: %s\nContent is: %s\n", Sender, Content)

		handlePostMsg(strings.ToLower(Content), Sender)

		// SendMessage(Sender, GetCurrentIPv4())
		// SendMarkdownMesg(Sender, GetCurrentIPv4())

		log.Println(strings.Repeat("*", 100))
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func handlePostMsg(content, sender string) {

	if strings.Contains(content, "ip") {
		// --- 此处使用markdown格式传递
		getIPMarkdown := " \n #### 🌏您好,今日公网IPv4地址是: \n\n ### " + strings.Repeat("&nbsp;", 15) + GetCurrentIPv4()
		SendMarkdownMesg(Sender, getIPMarkdown)
	} else if strings.Contains(content, "help") {
		getHelpMarkdown := " \n #### **ℹ️帮助信息(请包含以下关键字):** \n\n - **ip: 获取公司的IPv4公网地址;** \n - **测试: 获取当前测试环境占用情况链接;** \n - **help: 获取帮助信息;** \n - **about: 关于该机器人;** \n\n " + GetCurrentIPv4()
		SendMarkdownMesg(sender, getHelpMarkdown)
	} else if strings.Contains(content, "about") {
		SendActionCardMesg("### 🤖关于该机器人\n ![logo](https://agou-images.oss-cn-qingdao.aliyuncs.com/others/robot128_128.png) \n\n Robot version:"+version+" \n\n Backend: go1.17.2 darwin/arm64 \n\n  [Source Code](https://github.com/AGou-ops/dingtalk_robot_sample) \n\n  > 更新日志:\n > 1. 初始化,实现基础功能,获取公网IPv4;", "Read More...", "https://alidocs.dingtalk.com/i/team/O5pXB64OMkEoX7Zv/docs/O5pXBZxL6ZEAoX7Z")
	} else if strings.Contains(content, "测试") {
		getTestEnvLink := " \n [🔗测试环境占用情况](http://192.168.10.203:8088/#/environment/environment)"
		SendMarkdownMesg(sender, getTestEnvLink)
	} else {
		noKeyMarkdown := " ⚠️*抱歉,您的指令有误!* \n\n #### **帮助信息(请包含以下关键字):** \n - **ip: 获取公司的IPv4公网地址;** \n - **测试: 获取当前测试环境占用情况链接;** \n - **help: 获取帮助信息;** \n - **about: 关于该机器人;** \n\n" + GetCurrentIPv4()
		SendMarkdownMesg(sender, noKeyMarkdown)
	}

	// 回头想一想用switch
	// switch content {
	// case strings.Contains(content, "ip"):
	// 	SendMarkdownMesg(sender, ipv4addr)
	// }
}

//
// resolveResp
// @Description: 解析请求体,并传入变量
// @param j
// @return sender	  消息发送者
// @return content	  消息内容
//
func resolveResp(j string) (sender, content string) {
	var rb ReqBody
	json.Unmarshal([]byte(j), &rb)
	Sender, Content = rb.SenderStaffID, rb.Text.Content
	return rb.SenderStaffID, rb.Text.Content
}

//
// SendMessage
// @Description: 发送纯文本信息
// @param atUser	艾特回消息发送者
// @param ipv4addr	  从api中获取的ipv4地址
//
func SendMessage(atUser, mesg string) {

	robot := dingtalk_notify.NewRobot(accessToken, "")
	if err := robot.SendTextMessage(mesg, []string{atUser}, false); err != nil {
		log.Fatal(err)
	}
}

//
// SendMarkdownMesg
// @Description: 发送markdown类型的消息
// @param atUser
// @param ipv4addr
//
func SendMarkdownMesg(atUser, mesg string) {
	robot := dingtalk_notify.NewRobot(accessToken, "")
	err := robot.SendMarkdownMessage(
		"title here.",
		// "@"+atUser+"\n *Your Public IPv4 Address is:*\n### **"+ipv4addr+"**",
		"@"+atUser+mesg,
		[]string{
			atUser,
		},
		false,
	)
	if err != nil {
		log.Fatal(err)
	}

}

func SendActionCardMesg(text, singleTitle, singleURL string) {
	robot := dingtalk_notify.NewRobot(accessToken, "")
	err := robot.SendActionCardMessage("title here.", text, singleTitle, singleURL)
	if err != nil {
		log.Fatal(err)
	}
}

//
// GetCurrentIPv4
// @Description: 调用api获取当前IPv4地址
// @return string
//
func GetCurrentIPv4() string {
	// responseClient, errClient := http.Get(getIPAPI)

	// if errClient != nil {
	// 	log.Printf("Failed to get current public IP address!\n")
	// 	panic(errClient)
	// }
	// defer responseClient.Body.Close()

	// // 获取 http response 的 body
	// body, _ := ioutil.ReadAll(responseClient.Body)
	// clientIP := string(body)
	// return clientIP

	getIPFromFile, err := ioutil.ReadFile("./currentIPv4.txt")
	if err != nil {
		log.Fatal(err)
	}
	return string(getIPFromFile)
}
