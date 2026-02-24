package service

import (
	"context"

	"github.com/0DayMonxrch/project-management-system/internal/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type noteService struct {
	noteRepo    domain.NoteRepository
	projectRepo domain.ProjectRepository
}

func NewNoteService(noteRepo domain.NoteRepository, projectRepo domain.ProjectRepository) domain.NoteService {
	return &noteService{noteRepo: noteRepo, projectRepo: projectRepo}
}

func (s *noteService) CreateNote(ctx context.Context, projectID, requesterID, title, content string) (*domain.Note, error) {
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if !hasRole(project, requesterID, domain.RoleAdmin) {
		return nil, domain.ErrForbidden
	}

	projectOID, _ := bson.ObjectIDFromHex(projectID)
	requesterOID, _ := bson.ObjectIDFromHex(requesterID)

	note := &domain.Note{
		ProjectID: projectOID,
		Title:     title,
		Content:   content,
		CreatedBy: requesterOID,
	}

	if err := s.noteRepo.Create(ctx, note); err != nil {
		return nil, err
	}
	return note, nil
}

func (s *noteService) GetNote(ctx context.Context, projectID, noteID string) (*domain.Note, error) {
	note, err := s.noteRepo.FindByID(ctx, noteID)
	if err != nil {
		return nil, err
	}
	if note.ProjectID.Hex() != projectID {
		return nil, domain.ErrNotFound
	}
	return note, nil
}

func (s *noteService) ListNotes(ctx context.Context, projectID string) ([]domain.Note, error) {
	return s.noteRepo.FindByProjectID(ctx, projectID)
}

func (s *noteService) UpdateNote(ctx context.Context, noteID, requesterID, title, content string) (*domain.Note, error) {
	note, err := s.noteRepo.FindByID(ctx, noteID)
	if err != nil {
		return nil, err
	}
	project, err := s.projectRepo.FindByID(ctx, note.ProjectID.Hex())
	if err != nil {
		return nil, err
	}
	if !hasRole(project, requesterID, domain.RoleAdmin) {
		return nil, domain.ErrForbidden
	}

	note.Title = title
	note.Content = content
	if err := s.noteRepo.Update(ctx, note); err != nil {
		return nil, err
	}
	return note, nil
}

func (s *noteService) DeleteNote(ctx context.Context, noteID, requesterID string) error {
	note, err := s.noteRepo.FindByID(ctx, noteID)
	if err != nil {
		return err
	}
	project, err := s.projectRepo.FindByID(ctx, note.ProjectID.Hex())
	if err != nil {
		return err
	}
	if !hasRole(project, requesterID, domain.RoleAdmin) {
		return domain.ErrForbidden
	}
	return s.noteRepo.Delete(ctx, noteID)
}