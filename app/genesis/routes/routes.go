package routes
	
	import (
	"solar-system/genesis/handler"

	"github.com/go-chi/chi/v5"
)

func MountRoutes(r *chi.Mux, h *handler.Handler) {
	
		r.Get("/api/v1/example", h.GetExample)

}
	