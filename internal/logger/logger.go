package logger

import (
	"github.com/wb-go/wbf/zlog"
)

func Init() {
	zlog.InitConsole()
}

func Info(msg string, args ...any) {
	event := zlog.Logger.Info()

	for i := 0; i+1 < len(args); i += 2 {
		key, ok := args[i].(string)
		if !ok {
			continue
		}
		event = event.Interface(key, args[i+1])
	}

	event.Msg(msg)
}

func Error(msg string, args ...any) {
	event := zlog.Logger.Error()

	for i := 0; i+1 < len(args); i += 2 {
		key, ok := args[i].(string)
		if !ok {
			continue
		}
		event = event.Interface(key, args[i+1])
	}

	event.Msg(msg)
}

func Debug(msg string, args ...any) {
	event := zlog.Logger.Debug()

	for i := 0; i+1 < len(args); i += 2 {
		key, ok := args[i].(string)
		if !ok {
			continue
		}
		event = event.Interface(key, args[i+1])
	}

	event.Msg(msg)
}

func Warn(msg string, args ...any) {
	event := zlog.Logger.Warn()

	for i := 0; i+1 < len(args); i += 2 {
		key, ok := args[i].(string)
		if !ok {
			continue
		}
		event = event.Interface(key, args[i+1])
	}

	event.Msg(msg)
}