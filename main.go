package main

import (
	"embed"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"go.uber.org/zap"
	"log"
	"os"
)

//go:embed all:frontend/dist/spa
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

var AppLogger *zap.Logger
var BotLogger *zap.Logger
var PlayerLogger *zap.Logger

func initAppLogger() {
	var err error

	err = os.MkdirAll("./logs/App/", os.ModePerm)
	if err != nil {
		panic(err)
	}

	_, err = os.Create("./logs/App/app.log")
	if err != nil {
		panic(err)
	}

	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"./logs/App/app.log"}
	AppLogger, err = config.Build()

	if err != nil {
		log.Fatal(err)
	}
}
func initBotLogger() {
	var err error

	err = os.MkdirAll("./logs/Bot/", os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	_, err = os.Create("./logs/Bot/bot.log")
	if err != nil {
		panic(err)
	}

	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"./logs/Bot/bot.log"}
	BotLogger, err = config.Build()

	if err != nil {
		log.Fatal(err)
	}
}
func initPlayerLogger() {
	var err error

	err = os.MkdirAll("./logs/Player/", os.ModePerm)
	if err != nil {
		panic(err)
	}

	_, err = os.Create("./logs/Player/player.log")
	if err != nil {
		panic(err)
	}

	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"./logs/Player/player.log"}
	PlayerLogger, err = config.Build()

	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	initAppLogger()
	initBotLogger()
	initPlayerLogger()

	UpdateChan = make(chan Track)
}

func main() {

	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "Это лечится",
		Width:  1024,
		Height: 768,
		Linux: &linux.Options{
			Icon: icon,
		},
		MinWidth:          1024,
		MinHeight:         768,
		MaxWidth:          1280,
		MaxHeight:         800,
		DisableResize:     false,
		Fullscreen:        false,
		Frameless:         false,
		StartHidden:       false,
		HideWindowOnClose: false,
		BackgroundColour:  &options.RGBA{R: 255, G: 255, B: 255, A: 255},
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Menu:             nil,
		Logger:           nil,
		LogLevel:         logger.DEBUG,
		OnStartup:        app.startup,
		OnDomReady:       app.domReady,
		OnBeforeClose:    app.beforeClose,
		OnShutdown:       app.shutdown,
		WindowStartState: options.Normal,
		Bind: []interface{}{
			app,
		},
		// Windows platform specific options
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
			// DisableFramelessWindowDecorations: false,
			WebviewUserDataPath: "",
			ZoomFactor:          1.0,
		},
		// Mac platform specific options
		Mac: &mac.Options{
			TitleBar: &mac.TitleBar{
				TitlebarAppearsTransparent: false,
				HideTitle:                  true,
				HideTitleBar:               false,
				FullSizeContent:            false,
				UseToolbar:                 false,
				HideToolbarSeparator:       true,
			},
			Appearance:           mac.NSAppearanceNameDarkAqua,
			WebviewIsTransparent: false,
			WindowIsTranslucent:  true,
			About: &mac.AboutInfo{
				Title:   "love",
				Message: "HELLO BOY",
				Icon:    icon,
			},
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}
