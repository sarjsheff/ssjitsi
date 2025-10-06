package ssjitsi

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"sheff.online/ssjitsi/docs"
)

func newError(c *gin.Context, s int, err error) {
	c.JSON(s, gin.H{
		"error":   true,
		"message": err.Error(),
	})
}

// BotInfo содержит информацию о боте для API
type BotInfo struct {
	ID         string    `json:"id"`
	Room       string    `json:"room"`
	BotName    string    `json:"botName"`
	Server     string    `json:"server"`
	AuthMethod string    `json:"authMethod"`
	Status     string    `json:"status"` // Статус бота: running, stopped, starting, stopping
	LastUpdate time.Time `json:"lastUpdate"`
}

type HttpServer struct {
	bots   map[string]*Bot
	router *gin.Engine
}

func (h *HttpServer) AddBot(b *Bot) {
	h.bots[b.ID] = b
}

// getAuthMethod определяет метод авторизации бота
func getAuthMethod(b *Bot) string {
	if b.JWTAppID != "" && b.JWTAppSecret != "" {
		return "JWT"
	}
	if b.Username != "" || b.Pass != "" {
		return "Password"
	}
	return "None"
}

func (h *HttpServer) Start(l string) error {
	return h.router.Run(l)
}

// @BasePath /api/v1

// ListBots godoc
// @Summary      List bots
// @Description  get bots with full information
// @Tags         main
// @Accept       json
// @Produce      json
// @Success      200  {array}   BotInfo
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /bots [get]
func (h *HttpServer) ListBots(c *gin.Context) {
	botInfos := make([]BotInfo, 0, len(h.bots))
	for _, bot := range h.bots {
		botInfos = append(botInfos, BotInfo{
			ID:         bot.ID,
			Room:       bot.Room,
			BotName:    bot.BotName,
			Server:     bot.JitsiServer,
			AuthMethod: getAuthMethod(bot),
			Status:     bot.GetStatus(),
			LastUpdate: time.Now(),
		})
	}
	c.JSON(http.StatusOK, botInfos)
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
	bot, ok := h.bots[id]
	if !ok {
		newError(c, http.StatusBadRequest, errors.New("not found"))
		return
	}

	var res string
	err := chromedp.Run(bot.Ctx,
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
	bot, ok := h.bots[id]
	if !ok {
		newError(c, http.StatusBadRequest, errors.New("not found"))
		return
	}

	// Проверяем статус бота
	status := bot.GetStatus()
	if status != "running" {
		log.Printf("Попытка сделать скриншот бота %s в статусе %s", id, status)
		newError(c, http.StatusServiceUnavailable, fmt.Errorf("bot is not running (status: %s)", status))
		return
	}

	var buf []byte
	err := chromedp.Run(bot.Ctx,
		chromedp.FullScreenshot(&buf, 100),
	)
	if err != nil {
		log.Printf("Ошибка создания скриншота для бота %s: %v", id, err)
		newError(c, http.StatusInternalServerError, fmt.Errorf("failed to capture screenshot: %v", err))
		return
	}

	c.Data(http.StatusOK, "image/png", buf)
}

// StopBot godoc
// @Summary      Stop bot
// @Description  stop a bot by ID
// @Tags         bot
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Bot ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /:id/stop [post]
func (h *HttpServer) StopBot(c *gin.Context) {
	id := c.Param("id")
	log.Printf("Получен запрос на остановку бота с ID: %s", id)

	if id == "" {
		newError(c, http.StatusBadRequest, errors.New("id required"))
		return
	}
	bot, ok := h.bots[id]
	if !ok {
		log.Printf("Бот с ID %s не найден", id)
		newError(c, http.StatusNotFound, errors.New("bot not found"))
		return
	}

	log.Printf("Остановка бота %s (%s), текущий статус: %s", bot.BotName, id, bot.GetStatus())
	err := bot.Stop()
	if err != nil {
		log.Printf("Ошибка остановки бота %s: %v", id, err)
		newError(c, http.StatusInternalServerError, err)
		return
	}

	log.Printf("Бот %s успешно остановлен, новый статус: %s", id, bot.GetStatus())
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Bot stopped successfully",
		"status":  bot.GetStatus(),
	})
}

// RestartBot godoc
// @Summary      Restart bot
// @Description  restart a bot by ID
// @Tags         bot
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Bot ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /:id/restart [post]
func (h *HttpServer) RestartBot(c *gin.Context) {
	id := c.Param("id")
	log.Printf("Получен запрос на перезапуск бота с ID: %s", id)

	if id == "" {
		newError(c, http.StatusBadRequest, errors.New("id required"))
		return
	}
	bot, ok := h.bots[id]
	if !ok {
		log.Printf("Бот с ID %s не найден", id)
		newError(c, http.StatusNotFound, errors.New("bot not found"))
		return
	}

	log.Printf("Перезапуск бота %s (%s), текущий статус: %s", bot.BotName, id, bot.GetStatus())

	// Запускаем перезапуск в горутине, чтобы не блокировать HTTP ответ
	go func() {
		err := bot.Restart()
		if err != nil {
			log.Printf("Ошибка перезапуска бота %s: %v", id, err)
		}
	}()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Bot restart initiated",
		"status":  bot.GetStatus(),
	})
}

// BasicAuthMiddleware создает middleware для базовой авторизации
func BasicAuthMiddleware(username, password string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Если логин и пароль не установлены, пропускаем авторизацию
		if username == "" || password == "" {
			c.Next()
			return
		}

		user, pass, hasAuth := c.Request.BasicAuth()
		if !hasAuth || user != username || pass != password {
			c.Header("WWW-Authenticate", "Basic realm=\"Authorization Required\"")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	}
}

func NewHttpServer(webUsername, webPassword string) *HttpServer {
	srv := HttpServer{bots: map[string]*Bot{}, router: gin.Default()}

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
	// Применяем BasicAuth middleware к API endpoints
	v1.Use(BasicAuthMiddleware(webUsername, webPassword))
	{
		v1.GET("/bots", srv.ListBots)
		v1.GET("/:id/html", srv.HTML)
		v1.GET("/:id/screenshot", srv.Screenshot)
		v1.POST("/:id/stop", srv.StopBot)
		v1.POST("/:id/restart", srv.RestartBot)
	}
	srv.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return &srv
}
