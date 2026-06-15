package main

import (
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/Hayversong/questboard/internal/storage"
)

func main() {
	from := flag.String("from", storage.DataFilePath(), "caminho do JSON de origem")
	to := flag.String("to", storage.SQLiteFilePath(), "caminho do SQLite de destino")
	force := flag.Bool("force", false, "sobrescreve o banco SQLite se ele ja existir")

	flag.Parse()

	summary, err := storage.MigrateJSONToSQLite(storage.JSONToSQLiteMigration{
		JSONPath:   *from,
		SQLitePath: *to,
		Overwrite:  *force,
	})
	if err != nil {
		if errors.Is(err, storage.ErrSQLiteDatabaseExists) {
			log.Fatalf("banco SQLite ja existe em %s; use -force para sobrescrever", *to)
		}

		log.Fatal(err)
	}

	fmt.Printf(
		"Migracao concluida: %d projetos, %d cards, %d atividades -> %s\n",
		summary.Projects,
		summary.Cards,
		summary.Activities,
		*to,
	)
}
