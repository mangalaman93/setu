package main

import (
	"log"
	"os"
	"os/signal"
	"path"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/sevlyar/go-daemon"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	days = []time.Time{time.Now(), time.Now().Add(24 * time.Hour)}
)

func main() {
	districtIDstr := os.Getenv("DISTRICT_ID")
	districtID, err := strconv.ParseInt(districtIDstr, 10, 64)
	if err != nil {
		panic(err)
	}

	centersStr := os.Getenv("CENTERS")
	centers := make([]int, 0)
	for _, s := range strings.Split(centersStr, ",") {
		center, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			panic(err)
		}
		centers = append(centers, int(center))
	}

	// setup setu dir
	homeFolder := os.Getenv("HOME")
	if homeFolder == "" {
		panic("home folder env not found")
	}
	setuDir := path.Join(homeFolder, "setu")

	// setup log dir
	logDir := path.Join(setuDir, "logs")
	if err := os.MkdirAll(logDir, 0744); err != nil {
		panic(err)
	}

	// daemonize
	ctx := &daemon.Context{
		PidFileName: path.Join(setuDir, "setu.pid"),
		PidFilePerm: 0644,
	}
	child, err := ctx.Reborn()
	if err != nil {
		panic(err)
	}

	if child != nil {
		log.Println("[INFO] running the service as a daemon")
	} else {
		defer ctx.Release()
		runChild(int(districtID), centers, logDir)
	}
}

func runChild(districtID int, centers []int, logDir string) {
	// log setup
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	logFilePath := path.Join(logDir, "setu.log")
	log.SetOutput(&lumberjack.Logger{
		Filename:  logFilePath,
		LocalTime: true,
	})
	log.Println("#################### BEGIN OF LOG ##########################")

	// register ctrl+c
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	log.Println("[INFO] adding signal handler for SIGTERM")

	// loop
	ticker := time.NewTicker(5 * time.Minute)
	for {
		select {
		case <-ticker.C:
			slots, err := getSlotsForDays(districtID, centers, days)
			if err != nil {
				log.Printf("[ERROR] error ocurred: %v", err)
				email(0, err)
			} else if slots > 0 {
				log.Printf("[INFO] found %v empty slots", slots)
				email(slots, nil)
			} else {
				log.Println("[INFO] no available slots found")
			}
		case <-sigs:
			return
		}
	}
}
