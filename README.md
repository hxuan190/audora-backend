# Audora - Containerized Infrastructure

A containerized music platform built with Go, PostgreSQL, MinIO, and Ory Kratos.

## Architecture Overview

This project uses Docker Compose to orchestrate a complete development environment with the following services:

- **Audora API** (Go Backend): Core application logic
- **PostgreSQL (Audora)**: Main application database
- **PostgreSQL (Kratos)**: Identity service database  
- **MinIO**: Object storage for audio files and images
- **Ory Kratos**: Identity and authentication service

## Quick Start

1. **Clone and setup environment:**
   ```bash
   git clone <repository-url>
   cd aura
   cp env.example .env
   # Edit .env with your preferred values
   ```

2. **Start all services:**
   ```bash
   docker-compose up -d
   ```

3. **Access services:**
   - **Audora API**: http://localhost:8080
   - **MinIO Console**: http://localhost:9001 (admin/admin)
   - **Kratos Public API**: http://localhost:4433
   - **Kratos Admin API**: http://localhost:4434

## Service Details

### PostgreSQL (Audora Database)
- **Port**: 5432
- **Database**: audora
- **Purpose**: Stores artists, tracks, playlists, likes, and tips
- **Persistence**: Docker volume `audora_data`

### PostgreSQL (Kratos Database)
- **Port**: 5433
- **Database**: kratos
- **Purpose**: Stores user identities, sessions, and authentication data
- **Persistence**: Docker volume `kratos_data`

### MinIO Object Storage
- **Ports**: 9000 (API), 9001 (Console)
- **Bucket**: audora (auto-created)
- **Purpose**: Stores audio files and cover art images
- **Access**: Public read, authenticated write
- **Persistence**: Docker volume `minio_data`

### Ory Kratos Identity Service
- **Ports**: 4433 (Public), 4434 (Admin)
- **Purpose**: Handles user registration, login, and session management
- **Integration**: Webhook to Audora API after user registration
- **Features**: Email verification, password recovery, MFA support

## Configuration

### Environment Variables

Key environment variables in `.env`:

```bash
# Database Configuration
AUDORA_DB_USER=postgres
AUDORA_DB_PASSWORD=postgres
AUDORA_DB_NAME=audora

# MinIO Configuration  
MINIO_ROOT_USER=minioadmin
MINIO_ROOT_PASSWORD=minioadmin

# Service URLs
AUDORA_APP_URL=http://localhost:3000
AUDORA_API_URL=http://localhost:8080
```

### Kratos Configuration

Kratos is configured via files in the `kratos/` directory:

- `kratos.yml`: Main configuration file
- `identity.schema.json`: User data schema
- `after-registration.jsonnet`: Webhook payload template

## Database Schema

The Audora database includes these tables:

- **artists**: User profiles with Kratos integration
- **tracks**: Music tracks with metadata
- **user_likes**: Track likes by users
- **tips**: Monetary tips between users
- **playlists**: User-created playlists
- **playlist_tracks**: Playlist-track relationships

## Development Workflow

1. **Start services**: `docker-compose up -d`
2. **View logs**: `docker-compose logs -f <service-name>`
3. **Stop services**: `docker-compose down`
4. **Reset data**: `docker-compose down -v && docker-compose up -d`

## Security Notes

⚠️ **Development Only**: This setup uses default credentials and development settings. For production:

- Change all default passwords
- Use proper secrets management
- Enable HTTPS/TLS
- Configure proper CORS policies
- Use production-grade PostgreSQL and MinIO configurations

## Troubleshooting

### Service Health Checks
```bash
# Check service status
docker-compose ps

# View service logs
docker-compose logs audora-api
docker-compose logs kratos
docker-compose logs minio
```

### Database Access
```bash
# Connect to Audora database
docker-compose exec audora-db psql -U postgres -d audora

# Connect to Kratos database  
docker-compose exec kratos-db psql -U kratos -d kratos
```

### MinIO Access
- **Console**: http://localhost:9001
- **API Endpoint**: http://localhost:9000
- **Default credentials**: minioadmin/minioadmin
