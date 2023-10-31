package main

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin/render"
)

type TemplRenderer struct {
	Code int
	Data templ.Component
}

func (t TemplRenderer) Render(w http.ResponseWriter) error {
	w.WriteHeader(t.Code)
	if t.Data != nil {
		return t.Data.Render(context.Background(), w)
	}

	return nil
}

func (t TemplRenderer) WriteContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}

func (t *TemplRenderer) Instance(name string, data interface{}) render.Render {
	if templData, ok := data.(templ.Component); ok {
		return &TemplRenderer{
			Code: http.StatusOK,
			Data: templData,
		}
	}

	return nil
}
