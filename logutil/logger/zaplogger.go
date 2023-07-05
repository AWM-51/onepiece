package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

const (
	logTmFmt = "2023-01-02 15:04:05"
)

var errorLogger *zap.SugaredLogger

var levelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

func getLoggerLevel(lvl string) zapcore.Level {
	if level, ok := levelMap[lvl]; ok {
		return level
	}
	return zapcore.InfoLevel
}

func init() {
	level := getLoggerLevel("debug")
	syncWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./logs/lesson.log",
		MaxSize:    100,
		MaxBackups: 30,
		MaxAge:     30,
		LocalTime:  true,
		Compress:   true,
	})

	encoder := zap.NewProductionEncoderConfig()
	encoder.EncodeTime = zapcore.ISO8601TimeEncoder

	ConsoleEncoder := GetConsoleEncoder()
	core := zapcore.NewTee(
		zapcore.NewCore(ConsoleEncoder, syncWriter, zap.NewAtomicLevelAt(level)),     // 写入文件
		zapcore.NewCore(ConsoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel), // 写入控制台
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	errorLogger = logger.Sugar()
}

func GetEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(
		zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller_line",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     "  ",
			EncodeLevel:    cEncodeLevel,
			EncodeTime:     cEncodeTime,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   cEncodeCaller,
		})
}

// GetConsoleEncoder 输出日志到控制台
func GetConsoleEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
}

// GetLevelEnabler 自定义的LevelEnabler
func GetLevelEnabler() zapcore.Level {
	return zapcore.InfoLevel
}

// cEncodeLevel 自定义日志级别显示
func cEncodeLevel(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + level.CapitalString() + "]")
}

// cEncodeTime 自定义时间格式显示
func cEncodeTime(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + t.Format(logTmFmt) + "]")
}

// cEncodeCaller 自定义行号显示
func cEncodeCaller(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + caller.TrimmedPath() + "]")
}
func Debug(args ...interface{}) {
	errorLogger.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	errorLogger.Debugf(template, args...)
}

func Info(args ...interface{}) {
	errorLogger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	errorLogger.Infof(template, args...)
}

func Warn(args ...interface{}) {
	errorLogger.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	errorLogger.Warnf(template, args...)
}

func Error(args ...interface{}) {
	errorLogger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	errorLogger.Errorf(template, args...)
}

func DPanic(args ...interface{}) {
	errorLogger.DPanic(args...)
}

func DPanicf(template string, args ...interface{}) {
	errorLogger.DPanicf(template, args...)
}

func Panic(args ...interface{}) {
	errorLogger.Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	errorLogger.Panicf(template, args...)
}

func Fatal(args ...interface{}) {
	errorLogger.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	errorLogger.Fatalf(template, args...)
}
