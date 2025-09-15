# JitsiBot - Audio Recording Bot for Jitsi Meet

JitsiBot is a Go-based automation tool that joins Jitsi Meet conferences and records audio from all participants. It uses Chrome automation via chromedp to interact with the Jitsi Meet web interface and captures audio streams from remote participants.

## Features

- 🎤 **Automatic audio recording** of all conference participants
- 🔍 **Real-time monitoring** of participant join/leave events
- 💾 **Automatic file organization** by user and room
- 🎯 **WebM format** with Opus codec for efficient storage
- 🔧 **Configurable** through command-line options
- 🛡️ **Headless Chrome automation** for reliable operation

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

### Dependencies

The project uses the following Go dependencies:
- `github.com/chromedp/chromedp` - Chrome automation
- `github.com/chromedp/cdproto` - Chrome DevTools Protocol

## Usage

### Basic Usage

```bash
./jitsibot -room "my-conference-room" -botname "RecordingBot"
```

### Command Line Options

| Option | Description | Default |
|--------|-------------|---------|
| `-room` | Conference room name | `ssjitsi-test` |
| `-botname` | Bot display name in conference | `SSJitsiBot` |
| `-datadir` | Directory for saving recordings | `../data/` |
| `-jitsi` | Jitsi server URL | `https://meet.jit.si/` |
| `-username` | Username for authentication | (empty) |
| `-pass` | Password for authentication | (empty) |
| `-help` | Show help information | `false` |

### Examples

**Join public room:**
```bash
./jitsibot -room "team-meeting" -botname "TeamRecorder" -datadir "./recordings"
```

**Join password-protected room:**
```bash
./jitsibot -room "private-meeting" -botname "SecureRecorder" -username "user" -pass "password"
```

**Use custom Jitsi server:**
```bash
./jitsibot -room "conference" -jitsi "https://my-jitsi-server.com/" -datadir "/data/recordings"
```

## Data Storage Format

Recordings are organized in the following structure:

```
data/
└── {bot-user-id}/
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
