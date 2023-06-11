package dao

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func OpenDb() (*sql.DB, error) {
	Db, err := sql.Open("mysql", "root:admin@tcp(127.0.0.1:3306)/dataframe")
	if err != nil {
		log.Println(1)
		return nil, err
	}
	err = Db.Ping()
	if err != nil {
		log.Println(2)
		return nil, err
	}
	return Db, err
}

func AllUserInfo(name string) (*sql.Rows, error) {
	Db, err := OpenDb()
	stmt, err := Db.Prepare("select * from user where phoneoremail= ?")
	if err != nil {
		log.Print("1:")
		log.Println(err)
		return nil, err
	}
	row, err := stmt.Query(name)
	if err != nil {
		log.Println("2:")
		log.Println(err)
		return nil, err
	}
	return row, nil
}
