#!/bin/sh

echo "Setting up MinIO..."

# Install MinIO client
echo "Installing MinIO client..."
apk add --no-cache curl
wget https://dl.min.io/client/mc/release/linux-amd64/mc -O /usr/local/bin/mc
chmod +x /usr/local/bin/mc

# Wait for MinIO to be ready
echo "Waiting for MinIO to be ready..."
until mc alias set myminio http://minio:9000 ${MINIO_ROOT_USER:-minioadmin} ${MINIO_ROOT_PASSWORD:-minioadmin}; do
  echo "Waiting for MinIO to be ready..."
  sleep 2
done

# Create audora-tracks bucket for original files (pipeline input)
echo "Creating audora-tracks bucket..."
mc mb myminio/audora-tracks 2>/dev/null || echo "Bucket already exists"

# Create processed-tracks bucket for processed files (pipeline output, used for streaming)
echo "Creating processed-tracks bucket..."
mc mb myminio/processed-tracks 2>/dev/null || echo "Bucket already exists"

# Set public read policy for audora bucket
echo "Setting public read policy for audora bucket..."
mc anonymous set download myminio/audora

# Set public read policy for processed-tracks bucket (for streaming access)
echo "Setting public read policy for processed-tracks bucket..."
mc anonymous set download myminio/processed-tracks

echo "MinIO setup completed successfully!" 