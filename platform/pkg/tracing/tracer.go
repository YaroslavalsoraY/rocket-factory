package tracing

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	// DefaultCompressor - алгоритм сжатия по умолчанию
	DefaultCompressor = "gzip"
	// DefaultRetryEnabled - включение повторных попыток по умолчанию
	DefaultRetryEnabled = true
	// DefaultRetryInitialInterval - начальный интервал между повторными попытками
	DefaultRetryInitialInterval = 500 * time.Millisecond
	// DefaultRetryMaxInterval - максимальный интервал между повторными попытками
	DefaultRetryMaxInterval = 5 * time.Second
	// DefaultRetryMaxElapsedTime - максимальное время на все повторные попытки
	DefaultRetryMaxElapsedTime = 30 * time.Second
	// DefaultTimeout - таймаут по умолчанию для операций
	DefaultTimeout = 5 * time.Second
)

var serviceName string

type Config interface {
	CollectorEndpoint() string
	ServiceName() string
	Environment() string
	ServiceVersion() string
}

func InitTracer(ctx context.Context, cfg Config) error {
	serviceName = cfg.ServiceName()
	exporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithEndpoint(cfg.CollectorEndpoint()),
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithTimeout(DefaultTimeout),
		otlptracegrpc.WithCompressor(DefaultCompressor),
		otlptracegrpc.WithRetry(otlptracegrpc.RetryConfig{
			Enabled:         DefaultRetryEnabled,
			InitialInterval: DefaultRetryInitialInterval,
			MaxInterval:     DefaultRetryMaxInterval,
			MaxElapsedTime:  DefaultRetryMaxElapsedTime,
		}),
	)
	if err != nil {
		return err
	}

	attributeResource, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName()),
			semconv.ServiceVersion(cfg.ServiceVersion()),
			attribute.String("deployment.environment", cfg.Environment()),
		),
		resource.WithHost(),
		resource.WithOS(),
		resource.WithProcess(),
		resource.WithContainer(),
		resource.WithTelemetrySDK(),
	)
	if err != nil {
		return err
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(attributeResource),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(1.0))),
	)

	otel.SetTracerProvider(tracerProvider)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return nil
}

// ShutdownTracer закрывает глобальный трейсер OpenTelemetry.
// Функцию следует вызывать при завершении работы приложения.
func ShutdownTracer(ctx context.Context) error {
	provider := otel.GetTracerProvider()
	if provider == nil {
		return nil
	}

	// Приводим к конкретному типу для вызова Shutdown
	tracerProvider, ok := provider.(*sdktrace.TracerProvider)
	if !ok {
		return nil
	}

	// Закрываем провайдер трейсов при выходе
	err := tracerProvider.Shutdown(ctx)
	if err != nil {
		// Ошибки при закрытии не критичны, но могут привести к потере последних трейсов
		return err
	}

	return nil
}

func StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	// Получаем трейсер из глобального провайдера
	// Используем имя сервиса из конфигурации для лучшей идентификации в Jaeger
	return otel.Tracer(serviceName).Start(ctx, name, opts...)
}

// SpanFromContext возвращает текущий активный спан из контекста.
// Если спан не существует, возвращается NoopSpan.
func SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

// TraceIDFromContext извлекает trace ID из контекста.
// Возвращает строку с ID трейса или пустую строку, если трейс не найден.
func TraceIDFromContext(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return ""
	}

	return span.SpanContext().TraceID().String()
}
