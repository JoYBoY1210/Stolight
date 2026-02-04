package api

import (
	"net/http"

	"github.com/joyboy1210/stolight/api/handlers"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/upload", handlers.UploadHandlerAPI)
	mux.HandleFunc("/api/download/", handlers.DownloadHandler)

	mux.HandleFunc("/api/buckets/create", handlers.CreateBucketHandler)

	mux.HandleFunc("/api/files/all", handlers.ListFilesInBucketHandler)
	mux.HandleFunc("/api/files/delete", handlers.DeleteFile)
}
