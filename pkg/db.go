package pkg

import (
	"fmt"
	"gopkg.in/pg.v3"
	"time"
)

type Format struct {
	Info     string
	Dt       string
	Account  string
	Context  string
	Company  string
	Deal     string
	Serverip string
	Clientip string
	Cliaddr  string
	Note     string
	Dispnode string
}

func ReaddbConfig() {
	fmt.Println(time.Now(), "----reading the config !!!")
	fmt.Printf("host = %v\n", HostName())
	fmt.Printf("port = %v\n", Port())
	fmt.Printf("user = %v\n", Username())
	fmt.Printf("passwd = %v\n", Password())
	fmt.Printf("database = %v\n", Dbname())
	// fmt.Printf("nodes is %v ,node info is %v !\n", NodeCount(), NodeInfo())
	fmt.Printf("%v connect the database !\n", time.Now())
}
func DBConfig() *pg.DB {
	db := pg.Connect(&pg.Options{
		Host:     HostName(),
		Port:     Port(),
		User:     Username(),
		Password: Password(),
		Database: Dbname(),
	})
	return db
}

// func getdailysigntable() ([]dailysign, error) {
// 	var us []dailysign
// 	_, err := db.Query(&us, "SELECT * FROM dailysign")
// 	if err != nil {
// 		return nil, err
// 	}
// 	return us, nil
// }
func InsertFormat(db *pg.DB, a *Format) error {
	_, err := db.ExecOne(`insert into disp (dt,info,account,context,company,deal,serverip,clientip,cliaddr,dispnode,note)  values (?dt,?info,?account,?context,?company,?deal,?serverip,?clientip,?cliaddr,?dispnode,?note) `, a)
	return err
}
