package tracing

import (
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"

	"rhyme80/redis-server-rkgrpc/model"
)

type Tracing struct {
	logger      model.Logger
	timezone    string
	maxCharBody int
	isDebugMode bool
}

func NewTracing(logger model.Logger, timezone string, maxCharBody int, isDebugMode bool) Tracing {
	return Tracing{logger: logger, timezone: timezone, maxCharBody: maxCharBody, isDebugMode: isDebugMode}
}

func (t Tracing) Logger(ctx context.Context, dataLogger model.DataLogger) {
	if t.isDebugMode {
		t.loggerDebug(ctx, dataLogger)
	}
}

func (t Tracing) loggerDebug(ctx context.Context, dataLogger model.DataLogger) {
	dataMetadata, _ := metadata.FromIncomingContext(ctx)

	eventStatus := "Ended"
	resCode := "OK"
	if dataLogger.IsError {
		eventStatus = "Unfinished"
		resCode = "500"
	}

	info := []any{
		"startTime", dataLogger.StartTime,
		"endTime", dataLogger.EndTime,
		"elapseNano", dataLogger.ElapsedNano,
		"timezone", t.timezone,
		"error", dataLogger.Error,
		"headersRequest", dataMetadata,
		"bodyRequest", cutData(t.maxCharBody, dataLogger.Body),
		"remoteAddr", getFirstOfSlice(dataMetadata.Get("hostname")),
		"operation", dataLogger.Operation,
		"evenStatus", eventStatus,
		"requestId", getFirstOfSlice(dataMetadata.Get("request-id")),
		"path", dataLogger.Where,
		"resCode", resCode,
		"bodyResponse", dataLogger.BodyResponse,
	}

	if dataLogger.IsError {
		t.logger.Errorw("log", info...)
		return
	}

	t.logger.Infow("log", info...)
}

func (t Tracing) InitLogger(port string) {
	t.logger.Infow("Bootstrap Init", "port", port)
}

func cutData(maxChar int, value any) string {
	valueString := fmt.Sprintf("%v", value)
	if len(valueString) <= maxChar {
		return valueString
	}

	return valueString[:maxChar]
}

func getFirstOfSlice(data []string) string {
	first := ""
	if len(data) > 0 {
		first = data[0]
	}
	return first
}
