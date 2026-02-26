package api

import (
	"net/http"

	"github.com/joyboy1210/stolight/api/handlers"
	"github.com/joyboy1210/stolight/api/middlewares"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/upload/{bucket}", middlewares.CheckAuth(handlers.UploadHandlerAPI))
	mux.HandleFunc("GET /api/download/{bucket}/{fileId}", middlewares.CheckAuth(handlers.DownloadHandler))

	mux.HandleFunc("POST /api/buckets/", middlewares.CheckAuth(handlers.CreateBucketHandler))

	mux.HandleFunc("GET /api/buckets/{bucket}/files", middlewares.CheckAuth(handlers.ListFilesInBucketHandler))
	// mux.HandleFunc("DELETE /api/buckets/", middlewares.CheckAuth(handlers.DeleteFile))

	mux.HandleFunc("POST /api/login", handlers.Login)

	mux.HandleFunc("POST /api/admin/projects/create", middlewares.CheckAuth(handlers.CreateProjectHandler))
	mux.HandleFunc("POST /api/admin/projects/update", middlewares.CheckAuth(handlers.UpdateProject))
}
