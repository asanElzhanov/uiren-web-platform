package yandex_sender

import (
	"errors"
	"fmt"
	"net/smtp"
	"uiren/pkg/logger"

	"github.com/jordan-wright/email"
)

const (
	smtpAuthAddress   = "smtp.yandex.kz"
	smtpServerAddress = "smtp.yandex.kz:587"
)

var (
	ErrFailedToAttachFile = errors.New("failed to attach file")
)

type YandexSender struct {
	name              string
	fromEmailAddress  string
	fromEmailPassword string
}

var sender *YandexSender

func Init(name, fromEmailAddress, fromEmailPassword string) {
	sender = &YandexSender{
		name:              name,
		fromEmailAddress:  fromEmailAddress,
		fromEmailPassword: fromEmailPassword,
	}
	logger.Info("yandex email sender initialized")
}

func SendEmail(
	subject string,
	content string,
	to []string,
	cc []string,
	bcc []string,
	attachFiles []string,
) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", sender.name, sender.fromEmailAddress)
	e.Subject = subject
	e.HTML = []byte(content)
	e.To = to
	e.Cc = cc
	e.Bcc = bcc

	for _, f := range attachFiles {
		_, err := e.AttachFile(f)
		if err != nil {
			return ErrFailedToAttachFile
		}
	}

	smtpAuth := smtp.PlainAuth("", sender.fromEmailAddress, sender.fromEmailPassword, smtpAuthAddress)
	return e.Send(smtpServerAddress, smtpAuth)
}
