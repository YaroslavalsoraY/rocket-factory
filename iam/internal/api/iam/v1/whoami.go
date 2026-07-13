package iam_v1

import (
	"context"
	"errors"
	"iam/internal/converter"
	"iam/internal/model"
	iam_v1 "shared/pkg/proto/iam/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (iam *IamApi) Whoami(ctx context.Context, req *iam_v1.WhoamiRequest) (*iam_v1.WhoamiResponse, error) {
	session, err := iam.sessionService.WhoAmI(ctx, req.SessionUuid)
	if err != nil {
		if errors.Is(err, model.ErrSessionNotFound) {
			return nil, status.Error(codes.Unauthenticated, model.ErrSessionNotFound.Error())
		}

		return nil, status.Error(codes.Internal, model.ErrInternal.Error())
	}

	return &iam_v1.WhoamiResponse{Session: converter.SessionToProto(session)}, nil
}