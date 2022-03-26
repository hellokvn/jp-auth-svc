package sender

import (
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
)

type SendMailBody struct {
	Template string `json:"template"`
	To       string `json:"to"`
}

type Handler struct {
	MailSvc *amqp.Channel
}

func Init(url string) Handler {
	connectRabbitMQ, err := amqp.Dial(url)

	if err != nil {
		panic(err)
	}

	channelRabbitMQ, err := connectRabbitMQ.Channel()

	if err != nil {
		panic(err)
	}

	_, err = channelRabbitMQ.QueueDeclare(
		"QueueService1",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		panic(err)
	}

	return Handler{channelRabbitMQ}
}

func (h Handler) SendMail(b *SendMailBody) {
	body, _ := json.Marshal(&b)

	message := amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
		Type:        "send_mail",
	}

	if err := h.MailSvc.Publish(
		"",
		"QueueService1",
		false,
		false,
		message,
	); err != nil {
		fmt.Println(err)
	}
}
