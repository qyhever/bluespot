package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Config 应用配置结构体
type Config struct {
	Mode              string         `mapstructure:"mode"`
	ViewAttachBaseURL string         `mapstructure:"view_attach_base_url"`
	UploadDirPath     string         `mapstructure:"upload_dir_path"`
	Server            ServerConfig   `mapstructure:"server"`
	Logger            LoggerConfig   `mapstructure:"logger"`
	Database          DatabaseConfig `mapstructure:"database"`
	JWT               JWTConfig      `mapstructure:"jwt"`
	Auth              AuthConfig     `mapstructure:"auth"`
	Postal            PostalConfig   `mapstructure:"postal"`
	Attach            AttachConfig   `mapstructure:"attach"`
	TG                TelegramConfig `mapstructure:"tg"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port int `mapstructure:"port"`
}

// LoggerConfig 日志配置
type LoggerConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	MySQL MySQLConfig `mapstructure:"mysql"`
}

// MySQLConfig MySQL数据库配置
type MySQLConfig struct {
	Addr     string `mapstructure:"addr"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"db_name"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret           string `mapstructure:"secret"`
	AccessExpiresIn  string `mapstructure:"access_expires_in"`
	RefreshExpiresIn string `mapstructure:"refresh_expires_in"`
}

// AuthConfig 认证相关配置
// 包括访问白名单（支持 "METHOD:/path" 或仅路径形式）。
type AuthConfig struct {
	Whitelist   []string          `mapstructure:"whitelist"`
	DefaultUser DefaultUserConfig `mapstructure:"default_user"`
}

// Postal配置
type PostalConfig struct {
	SmtpServer string `mapstructure:"smtp_server"`
	SmtpPort   string `mapstructure:"smtp_port"`
	FromEmail  string `mapstructure:"from_email"`
	FromPass   string `mapstructure:"from_pass"`
	FromName   string `mapstructure:"from_name"`
}

type AttachConfig struct {
	ViewAttachBaseURL    string `mapstructure:"view_attach_base_url"`
	UploadDirPath        string `mapstructure:"upload_dir_path"`
	ViewLargeFileBaseURL string `mapstructure:"view_large_file_base_url"`
	UploadLargeFilePath  string `mapstructure:"upload_large_file_path"`
	ChunkDirPath         string `mapstructure:"chunk_dir_path"`
	ChunkDirSalt         string `mapstructure:"chunk_dir_salt"`
}

// TelegramConfig Telegram 机器人配置。
type TelegramConfig struct {
	BotToken string `mapstructure:"bot_token"`
	ChatID   string `mapstructure:"chat_id"`
}

type ThirdPartyConfig struct {
	FileAPI FileAPIConfig `mapstructure:"file_api"`
}

type FileAPIConfig struct {
	BaseURL        string `mapstructure:"base_url"`
	Secret         string `mapstructure:"secret"`
	TimeoutSeconds int    `mapstructure:"timeout_seconds"`
}

