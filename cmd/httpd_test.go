package cmd

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const testCSV = `Date,Visitor,,Details,,Home,Location
"2025-03-23T19:20:00.000Z","Mid Ice Crisis","Mid Ice Crisis","L 6 - 1","HatTrick Swayzes","HatTrick Swayzes","North"
`

func TestHandleFormGet(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handleForm(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `enctype="multipart/form-data"`) {
		t.Error("form missing multipart enctype")
	}
}

func TestHandleConvertRawPost(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/convert", strings.NewReader(testCSV))
	rec := httptest.NewRecorder()
	handleConvert(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "text/csv" {
		t.Errorf("Content-Type = %q, want text/csv", ct)
	}
	if !strings.Contains(rec.Body.String(), "Mid Ice Crisis") {
		t.Errorf("response missing team name: %s", rec.Body.String())
	}
}

func TestHandleConvertMultipartPost(t *testing.T) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	fw, err := w.CreateFormFile("file", "input.csv")
	if err != nil {
		t.Fatal(err)
	}
	fw.Write([]byte(testCSV))
	w.Close()

	req := httptest.NewRequest(http.MethodPost, "/convert", &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	rec := httptest.NewRecorder()
	handleConvert(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Mid Ice Crisis") {
		t.Errorf("response missing team name: %s", rec.Body.String())
	}
	if cd := rec.Header().Get("Content-Disposition"); !strings.Contains(cd, "schedule.csv") {
		t.Errorf("Content-Disposition = %q, want attachment with schedule.csv", cd)
	}
}

func TestHandleConvertGetRejects(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/convert", nil)
	rec := httptest.NewRecorder()
	handleConvert(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want 405", rec.Code)
	}
}
