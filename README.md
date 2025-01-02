# MyST (My SecreTs)

A secure command-line secret manager for storing and managing sensitive information like passwords, API keys, and other credentials.

## Features

- 🔐 **Secure Storage**: AES-GCM encryption with master passphrase
- 🔍 **Fast Search**: Full-text search through secrets
- 📋 **Clipboard Support**: Copy secret values directly to clipboard
- 🖥️ **Interactive CLI**: Arrow key navigation and command typing
- 🔒 **Local Storage**: All data stored locally in `~/.myst/`

## Installation

### From Source
```bash
# Clone the repository
git clone https://github.com/Isaac-Fate/myst.git
cd myst

# Build for your platform
make build

# Or build for all platforms
make build-all
```

### Using Go
```bash
go install github.com/Isaac-Fate/myst@latest
```

## Usage

### First Time Setup
```bash
myst
```
You'll be prompted to set a master passphrase (minimum 8 characters) which will be used to encrypt/decrypt your secrets.

### Available Commands

- `add`: Add a new secret
  ```
  Key: my-api-key
  Value: [hidden]
  Website (optional): api.example.com
  Notes (optional): Production API key
  ```

- `find`: Search secrets by key/website/notes
  - Options for each secret:
    - Display in terminal
    - Copy to clipboard
    - Skip

- `list`: View all secrets
  - Select any secret to:
    - Display value
    - Copy to clipboard
    - Skip

- `update`: Modify existing secrets
  - Update value/website/notes
  - Requires confirmation

- `remove`: Delete secrets
  - Select secret to remove
  - Requires confirmation

### Navigation

- ↑/↓: Navigate options
- Type: Quick command access
- Enter: Select option
- Ctrl+C: Cancel operation

## Development

### Requirements
- Go 1.21+
- Make

### Build Commands
```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Clean build artifacts
make clean
```

## License

MIT License
