package handler

import (
	"net/http"

	"github.com/0DayMonxrch/project-management-system/internal/middleware"
)

func RegisterRoutes(
	mux *http.ServeMux,
	auth *AuthHandler,
	project *ProjectHandler,
	task *TaskHandler,
	note *NoteHandler,
	jwtSecret string,
) {
	protected := middleware.Authenticate(jwtSecret)

	// Health check
	mux.HandleFunc("GET /api/v1/healthcheck/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// Auth routes
	mux.HandleFunc("POST /api/v1/auth/register", auth.Register)
	mux.HandleFunc("POST /api/v1/auth/login", auth.Login)
	mux.HandleFunc("POST /api/v1/auth/refresh-token", auth.RefreshToken)
	mux.HandleFunc("GET /api/v1/auth/verify-email/{verificationToken}", auth.VerifyEmail)
	mux.HandleFunc("POST /api/v1/auth/forgot-password", auth.ForgotPassword)
	mux.HandleFunc("POST /api/v1/auth/reset-password/{resetToken}", auth.ResetPassword)

	// Auth routes (protected)
	mux.Handle("POST /api/v1/auth/logout", protected(http.HandlerFunc(auth.Logout)))
	mux.Handle("GET /api/v1/auth/current-user", protected(http.HandlerFunc(auth.GetCurrentUser)))
	mux.Handle("POST /api/v1/auth/change-password", protected(http.HandlerFunc(auth.ChangePassword)))
	mux.Handle("POST /api/v1/auth/resend-email-verification", protected(http.HandlerFunc(auth.ResendVerificationEmail)))

	// Project routes (protected)
	mux.Handle("GET /api/v1/projects/", protected(http.HandlerFunc(project.ListProjects)))
	mux.Handle("POST /api/v1/projects/", protected(http.HandlerFunc(project.CreateProject)))
	mux.Handle("GET /api/v1/projects/{projectId}", protected(http.HandlerFunc(project.GetProject)))
	mux.Handle("PUT /api/v1/projects/{projectId}", protected(http.HandlerFunc(project.UpdateProject)))
	mux.Handle("DELETE /api/v1/projects/{projectId}", protected(http.HandlerFunc(project.DeleteProject)))
	mux.Handle("GET /api/v1/projects/{projectId}/members", protected(http.HandlerFunc(project.ListMembers)))
	mux.Handle("POST /api/v1/projects/{projectId}/members", protected(http.HandlerFunc(project.AddMember)))
	mux.Handle("PUT /api/v1/projects/{projectId}/members/{userId}", protected(http.HandlerFunc(project.UpdateMemberRole)))
	mux.Handle("DELETE /api/v1/projects/{projectId}/members/{userId}", protected(http.HandlerFunc(project.RemoveMember)))

	// Task routes (protected)
	mux.Handle("GET /api/v1/tasks/{projectId}", protected(http.HandlerFunc(task.ListTasks)))
	mux.Handle("POST /api/v1/tasks/{projectId}", protected(http.HandlerFunc(task.CreateTask)))
	mux.Handle("GET /api/v1/tasks/{projectId}/t/{taskId}", protected(http.HandlerFunc(task.GetTask)))
	mux.Handle("PUT /api/v1/tasks/{projectId}/t/{taskId}", protected(http.HandlerFunc(task.UpdateTask)))
	mux.Handle("DELETE /api/v1/tasks/{projectId}/t/{taskId}", protected(http.HandlerFunc(task.DeleteTask)))
	mux.Handle("POST /api/v1/tasks/{projectId}/t/{taskId}/subtasks", protected(http.HandlerFunc(task.CreateSubTask)))
	mux.Handle("PUT /api/v1/tasks/{projectId}/st/{subTaskId}", protected(http.HandlerFunc(task.UpdateSubTask)))
	mux.Handle("DELETE /api/v1/tasks/{projectId}/st/{subTaskId}", protected(http.HandlerFunc(task.DeleteSubTask)))

	// Note routes (protected)
	mux.Handle("GET /api/v1/notes/{projectId}", protected(http.HandlerFunc(note.ListNotes)))
	mux.Handle("POST /api/v1/notes/{projectId}", protected(http.HandlerFunc(note.CreateNote)))
	mux.Handle("GET /api/v1/notes/{projectId}/n/{noteId}", protected(http.HandlerFunc(note.GetNote)))
	mux.Handle("PUT /api/v1/notes/{projectId}/n/{noteId}", protected(http.HandlerFunc(note.UpdateNote)))
	mux.Handle("DELETE /api/v1/notes/{projectId}/n/{noteId}", protected(http.HandlerFunc(note.DeleteNote)))
}