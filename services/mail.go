package services

import (
	"github.com/Ekenzy-101/Pentahire-API/config"
	"github.com/gin-gonic/gin"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type MailOption struct {
	To         *mail.Email
	Data       gin.H
	TemplateID string
}

func SendMail(options MailOption) (*rest.Response, error) {
	if config.IsTesting {
		return &rest.Response{StatusCode: 200}, nil
	}

	from := mail.NewEmail("PentaHire", config.SendgridSender)
	personalization := mail.NewPersonalization()
	personalization.AddTos(options.To)
	for key, value := range options.Data {
		personalization.SetDynamicTemplateData(key, value)
	}

	v3Mail := mail.NewV3Mail()
	v3Mail.SetTemplateID(options.TemplateID).SetFrom(from)
	v3Mail.AddPersonalizations(personalization)

	request := sendgrid.GetRequest(config.SendgridAPIKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(v3Mail)
	return sendgrid.API(request)
}
