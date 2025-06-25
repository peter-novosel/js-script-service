package admin

import (
	"encoding/json"
	"net/http"

	"github.com/dop251/goja"

	"github.com/peter-novosel/js-script-service/internal/db"
	"github.com/peter-novosel/js-script-service/internal/logger"
	
)

type ScriptRequest struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
	Code string `json:"code"`
	Enabled *bool   `json:"enabled"` // Optional, defaults to true
}

func CreateOrUpdateScript(w http.ResponseWriter, r *http.Request) {
	var req ScriptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Slug == "" {
		http.Error(w, "slug is required", http.StatusBadRequest)
		return
	}

	if req.Enabled == nil || *req.Enabled {
		// If creating or enabling a script, require code
		if req.Code == "" {
			http.Error(w, "code is required when enabled", http.StatusBadRequest)
			return
		}
	}

	// ✅ Validate JS syntax
	if _, err := goja.Compile("", req.Code, false); err != nil {
		logger.Init().WithError(err).Warn("invalid JavaScript")
		http.Error(w, "invalid JavaScript: "+err.Error(), http.StatusBadRequest)
		return
	}

	// ✅ Save to DB
	isEnabled := true
	if req.Enabled != nil {
		isEnabled = *req.Enabled
	}
	err := db.UpsertScript(r.Context(), req.Name, req.Slug, req.Code, isEnabled)
	if err != nil {
		logger.Init().WithError(err).Error("failed to upsert script")
		http.Error(w, "could not save script", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"status":"ok"}`))
}


func ListScripts(w http.ResponseWriter, r *http.Request) {
	scripts, err := db.ListScripts(r.Context())
	if err != nil {
		logger.Init().WithError(err).Error("failed to list scripts")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(scripts)
}

