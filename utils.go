package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
)

func LoggerSetup() {
	homeDir, _ := os.UserHomeDir()
	file, err := os.Create(path.Join(homeDir, ".rabbitmq.log"))
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
}

func ReadConfig() (*Config, error) {
	homeDir, _ := os.UserHomeDir()

	b, err := ioutil.ReadFile(path.Join(homeDir, ".rabbitmq.config"))
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	var data Config
	err = json.Unmarshal(b, &data)
	if err != nil {
		return nil, fmt.Errorf("error unmarshal config data: %v", err)
	}
	return &data, nil
}

func Notify(title, message string) {
	osa, err := exec.LookPath("osascript")
	if err != nil {
		fmt.Println("Notification is not working:", err)
	}

	cmd := exec.Command(osa, "-e", `display notification "`+message+`" with title "`+title+`"`)
	cmd.Run()
}

func Alert(title, message string) {
	osa, err := exec.LookPath("osascript")
	if err != nil {
		fmt.Println("Notification is not working:", err)
	}

	cmd := exec.Command(osa, "-e", `display alert "`+title+`" message "`+message+`" as critical buttons {"Ok"}`)
	cmd.Run()
}
