package telemetry

import (
	"context"
	"log/slog"
	"reflect"
)

func errorFormattingMiddleware(
	ctx context.Context,
	record slog.Record,
	next func(context.Context, slog.Record) error,
) error {
	attrs := []slog.Attr{}

	record.Attrs(func(attr slog.Attr) bool {
		key := attr.Key
		value := attr.Value
		kind := attr.Value.Kind()

		if (key == "error" || key == "err") && kind == slog.KindAny {
			if err, ok := value.Any().(error); ok {
				errType := reflect.TypeOf(err).String()
				msg := err.Error()

				attrs = append(
					attrs,
					slog.Group("error",
						slog.String("type", errType),
						slog.String("message", msg),
					),
				)
			} else {
				attrs = append(attrs, attr)
			}
		} else {
			attrs = append(attrs, attr)
		}

		return true
	})

	record = slog.NewRecord(record.Time, record.Level, record.Message, record.PC)
	record.AddAttrs(attrs...)

	return next(ctx, record)
}
