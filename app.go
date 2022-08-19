package main

import (
	"flag"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/michael34435/lnurlp/models"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

func createSchema(db *pg.DB) error {
	models := []interface{}{
		(*models.Alias)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	db := pg.Connect(&pg.Options{
		User:     os.Getenv("PG_USER"),
		Password: os.Getenv("PG_PASS"),
		Database: os.Getenv("PG_DB"),
		Addr:     os.Getenv("PG_HOST"),
	})
	defer db.Close()

	err := createSchema(db)
	if err != nil {
		panic(err)
	}

	prod := flag.Bool("prod", false, "for production")
	flag.Parse()

	if *prod {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.Use(gin.ErrorLogger())
	r.GET("/.well-known/lnurlp/:key", func(c *gin.Context) {
		key := c.Param("key")

		alias := models.Alias{}
		err := db.Model(&alias).Where("key = ?", key).Select()

		if err != nil {
			c.AbortWithError(400, err)
			return
		}

		c.Redirect(301, "https://"+alias.Domain+"/.well-known/lnurlp/"+alias.User)
	})
	r.Run()
}
