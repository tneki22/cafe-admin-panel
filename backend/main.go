package main

import (
	"flag"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	// Чтение пути к файлу конфигурации
	configPath := flag.String("config-path", "./auth.yaml", "путь к файлу конфигурации")
	flag.Parse()

	// Загрузка конфигурации
	cfg, err := LoadConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	// Инициализация базы данных
	dbProvider := NewProvider(cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.DBname)
	if dbProvider == nil {
		log.Fatal("Failed to initialize database provider")
	}

	// Инициализация JWT провайдера
	jwtProvider := NewJWTProvider(cfg.JWT.Secret)

	// Инициализация бизнес-логики
	usecase := NewUsecase(cfg.Usecase.DefaultMessage, *dbProvider, *jwtProvider)

	// Инициализация сервера
	server := NewServer(cfg.IP, cfg.Port, cfg.API.MinPasswordSize, cfg.API.MaxPasswordSize, cfg.API.MinUsernameSize, cfg.API.MaxUsernameSize, cfg.JWT.Secret, *usecase)

	// Запуск сервера
	server.Run()
}
