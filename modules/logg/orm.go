package logg

import (
	"context"
	"fmt"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"time"
)

var logStr = "(%s) %s %s %s %s"

var levels = map[logger.LogLevel]string{
	logger.Silent: "trace",
	logger.Info:   "info",
	logger.Warn:   "warn",
	logger.Error:  "error",
}

func ForGorm(logger Logger) GormLogAdapter {
	return GormLogAdapter{logger}
}

type GormLogAdapter struct {
	logger Logger
}

func (g GormLogAdapter) LogMode(level logger.LogLevel) logger.Interface {
	return &GormLogAdapter{
		logger: g.logger.WithLevel(levels[level]),
	}
}

func (g GormLogAdapter) Info(ctx context.Context, s string, i ...interface{}) {
	g.logger.WithContext(ctx).Infof(s, i)
}

func (g GormLogAdapter) Warn(ctx context.Context, s string, i ...interface{}) {
	g.logger.WithContext(ctx).Warnf(s, i)
}

func (g GormLogAdapter) Error(ctx context.Context, s string, i ...interface{}) {
	g.logger.WithContext(ctx).Errorf(s, i)
}

func (g GormLogAdapter) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if g.logger.Level() != "trace" {
		elapsed := time.Since(begin)
		l := g.logger.WithContext(ctx)
		switch {
		case err != nil && l.Level() != "error":
			sql, rows := fc()
			if rows == -1 {
				l.Tracef(logStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
			} else {
				l.Tracef(logStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
			}
		case elapsed > 2500*time.Millisecond && l.Level() != "warn":
			sql, rows := fc()
			slowLog := fmt.Sprintf("SLOW SQL >= %v", 2500*time.Millisecond)
			if rows == -1 {
				l.Warnf(logStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
			} else {
				l.Warnf(logStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
			}
		case l.Level() == "info":
			sql, rows := fc()
			if rows == -1 {
				l.Infof(logStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
			} else {
				l.Infof(logStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
			}
		}
	}
}
