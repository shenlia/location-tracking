# Location Tracking Shortlink

A short URL generation service with GPS tracking capabilities. Generate short links to track visitor locations, devices, and browsing duration.

## Features

- **Short URL Generation**: Convert any article URL into a trackable short link
- **GPS Location Tracking**: Obtain precise GPS coordinates via诱导授权 (induced authorization)
- **IP Geolocation**: Fallback to IP-based city-level location
- **Device Tracking**: Track OS, browser, and device type
- **Duration Tracking**: Monitor page停留时长 (visit duration)
- **Admin Dashboard**: View detailed visit records and statistics

## Quick Start

```bash
# Start the server
./start.sh

# Or manually
mkdir -p data
./track
```

Server will start at `http://localhost:8080`

## API Usage

### Create Shortlink

```bash
curl -X POST http://localhost:8080/api/shortlinks/create \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com/article/123"}'
```

Response:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "short_url": "http://localhost:8080/abc123",
    "code": "abc123",
    "original_url": "https://example.com/article/123"
  }
}
```

### Access Shortlink

Visit `http://localhost:8080/{code}` to access the original URL and trigger tracking.

## Admin Dashboard

- Main panel: `http://localhost:8080/admin`
- Statistics: `http://localhost:8080/stats`

Default credentials (configure in config.yaml):
- Username: admin
- Password: changeme123

## Configuration

Edit `config.yaml`:

```yaml
server:
  host: "0.0.0.0"
  port: 8080

database:
  path: "./data/track.db"

admin:
  username: "admin"
  password: "changeme123"
```

## How Tracking Works

1. User creates a short link via API or admin panel
2. When someone visits the short link:
   - An诱导页面 (inducement page) appears with "Click to view details"
   - Clicking triggers GPS authorization request
   - If user grants permission, precise GPS coordinates are captured
   - If denied, IP-based location is used as fallback
3. Visit data (IP, location, device, duration) is recorded
4. View all data in the admin dashboard

## Tech Stack

- Backend: Go + Gin
- Database: SQLite
- Frontend: HTML5 + Vanilla JS
- Maps: Leaflet.js + OpenStreetMap
