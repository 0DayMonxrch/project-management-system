package repository

import (
	"context"
	"errors"
	"time"

	"github.com/0DayMonxrch/project-management-system/internal/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type noteRepository struct {
	col *mongo.Collection
}

func NewNoteRepository(db *mongo.Database) domain.NoteRepository {
	return &noteRepository{col: db.Collection("notes")}
}

func (r *noteRepository) Create(ctx context.Context, note *domain.Note) error {
	note.ID = bson.NewObjectID()
	note.CreatedAt = time.Now()
	note.UpdatedAt = time.Now()

	_, err := r.col.InsertOne(ctx, note)
	return err
}

func (r *noteRepository) FindByID(ctx context.Context, id string) (*domain.Note, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.ErrInvalidInput
	}

	var note domain.Note
	err = r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&note)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, domain.ErrNotFound
	}
	return &note, err
}

func (r *noteRepository) FindByProjectID(ctx context.Context, projectID string) ([]domain.Note, error) {
	oid, err := bson.ObjectIDFromHex(projectID)
	if err != nil {
		return nil, domain.ErrInvalidInput
	}

	cursor, err := r.col.Find(ctx, bson.M{"project_id": oid})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var notes []domain.Note
	if err := cursor.All(ctx, &notes); err != nil {
		return nil, err
	}
	return notes, nil
}

func (r *noteRepository) Update(ctx context.Context, note *domain.Note) error {
	note.UpdatedAt = time.Now()
	_, err := r.col.ReplaceOne(ctx, bson.M{"_id": note.ID}, note)
	return err
}

func (r *noteRepository) Delete(ctx context.Context, id string) error {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrInvalidInput
	}
	_, err = r.col.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}