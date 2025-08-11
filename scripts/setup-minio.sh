#!/bin/sh

echo "🎵 Setting up MinIO for Audora..."

# Install MinIO client
echo "📦 Installing MinIO client..."
apk add --no-cache curl wget
wget https://dl.min.io/client/mc/release/linux-amd64/mc -O /usr/local/bin/mc
chmod +x /usr/local/bin/mc

# Wait for MinIO to be ready
echo "⏳ Waiting for MinIO to be ready..."
until mc alias set audora-minio http://minio:9000 ${MINIO_ROOT_USER:-minioadmin} ${MINIO_ROOT_PASSWORD:-minioadmin}; do
  echo "Waiting for MinIO server..."
  sleep 2
done

echo "✅ MinIO server is ready!"

# Create buckets
echo "🪣 Creating buckets..."

# Main audora bucket (general uploads, user content)
mc mb audora-minio/audora 2>/dev/null || echo "  ℹ️  Bucket 'audora' already exists"

# Original tracks bucket (pipeline input - artist uploads)
mc mb audora-minio/audora-tracks 2>/dev/null || echo "  ℹ️  Bucket 'audora-tracks' already exists"

# Processed tracks bucket (pipeline output - streaming files)
mc mb audora-minio/processed-tracks 2>/dev/null || echo "  ℹ️  Bucket 'processed-tracks' already exists"

# Create pipeline service user
echo "👤 Creating pipeline service user..."
mc admin user add audora-minio pipeline-user pipeline-secret-key 2>/dev/null || echo "  ℹ️  User 'pipeline-user' already exists"

# Apply custom policy if available
if [ -f "/minio-policy.json" ]; then
    echo "📋 Applying custom Audora policy..."
    
    # Create the policy from JSON file
    mc admin policy create audora-minio audora-policy /minio-policy.json 2>/dev/null || echo "  ℹ️  Policy 'audora-policy' already exists"
    
    # Attach policy to pipeline user
    mc admin policy attach audora-minio audora-policy --user pipeline-user 2>/dev/null || echo "  ℹ️  Policy already attached to user"
    
    echo "✅ Custom policy applied successfully!"
else
    echo "⚠️  Policy file not found at /minio-policy.json, using default policies..."
fi

# Set bucket access policies for different use cases
echo "🔒 Setting bucket access policies..."

# For audora bucket: Allow public read for public assets only
echo "  📁 Setting audora bucket policies..."
mc anonymous set download audora-minio/audora/public/ 2>/dev/null || echo "    ℹ️  Public read policy already set for audora/public/"

# For audora-tracks: Private bucket for original uploads (artists only)
echo "  📁 Setting audora-tracks bucket policies..."
mc anonymous set none audora-minio/audora-tracks 2>/dev/null || echo "    ℹ️  Private policy already set for audora-tracks"

# For processed-tracks: Allow signed URL access for streaming
echo "  📁 Setting processed-tracks bucket policies..."
mc anonymous set none audora-minio/processed-tracks 2>/dev/null || echo "    ℹ️  Private policy already set for processed-tracks"

# Set up bucket versioning for important data
echo "🔄 Enabling versioning..."
mc version enable audora-minio/audora-tracks 2>/dev/null || echo "  ℹ️  Versioning already enabled for audora-tracks"

# Set up lifecycle policies for cost optimization
echo "⏰ Setting up lifecycle policies..."
cat > /tmp/lifecycle-policy.json << 'EOF'
{
    "Rules": [
        {
            "ID": "ProcessedTracksTransition",
            "Status": "Enabled",
            "Filter": {
                "Prefix": ""
            },
            "Transition": {
                "Days": 30,
                "StorageClass": "REDUCED_REDUNDANCY"
            }
        },
        {
            "ID": "OldVersionCleanup",
            "Status": "Enabled",
            "Filter": {
                "Prefix": ""
            },
            "NoncurrentVersionExpiration": {
                "NoncurrentDays": 90
            }
        }
    ]
}
EOF

mc ilm import audora-minio/processed-tracks < /tmp/lifecycle-policy.json 2>/dev/null || echo "  ℹ️  Lifecycle policy already applied"

# Create some test folders structure
echo "📂 Creating folder structure..."
mc cp /dev/null audora-minio/audora/public/.keep 2>/dev/null || true
mc cp /dev/null audora-minio/audora/avatars/.keep 2>/dev/null || true
mc cp /dev/null audora-minio/audora/artwork/.keep 2>/dev/null || true
mc cp /dev/null audora-minio/audora-tracks/uploads/.keep 2>/dev/null || true
mc cp /dev/null audora-minio/processed-tracks/mp3/.keep 2>/dev/null || true
mc cp /dev/null audora-minio/processed-tracks/flac/.keep 2>/dev/null || true
mc cp /dev/null audora-minio/processed-tracks/hires/.keep 2>/dev/null || true

# Display bucket information
echo "📊 Bucket information:"
echo "================================"
mc ls audora-minio/
echo "================================"

# Display policy information
echo "📋 Active policies:"
echo "================================"
mc admin policy list audora-minio
echo "================================"

# Display user information
echo "👥 Service users:"
echo "================================"
mc admin user list audora-minio
echo "================================"

echo "🎉 MinIO setup completed successfully!"
echo ""
echo "📍 Access Information:"
echo "  🌐 MinIO Console: http://localhost:9001"
echo "  🔑 Username: ${MINIO_ROOT_USER:-minioadmin}"
echo "  🔑 Password: ${MINIO_ROOT_PASSWORD:-minioadmin}"
echo ""
echo "📦 Buckets created:"
echo "  🎵 audora - General uploads and public assets"
echo "  🎤 audora-tracks - Original artist uploads (private)"
echo "  🎧 processed-tracks - Processed audio files for streaming (private)"
echo ""
echo "👤 Service Account:"
echo "  👤 User: pipeline-user"
echo "  🔑 Key: pipeline-secret-key"
echo "  📋 Policy: audora-policy (custom bucket access)"