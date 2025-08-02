# Audora

- **Set environment variables first**
- ```export GOOSE_DRIVER=postgres```
- ```export GOOSE_DBSTRING="postgres://postgres:postgres@localhost:5432/audora?sslmode=disable"```

- **Now you can run shorter commands**
- ```goose -dir ./migrations status```
- ```goose -dir ./migrations up```
- ```goose -dir ./migrations down```
- ```goose -dir ./migrations reset```

- **Rebuild docker image**
- ```docker build -t audora-api .```
- ```
  docker run -d \
  --name audora-api \
  -p 8080:8080 \
  -v $(pwd):/app \
  --network music-app-backend_audora-network \ 
  -e APP_ENV=development \
  -e APP_PORT=8080 \
  -e AUDORA_DB_HOST=audora-db \
  -e AUDORA_DB_PORT=5432 \
  -e AUDORA_DB_USER=postgres \
  -e AUDORA_DB_PASSWORD=postgres \
  -e AUDORA_DB_NAME=audora \
  -e MINIO_ENDPOINT=minio:9000 \
  -e MINIO_ACCESS_KEY=minioadmin \
  -e MINIO_SECRET_KEY=minioadmin \
  -e MINIO_BUCKET_NAME=audora \
  -e KRATOS_PUBLIC_URL=http://kratos:4433 \
  -e KRATOS_ADMIN_URL=http://kratos:4434 \
  -e JWT_SECRET=your-jwt-secret-here \
  audora-api
  ```

- **Audora API**: http://localhost:8080
- **MinIO Console**: http://localhost:9001 (admin/admin)
- **Kratos Public API**: http://localhost:4433
- **Kratos Admin API**: http://localhost:4434