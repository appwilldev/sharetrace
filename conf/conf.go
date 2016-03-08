package conf

import (
	"flag"
	//"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/gpmgo/gopm/modules/goconfig"
)

type DBConfig struct {
	BaseDir string
	Host    string
	Port    int
	DBName  string
	User    string
	PassWd  string
}

type RedisConfig struct {
	Host string
	Port int
}

const (
	DEFAULT_CACHE_DB_NAME          = "default"
	DEFAULT_CACHE_EXPIRE           = 7 * 24 * 60 * 60
	UserExpires                    = 7 * 24 * 60 * 60
	DumpExpiresDuration            = 60
	CLICK_SESSION_STATUS_INSTALLED = 1
	CLICK_SESSION_STATUS_BUTTON    = 2
	CLICK_SESSION_STATUS_CLICK     = 0
	CLICK_TYPE_COOKIE              = 0
	CLICK_TYPE_IP                  = 1
	COOKIE_PREFIX                  = "st"
)

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
	debugMode    = flag.Bool("debug", true, "debug mode")
	webDebugMode = flag.Bool("web-debug", true, "web debug mode")
	logLevel     = flag.String("log-level", "INFO", "DEBUG | INFO | WARN | ERROR | FATAL | PANIC")

	RedConfig = make(map[string]*RedisConfig)
)

func init() {
	flag.Parse()

	DebugMode = *debugMode
	LogLevel = *logLevel
	WebDebugMode = *webDebugMode

	if len(os.Args) == 2 {
		if os.Args[1] == "reload" {
			wd, _ := os.Getwd()
			pidFile, err := os.Open(filepath.Join(wd, "sharetrace.pid"))
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
				log.Printf("Failed to restart service: %s", err.Error())
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

	Mode, _ = config.GetValue("http", "mode")

	HttpAddr, _ = config.GetValue("http", "addr")
	UserPassCodeEncryptKey, _ = config.GetValue("http", "user_passcode_encrypt_key")

	DatabaseConfig.BaseDir, _ = config.GetValue("postgres", "base_dir")
	DatabaseConfig.Host, _ = config.GetValue("postgres", "host")
	port, _ := config.GetValue("postgres", "port")
	DatabaseConfig.Port, err = strconv.Atoi(port)
	if err != nil {
		log.Printf("DB port is not correct: %s - %s", *configFile, err.Error())
		os.Exit(1)
	}
	DatabaseConfig.DBName, _ = config.GetValue("postgres", "db_name")
	DatabaseConfig.User, _ = config.GetValue("postgres", "user")
	DatabaseConfig.PassWd, _ = config.GetValue("postgres", "passwd")

	redisConfig := make(map[string]*RedisConfig, 1)
	dbName := "default"
	redisConfig[dbName] = new(RedisConfig)
	redisConfig[dbName].Host, _ = config.GetValue("redis", "host")
	redPort, _ := config.GetValue("redis", "port")
	redisConfig[dbName].Port, _ = strconv.Atoi(redPort)
	RedConfig = redisConfig

	if !DebugMode {
		// disable all console log
		nullFile, _ := os.Open(os.DevNull)
		log.SetOutput(nullFile)
		os.Stdout = nullFile
	}

}
