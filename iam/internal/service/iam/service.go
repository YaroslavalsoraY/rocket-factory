package iam

import "iam/internal/repository"

type IamService struct {
	userRepository repository.UserRepository
	sessionRepository repository.SessionRepository
}

func NewIamService(userRepository repository.UserRepository, sessionRepository repository.SessionRepository) *IamService {
	return &IamService{
		userRepository: userRepository,
		sessionRepository: sessionRepository,
	}
}