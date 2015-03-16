package main

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const (
	// name of the table that is used for the message cache
	MsgTableName = "MessageCache"

	// statement used to create the cache table
	MsgTableCreateStmt = "CREATE TABLE `MessageCache` (`user_id` TEXT, `message_id` TEXT)"

	// statement used to insert new messages into the cache
	MsgTableInsertStmt = "INSERT INTO `MessageCache` (`user_id`, `message_id`) VALUES ('%s', '%s')"

	// statement used to remove messages from the cache
	MsgTableRemoveStmt = "DELETE FROM `MessageCache` WHERE user_id='%s' AND message_id='%s'"

	// statement used to get all records in the message cache
	MsgTableGetStmt = "SELECT `message_id` FROM `MessageCache` WHERE user_id='%s'"

	// statement used to check if table exists
	TableExistStmt = "SELECT `name` FROM `sqlite_master` WHERE name='%s'"
)

// some errors
var (
	ErrDBNotAlive   = errors.New("DB not alive")
	ErrDoesNotExist = errors.New("Database does not exist")
)

type Database struct{ *sql.DB }

func InitDB(p string) (db *Database, err error) {

	db = &Database{}

	db.DB, err = sql.Open("sqlite3", DatabaseFile)
	if err != nil {

		return
	}

	if err := db.DB.Ping(); err != nil {

		return nil, ErrDBNotAlive
	}

	ok, err := db.TableExists(MsgTableName)
	if err != nil {

		return
	}
	if !ok {

		if _, err := db.DB.Exec(MsgTableCreateStmt); err != nil {

			return nil, err
		}
	}

	return
}

func (db *Database) TableExists(n string) (ok bool, err error) {

	if err := db.DB.Ping(); err != nil {

		return false, ErrDBNotAlive
	}

	rows, err := db.DB.Query(fmt.Sprintf(TableExistStmt, n))
	if err != nil {

		return
	}

	return rows.Next(), nil
}

func (db *Database) AddMessageToCache(userid, msgid string) error {

	if err := db.DB.Ping(); err != nil {

		return ErrDBNotAlive
	}

	if _, err := db.DB.Exec(fmt.Sprintf(MsgTableInsertStmt, userid, msgid)); err != nil {

		return err
	}

	return nil
}

func (db *Database) RemoveMessageFromCache(userid, msgid string) (err error) {

	if err := db.DB.Ping(); err != nil {

		return ErrDBNotAlive
	}

	if _, err := db.DB.Exec(fmt.Sprintf(MsgTableRemoveStmt, userid, msgid)); err != nil {

		return err
	}

	return
}

func (db *Database) GetMessageCache(userid string) (c []string, err error) {

	if err := db.DB.Ping(); err != nil {

		return nil, ErrDBNotAlive
	}

	rows, err := db.DB.Query(fmt.Sprintf(MsgTableGetStmt, userid))
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
