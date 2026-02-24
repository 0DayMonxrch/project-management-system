package handler

import (
	"encoding/json"
	"net/http"

	"github.com/0DayMonxrch/project-management-system/internal/domain"
	"github.com/0DayMonxrch/project-management-system/internal/middleware"
)

type TaskHandler struct {
	svc domain.TaskService
}

func NewTaskHandler(svc domain.TaskService) *TaskHandler {
	return &TaskHandler{svc: svc}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		AssignedTo  string `json:"assigned_to"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	userID, _ := middleware.GetUserID(r)
	projectID := r.PathValue("projectId")

	task, err := h.svc.CreateTask(r.Context(), projectID, userID, body.Title, body.Description, body.AssignedTo)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, task)
}

func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")
	tasks, err := h.svc.ListTasks(r.Context(), projectID)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, tasks)
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")
	taskID := r.PathValue("taskId")

	task, err := h.svc.GetTask(r.Context(), projectID, taskID)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	var body map[string]any
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	userID, _ := middleware.GetUserID(r)
	taskID := r.PathValue("taskId")

	task, err := h.svc.UpdateTask(r.Context(), taskID, userID, body)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r)
	taskID := r.PathValue("taskId")

	if err := h.svc.DeleteTask(r.Context(), taskID, userID); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "task deleted successfully"})
}

func (h *TaskHandler) CreateSubTask(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	userID, _ := middleware.GetUserID(r)
	taskID := r.PathValue("taskId")

	task, err := h.svc.CreateSubTask(r.Context(), taskID, userID, body.Title)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, task)
}

func (h *TaskHandler) UpdateSubTask(w http.ResponseWriter, r *http.Request) {
	var body struct {
		IsCompleted bool `json:"is_completed"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	userID, _ := middleware.GetUserID(r)
	taskID := r.PathValue("taskId")
	subTaskID := r.PathValue("subTaskId")

	task, err := h.svc.UpdateSubTask(r.Context(), taskID, subTaskID, userID, body.IsCompleted)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) DeleteSubTask(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r)
	taskID := r.PathValue("taskId")
	subTaskID := r.PathValue("subTaskId")

	if err := h.svc.DeleteSubTask(r.Context(), taskID, subTaskID, userID); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "subtask deleted successfully"})
}