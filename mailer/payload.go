package mailer

import (
	"fmt"
	"html/template"

	"cherry/config"
	. "cherry/orm/db/models"
)

type Payload struct {
	Method  string           `json:"method"`
	Locale  string           `json:"locale"`
	Subject string           `json:"subject"`
	Args    []string         `json:"args"`
	Emails  []string         `json:"emails"`
	FuncMap template.FuncMap `json:"-"`
}

func (mp *Payload) foot() template.HTML {
	return template.HTML(fmt.Sprint(I18n.T(mp.Locale, "mailer.footer", map[string]interface{}{"contact_mail": config.CurrentEnv.Email.SupportMail})))
}
