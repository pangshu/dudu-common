package common

import (
	"bytes"
	"os"
	"path"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Option func(*LogConfig)
type LogConfig struct {
	Path		string			`yaml:"path"`
	Name 		string			`yaml:"name"`
	Level 		string			`yaml:"level"`
	MaxSize		int				`yaml:"max_size"`
	MaxAge		int				`yaml:"max_age"`

	Stacktrace	string			`yaml:"stacktrace"`
	Stdout		bool 			`yaml:"stdout"`
}

var instance *zap.Logger

// Instance 唯一实例
func (*DuduLog) Instance() *zap.Logger {
	return instance
}

// Init 初始化,生成的日志文件夹名字
func (*DuduLog) Init(conf LogConfig) *zap.Logger {
	instance = Log.NewLogger(func(o *LogConfig) {
		o.Path = conf.Path
		o.Name = conf.Name
		o.Level = conf.Level
		o.MaxSize = conf.MaxSize
		o.MaxAge = conf.MaxAge
		o.Stacktrace = conf.Stacktrace
		o.Stdout = conf.Stdout
	})
	return instance
}

// NewLogger 新建日志
func (*DuduLog) NewLogger(opts ...Option) *zap.Logger {
	o := &LogConfig{
		Path:"./log/",
		Name:"output",
		Level:"debug",
		MaxSize:100,
		MaxAge:7,
		Stacktrace: "error",
		Stdout: true,
	}
	for _,opt := range opts {
		opt(o)
	}
	writers := []zapcore.WriteSyncer{newRollingFile(o.Path,o.Name,o.MaxSize,o.MaxAge)}
	if o.Stdout == true {
		writers = append(writers, os.Stdout)
	}
	logger := newZapLogger(getLevel(o.Level),getLevel(o.Stacktrace), zapcore.NewMultiWriteSyncer(writers...))
	zap.RedirectStdLog(logger)

	return logger
}

func getLevel(level string) zapcore.Level{
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.ErrorLevel
	}
}

//创建分割日志的writer
func newRollingFile(logPath,logName string,maxSize,maxAge int) zapcore.WriteSyncer {
	if err := os.MkdirAll(logPath, 0766); err != nil {
		panic(err)
		return nil
	}

	return newLumberjackWriteSyncer(&lumberjack.Logger{
		Filename:  path.Join(logPath,logName + ".log"),
		MaxSize:   maxSize, //megabytes
		MaxAge:    maxAge,   //days
		LocalTime: true,
		Compress:  false,
	})
}

func newZapLogger(level,stacktrace zapcore.Level, output zapcore.WriteSyncer) (*zap.Logger) {
	encCfg := zapcore.EncoderConfig{
		TimeKey:        "@timestamp",
		LevelKey:       "level",
		NameKey:        "app",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeDuration: zapcore.NanosDurationEncoder,
		//EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		//	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
		//},
		EncodeTime: zapcore.ISO8601TimeEncoder,
	}

	var encoder zapcore.Encoder
	dyn := zap.NewAtomicLevel()
	//encCfg.EncodeLevel = zapcore.LowercaseLevelEncoder
	//encoder = zapcore.NewJSONEncoder(encCfg) // zapcore.NewConsoleEncoder(encCfg)
	dyn.SetLevel(level)
	encCfg.EncodeLevel = zapcore.LowercaseLevelEncoder
	encoder = zapcore.NewJSONEncoder(encCfg)

	return zap.New(zapcore.NewCore(encoder, output, dyn), zap.AddCaller(),zap.AddStacktrace(stacktrace),zap.AddCallerSkip(2))
}

type lumberjackWriteSyncer struct {
	*lumberjack.Logger
	buf       *bytes.Buffer
	logChan   chan []byte
	closeChan chan interface{}
	maxSize   int
}

func newLumberjackWriteSyncer(l *lumberjack.Logger) *lumberjackWriteSyncer {
	ws := &lumberjackWriteSyncer{
		Logger:    l,
		buf:       bytes.NewBuffer([]byte{}),
		logChan:   make(chan []byte, 5000),
		closeChan: make(chan interface{}),
		maxSize:   1024,
	}
	go ws.run()
	return ws
}

func (l *lumberjackWriteSyncer) run() {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			if l.buf.Len() > 0 {
				l.sync()
			}
		case bs := <-l.logChan:
			_, err := l.buf.Write(bs)
			if err != nil {
				continue
			}
			if l.buf.Len() > l.maxSize {
				l.sync()
			}
		case <-l.closeChan:
			l.sync()
			return
		}
	}
}

func (l *lumberjackWriteSyncer) Stop() {
	close(l.closeChan)
}

func (l *lumberjackWriteSyncer) Write(bs []byte) (int, error) {
	b := make([]byte, len(bs))
	for i, c := range bs {
		b[i] = c
	}
	l.logChan <- b
	return 0, nil
}

func (l *lumberjackWriteSyncer) Sync() error {
	return nil
}

func (l *lumberjackWriteSyncer) sync() error {
	defer l.buf.Reset()
	_, err := l.Logger.Write(l.buf.Bytes())
	if err != nil {
		return err
	}
	return nil
}
