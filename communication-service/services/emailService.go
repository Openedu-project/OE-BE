package services

import (
	"fmt"
	"net/smtp"

	"communication-service/configs"
)

func SendWelcomeEmail(to, name string) error {
	from := configs.Env.SMTPUser
	pass := configs.Env.SMTPPass
	host := configs.Env.SMTPHost
	port := configs.Env.SMTPPort

	auth := smtp.PlainAuth("", from, pass, host)
	subject := "Subject: Welcome to OpenEdu!\n"
	body := fmt.Sprintf("Hi %s,\n\nWelcome to OpenEdu! Your registration was successful.\n\nLogin here: http://example.com/login\n", name)
	msg := []byte(subject + "\n" + body)

	addr := fmt.Sprintf("%s:%s", host, port)
	return smtp.SendMail(addr, auth, from, []string{to}, msg)
}
