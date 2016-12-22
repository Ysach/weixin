package mysqldb

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
	"strings"
)

func checkErr(err error) {
	if err != nil {
		fmt.Println("数据库更新数据失败: ", err)
		return
	}
}

func MysqlData(content string) error{

	db, err := sql.Open("mysql", "wx:wx@/data?charset=utf8")
	checkErr(err)

	conn, err := db.Prepare("INSERT msg SET msg=?")
	if err != nil {
		fmt.Println("数据库更新失败：", err)
		db.Close()
	}

	res, err := conn.Exec(content)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)
	fmt.Println(affect)
	db.Close()
	return nil

}

func MysqlSel() string{
	phone := []string{}
	db, err := sql.Open("mysql", "wx:wx@/data?charset=utf8")
	checkErr(err)
	conn, err := db.Query("select * from phone")
	if err != nil {
		fmt.Println("数据库更新失败：", err)
		db.Close()
	}
	for conn.Next() {
		var id int
		var phone_num string
		err = conn.Scan(&id, &phone_num)
		checkErr(err)
		//fmt.Println(id)
		phone = append(phone, phone_num)
		//fmt.Println(msg)
		//fmt.Println(create)
		//fmt.Println(phone)
	}
	s := strings.Join(phone, ",")
	db.Close()
	return s
}