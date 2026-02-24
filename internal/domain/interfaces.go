package domain

import "context"

// --- Repository Interfaces ---

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByVerificationToken(ctx context.Context, token string) (*User, error)
	FindByResetToken(ctx context.Context, token string) (*User, error)
	Update(ctx context.Context, user *User) error
}

type ProjectRepository interface {
	Create(ctx context.Context, project *Project) error
	FindByID(ctx context.Context, id string) (*Project, error)
	FindByUserID(ctx context.Context, userID string) ([]Project, error)
	Update(ctx context.Context, project *Project) error
	Delete(ctx context.Context, id string) error
}

type TaskRepository interface {
	Create(ctx context.Context, task *Task) error
	FindByID(ctx context.Context, id string) (*Task, error)
	FindByProjectID(ctx context.Context, projectID string) ([]Task, error)
	Update(ctx context.Context, task *Task) error
	Delete(ctx context.Context, id string) error
}

type NoteRepository interface {
	Create(ctx context.Context, note *Note) error
	FindByID(ctx context.Context, id string) (*Note, error)
	FindByProjectID(ctx context.Context, projectID string) ([]Note, error)
	Update(ctx context.Context, note *Note) error
	Delete(ctx context.Context, id string) error
}

// --- Service Interfaces ---

type AuthService interface {
	Register(ctx context.Context, name, email, password string) error
	Login(ctx context.Context, email, password string) (accessToken, refreshToken string, err error)
	Logout(ctx context.Context, userID string) error
	VerifyEmail(ctx context.Context, token string) error
	RefreshToken(ctx context.Context, refreshToken string) (accessToken string, err error)
	ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, newPassword string) error
	ResendVerificationEmail(ctx context.Context, userID string) error
	GetCurrentUser(ctx context.Context, userID string) (*User, error)
}

type ProjectService interface {
	CreateProject(ctx context.Context, userID, name, description string) (*Project, error)
	GetProject(ctx context.Context, projectID, userID string) (*Project, error)
	ListProjects(ctx context.Context, userID string) ([]Project, error)
	UpdateProject(ctx context.Context, projectID, userID, name, description string) (*Project, error)
	DeleteProject(ctx context.Context, projectID, userID string) error
	AddMember(ctx context.Context, projectID, requesterID, email string, role Role) error
	ListMembers(ctx context.Context, projectID string) ([]ProjectMember, error)
	UpdateMemberRole(ctx context.Context, projectID, requesterID, targetUserID string, role Role) error
	RemoveMember(ctx context.Context, projectID, requesterID, targetUserID string) error
}

type TaskService interface {
	CreateTask(ctx context.Context, projectID, requesterID, title, description, assigneeID string) (*Task, error)
	GetTask(ctx context.Context, projectID, taskID string) (*Task, error)
	ListTasks(ctx context.Context, projectID string) ([]Task, error)
	UpdateTask(ctx context.Context, taskID, requesterID string, updates map[string]any) (*Task, error)
	DeleteTask(ctx context.Context, taskID, requesterID string) error
	CreateSubTask(ctx context.Context, taskID, requesterID, title string) (*Task, error)
	UpdateSubTask(ctx context.Context, taskID, subTaskID, requesterID string, isCompleted bool) (*Task, error)
	DeleteSubTask(ctx context.Context, taskID, subTaskID, requesterID string) error
}

type NoteService interface {
	CreateNote(ctx context.Context, projectID, requesterID, title, content string) (*Note, error)
	GetNote(ctx context.Context, projectID, noteID string) (*Note, error)
	ListNotes(ctx context.Context, projectID string) ([]Note, error)
	UpdateNote(ctx context.Context, noteID, requesterID, title, content string) (*Note, error)
	DeleteNote(ctx context.Context, noteID, requesterID string) error
}

type EmailService interface {
	SendVerificationEmail(to, token string) error
	SendPasswordResetEmail(to, token string) error
}