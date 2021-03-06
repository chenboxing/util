package log

import (
	"io"
	"sync"

	"github.com/chenboxing/util/log/formatter/logstash"

	rus "github.com/Sirupsen/logrus"
)

var std = newLogger()

type logger struct {
	formatter string
	mode      Mode
	release   string
	l         *rus.Logger
	daemon    daemon
	mu        sync.Mutex
}

func newLogger() *logger {
	return &logger{l: rus.New()}
}

func (l *logger) print(fields Fields, msg interface{}, f logFunc) {
	l.mu.Lock()
	defer l.mu.Unlock()
	f(l.l.WithFields(rus.Fields(fields)), msg)
}

func (l *logger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.l.Level = rus.Level(level)
}

func (l *logger) GetLevel() Level {
	l.mu.Lock()
	defer l.mu.Unlock()
	return Level(l.l.Level)
}

func (l *logger) SetMode(mode Mode) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.mode = mode
}

func (l *logger) GetMode() Mode {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.mode
}

func (l *logger) SetRelease(release string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.release = release
}

func (l *logger) GetRelease() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.release
}

func (l *logger) SetPort(port int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.daemon.Port = port
}

func (l *logger) GetPort() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.daemon.Port
}

func (l *logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.l.Out = w
}

func (l *logger) SetFormatter(formatter string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	switch formatter {
	case "logstash":
		l.l.Formatter = new(logstash.LogstashFormatter)
		l.formatter = formatter
	case "json":
		l.l.Formatter = new(rus.JSONFormatter)
		l.formatter = formatter
	case "text":
		l.l.Formatter = new(rus.TextFormatter)
		l.formatter = formatter
	default:
		l.l.Formatter = new(rus.TextFormatter)
		l.formatter = "text"
	}
}

func (l *logger) GetFormatter() string {
	return l.formatter
}

func (l *logger) startDaemon() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.daemon.start()
}

func (l *logger) AddHook(hook rus.Hook) {
	l.mu.Lock()
	l.mu.Unlock()

	l.l.Hooks.Add(hook)
}
