package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/averak/gamebox/app/core/config"
	"github.com/averak/gamebox/app/core/ctxval"
	pb "github.com/averak/gamebox/protobuf/config"
)

var (
	logger = log.New(os.Stdout, "", 0)

	sctx  serviceContext
	mutex = &sync.Mutex{}
)

func Init(service, version string) {
	mutex.Lock()
	defer mutex.Unlock()

	sctx = serviceContext{
		service: service,
		version: version,
	}
}

type serviceContext struct {
	service string
	version string
}

// デバッグ用の情報を出力する。
func Debug(ctx context.Context, msg any) {
	_, file, line, _ := runtime.Caller(1)
	output(ctx, pb.Logging_SEVERITY_DEBUG, msg, file, line)
}

// 特に対応は不要だが、対応が必要なログのヒントになることを期待する情報になる。
func Info(ctx context.Context, msg any) {
	_, file, line, _ := runtime.Caller(1)
	output(ctx, pb.Logging_SEVERITY_INFO, msg, file, line)
}

// 起動、シャットダウン、設定変更など、通常だが重要なイベント。
func Notice(ctx context.Context, msg any) {
	_, file, line, _ := runtime.Caller(1)
	output(ctx, pb.Logging_SEVERITY_NOTICE, msg, file, line)
}

// 直ちに対応は不要だが、近いうちに対応が必要になる。
func Warning(ctx context.Context, msg any) {
	_, file, line, _ := runtime.Caller(1)
	output(ctx, pb.Logging_SEVERITY_WARNING, msg, file, line)
}

// 一定の期間以上継続する場合、直ちに対応が求められる。
func Error(ctx context.Context, msg any) {
	_, file, line, _ := runtime.Caller(1)
	output(ctx, pb.Logging_SEVERITY_ERROR, msg, file, line)
}

// より深刻な問題を誘発する可能性がある状態。
// 1件でも発生したら、直ちに対応が求められる。(サービスは継続して良い)
func Critical(ctx context.Context, msg any) {
	_, file, line, _ := runtime.Caller(1)
	output(ctx, pb.Logging_SEVERITY_CRITICAL, msg, file, line)
}

// システムの一部が損なわれている可能性があるが、全体としては稼働している状態。
// 1件でも発生したら、一度サービスを停止またはフォールバックして、直ちに対応が求められる。
func Alert(ctx context.Context, msg any) {
	_, file, line, _ := runtime.Caller(1)
	output(ctx, pb.Logging_SEVERITY_ALERT, msg, file, line)
}

// システム全体が利用できなくなっている状態。
// 1件でも発生したら、一度サービスを停止またはフォールバックして、直ちに対応が求められる。
func Emergency(ctx context.Context, msg any) {
	_, file, line, _ := runtime.Caller(1)
	output(ctx, pb.Logging_SEVERITY_EMERGENCY, msg, file, line)
}

func output(ctx context.Context, severity pb.Logging_Severity, msg any, file string, line int) {
	if config.Get().GetLogging().GetSeverity() > severity {
		return
	}

	data := map[string]any{
		"timestamp": time.Now().Format(time.RFC3339),
		"severity":  severityString(severity),
		"message":   msg,
		"context": map[string]any{
			"reportLocation": map[string]any{
				"filePath":   file,
				"lineNumber": line,
			},
		},
		"serviceContext": map[string]string{
			"service": sctx.service,
			"version": sctx.version,
		},
	}

	traceID, ok := ctxval.GetTraceID(ctx)
	if ok {
		data["trace"] = fmt.Sprintf("projects/%s/traces/%s", config.Get().GetGoogleCloud().GetProjectId(), traceID)
	}

	result, err := json.Marshal(data)
	if err != nil {
		logger.Printf("unable to marshal data: %v: error: %v", data, err)
		return
	}
	logger.Println(string(result))
}

func severityString(severity pb.Logging_Severity) string {
	var result string
	switch severity {
	case pb.Logging_SEVERITY_UNSPECIFIED:
		result = ""
	case pb.Logging_SEVERITY_DEBUG:
		result = "DEBUG"
	case pb.Logging_SEVERITY_INFO:
		result = "INFO"
	case pb.Logging_SEVERITY_NOTICE:
		result = "NOTICE"
	case pb.Logging_SEVERITY_WARNING:
		result = "WARNING"
	case pb.Logging_SEVERITY_ERROR:
		result = "ERROR"
	case pb.Logging_SEVERITY_CRITICAL:
		result = "CRITICAL"
	case pb.Logging_SEVERITY_ALERT:
		result = "ALERT"
	case pb.Logging_SEVERITY_EMERGENCY:
		result = "EMERGENCY"
	}
	return result
}
