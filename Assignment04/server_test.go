package main

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"hospitaladmission/internal/equipmentreserver"
	"hospitaladmission/internal/facades"
	"hospitaladmission/internal/insuranceverifier"
	"hospitaladmission/internal/pharmacydispatcher"
	"hospitaladmission/internal/store"
)

func TestWebDemo(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "web.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	i := insuranceverifier.Verifier{Rule: insuranceverifier.Rule{}, Repository: insuranceverifier.Repository{DB: db}}
	e := equipmentreserver.Reserver{Rule: equipmentreserver.Rule{}, Repository: equipmentreserver.Repository{DB: db}}
	p := pharmacydispatcher.Dispatcher{Rule: pharmacydispatcher.Rule{}, Repository: pharmacydispatcher.Repository{DB: db}}
	f := facades.IntakeFacade{Gateway: facades.Gateway{DB: db}, Insurance: i, Equipment: e, Pharmacy: p}
	handler := newServer(db, f, facades.AdminFacade{Intake: f})
	get := func(path string) (int, map[string]any) {
		t.Helper()
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, httptest.NewRequest("GET", path, nil))
		var body map[string]any
		if strings.HasPrefix(response.Header().Get("Content-Type"), "application/json") {
			if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
				t.Fatal(err)
			}
		}
		return response.Code, body
	}
	post := func(path string, data any) (int, map[string]any) {
		t.Helper()
		raw, err := json.Marshal(data)
		if err != nil {
			t.Fatal(err)
		}
		request := httptest.NewRequest("POST", path, bytes.NewReader(raw))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)
		var body map[string]any
		if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		return response.Code, body
	}
	for _, path := range []string{"/", "/app.css", "/app.js"} {
		if status, _ := get(path); status != 200 {
			t.Fatalf("%s: %d", path, status)
		}
	}
	if status, _ := get("/missing"); status != 404 {
		t.Fatalf("missing page: %d", status)
	}
	for _, item := range []map[string]any{{"kind": "equipment", "name": "wheelchair", "quantity": 1}, {"kind": "pharmacy", "name": "paracetamol", "quantity": 0}} {
		if status, _ := post("/api/stock", item); status != 200 {
			t.Fatalf("set stock: %d", status)
		}
	}
	request := map[string]any{"action": "admit", "id": "C1", "iin": "990101123456", "coverage": "OSMS", "equipment": "wheelchair", "medicine": "paracetamol"}
	status, body := post("/api/admissions", request)
	if status != 422 || body["view"].(map[string]any)["Admission"].(map[string]any)["Status"] != facades.Compensated {
		t.Fatalf("C failure: %d %+v", status, body)
	}
	if status, body = get("/api/admissions/C1"); status != 200 || body["view"].(map[string]any)["Equipment"] != equipmentreserver.Released {
		t.Fatalf("inspect: %d %+v", status, body)
	}
	if status, body = get("/api/state"); status != 200 || len(body["admissions"].([]any)) != 1 {
		t.Fatalf("state: %d %+v", status, body)
	}
	if status, _ = post("/api/stock", map[string]any{"kind": "pharmacy", "name": "paracetamol", "quantity": 1}); status != 200 {
		t.Fatalf("set pharmacy stock: %d", status)
	}
	request["id"] = "R1"
	request["action"] = "start"
	if status, body = post("/api/admissions", request); status != 200 || body["view"].(map[string]any)["Admission"].(map[string]any)["Status"] != facades.PendingB {
		t.Fatalf("start: %d %+v", status, body)
	}
	if status, _ = post("/api/admissions/R1/resume", map[string]any{}); status != 200 {
		t.Fatalf("resume: %d", status)
	}
	if status, _ = post("/api/admissions/R1/compensate", map[string]any{}); status != 422 {
		t.Fatalf("compensate completed: %d", status)
	}
	if status, _ = get("/api/admissions/missing"); status != 404 {
		t.Fatalf("missing admission: %d", status)
	}
}

func TestWebNormalAdmissionUsesInitialStock(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "normal.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	f := facades.IntakeFacade{
		Gateway:   facades.Gateway{DB: db},
		Insurance: insuranceverifier.Verifier{Rule: insuranceverifier.Rule{}, Repository: insuranceverifier.Repository{DB: db}},
		Equipment: equipmentreserver.Reserver{Rule: equipmentreserver.Rule{}, Repository: equipmentreserver.Repository{DB: db}},
		Pharmacy:  pharmacydispatcher.Dispatcher{Rule: pharmacydispatcher.Rule{}, Repository: pharmacydispatcher.Repository{DB: db}},
	}
	handler := newServer(db, f, facades.AdminFacade{Intake: f})
	body := []byte(`{"action":"admit","iin":"990101123456","coverage":"OSMS","equipment":"wheelchair","medicine":"paracetamol"}`)
	req := httptest.NewRequest("POST", "/api/admissions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, req)
	if response.Code != 200 {
		t.Fatalf("normal admission: %d %s", response.Code, response.Body.String())
	}
	var result struct {
		View facades.View `json:"view"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result.View.Admission.Status != facades.Completed || result.View.Pharmacy != pharmacydispatcher.Dispatched {
		t.Fatalf("normal admission state: %+v", result.View)
	}
	firstID := result.View.Admission.ID
	if !strings.HasPrefix(firstID, "ADM-") || len(firstID) != 20 {
		t.Fatalf("generated admission ID: %q", firstID)
	}
	post := func(path string, payload any) *httptest.ResponseRecorder {
		t.Helper()
		data, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}
		req := httptest.NewRequest("POST", path, bytes.NewReader(data))
		req.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, req)
		return response
	}
	for _, item := range []map[string]any{
		{"kind": "equipment", "name": "oxygen monitor", "quantity": 2},
		{"kind": "pharmacy", "name": "ibuprofen", "quantity": 3},
	} {
		if response := post("/api/stock", item); response.Code != 200 {
			t.Fatalf("add stock: %d %s", response.Code, response.Body.String())
		}
	}
	response = post("/api/admissions", map[string]any{"action": "admit", "iin": "990101123456", "coverage": "OSMS", "equipment": "oxygen monitor", "medicine": "ibuprofen"})
	if response.Code != 200 {
		t.Fatalf("custom stock admission: %d %s", response.Code, response.Body.String())
	}
	if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result.View.Admission.Status != facades.Completed || result.View.Admission.Equipment != "oxygen monitor" || result.View.Admission.Medicine != "ibuprofen" {
		t.Fatalf("custom stock state: %+v", result.View)
	}
	if result.View.Admission.ID == firstID || !strings.HasPrefix(result.View.Admission.ID, "ADM-") {
		t.Fatalf("second generated admission ID: %q", result.View.Admission.ID)
	}
}
