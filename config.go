package log

import (
	"fmt"

	"github.com/creasty/defaults"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ConfigMode 配置预设 preset
type ConfigMode string

const (
	CDevelopment ConfigMode = "dev"  // 开发模式
	CProduction  ConfigMode = "prod" // 生产模式
)

type Config struct {
	Mode ConfigMode `default:"dev"`
	Dev  ConfigEnv
	Prod ConfigEnv
}

type ConfigEnv struct {
	Console struct {
		ConfigConsole
		On bool `default:"true"`
	}
	File struct {
		ConfigFile
		On bool
	}
	Redis struct {
		ConfigRedis
		On bool
	}
	CallerSkip int `default:"1"`
}

type ConfigConsole struct {
	Priority int `default:"-1"`
}

type ConfigFile struct {
	Path       string `default:"/tmp/cspdls/logs"`
	SplitSize  int    `default:"1024"`
	LocalTime  bool   `default:"true"`
	MaxBackups int    `default:"5"`
	Compress   bool   `default:"true"`
	Priority   int    `default:"-1"`
}

type ConfigRedis struct {
	Addr     string `default:"127.0.0.1:6379"`
	Password string
	DB       int    `default:"0"`
	KeyName  string `default:"logs"`
	Priority int    `default:"-1"`
}

var DefaultConfig = Config{}

func SetupConfig(cfgm map[string]interface{}) (logger *Logger, err error) {
	var cfg = &ConfigEnv{}
	defaults.Set(cfg)

	if err = mapstructure.Decode(cfgm, cfg); err != nil {
		return nil, err
	}

	var cores = make([]zapcore.Core, 0)

	if cfg.Console.On {
		cores = append(cores, ConsoleCore(cfg.Console.ConfigConsole))
	}

	if cfg.File.On {
		cores = append(cores, FileCore(cfg.File.ConfigFile))
	}

	if cfg.Redis.On {
		cores = append(cores, RedisCore(cfg.Redis.ConfigRedis))
	}
	core := zapcore.NewTee(cores...)

	fmt.Printf("caller skip %d", cfg.CallerSkip)
	zlog := zap.New(core).WithOptions(zap.AddCaller(), zap.AddCallerSkip(cfg.CallerSkip))
	return &Logger{Logger: zlog}, nil
}

func ConfigLogger(cfgm map[string]interface{}) {
	_logger, err := SetupConfig(cfgm)
	if err != nil {
		return
	}

	logger = _logger

	sugar = &Sugar{SugaredLogger: logger.Sugar()}
}

func MultiCore(cfg *ConfigEnv, cores ...zapcore.Core) *Logger {
	core := zapcore.NewTee(cores...)
	var opts = []zap.Option{zap.AddCaller()}
	if cfg != nil {
		opts = append(opts, zap.AddCallerSkip(cfg.CallerSkip))
	}
	zlog := zap.New(core).WithOptions(opts...)

	return &Logger{Logger: zlog}
}
