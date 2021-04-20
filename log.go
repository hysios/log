package log

import (
	"os"
	"strconv"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger *Logger
	sugar  *Sugar
	atom   = zap.NewAtomicLevel()
)

// Logger 日志
type Logger struct {
	*zap.Logger
}

// Sugar 日志糖
type Sugar struct {
	*zap.SugaredLogger
}

type Level int

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel = Level(zapcore.DebugLevel)
	// InfoLevel is the default logging priority.
	InfoLevel = Level(zapcore.InfoLevel)
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel = Level(zapcore.WarnLevel)
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel = Level(zapcore.ErrorLevel)
	// DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	DPanicLevel = Level(zapcore.DPanicLevel)
	// PanicLevel logs a message, then panics.
	PanicLevel = Level(zapcore.PanicLevel)
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel = Level(zapcore.FatalLevel)

	_minLevel = DebugLevel
	_maxLevel = FatalLevel
)

// SetLevel 设置输出级别
func SetLevel(l Level) {
	atom.SetLevel(zapcore.Level(l))
}

// NewLogger 创建日志对象
func NewLogger() *Sugar {

	cfg := &zap.Config{
		Level:            atom,
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	// cfg.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder

	zlog, _ := cfg.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	logger = &Logger{Logger: zlog}
	defer logger.Sync()
	return &Sugar{SugaredLogger: logger.Sugar()}
}

func ConsoleCore(cfg ConfigConsole) zapcore.Core {
	lowPriority := zap.LevelEnablerFunc(func(lv zapcore.Level) bool {
		return lv >= zapcore.Level(cfg.Priority)
	})

	// 使用 stdout 格式日志
	zcfg := zap.NewDevelopmentEncoderConfig()
	zcfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	stdEnc := zapcore.NewConsoleEncoder(zcfg)
	core := zapcore.NewCore(stdEnc, zapcore.Lock(os.Stdout), lowPriority)
	return core
}

func FileCore(cfg ConfigFile) zapcore.Core {
	lowPriority := zap.LevelEnablerFunc(func(lv zapcore.Level) bool {
		return lv >= zapcore.Level(cfg.Priority)
	})

	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   cfg.Path,
		MaxSize:    cfg.SplitSize,
		MaxBackups: cfg.MaxBackups,
		LocalTime:  cfg.LocalTime,
		Compress:   cfg.Compress,
	})

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	return zapcore.NewCore(
		zapcore.NewJSONEncoder(config),
		w,
		lowPriority,
	)
}

func RedisCore(cfg ConfigRedis) zapcore.Core {
	lowPriority := zap.LevelEnablerFunc(func(lv zapcore.Level) bool {
		return lv >= zapcore.Level(cfg.Priority)
	})

	cli := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	writer := NewRedisWriter(cfg.KeyName, cli)

	// 使用 JSON 格式日志
	jsonEnc := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	syncer := zapcore.AddSync(writer)
	return zapcore.NewCore(jsonEnc, syncer, lowPriority)
}

func init() {
	l, err := strconv.Atoi(os.Getenv("DEBUG_LEVEL"))
	if err == nil {
		SetLevel(Level(l))
	}
	sugar = NewLogger()
}

// SetLogger 设置新 Logger 句柄
func SetLogger(logger *Logger) {
	defer logger.Sync()
	sugar = &Sugar{SugaredLogger: logger.Sugar()}
}

func DPanic(args ...interface{}) {
	sugar.DPanic(args...)
}

func DPanicf(template string, args ...interface{}) {
	sugar.DPanicf(template, args)
}
func DPanicw(msg string, keysAndValues ...interface{}) {
	sugar.DPanicw(msg, keysAndValues...)
}
func Debug(args ...interface{}) {
	sugar.Debug(args...)
}
func Debugf(template string, args ...interface{}) {
	sugar.Debugf(template, args...)
}
func Debugw(msg string, keysAndValues ...interface{}) {
	sugar.Debugw(msg, keysAndValues...)
}
func Desugar() *Logger {
	return &Logger{Logger: sugar.Desugar()}
}
func Error(args ...interface{}) {
	sugar.Error(args...)
}
func Errorf(template string, args ...interface{}) {
	sugar.Errorf(template, args...)
}
func Errorw(msg string, keysAndValues ...interface{}) {
	sugar.Errorw(msg, keysAndValues...)
}
func Fatal(args ...interface{}) {
	sugar.Fatal(args...)
}
func Fatalf(template string, args ...interface{}) {
	sugar.Fatalf(template, args...)
}
func Fatalw(msg string, keysAndValues ...interface{}) {
	sugar.Fatalw(msg, keysAndValues...)
}
func Info(args ...interface{}) {
	sugar.Info(args...)
}
func Infof(template string, args ...interface{}) {
	sugar.Infof(template, args...)
}
func Infow(msg string, keysAndValues ...interface{}) {
	sugar.Infow(msg, keysAndValues...)
}
func Named(name string) *Sugar {
	return &Sugar{SugaredLogger: sugar.Named(name)}
}
func Panic(args ...interface{}) {
	sugar.Panic(args...)
}
func Panicf(template string, args ...interface{}) {
	sugar.Panicf(template, args...)
}
func Panicw(msg string, keysAndValues ...interface{}) {
	sugar.Panicw(msg, keysAndValues...)
}
func Sync() error {
	return sugar.Sync()
}
func Warn(args ...interface{}) {
	sugar.Warn(args...)
}
func Warnf(template string, args ...interface{}) {
	sugar.Warnf(template, args)
}
func Warnw(msg string, keysAndValues ...interface{}) {
	sugar.Warnw(msg, keysAndValues...)
}

func With(args ...interface{}) *Sugar {
	return &Sugar{SugaredLogger: sugar.With(args...)}
}

func LoggerFile() *Sugar {
	return sugar
}
