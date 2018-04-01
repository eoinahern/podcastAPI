package util

import (
	"bytes"
	"html/template"
	"log"
	"net/smtp"

	"github.com/eoinahern/my_podcast_api/models"
)

const (
	mime string = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	//subject string = "Subject: Registration\n"
)

//MailRequest : send a mail to a specific email using html template.
type MailRequest struct {
	SenderId     string
	ToId         string
	BodyLocation string
	BodyParams   *models.TemplateParams
}

//SendMail : send a mail via smtp server using plainauth
func (m *MailRequest) SendMail() (bool, error) {

	smtpConf := &models.SmtpConfig{}
	smtpConf.ReadFromFile("smtpConfig.json")
	auth := smtp.PlainAuth("", smtpConf.Username, smtpConf.Password, smtpConf.Server)
	err := smtp.SendMail(smtpConf.Server+":"+smtpConf.Port, auth, m.SenderId, []string{m.ToId}, []byte(m.buildMail()))

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
	err = template.Execute(buf, m.BodyParams)

	if err != nil {
		log.Println(err)
		return ""
	}

	return buf.String()
}
