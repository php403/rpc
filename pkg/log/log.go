package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level int8

const (
	// LevelDebug is logger debug level.
	LevelDebug Level = iota - 1
	// LevelInfo is logger info level.
	LevelInfo
	// LevelWarn is logger warn level.
	LevelWarn
	// LevelError is logger error level.
	LevelError
	// LevelFatal is logger fatal level
	LevelFatal
)


type Logger interface {
	Log(level Level, keyvals ...interface{}) error
	Debug(msg string, keyvals ...interface{})
	Info(msg string, keyvals ...interface{})
	Warn(msg string, keyvals ...interface{})
	Error(msg string, keyvals ...interface{})
	Debugf(msg string, a ...interface{})
	Infof(msg string, a ...interface{})
	Warnf(msg string, a ...interface{})
	Errorf(msg string, a ...interface{})
}

type Config struct {
	test string
}

type ZapLogger struct {
	log  *zap.Logger
}



func (l *ZapLogger) Log(level Level, keyvals ...interface{}) error {
	if len(keyvals) == 0 || len(keyvals)%2 != 0 {
		l.log.Warn(fmt.Sprint("Keyvalues must appear in pairs: ", keyvals))
		return nil
	}

	// Zap.Field is used when keyvals pairs appear
	var data []zap.Field
	for i := 0; i < len(keyvals); i += 2 {
		data = append(data, zap.Any(fmt.Sprint(keyvals[i]), fmt.Sprint(keyvals[i+1])))
	}
	switch level {
	case LevelDebug:
		l.log.Debug("", data...)
	case LevelInfo:
		l.log.Info("", data...)
	case LevelWarn:
		l.log.Warn("", data...)
	case LevelError:
		l.log.Error("", data...)
	}
	return nil
}

func(l *ZapLogger) Logs(level Level,msg string)  {
	switch level {
	case LevelDebug:
		l.log.Debug(msg)
	case LevelInfo:
		l.log.Info(msg)
	case LevelWarn:
		l.log.Warn(msg)
	case LevelError:
		l.log.Error(msg)
	}
}

func(l *ZapLogger) Debug(msg string,keyvals ...interface{}) ()  {
	l.log.Debug(msg)
}
func(l *ZapLogger) Info(msg string,keyvals ...interface{}) ()  {
	l.log.Debug(msg)
}
func(l *ZapLogger) Warn(msg string,keyvals ...interface{}) ()  {
	l.log.Debug(msg)
}
func(l *ZapLogger) Error(msg string,keyvals ...interface{}) ()  {
	l.log.Debug(msg)
}

func(l *ZapLogger) Debugf(msg string,a ...interface{}) ()  {
	l.log.Debug(fmt.Sprintf(msg,a))
}
func(l *ZapLogger) Infof(msg string,a ...interface{}) ()  {
	l.log.Debug(fmt.Sprintf(msg,a))
}
func(l *ZapLogger) Warnf(msg string,a ...interface{}) ()  {
	l.log.Debug(fmt.Sprintf(msg,a))
}
func(l *ZapLogger) Errorf(msg string,a ...interface{}) ()  {
	l.log.Debug(fmt.Sprintf(msg,a))
}

func NewZapLogger(encoder zapcore.EncoderConfig, level zap.AtomicLevel) *ZapLogger {
	config := zap.Config{
		Level:            level,                                                // 日志级别
		Development:      true,                                                // 开发模式，堆栈跟踪
		Encoding:         "json",                                              // 输出格式 console 或 json
		EncoderConfig:    encoder,                                       // 编码器配置
		InitialFields:    map[string]interface{}{"serviceName": "im"}, // 初始化字段，如：添加一个服务器名称
		OutputPaths:      []string{"stdout"},         // 输出到指定文件 stdout（标准输出，正常颜色） stderr（错误输出，红色）
		ErrorOutputPaths: []string{"stderr"},
	}
	// 构建日志
	logger, err := config.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	if err != nil {
		panic(fmt.Sprintf("log 初始化失败: %v", err))
	}
	return &ZapLogger{log: logger}
}

func InitConfig() (encoderConfig zapcore.EncoderConfig,atomicLevel zap.AtomicLevel) {
	encoderConfig = zapcore.EncoderConfig{
		LevelKey:       "level",
		MessageKey:     "msg",
		CallerKey:      "caller",
		TimeKey:        "time",
		NameKey:        "logger",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
	}
	// 设置日志级别
	atomicLevel = zap.NewAtomicLevelAt(zap.DebugLevel)
	return
}