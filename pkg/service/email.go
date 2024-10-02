package service

import (
	"log"

	"gopkg.in/gomail.v2"
)

var (
From_mail     = "smoggy46@rustyload.com"
Mail_password = "/N{.j;;BT8"
SMTP_host = "smtp.mail.tm"
)


func SendWarningEmail(email string, message string){
	mess:=gomail.NewMessage()

	mess.SetHeader("From",From_mail)
	mess.SetHeader("To", email)
	mess.SetHeader("Subject", "Warning")
	mess.SetBody("text/plain",message)

	sender:=gomail.NewDialer(SMTP_host,587,From_mail,Mail_password)

	if err:=sender.DialAndSend(mess);err!=nil{
		log.Print(err)
		return
	}
}