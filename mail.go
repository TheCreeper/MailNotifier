package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/TheCreeper/go-notify"
	"github.com/TheCreeper/go-pop3"
)

func (cfg *ClientConfig) LaunchPOP3Client(wg *sync.WaitGroup, acc *Account) {

	for {

		c, err := pop3.DialTLS(acc.Host)
		if err != nil {

			log.Fatal(err)
		}
		defer c.Quit()

		if err = c.Auth(acc.Username, acc.Password); err != nil {

			log.Fatal(err)
		}

		if Debug {

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

			ok, err := db.IsInCache(HashString(acc.Username+acc.Host), v.UID)
			if err != nil {

				log.Fatal(err)
			}
			if ok {

				continue
			}

			if err := db.AddMessageToCache(HashString(acc.Username+acc.Host), v.UID); err != nil {

				log.Fatal(err)
			}

			m, err := c.Top(v.ID, 0)
			if err != nil {

				log.Fatal(err)
			}

			n := &notify.Message{

				Title:     fmt.Sprintf("From: %s ", m.Header.Get("From")),
				Body:      fmt.Sprintf("Subject: %s ", m.Header.Get("Subject")),
				Icon:      "/usr/share/icons/gnome/32x32/status/mail-unread.png",
				SoundPipe: LetterArriveSound,
			}
			if err = n.Send(); err != nil {

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
