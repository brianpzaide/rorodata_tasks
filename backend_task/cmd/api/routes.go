package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func routes(app *application) http.Handler {
	router := chi.NewRouter()

	router.NotFound(app.notFoundResponse)
	router.MethodNotAllowed(app.methodNotAllowedResponse)

	router.Route("/api", func(r1 chi.Router) {
		r1.Post("/machines/operations", app.changeMachineState)
		r1.Route("/clusters", func(r2 chi.Router) {
			r2.Post("/", app.createCluster)
			r2.Route("/{cluster_id}", func(r3 chi.Router) {
				r3.Delete("/", app.deleteCluster)
				r3.Route("/vms", func(r4 chi.Router) {
					r4.Post("/", app.createMachine)
					r4.Delete("/{vm_id}", app.deleteMachine)
				})
			})
		})
	})

	return router
}
