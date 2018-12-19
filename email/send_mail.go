package main

import (
	"fmt"
	"net/smtp"
	"strings"
)

func main() {
	auth := smtp.PlainAuth("", "123@qq.com", "123", "smtp.qq.com")
	to := []string{"zysdio77@sina.com","zysdio77@163.com"}
	nickname := "昵称"
	user := "123@qq.com"
	subject := "go测试发送邮件"
	content_type := "Content-Type: text/plain; charset=UTF-8"
	body := "This is the email body."
	msg := []byte("To: " + strings.Join(to, ",") + "\r\nFrom: " + nickname +
		"<" + user + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	err := smtp.SendMail("smtp.qq.com:25", auth, user, to, msg)
	if err != nil {
		fmt.Printf("send mail error: %v", err)
	} else {
		fmt.Println("send mail successful")
	}
}
