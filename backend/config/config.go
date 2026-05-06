package config

import (
	"dragonroll/backend/utils"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Server ServerConfig
	Mongo  MongoConfig
	JWT    JWTConfig
	Logger LoggerConfig
	Cors   CorsConfig
}

type ServerConfig struct {
	Mode string // "debug" | "release"
	Port string
}

type MongoConfig struct {
	User string
	Pass string
	Host string
	Port string
	Name string
	URI  string
}

type JWTConfig struct {
	SigningAlg string
	Key        string
}

type LoggerConfig struct {
	Filepath string
	Filename string
	Level    zapcore.Level
	Tz       string
}

type CorsConfig struct {
	CorsUrl string
}

var Cfg *Config

func LoadConfig() {
	Cfg = &Config{}

	Cfg.Server.Mode = utils.GetEnvWithDefault("GIN_MODE", "debug")
	gin.SetMode(Cfg.Server.Mode)
	Cfg.Server.Port = utils.GetEnvWithDefault("PORT", "8080")

	loadMongoConfig(&Cfg.Mongo)
	loadJWTConfig(&Cfg.JWT)
	loadLoggerConfig(&Cfg.Logger)
	loadCorsConfig(&Cfg.Cors)
}

func loadMongoConfig(m *MongoConfig) {
	m.User = os.Getenv("DB_USER_MONGO")
	m.Pass = os.Getenv("DB_PASS_MONGO")
	m.Host = utils.GetEnvWithDefault("DB_HOST_MONGO", "localhost")
	m.Port = utils.GetEnvWithDefault("DB_PORT_MONGO", "27017")
	m.Name = utils.GetEnvWithDefault("DB_NAME_MONGO", "dragonroll")

	// URI con o sin auth
	if m.User != "" && m.Pass != "" {
		m.URI = fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=admin",
			m.User, m.Pass, m.Host, m.Port, m.Name)
	} else {
		m.URI = fmt.Sprintf("mongodb://%s:%s/%s", m.Host, m.Port, m.Name)
	}
}

func loadJWTConfig(j *JWTConfig) {
	// Solo soportamos HS256/HS384/HS512 — simétricos. Es lo que necesita
	// una app self-hosted; para asimétrico se reintroduce el path de keys.
	j.SigningAlg = utils.GetEnvWithDefault("JWT_ALG", "HS256")
	j.Key = os.Getenv("JWT_KEY")
	if j.Key == "" {
		log.Fatal("CONFIG: JWT_KEY es obligatorio. Generá uno con: openssl rand -hex 32")
	}
	if len(j.Key) < 32 {
		log.Fatal("CONFIG: JWT_KEY debe tener al menos 32 caracteres.")
	}
}

func loadLoggerConfig(l *LoggerConfig) {
	l.Filepath = utils.GetEnvWithDefault("LOG_FILEPATH", "/app/logs/")
	l.Filename = utils.GetEnvWithDefault("LOG_FILENAME", "dragonroll")
	l.Tz = utils.GetEnvWithDefault("TZ", "America/Santiago")

	if Cfg.Server.Mode == "debug" {
		l.Level = zapcore.DebugLevel
	} else {
		l.Level = zapcore.InfoLevel
	}
}

func loadCorsConfig(c *CorsConfig) {
	c.CorsUrl = utils.GetEnvWithDefault("CORS_URLS", "http://localhost:3000")
}
