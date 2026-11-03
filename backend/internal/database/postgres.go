package database

impost (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(databaseURL string) *pgxpool.Pool {

	config, err := pgxpool.ParseConfig(databaseURL)

	if err != nil {
		log.Fatal("Unable to parse database configuration:", err)
	}

	config.MaxConns = 10
	config.Minconns = 2
	config.MaxConnLifetime = time.Hour

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)

	defer cancel()

	pool, err := pgxpool.NewWithconfig(ctx, config)

	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		log.Fatal("Database ping failed:", err)
	}

	log.Println("Connected to database successfully")

	return pool
}