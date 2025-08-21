package http

import (
	"encoding/json"
	"net/http"
)

type AddAppRequest struct {
	AppID string `json:"appId"`
}

func (h *Handlers) AddApp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request AddAppRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if request.AppID == "" {
		http.Error(w, "AppID is required", http.StatusBadRequest)
		return
	}

	err := h.addAppUseCase.Execute(request.AppID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	w.WriteHeader(http.StatusCreated)
}
