package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type TaskStatus string

const (
	StatusTodo       TaskStatus = "todo"
	StatusInProgress TaskStatus = "in_progress"
	StatusDone       TaskStatus = "done"
)

type Attachment struct {
	URL      string `bson:"url"       json:"url"`
	MimeType string `bson:"mime_type" json:"mime_type"`
	Size     int64  `bson:"size"      json:"size"`
}

type SubTask struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title"         json:"title"`
	IsCompleted bool               `bson:"is_completed"  json:"is_completed"`
	CreatedAt   time.Time          `bson:"created_at"    json:"created_at"`
}

type Task struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id"`
	ProjectID   bson.ObjectID `bson:"project_id"    json:"project_id"`
	Title       string             `bson:"title"         json:"title"`
	Description string             `bson:"description"   json:"description"`
	Status      TaskStatus         `bson:"status"        json:"status"`
	AssignedTo  bson.ObjectID `bson:"assigned_to"   json:"assigned_to"`
	Attachments []Attachment       `bson:"attachments"   json:"attachments"`
	SubTasks    []SubTask          `bson:"subtasks"      json:"subtasks"`
	CreatedBy   bson.ObjectID `bson:"created_by"    json:"created_by"`
	CreatedAt   time.Time          `bson:"created_at"    json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"    json:"updated_at"`
}