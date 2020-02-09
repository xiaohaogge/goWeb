package utils

import (
	"bytes"
	"crypto/md5"
	"database/sql"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/russross/blackfriday"
	"github.com/sourcegraph/syntaxhighlight"
	"html/template"
	"log"
)

var db *sql.DB

func InitMysql() {
	fmt.Println("init mysql........")
	driverName := beego.AppConfig.String("driverName")

	//注册数据库驱动
	orm.RegisterDriver(driverName, orm.DRMySQL)
	//数据库连接
	user := beego.AppConfig.String("mysqluser")
	pwd := beego.AppConfig.String("mysqlpwd")
	host := beego.AppConfig.String("host")
	port := beego.AppConfig.String("port")
	dbname := beego.AppConfig.String("dbname")

	//dbConn := "root:yu271400@tcp(127.0.0.1:3306)/cmsproject?charset=utf8"
	dbConn := user + ":" + pwd + "@tcp(" + host + ":" + port + ")/" + dbname + "?charset=utf8"
	fmt.Println("dbconn:", dbConn)
	// err := orm.RegisterDataBase("default", driverName, dbConn)
	// if err != nil {
	// 	fmt.Println("连接数据库出错......")
	// } else {
	// 	fmt.Println("连接数据库成功......")
	// }
	if db == nil {
		db, err := sql.Open(driverName,dbConn)
		if err != nil {
			fmt.Println("连接数据库出错，err=",err)
		} else {
			// 创建用户表
			CreateTableWithUser()
			CreateTableWithArticle()
			CreateTableWithAlbum()
		}
	}
}

//创建用户表；
func CreateTableWithUser() {
	sql := `CREATE TABLE IF NOT EXISTS users(
		id INT(4) PRIMARY KEY AUTO_INCREMENT NOT NULL,
		username VARCHAR(64),
		password VARCHAR(64),
		status INT(4),
		createtime INT(10)
	);`
	ModifyDB(sql)
}

// 创建文章表
func CreateTableWithArticle(){
	sql := `create table if noe exists article(
			id int(4) primary key auto_increment not null,
			title varchar(30),
			author varchar(20),
			tags varchar(30),
			short varchar(255),
			content longtext,
			createtime int(10)
			);`
	ModifyDB(sql)
}

// 创建图片表
func CreateTableWithAlbum() {
	sql := `create table if not exists album(
		id int(4) primary key auto_increment not null,
		filepath varchar(255),
		filename varchar(64),
		status int(4),
		createtime int(10)
	);`
	ModifyDB(sql)
}

//操作数据库；
func ModifyDB(sql string, args ...interface{}) (int64, error) {
	result, err := db.Exec(sql, args...)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	count, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return count, nil
}

//查询数据库
func QueryRowDB(sql string) *sql.Row {
	row := db.QueryRow(sql)
	fmt.Println("row：", row)
	return row
}

func QueryDB(sql string)(*sql.Rows, error){
	return db.Query(sql)
}

// 添加工具方法
func MD5(str string) string {
	md5str := fmt.Sprintf("%x", md5.Sum([]byte(str)))
	return md5str
}

// 将文章详情的内容 装换成html语句；
func SwitchTimeStampToData(content string) template.HTML{
	markdown := blackfriday.MarkdownCommon([]byte(content))

	// 获取到html 文档；
	doc,_ := goquery.NewDocumentFromReader(bytes.NewReader(markdown))
	/*
	对document 进程查询，选择器和css的语法一样；第一个参数 i是查询到的第几个元素；
	第二个参数：selection 就是查询到的元素
	*/
	doc.Find("code").Each(func(i int, selection *goquery.Selection) {
		light,_ := syntaxhighlight.AsHTML([]byte(selection.Text()))
		selection.SetHtml(string(light))
		fmt.Println(selection.Html())
		fmt.Println("light:",string(light))
		fmt.Println("\n\n\n")
	})
	htmlString, _ := doc.Html()
	return template.HTML(htmlString)
}