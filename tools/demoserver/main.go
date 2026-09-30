// Command demoserver supplies invented API responses for an account-free VHS recording.
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

func main() {
	addr := "127.0.0.1:8645"
	if len(os.Args) > 1 {
		addr = os.Args[1]
	}
	server := &http.Server{
		Addr: addr, Handler: demoHandler(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
	fmt.Fprintln(os.Stderr, "demo API listening on", addr)
	if err := server.ListenAndServe(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func writeJSON(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(value)
}

func demoHandler() http.Handler {
	var mu sync.Mutex
	course := map[string]any{"id": 101, "name": "Demo: Practical Astronomy", "course_code": "DEMO-AST101", "workflow_state": "available", "account_id": 1}
	assignments := []any{
		map[string]any{"id": 201, "course_id": 101, "name": "Demo observation journal", "points_possible": 20, "grading_type": "points", "published": true, "submission_types": []string{"online_upload"}},
		map[string]any{"id": 202, "course_id": 101, "name": "Demo lunar phases quiz", "points_possible": 10, "grading_type": "points", "published": true, "submission_types": []string{"online_text_entry"}},
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		switch r.Method + " " + r.URL.Path {
		case "GET /api/v1/accounts":
			// The client probes accounts during version detection before any resource request.
			w.Header().Set("X-Canvas-Meta", `{"version":"canvas-2026.09.01"}`)
			writeJSON(w, []any{map[string]any{"id": 1, "name": "Demo Learning Lab"}})
		case "GET /api/v1/courses":
			writeJSON(w, []any{course, map[string]any{"id": 102, "name": "Demo: Creative Writing", "course_code": "DEMO-WR102", "workflow_state": "available", "account_id": 1}})
		case "GET /api/v1/courses/101":
			// Assignment commands validate the course before accessing its assignments.
			writeJSON(w, course)
		case "GET /api/v1/courses/101/assignments":
			writeJSON(w, assignments)
		case "POST /api/v1/courses/101/assignments":
			var body struct {
				Assignment map[string]any `json:"assignment"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Assignment == nil {
				http.Error(w, "expected an assignment object", http.StatusBadRequest)
				return
			}
			assignment := body.Assignment
			assignment["id"] = 203
			assignment["course_id"] = 101
			assignment["grading_type"] = "points"
			assignments = append(assignments, assignment)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			writeJSON(w, assignment)
		default:
			http.NotFound(w, r)
		}
	})
}
