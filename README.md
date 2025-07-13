# Audora

# Set environment variables first
export GOOSE_DRIVER=postgres
export GOOSE_DBSTRING="postgres://postgres:postgres@localhost:5432/audora?sslmode=disable"

# Now you can run shorter commands
goose -dir ./migrations status
goose -dir ./migrations up
goose -dir ./migrations down
goose -dir ./migrations reset