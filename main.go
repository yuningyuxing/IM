package main

import "main/router"

func main() {
	r := router.Router()

	r.Run(":9090")
}
