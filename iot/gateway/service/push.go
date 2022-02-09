package service

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

//amqp:// 账号 密码@地址:端口号/vhost
const MQURL = "amqp://guest:guest@10.10.10.137:5672/"
const Router_prop = "prop"
const Router_event = "event"
const Router_alarm = "aievent"
const exchange = "aiot"

var conn *amqp.Connection
var ch *amqp.Channel

func Connect() {
	conn, err := amqp.Dial(MQURL)
	if err != nil {
		log.Fatal(err)
	}
	ch, err = conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	// q, err := ch.QueueDeclare(
	// 	"hello", // name
	// 	false,   // durable
	// 	false,   // delete when unused
	// 	false,   // exclusive
	// 	false,   // no-wait
	// 	nil,     // arguments
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }
}

func Close() {
	if conn != nil {
		conn.Close()
	}
	if ch != nil {
		ch.Close()
	}
}

func PublishData(data interface{}, router string) bool {

	str, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
		return false
	}
	err = ch.Publish(
		exchange, // exchange
		router,   // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "application/json", //application/json,text/plain
			Body:        str,
		})
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}
