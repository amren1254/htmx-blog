package router

import (
	"log"
	"net/http"

	"github.com/amren1254/htmx-blog/controller"
	"github.com/amren1254/htmx-blog/database"
	"github.com/amren1254/htmx-blog/entity"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func InitRoute() {

	//create db connection and ping
	db, err := database.InitDb()
	defer db.Close()
	if err != nil {
		log.Println("Error connecting database", err)
	}

	config, err := entity.LoadConfig()
	if err != nil {
		log.Println("Error loading config")
	}
	blogHandler := controller.NewDbEnv(db)

	router := gin.Default()
	router.Static("/static", "./static")
	// router.LoadHTMLGlob("htmx-template/*.html")

	router.GET("/", logger.SetLogger(
		logger.WithLogger(func(_ *gin.Context, l zerolog.Logger) zerolog.Logger {
			return l.Output(gin.DefaultWriter).With().Logger()
		}),
	), func(c *gin.Context) {
		c.File("htmx-template/index.html")
	})
	type formStruct struct {
		form string
	}
	router.GET("/newblog", logger.SetLogger(
		logger.WithLogger(func(_ *gin.Context, l zerolog.Logger) zerolog.Logger {
			return l.Output(gin.DefaultWriter).With().Logger()
		}),
	), func(c *gin.Context) {
		c.File("htmx-template/blog.html")
	})
	router.GET("/blog", logger.SetLogger(
		logger.WithLogger(func(_ *gin.Context, l zerolog.Logger) zerolog.Logger {
			return l.Output(gin.DefaultWriter).With().Logger()
		}),
	), blogHandler.GetBlog)

	router.POST("/blog", logger.SetLogger(
		logger.WithLogger(func(_ *gin.Context, l zerolog.Logger) zerolog.Logger {
			return l.Output(gin.DefaultWriter).With().Logger()
		}),
	), blogHandler.InsertBlog)

	router.POST("/author", logger.SetLogger(
		logger.WithLogger(func(_ *gin.Context, l zerolog.Logger) zerolog.Logger {
			return l.Output(gin.DefaultWriter).With().Logger()
		}),
	), blogHandler.InsertAuthor)

	router.GET("/author", logger.SetLogger(
		logger.WithLogger(func(_ *gin.Context, l zerolog.Logger) zerolog.Logger {
			return l.Output(gin.DefaultWriter).With().Logger()
		}),
	), func(c *gin.Context) {
		c.File("htmx-template/author.html")
	})

	log.Println("Started listening at ", config.ServerAddress())
	log.Fatal(http.ListenAndServe(config.ServerAddress(), router))
}
