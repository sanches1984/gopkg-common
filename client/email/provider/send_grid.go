package provider

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/sanches1984/gopkg-errors"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type sendGridProvider struct {
	api    *sendgrid.Client
	from   *mail.Email
	logger zerolog.Logger
}

func NewSendGridProvider(key string, logger zerolog.Logger) IProvider {
	return &sendGridProvider{
		api:    sendgrid.NewSendClient(key),
		logger: logger,
	}
}

func (c sendGridProvider) Connect(fromAddress, fromName string) (ISender, func() error, error) {
	c.from = mail.NewEmail(fromName, fromAddress)
	return c, nil, nil
}

func (c sendGridProvider) Send(ctx context.Context, msg *Message) error {
	m := mail.NewV3Mail()
	e := mail.NewEmail(c.from.Name, c.from.Address)
	m.SetFrom(e)

	m.Subject = msg.Subject
	p := mail.NewPersonalization()
	toList := make([]*mail.Email, 0, len(msg.To))
	used := make(map[string]struct{}, len(msg.To))
	for _, contact := range msg.To {
		if contact.Address == "" {
			continue
		}
		if _, ok := used[contact.Address]; !ok {
			toList = append(toList, mail.NewEmail(contact.Name, contact.Address))
			used[contact.Address] = struct{}{}
		}
	}

	if len(toList) == 0 {
		return nil
	}

	p.AddTos(toList...)
	m.AddPersonalizations(p)

	if msg.calendarCard != nil {
		content, err := msg.calendarCard.GetContent()
		if err != nil {
			return errors.Internal.ErrWrap(ctx, "Can't create calendar card content by template", err)
		}
		m.Attachments = make([]*mail.Attachment, 0, 1)
		m.Attachments = append(m.Attachments, &mail.Attachment{
			Content:     content,
			Type:        "application/ics",
			Name:        "invite.ics",
			Filename:    "invite.ics",
			Disposition: "attachment",
		})
	}

	content := mail.NewContent("text/plain", msg.bodyPlain)
	m.AddContent(content)
	content = mail.NewContent("text/html", msg.bodyHTML)
	m.AddContent(content)

	response, err := c.api.Send(m)
	if err != nil {
		c.logger.Error().Err(err).Str("to", msg.To.String()).Str("subject", msg.Subject).Msg("send email error")
		return errors.Internal.Err(ctx, "Ошибка при отправке письма").
			WithLogKV("err", err.Error())
	}
	if response.StatusCode >= 400 {
		c.logger.Error().Err(err).
			Str("to", msg.To.String()).
			Str("subject", msg.Subject).
			Int("code", response.StatusCode).
			Str("body", response.Body).
			Msg("response email error")

		return errors.Internal.Err(ctx, "Ошибка при отправке письма").
			WithLogKV("status", response.StatusCode, "body", response.Body)
	}

	c.logger.Info().Str("to", msg.To.String()).Str("subject", msg.Subject).Msg("email sent successfully")
	return nil
}
