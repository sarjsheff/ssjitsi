package ssjitsi

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/golang-jwt/jwt/v5"
)

type Bot struct {
	ID           string             `yaml:"ID,omitempty"`
	Room         string             `yaml:"Room"`
	BotName      string             `yaml:"BotName"`
	DataDir      string             `yaml:"DataDir"`
	JitsiServer  string             `yaml:"JitsiServer"`
	Username     string             `yaml:"Username"`
	Pass         string             `yaml:"Pass"`
	JWTAppID     string             `yaml:"JWTAppID"`
	JWTAppSecret string             `yaml:"JWTAppSecret"`
	Headless     bool               `yaml:"Headless"`
	Ctx          context.Context    `yaml:"-"`
	CtxCancel    context.CancelFunc `yaml:"-"`
}
type Record struct {
	U      string `json:"u"`
	D      string `json:"d"`
	User   string `json:"user"`
	UserId string `json:"userid"`
	Room   string `json:"room"`
	Myid   string `json:"myid"`
}

// GenerateJitsiJWT генерирует JWT токен для авторизации в Jitsi Meet
func GenerateJitsiJWT(appID, appSecret, jitsiServer, room, userName string) (string, error) {
	// Извлекаем домен из URL сервера Jitsi
	parsedURL, err := url.Parse(jitsiServer)
	if err != nil {
		return "", fmt.Errorf("failed to parse Jitsi server URL: %v", err)
	}
	domain := parsedURL.Hostname()

	// Создаем claims для токена
	now := time.Now()
	claims := jwt.MapClaims{
		"iss":  appID,                         // Issuer - идентификатор приложения
		"aud":  appID,                         // Audience - идентификатор приложения
		"sub":  domain,                        // Subject - домен Jitsi сервера
		"room": room,                          // Имя комнаты
		"exp":  now.Add(2 * time.Hour).Unix(), // Expiration - токен действителен 2 часа
		"nbf":  now.Unix(),                    // Not Before - токен действителен с текущего момента
		"context": map[string]interface{}{
			"user": map[string]interface{}{
				"name": userName,
			},
		},
	}

	// Создаем токен с алгоритмом HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем токен секретным ключом
	tokenString, err := token.SignedString([]byte(appSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT token: %v", err)
	}

	return tokenString, nil
}

func (bot *Bot) Start() error {
	ctx, _ := chromedp.NewContext(context.Background())
	// defer cancel()

	jsContent, err := os.ReadFile("script.js")
	if err != nil {
		fmt.Println(err)
		return err
	}

	allocCtx, _ := chromedp.NewExecAllocator(
		ctx,
		append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("use-fake-ui-for-media-stream", true),
			chromedp.Flag("headless", bot.Headless),
		)...,
	)
	// defer cancelAlloc()

	bot.Ctx, bot.CtxCancel = chromedp.NewContext(allocCtx)
	// defer cancelTask()

	chromedp.ListenTarget(bot.Ctx, func(ev interface{}) {
		switch ev := ev.(type) {
		case *runtime.EventConsoleAPICalled:
			if ev.Type == "error" {
				fmt.Printf("Console message: Type=%s, Args=%v\n", ev.Type, ev.Args)

				for _, arg := range ev.Args {
					fmt.Printf("  Arg: Type=%s, Value=%v\n", arg.Type, arg.Value)
				}
			}
		case *runtime.EventBindingCalled:
			if ev.Name == "ssbot_writeSound" {

				var p Record
				err := json.Unmarshal([]byte(ev.Payload), &p)
				if err != nil {
					fmt.Println("Error unmarshaling JSON:", err)
					return
				}
				writeRecordToFile(p, bot.DataDir, bot.ID)
			}
		case *browser.EventDownloadProgress:
			if ev.State == browser.DownloadProgressStateCompleted {
				log.Println(ev.GUID)
			}
		case *runtime.EventExceptionThrown:
			fmt.Printf("Exception thrown: %s\n", ev.ExceptionDetails.Text)
		}
	})

	params := &browser.SetPermissionParams{
		Permission: &browser.PermissionDescriptor{
			Name: "microphone",
		},
		Setting: browser.PermissionSettingDenied, //browser.PermissionSettingGranted,
	}

	// run task list
	var res string
	// var buf []byte

	// Проверяем, нужна ли JWT авторизация
	if bot.JWTAppID != "" && bot.JWTAppSecret != "" {
		// Используем JWT авторизацию
		log.Println("Используем JWT авторизацию")

		// Генерируем JWT токен
		token, err := GenerateJitsiJWT(bot.JWTAppID, bot.JWTAppSecret, bot.JitsiServer, bot.Room, bot.BotName)
		if err != nil {
			return fmt.Errorf("failed to generate JWT token: %v", err)
		}

		// Формируем URL с JWT токеном
		jitsiURL := strings.TrimRight(bot.JitsiServer, "/") + "/" + bot.Room + "?jwt=" + token
		log.Printf("Переходим на URL: %s", strings.TrimRight(bot.JitsiServer, "/")+"/"+bot.Room+"?jwt=***")

		err = chromedp.Run(bot.Ctx,
			chromedp.Navigate(jitsiURL),
			chromedp.ActionFunc(func(ctx context.Context) error {
				// Параметры: разрешение, настройка, источник (опционально)
				return params.Do(ctx)
			}),
			chromedp.Sleep(2*time.Second), // Даем время на загрузку страницы
			chromedp.Click(`[aria-label="Join meeting"]`, chromedp.ByQuery),
			chromedp.Sleep(2*time.Second), // Даем время на подключение к конференции
		)
		if err != nil {
			return err
		}
	} else {
		// Используем старый метод с формами
		log.Println("Используем авторизацию с формами")

		var nodes []*cdp.Node
		err = chromedp.Run(bot.Ctx,
			chromedp.Navigate(bot.JitsiServer),
			chromedp.Click(`[aria-label="Meeting name input"]`, chromedp.ByQuery),
			chromedp.SendKeys(`[aria-label="Meeting name input"]`, bot.Room, chromedp.ByQuery),
			chromedp.Click("#enter_room_button", chromedp.ByQuery),
			chromedp.ActionFunc(func(ctx context.Context) error {
				// Параметры: разрешение, настройка, источник (опционально)
				return params.Do(ctx)
			}),
			chromedp.Sleep(1*time.Second),
			chromedp.SendKeys(`[aria-label="Enter your name"]`, bot.BotName, chromedp.ByQuery),
			chromedp.Click(`[aria-label="Join meeting"]`, chromedp.ByQuery),
			chromedp.Click(`[aria-label="Join meeting"]`, chromedp.ByQuery),
			chromedp.Sleep(2*time.Second),
			chromedp.Nodes("#login-dialog-username", &nodes, chromedp.AtLeast(0)),
		)
		if err != nil {
			return err
		}

		// Нужна авторизация?
		if len(nodes) > 0 {
			log.Println("Авторизуемся.")
			err = chromedp.Run(bot.Ctx,
				chromedp.SendKeys("#login-dialog-username", bot.Username, chromedp.ByQuery),
				chromedp.SendKeys("#login-dialog-password", bot.Pass, chromedp.ByQuery),
				chromedp.Click(`[aria-label="Login"]`, chromedp.ByQuery),
				chromedp.Sleep(1*time.Second),
			)
			if err != nil {
				return err
			}
		}
	}

	err = chromedp.Run(bot.Ctx,
		runtime.AddBinding("ssbot_writeSound"),
		chromedp.Evaluate(string(jsContent), &res),
	)

	return err
}

