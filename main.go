package main

import (
	"changeme/internal/client/win32"
	"changeme/internal/client/wmi"
	"changeme/internal/database"
	"changeme/internal/repository/sqlite"
	"changeme/internal/usecase"
	"context"
	"embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func init() {
	//application.RegisterEvent[string]("name")
}

func main() {
	// Пример события wails
	//app.Event.Emit("name", data)

	// Инициализация БД
	db, err := database.Open()
	if err != nil {
		log.Fatal(err)
	}
	// Создание таблиц в БД
	if err = database.Migrate(context.Background(), db); err != nil {
		log.Fatal(err)
	}

	// Репозиторий
	processRepository := sqlite.NewProcessRepository(db)

	// Клиент
	win32Client := win32.NewClient()
	wmiClient := wmi.NewClient()

	// Бизнес логика
	InitialStateCollectorUsecase := usecase.NewInitialStateCollector(win32Client, wmiClient, processRepository)
	err = InitialStateCollectorUsecase.CollectProcesses()
	if err != nil {
		log.Fatal(err)
	}

	// Мост с пользовательской частью
	frontend := NewApp()

	// Конструктор приложения
	app := application.New(application.Options{
		Name:        "net-monitor",
		Description: "",
		Services: []application.Service{
			application.NewService(frontend),
		},
		OnShutdown: func() {
			// Закрытие БД
			if err = database.Close(db); err != nil {
				log.Fatal(err)
			}
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
	})

	// Основное окно приложения
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:            "Сетевой монитор",
		Width:            1000,
		Height:           618,
		BackgroundColour: application.NewRGB(255, 255, 255),
		URL:              "/",
	})

	// Запуск приложения
	if err = app.Run(); err != nil {
		log.Fatal(err)
	}
}
