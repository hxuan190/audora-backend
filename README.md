# Audora

- **Set environment variables first**
- ```export GOOSE_DRIVER=postgres```
- ```export GOOSE_DBSTRING="postgres://postgres:postgres@localhost:5432/audora?sslmode=disable"```

- **Now you can run shorter commands
- ```goose -dir ./migrations status```
- ```goose -dir ./migrations up```
- ```goose -dir ./migrations down```
- ```goose -dir ./migrations reset```

- **Audora API**: http://localhost:8080
- **MinIO Console**: http://localhost:9001 (admin/admin)
- **Kratos Public API**: http://localhost:4433
- **Kratos Admin API**: http://localhost:4434