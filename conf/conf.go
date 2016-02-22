package conf

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/gpmgo/gopm/modules/goconfig"
)

type DBConfig struct {
	Driver string
	Host   string
	Port   int
	User   string
	DBName string
	PassWd string
}

var (
	Mode           string
	HttpAddr       string
	DatabaseConfig = &DBConfig{}
	DataExpires    int

	UserPassCodeEncryptKey string

	WebDebugMode bool
	DebugMode    bool
	LogLevel     string

	configFile   = flag.String("config", "__unset__", "service config file")
	maxThreadNum = flag.Int("max-thread", 0, "max threads of service")
	debugMode    = flag.Bool("debug", false, "debug mode")
	webDebugMode = flag.Bool("web-debug", false, "web debug mode")
	logLevel     = flag.String("log-level", "INFO", "DEBUG | INFO | WARN | ERROR | FATAL | PANIC")
)

func init() {
	flag.Parse()

	DebugMode = *debugMode
	LogLevel = *logLevel
	WebDebugMode = *webDebugMode

	if len(os.Args) == 2 {
		if os.Args[1] == "reload" {
			wd, _ := os.Getwd()
			pidFile, err := os.Open(filepath.Join(wd, "instafig.pid"))
			if err != nil {
				log.Printf("Failed to open pid file: %s", err.Error())
				os.Exit(1)
			}
			pids := make([]byte, 10)
			n, err := pidFile.Read(pids)
			if err != nil {
				log.Printf("Failed to read pid file: %s", err.Error())
				os.Exit(1)
			}
			if n == 0 {
				log.Printf("No pid in pid file: %s", err.Error())
				os.Exit(1)
			}
			_, err = exec.Command("kill", "-USR2", string(pids[:n])).Output()
			if err != nil {
				log.Printf("Failed to restart Instafig service: %s", err.Error())
				os.Exit(1)
			}
			pidFile.Close()
			os.Exit(0)
		}
	}

	if DebugMode {
		LogLevel = "DEBUG"
	}

	if *maxThreadNum == 0 {
		*maxThreadNum = runtime.NumCPU()
	}
	runtime.GOMAXPROCS(*maxThreadNum)

	if *configFile == "__unset__" {
		p, _ := os.Getwd()
		*configFile = filepath.Join(p, "conf/config.ini")
	}

	confFile, err := filepath.Abs(*configFile)
	if err != nil {
		log.Printf("No correct config file: %s - %s", *configFile, err.Error())
		os.Exit(1)
	}

	config, err := goconfig.LoadConfigFile(confFile)
	if err != nil {
		log.Printf("No correct config file: %s - %s", *configFile, err.Error())
		os.Exit(1)
	}

	Mode, _ = config.GetValue("", "mode")

	HttpAddr, _ = config.GetValue("", "addr")
	UserPassCodeEncryptKey, _ = config.GetValue("", "user_passcode_encrypt_key")

	if !DebugMode {
		// disable all console log
		nullFile, _ := os.Open(os.DevNull)
		log.SetOutput(nullFile)
		os.Stdout = nullFile
	}
}
