package provider

import (
	"bytes"
	"context"
	"encoding/base64"
	"github.com/k3a/html2text"
	"github.com/sanches1984/gopkg-common/types"
	errors "github.com/sanches1984/gopkg-errors"
	uuid "gopkg.in/satori/go.uuid.v1"
	"html/template"
	"strings"
	"time"
)

var nl2br = strings.NewReplacer("\r\n", "<br/>", "\n", "<br/>")
var html2plain = strings.NewReplacer("\n", "", "\r", "", "<br>", "<br/>", "</br>", "")

// message

type Message struct {
	Subject      string
	To           ContactList
	bodyPlain    string
	bodyHTML     string
	calendarCard *CalendarCard
}

func (m *Message) Prepare(ctx context.Context) error {
	if m.bodyPlain == "" && m.bodyHTML == "" {
		if m.calendarCard != nil {
			m.bodyPlain = "Вам назначена встреча (во вложении)"
		} else {
			return errors.Internal.Err(ctx, "Ошибка при отправке письма").
				WithLogKV("err", "body can't be empty")
		}
	}
	if m.bodyPlain == "" {
		m.bodyPlain = html2text.HTML2Text(html2plain.Replace(m.bodyHTML))
	}
	if m.bodyHTML == "" {
		m.bodyHTML = nl2br.Replace(m.bodyPlain)
	}
	return nil
}

func (m *Message) WithPlain(body string) *Message {
	m.bodyPlain = body
	return m
}

func (m *Message) WithHTML(body string) *Message {
	m.bodyHTML = body
	return m
}

func (m *Message) WithCalendarCard(card *CalendarCard) *Message {
	m.calendarCard = card
	return m
}

// calendar card

var vCardMutliLineReplacer = strings.NewReplacer("\n", "\\n")

type CalendarCard struct {
	UID               string
	Name              string
	Location          string
	AttendeeEmailList []string
	Description       string
	IsAllDay          bool
	Start             time.Time
	Finish            time.Time
}

func (i CalendarCard) GetContent() (string, error) {
	var start, end string
	if i.IsAllDay {
		start = i.Start.Format("20060102")
		end = i.Finish.Format("20060102")
	} else {
		start = types.DateTimeToYMDTHms(i.Start)
		if !i.Finish.IsZero() {
			end = types.DateTimeToYMDTHms(i.Finish)
		}
	}
	var eventUUID string
	if i.UID != "" {
		eventUUID = uuid.NewV3(uuid.Nil, i.UID).String()
	}
	var organizer string
	var attendeeList []string
	if len(i.AttendeeEmailList) != 0 {
		organizer = i.AttendeeEmailList[0]
	}
	if len(i.AttendeeEmailList) > 1 {
		attendeeList = i.AttendeeEmailList[1:]
	}
	bodyTemplate := `BEGIN:VCALENDAR
VERSION:2.0
CALSCALE:GREGORIAN
BEGIN:VEVENT{{if .UUID}}
UID:{{.UUID}}{{end}}
SUMMARY:{{.Name}}{{if .Organizer}}
ORGANIZER:mailto:{{.Organizer}}{{end}}{{range $Attendee := .AttendeeList}}
ATTENDEE;RSVP=TRUE:mailto:{{$Attendee}}{{end}}
DTSTART;TZID=Europe/Moscow:{{.Start}}{{if .End}}
DTEND;TZID=Europe/Moscow:{{.End}}{{end}}
LOCATION:{{.Location}}{{if .Description}}
DESCRIPTION:{{.Description}}{{end}}
STATUS:CONFIRMED
SEQUENCE:3
BEGIN:VALARM
TRIGGER:-PT10M
ACTION:DISPLAY
END:VALARM
END:VEVENT
END:VCALENDAR`
	bodyData := map[string]interface{}{
		"UUID":         eventUUID,
		"Name":         i.Name,
		"Organizer":    organizer,
		"AttendeeList": attendeeList,
		"Description":  vCardMutliLineReplacer.Replace(i.Description),
		"Location":     i.Location,
		"Start":        start,
		"End":          end,
	}
	t := template.Must(template.New("card").Parse(bodyTemplate))
	buf := &bytes.Buffer{}
	if err := t.Execute(buf, bodyData); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}
