package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"hospitaladmission/internal/admission"
)

type webServer struct {
	db     *sql.DB
	intake admission.IntakeFacade
	admin  admission.AdminFacade
}

func newServer(db *sql.DB, intake admission.IntakeFacade, admin admission.AdminFacade) http.Handler {
	s := webServer{db: db, intake: intake, admin: admin}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", s.index)
	mux.HandleFunc("GET /app.css", s.asset)
	mux.HandleFunc("GET /app.js", s.asset)
	mux.HandleFunc("GET /api/state", s.state)
	mux.HandleFunc("POST /api/stock", s.stock)
	mux.HandleFunc("POST /api/admissions", s.admit)
	mux.HandleFunc("GET /api/admissions/{id}", s.inspect)
	mux.HandleFunc("POST /api/admissions/{id}/resume", s.resume)
	mux.HandleFunc("POST /api/admissions/{id}/compensate", s.compensate)
	return mux
}

func (s webServer) index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	data, err := os.ReadFile("web/index.html")
	if err != nil {
		http.Error(w, "page unavailable", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(data)
}
func (s webServer) asset(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/")
	data, err := os.ReadFile("web/" + name)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if strings.HasSuffix(name, ".css") {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
	} else {
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
	}
	w.Write(data)
}
func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func readJSON(r *http.Request, dst any) error {
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		return errors.New("send JSON content")
	}
	decoder := json.NewDecoder(io.LimitReader(r.Body, 65536))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dst); err != nil {
		return err
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return errors.New("send one JSON object")
	}
	return nil
}
func bad(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]any{"error": err.Error()})
}

func (s webServer) state(w http.ResponseWriter, r *http.Request) {
	items, err := s.intake.Gateway.List()
	if err != nil {
		bad(w, 500, err)
		return
	}
	views := make([]admission.View, 0, len(items))
	for _, item := range items {
		v, err := s.admin.Inspect(item.ID)
		if err != nil {
			bad(w, 500, err)
			return
		}
		views = append(views, v)
	}
	rows, err := s.db.Query(`SELECT kind,name,quantity FROM stock ORDER BY kind,name`)
	if err != nil {
		bad(w, 500, err)
		return
	}
	defer rows.Close()
	stock := []struct {
		Kind     string `json:"kind"`
		Name     string `json:"name"`
		Quantity int    `json:"quantity"`
	}{}
	for rows.Next() {
		var item struct {
			Kind     string `json:"kind"`
			Name     string `json:"name"`
			Quantity int    `json:"quantity"`
		}
		if err := rows.Scan(&item.Kind, &item.Name, &item.Quantity); err != nil {
			bad(w, 500, err)
			return
		}
		stock = append(stock, item)
	}
	if err := rows.Err(); err != nil {
		bad(w, 500, err)
		return
	}
	writeJSON(w, 200, map[string]any{"admissions": views, "stock": stock})
}
func (s webServer) stock(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Kind     string `json:"kind"`
		Name     string `json:"name"`
		Quantity int    `json:"quantity"`
	}
	if err := readJSON(r, &req); err != nil {
		bad(w, 400, err)
		return
	}
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		bad(w, 400, errors.New("item name is required"))
		return
	}
	var err error
	switch req.Kind {
	case "equipment":
		err = s.intake.Equipment.Repository.SetStock(req.Name, req.Quantity)
	case "pharmacy":
		err = s.intake.Pharmacy.Repository.SetStock(req.Name, req.Quantity)
	default:
		bad(w, 400, errors.New("kind must be equipment or pharmacy"))
		return
	}
	if err != nil {
		bad(w, 400, err)
		return
	}
	writeJSON(w, 200, map[string]any{"message": fmt.Sprintf("%s stock updated", req.Kind)})
}
func (s webServer) admit(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Action    string `json:"action"`
		ID        string `json:"id"`
		IIN       string `json:"iin"`
		Policy    string `json:"policy"`
		Equipment string `json:"equipment"`
		Medicine  string `json:"medicine"`
	}
	if err := readJSON(r, &req); err != nil {
		bad(w, 400, err)
		return
	}
	a := admission.Admission{ID: strings.TrimSpace(req.ID), IIN: strings.TrimSpace(req.IIN), Policy: strings.TrimSpace(req.Policy), Equipment: strings.TrimSpace(req.Equipment), Medicine: strings.TrimSpace(req.Medicine)}
	var err error
	switch req.Action {
	case "start":
		err = s.intake.Start(a)
	case "admit":
		err = s.intake.Admit(a)
	default:
		bad(w, 400, errors.New("action must be start or admit"))
		return
	}
	s.result(w, a.ID, err)
}
func (s webServer) inspect(w http.ResponseWriter, r *http.Request) {
	v, err := s.admin.Inspect(r.PathValue("id"))
	if err == sql.ErrNoRows {
		bad(w, 404, errors.New("admission not found"))
		return
	}
	if err != nil {
		bad(w, 500, err)
		return
	}
	writeJSON(w, 200, map[string]any{"view": v})
}
func (s webServer) resume(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	s.result(w, id, s.intake.Resume(id))
}
func (s webServer) compensate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	s.result(w, id, s.admin.Compensate(id))
}
func (s webServer) result(w http.ResponseWriter, id string, actionErr error) {
	v, viewErr := s.admin.Inspect(id)
	if actionErr != nil {
		status := 422
		if actionErr == sql.ErrNoRows {
			status = 404
		}
		body := map[string]any{"error": actionErr.Error()}
		if viewErr == nil {
			body["view"] = v
		}
		writeJSON(w, status, body)
		return
	}
	if viewErr != nil {
		bad(w, 500, viewErr)
		return
	}
	writeJSON(w, 200, map[string]any{"view": v})
}
