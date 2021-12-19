package main

import "github.com/getlantern/systray"

type Queue struct {
	Key         string `json:"key"`
	Display     string `json:"display"`
	Endpoint    string `json:"endpoint"`
	VirtualHost string `json:"virtualHost"`
	Interval    int    `json:"interval"`
}

type Config struct {
	Username string  `json:"username"`
	Password string  `json:"password"`
	Queues   []Queue `json:"queues"`
}

type Item struct {
	Key     string
	Display string
	I       *systray.MenuItem
}

type Message struct {
	Key   string
	value int
}

type RabbitMqMessage struct {
	Messages int    `json:"messages"`
	Name     string `json:"name"`
	Vhost    string `json:"vhost"`
}
