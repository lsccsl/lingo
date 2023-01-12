package main

import (
	"database/sql"

	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func test_mysql() {
	db, err := sql.Open("mysql", "root:123456@tcp(192.168.15.146:3306)/test?charset=utf8")
	if err != nil {
		fmt.Printf("connect mysql fail ! [%s]", err)
	} else {
		fmt.Println("connect to mysql success")
	}

	rows, err := db.Query("select aa from test")
	if err != nil {
		fmt.Printf("select fail [%s]", err)
	}

	fmt.Println("err:", err)
	//fmt.Println("rows:", rows)

	for rows.Next() {
		//fmt.Println("rows:", rows)
		var aa string
		rows.Columns()
		err := rows.Scan(&aa)
		if err != nil {
			fmt.Printf("get user info error [%s]", err)
		}
		fmt.Println(" aa:", aa)
	}
}

func main() {
	test_mysql()
}