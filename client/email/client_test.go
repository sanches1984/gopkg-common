// +build payment

package email

import (
	"context"
	"github.com/rs/zerolog"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/sanches1984/gopkg-common/client/email/provider"
	"github.com/stretchr/testify/assert"
)

func TestSendEmail(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		panic("Failed to load env variables, error: " + err.Error())
	}

	to := []string{
		os.Getenv("EMAIL_TO1_ADDRESS"), os.Getenv("EMAIL_TO1_NAME"),
		os.Getenv("EMAIL_TO2_ADDRESS"), os.Getenv("EMAIL_TO2_NAME"),
	}

	client, _, err := NewClient(
		provider.NewSendGridProvider(os.Getenv("EMAIL_SEND_GRID_KEY"), zerolog.Nop()),
		os.Getenv("EMAIL_FROM_ADDRESS"),
		os.Getenv("EMAIL_FROM_NAME"),
	)
	assert.Nil(t, err)
	testClient(t, to, client)
}

func testClient(t *testing.T, to []string, email *Client) {
	t.Run("Send plain", func(t *testing.T) {
		err := email.Send(context.Background(), "subject", map[string]string{"email@yandex.ru": "user name"}, NewMessage().WithPlain("test"))
		assert.Nil(t, err)
	})

	t.Run("Send invite", func(t *testing.T) {
		err := email.Send(context.Background(), "subject", map[string]string{"email@yandex.ru": "user name"}, NewMessage().
			WithCalendarCard(&provider.CalendarCard{
				Name:        "Супервстреча",
				Location:    "Переговорка 1",
				Description: "Описание встречи 111",
				Start:       time.Now().Add(3 * time.Hour),
				Finish:      time.Now().Add(4 * time.Hour),
			}))
		assert.Nil(t, err)
	})

	t.Run("Send invite all day", func(t *testing.T) {
		dt := time.Date(2020, 5, 2, 0, 0, 0, 0, time.UTC)
		err := email.Send(context.Background(), "subject", map[string]string{"email@yandex.ru": "user name"}, NewMessage().
			WithCalendarCard(&provider.CalendarCard{
				Name:        "Супервстреча на весь день",
				Location:    "Переговорка 1",
				Description: "Описание встречи 111",
				IsAllDay:    true,
				Start:       dt,
				Finish:      dt,
			}))
		assert.Nil(t, err)
	})

	t.Run("Send email with form", func(t *testing.T) {
		html := `<h1 style="font: 22px arial; text-align: center; margin: 0 0 10px 0; padding: 0;">Ваша оценка</h1>
<form style="width: 99%;"
      method="get" action="http://127.0.0.1/test/form">
<ul style="list-style-type: none; margin: 0 auto 2em auto; padding: 0; overflow: hidden; zoom: 1;">
	<li style="float: left; margin-right: 2em;"><label><input type="radio" name="rate" value="1"/> <span style="font-size: 2em;">👍</span></label></li>
	<li style="float: left; margin-right: 2em;"><label><input type="radio" name="rate" value="2"/> <span style="font-size: 2em;">😍</span></label></li>
	<li style="float: left; margin-right: 2em;"><label><input type="radio" name="rate" value="3" checked/> <span style="font-size: 2em;">😐</span></label></li>
	<li style="float: left; margin-right: 2em;"><label><input type="radio" name="rate" value="4"/> <span style="font-size: 2em;">😊</span></label></li>
	<li style="float: left; margin-right: 2em;"><label><input type="radio" name="rate" value="5"/> <span style="font-size: 2em;">😬</span></label></li>
</ul>
<div style="margin: 0 0 1em 0">Прокомментируйте отчет сотрудника</div>
<textarea style="background: #ddd; padding: .8em 0; display: block; width: 100%; height: 80px; border-style: none; margin-bottom: 1em;" 
          name="comment"></textarea>
<button style="background: #060; color: #fff; border-style: none; padding: 1em 0; text-align: center; display: block; width: 100%;" 
        type="submit">Отправить</button>
</form>`

		err := email.Send(context.Background(), "subject", map[string]string{"email@yandex.ru": "user name"}, NewMessage().WithHTML(html))
		assert.Nil(t, err)
	})
}
