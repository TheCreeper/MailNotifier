package main

import (
	"fmt"
	"log"
	"time"

	"github.com/TheCreeper/go-notify"
	"github.com/TheCreeper/go-pop3"
)

func (cfg *ClientConfig) LaunchPOP3Client(acc *Account) {

	for {

		c, err := pop3.DialTLS(acc.Address)
		if err != nil {

			log.Print(err)
			time.Sleep(time.Duration(cfg.CheckFrequency) * time.Minute)
			continue
		}

		if err = c.Auth(acc.User, acc.Password); err != nil {

			log.Fatal(err)
		}

		if Verbose {

			count, size, err := c.Stat()
			if err != nil {

				log.Fatal(err)
			}
			log.Printf("DEBUG: Message Count: %d, Mailbox Space Consumed: %d bytes", count, size)
		}

		// Check for new messages
		messages, err := c.UidlAll()
		if err != nil {

			log.Fatal(err)
		}
		for _, v := range messages {

			ok, err := db.IsInCache(HashString(acc.User+acc.Address), v.UID)
			if err != nil {

				log.Fatal(err)
			}
			if ok {

				continue
			}

			if err := db.AddMessageToCache(HashString(acc.User+acc.Address), v.UID); err != nil {

				log.Fatal(err)
			}

			m, err := c.Top(v.ID, 0)
			if err != nil {

				log.Fatal(err)
			}

			ntf := &notify.Notification{

				Summary: fmt.Sprintf("From: %s", m.Header.Get("From")),
				Body:    fmt.Sprintf("Subject: %s", m.Header.Get("Subject")),
				AppIcon: cfg.NotificationIcon,
				Hints:   map[string]string{"sound-name": cfg.NotificationSound},
				Timeout: notify.ExpiresDefault,
			}
			if _, err = ntf.Show(); err != nil {

				log.Fatal(err)
			}
		}

		if err = c.Quit(); err != nil {

			log.Fatal(err)
		}

		time.Sleep(time.Duration(cfg.CheckFrequency) * time.Minute)
	}

	wg.Done()
}
