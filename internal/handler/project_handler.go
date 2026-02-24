package handler

import (
	"encoding/json"
	"net/http"

	"github.com/0DayMonxrch/project-management-system/internal/domain"
	"github.com/0DayMonxrch/project-management-system/internal/middleware"
	"github.com/0DayMonxrch/project-management-system/pkg/validator"
)

type ProjectHandler struct {
	svc domain.ProjectService
}

func NewProjectHandler(svc domain.ProjectService) *ProjectHandler {
	return &ProjectHandler{svc: svc}
}

func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if err := validator.New().
		Required("name", body.Name).
		MaxLength("name", body.Name, 100).
		Required("description", body.Description).
		Validate(); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	userID, _ := middleware.GetUserID(r)
	project, err := h.svc.CreateProject(r.Context(), userID, body.Name, body.Description)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, project)
}

func (h *ProjectHandler) ListProjects(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r)
	projects, err := h.svc.ListProjects(r.Context(), userID)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, projects)
}

func (h *ProjectHandler) GetProject(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r)
	projectID := r.PathValue("projectId")

	project, err := h.svc.GetProject(r.Context(), projectID, userID)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, project)
}

func (h *ProjectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	userID, _ := middleware.GetUserID(r)
	projectID := r.PathValue("projectId")

	project, err := h.svc.UpdateProject(r.Context(), projectID, userID, body.Name, body.Description)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, project)
}

func (h *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r)
	projectID := r.PathValue("projectId")

	if err := h.svc.DeleteProject(r.Context(), projectID, userID); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "project deleted successfully"})
}

func (h *ProjectHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email string      `json:"email"`
		Role  domain.Role `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if err := validator.New().
		Required("email", body.Email).
		Email("email", body.Email).
		Required("role", string(body.Role)).
		OneOf("role", string(body.Role), "admin", "project_admin", "member").
		Validate(); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	userID, _ := middleware.GetUserID(r)
	projectID := r.PathValue("projectId")

	if err := h.svc.AddMember(r.Context(), projectID, userID, body.Email, body.Role); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "member added successfully"})
}

func (h *ProjectHandler) ListMembers(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")
	members, err := h.svc.ListMembers(r.Context(), projectID)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, members)
}

func (h *ProjectHandler) UpdateMemberRole(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Role domain.Role `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if err := validator.New().
		Required("role", string(body.Role)).
		OneOf("role", string(body.Role), "admin", "project_admin", "member").
		Validate(); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	userID, _ := middleware.GetUserID(r)
	projectID := r.PathValue("projectId")
	targetUserID := r.PathValue("userId")

	if err := h.svc.UpdateMemberRole(r.Context(), projectID, userID, targetUserID, body.Role); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "member role updated successfully"})
}

func (h *ProjectHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r)
	projectID := r.PathValue("projectId")
	targetUserID := r.PathValue("userId")

	if err := h.svc.RemoveMember(r.Context(), projectID, userID, targetUserID); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "member removed successfully"})
}