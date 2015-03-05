package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type ClientConfig struct {

	// How often to query the POP3 server
	CheckFrequency int

	Proxys   []Proxy
	Accounts []Account
}

type Proxy struct {
	Name     string
	Type     string
	Address  string
	Username string
	Password string
	Timeout  int
}

type Account struct {

	// Connection Info
	Host     string
	Username string
	Password string
	Proxy    string
}

func GenerateConfig() ([]byte, error) {

	cfg := &ClientConfig{

		CheckFrequency: 1,

		Proxys: []Proxy{

			{
				Name:     "tor",
				Type:     "socks5",
				Address:  "127.0.0.1:9050",
				Username: "",
				Password: "",
				Timeout:  5,
			},
		},

		Accounts: []Account{

			{
				Host:     "pop3.riseup.net:995",
				Username: "example",
				Password: "password",
				Proxy:    "",
			},
		},
	}
	return json.MarshalIndent(cfg, "", "	")
}

func (cfg *ClientConfig) GetProxyInfo(n string) (Proxy, error) {

	for _, v := range cfg.Proxys {

		if !(v.Name == n) {

			continue
		}

		return v, nil
	}

	return Proxy{}, fmt.Errorf("No proxy found by that name for %s", n)
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

	return
}
