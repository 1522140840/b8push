package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"b8push/utils"
	"b8push/conf"
	"strconv"
)

var db *sql.DB

func init() {
	url:=conf.GetVal("mysql","jdbc-url")
	maxOpenConns:=conf.GetVal("mysql","MaxOpenConns")
	maxIdleConns:=conf.GetVal("mysql","MaxIdleConns")


	fmt.Printf("初始化mysql信息，url:[%s],maxOpenConns:[%s],maxIdleConns:[%s]\n",url,maxOpenConns,maxIdleConns)
	db, _ = sql.Open("mysql", url)
	openConns,_:=strconv.Atoi(maxOpenConns)
	idleConns,_:=strconv.Atoi(maxIdleConns)
	db.SetMaxOpenConns(openConns)
	db.SetMaxIdleConns(idleConns)
	db.Ping()
}


func QuerySymbol() (symbols map[string]bool,err error){
	rows, err := db.Query("SELECT exchange,symbol FROM exchange_symbol")
	defer rows.Close()

	if err!=nil{
		fmt.Printf("read mysql symbol error:%s\n",err)
		return
	}

	symbols=make(map[string]bool)
	for rows.Next() {
		var exchange string
		var symbol string
		err = rows.Scan(&exchange, &symbol)
		if err==nil&&!util.StrIsBlank(exchange)&&!util.StrIsBlank(symbol){
			symbols[exchange+"."+symbol]=true
		}
	}

	return

}

func QueryExchange() (exchanges map[string]bool,err error){
	rows, err := db.Query("SELECT name FROM exchange")
	defer rows.Close()

	if err!=nil{
		fmt.Printf("read mysql symbol error:%s\n",err)
		return
	}

	exchanges=make(map[string]bool)
	for rows.Next() {
		var exchange string
		err = rows.Scan(&exchange)
		if err==nil&&!util.StrIsBlank(exchange){
			exchanges[exchange+".overview"]=true
		}
	}

	return

}


