// Code generated by hertz generator. DO NOT EDIT.

package main

import (
	router "github.com/Linda-ui/orbital_HeBao/hertz_gateway/router"
	"github.com/cloudwego/hertz/pkg/app/server"
)

// register registers all routers.
func register(r *server.Hertz) {

	router.GeneratedRegister(r)

	customizedRegister(r)
}
