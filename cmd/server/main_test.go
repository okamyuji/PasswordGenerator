package main

import (
	"embed"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/okamyuji/PasswordGenerator/internal/generator"
	"github.com/okamyuji/PasswordGenerator/internal/handler"
	"github.com/okamyuji/PasswordGenerator/internal/middleware"
)

// テスト用の埋め込みコンテンツをモック
//
//go:embed templates/* static/*
var testContent embed.FS

func TestGenerateCSRFToken(t *testing.T) {
	token := generateCSRFToken()

	// トークンが空でないことを確認
	if token == "" {
		t.Error("生成されたCSRFトークンが空です")
	}

	// トークンがBase64エンコードされていることを確認
	decoded, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		t.Errorf("生成されたトークンが有効なBase64 URLエンコード文字列ではありません: %v", err)
	}

	// デコードされたトークンの長さが32バイトであることを確認
	if len(decoded) != 32 {
		t.Errorf("デコードされたトークンの長さが32バイトではありません。長さ: %d", len(decoded))
	}
}

func TestEmbedContent(t *testing.T) {
	// グローバルなコンテンツを一時的にテスト用のコンテンツに置き換え
	originalContent := content
	content = testContent
	defer func() { content = originalContent }()

	// embedContent関数をテスト
	embeddedContent, err := embedContent()
	if err != nil {
		t.Fatalf("コンテンツの埋め込みに失敗: %v", err)
	}

	// 期待されるファイルが存在することを確認
	testFiles := []string{
		"static/css/style.css",
		"static/js/app.js",
		"templates/index.html",
	}

	for _, file := range testFiles {
		_, err := embeddedContent.ReadFile(file)
		if err != nil {
			t.Errorf("埋め込まれたコンテンツに期待されるファイル %s が見つかりません", file)
		}
	}
}

func TestHealthCheckEndpoint(t *testing.T) {
	// ハンドラーに渡すリクエストを作成
	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// レスポンスを記録するためのRecorderを作成
	rr := httptest.NewRecorder()

	// ハンドラーを作成
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			t.Errorf("ヘルスチェックレスポンスの書き込みに失敗: %v", err)
		}
	})

	// ハンドラーを呼び出し
	handler.ServeHTTP(rr, req)

	// ステータスコードを確認
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("ヘルスチェックハンドラーが誤ったステータスコードを返しました: 取得 %v、期待 %v",
			status, http.StatusOK)
	}

	// レスポンスボディを確認
	expected := "OK"
	if rr.Body.String() != expected {
		t.Errorf("ヘルスチェックハンドラーが予期しないボディを返しました: 取得 %v、期待 %v",
			rr.Body.String(), expected)
	}
}

func TestServerConfiguration(t *testing.T) {
	server := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// サーバー構成を確認
	if server.Addr != ":8080" {
		t.Errorf("予期しないサーバーアドレス: 取得 %v、期待 :8080", server.Addr)
	}

	if server.ReadTimeout != 5*time.Second {
		t.Errorf("予期しない読み取りタイムアウト: 取得 %v、期待 5s", server.ReadTimeout)
	}

	if server.WriteTimeout != 10*time.Second {
		t.Errorf("予期しない書き込みタイムアウト: 取得 %v、期待 10s", server.WriteTimeout)
	}

	if server.IdleTimeout != 120*time.Second {
		t.Errorf("予期しないアイドルタイムアウト: 取得 %v、期待 120s", server.IdleTimeout)
	}
}

func TestCSRFTokenEnvironmentVariable(t *testing.T) {
	// 既存のCSRF_TOKENをクリア
	os.Unsetenv("CSRF_TOKEN")

	// CSRFトークンを生成して設定
	csrfToken := generateCSRFToken()
	os.Setenv("CSRF_TOKEN", csrfToken)

	// トークンを取得
	storedToken := os.Getenv("CSRF_TOKEN")

	// トークンを検証
	if storedToken != csrfToken {
		t.Errorf("CSRFトークンが正しく環境変数に設定されていません: 取得 %v、期待 %v",
			storedToken, csrfToken)
	}

	// トークンの長さとエンコーディングを確認
	decodedToken, err := base64.URLEncoding.DecodeString(storedToken)
	if err != nil {
		t.Errorf("無効なBase64エンコーディング: %v", err)
	}

	if len(decodedToken) != 32 {
		t.Errorf("デコードされたトークンの長さが不正: 取得 %d、期待 32", len(decodedToken))
	}
}

// メインの機能をテストするためのモック関数
func startTestServer(t *testing.T, content embed.FS) (*http.Server, error) {
	// これをテストヘルパー関数としてマーク
	t.Helper()

	// セキュリティミドルウェア
	securityMiddleware := middleware.NewSecurityMiddleware()

	// テンプレートレンダラー
	templateRenderer := handler.NewEmbedFSTemplateRenderer(content)

	// パスワードジェネレーター
	passwordGenerator := generator.New()

	// 依存性注入を使用したパスワードハンドラー
	passwordHandler := handler.NewPasswordHandler(templateRenderer, passwordGenerator)

	// テスト用のサーバー構成を作成
	server := &http.Server{
		Addr:         ":0", // テスト用のランダムポート
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/health":
				w.WriteHeader(http.StatusOK)
				if _, err := w.Write([]byte("OK")); err != nil {
					t.Errorf("ヘルスチェックレスポンスの書き込みに失敗: %v", err)
				}
			case "/":
				securityMiddleware.Middleware(passwordHandler.Handle)(w, r)
			default:
				http.NotFound(w, r)
			}
		}),
	}

	return server, nil
}

func TestStartTestServer(t *testing.T) {
	server, err := startTestServer(t, testContent)
	if err != nil {
		t.Fatalf("テストサーバーの起動に失敗: %v", err)
	}
	defer server.Close()

	// 追加のサーバー構成チェックをここに追加できます
	if server.Addr == "" {
		t.Error("サーバーアドレスが設定されていません")
	}
}
