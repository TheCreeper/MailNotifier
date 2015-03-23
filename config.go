package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (

	// Default app icon in the notification
	DefaultNotificationIcon = "mail-unread"

	// Default sound
	DefaultNotificationSound = "message-new-email"

	// File prefix for specifying files in the config
	FilePrefix = "file://"
)

type ClientConfig struct {

	// How often to query the POP3 server
	CheckFrequency int

	// The path to the notification icon. Will be set to the default if nothing is specified
	NotificationIcon string

	// The path to the notification sound. Default will be used if nothing is set
	NotificationSound string

	Proxys   []Proxy
	Accounts []Account
}

type Proxy struct {
	Name     string
	Address  string
	User     string
	Password string
}

type Account struct {

	// Connection Info
	Address  string
	User     string
	Password string
	Proxy    string
}

func GenerateConfig() ([]byte, error) {

	cfg := &ClientConfig{

		CheckFrequency:    20,
		NotificationIcon:  DefaultNotificationIcon,
		NotificationSound: DefaultNotificationSound,

		Proxys: []Proxy{

			{
				Name:     "tor",
				Address:  "127.0.0.1:9050",
				User:     "",
				Password: "",
			},
		},

		Accounts: []Account{

			{
				Address:  "pop3.riseup.net:995",
				User:     "example",
				Password: "password",
				Proxy:    "",
			},
		},
	}
	return json.MarshalIndent(cfg, "", "	")
}

func (cfg *ClientConfig) GetProxyInfo(n string) (Proxy, error) {

	for _, v := range cfg.Proxys {

		if v.Name != n {

			continue
		}

		return v, nil
	}

	return Proxy{}, fmt.Errorf("No proxy found by that name for %s", n)
}

func (cfg *ClientConfig) Validate() (err error) {

	// Set to default if icon not specified
	if len(cfg.NotificationIcon) == 0 {

		cfg.NotificationIcon = DefaultNotificationIcon
	}

	if strings.HasPrefix(cfg.NotificationIcon, FilePrefix) {

		trimmed := strings.TrimPrefix(cfg.NotificationIcon, FilePrefix)
		cleaned := filepath.Clean(trimmed)

		cfg.NotificationIcon = os.ExpandEnv(cleaned)
	}

	// Set to default if sound not specified
	if len(cfg.NotificationSound) == 0 {

		cfg.NotificationSound = DefaultNotificationSound
	}

	if strings.HasPrefix(cfg.NotificationSound, FilePrefix) {

		trimmed := strings.TrimPrefix(cfg.NotificationSound, FilePrefix)
		cleaned := filepath.Clean(trimmed)

		cfg.NotificationSound = os.ExpandEnv(cleaned)
	}

	return
}

func GetCFG() (cfg *ClientConfig, err error) {

	b, err := ioutil.ReadFile(ConfigFile)
	if err != nil {

		return
	}

	err = json.Unmarshal(b, &cfg)
	if err != nil {

		return
	}

	if err = cfg.Validate(); err != nil {

		return
	}

	return
}
