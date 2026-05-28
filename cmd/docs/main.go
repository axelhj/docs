ackage main

import (
	"log"

	"github.com/yourusername/docsrv/internal/config"
	"github.com/yourusername/docsrv/internal/repository"
	"github.com/yourusername/docsrv/internal/server"
)

func main() {
	cfg := config.Load()

	db := repository.NewCouchDB(cfg.DB.Host, cfg.DB.User, cfg.DB.Pass)
	srv := server.NewDocServer(db)

	if err := srv.Start(cfg.Server.Port); err != nil {
		log.Fatal(err)
	}
}
