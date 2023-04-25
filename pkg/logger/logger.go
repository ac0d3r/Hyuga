/*
 * Copyright 2022 by Mel2oo <https://github.com/saferun/owl>
 *
 * Licensed under the GNU General Public License version 3 (GPLv3)
 *
 * If you distribute GPL-licensed software the license requires
 * that you also distribute the complete, corresponding source
 * code (as defined by GPL) to that GPL-licensed software.
 *
 * You should have received a copy of the GNU General Public License
 * with this program. If not, see <https://www.gnu.org/licenses/>
 */

package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
	AppName    string `mapstructure:"app"`
	Level      string `mapstructure:"level"`
	MaxSize    int    `mapstructure:"max-size"`
	MaxAge     int    `mapstructure:"max-age"`
	MaxBackups int    `mapstructure:"max-backup"`
	Output     string `mapstructure:"output"`
	Stdout     bool   `mapstructure:"stdout"`
	ShowFile   bool   `mapstructure:"show-file"`
}

func Init(cfg *Config) error {
	// setup log formatter
	logrus.SetFormatter(&nested.Formatter{
		TimestampFormat: time.RFC3339,
		HideKeys:        true,
	})

	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		return err
	}
	logrus.SetLevel(level)

	// disable writing to stdout
	if !cfg.Stdout {
		logrus.SetOutput(io.Discard)
	}

	// initialize log rotate hook
	if cfg.AppName == "" {
		cfg.AppName = os.Args[0]
	}
	rhook, err := NewHook(cfg.AppName, cfg, level)
	if err != nil {
		return err
	}
	logrus.AddHook(rhook)

	return err
}

// File represents the rotate file hook.
type File struct {
	app       string
	level     logrus.Level
	file      bool
	formatter logrus.Formatter
	w         io.Writer
}

// NewHook builds a new rotate file hook.
func NewHook(app string, cfg *Config, level logrus.Level) (logrus.Hook, error) {
	name := os.Getenv("LOG_NAME")
	if len(name) == 0 {
		name = filepath.Base(os.Args[0])
	}

	return &File{
		app:   app,
		level: level,
		file:  cfg.ShowFile,
		formatter: &nested.Formatter{
			TimestampFormat: time.RFC3339,
			NoColors:        true,
			HideKeys:        true,
		},
		w: &lumberjack.Logger{
			Filename:   filepath.Join(cfg.Output, name+".log"),
			MaxSize:    cfg.MaxSize,
			MaxAge:     cfg.MaxAge,
			MaxBackups: cfg.MaxBackups,
		},
	}, nil
}

// Levels determines log levels that for which the logs are written.
func (hook *File) Levels() []logrus.Level {
	return logrus.AllLevels[:hook.level+1]
}

// Fire is called by logrus when it is about to write the log entry.
func (hook *File) Fire(entry *logrus.Entry) (err error) {
	entry.Data["app"] = hook.app
	if hook.file {
		entry.Data["file"] = FindCaller()
	}
	b, err := hook.formatter.Format(entry)
	if err != nil {
		return err
	}
	_, err = hook.w.Write(b)
	return err
}

func FindCaller() string {
	var (
		file string
		line int
	)

	for i := 0; i < 20; i++ {
		file, line = getCaller(i + 5)
		if !skipFile(file) {
			break
		}
	}

	return fmt.Sprintf("%s:%d", file, line)
}

var skipPrefixes = []string{"logrus/", "logrus@", "v4@", "logger/"}

func skipFile(file string) bool {
	for i := range skipPrefixes {
		if strings.HasPrefix(file, skipPrefixes[i]) {
			return true
		}
	}
	return false
}

func getCaller(skip int) (string, int) {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "", 0
	}

	n := 0
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			n++
			if n >= 2 {
				file = file[i+1:]
				break
			}
		}
	}

	return file, line
}
