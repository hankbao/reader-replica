// handler.go
// author: hankbao

package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/hankbao/reader-replica/scrape"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db}
}

func (h *Handler) GetFeedById(c *gin.Context) {
	feedId := c.Param("id")
	if feedId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Feed ID is required"})
		return
	}

	var feed scrape.Feed
	result := h.db.First(&feed, feedId)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Feed not found"})
		} else {
			log.Printf("Query database error: %v", result.Error)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, feed)
}

func (h *Handler) GetArticleById(c *gin.Context) {
	articleId := c.Param("id")
	if articleId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Article ID is required"})
		return
	}

	var article scrape.Article
	result := h.db.First(&article, articleId)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		} else {
			log.Printf("Query database error: %v", result.Error)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, article)
}

func (h *Handler) Subscribe(c *gin.Context) {
	link := c.PostForm("link")
	if link == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Link is required"})
		return
	}

	scraper := scrape.NewScraper(30)
	feed, err := scraper.ScrapeFeed(link)
	if err != nil {
		log.Printf("Scrape feed error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if result := h.db.Create(&feed); result.Error != nil {
		log.Printf("Insert feed into database error: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, feed)
}
