#!/bin/sh

echo "Setting up MinIO..."

# Wait for MinIO to be ready
until mc alias set myminio http://minio:9000 ${MINIO_ROOT_USER:-minioadmin} ${MINIO_ROOT_PASSWORD:-minioadmin}; do
  echo "Waiting for MinIO to be ready..."
  sleep 2
done

# Create audora bucket
echo "Creating audora bucket..."
mc mb myminio/audora

# Set public read policy for audora bucket
echo "Setting public read policy for audora bucket..."
mc policy set download myminio/audora

echo "MinIO setup completed successfully!" 