package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/getlantern/systray"
)

const (
	URL = "https://bobbish-hedgehog.rmq.cloudamqp.com/api/queues"
)

var (
	//go:embed assets/icon.png
	d           []byte
	config      *Config
	itemMap     map[string]Item
	messageChan chan Message
	itemChan    chan string
)

// Call RabbitMQ API in an infinite Time Ticker loop
func callApi(queue Queue) {
	for range time.NewTicker(time.Duration(queue.Interval) * time.Second).C {
		if !itemMap[queue.Key].I.Checked() {
			continue
		}
		client := http.Client{Timeout: 4 * time.Second}
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%v/%v/%v", URL, queue.VirtualHost, queue.Endpoint), http.NoBody)
		if err != nil {
			log.Printf("ERROR: Not able to create new request: %v\n", err)
			continue
		}
		req.SetBasicAuth(config.Username, config.Password)
		res, err := client.Do(req)
		if err != nil {
			log.Printf("ERROR: Calling API: %v\n", err)
			continue
		}
		if res.StatusCode == http.StatusUnauthorized {
			go Alert("RabbitMQ Notify Unauthorized", "Please check if you have access to RabbitMq or if username/password in config file is correct")
			systray.Quit()
		}

		var msg RabbitMqMessage
		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			log.Printf("ERROR: Converting response body: %v\n", err)
		} else {
			json.Unmarshal(resBody, &msg)
			messageChan <- Message{queue.Key, msg.Messages}
		}
		res.Body.Close()
	}
}

func createItems() {
	for _, queue := range config.Queues {
		item := systray.AddMenuItemCheckbox(queue.Display, queue.Display, true)
		itemMap[queue.Key] = Item{queue.Key, queue.Display, item}
		go callApi(queue)
		go func(item *systray.MenuItem, key string) {
			for range item.ClickedCh {
				itemChan <- key
			}
		}(item, queue.Key)
	}
}

func onReady() {
	systray.SetTemplateIcon(d, d)

	createItems()
	quit := systray.AddMenuItem("Quit", "Quit App")
	for {
		select {
		case val := <-messageChan:
			itemMap[val.Key].I.SetTitle(fmt.Sprintf("%v: %v", itemMap[val.Key].Display, val.value))
		case key := <-itemChan:
			if itemMap[key].I.Checked() {
				itemMap[key].I.SetTitle(itemMap[key].Display)
				itemMap[key].I.Uncheck()
			} else {
				itemMap[key].I.Check()
			}
		case <-quit.ClickedCh:
			systray.Quit()
		}
	}
}

func main() {
	// Setup log file
	LoggerSetup()

	log.Println("Starting App")

	var err error
	config, err = ReadConfig()
	if err != nil {
		log.Fatal(err)
	}

	itemMap = make(map[string]Item)
	messageChan = make(chan Message, len(config.Queues))
	itemChan = make(chan string, len(config.Queues))

	onExit := func() {
		log.Println("Quitting App")
	}

	systray.Run(onReady, onExit)
}
