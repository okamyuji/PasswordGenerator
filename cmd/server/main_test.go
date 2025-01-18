package main

import (
	"mime"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/okamyuji/PasswordGenerator/internal/generator"
	"github.com/okamyuji/PasswordGenerator/internal/handler"
)

func TestFileServer(t *testing.T) {
	// テスト前にMIMEタイプを設定
	if err := mime.AddExtensionType(".css", "text/css"); err != nil {
		t.Fatalf("Failed to add MIME type for .css: %v", err)
	}
	if err := mime.AddExtensionType(".js", "application/javascript"); err != nil {
		t.Fatalf("Failed to add MIME type for .js: %v", err)
	}

	tests := []struct {
		name       string
		path       string
		wantStatus int
		wantType   string
	}{
		{
			name:       "CSS file",
			path:       "/static/css/style.css",
			wantStatus: http.StatusOK,
			wantType:   "text/css",
		},
		{
			name:       "JavaScript file",
			path:       "/static/js/app.js",
			wantStatus: http.StatusOK,
			wantType:   "application/javascript",
		},
		{
			name:       "Non-existent file",
			path:       "/static/nonexistent.txt",
			wantStatus: http.StatusNotFound,
			wantType:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rec := httptest.NewRecorder()

			// 直接http.FileServerを使用してテスト
			fs := http.FileServer(http.FS(content))
			fs.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					rec.Code, tt.wantStatus)
			}

			if tt.wantType != "" && rec.Code == http.StatusOK {
				gotType := rec.Header().Get("Content-Type")
				if !strings.Contains(gotType, tt.wantType) {
					t.Errorf("handler returned wrong content type: got %v want %v",
						gotType, tt.wantType)
				}
			}
		})
	}
}

func TestServerStartup(t *testing.T) {
	// サーバー起動をゴルーチンで実行
	go func() {
		if err := startServer(":0"); err != http.ErrServerClosed {
			t.Errorf("unexpected server error: %v", err)
		}
	}()

	// サーバーの起動を少し待つ
	time.Sleep(100 * time.Millisecond)
}

// テスト用のヘルパー関数
func startServer(addr string) error {
	mux := http.NewServeMux()

	// DIを使用したコンポーネントの作成
	templateRenderer := handler.NewEmbedFSTemplateRenderer(content)
	passwordGenerator := generator.New()
	passwordHandler := handler.NewPasswordHandler(templateRenderer, passwordGenerator)

	mux.HandleFunc("/", passwordHandler.Handle)

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return server.ListenAndServe()
}

func TestMIMETypeInitialization(t *testing.T) {
	// テスト前にMIMEタイプを設定
	if err := mime.AddExtensionType(".css", "text/css"); err != nil {
		t.Fatalf("Failed to add MIME type for .css: %v", err)
	}
	if err := mime.AddExtensionType(".js", "application/javascript"); err != nil {
		t.Fatalf("Failed to add MIME type for .js: %v", err)
	}

	tests := []struct {
		name     string
		ext      string
		wantType string
	}{
		{
			name:     "CSS MIME type",
			ext:      ".css",
			wantType: "text/css",
		},
		{
			name:     "JavaScript MIME type",
			ext:      ".js",
			wantType: "application/javascript",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotType := mime.TypeByExtension(tt.ext)
			if !strings.Contains(gotType, tt.wantType) {
				t.Errorf("wrong MIME type for %s: got %v want %v",
					tt.ext, gotType, tt.wantType)
			}
		})
	}
}

func TestEmbedFileSystem(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "Template exists",
			path:    "templates/index.html",
			wantErr: false,
		},
		{
			name:    "CSS file exists",
			path:    "static/css/style.css",
			wantErr: false,
		},
		{
			name:    "JavaScript file exists",
			path:    "static/js/app.js",
			wantErr: false,
		},
		{
			name:    "Non-existent file",
			path:    "nonexistent.txt",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := content.ReadFile(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("content.ReadFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
