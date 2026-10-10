package http

import (
	"context"
	"net/http"

	grpcAuth "platform/pkg/middleware/grpc"
	commonV1 "shared/pkg/proto/common/v1"
	iamV1 "shared/pkg/proto/iam/v1"
)

const SessionUUIDHeader = "X-Session-Uuid"

// IAMClient это алиас для сгенерированного gRPC клиента
type IAMClient = iamV1.IamServiceClient

// AuthMiddleware middleware для аутентификации HTTP запросов
type AuthMiddleware struct {
	iamClient IAMClient
}

// NewAuthMiddleware создает новый middleware аутентификации
func NewAuthMiddleware(iamClient IAMClient) *AuthMiddleware {
	return &AuthMiddleware{
		iamClient: iamClient,
	}
}

// Handle обрабатывает HTTP запрос с аутентификацией
func (m *AuthMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Извлекаем session UUID из заголовка
		sessionUUID := r.Header.Get(SessionUUIDHeader)
		if sessionUUID == "" {
			writeErrorResponse(w, http.StatusUnauthorized, "MISSING_SESSION", "Authentication required")
			return
		}

		// Валидируем сессию через IAM сервис
		whoamiRes, err := m.iamClient.Whoami(r.Context(), &iamV1.WhoamiRequest{
			SessionUuid: sessionUUID,
		})
		if err != nil {
			writeErrorResponse(w, http.StatusUnauthorized, "INVALID_SESSION", "Authentication failed")
			return
		}

		// Добавляем пользователя и session UUID в контекст используя функции из grpc middleware
		ctx := r.Context()
		ctx = grpcAuth.AddSessionUUIDToContext(ctx, sessionUUID)
		// Также добавляем пользователя в контекст
		ctx = context.WithValue(ctx, grpcAuth.GetUserContextKey(), whoamiRes.Session)

		// Передаем управление следующему handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserFromContext извлекает пользователя из контекста
func GetSessionFromContext(ctx context.Context) (*commonV1.Session, bool) {
	return grpcAuth.GetSessionFromContext(ctx)
}

// GetSessionUUIDFromContext извлекает session UUID из контекста
func GetSessionUUIDFromContext(ctx context.Context) (string, bool) {
	return grpcAuth.GetSessionUUIDFromContext(ctx)
}
