package ssjitsi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/chromedp/chromedp"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"sheff.online/jitsibot/docs"
)

func newError(c *gin.Context, s int, err error) {
	c.JSON(s, gin.H{
		"error":   true,
		"message": err.Error(),
	})
}

type HttpServer struct {
	bots   map[string]context.Context
	router *gin.Engine
}

func (h *HttpServer) AddBot(b Bot) {
	h.bots[b.ID] = b.Ctx
}

func (h *HttpServer) Start(l string) error {
	return h.router.Run(l)
}

// @BasePath /api/v1

// ListBots godoc
// @Summary      List bots
// @Description  get bots
// @Tags         main
// @Accept       json
// @Produce      json
// @Success      200  {array}   string
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /bots [get]
func (h *HttpServer) ListBots(c *gin.Context) {
	keys := make([]string, 0, len(h.bots))
	for k := range h.bots {
		keys = append(keys, k)
	}
	c.JSON(http.StatusOK, keys)
}

// HTML endpoint
// @Summary html
// @Schemes
// @Description do main
// @Tags bot
// @Accept json
// @Param   id   path  string  true  "Bot ID"
// @Produce json
// @Success 200 {string} Content
// @Router /{id}/html [get]
func (h *HttpServer) HTML(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		newError(c, http.StatusBadRequest, errors.New("id required"))
		return
	}
	b, ok := h.bots[id]
	if !ok {
		newError(c, http.StatusBadRequest, errors.New("not found"))
		return
	}

	var res string
	err := chromedp.Run(b,
		// chromedp.FullScreenshot(&buf, 100),
		chromedp.OuterHTML("body", &res, chromedp.ByQuery),
	)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	fmt.Println(strings.TrimSpace(res))
	c.JSON(http.StatusOK, res)
}

// screenshot
// @Summary screenshot
// @Schemes
// @Description do screenshot
// @Tags bot
// @Accept json
// @Param   id   path  string  true  "Bot ID"
// @Produce png
// @Success 200 {file} Screenshot
// @Router /{id}/screenshot [get]
func (h *HttpServer) Screenshot(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		newError(c, http.StatusBadRequest, errors.New("id required"))
		return
	}
	b, ok := h.bots[id]
	if !ok {
		newError(c, http.StatusBadRequest, errors.New("not found"))
		return
	}

	var buf []byte
	err := chromedp.Run(b,
		chromedp.FullScreenshot(&buf, 100),
	)
	if err != nil {
		fmt.Println(err)
	}

	c.Data(http.StatusOK, "image/png", buf)
}

func NewHttpServer() *HttpServer {
	srv := HttpServer{bots: map[string]context.Context{}, router: gin.Default()}

	// Настройка CORS middleware
	srv.router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60, // 12 hours
	}))

	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := srv.router.Group("/api/v1")
	{
		v1.GET("/bots", srv.ListBots)
		v1.GET("/:id/html", srv.HTML)
		v1.GET("/:id/screenshot", srv.Screenshot)
	}
	srv.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return &srv
}
