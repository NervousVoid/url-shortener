package tools

import "golang.org/x/exp/slog"

func LogAttr(key, message string) slog.Attr {
	return slog.Attr{
		Key:   key,
		Value: slog.StringValue(message),
	}
}
