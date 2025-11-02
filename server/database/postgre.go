package database

import (
	"context"
	"log"
	"time"

	"task-manager/ent"

	_ "github.com/lib/pq"
)

func ConnectPostgre(dsn string, maxRetries int, delay time.Duration) (*ent.Client, error) {
	var client *ent.Client
	var err error

	for i := 0; i < maxRetries; i++ {
		client, err = ent.Open("postgres", dsn)
		if err == nil {

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if errPing := client.Schema.Create(ctx); errPing == nil {
				log.Println("Conexión a base de datos establecida")
				return client, nil
			} else {
				log.Printf("Ping/schema fallido en intento %d/%d: %v", i+1, maxRetries, errPing)
				client.Close()
				err = errPing
			}
		} else {
			log.Printf("Error al abrir la base de datos en intento %d/%d: %v", i+1, maxRetries, err)
		}

		time.Sleep(delay)
	}

	return nil, err
}
