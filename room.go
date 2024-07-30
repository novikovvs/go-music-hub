package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/tnychn/gotube"
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
var quitState chan bool
var roomAddr string

type SubmitRequest struct {
	VideoId string `json:"video_id"`
	Url     string `json:"url"`
}

func StartRoom() string {
	quitState = make(chan bool, 1)
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
	return <-quitState
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
	router.POST("/submit", func(c *gin.Context) {

		video, err := gotube.NewVideo("9vc-I9rvGsw", true)

		c.JSON(200, gin.H{
			"message": "pong",
		})
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
	quitState <- true
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
