package router

import (
	"fmt"
	"log"
	"net/http"

	"github.com/amren1254/htmx-blog/controller"
	"github.com/amren1254/htmx-blog/database"
	"github.com/gin-gonic/gin"
)

func InitRoute() {
	//create db connection and ping
	db := database.InitDb()
	blogHandler := controller.NewDbEnv(db)

	router := gin.Default()
	router.Static("/static", "./static")
	// router.LoadHTMLGlob("htmx-template/*.html")

	router.GET("/", func(c *gin.Context) {
		c.File("htmx-template/index.html")
	})
	type formStruct struct {
		form string
	}
	router.GET("/newblog", func(c *gin.Context) {
		var form formStruct
		form.form = fmt.Sprintf("<form id='inputForm' hx-post='/blog' hx-trigger='submit'><div>					<label for='id'>ID:</label>					<input type='text' name='id' id='id' required>					<label for='title'>Title:</label>					<input type='text' name='title' id='title' required>					<label for='content'>Content:</label>					<textarea name='content' id='content' rows='4' cols='50' required></textarea>					<label for='image_path'>Image Path:</label>					<input type='text' name='image_path' id='image_path' required>					<label for='author_id'>Author ID:</label>					<input type='text' name='author_id' id='author_id' required>					<button type='submit'>Submit</button>				</div>			</form>	")
		c.JSON(200, gin.H{"form": form.form})
	})
	router.GET("/blog", blogHandler.GetBlog)
	router.POST("/blog", blogHandler.InsertBlog)
	port := "8081" //os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	log.Fatal(http.ListenAndServe(":"+port, router))
}
