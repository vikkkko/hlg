package dao

import (
	"fmt"
	"github.com/tpkeeper/pitaya/models"
	"gorm.io/gorm"
	"os"
	"testing"
)

var DB *gorm.DB
var dao *Dao

func init() {
	db, err := newDB()
	if err != nil {
		fmt.Println("db init err", err)
		os.Exit(0)
	}
	DB = db
	dao = New(DB)
}
func newDB() (db *gorm.DB, err error) {
	db, err = models.NewDB(&models.Config{
		"127.0.0.1", "3306", "root", "123456", "hbt_dev", 1})

	//db, err = models.NewDB(&models.Config{"161.117.178.50", "3306", "root", "H3s63jhzPgOcSmwS",
	//	"v2mainnet", "release"})
	return
}

func TestDao_GetTotalAllocPoint(t *testing.T) {
	amount, err := dao.GetTotalAllocPoint()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(amount)
}
func TestDao_GetLogsTopBlockNumber(t *testing.T) {
	number,err:=dao.GetLogsTopBlockNumber()
	if err!=nil{
		t.Fatal(err)
	}
	t.Log(number)
}
func TestDao_GetLogsTotalNumberByUser(t *testing.T) {
	number,err:=dao.GetLogsTotalNumberByUser("0x6555d470d605531fcee94f48ba097355e940b869")
	if err!=nil{
		t.Fatal(err)
	}
	t.Log(number)
}