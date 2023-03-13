package mailer

import (
	"bytes"
	"fmt"
	"log"
	"text/template"
	"time"

	"github.com/izaakdale/service-event-order/pkg/proto/order"
	"github.com/kelseyhightower/envconfig"
	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail"
)

type specification struct {
	Host     string
	Port     int
	User     string
	Password string
}

var (
	client *mail.SMTPClient
	spec   specification
)

func init() {
	envconfig.Process("MAIL", &spec)

	srv := mail.SMTPServer{
		Host:           spec.Host,
		Port:           spec.Port,
		Username:       spec.User,
		Password:       spec.Password,
		Encryption:     mail.EncryptionNone,
		KeepAlive:      true,
		ConnectTimeout: 10 * time.Second,
		SendTimeout:    10 * time.Second,
	}

	var err error
	client, err = srv.Connect()
	if err != nil {
		panic(err)
	}
}

func Send(recipient, orderID string, tickets []*order.Ticket) error {
	fmtd, err := buildHTMLMessage(orderID, tickets)
	if err != nil {
		return err
	}

	email := mail.NewMSG()
	email.SetFrom("awesome_tickets@test.com").AddTo(recipient).SetSubject("your tickets...")
	email.AddAlternative(mail.TextHTML, fmtd)

	for _, v := range tickets {
		email.AddAttachment(fmt.Sprintf("tmp/%s.jpg", v.TicketId), fmt.Sprintf("ticket#%s", v.TicketId))
	}

	err = email.Send(client)
	if err != nil {
		return err
	}

	return nil
}

func buildHTMLMessage(orderID string, tickets []*order.Ticket) (string, error) {
	templateToRender := "./internal/mailer/mail-templates/mail.html.gohtml"

	t, err := template.New("email-html").ParseFiles(templateToRender)
	if err != nil {
		log.Println("Error parsing html file template")
		return "", err
	}

	data := map[string]any{
		"order":   orderID,
		"tickets": tickets,
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", data); err != nil {
		log.Println("Error executing template")
		return "", err
	}

	fmtdMsg := tpl.String()
	fmtdMsg, err = InlineCSS(fmtdMsg)
	if err != nil {
		log.Println("Error CSS")
		return "", nil
	}

	return fmtdMsg, nil
}

// func buildPlainTextMsg(msg Message) (string, error) {
// 	templateToRender := "./templates/mail.plain.gohtml"

// 	t, err := template.New("email-plain").ParseFiles(templateToRender)
// 	if err != nil {
// 		log.Println("Error Parsing")
// 		return "", err
// 	}

// 	var tpl bytes.Buffer
// 	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
// 		log.Println("Error Executing template")
// 		return "", err
// 	}

// 	plainMsg := tpl.String()

// 	return plainMsg, nil
// }

func InlineCSS(s string) (string, error) {
	opts := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}

	prem, err := premailer.NewPremailerFromString(s, &opts)
	if err != nil {
		log.Println("Error Premailer")
		return "", err
	}

	html, err := prem.Transform()
	if err != nil {
		log.Println("Error Transform")
		return "", nil
	}

	return html, nil
}
