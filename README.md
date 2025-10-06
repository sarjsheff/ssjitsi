# JitsiBot Server - Audio Recording Bot for Jitsi Meet

JitsiBot Server is a Go-based automation tool that manages multiple bots to join Jitsi Meet conferences and record audio from all participants. It uses Chrome automation via chromedp to interact with the Jitsi Meet web interface and captures audio streams from remote participants. The system includes a web interface for real-time monitoring of bot status and screenshots.

## Features

- 🎤 **Automatic audio recording** of all conference participants
- 🔍 **Real-time monitoring** of participant join/leave events
- 💾 **Automatic file organization** by user and room
- 🎯 **WebM format** with Opus codec for efficient storage
- 🔧 **Configurable** through YAML configuration
- 🛡️ **Headless Chrome automation** for reliable operation
- 🔐 **Flexible authentication** - supports both username/password and JWT
- 🌐 **Web interface** for real-time bot monitoring and screenshots
- 🚀 **Multiple bot management** for simultaneous recordings

## How It Works

### 1. Browser Automation
The Go application uses `chromedp` to automate Chrome browser:
- Navigates to Jitsi Meet server
- Joins specified conference room
- Handles authentication (username/password or JWT)
- Injects custom JavaScript for audio recording

### 2. Audio Capture
The injected JavaScript (`script.js`):
- Monitors DOM for audio elements with `remoteAudio_` prefix
- Creates custom Web Components for each participant
- Uses `MediaRecorder` API to capture audio streams
- Converts audio chunks to base64 and sends to Go backend

### 3. Data Storage
The Go backend:
- Receives audio data via JavaScript binding
- Decodes base64 audio chunks
- Organizes files by user ID and room
- Saves WebM audio files with metadata

## Installation

### Prerequisites

1. **Go** (version 1.24.4 or later)
2. **Chrome/Chromium** browser
3. **Node.js** and **npm** (for building the web interface)

### Dependencies

The project uses the following key Go dependencies:
- `github.com/chromedp/chromedp` - Chrome automation
- `github.com/chromedp/cdproto` - Chrome DevTools Protocol
- `github.com/golang-jwt/jwt/v5` - JWT token generation and signing
- `github.com/gin-gonic/gin` - HTTP web framework
- `gopkg.in/yaml.v2` - YAML configuration parsing

### Building from Source

```bash
# Build the complete application (includes web UI)
make all

# Or build step by step:
make docs      # Generate API documentation
make build-ui  # Build React web interface
make ssjitsi   # Build the main server binary

# Clean build artifacts
make clean
```

The build process will create a `ssjitsi` binary that includes the embedded web interface.

## Usage

### Server Mode (Multiple Bots)

```bash
./ssjitsi -config ssjitsi.yaml
```

The server includes a built-in web interface available at: http://localhost:8080/

### Command Line Options

#### Server (Multiple Bots)

| Option | Description | Default |
|--------|-------------|---------|
| `-config` | Path to configuration file | `ssjitsi.yaml` |
| `-help` | Show help information | `false` |

### Examples

**Run server with multiple bots:**
```bash
./ssjitsi -config ssjitsi.yaml
```

**Access web interface:**
Open http://localhost:8080/ in your browser to monitor bot status and view screenshots.

### Configuration File Format

The server uses a YAML configuration file to manage multiple bots. The bot supports two authentication methods:

#### Traditional Authentication (Username/Password)

For Jitsi Meet servers with basic authentication:

```yaml
http: ":8080"

bots:
  - Room: "conference-room-1"
    BotName: "Recording Bot 1"
    DataDir: "./data"
    JitsiServer: "https://meet.jit.si/"
    Username: "user"
    Pass: "password"
    Headless: true
```

#### JWT Authentication (Recommended for self-hosted Jitsi)

For Jitsi Meet servers with JWT token authentication:

```yaml
http: ":8080"

bots:
  - Room: "my-room"
    BotName: "JWT Bot"
    DataDir: "./data"
    JitsiServer: "https://jitsi.example.com"
    JWTAppID: "your_app_id"
    JWTAppSecret: "your_app_secret"
    Headless: true
```

**JWT Authentication Details:**
- Uses JSON Web Tokens (JWT) for secure authentication
- Tokens are generated automatically using HS256 signing algorithm
- Token includes claims: `iss`, `aud`, `sub`, `room`, `exp`, `context`
- Token validity: 2 hours
- Bot navigates directly to `{JitsiServer}/{Room}?jwt={token}`

**Configuration Fields:**

| Field | Required | Description |
|-------|----------|-------------|
| `Room` | Yes | Name of the conference room |
| `BotName` | Yes | Display name for the bot |
| `DataDir` | Yes | Directory to store recordings |
| `JitsiServer` | Yes | URL of the Jitsi Meet server |
| `Username` | No | Username for basic auth |
| `Pass` | No | Password for basic auth |
| `JWTAppID` | No | JWT application ID |
| `JWTAppSecret` | No | JWT secret key for signing |
| `Headless` | Yes | Run in headless mode (true/false) |

**Note:** If both JWT credentials and Username/Pass are provided, JWT takes precedence.

## Web Interface

The server includes a built-in React web application for monitoring bot status:

