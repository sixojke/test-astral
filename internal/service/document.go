package service

import (
	"errors"
	"os"

	"github.com/sixojke/test-astral/domain"
	"github.com/sixojke/test-astral/internal/repository"
	"github.com/sixojke/test-astral/pkg/logger"
)

type DocumentService struct {
	repo     repository.Document
	repoUser repository.User
}

func NewDocumentService(repo repository.Document, repoUser repository.User) *DocumentService {
	return &DocumentService{
		repo:     repo,
		repoUser: repoUser,
	}
}

func (s *DocumentService) Create(document *domain.Document, userId string) error {
	return s.repo.Create(document, userId)
}

func (s *DocumentService) GetByUser(userLogin, currentUserId string, params *domain.FilterParams) (*[]domain.Document, error) {
	logger.Debugf("login=%v", userLogin)
	if userLogin == "" {
		return s.repo.GetCurrentUserDocuments(currentUserId, params)
	}

	userId, err := s.repoUser.GetUserIdByLogin(userLogin)
	if err != nil {
		logger.Errorf("failed to get userId by login: %v", err)
		return nil, err
	}

	return s.repo.GetOtherUserDocuments(userId, currentUserId, params)
}

func (s *DocumentService) GetById(documentId, userId string) (*domain.Document, error) {
	return s.repo.GetById(documentId, userId)
}

func (s *DocumentService) CheckById(documentId, userId string) (bool, error) {
	return s.repo.CheckById(documentId, userId)
}

func (s *DocumentService) Delete(documentId, userId string) error {
	filePath, err := s.repo.Delete(documentId, userId)
	if err != nil {
		if errors.Is(err, domain.ErrDocumentNotFound) {
			return err
		}

		logger.Errorf("failed to delete document: %v", err)
		return err
	}

	if err = os.Remove(filePath); err != nil {
		logger.Errorf("failed to delete file: %v", err)
	}

	return nil
}
