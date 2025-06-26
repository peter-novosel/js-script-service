package executor

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	v8 "rogchap.com/v8go"
	"go.kuoruan.net/v8go-polyfills/fetch"

	"github.com/peter-novosel/js-script-service/internal/db"
	"github.com/peter-novosel/js-script-service/internal/logger"
)

// ExecInput holds the parameters and body for script execution.
type ExecInput struct {
	Params map[string]string      `json:"params"`
	Body   map[string]interface{} `json:"body"`
}

// HandleExecuteScript runs a stored JS script via V8Go, with fetch support.
func HandleExecuteScript(w http.ResponseWriter, r *http.Request) {
	// 1) Load script
	slug := chi.URLParam(r, "slug")
	script, err := db.GetScriptByPath(r.Context(), "/scripts/"+slug)
	if err != nil {
		logger.Init().WithError(err).Error("Script not found")
		http.Error(w, "Script not found", http.StatusNotFound)
		return
	}

	// 2) Build input JSON
	input := ExecInput{Params: map[string]string{}, Body: map[string]interface{}{}}
	for k := range r.URL.Query() {
		input.Params[k] = r.URL.Query().Get(k)
	}
	if r.Method == http.MethodPost {
		if err := json.NewDecoder(r.Body).Decode(&input.Body); err != nil {
			logger.Init().WithError(err).Warn("Invalid JSON body")
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
	}
	rawInput, _ := json.Marshal(input)

	// 3) Initialize V8 isolate, inject fetch polyfill
	iso := v8.NewIsolate()
	global := v8.NewObjectTemplate(iso)
	if err := fetch.InjectTo(iso, global); err != nil {
		logger.Init().WithError(err).Error("Failed to inject fetch polyfill")
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	// 4) Create V8 context
	ctx := v8.NewContext(iso, global)
	ctx.Global().Set("input", string(rawInput))

	// 5) Wrap user script in an async IIFE that JSON.stringify's the result
	source := `(function(){
	return (async () => {
	` + script.Code + `
	})().then(res => JSON.stringify(res));
})();`

	// 6) Run and unwrap the Promise
	val, err := ctx.RunScript(source, "user.js")
	if err != nil {
		logger.Init().WithError(err).Error("Script execution failed")
		db.LogExecution(r.Context(), script.ID, input, nil, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	promise, err := val.AsPromise()
	if err != nil {
		logger.Init().WithError(err).Error("Expected a Promise")
		http.Error(w, "Execution did not return a Promise", http.StatusInternalServerError)
		return
	}

	// wait for resolution
	for promise.State() == v8.Pending {
		// spin until done
	}

	var output interface{}
	switch promise.State() {
	case v8.Fulfilled:
		jsonStr := promise.Result().String()
		if err := json.Unmarshal([]byte(jsonStr), &output); err != nil {
			logger.Init().WithError(err).Error("Invalid JSON from script result")
			http.Error(w, "Invalid script output", http.StatusInternalServerError)
			return
		}
	case v8.Rejected:
		execErr := promise.Result().String()
		logger.Init().Error("Script rejected: " + execErr)
		http.Error(w, execErr, http.StatusInternalServerError)
		return
	}

	// 7) Log and return
	db.LogExecution(r.Context(), script.ID, input, output, "")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}
