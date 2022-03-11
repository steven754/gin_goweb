package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"net/http"
)

type Todo struct {
	ID int	`gorm:"primaryKey" json:"id"`
	Title string	`json:"title"`
	Status bool	`json:"status"`
}

var (
	DB *gorm.DB
)

const (
	user     string = "root"
	password string = "123456"
	host     string = "stevenwang.top"
	dbName   string = "Todo"
)

func ConnectMysql() (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		user,
		password,
		host,
		dbName)

	//dsn := "root:123456@tcp(stevenwang.top:3306)/Todo?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
		panic("连接数据库失败")
	}
	DB.LogMode(true)  //开启调试模式-用于显示sql语句
	DB.SingularTable(true)  //禁止表名复数形式
	return DB.DB().Ping()
}

func main(){
	//1.创建数据库CREATE DATABASE Todo;
	//2.连接数据库
	err := ConnectMysql()
	if err != nil {
		panic(err)
	}
	defer DB.Close()
	DB.AutoMigrate(&Todo{})

	//获取gin默认引擎
	r := gin.Default()
	//3.gin框架静态文件
	r.Static("/static","static")
	//4.gin框架模板文件
	r.LoadHTMLGlob("templates/*")

	//获取首页静态地址资源
	r.GET("/", func (c *gin.Context){
		c.HTML(http.StatusOK, "index.html", nil)
	})
	//5.添加路由组
	v1Group := r.Group("/v1")
	{
		//获取待办事件列表
		v1Group.GET("/todo", func(c *gin.Context) {
			var todoList [] Todo
			if err = DB.Find(&todoList).Error; err!= nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}else {
				c.JSON(http.StatusOK, todoList)
			}
		})

		//添加待办事件列表
		v1Group.POST("/todo", func(c *gin.Context) {
			var todo Todo
			c.BindJSON(&todo)
			err = DB.Create(&todo).Error
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, todo)}
		})

		//修改待办事件列表
		v1Group.PUT("/todo/:id", func(c *gin.Context) {
			id, ok := c.Params.Get("id")
			if !ok {
				c.JSON(http.StatusBadRequest, gin.H{"error": "无效的id"})
				return
			}
			var todo Todo
			if err = DB.Where("id=?", id).First(&todo).Error; err!=nil{
				c.JSON(http.StatusOK, gin.H{"error": err.Error()})
				return
			}
			c.BindJSON(&todo)
			if err = DB.Save(&todo).Error; err!= nil{
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}else{
				c.JSON(http.StatusOK, todo)
			}
		})

		//删除待办事件列表
		v1Group.DELETE("/todo/:id", func(c *gin.Context) {
			id, ok := c.Params.Get("id")
			if !ok {
				c.JSON(http.StatusBadRequest, gin.H{"error": "无效的id"})
				return
			}
			if err = DB.Where("id=?", id).Delete(Todo{}).Error;err!=nil{
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}else{
				c.JSON(http.StatusOK, gin.H{id:"已被删除"})
			}
		})
	}

	r.Run(":9090")
	}
