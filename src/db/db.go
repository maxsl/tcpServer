package db

import (
	"fmt"
	"database/sql" //这包一定要引用，是底层的sql驱动
	_ "github.com/go-sql-driver/mysql"
	)

var serverIP = "45.56.92.95"
var serverPort = 8080
var DBPassword = "1986814"

var talkDB *sql.DB

func Get() *sql.DB {
	if(talkDB == nil) {
		args := fmt.Sprintf("root:%s@tcp(%s:3306)/secretchat?charset=utf8mb4", DBPassword, serverIP)
		db, err := sql.Open("mysql", args)
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}
		talkDB = db
		talkDB.SetMaxIdleConns(16)
		err = talkDB.Ping()
		if err != nil {
			fmt.Print("Can not connect database...")
			return nil
		}
	}
	return talkDB;
}
