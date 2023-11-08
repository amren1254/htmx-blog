package controller

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/amren1254/htmx-blog/database"
	"github.com/amren1254/htmx-blog/entity"
	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday/v2"
)

type DbEnv struct {
	db *sql.DB
}

func NewDbEnv(db *sql.DB) *DbEnv {
	return &DbEnv{db: db}
}

// get the requested blog by author
func (env *DbEnv) GetBlog(c *gin.Context) {
	// check for db connection by PING() method
	if err := env.db.Ping(); err != nil {
		log.Println("Db connection error ", err)
		c.JSON(500, gin.H{"message": err})
		return
	}
	var blogRequest entity.BlogRequest
	var err error
	blogRequest.AuthorId, err = strconv.Atoi(c.Query("author_id"))
	if err != nil {
		log.Println("Error getting author id from query string paramter")
		c.JSON(500, gin.H{"message": err})
		return
	}
	blogRequest.PostId, err = strconv.Atoi(c.Query("post_id"))
	if err != nil {
		log.Println("Error getting post id from query string paramter")
		c.JSON(500, gin.H{"message": err})
		return
	}
	if blogRequest.PostId == 0 && blogRequest.AuthorId != 0 {
		blogData, err := database.GetBlogDataForAuthor(env.db, blogRequest)
		if err != nil {
			c.JSON(500, gin.H{"message": err})
			return
		}
		c.JSON(200, gin.H{"data": blogData})
		return
	} else {
		blogData, err := database.GetBlogData(env.db, blogRequest)
		if err != nil {
			c.JSON(500, gin.H{"message": err})
			return
		}
		c.JSON(200, gin.H{"data": blogData})
		return
	}
}

// insert blogs in database
func (env *DbEnv) InsertBlog(c *gin.Context) {
	//check for db connection by PING() method
	if err := env.db.Ping(); err != nil {
		log.Println("Db connection error ", err)
		c.JSON(500, gin.H{"message": err})
		return
	}
	err := c.Request.ParseForm()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	//bind input json data to blog entity
	var blog entity.Blog
	blog.Id, err = strconv.Atoi(c.PostForm("id"))
	blog.Title = c.PostForm("title")
	blog.Content = c.PostForm("content")
	htmlContent := blackfriday.Run([]byte(blog.Content))
	log.Println(htmlContent)
	blog.AuthorId, err = strconv.Atoi(c.PostForm("author_id"))

	if err != nil {
		log.Println("Input form parsing error", err)
		c.JSON(500, gin.H{"message": err})
		return
	}

	//call db to insert recvd blog data
	status, err := database.InsertBlogInDatabase(env.db, blog)
	if err != nil || status == false {
		c.JSON(500, gin.H{"message": err})
		return
	}

	c.JSON(201, gin.H{"Message": "Blog created successfully"})
}