type DefaultUserConfig struct {
	UserID   int64  `mapstructure:"userid"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Nickname string `mapstructure:"nickname"`
}

// GlobalConfig 全局配置实例
var GlobalConfig *Config

// Init 初始化配置
func Init() error {
	// 获取环境变量，默认为 dev
	env := os.Getenv("BLUESPOT_ENV")
	if env == "" {
		env = "dev"
	}

	// 验证环境变量值，只允许 dev、test、prod
	validEnvs := map[string]bool{
		"dev":  true,
		"test": true,
		"prod": true,
	}
	if !validEnvs[env] {
		return fmt.Errorf("无效的环境变量 BLUESPOT_ENV=%s，只允许: dev, test, prod", env)
	}

	// 获取当前工作目录
	workDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取工作目录失败: %w", err)
	}

	loader := viper.New()
	addConfigPaths(loader, workDir)
	bindEnvVars(loader)

	if err := readRequiredConfig(loader, "app"); err != nil {
		return fmt.Errorf("读取配置文件 app.yml 失败: %w", err)
	}

	if err := mergeRequiredConfig(loader, env); err != nil {
		log.Printf("未找到配置文件 %s.yml，跳过: %v", env, err)
	}

	if merged, err := mergeOptionalConfig(loader, env+".local"); err != nil {
		return fmt.Errorf("读取配置文件 %s.local.yml 失败: %w", env, err)
	} else if merged {
		log.Printf("已合并本地配置文件: %s.local.yml", env)
	}

	// 将配置解析到结构体
	GlobalConfig = &Config{}
	if err := loader.Unmarshal(GlobalConfig); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	log.Printf("当前环境: %s", env)
	log.Printf("配置文件加载成功: app.yml -> %s.yml", env)
	return nil
}

func addConfigPaths(loader *viper.Viper, workDir string) {
	loader.AddConfigPath(filepath.Join(workDir, "internal/config"))
	loader.AddConfigPath("./internal/config")
	loader.AddConfigPath(".")
}

func bindEnvVars(loader *viper.Viper) {
	loader.SetEnvPrefix("BLUESPOT")
	loader.AutomaticEnv()

	loader.BindEnv("mode", "BLUESPOT_MODE")
	loader.BindEnv("view_attach_base_url", "BLUESPOT_VIEW_ATTACH_BASE_URL")
	loader.BindEnv("upload_dir_path", "BLUESPOT_UPLOAD_DIR_PATH")
	loader.BindEnv("attach.view_attach_base_url", "BLUESPOT_ATTACH_VIEW_ATTACH_BASE_URL")
	loader.BindEnv("attach.upload_dir_path", "BLUESPOT_ATTACH_UPLOAD_DIR_PATH")
	loader.BindEnv("attach.view_large_file_base_url", "BLUESPOT_ATTACH_VIEW_LARGE_FILE_BASE_PATH")
	loader.BindEnv("attach.upload_large_file_path", "BLUESPOT_ATTACH_UPLOAD_LARGE_FILE_PATH")
	loader.BindEnv("attach.chunk_dir_path", "BLUESPOT_ATTACH_CHUNK_DIR_PATH")
	loader.BindEnv("attach.chunk_dir_salt", "BLUESPOT_ATTACH_CHUNK_DIR_SALT")

	loader.BindEnv("server.port", "BLUESPOT_SERVER_PORT")

	loader.BindEnv("logger.level", "BLUESPOT_LOGGER_LEVEL")
	loader.BindEnv("logger.filename", "BLUESPOT_LOGGER_FILENAME")
	loader.BindEnv("logger.max_size", "BLUESPOT_LOGGER_MAX_SIZE")
	loader.BindEnv("logger.max_age", "BLUESPOT_LOGGER_MAX_AGE")
	loader.BindEnv("logger.max_backups", "BLUESPOT_LOGGER_MAX_BACKUPS")

	loader.BindEnv("database.mysql.addr", "BLUESPOT_DATABASE_MYSQL_ADDR")
	loader.BindEnv("database.mysql.user", "BLUESPOT_DATABASE_MYSQL_USER")
	loader.BindEnv("database.mysql.password", "BLUESPOT_DATABASE_MYSQL_PASSWORD")
	loader.BindEnv("database.mysql.db_name", "BLUESPOT_DATABASE_MYSQL_DB_NAME")

	viper.BindEnv("jwt.secret", "BLUESPOT_JWT_SECRET")
	viper.BindEnv("jwt.access_expires_in", "BLUESPOT_JWT_ACCESS_EXPIRES_IN")
	viper.BindEnv("jwt.refresh_expires_in", "BLUESPOT_JWT_REFRESH_EXPIRES_IN")

	loader.BindEnv("auth.default_user.userid", "BLUESPOT_AUTH_DEFAULT_USER_ID")
	loader.BindEnv("auth.default_user.username", "BLUESPOT_AUTH_DEFAULT_USER_USERNAME")
	loader.BindEnv("auth.default_user.password", "BLUESPOT_AUTH_DEFAULT_USER_PASSWORD")
	loader.BindEnv("auth.default_user.nickname", "BLUESPOT_AUTH_DEFAULT_USER_NICKNAME")

	loader.BindEnv("tg.bot_token", "BLUESPOT_TG_BOT_TOKEN")
	loader.BindEnv("tg.chat_id", "BLUESPOT_TG_CHAT_ID")

}

func readRequiredConfig(loader *viper.Viper, name string) error {
	loader.SetConfigName(name)
	loader.SetConfigType("yml")
	return loader.ReadInConfig()
}

func mergeRequiredConfig(loader *viper.Viper, name string) error {
	loader.SetConfigName(name)
	loader.SetConfigType("yml")
	return loader.MergeInConfig()
}

func mergeOptionalConfig(loader *viper.Viper, name string) (bool, error) {
	loader.SetConfigName(name)
	loader.SetConfigType("yml")
	if err := loader.MergeInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetConfig 获取全局配置
func GetConfig() *Config {
	return GlobalConfig
}

// GetAttachUploadDirPath 获取上传文件存放目录。
func GetAttachUploadDirPath() string {
	if GlobalConfig == nil {
		return ""
	}
	if strings.TrimSpace(GlobalConfig.Attach.UploadDirPath) != "" {
		return GlobalConfig.Attach.UploadDirPath
	}
	return GlobalConfig.UploadDirPath
}

// GetAttachViewBaseURL 获取上传文件访问地址前缀。
func GetAttachViewBaseURL() string {
	if GlobalConfig == nil {
		return ""
	}
	if strings.TrimSpace(GlobalConfig.Attach.ViewAttachBaseURL) != "" {
		return GlobalConfig.Attach.ViewAttachBaseURL
	}
	return GlobalConfig.ViewAttachBaseURL
}

// GetAttachLargeFileUploadPath 获取大文件上传存放目录。
func GetAttachLargeFileUploadPath() string {
	if GlobalConfig == nil {
		return ""
	}
	return GlobalConfig.Attach.UploadLargeFilePath
}

// GetAttachLargeFileViewBaseURL 获取大文件访问地址前缀。
func GetAttachLargeFileViewBaseURL() string {
	if GlobalConfig == nil {
		return ""
	}
	return GlobalConfig.Attach.ViewLargeFileBaseURL
}

// GetAttachChunkDirPath 获取分片文件存放目录。
func GetAttachChunkDirPath() string {
	if GlobalConfig == nil {
		return ""
	}
	return GlobalConfig.Attach.ChunkDirPath
}

// GetAttachChunkDirSalt 获取分片上传 ID 计算盐值。
func GetAttachChunkDirSalt() string {
	if GlobalConfig == nil {
		return ""
	}
	return GlobalConfig.Attach.ChunkDirSalt
}

// GetServerAddr 获取服务器地址
func GetServerAddr() string {
	if GlobalConfig == nil {
		return ":6306" // 默认端口
	}
	return fmt.Sprintf(":%d", GlobalConfig.Server.Port)
}

// GetMySQLDSN 获取MySQL连接字符串
func GetMySQLDSN() string {
	if GlobalConfig == nil {
		return ""
	}
	return BuildMySQLDSN(GlobalConfig)
}

// BuildMySQLDSN 从指定配置生成 MySQL 连接字符串
func BuildMySQLDSN(cfg *Config) string {
	if cfg == nil {
		return ""
	}

	mysql := cfg.Database.MySQL
	if strings.TrimSpace(mysql.Addr) == "" ||
		strings.TrimSpace(mysql.User) == "" ||
		strings.TrimSpace(mysql.DBName) == "" {
		return ""
	}
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysql.User, mysql.Password, mysql.Addr, mysql.DBName)
}

// GetEnv 获取当前环境（dev/test/prod）
func GetEnv() string {
	env := os.Getenv("BLUESPOT_ENV")
	if env == "" {
		return "dev"
	}
	return env
}

// IsProduction 判断是否为生产环境
func IsProduction() bool {
	return GetEnv() == "prod"
}

// IsDevelopment 判断是否为开发环境
func IsDevelopment() bool {
	return GetEnv() == "dev"
}

// IsTest 判断是否为测试环境
func IsTest() bool {
	return GetEnv() == "test"
}
