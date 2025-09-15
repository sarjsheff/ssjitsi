package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

type Record struct {
	U      string `json:"u"`
	D      string `json:"d"`
	User   string `json:"user"`
	UserId string `json:"userid"`
	Room   string `json:"room"`
	Myid   string `json:"myid"`
}

func main() {

	room := flag.String("room", "ssjitsi-test", "Название комнаты для подключения")
	botname := flag.String("botname", "SSJitsiBot", "Имя бота в комнате")
	datadir := flag.String("datadir", "../data/", "Директория для сохранения данных")
	jitsiServer := flag.String("jitsi", "https://meet.jit.si/", "Адрес jitsi сервера")
	username := flag.String("username", "", "Имя пользователя (если нужна авторизация)")
	pass := flag.String("pass", "", "Пароль (если нужна авторизация)")
	help := flag.Bool("help", false, "Показать справку")

	flag.Parse()

	if *help {
		fmt.Println("JitsiBot - бот для записи звука с Jitsi Meet")
		fmt.Println("\nИспользование:")
		flag.PrintDefaults()
		fmt.Println("\nПример:")
		fmt.Println("  jitsibot -room myroom -botname 'My Bot' -datadir ./data")
		os.Exit(0)
	}

	ch := make(chan Record, 100)
	go func(chn chan Record) {
		for {
			val, ok := <-chn
			if !ok {
				fmt.Println("Channel is closed and empty.")
				break
			}
			writeRecordToFile(val, *datadir)
		}
	}(ch)

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	downloadDir, err := os.MkdirTemp("", "chromedp-downloads")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Downloads will be saved to: %s", downloadDir)

	jsContent, err := os.ReadFile("script.js")
	if err != nil {
		fmt.Println(err)
		return
	}

	allocCtx, cancelAlloc := chromedp.NewExecAllocator(
		ctx,
		append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("use-fake-ui-for-media-stream", true),
			chromedp.Flag("headless", false),
		)...,
	)
	defer cancelAlloc()

	taskCtx, cancelTask := chromedp.NewContext(allocCtx)
	defer cancelTask()

	chromedp.ListenTarget(taskCtx, func(ev interface{}) {
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
				writeRecordToFile(p, *datadir)
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
	var buf []byte
	var nodes []*cdp.Node
	err = chromedp.Run(taskCtx,
		chromedp.Navigate(*jitsiServer),
		chromedp.Click(`[aria-label="Meeting name input"]`, chromedp.ByQuery),
		chromedp.SendKeys(`[aria-label="Meeting name input"]`, *room, chromedp.ByQuery),
		chromedp.Click("#enter_room_button", chromedp.ByQuery),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// Параметры: разрешение, настройка, источник (опционально)
			return params.Do(ctx)
		}),
		chromedp.Sleep(1*time.Second),
		chromedp.SendKeys(`[aria-label="Enter your name"]`, *botname, chromedp.ByQuery),
		chromedp.Click(`[aria-label="Join meeting"]`, chromedp.ByQuery),
		chromedp.Click(`[aria-label="Join meeting"]`, chromedp.ByQuery),
		chromedp.Sleep(2*time.Second),
		chromedp.Nodes("#login-dialog-username", &nodes, chromedp.AtLeast(0)),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Printf("Авторизуемся? %d\n", len(nodes))

	// Нужна авторизация?
	if len(nodes) > 0 {
		log.Println("Авторизуемся")
		err = chromedp.Run(taskCtx,
			chromedp.SendKeys("#login-dialog-username", *username, chromedp.ByQuery),
			chromedp.SendKeys("#login-dialog-password", *pass, chromedp.ByQuery),
			chromedp.Click(`[aria-label="Login"]`, chromedp.ByQuery),
			chromedp.Sleep(1*time.Second),
		)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	err = chromedp.Run(taskCtx,
		browser.
			SetDownloadBehavior(browser.SetDownloadBehaviorBehaviorAllowAndName).
			WithDownloadPath(downloadDir).
			WithEventsEnabled(true),
		runtime.AddBinding("ssbot_writeSound"),
		chromedp.Evaluate(string(jsContent), &res),
		chromedp.Sleep(10000*time.Hour),
		// chromedp.FullScreenshot(&buf, 100),
		// chromedp.OuterHTML("body", &res, chromedp.ByQuery),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := os.WriteFile("screen.png", buf, 0o644); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(strings.TrimSpace(res))
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

func writeRecordToFile(p Record, datadir string) error {
	log.Println("in2")
	// Декодируем base64 строку
	data, err := base64.StdEncoding.DecodeString(p.D)
	if err != nil {
		return fmt.Errorf("ошибка декодирования base64: %v", err)
	}

	udir := filepath.Join(datadir, p.Myid)
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
