package main

import (
	"log"
	"vfinance-api/internal/config"
	"vfinance-api/internal/database"
	"vfinance-api/internal/server"
)

func main() {
	// Carregar configurações
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Erro ao carregar configurações:", err)
	}

	// Conectar ao banco de dados
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Erro ao conectar ao banco:", err)
	}

	// Executar migrações
	if err := database.Migrate(db); err != nil {
		log.Fatal("Erro nas migrações:", err)
	}

	// Iniciar servidor
	srv := server.New(cfg, db)
	log.Printf("Servidor iniciado na porta %s", cfg.Port)
	if err := srv.Start(); err != nil {
		log.Fatal("Erro ao iniciar servidor:", err)
	}
}