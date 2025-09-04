package main

import (
	"gateway/configs"
)

func init() {
	configs.InitEnv()
	configs.ConnectDatabase()
}

func main() {

	// routersInit := routes.InitRouter()
}
