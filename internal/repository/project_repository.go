package repository

import (
	"context"
	"errors"
	"time"

	"github.com/0DayMonxrch/project-management-system/internal/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type projectRepository struct {
	col *mongo.Collection
}

func NewProjectRepository(db *mongo.Database) domain.ProjectRepository {
	return &projectRepository{col: db.Collection("projects")}
}

func (r *projectRepository) Create(ctx context.Context, project *domain.Project) error {
	project.ID = bson.NewObjectID()
	project.CreatedAt = time.Now()
	project.UpdatedAt = time.Now()

	_, err := r.col.InsertOne(ctx, project)
	return err
}

func (r *projectRepository) FindByID(ctx context.Context, id string) (*domain.Project, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.ErrInvalidInput
	}

	var project domain.Project
	err = r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&project)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, domain.ErrNotFound
	}
	return &project, err
}

func (r *projectRepository) FindByUserID(ctx context.Context, userID string) ([]domain.Project, error) {
	oid, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, domain.ErrInvalidInput
	}

	cursor, err := r.col.Find(ctx, bson.M{"members.user_id": oid})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var projects []domain.Project
	if err := cursor.All(ctx, &projects); err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *projectRepository) Update(ctx context.Context, project *domain.Project) error {
	project.UpdatedAt = time.Now()
	_, err := r.col.ReplaceOne(ctx, bson.M{"_id": project.ID}, project)
	return err
}

func (r *projectRepository) Delete(ctx context.Context, id string) error {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrInvalidInput
	}
	_, err = r.col.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}