package models

import (
	"fmt"
	"html/template"
	"strings"

	envConfig "cherry/config"
)

type MailerPayload struct {
	MailerClass string   `json:"mailer_class"`
	Method      string   `json:"method"`
	Args        []string `json:"args"`
	Locale      string   `json:"locale"`
	Subject     string   `json:"subject"`
	Email       string   `json:"email"`
	Times       int      `json:"times"`
	Content     string
	FuncMap     template.FuncMap
}

func (mp *MailerPayload) I18nFuncName() string {
	arr := strings.Split(mp.Method, "_")
	for i, ar := range arr {
		arr[i] = strings.Title(ar)
	}
	return strings.Join(arr, "")
}

func (mp *MailerPayload) AppResetPassword() {
	mp.Subject = fmt.Sprint(I18n.T(mp.Locale, "token_mailer.app_reset_password.subject"))
	mp.Content = mp.Args[0]
	mp.Subject = fmt.Sprint(I18n.T(mp.Locale, "token_mailer.app_reset_password.subject", map[string]interface{}{"token": mp.Content}))
	mp.FuncMap = template.FuncMap{
		"hi": func() string {
			return fmt.Sprint(I18n.T(mp.Locale, "token_mailer.app_reset_password.hi", map[string]interface{}{"email": mp.Email}))
		},
		"follow_token": func() string {
			return fmt.Sprint(I18n.T(mp.Locale, "token_mailer.app_reset_password.follow_token", map[string]interface{}{"token": mp.Content}))
		},
		"foot": func() template.HTML {
			return template.HTML(fmt.Sprint(I18n.T(mp.Locale, "mailer.footer", map[string]interface{}{"contact_mail": envConfig.CurrentEnv.Email.SupportMail})))
		},
	}
}

func (mp *MailerPayload) BindEmail() {
	mp.Subject = fmt.Sprint(I18n.T(mp.Locale, "token_mailer.bind_email.subject"))
	mp.Content = mp.Args[0]
	mp.Subject = fmt.Sprint(I18n.T(mp.Locale, "token_mailer.bind_email.subject", map[string]interface{}{"token": mp.Content}))
	mp.FuncMap = template.FuncMap{
		"hi": func() string {
			return fmt.Sprint(I18n.T(mp.Locale, "token_mailer.bind_email.hi", map[string]interface{}{"email": mp.Email}))
		},
		"follow_token": func() string {
			return fmt.Sprint(I18n.T(mp.Locale, "token_mailer.bind_email.follow_token", map[string]interface{}{"token": mp.Content}))
		},
		"foot": func() template.HTML {
			return template.HTML(fmt.Sprint(I18n.T(mp.Locale, "mailer.footer", map[string]interface{}{"contact_mail": envConfig.CurrentEnv.Email.SupportMail})))
		},
	}
}

func (mp *MailerPayload) SetFundPassword() {
	mp.Subject = fmt.Sprint(I18n.T(mp.Locale, "token_mailer.set_fund_password.subject"))
	mp.Content = mp.Args[0]
	mp.Subject = fmt.Sprint(I18n.T(mp.Locale, "token_mailer.set_fund_password.subject", map[string]interface{}{"token": mp.Content}))
	mp.FuncMap = template.FuncMap{
		"hi": func() string {
			return fmt.Sprint(I18n.T(mp.Locale, "token_mailer.set_fund_password.hi", map[string]interface{}{"email": mp.Email}))
		},
		"follow_token": func() string {
			return fmt.Sprint(I18n.T(mp.Locale, "token_mailer.set_fund_password.follow_token", map[string]interface{}{"token": mp.Content}))
		},
		"foot": func() template.HTML {
			return template.HTML(fmt.Sprint(I18n.T(mp.Locale, "mailer.footer", map[string]interface{}{"contact_mail": envConfig.CurrentEnv.Email.SupportMail})))
		},
	}
}

func (mp *MailerPayload) EmailCodeVerified() {
	mp.Subject = fmt.Sprint(I18n.T(mp.Locale, "token_mailer.email_code_verified.subject"))
	mp.Content = mp.Args[0]
	mp.Subject = fmt.Sprint(I18n.T(mp.Locale, "token_mailer.email_code_verified.subject", map[string]interface{}{"token": mp.Content}))
	mp.FuncMap = template.FuncMap{
		"hi": func() string {
			return fmt.Sprint(I18n.T(mp.Locale, "token_mailer.email_code_verified.hi", map[string]interface{}{"email": mp.Email}))
		},
		"follow_token": func() string {
			return fmt.Sprint(I18n.T(mp.Locale, "token_mailer.email_code_verified.follow_token", map[string]interface{}{"token": mp.Content}))
		},
		"foot": func() template.HTML {
			return template.HTML(fmt.Sprint(I18n.T(mp.Locale, "mailer.footer", map[string]interface{}{"contact_mail": envConfig.CurrentEnv.Email.SupportMail})))
		},
	}
}
