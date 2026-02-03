package api

import (
	"net/http"

	"github.com/joyboy1210/stolight/api/handlers"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/upload", handlers.UploadHandlerAPI)
	mux.HandleFunc("/download/", handlers.DownloadHandler)
}
