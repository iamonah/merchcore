package mailer

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"time"

	"gopkg.in/gomail.v2"
)

const (
	numOfEmailRetries = 3
)

//go:embed templates
var templateFS embed.FS

type Mail struct {
	dialer *gomail.Dialer
	sender string
}

func NewMailTrap(sender, stmthost, password, username string, smtpport int) *Mail {
	dialer := gomail.NewDialer(stmthost, smtpport, username, password)

	return &Mail{
		dialer: dialer,
		sender: sender,
	}
}

func (m *Mail) Send(templateFile string, reciever string, data any) error {
	var err error
	message := gomail.NewMessage()

	template, err := template.ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return fmt.Errorf("parsetemplate: %w", err)
	}

	subjectBuffer := bytes.NewBuffer([]byte{})
	if err = template.ExecuteTemplate(subjectBuffer, "subject", data); err != nil {
		return fmt.Errorf("subjecttemplate: %w", err)
	}

	plainbodyBuffer := bytes.NewBuffer([]byte{})
	if err = template.ExecuteTemplate(plainbodyBuffer, "plainBody", data); err != nil {
		return fmt.Errorf("plainBodytemplate: %w", err)
	}

	htmlbodyBuffer := bytes.NewBuffer([]byte{})
	if err = template.ExecuteTemplate(htmlbodyBuffer, "htmlBody", data); err != nil {
		return fmt.Errorf("parsehtmlBodytemplate: %w", err)
	}

	message.SetHeader("From", m.sender)
	message.SetHeader("To", reciever)
	message.SetHeader("Subject", subjectBuffer.String())
	message.SetHeader("List-Unsubscribe", fmt.Sprintf(`<mailto:unsubscribe@yourapp.com?subject=unsub-%s>, <https://yourapp.com/unsubscribe?email=%s>`, reciever, reciever))  // Add this header
	message.SetBody("text/plain", plainbodyBuffer.String())

	message.AddAlternative("text/html", htmlbodyBuffer.String())

	for i := 0; i < numOfEmailRetries; i++ {
		if err = m.dialer.DialAndSend(message); err != nil {
			time.Sleep(time.Millisecond * time.Duration(numOfEmailRetries))
			continue
		}
		return nil
	}

	return fmt.Errorf("failed to send email after %d retries : %w", numOfEmailRetries, err)
}