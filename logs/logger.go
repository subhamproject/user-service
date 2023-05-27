package logs

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/subhamproject/user-service/consts"
	"go.opentelemetry.io/otel/trace"
)

var Log *logrus.Logger

func init() {
	// Create a new instance of the logger. You can have any number of instances.
	Log = logrus.New()

	Log.Formatter = new(logrus.JSONFormatter)
	// Output to stdout instead of the default stderr
	Log.Level = logrus.TraceLevel
	Log.Out = os.Stdout
}

func AddFieldsToLogger(span trace.Span) logrus.Fields {
	if span == nil {
		return logrus.Fields{}
	}

	return logrus.Fields{
		"dd.trace_id": convertTraceID(span.SpanContext().TraceID().String()),
		"dd.span_id":  convertTraceID(span.SpanContext().SpanID().String()),
		"dd.service":  consts.ServiceName,
	}
}

// var standardFields = logrus.Fields{
// 	"dd.trace_id": convertTraceID(span.SpanContext().TraceID().String()),
// 	"dd.span_id":  convertTraceID(span.SpanContext().SpanID().String()),
// 	"dd.service":  "serviceName",
// 	"dd.env":      "serviceEnv",
// 	"dd.version":  "serviceVersion",
// }

func Info(msg string) {
	Log.Info(msg)
}

func Debug(msg string) {
	Log.Debug(msg)
}

func Error(msg string) {
	Log.Error(msg)
}

func Warn(msg string) {
	Log.Warn(msg)
}

func InfoTrace(ctx context.Context, span trace.Span, msg string) {
	Log.WithFields(AddFieldsToLogger(span)).WithContext(ctx).Info(msg)
}

func DebugTrace(ctx context.Context, span trace.Span, msg string) {
	Log.WithFields(AddFieldsToLogger(span)).WithContext(ctx).Debug(msg)
}

func ErrorTrace(ctx context.Context, span trace.Span, msg string) {
	Log.WithFields(AddFieldsToLogger(span)).WithContext(ctx).Error(msg)
}

func WarnTrace(ctx context.Context, span trace.Span, msg string) {
	Log.WithFields(AddFieldsToLogger(span)).WithContext(ctx).Warn(msg)
}

func convertTraceID(id string) string {
	if len(id) < 16 {
		return ""
	}
	return id
}
