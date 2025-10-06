package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"sheff.online/ssjitsi/internal/pkg/ssjitsi"
)

// Читаем конфиг из файла (по умолчанию ssjitsi.yaml), создаем ботов и запускаем http сервер.
func main() {
	configFile := flag.String("config", "ssjitsi.yaml", "Путь к файлу конфигурации")
	help := flag.Bool("help", false, "Показать справку")
	flag.Parse()

	if *help {
		fmt.Println("JitsiBot Server - сервер для управления ботами Jitsi Meet")
		fmt.Println("\nИспользование:")
		flag.PrintDefaults()
		fmt.Println("\nПример:")
		fmt.Println("  server -config ssjitsi.yaml")
		os.Exit(0)
	}

	// Загружаем конфигурацию
	config, err := ssjitsi.LoadConfig(*configFile)
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Создаем HTTP сервер с авторизацией
	server := ssjitsi.NewHttpServer(config.WebUsername, config.WebPassword)

	// Создаем embedded сервер с встроенным UI и авторизацией
	router := ssjitsi.NewEmbeddedServer(server, config.WebUsername, config.WebPassword)

	// Запускаем HTTP сервер в отдельной горутине
	log.Printf("Запуск HTTP сервера на %s", config.HTTP)
	log.Printf("Web UI доступен по адресу http://localhost%s", config.HTTP)

	go func() {
		err := router.Run(config.HTTP)
		if err != nil {
			log.Fatalf("Ошибка запуска HTTP сервера: %v", err)
		}
	}()

	// Создаем и запускаем ботов
	for i := range config.Bots {
		botConfig := &config.Bots[i]
		// Генерируем уникальный ID для бота
		botConfig.ID = uuid.New().String()

		// Устанавливаем начальный статус
		botConfig.SetStatus("stopped")

		log.Printf("Запуск бота %d: комната '%s', имя '%s'", i+1, botConfig.Room, botConfig.BotName)

		// Добавляем бота в сервер сразу
		server.AddBot(botConfig)

		// Запускаем бота в горутине
		go func(bot *ssjitsi.Bot, index int) {
			err := bot.Start()
			if err != nil {
				log.Printf("Ошибка работы бота %d: %v", index+1, err)
			}
		}(botConfig, i)

		log.Printf("Бот %d запускается (ID: %s)", i+1, botConfig.ID)
	}

	// Ждем завершения (блокируем main)
	select {}
}
