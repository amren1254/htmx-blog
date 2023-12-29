package controller

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/amren1254/htmx-blog/database"
	"github.com/amren1254/htmx-blog/entity"
	"github.com/gin-gonic/gin"
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
	var err error
	err = testDBConnection(env.db)
	if err != nil {
		c.JSON(500, gin.H{"message": err})
	}
	var blogRequest entity.BlogRequest

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
		// Parse the date string into a time.Time object
		parsedTime, err := time.Parse(time.RFC3339, blogData[0].PublishedDate)
		if err != nil {
			fmt.Println("Error parsing date:", err)
			return
		}

		// Extract date and month components
		_, _, day := parsedTime.Date()

		monthString := parsedTime.Month().String()
		blogData[0].PublishedDate = monthString + fmt.Sprintf("%d", day)
		c.JSON(200, gin.H{"data": blogData})
		return
	} else {
		blogData, err := database.GetBlogData(env.db, blogRequest)
		if err != nil {
			c.JSON(500, gin.H{"message": err})
			return
		}
		parsedTime, _ := time.Parse(time.RFC3339, blogData.PublishedDate)
		_, _, day := parsedTime.Date()

		monthString := parsedTime.Month().String()
		blogData.PublishedDate = monthString[:3] + fmt.Sprintf(" %d", day)
		c.JSON(200, gin.H{"data": blogData})
		return
	}
}

// insert blogs in database
func (env *DbEnv) InsertBlog(c *gin.Context) {
	//check for db connection by PING() method
	err := testDBConnection(env.db)
	if err != nil {
		c.JSON(500, gin.H{"message": err})
	}
	err = c.Request.ParseForm()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	//bind input json data to blog entity
	var blog entity.Blog
	blog.Title = c.PostForm("title")
	blog.Content = c.PostForm("content")
	blog.AuthorId, err = strconv.Atoi(c.PostForm("author_id"))
	if blog.AuthorId == 0 {
		blog.AuthorId = 1
	}
	// if err != nil {
	// 	log.Println("Input form parsing error", err)
	// 	c.JSON(500, gin.H{"message": err})
	// 	return
	// }

	//call db to insert recvd blog data
	status, err := database.InsertBlogInDatabase(env.db, blog)
	if err != nil || status == false {
		c.JSON(500, gin.H{"message": err})
		return
	}

	c.JSON(201, gin.H{"Message": "Blog created successfully"})
}

func testDBConnection(db *sql.DB) error {
	if err := db.Ping(); err != nil {
		log.Println("Db connection error ", err)
		return err
	}
	return nil
}

func (env *DbEnv) InsertAuthor(c *gin.Context) {
	//check for db connection by PING() method
	err := testDBConnection(env.db)
	if err != nil {
		c.JSON(500, gin.H{"message": err})
		return
	}
	err = c.Request.ParseForm()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	var author entity.Author
	author.Name = c.PostForm("name")
	author.Email = c.PostForm("email")
	author.Password = c.PostForm("password")

	status, err := database.InsertAuthor(env.db, author)
	if err != nil || status == false {
		c.JSON(500, gin.H{"message": err})
		return
	}
	c.JSON(201, gin.H{"Message": "Author added!"})
}
