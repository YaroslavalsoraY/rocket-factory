package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	// TraceIDHeader - заголовок для передачи trace ID
	TraceIDHeader = "x-trace-id"
)

func UnaryServerInterceptor(serviceName string) grpc.UnaryServerInterceptor {
	tracer := otel.GetTracerProvider().Tracer(serviceName)
	propagator := otel.GetTextMapPropagator()

	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}

		ctx = propagator.Extract(ctx, metadataCarrier(md))

		ctx, span := tracer.Start(
			ctx,
			info.FullMethod,
			trace.WithSpanKind(trace.SpanKindServer),
		)

		defer span.End()

		ctx = AddTraceIDToResponse(ctx)

		resp, err := handler(ctx, req)
		if err != nil {
			span.RecordError(err)
		}

		return resp, err
	}
}

func UnaryClientInterceptor(serviceName string) grpc.UnaryClientInterceptor {
	// Получаем текущий трейсер и пропагатор
	tracer := otel.GetTracerProvider().Tracer(serviceName)
	propagator := otel.GetTextMapPropagator()

	// Возвращаем функцию-interceptor
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// Определяем имя спана в зависимости от наличия контекста трассировки
		spanName := formatSpanName(ctx, method)

		// Создаем новый спан с подготовленным именем
		ctx, span := tracer.Start(
			ctx,
			spanName,
			trace.WithSpanKind(trace.SpanKindClient),
		)
		defer span.End()

		// Создаем переносчик метаданных для пропагации трассировки
		carrier := metadataCarrier(extractOutgoingMetadata(ctx))

		// Внедряем контекст трассировки в метаданные
		// Пропагатор добавляет информацию о текущем трейсе в метаданные запроса
		propagator.Inject(ctx, carrier)

		// Обновляем метаданные в контексте
		ctx = metadata.NewOutgoingContext(ctx, metadata.MD(carrier))

		// Вызываем следующий обработчик с обогащенным контекстом
		err := invoker(ctx, method, req, reply, cc, opts...)
		// Если произошла ошибка, записываем её в спан
		if err != nil {
			trace.SpanFromContext(ctx).RecordError(err)
		}

		return err
	}
}

// formatSpanName формирует имя спана в зависимости от наличия контекста трассировки
func formatSpanName(ctx context.Context, method string) string {
	if !trace.SpanContextFromContext(ctx).IsValid() {
		return "client." + method
	}

	return method
}

// extractOutgoingMetadata извлекает исходящие метаданные из контекста и создает их копию
func extractOutgoingMetadata(ctx context.Context) metadata.MD {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return metadata.New(nil)
	}

	return md.Copy()
}

// GetTraceIDFromContext извлекает trace ID из контекста.
// Полезно для логирования и возврата trace ID клиенту.
func GetTraceIDFromContext(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return ""
	}

	return span.SpanContext().TraceID().String()
}

// AddTraceIDToResponse добавляет trace ID в исходящие метаданные gRPC ответа.
// Это позволяет клиенту получить trace ID для последующего поиска в системе трассировки.
func AddTraceIDToResponse(ctx context.Context) context.Context {
	traceID := GetTraceIDFromContext(ctx)
	if traceID == "" {
		return ctx
	}

	// Используем вспомогательную функцию для извлечения метаданных
	md := extractOutgoingMetadata(ctx)

	md.Set(TraceIDHeader, traceID)
	return metadata.NewOutgoingContext(ctx, md)
}
