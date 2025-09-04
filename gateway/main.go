package main

import (
	"fmt"
	"gateway/configs"
)

func init() {
	configs.InitEnv()
	fmt.Println("👉 init() chạy trước 123456")
}

func main() {

	// routersInit := routes.InitRouter()
}
