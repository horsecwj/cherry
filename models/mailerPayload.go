package models

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/spf13/viper"
)

type MailerPayload struct {
	Method  string   `json:"method"`
	Args    []string `json:"args"`
	Locale  string   `json:"locale"`
	Subject string   `json:"subject"`
	Emails  []string `json:"emails"`
	Content string
	FuncMap template.FuncMap `json:"-"`
}

func (mp *MailerPayload) I18nFuncName() string {
	arr := strings.Split(mp.Method, "_")
	for i, ar := range arr {
		arr[i] = strings.Title(ar)
	}
	return strings.Join(arr, "")
}

func (mp *MailerPayload) Activation() {
	mp.Subject = fmt.Sprint(I18n.T(mp.Locale, "token_mailer.activation.subject"))
	mp.Emails = []string{mp.Args[0]}
	mp.Content = mp.Args[1]
	mp.FuncMap = template.FuncMap{
		"title": func() string {
			return fmt.Sprint(I18n.T(mp.Locale, "token_mailer.activation.title"))
		},
		"content": func() string {
			return fmt.Sprint(I18n.T(mp.Locale, "token_mailer.activation.content", map[string]interface{}{"link": mp.Content}))
		},
		"foot": func() template.HTML {
			return template.HTML(fmt.Sprint(I18n.T(mp.Locale, "mailer.footer", map[string]interface{}{"contact_mail": viper.GetString("email.support_email")})))
		},
	}
}

func (mp *MailerPayload) AppResetPassword() {
	mp.Subject = fmt.Sprint(I18n.T(mp.Locale, "token_mailer.app_reset_password.subject"))
	mp.Emails = []string{mp.Args[0]}
	mp.Content = mp.Args[1]
	mp.FuncMap = template.FuncMap{
		"title": func() string {
			return fmt.Sprint(I18n.T(mp.Locale, "token_mailer.app_reset_password.title", map[string]interface{}{"email": mp.Emails[0]}))
		},
		"content": func() string {
			return fmt.Sprint(I18n.T(mp.Locale, "token_mailer.app_reset_password.content", map[string]interface{}{"token": mp.Content}))
		},
		"foot": func() template.HTML {
			return template.HTML(fmt.Sprint(I18n.T(mp.Locale, "mailer.footer", map[string]interface{}{"contact_mail": viper.GetString("email.support_email")})))
		},
	}
}

func (mp *MailerPayload) SetFundPassword() {
	mp.Subject = fmt.Sprint(I18n.T(mp.Locale, "token_mailer.set_fund_password.subject"))
	mp.Emails = []string{mp.Args[0]}
	mp.Content = mp.Args[1]
	mp.FuncMap = template.FuncMap{
		"title": func() string {
			return fmt.Sprint(I18n.T(mp.Locale, "token_mailer.set_fund_password.title", map[string]interface{}{"email": mp.Emails[0]}))
		},
		"content": func() string {
			return fmt.Sprint(I18n.T(mp.Locale, "token_mailer.set_fund_password.content", map[string]interface{}{"token": mp.Content}))
		},
		"foot": func() template.HTML {
			return template.HTML(fmt.Sprint(I18n.T(mp.Locale, "mailer.footer", map[string]interface{}{"contact_mail": viper.GetString("email.support_email")})))
		},
	}
}

func (mp *MailerPayload) BindEmail() {
	mp.Subject = fmt.Sprint(I18n.T(mp.Locale, "token_mailer.bind_email.subject"))
	mp.Emails = []string{mp.Args[0]}
	mp.Content = mp.Args[1]
	mp.FuncMap = template.FuncMap{
		"title": func() string {
			return fmt.Sprint(I18n.T(mp.Locale, "token_mailer.bind_email.title", map[string]interface{}{"email": mp.Emails[0]}))
		},
		"content": func() string {
			return fmt.Sprint(I18n.T(mp.Locale, "token_mailer.bind_email.content", map[string]interface{}{"token": mp.Content}))
		},
		"foot": func() template.HTML {
			return template.HTML(fmt.Sprint(I18n.T(mp.Locale, "mailer.footer", map[string]interface{}{"contact_mail": viper.GetString("email.support_email")})))
		},
	}
}

func (mp *MailerPayload) EmailCodeVerified() {
	mp.Subject = fmt.Sprint(I18n.T(mp.Locale, "token_mailer.email_code_verified.subject"))
	mp.Emails = []string{mp.Args[0]}
	mp.Content = mp.Args[1]
	mp.FuncMap = template.FuncMap{
		"title": func() string {
			return fmt.Sprint(I18n.T(mp.Locale, "token_mailer.email_code_verified.title", map[string]interface{}{"email": mp.Emails[0]}))
		},
		"content": func() string {
			return fmt.Sprint(I18n.T(mp.Locale, "token_mailer.email_code_verified.content", map[string]interface{}{"token": mp.Content}))
		},
		"foot": func() template.HTML {
			return template.HTML(fmt.Sprint(I18n.T(mp.Locale, "mailer.footer", map[string]interface{}{"contact_mail": viper.GetString("email.support_email")})))
		},
	}
}
