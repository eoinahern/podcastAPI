package util

import (
	"bytes"
	"html/template"
	"log"
	"net/smtp"

	"github.com/eoinahern/podcastAPI/models"
)

//MailRequestInt interface facilitate mocking
type MailRequestInt interface {
	SendMail() (bool, error)
	SetBodyParams(params *models.TemplateParams)
	SetToID(toid string)
}

const (
	mime string = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	//subject string = "Subject: Registration\n"
)

//MailRequest : send a mail to a specific email using html template.
type MailRequest struct {
	SenderId     string
	toId         string
	BodyLocation string
	bodyParams   *models.TemplateParams
}

//SetBodyParams set bodyParams
func (m *MailRequest) SetBodyParams(bodyParams *models.TemplateParams) {
	m.bodyParams = bodyParams
}

//SetToID set toId
func (m *MailRequest) SetToID(toid string) {
	m.toId = toid
}

//SendMail : send a mail via smtp server using plainauth
func (m *MailRequest) SendMail() (bool, error) {

	smtpConf := &models.SmtpConfig{}
	smtpConf.ReadFromFile("config/smtpConfig.json")
	auth := smtp.PlainAuth("", smtpConf.Username, smtpConf.Password, smtpConf.Server)
	err := smtp.SendMail(smtpConf.Server+":"+smtpConf.Port, auth, m.SenderId, []string{m.toId}, []byte(m.buildMail()))

	if err != nil {
		return false, err
	}

	return true, nil
}

//internal helper method. just build the mail string
func (m *MailRequest) buildMail() string {
	return mime + m.constructTemplate()
}

//helper. create the template from bodyloaction and bodyParams
func (m *MailRequest) constructTemplate() string {

	template, err := template.ParseFiles(m.BodyLocation)

	if err != nil {
		log.Println(err)
		return ""
	}

	buf := new(bytes.Buffer)
	err = template.Execute(buf, m.bodyParams)

	if err != nil {
		log.Println(err)
		return ""
	}

	return buf.String()
}
