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

	// Создаем HTTP сервер
	server := ssjitsi.NewHttpServer()

	// Создаем embedded сервер с встроенным UI
	router := ssjitsi.NewEmbeddedServer(server)

	// Создаем и запускаем ботов
	for i, botConfig := range config.Bots {
		// Генерируем уникальный ID для бота
		botConfig.ID = uuid.New().String()

		log.Printf("Запуск бота %d: комната '%s', имя '%s'", i+1, botConfig.Room, botConfig.BotName)

		// Запускаем бота
		err := botConfig.Start()
		if err != nil {
			log.Printf("Ошибка запуска бота %d: %v", i+1, err)
			continue
		}

		// Добавляем бота в сервер
		server.AddBot(botConfig)
		log.Printf("Бот %d успешно запущен (ID: %s)", i+1, botConfig.ID)
	}

	log.Printf("Запуск embedded сервера на %s", config.HTTP)
	log.Printf("Web UI доступен по адресу http://localhost%s", config.HTTP)

	// Запускаем embedded сервер
	err = router.Run(config.HTTP)
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
