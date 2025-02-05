package usecase

import (
	"github.com/kihyun1998/dupdater/internal/file/domain/ports"
	"github.com/kihyun1998/dupdater/internal/file/domain/repository"
)

type fileService struct {
	repo repository.FileRepository
}

func NewFileService(repo repository.FileRepository) ports.FilePort {
	return &fileService{
		repo: repo,
	}
}

func (s *fileService) Backup() error {
	return s.repo.Backup()
}

func (s *fileService) Restore() error {
	return s.repo.Restore()
}

func (s *fileService) ExtractZip(zipFile string) error {
	return s.repo.ExtractZip(zipFile)
}

func (s *fileService) DeleteFile(path string) error {
	return s.repo.DeleteFile(path)
}
