package main

import (
	"fmt"
	"io"

	"github.com/gin-gonic/gin"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"

	"github.com/gothinkster/golang-gin-realworld-example-app/articles"
	"github.com/gothinkster/golang-gin-realworld-example-app/common"
	"github.com/gothinkster/golang-gin-realworld-example-app/users"
	"github.com/jinzhu/gorm"
)

func Migrate(db *gorm.DB) {
	users.AutoMigrate()
	db.AutoMigrate(&articles.ArticleModel{})
	db.AutoMigrate(&articles.TagModel{})
	db.AutoMigrate(&articles.FavoriteModel{})
	db.AutoMigrate(&articles.ArticleUserModel{})
	db.AutoMigrate(&articles.CommentModel{})
}

func NewReader(rd io.Reader) io.Reader {
	// Make an tranformer that decodes MS-Windows (16LE) UTF files:
	winutf := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	// Make a transformer that is like winutf, but abides by BOM if found:
	decoder := winutf.NewDecoder()
	// Make a Reader that uses decoder:
	return transform.NewReader(rd, unicode.BOMOverride(decoder))
}

func main() {

	db := common.Init()
	Migrate(db)
	defer db.Close()

	r := gin.Default()

	v1 := r.Group("/api")
	users.UsersRegister(v1.Group("/users"))
	v1.Use(users.AuthMiddleware(false))
	articles.ArticlesAnonymousRegister(v1.Group("/articles"))
	articles.TagsAnonymousRegister(v1.Group("/tags"))

	v1.Use(users.AuthMiddleware(true))
	users.UserRegister(v1.Group("/user"))
	users.ProfileRegister(v1.Group("/profiles"))

	articles.ArticlesRegister(v1.Group("/articles"))

	testAuth := r.Group("/api/ping")

	testAuth.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// test 1 to 1
	tx1 := db.Begin()
	userA := users.UserModel{
		Username: "AAAAAAAAAAAAAAAA",
		Email:    "aaaa@g.cn",
		Bio:      "hehddeda",
		Image:    nil,
	}
	tx1.Save(&userA)
	tx1.Commit()
	fmt.Println(userA)

	//db.Save(&ArticleUserModel{
	//    UserModelID:userA.ID,
	//})
	//var userAA ArticleUserModel
	//db.Where(&ArticleUserModel{
	//    UserModelID:userA.ID,
	//}).First(&userAA)
	//fmt.Println(userAA)

	r.Run() // listen and serve on 0.0.0.0:8080
}
