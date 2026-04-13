package handlers

import (
	"html/template"
	"net/http"
	"sibir2025/service/app_logic/utils"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
