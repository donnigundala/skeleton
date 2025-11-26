package routes

import (
	"net/http"

	contractHTTP "github.com/donnigundala/dg-core/contracts/http"
	"github.com/donnigundala/dg-core/http/response"
)

// Register registers the web routes.
func Register(router contractHTTP.Router) {
	// Welcome route
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		response.Success(w, map[string]string{
			"message": "Welcome to Skeleton V2!",
			"version": "2.0.0",
		}, "DG Framework Skeleton Application")
	})

	// API group
	router.Group(contractHTTP.GroupAttributes{
		Prefix: "/api/v1",
	}, func(r contractHTTP.Router) {
		r.Get("/status", func(w http.ResponseWriter, r *http.Request) {
			response.JSON(w, http.StatusOK, map[string]interface{}{
				"status": "online",
				"app":    "skeleton-v2",
			})
		})
	})
}
