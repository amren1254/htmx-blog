package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/amren1254/htmx-blog/entity"
	_ "github.com/lib/pq"
)

// Define the PostgreSQL connection string
func createConnString() string {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	password := os.Getenv("DB_PASSWORD")
	user := os.Getenv("DB_USERNAME")
	database := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, database, sslmode)
}

func InitDb() (*sql.DB, error) {

	// Open a database connection
	DB, err := sql.Open("postgres", createConnString())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	// defer DB.Close() // Close the database connection when done

	// Check the connection
	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	log.Println("Connected to PostgreSQL")
	return DB, nil
}

func InsertBlogInDatabase(db *sql.DB, blog entity.Blog) (bool, error) {
	query := `INSERT INTO blogpost(title, content, image_path, author_id, published_date)
				VALUES($1,$2,$3,$4,$5)`
	insert, err := db.Prepare(query)
	if err != nil {
		log.Println("Query Error : ", err)
		return false, err
	}
	defer insert.Close()

	var datetime = time.Now().UTC()
	queryResponse, err := insert.Exec(
		blog.Title,
		blog.Content,
		blog.ImagePath,
		blog.AuthorId,
		datetime.Format(time.RFC3339),
	)
	if err != nil {
		log.Println("Query Response error :", err)
		return false, err
	}
	if affectedRows, err := queryResponse.RowsAffected(); err != nil || affectedRows != 1 {
		log.Println("Blog insertion failed : ", err)
		return false, err
	}
	return true, nil
}

func GetBlogData(db *sql.DB, blogReq entity.BlogRequest) (entity.Blog, error) {
	var blog entity.Blog
	if err := db.QueryRow(`SELECT * FROM blogpost WHERE post_id = $1 AND author_id = $2`, blogReq.PostId, blogReq.AuthorId).
		Scan(
			&blog.Id,
			&blog.Title,
			&blog.Content,
			&blog.ImagePath,
			&blog.AuthorId,
			&blog.PublishedDate); err != nil {
		if err == sql.ErrNoRows {
			log.Println("Error getting rows: no rows found", err)
			return entity.Blog{}, err
		}
		log.Println("Error querying select statement : ", err)
		return entity.Blog{}, err
	}

	return blog, nil
}

func GetBlogDataForAuthor(db *sql.DB, blogReq entity.BlogRequest) ([]entity.Blog, error) {
	var blog []entity.Blog
	rows, err := db.Query(`SELECT * FROM blogpost WHERE author_id = $1`, blogReq.AuthorId)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Error getting rows: no rows found", err)
			return []entity.Blog{}, err
		}
		log.Println("Error querying select statement : ", err)
		return []entity.Blog{}, err
	}
	defer rows.Close()
	// Iterate through the result rows
	for rows.Next() {
		var newBlog entity.Blog
		if err := rows.Scan(&newBlog.Id, &newBlog.Title, &newBlog.Content, &newBlog.ImagePath, &newBlog.AuthorId, &newBlog.PublishedDate); err != nil {
			log.Fatal(err)
		}
		blog = append(blog, newBlog)
	}
	return blog, nil
}

func InsertAuthor(db *sql.DB, author entity.Author) (bool, error) {
	query := `INSERT INTO author(author_name, email)
			VALUES($1,$2)`
	insert, err := db.Prepare(query)
	if err != nil {
		log.Println("Query Error : ", err)
		return false, err
	}
	defer insert.Close()

	queryResponse, err := insert.Exec(
		author.Name,
		author.Email,
	)
	if err != nil {
		log.Println("Query Response error :", err)
		return false, err
	}
	if affectedRows, err := queryResponse.RowsAffected(); err != nil || affectedRows != 1 {
		log.Println("Author insertion failed : ", err)
		return false, err
	}
	return true, nil
}
