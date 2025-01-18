package handler

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/okamyuji/PasswordGenerator/internal/config"
)

// モックPasswordGeneratorの作成
type MockPasswordGenerator struct{}

func (m *MockPasswordGenerator) Generate(cfg config.PasswordConfig) (string, error) {
	// テスト用のパスワード生成ロジック
	if cfg.Length <= 0 {
		return "", fmt.Errorf("invalid length")
	}
	return strings.Repeat("A", cfg.Length), nil
}

// モックTemplateRendererの作成
type MockTemplateRenderer struct {
	tmpl *template.Template
}

func (m *MockTemplateRenderer) ExecuteTemplate(w http.ResponseWriter, name string, data interface{}) error {
	return m.tmpl.ExecuteTemplate(w, name, data)
}

func TestPasswordHandler_Handle(t *testing.T) {
	// テスト用の簡易テンプレート文字列を作成
	tmplStr := `
        <!DOCTYPE html>
        <html>
            <head>
                <title>パスワード生成ツール</title>
            </head>
            <body>
                <h1>パスワード生成ツール</h1>
            </body>
        </html>
    `
	tmpl := template.Must(template.New("index.html").Parse(tmplStr))

	// モックテンプレートレンダラーの作成
	mockRenderer := &MockTemplateRenderer{tmpl: tmpl}

	// モックパスワードジェネレーターの作成
	mockGenerator := &MockPasswordGenerator{}

	// DIを使用したハンドラーの作成
	h := NewPasswordHandler(mockRenderer, mockGenerator)

	tests := []struct {
		name         string
		method       string
		formData     url.Values
		wantStatus   int
		wantContains string
	}{
		{
			name:         "GET request",
			method:       http.MethodGet,
			wantStatus:   http.StatusOK,
			wantContains: "パスワード生成ツール",
		},
		{
			name:   "POST request - valid data",
			method: http.MethodPost,
			formData: url.Values{
				"length":    {"12"},
				"uppercase": {"true"},
				"lowercase": {"true"},
				"numbers":   {"true"},
				"symbols":   {"true"},
			},
			wantStatus: http.StatusOK,
		},
		{
			name:   "POST request - invalid length",
			method: http.MethodPost,
			formData: url.Values{
				"length":    {"0"},
				"uppercase": {"true"},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Invalid method",
			method:     http.MethodPut,
			wantStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body io.Reader
			if tt.formData != nil {
				body = strings.NewReader(tt.formData.Encode())
			}

			req := httptest.NewRequest(tt.method, "/", body)
			if tt.method == http.MethodPost {
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			}
			rec := httptest.NewRecorder()

			h.Handle(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("PasswordHandler.Handle() status = %v, want %v", rec.Code, tt.wantStatus)
			}

			if tt.wantContains != "" && !strings.Contains(rec.Body.String(), tt.wantContains) {
				t.Errorf("PasswordHandler.Handle() response does not contain %q", tt.wantContains)
			}

			if tt.method == http.MethodPost && rec.Code == http.StatusOK {
				// パスワードが生成されていることを確認
				if rec.Body.String() == "" {
					t.Error("Expected password to be generated but got empty response")
				}
			}
		})
	}
}
