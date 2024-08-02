package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var quit chan os.Signal
var roomAddr string
var botStarted chan bool
var Bot *tgbotapi.BotAPI

type fileInfoResp struct {
	Ok     bool `json:"ok"`
	Result struct {
		FileId       string `json:"file_id"`
		FileUniqueId string `json:"file_unique_id"`
		FileSize     int    `json:"file_size"`
		FilePath     string `json:"file_path"`
	} `json:"result"`
}

func StartRoomBot(apiKey string) string {
	if len(roomAddr) > 0 {
		return roomAddr
	}

	botStarted = make(chan bool, 1)
	go startBot(apiKey)

	if <-botStarted {
		return roomAddr
	} else {
		return ""
	}
}

func StartRoom() string {
	port, _ := getFreePort()

	go startServer(strconv.Itoa(port))

	roomAddr = getAddr(strconv.Itoa(port))

	return roomAddr
}

func RoomExist() string {
	if len(roomAddr) > 0 {
		return roomAddr
	}
	return ""
}

func StopServer() bool {
	quit <- os.Kill
	return true
}

func startBot(apiKey string) {
	defer func() {
		roomAddr = ""
	}()

	var err error
	Bot, err = tgbotapi.NewBotAPI(apiKey)
	if err != nil {
		BotLogger.Error(err.Error())
		botStarted <- false
		return
	}

	BotLogger.Info(Bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 5

	updates := Bot.GetUpdatesChan(u)

	go func() {
		botMessageHandle(updates)
	}()

	botStarted <- true

	roomAddr = apiKey

	quit = make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	BotLogger.Info("Shutdown Bot ...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	Bot.StopReceivingUpdates()

	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		BotLogger.Info("timeout of 5 seconds.")
	}

	BotLogger.Info("Bot exiting")
}

func startServer(port string) {
	defer func() {
		roomAddr = ""
	}()

	router := gin.Default()
	router.LoadHTMLFiles("./user-app/dist/spa/index.html")
	router.Static("/assets", "./user-app/dist/spa/assets/")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router.Handler(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit = make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
}

func getAddr(port string) string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String() + ":" + port
			}
		}
	}
	return ""
}

func getFreePort() (port int, err error) {
	var a *net.TCPAddr
	if a, err = net.ResolveTCPAddr("tcp", "localhost:0"); err == nil {
		var l *net.TCPListener
		if l, err = net.ListenTCP("tcp", a); err == nil {
			defer l.Close()
			return l.Addr().(*net.TCPAddr).Port, nil
		}
	}
	return
}

func botMessageHandle(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message != nil {
			var msg tgbotapi.MessageConfig

			switch update.Message.Text {
			case "/start":
				{
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Перишли файл или сообщение с файлом музыки которую нужно добавить в очередь)")
					_, err := Bot.Send(msg)
					if err != nil {
						return
					}
				}
			default:
				{
					if update.Message.Audio != nil {
						BotLogger.Info(update.Message.Audio.FileID)
						go downloadFile(update.Message.Audio)
						msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Спасибо! Ваш аудиофайл поставлен на скачивание")
						_, err := Bot.Send(msg)
						if err != nil {
							return
						}
					}

					if update.Message.Document != nil {
						msg = tgbotapi.NewMessage(update.Message.Chat.ID, "К сожалению, поддерживается только mp3 :(")
						_, err := Bot.Send(msg)
						if err != nil {
							return
						}
					}

				}
			}

		}
	}

	return
}

func downloadFile(file *tgbotapi.Audio) {
	getFilePathUrl := fmt.Sprintf("https://api.telegram.org/bot%s/getFile?file_id=%s", roomAddr, file.FileID)
	BotLogger.Info(getFilePathUrl)
	resp, err := http.Get(getFilePathUrl)
	object := fileInfoResp{}

	if err != nil {
		BotLogger.Error(err.Error())
	}

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		BotLogger.Error(err.Error())
	}
	BotLogger.Info(string(resBody))

	json.Unmarshal(resBody, &object)

	getFilePathUrl = fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", roomAddr, object.Result.FilePath)
	BotLogger.Info(getFilePathUrl)
	resp, err = http.Get(getFilePathUrl)

	if err != nil {
		BotLogger.Error(err.Error())
	}
	var filePath string
	if len(file.Performer) > 0 {
		filePath = "downloads/" + file.Performer + " - " + file.Title + ".mp3"
	} else {
		filePath = "downloads/" + file.FileName
	}

	out, _ := os.Create(filePath)

	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		BotLogger.Info(err.Error())
	}

	UpdateChan <- Track{
		Label: out.Name(),
		Path:  filePath,
	}

}