### Features
- 📊 **Real-time bot list** with connection status
- 🖼️ **Live screenshots** from each bot's browser
- 🔄 **Auto-refresh** (bots: 10s, screenshots: 30s)
- 📱 **Responsive design** with Bootstrap
- ⚡ **Error handling** and loading states

### Access
After starting the server, open http://localhost:8080/ in your browser.

### Architecture
- **Frontend**: React 19.1.1 with Bootstrap 5
- **Backend**: Go server with Gin framework
- **Static files**: Embedded using `go:embed`
- **API**: REST endpoints for bot data and screenshots

## Data Storage Format

Recordings are organized in the following structure:

```
data/
└── {room-name}/                    # Room directory (safe filename)
    └── {bot-session-id}/           # Bot session directory
        ├── {participant-user-id}_{audio-element-id}.webm     # Audio recordings
        ├── {participant-user-id}_{audio-element-id}.json     # Start timestamp
        ├── {participant-user-id}.json                        # Participant display name
        └── room.json                                         # Room name
```

### File Types

1. **`.webm`** - Audio recordings in WebM format with Opus codec
2. **`.json`** - Metadata files containing:
   - Start timestamps (Unix milliseconds)
   - Participant display names
   - Room information

### Directory Structure Details

- **`{room-name}/`** - Directory named after the room (sanitized for filesystem safety)
- **`{bot-session-id}/`** - Unique directory for each bot session
- **Audio files** - Named with participant user ID and audio element ID
- **Metadata files** - JSON files with timestamps and participant information

## JavaScript Components

### Custom Web Components

- **`ssbot-info`**: Status panel showing bot connection information
- **`ssbot-audio`**: Audio recording component for each participant

### Key Functions

- **`observeElements()`**: Monitors DOM for audio element changes
- **`handleElementAppeared()`**: Starts recording when participant joins
- **`handleElementDisappeared()`**: Stops recording when participant leaves
- **`window.ssbot_writeSound`**: Binding to send audio data to Go backend

## Troubleshooting

### Common Issues

1. **Chrome not found**: Ensure Chrome/Chromium is installed
2. **Permission denied**: Run with appropriate permissions for data directory
3. **Authentication failed**:
   - For username/password: Verify credentials are correct
   - For JWT: Check that `JWTAppID` and `JWTAppSecret` match your Jitsi configuration
   - Verify Jitsi server JWT authentication is properly configured
4. **Audio not recording**: Check browser console for errors
5. **Bot joins but doesn't record**:
   - Check that other participants have audio enabled
   - Verify `script.js` is being injected properly
   - Monitor Go logs for binding events
6. **JWT token errors**:
   - Ensure JWT secret matches the one configured on Jitsi server
   - Check token hasn't expired (default: 2 hours validity)
   - Verify room name matches the token claim

### Debug Mode

Set `Headless: false` in the configuration file to see the browser in action. This allows you to:
- Visually confirm the bot joins the conference
- Check browser console for JavaScript errors
- Monitor network requests
- Verify authentication flow

### Checking Logs

The bot outputs useful information:
```bash
2025/10/06 18:09:46 Используем JWT авторизацию
2025/10/06 18:09:46 Переходим на URL: https://test-jitsi.aisa.ru/test_aiplan?jwt=***
2025/10/06 18:09:50 Бот 1 успешно запущен (ID: bc4aa902-1172-4fd8-800a-91ac42e31592)
```

Expected console messages:
- "Failed to create local tracks" - Normal, bot has no microphone
- "Video track creation failed" - Normal, bot has no camera
- "APP.conference not ready yet, will retry..." - Normal during initialization

## JWT Authentication Setup

### For Self-Hosted Jitsi Meet Servers

To use JWT authentication, your Jitsi Meet server must be configured with JWT support. Here's a brief overview:

#### Server Requirements

1. **Prosody Configuration** (`/etc/prosody/conf.avail/your-domain.cfg.lua`):
   ```lua
   VirtualHost "your-domain.com"
       authentication = "token"
       app_id = "your_app_id"
       app_secret = "your_app_secret"
   ```

2. **Install JWT Plugin**:
   ```bash
   apt-get install lua-basexx lua-cjson luarocks
   luarocks install lua-sec
   ```

3. **Configure Jitsi Meet** (`/etc/jitsi/meet/your-domain-config.js`):
   ```javascript
   var config = {
       hosts: {
           domain: 'your-domain.com',
           muc: 'conference.your-domain.com'
       },
       // ... other settings
   };
   ```

#### Bot Configuration

Once your Jitsi server is configured for JWT, use these settings in `ssjitsi.yaml`:

```yaml
bots:
  - Room: "test-room"
    BotName: "Recording Bot"
    DataDir: ./data
    JitsiServer: https://your-domain.com
    JWTAppID: your_app_id          # Must match Prosody app_id
    JWTAppSecret: your_app_secret  # Must match Prosody app_secret
    Headless: true
```

The bot will automatically:
1. Generate a JWT token with the correct claims
2. Sign it with the provided secret
3. Navigate to the room URL with the token: `https://your-domain.com/room?jwt=<token>`
4. Join the conference without additional authentication prompts

For more information on Jitsi JWT setup, see: https://jitsi.github.io/handbook/docs/devops-guide/secure-domain
