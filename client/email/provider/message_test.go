package provider

import (
	"context"
	"encoding/base64"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPrepareMessage(t *testing.T) {
	t.Run("Plain to HTML", func(t *testing.T) {
		msg := &Message{
			bodyPlain: `Hello!
Second, string

Bye`,
		}
		err := msg.Prepare(context.Background())
		assert.Nil(t, err)
		assert.Equal(t, "Hello!<br/>Second, string<br/><br/>Bye", msg.bodyHTML)
	})

	t.Run("HTML to plain", func(t *testing.T) {
		msg := &Message{
			bodyHTML: `<h1>Hello!</h1>
<p>Second, string</p>
Some other string<br/>

Second other string<br>


Bye`,
		}
		err := msg.Prepare(context.Background())
		assert.Nil(t, err)
		assert.Equal(t, "Hello!\r\n\r\nSecond, string\r\n\r\nSome other string\r\nSecond other string\r\nBye", msg.bodyPlain)
	})
}

func TestCalendarCard(t *testing.T) {
	cs := []struct {
		name   string
		card   CalendarCard
		result string
	}{
		{
			name: "empty",
			card: CalendarCard{},
			result: `BEGIN:VCALENDAR
VERSION:2.0
CALSCALE:GREGORIAN
BEGIN:VEVENT
SUMMARY:
DTSTART;TZID=Europe/Moscow:00010101T000000
LOCATION:
STATUS:CONFIRMED
SEQUENCE:3
BEGIN:VALARM
TRIGGER:-PT10M
ACTION:DISPLAY
END:VALARM
END:VEVENT
END:VCALENDAR`,
		},
		{
			name: "filled",
			card: CalendarCard{
				UID:               "tested",
				Name:              "Называние встречи",
				Location:          "Место проведения встречи",
				AttendeeEmailList: []string{"first@mail.ru", "second@mail.ru", "third@mail.ru"},
				Description:       "Описание встречи\nВторая строка",
				IsAllDay:          false,
				Start:             time.Date(2020, 1, 1, 12, 15, 0, 0, time.UTC),
				Finish:            time.Date(2020, 1, 1, 13, 45, 0, 0, time.UTC),
			},
			result: `BEGIN:VCALENDAR
VERSION:2.0
CALSCALE:GREGORIAN
BEGIN:VEVENT
UID:99e9aaf2-889d-3ff2-a377-bc8230a762ea
SUMMARY:Называние встречи
ORGANIZER:mailto:first@mail.ru
ATTENDEE;RSVP=TRUE:mailto:second@mail.ru
ATTENDEE;RSVP=TRUE:mailto:third@mail.ru
DTSTART;TZID=Europe/Moscow:20200101T121500
DTEND;TZID=Europe/Moscow:20200101T134500
LOCATION:Место проведения встречи
DESCRIPTION:Описание встречи\nВторая строка
STATUS:CONFIRMED
SEQUENCE:3
BEGIN:VALARM
TRIGGER:-PT10M
ACTION:DISPLAY
END:VALARM
END:VEVENT
END:VCALENDAR`,
		},
	}
	for _, c := range cs {
		gotResult, err := c.card.GetContent()
		assert.Nil(t, err)
		assert.Equal(t,
			base64.StdEncoding.EncodeToString([]byte(c.result)),
			gotResult,
			c.name)
	}
}
