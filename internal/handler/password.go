package handler

import (
	"embed"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/okamyuji/PasswordGenerator/internal/config"
)

// PasswordGeneratorInterface defines the contract for password generation
type PasswordGeneratorInterface interface {
	Generate(cfg config.PasswordConfig) (string, error)
}

// TemplateRendererInterface abstracts template rendering
type TemplateRendererInterface interface {
	ExecuteTemplate(w http.ResponseWriter, name string, data interface{}) error
}

// EmbedFSTemplateRenderer implements TemplateRendererInterface using html/template and embed.FS
type EmbedFSTemplateRenderer struct {
	tmpl *template.Template
}

func NewEmbedFSTemplateRenderer(content embed.FS) *EmbedFSTemplateRenderer {
	return &EmbedFSTemplateRenderer{
		tmpl: template.Must(template.ParseFS(content, "templates/*.html")),
	}
}

func (r *EmbedFSTemplateRenderer) ExecuteTemplate(w http.ResponseWriter, name string, data interface{}) error {
	return r.tmpl.ExecuteTemplate(w, name, data)
}

// PasswordHandler now depends on interfaces, not concrete implementations
type PasswordHandler struct {
	renderer  TemplateRendererInterface
	generator PasswordGeneratorInterface
}

// NewPasswordHandler creates a new PasswordHandler with dependency injection
func NewPasswordHandler(
	renderer TemplateRendererInterface,
	generator PasswordGeneratorInterface,
) *PasswordHandler {
	return &PasswordHandler{
		renderer:  renderer,
		generator: generator,
	}
}

func (h *PasswordHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.handleGet(w, r)
		return
	}
	if r.Method == http.MethodPost {
		h.handlePost(w, r)
		return
	}
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (h *PasswordHandler) handleGet(w http.ResponseWriter, _ *http.Request) {
	if err := h.renderer.ExecuteTemplate(w, "index.html", nil); err != nil {
		slog.Error("Template execution error", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *PasswordHandler) handlePost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	length, _ := strconv.Atoi(r.FormValue("length"))
	// 長さのバリデーション
	if length <= 0 {
		http.Error(w, "Invalid length", http.StatusBadRequest)
		return
	}

	pwdConfig := config.PasswordConfig{
		Length:        length,
		UseUppercase:  r.Form.Get("uppercase") == "true",
		UseLowercase:  r.Form.Get("lowercase") == "true",
		UseNumbers:    r.Form.Get("numbers") == "true",
		UseSymbols:    r.Form.Get("symbols") == "true",
		CustomSymbols: strings.TrimSpace(r.Form.Get("customSymbols")),
	}

	password, err := h.generator.Generate(pwdConfig)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	if _, err := w.Write([]byte(password)); err != nil {
		slog.Error("Failed to write password", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
