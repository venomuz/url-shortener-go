package router

import (
	"context"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/venomuz/url-shortener-go/config"
	"github.com/venomuz/url-shortener-go/pkg/logger"
	"github.com/venomuz/url-shortener-go/router/cors"
	"github.com/venomuz/url-shortener-go/storage"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type CreateShortUrl struct {
	Link string `json:"link" binding:"required" example:"https://github.com/"`
}
type ResOk struct {
	Url     string `json:"url"`
	Message string `json:"message"`
}
type ResError struct {
	Message string `json:"message" example:"error"`
}
type Option struct {
	Conf config.Config
	Log  logger.Logger
	Rds  storage.RepositoryStorage
}
type OptionConfig struct {
	Conf      config.Config
	Logger    logger.Logger
	RedisRepo storage.RepositoryStorage
}

func Opt(c *OptionConfig) *Option {
	return &Option{
		Log:  c.Logger,
		Conf: c.Conf,
		Rds:  c.RedisRepo,
	}
}
func New(option Option) *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(cors.New(cors.GinCorsMiddleware()))
	handler := Opt(&OptionConfig{
		Logger:    option.Log,
		Conf:      option.Conf,
		RedisRepo: option.Rds,
	})

	router.POST("/url", handler.CreateUrl)
	router.GET("/:url", handler.GetUrl)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return router
}

// CreateUrl This api short your url
// @Summary     Create short url
// @Description This api short your url
// @Tags        My-API
// @Accept      json
// @Produce     json
// @Param       data body     CreateShortUrl true "data body"
// @Success     200  {object} ResOk
// @Failure     400 {object} ResError
// @Failure     500  {object} ResError
// @Router      /url [POST]
func (h *Option) CreateUrl(c *gin.Context) {
	var body CreateShortUrl
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.Log.Error("failed to bind json", logger.Error(err))
		return
	}
	if !validLink(body.Link) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "your url is not valid ex: https://github.com",
		})
		h.Log.Error("Wrong url name", logger.Error(err))
		return
	}

	newLink := randStringBytes(9)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.Conf.CtxTimeout))
	defer cancel()

	err = h.Rds.Set(newLink, body.Link, ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "error while save url",
		})
		h.Log.Error("error while save url", logger.Error(err))
		return
	}
	c.JSON(http.StatusOK, ResOk{
		Url:     "http://52.42.75.134" + h.Conf.HTTPPort + "/" + newLink,
		Message: "success",
	})
}

func (h *Option) GetUrl(c *gin.Context) {
	name := c.Param("url")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.Conf.CtxTimeout))
	defer cancel()
	get, err := h.Rds.Get(name, ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Not fund your url",
		})
		h.Log.Error("Error while get url", logger.Error(err))
		return
	}
	c.Redirect(http.StatusMovedPermanently, get)
}
func validLink(link string) bool {
	r, err := regexp.Compile("^(http|https)://")
	if err != nil {
		return false
	}
	link = strings.TrimSpace(link)
	log.Printf("Checking for valid link: %s", link)
	// Check if string matches the regex
	if r.MatchString(link) {
		return true
	}
	return false
}
func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
