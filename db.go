package main

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// Some errors
var (
	ErrDBNotAlive   = errors.New("DB not alive")
	ErrDoesNotExist = errors.New("Database does not exist")
)

type Database struct {

	// Database object
	*sql.DB
}

func InitDB(p string) (db *Database, err error) {

	db = &Database{}

	db.DB, err = sql.Open("sqlite3", DatabaseFile)
	if err != nil {

		return
	}

	if err := db.DB.Ping(); err != nil {

		return nil, ErrDBNotAlive
	}

	ok, err := db.TableExists("MessageCache")
	if err != nil {

		return
	}
	if !ok {

		if _, err := db.DB.Exec("CREATE TABLE `MessageCache` (`user_id` TEXT, `message_id` TEXT)"); err != nil {

			return nil, err
		}
	}

	return
}

func (db *Database) TableExists(n string) (ok bool, err error) {

	if err := db.DB.Ping(); err != nil {

		return false, ErrDBNotAlive
	}

	rows, err := db.DB.Query(fmt.Sprintf("SELECT `name` FROM `sqlite_master` WHERE name='%s'", n))
	if err != nil {

		return
	}

	return rows.Next(), nil
}

func (db *Database) AddMessageToCache(userid, msgid string) error {

	if err := db.DB.Ping(); err != nil {

		return ErrDBNotAlive
	}

	if _, err := db.DB.Exec(fmt.Sprintf("INSERT INTO `MessageCache` (`user_id`, `message_id`) VALUES ('%s', '%s')", userid, msgid)); err != nil {

		return err
	}

	return nil
}

func (db *Database) RemoveMessageFromCache(userid, msgid string) (err error) {

	if err := db.DB.Ping(); err != nil {

		return ErrDBNotAlive
	}

	if _, err := db.DB.Exec(fmt.Sprintf("DELETE FROM `MessageCache` WHERE user_id='%s' AND message_id='%s'", userid, msgid)); err != nil {

		return err
	}

	return
}

func (db *Database) GetMessageCache(userid string) (c []string, err error) {

	if err := db.DB.Ping(); err != nil {

		return nil, ErrDBNotAlive
	}

	rows, err := db.DB.Query(fmt.Sprintf("SELECT `message_id` FROM `MessageCache` WHERE user_id='%s'", userid))
	if err != nil {

		return
	}
	defer rows.Close()

	for rows.Next() {

		var msgid string

		if err := rows.Scan(&msgid); err != nil {

			return nil, err
		}

		c = append(c, msgid)
	}
	if err := rows.Err(); err != nil {

		return nil, err
	}

	return
}

func (db *Database) IsInCache(userid, msgid string) (ok bool, err error) {

	if err := db.DB.Ping(); err != nil {

		return false, ErrDBNotAlive
	}

	c, err := db.GetMessageCache(userid)
	if err != nil {

		return
	}
	for _, v := range c {

		if v == msgid {

			return true, nil
		}
	}

	return
}
