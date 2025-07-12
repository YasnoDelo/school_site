package server

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations подключается к DB и прогоняет все up‑миграции из db/migrations.
func RunMigrations(db *sql.DB) {
	// 1) Создаём драйвер для postgres
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("migrate: cannot create DB driver: %v", err)
	}

	// 2) Указываем source файловых миграций и target‑БД
	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres", driver,
	)
	if err != nil {
		log.Fatalf("migrate: initialization failed: %v", err)
	}

	// 3) Применяем миграции «вперёд»
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migrate: up failed: %v", err)
	}

	log.Printf("migrate: completed (or no change needed)")
}
