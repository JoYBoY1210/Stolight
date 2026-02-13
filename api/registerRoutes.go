package api

import (
	"net/http"

	"github.com/joyboy1210/stolight/api/handlers"
	"github.com/joyboy1210/stolight/api/middlewares"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/upload/", middlewares.CheckAuth(handlers.UploadHandlerAPI))
	mux.HandleFunc("/api/download/", middlewares.CheckAuth(handlers.DownloadHandler))

	mux.HandleFunc("/api/buckets/", middlewares.CheckAuth(handlers.CreateBucketHandler))

	mux.HandleFunc("/api/files/all/", middlewares.CheckAuth(handlers.ListFilesInBucketHandler))
	mux.HandleFunc("/api/files/delete/", middlewares.CheckAuth(handlers.DeleteFile))

	mux.HandleFunc("/api/login", handlers.Login)

	mux.HandleFunc("/api/admin/projects/create", middlewares.CheckAuth(handlers.CreateProjectHandler))
}
