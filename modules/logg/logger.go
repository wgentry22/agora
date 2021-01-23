package logg

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/wgentry22/agora/types/config"
	"io"
	"io/ioutil"
	"os"
	"time"
)

var (
	rootLogger Logger
)

type Logger interface {
	Level() string
	WithWriter(writer io.Writer) Logger
	WithContext(ctx context.Context) Logger
	WithField(key string, value interface{}) Logger
	WithError(err error) Logger
	WithLevel(level string) Logger
	Trace(args ...interface{})
	Tracef(msg string, args ...interface{})
	Debug(args ...interface{})
	Debugf(msg string, args ...interface{})
	Info(args ...interface{})
	Infof(msg string, args ...interface{})
	Warning(args ...interface{})
	Warningf(msg string, args ...interface{})
	Warn(args ...interface{})
	Warnf(msg string, args ...interface{})
	Error(args ...interface{})
	Errorf(msg string, args ...interface{})
	Panic(args ...interface{})
	Panicf(msg string, args ...interface{})
}

func Writers(conf config.Logging) io.Writer {
	if len(conf.OutputPaths) == 0 {
		return os.Stdout
	}

	writers := make([]io.Writer, 0)

	for _, path := range conf.OutputPaths {
		switch {
		case path == "stdout":
			writers = append(writers, os.Stdout)
		case path == "stderr":
			writers = append(writers, os.Stderr)
		case path == "/dev/null":
			writers = append(writers, ioutil.Discard)
		default:
			file, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0600)
			if err != nil {
				rootLogger.
					WithField("path", path).
					WithError(err).
					Warning("Skipping entry")
			} else {
				writers = append(writers, file)
			}
		}
	}

	if len(writers) == 0 {
		writers = append(writers, os.Stdout)
	}

	return io.MultiWriter(writers...)
}

func NewLogrusLogger(conf config.Logging) Logger {
	return &logrusAdapter{
		conf:   conf,
		logger: buildLogger(conf),
	}
}

func buildLogger(conf config.Logging) *logrus.Entry {
	return buildLoggerWithWriter(conf, Writers(conf))
}

func buildLoggerWithWriter(conf config.Logging, writer io.Writer) *logrus.Entry {
	logger := logrus.New()

	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat:  time.UnixDate,
		DisableTimestamp: false,
		PrettyPrint:      true,
	})

	logger.Out = writer

	level, err := logrus.ParseLevel(conf.Level)
	if err != nil {
		logger.
			WithError(err).
			WithTime(time.Now()).
			WithField("level", conf.Level).
			Warningln("cannot use: rolling back to `debug`")

		level = logrus.DebugLevel
	}

	logger.Level = level

	return logger.WithFields(conf.Fields)
}

type logrusAdapter struct {
	conf   config.Logging
	logger *logrus.Entry
}

func (l logrusAdapter) Level() string {
	return l.conf.Level
}

func (l logrusAdapter) WithWriter(writer io.Writer) Logger {
	return &logrusAdapter{
		conf:   l.conf,
		logger: buildLoggerWithWriter(l.conf, writer),
	}
}

func (l logrusAdapter) WithContext(ctx context.Context) Logger {
	return &logrusAdapter{
		conf:   l.conf,
		logger: l.logger.WithContext(ctx),
	}
}

func (l logrusAdapter) WithField(key string, value interface{}) Logger {
	return &logrusAdapter{
		conf:   l.conf,
		logger: l.logger.WithField(key, value),
	}
}

func (l logrusAdapter) WithError(err error) Logger {
	return &logrusAdapter{
		conf:   l.conf,
		logger: l.logger.WithError(err),
	}
}

func (l logrusAdapter) WithLevel(level string) Logger {
	logrusLevel, err := logrus.ParseLevel(level)
	if err != nil {
		l.logger.
			WithField("level", level).
			Warningln("Failed to parse")

		logrusLevel = logrus.DebugLevel
	}

	conf := config.Logging{
		Level:       logrusLevel.String(),
		OutputPaths: l.conf.OutputPaths,
		Fields:      l.conf.Fields,
	}

	return &logrusAdapter{
		conf:   conf,
		logger: buildLogger(conf),
	}
}

func (l logrusAdapter) Trace(args ...interface{}) {
	l.logger.Traceln(args...)
}

func (l logrusAdapter) Tracef(msg string, args ...interface{}) {
	l.logger.Tracef(msg, args...)
}

func (l logrusAdapter) Debug(args ...interface{}) {
	l.logger.Debugln(args...)
}

func (l logrusAdapter) Debugf(msg string, args ...interface{}) {
	l.logger.Debugf(msg, args...)
}

func (l logrusAdapter) Info(args ...interface{}) {
	l.logger.Infoln(args...)
}

func (l logrusAdapter) Infof(msg string, args ...interface{}) {
	l.logger.Infof(msg, args...)
}

func (l logrusAdapter) Warning(args ...interface{}) {
	l.logger.Warningln(args...)
}

func (l logrusAdapter) Warningf(msg string, args ...interface{}) {
	l.logger.Warningf(msg, args...)
}

func (l logrusAdapter) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l logrusAdapter) Warnf(msg string, args ...interface{}) {
	l.logger.Warnf(msg, args...)
}

func (l logrusAdapter) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l logrusAdapter) Errorf(msg string, args ...interface{}) {
	l.logger.Errorf(msg, args...)
}

func (l logrusAdapter) Panic(args ...interface{}) {
	l.logger.Panic(args...)
}

func (l logrusAdapter) Panicf(msg string, args ...interface{}) {
	l.logger.Panicf(msg, args...)
}
