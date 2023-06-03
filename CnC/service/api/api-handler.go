package api

import (
	"net/http"
)

func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.Welcome)
	rt.router.POST("/zombie", rt.Zombie)
	rt.router.POST("/action", rt.Action)
	rt.router.PUT("/:id/delete", rt.Delete)

	return rt.router
}
