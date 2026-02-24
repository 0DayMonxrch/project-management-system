package service

import (
	"context"
	"fmt"

	"github.com/0DayMonxrch/project-management-system/internal/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type projectService struct {
	projectRepo domain.ProjectRepository
	userRepo    domain.UserRepository
}

func NewProjectService(projectRepo domain.ProjectRepository, userRepo domain.UserRepository) domain.ProjectService {
	return &projectService{projectRepo: projectRepo, userRepo: userRepo}
}

func (s *projectService) CreateProject(ctx context.Context, userID, name, description string) (*domain.Project, error) {
	oid, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, domain.ErrInvalidInput
	}

	project := &domain.Project{
		Name:        name,
		Description: description,
		CreatedBy:   oid,
		Members: []domain.ProjectMember{
			{UserID: oid, Role: domain.RoleAdmin},
		},
	}

	if err := s.projectRepo.Create(ctx, project); err != nil {
		return nil, err
	}
	return project, nil
}

func (s *projectService) GetProject(ctx context.Context, projectID, userID string) (*domain.Project, error) {
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if !isMember(project, userID) {
		return nil, domain.ErrForbidden
	}
	return project, nil
}

func (s *projectService) ListProjects(ctx context.Context, userID string) ([]domain.Project, error) {
	return s.projectRepo.FindByUserID(ctx, userID)
}

func (s *projectService) UpdateProject(ctx context.Context, projectID, userID, name, description string) (*domain.Project, error) {
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if !hasRole(project, userID, domain.RoleAdmin) {
		return nil, domain.ErrForbidden
	}

	project.Name = name
	project.Description = description
	if err := s.projectRepo.Update(ctx, project); err != nil {
		return nil, err
	}
	return project, nil
}

func (s *projectService) DeleteProject(ctx context.Context, projectID, userID string) error {
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return err
	}
	if !hasRole(project, userID, domain.RoleAdmin) {
		return domain.ErrForbidden
	}
	return s.projectRepo.Delete(ctx, projectID)
}

func (s *projectService) AddMember(ctx context.Context, projectID, requesterID, email string, role domain.Role) error {
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return err
	}
	if !hasRole(project, requesterID, domain.RoleAdmin) {
		return domain.ErrForbidden
	}

	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("user with email not found: %w", domain.ErrNotFound)
	}

	for _, m := range project.Members {
		if m.UserID == user.ID {
			return domain.ErrConflict
		}
	}

	project.Members = append(project.Members, domain.ProjectMember{UserID: user.ID, Role: role})
	return s.projectRepo.Update(ctx, project)
}

func (s *projectService) ListMembers(ctx context.Context, projectID string) ([]domain.ProjectMember, error) {
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return project.Members, nil
}

func (s *projectService) UpdateMemberRole(ctx context.Context, projectID, requesterID, targetUserID string, role domain.Role) error {
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return err
	}
	if !hasRole(project, requesterID, domain.RoleAdmin) {
		return domain.ErrForbidden
	}

	for i, m := range project.Members {
		if m.UserID.Hex() == targetUserID {
			project.Members[i].Role = role
			return s.projectRepo.Update(ctx, project)
		}
	}
	return domain.ErrNotFound
}

func (s *projectService) RemoveMember(ctx context.Context, projectID, requesterID, targetUserID string) error {
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return err
	}
	if !hasRole(project, requesterID, domain.RoleAdmin) {
		return domain.ErrForbidden
	}

	for i, m := range project.Members {
		if m.UserID.Hex() == targetUserID {
			project.Members = append(project.Members[:i], project.Members[i+1:]...)
			return s.projectRepo.Update(ctx, project)
		}
	}
	return domain.ErrNotFound
}

// --- helpers ---

func isMember(p *domain.Project, userID string) bool {
	for _, m := range p.Members {
		if m.UserID.Hex() == userID {
			return true
		}
	}
	return false
}

func hasRole(p *domain.Project, userID string, role domain.Role) bool {
	for _, m := range p.Members {
		if m.UserID.Hex() == userID && m.Role == role {
			return true
		}
	}
	return false
}