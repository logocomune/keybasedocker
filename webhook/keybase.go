package webhook

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/logocomune/keybasedocker/message"
)

const dockerWebhook = "docker-webhook"

type KB struct {
	webHookUlr string
	client     *http.Client
	formatter  formatter
}

type formatter interface {
	EventsToStr(message.EventsGroup) (string, bool)
}

//NewKB Initialize Keybase webhook sender
func NewKB(webHookUlr string, timeOut time.Duration, f formatter) *KB {
	tr := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: timeOut,
		}).DialContext,
	}

	return &KB{
		webHookUlr: webHookUlr,
		client: &http.Client{
			Transport: tr,
			Timeout:   timeOut,
		},
		formatter: f,
	}
}

//Send Send a group of docker events to keybase webhook
func (q *KB) Send(events map[string]message.EventsGroup) {
	msg := ""

	if e, ok := events[dockerWebhook]; ok {
		str, _ := q.formatter.EventsToStr(e)
		msg += str + "\n\n"

		delete(events, dockerWebhook)
	}

	for _, g := range events {
		str, _ := q.formatter.EventsToStr(g)
		msg += str + "\n\n"
	}

	if msg == "" {
		return
	}

	q.toWebHook(msg)
}
func (q *KB) toWebHook(msg string) {
	jsonStr, _ := json.Marshal(struct{ Msg string }{Msg: msg})
	req, err := http.NewRequest("POST", q.webHookUlr, bytes.NewBuffer(jsonStr))

	if err != nil {
		log.Println("Error", err.Error())
	}

	res, err := q.client.Do(req)
	if err == nil {
		log.Println(err)
		return
	}

	io.Copy(ioutil.Discard, res.Body)
	defer res.Body.Close()
}
