package entity

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	CurrentDirectory string
	StaticFileDir    string `env:"STATIC_FILE_DIR,default=./static"`
	Server           struct {
		Host string `env:"HOST default:localhost"`
		Port int    `env:"PORT default:8081"`
	}
}

type EntityID string

type EntityStatus int

type Author struct {
	ID            EntityID     `json:"ID,omitempty"`
	CreatedAt     time.Time    `json:"createdAt,omitempty"`
	UpdatedAt     *time.Time   `json:"updatedAt,omitempty"`
	LastLoginAt   *time.Time   `json:"lastLoginAt,omitempty"`
	Status        EntityStatus `json:"status,omitempty"`
	LoginAttempts int          `json:"loginAttempts,omitempty"`
	Name          string       `json:"name"`
	Email         string       `json:"email"`
	Password      string       `json:"password"`
}

func LoadConfig() (*Config, error) {
	var cfg Config
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	cfg.CurrentDirectory = cwd
	cfg.Server.Host = os.Getenv("HOST")
	cfg.Server.Port, err = strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *Config) ServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

type Blog struct {
	Id            int    `json:"id,omitempty"`
	Title         string `json:"title"`
	Content       string `json:"content"`
	ImagePath     string `json:"image_path,omitempty"`
	AuthorId      int    `json:"author_id"`
	PublishedDate string `json:"published_date,omitempty"`
}

type BlogRequest struct {
	AuthorId int `json:"author_id"`
	PostId   int `json:"post_id"`
}
