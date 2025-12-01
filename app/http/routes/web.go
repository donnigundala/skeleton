package routes

import (
	"fmt"
	"net/http"
	"time"

	"skeleton-v2/app/http/controllers"

	cache "github.com/donnigundala/dg-cache"
	"github.com/donnigundala/dg-core/contracts/foundation"
	contractHTTP "github.com/donnigundala/dg-core/contracts/http"
	"github.com/donnigundala/dg-core/http/response"
)

// Register registers the web routes.
func Register(app foundation.Application, router contractHTTP.Router) {
	// Welcome route
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		response.Success(w, map[string]string{
			"message": "Welcome to Skeleton V2!",
			"version": "2.0.0",
		}, "DG Framework Skeleton Application")
	})

	// Cache test endpoint
	router.Get("/test/cache", func(w http.ResponseWriter, r *http.Request) {
		cacheManagerInstance, err := app.Make("cache")
		if err != nil {
			response.Error(w, err, http.StatusInternalServerError)
			return
		}
		cacheManager := cacheManagerInstance.(*cache.Manager)

		ctx := r.Context()
		testKey := "test_key"
		testValue := fmt.Sprintf("cached_at_%d", time.Now().Unix())

		// 1. Put
		if err := cacheManager.Put(ctx, testKey, testValue, 60*time.Second); err != nil {
			response.Error(w, err, http.StatusInternalServerError)
			return
		}

		// 2. Get
		retrieved, err := cacheManager.Get(ctx, testKey)
		if err != nil {
			response.Error(w, err, http.StatusInternalServerError)
			return
		}

		// 3. Has
		exists, err := cacheManager.Has(ctx, testKey)
		if err != nil {
			response.Error(w, err, http.StatusInternalServerError)
			return
		}

		response.JSON(w, http.StatusOK, map[string]interface{}{
			"cache_working": true,
			"put":           testValue,
			"get":           retrieved,
			"has":           exists,
			"driver":        "memory",
		})
	})

	// Helper to resolve userController lazily on each request
	getUserController := func() *controllers.UserController {
		userControllerInstance, err := app.Make("userController")
		if err != nil {
			panic("failed to resolve userController: " + err.Error())
		}
		return userControllerInstance.(*controllers.UserController)
	}

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

		// User routes - controllers resolved on each request
		r.Post("/users", func(w http.ResponseWriter, r *http.Request) {
			getUserController().Create(w, r)
		})
		r.Get("/users", func(w http.ResponseWriter, r *http.Request) {
			getUserController().List(w, r)
		})
		r.Get("/users/:id", func(w http.ResponseWriter, r *http.Request) {
			getUserController().Get(w, r)
		})
		r.Put("/users/:id", func(w http.ResponseWriter, r *http.Request) {
			getUserController().Update(w, r)
		})
		r.Delete("/users/:id", func(w http.ResponseWriter, r *http.Request) {
			getUserController().Delete(w, r)
		})
	})
}
