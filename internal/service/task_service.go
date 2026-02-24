package service

import (
	"context"
	"time"

	"github.com/0DayMonxrch/project-management-system/internal/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type taskService struct {
	taskRepo    domain.TaskRepository
	projectRepo domain.ProjectRepository
}

func NewTaskService(taskRepo domain.TaskRepository, projectRepo domain.ProjectRepository) domain.TaskService {
	return &taskService{taskRepo: taskRepo, projectRepo: projectRepo}
}

func (s *taskService) CreateTask(ctx context.Context, projectID, requesterID, title, description, assigneeID string) (*domain.Task, error) {
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if !hasAdminOrProjectAdmin(project, requesterID) {
		return nil, domain.ErrForbidden
	}

	projectOID, _ := bson.ObjectIDFromHex(projectID)
	requesterOID, _ := bson.ObjectIDFromHex(requesterID)
	assigneeOID, _ := bson.ObjectIDFromHex(assigneeID)

	task := &domain.Task{
		ProjectID:   projectOID,
		Title:       title,
		Description: description,
		Status:      domain.StatusTodo,
		AssignedTo:  assigneeOID,
		CreatedBy:   requesterOID,
		Attachments: []domain.Attachment{},
		SubTasks:    []domain.SubTask{},
	}

	if err := s.taskRepo.Create(ctx, task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *taskService) GetTask(ctx context.Context, projectID, taskID string) (*domain.Task, error) {
	task, err := s.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	if task.ProjectID.Hex() != projectID {
		return nil, domain.ErrNotFound
	}
	return task, nil
}

func (s *taskService) ListTasks(ctx context.Context, projectID string) ([]domain.Task, error) {
	return s.taskRepo.FindByProjectID(ctx, projectID)
}

func (s *taskService) UpdateTask(ctx context.Context, taskID, requesterID string, updates map[string]any) (*domain.Task, error) {
	task, err := s.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	project, err := s.projectRepo.FindByID(ctx, task.ProjectID.Hex())
	if err != nil {
		return nil, err
	}

	// Members can only update status
	if !hasAdminOrProjectAdmin(project, requesterID) {
		if status, ok := updates["status"].(string); ok && len(updates) == 1 {
			task.Status = domain.TaskStatus(status)
			if err := s.taskRepo.Update(ctx, task); err != nil {
				return nil, err
			}
			return task, nil
		}
		return nil, domain.ErrForbidden
	}

	if title, ok := updates["title"].(string); ok {
		task.Title = title
	}
	if desc, ok := updates["description"].(string); ok {
		task.Description = desc
	}
	if status, ok := updates["status"].(string); ok {
		task.Status = domain.TaskStatus(status)
	}
	if assignee, ok := updates["assigned_to"].(string); ok {
		oid, _ := bson.ObjectIDFromHex(assignee)
		task.AssignedTo = oid
	}

	if err := s.taskRepo.Update(ctx, task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *taskService) DeleteTask(ctx context.Context, taskID, requesterID string) error {
	task, err := s.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return err
	}
	project, err := s.projectRepo.FindByID(ctx, task.ProjectID.Hex())
	if err != nil {
		return err
	}
	if !hasAdminOrProjectAdmin(project, requesterID) {
		return domain.ErrForbidden
	}
	return s.taskRepo.Delete(ctx, taskID)
}

func (s *taskService) CreateSubTask(ctx context.Context, taskID, requesterID, title string) (*domain.Task, error) {
	task, err := s.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	project, err := s.projectRepo.FindByID(ctx, task.ProjectID.Hex())
	if err != nil {
		return nil, err
	}
	if !hasAdminOrProjectAdmin(project, requesterID) {
		return nil, domain.ErrForbidden
	}

	subtask := domain.SubTask{
		ID:        bson.NewObjectID(),
		Title:     title,
		CreatedAt: time.Now(),
	}
	task.SubTasks = append(task.SubTasks, subtask)

	if err := s.taskRepo.Update(ctx, task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *taskService) UpdateSubTask(ctx context.Context, taskID, subTaskID, requesterID string, isCompleted bool) (*domain.Task, error) {
	task, err := s.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	project, err := s.projectRepo.FindByID(ctx, task.ProjectID.Hex())
	if err != nil {
		return nil, err
	}
	if !isMember(project, requesterID) {
		return nil, domain.ErrForbidden
	}

	for i, st := range task.SubTasks {
		if st.ID.Hex() == subTaskID {
			task.SubTasks[i].IsCompleted = isCompleted
			if err := s.taskRepo.Update(ctx, task); err != nil {
				return nil, err
			}
			return task, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (s *taskService) DeleteSubTask(ctx context.Context, taskID, subTaskID, requesterID string) error {
	task, err := s.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return err
	}
	project, err := s.projectRepo.FindByID(ctx, task.ProjectID.Hex())
	if err != nil {
		return err
	}
	if !hasAdminOrProjectAdmin(project, requesterID) {
		return domain.ErrForbidden
	}

	for i, st := range task.SubTasks {
		if st.ID.Hex() == subTaskID {
			task.SubTasks = append(task.SubTasks[:i], task.SubTasks[i+1:]...)
			return s.taskRepo.Update(ctx, task)
		}
	}
	return domain.ErrNotFound
}

// --- helpers ---

func hasAdminOrProjectAdmin(p *domain.Project, userID string) bool {
	for _, m := range p.Members {
		if m.UserID.Hex() == userID && (m.Role == domain.RoleAdmin || m.Role == domain.RoleProjectAdmin) {
			return true
		}
	}
	return false
}