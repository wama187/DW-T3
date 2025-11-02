package main

import (
	"log"
	"net/http"
	"os"
	"time"

    "task-manager/graph"
    "task-manager/shared/middleware"

	"task-manager/apps/projects"
	"task-manager/apps/tasks"
	"task-manager/apps/users"
	"task-manager/database"

	"github.com/rs/cors"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8080"


func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        log.Printf("Started %s %s", r.Method, r.URL.Path)
        
        next.ServeHTTP(w, r)
        
        log.Printf("Completed %s in %v", r.URL.Path, time.Since(start))
    })
}

func main() {
    port := defaultPort

    secretKey := "mZKk187ABBo911QIUYE989"

    dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "host=localhost port=5432 user=postgres password=postgres dbname=taskmanager sslmode=disable"
	}
    
    client, err := database.ConnectPostgre(dbPath, 10, 2*time.Second)

    if err != nil {
        log.Fatalf("Fue imposible conectar a la base de datos:%v", err)
    }

    defer client.Close()
    
    log.Println("Conexión a base de datos establecida")
    
    // Middleware de autenticación
    authMiddleware := middleware.NewAuthMiddleware(secretKey)
    
    // GraphQL server
    srv := handler.New(graph.NewExecutableSchema(graph.Config{
        Resolvers: graph.NewResolver(
			users.NewUserService(client, secretKey),
			projects.NewProjectService(client),
			tasks.NewTaskService(client),
		),
    }))
    
    srv.AddTransport(transport.Options{})
    srv.AddTransport(transport.GET{})
    srv.AddTransport(transport.POST{})
    srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))
    srv.Use(extension.Introspection{})
    srv.Use(extension.AutomaticPersistedQuery{
        Cache: lru.New[string](100),
    })
    
    // Router
    mux := http.NewServeMux()
    
    // Health check endpoint
    mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"status":"ok"}`))
    })
    
    // Playground (sin autenticación)
    mux.Handle("/", playground.Handler("GraphQL playground", "/graphql"))
    
    // GraphQL endpoint (con autenticación)
    mux.Handle("/graphql", authMiddleware.Middleware(srv))
    
    // CORS
    corsHandler := cors.New(cors.Options{
        AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:5173", "http://client:3000"},
        AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
        AllowedHeaders:   []string{"Authorization", "Content-Type"},
        AllowCredentials: true,
        Debug:            false,
    })
    
    handler := loggingMiddleware(corsHandler.Handler(mux))
    
    log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
    log.Printf("Server ready at http://localhost:%s/", port)
    log.Printf("GraphQL playground at http://localhost:%s/", port)
    log.Printf("Health check at http://localhost:%s/health", port)
    log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
    
    log.Fatal(http.ListenAndServe(":"+port, handler))
}