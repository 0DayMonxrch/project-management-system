package repository

import (
	"context"
	"errors"
	"time"

	"github.com/0DayMonxrch/project-management-system/internal/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type taskRepository struct {
	col *mongo.Collection
}

func NewTaskRepository(db *mongo.Database) domain.TaskRepository {
	return &taskRepository{col: db.Collection("tasks")}
}

func (r *taskRepository) Create(ctx context.Context, task *domain.Task) error {
	task.ID = bson.NewObjectID()
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	_, err := r.col.InsertOne(ctx, task)
	return err
}

func (r *taskRepository) FindByID(ctx context.Context, id string) (*domain.Task, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.ErrInvalidInput
	}

	var task domain.Task
	err = r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&task)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, domain.ErrNotFound
	}
	return &task, err
}

func (r *taskRepository) FindByProjectID(ctx context.Context, projectID string) ([]domain.Task, error) {
	oid, err := bson.ObjectIDFromHex(projectID)
	if err != nil {
		return nil, domain.ErrInvalidInput
	}

	cursor, err := r.col.Find(ctx, bson.M{"project_id": oid})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []domain.Task
	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *taskRepository) Update(ctx context.Context, task *domain.Task) error {
	task.UpdatedAt = time.Now()
	_, err := r.col.ReplaceOne(ctx, bson.M{"_id": task.ID}, task)
	return err
}

func (r *taskRepository) Delete(ctx context.Context, id string) error {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrInvalidInput
	}
	_, err = r.col.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}