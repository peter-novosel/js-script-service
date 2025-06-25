package executor

import (
	"encoding/json"
	"net/http"

	"github.com/dop251/goja"
	"github.com/go-chi/chi/v5"

	"github.com/peter-novosel/js-script-service/internal/db"
	"github.com/peter-novosel/js-script-service/internal/logger"
)

type ExecInput struct {
	Params map[string]string      `json:"params"`
	Body   map[string]interface{} `json:"body"`
}

func HandleExecuteScript(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	script, err := db.GetScriptByPath(r.Context(), "/scripts/"+slug)
	if err != nil {
		logger.Init().WithError(err).Error("Script not found")
		http.Error(w, "Script not found", http.StatusNotFound)
		return
	}

	// Prepare input
	input := ExecInput{
		Params: map[string]string{},
		Body:   map[string]interface{}{},
	}
	for k := range r.URL.Query() {
		input.Params[k] = r.URL.Query().Get(k)
	}
	if r.Method == http.MethodPost {
		if err := json.NewDecoder(r.Body).Decode(&input.Body); err != nil {
			logger.Init().WithError(err).Warn("Invalid JSON")
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
	}

	// Execute with Goja
	vm := goja.New()
	inputJSON, _ := json.Marshal(input)
	vm.Set("input", string(inputJSON))

	var result goja.Value
	result, err = vm.RunString(script.Code)
	if err != nil {
		logger.Init().WithError(err).Error("Script execution failed")
		db.LogExecution(r.Context(), script.ID, input, nil, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return result
	output := result.Export()
	db.LogExecution(r.Context(), script.ID, input, output, "")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}