func wrf(f string, d []byte) error {
	file, err := os.OpenFile(f, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()

	_, err = file.Write(d)
	if err != nil {
		return err
	}
	return nil
}

func writeRecordToFile(p Record, datadir string, sessionid string) error {
	// Декодируем base64 строку
	data, err := base64.StdEncoding.DecodeString(p.D)
	if err != nil {
		return fmt.Errorf("ошибка декодирования base64: %v", err)
	}

	udir := filepath.Join(datadir, SafeFilename(p.Room), sessionid)
	err = os.MkdirAll(udir, 0755)
	if err != nil {
		log.Println(err)
		return err
	}

	filename := filepath.Join(udir, p.UserId+"_"+p.U+".webm")
	starttime := filepath.Join(udir, p.UserId+"_"+p.U+".json")
	metadata := filepath.Join(udir, p.UserId+".json")
	room := filepath.Join(udir, "room.json")

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	_, err = os.Stat(starttime)
	if os.IsNotExist(err) {
		wrf(starttime, []byte(strconv.Itoa(int(time.Now().UnixMilli()))))
	}
	_, err = os.Stat(metadata)
	if os.IsNotExist(err) {
		wrf(metadata, []byte(p.User))
	}
	_, err = os.Stat(room)
	if os.IsNotExist(err) {
		wrf(room, []byte(p.Room))
	}
	return nil
}
