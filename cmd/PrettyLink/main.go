package main

import (
	"PrettyLinkBackend/internal/config"
	"fmt"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg) //in prod we need to remove all these debug prints

	//TODO init config: cleanenv

	//TODO init logger: slog (import from logs/slog)

	//TODO init storage: sqlite

	//TODO init router: chi - it is compatible with net/http
	//TODO also chi/render

	//TODO run server
}
