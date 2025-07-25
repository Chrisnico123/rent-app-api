package configs

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

const (
	MaxFileSize = 10 * 1024 * 1024 // 10MB
)

type Config struct {
	App        AppConfig        `yaml:"app"`
	DB         DBConfig         `yaml:"db"`
	Jwt        JwtConfig        `yaml:"jwt"`
	Cloudinary CloudinaryConfig `yaml:"cloudinary"`
	Permission Permission       `yaml:"permission"`
	Email      EmailConfig      `yaml:"email"`
	Xendit     XenditConfig     `yaml:"xendit"`
	SecretKey  SecretConfig     `yaml:"secret_key"`
}

type SecretConfig struct {
	Key string `yaml:"key"`
}

type XenditConfig struct {
	ApiKey string `yaml:"api_key"`
}

type EmailConfig struct {
	ApiKey string `yaml:"api_key"`
	From   string `yaml:"from"`
}

type AppConfig struct {
	Name string `yaml:"name"`
	Port string `yaml:"port"`
	Host string `yaml:"host"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	// Connection Pool Database
	ConnectionPool DBConnectionPool `yaml:"connection_pool"`
}

type DBConnectionPool struct {
	// Max Connection Waiting
	MaxIdleConnection uint8 `yaml:"max_idle_connection"`
	// Max Connection Open to Database
	MaxOpenConnection uint8 `yaml:"max_open_connection"`
	// Max Lifetime Connection to Open Connection to Database
	MaxLifeTimeConnection uint8 `yaml:"max_life_time_connection"`
	// Max Idle Lifetime Connection to Wait until connect into database
	MaxIdleTimeConnection uint8 `yaml:"max_idle_time_connection"`

	DB_POOL_MIN string `yaml:"db_pool_min"`
	DB_POOL_MAX string `yaml:"db_pool_max"`
	DB_TIMEOUT  string `yaml:"db_time_out"`
}

type JwtConfig struct {
	SessionLogin string `yaml:"sessionLogin"`
	SecretKey    string `yaml:"secretKey"`
}

type CloudinaryConfig struct {
	CloudName      string `yaml:"cloudinary_name"`
	CloudApiKey    string `yaml:"cloudinary_api_key"`
	CloudApiSecret string `yaml:"cloudinary_api_secret"`
}

type Permission struct {
	Admin     string `yaml:"admin"`
	User      string `yaml:"user"`
	UserAdmin string `yaml:"user_admin"`
}

var Cfg Config

// LoadConfig tries to load from YAML, if not found, load from env
func LoadConfig(filename string) (err error) {
	// Check if running in Docker/Production
	if os.Getenv("APP_ENV") != "development" {
		log.Println("Running in production mode, skipping .env file")
	} else {
		// Try YAML config first
		if configByte, err := os.ReadFile(filename); err == nil {
			return yaml.Unmarshal(configByte, &Cfg)
		}

		// Local development - try loading .env
		if err := loadDotEnv(); err != nil {
			log.Printf("Warning: couldn't load .env file: %v", err)
		}
	}

	loadFromEnv()
	return nil
}

// loadDotEnv tries to load .env file from project root
func loadDotEnv() error {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// If we're in the cmd directory, go up one level
	if filepath.Base(cwd) == "cmd" {
		cwd = filepath.Dir(cwd)
	}

	envPath := filepath.Join(cwd, ".env")
	return godotenv.Load(envPath)
}

func loadFromEnv() {
	Cfg.App.Name = getEnv("APP_NAME", "")
	Cfg.App.Port = getEnv("APP_PORT", "8080")
	Cfg.App.Host = getEnv("APP_HOST", "localhost")

	Cfg.DB.Host = getEnv("DB_HOST", "")
	Cfg.DB.Port = getEnv("DB_PORT", "")
	Cfg.DB.User = getEnv("DB_USER", "")
	Cfg.DB.Password = getEnv("DB_PASSWORD", "")
	Cfg.DB.Name = getEnv("DB_NAME", "")
	Cfg.DB.ConnectionPool.MaxIdleConnection = parseUint8Env("DB_MAX_IDLE_CONNECTION", 10)
	Cfg.DB.ConnectionPool.MaxOpenConnection = parseUint8Env("DB_MAX_OPEN_CONNECTION", 20)
	Cfg.DB.ConnectionPool.MaxLifeTimeConnection = parseUint8Env("DB_MAX_LIFE_TIME_CONNECTION", 20)
	Cfg.DB.ConnectionPool.MaxIdleTimeConnection = parseUint8Env("DB_MAX_IDLE_TIME_CONNECTION", 20)
	Cfg.DB.ConnectionPool.DB_POOL_MIN = getEnv("DB_POOL_MIN", "10")
	Cfg.DB.ConnectionPool.DB_POOL_MAX = getEnv("DB_POOL_MAX", "50")
	Cfg.DB.ConnectionPool.DB_TIMEOUT = getEnv("DB_TIME_OUT", "30s")

	Cfg.Email.ApiKey = getEnv("EMAIL_API_KEY", "")
	Cfg.Email.From = getEnv("EMAIL_FROM", "")

	Cfg.Jwt.SecretKey = getEnv("JWT_SECRET_KEY", "")
	Cfg.Jwt.SessionLogin = getEnv("JWT_SESSION_LOGIN", "")

	Cfg.Cloudinary.CloudName = getEnv("CLOUDINARY_NAME", "")
	Cfg.Cloudinary.CloudApiKey = getEnv("CLOUDINARY_API_KEY", "")
	Cfg.Cloudinary.CloudApiSecret = getEnv("CLOUDINARY_API_SECRET", "")

	Cfg.Permission.Admin = getEnv("PERMISSION_ADMIN", "admin")
	Cfg.Permission.User = getEnv("PERMISSION_USER", "user")
	Cfg.Permission.UserAdmin = getEnv("PERMISSION_USER_ADMIN", "user_admin")

	Cfg.Xendit.ApiKey = getEnv("API_KEY_XENDIT", "")
	Cfg.SecretKey.Key = getEnv("APP_ENCRYPTION_KEY", "")
}

// Helper function to get environment variable with default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func parseIntEnv(key string, def int) int {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	v, err := strconv.Atoi(val)
	if err != nil {
		return def
	}
	return v
}

func parseUint8Env(key string, def uint8) uint8 {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	v, err := strconv.ParseUint(val, 10, 8)
	if err != nil {
		return def
	}
	return uint8(v)
}

func NewServerConfig() string {
	return fmt.Sprintf(`%s:%s`, Cfg.App.Host, Cfg.App.Port)
}
