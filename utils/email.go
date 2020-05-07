package utils

import (
	"fmt"

	"gopkg.in/gomail.v2"

	envConfig "cherry/config"
)

func SendMail(to string, Subject, bodyMessage string) (err error) {
	d := gomail.NewDialer(envConfig.CurrentEnv.Email.SmtpAddress, envConfig.CurrentEnv.Email.SmtpPort, envConfig.CurrentEnv.Email.SmtpUsername, envConfig.CurrentEnv.Email.SmtpPassword)
	s, err := d.Dial()
	if err != nil {
		panic(err)
	}
	m := gomail.NewMessage()
	m.SetHeader("From", envConfig.CurrentEnv.Email.SystemMailFrom)
	m.SetHeader("To", to)
	m.SetHeader("Subject", Subject)
	m.SetBody("text/html", bodyMessage)
	if err = gomail.Send(s, m); err != nil {
		fmt.Printf("Could not send email to %q: %v", to, err)
	}
	return
}
