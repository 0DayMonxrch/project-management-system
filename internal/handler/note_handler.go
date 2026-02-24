package handler

import (
	"encoding/json"
	"net/http"

	"github.com/0DayMonxrch/project-management-system/internal/domain"
	"github.com/0DayMonxrch/project-management-system/internal/middleware"
)

type NoteHandler struct {
	svc domain.NoteService
}

func NewNoteHandler(svc domain.NoteService) *NoteHandler {
	return &NoteHandler{svc: svc}
}

func (h *NoteHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	userID, _ := middleware.GetUserID(r)
	projectID := r.PathValue("projectId")

	note, err := h.svc.CreateNote(r.Context(), projectID, userID, body.Title, body.Content)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, note)
}

func (h *NoteHandler) ListNotes(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")
	notes, err := h.svc.ListNotes(r.Context(), projectID)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, notes)
}

func (h *NoteHandler) GetNote(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")
	noteID := r.PathValue("noteId")

	note, err := h.svc.GetNote(r.Context(), projectID, noteID)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, note)
}

func (h *NoteHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	userID, _ := middleware.GetUserID(r)
	noteID := r.PathValue("noteId")

	note, err := h.svc.UpdateNote(r.Context(), noteID, userID, body.Title, body.Content)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, note)
}

func (h *NoteHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r)
	noteID := r.PathValue("noteId")

	if err := h.svc.DeleteNote(r.Context(), noteID, userID); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "note deleted successfully"})
}