# JitsiBot Server - Audio Recording Bot for Jitsi Meet

JitsiBot Server is a Go-based automation tool that manages multiple bots to join Jitsi Meet conferences and record audio from all participants. It uses Chrome automation via chromedp to interact with the Jitsi Meet web interface and captures audio streams from remote participants. The system includes a web interface for real-time monitoring of bot status and screenshots.

## Features

- 🎤 **Automatic audio recording** of all conference participants
- 🔍 **Real-time monitoring** of participant join/leave events
- 💾 **Automatic file organization** by user and room
- 🎯 **WebM format** with Opus codec for efficient storage
- 🔧 **Configurable** through YAML configuration
- 🛡️ **Headless Chrome automation** for reliable operation
- 🌐 **Web interface** for real-time bot monitoring and screenshots
- 🚀 **Multiple bot management** for simultaneous recordings

## How It Works

### 1. Browser Automation
The Go application uses `chromedp` to automate Chrome browser:
- Navigates to Jitsi Meet server
- Joins specified conference room
- Handles authentication if required
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

The project uses the following Go dependencies:
- `github.com/chromedp/chromedp` - Chrome automation
- `github.com/chromedp/cdproto` - Chrome DevTools Protocol

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

The server uses a YAML configuration file to manage multiple bots:

```yaml
http: ":8080"

bots:
  - Room: "conference-room-1"
    BotName: "Recording Bot 1"
    DataDir: "./data"
    JitsiServer: "https://meet.jit.si/"
    Username: ""
    Pass: ""
    Headless: true

  - Room: "conference-room-2"
    BotName: "Recording Bot 2"
    DataDir: "./data"
    JitsiServer: "https://meet.jit.si/"
    Username: "user"
    Pass: "password"
    Headless: true
```

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
3. **Authentication failed**: Verify username/password for protected rooms
4. **Audio not recording**: Check browser console for errors

### Debug Mode

For debugging, you can modify the script to run in non-headless mode (already configured) and check browser console for JavaScript errors.
