package mailer

import (
	"fmt"
	"html/template"
	"strings"

	"cherry/config"
	. "cherry/orm/db/models"
)

func (mp *Payload) I18nFuncName() string {
	arr := strings.Split(mp.Method, "_")
	for i, ar := range arr {
		arr[i] = strings.Title(ar)
	}
	return strings.Join(arr, "")
}

func (mp *Payload) Activation() {
	mp.Subject = fmt.Sprint(I18n.T(mp.Locale, "token_mailer.activation.subject"))
	mp.Emails = []string{mp.Args[0]}
	mp.FuncMap = template.FuncMap{
		"title": func() string {
			return fmt.Sprint(I18n.T(mp.Locale, "token_mailer.activation.title"))
		},
		"content": func() string {
			return fmt.Sprint(I18n.T(mp.Locale, "token_mailer.activation.content", map[string]interface{}{"link": mp.Args[1]}))
		},
		"foot": func() template.HTML {
			return template.HTML(fmt.Sprint(I18n.T(mp.Locale, "mailer.footer", map[string]interface{}{"contact_mail": config.CurrentEnv.Email.SupportMail})))
		},
	}
}

func (mp *Payload) AppResetPassword() {
	mp.Subject = fmt.Sprint(I18n.T(mp.Locale, "token_mailer.app_reset_password.subject"))
	mp.Emails = []string{mp.Args[0]}
	mp.FuncMap = template.FuncMap{
		"title": func() string {
			return fmt.Sprint(I18n.T(mp.Locale, "token_mailer.app_reset_password.title", map[string]interface{}{"email": mp.Emails[0]}))
		},
		"content": func() string {
			return fmt.Sprint(I18n.T(mp.Locale, "token_mailer.app_reset_password.content", map[string]interface{}{"token": mp.Args[1]}))
		},
		"foot": mp.foot,
	}
}

func (mp *Payload) SetFundPassword() {
	mp.Subject = fmt.Sprint(I18n.T(mp.Locale, "token_mailer.set_fund_password.subject"))
	mp.Emails = []string{mp.Args[0]}
	mp.FuncMap = template.FuncMap{
		"title": func() string {
			return fmt.Sprint(I18n.T(mp.Locale, "token_mailer.set_fund_password.title", map[string]interface{}{"email": mp.Emails[0]}))
		},
		"content": func() string {
			return fmt.Sprint(I18n.T(mp.Locale, "token_mailer.set_fund_password.content", map[string]interface{}{"token": mp.Args[1]}))
		},
		"foot": mp.foot,
	}
}

func (mp *Payload) BindEmail() {
	mp.Subject = fmt.Sprint(I18n.T(mp.Locale, "token_mailer.bind_email.subject"))
	mp.Emails = []string{mp.Args[0]}
	mp.FuncMap = template.FuncMap{
		"title": func() string {
			return fmt.Sprint(I18n.T(mp.Locale, "token_mailer.bind_email.title", map[string]interface{}{"email": mp.Emails[0]}))
		},
		"content": func() string {
			return fmt.Sprint(I18n.T(mp.Locale, "token_mailer.bind_email.content", map[string]interface{}{"token": mp.Args[1]}))
		},
		"foot": mp.foot,
	}
}

func (mp *Payload) EmailCodeVerified() {
	mp.Subject = fmt.Sprint(I18n.T(mp.Locale, "token_mailer.email_code_verified.subject"))
	mp.Emails = []string{mp.Args[0]}
	mp.FuncMap = template.FuncMap{
		"title": func() string {
			return fmt.Sprint(I18n.T(mp.Locale, "token_mailer.email_code_verified.title", map[string]interface{}{"email": mp.Emails[0]}))
		},
		"content": func() string {
			return fmt.Sprint(I18n.T(mp.Locale, "token_mailer.email_code_verified.content", map[string]interface{}{"token": mp.Args[1]}))
		},
		"foot": mp.foot,
	}
}
