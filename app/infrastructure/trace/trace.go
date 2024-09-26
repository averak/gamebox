package trace

import (
	"context"
	"errors"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"github.com/averak/gamebox/app/core/config"
	"github.com/averak/gamebox/app/core/ctxval"
	"github.com/averak/gamebox/app/core/logger"
	"go.opentelemetry.io/contrib/detectors/gcp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type Span struct {
	raw trace.Span
}

func (s *Span) End() {
	s.raw.End()
}

func StartSpan(ctx context.Context, name string) (context.Context, *Span) {
	tracer := otel.GetTracerProvider().Tracer("github.com/averak/gamebox/app/infrastructure/trace")
	ctx, sp := tracer.Start(ctx, name)
	if sp.SpanContext().IsValid() {
		ctx = ctxval.SetTraceID(ctx, sp.SpanContext().TraceID().String())
	}
	span := &Span{
		raw: sp,
	}
	return ctx, span
}

func Init(ctx context.Context, serviceName string, serviceVersion string, samplingRate float32) {
	exporter, err := texporter.New(texporter.WithProjectID(config.Get().GetGoogleCloud().GetProjectId()))
	if err != nil {
		// 計測できなくてもサービスは稼働できるので、ログだけ出してアプリケーションを続行する。
		logger.Emergency(ctx, map[string]interface{}{
			"message": "Failed to create trace exporter.",
			"error":   err.Error(),
		})
		return
	}

	res, err := resource.New(ctx,
		resource.WithDetectors(gcp.NewDetector()),
		resource.WithTelemetrySDK(),
		resource.WithAttributes(
			attribute.Key("service.name").String(serviceName),
			attribute.Key("service.version").String(serviceVersion),
		),
	)
	if errors.Is(err, resource.ErrPartialResource) || errors.Is(err, resource.ErrSchemaURLConflict) {
		logger.Error(ctx, map[string]interface{}{
			"message": "Failed to create resource.",
			"error":   err.Error(),
		})
	} else if err != nil {
		logger.Emergency(ctx, map[string]interface{}{
			"message": "Failed to create resource.",
			"error":   err.Error(),
		})
		return
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(float64(samplingRate))),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})
}
