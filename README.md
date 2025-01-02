# MyST (My SecreTs)

MyST is a secure command-line secret manager that helps you store and manage sensitive information like passwords, API keys, and other credentials.

## Features

- 🔐 **Secure Storage**: All secrets are encrypted using AES-GCM with a master passphrase
- 🔍 **Fast Search**: Full-text search capabilities for finding secrets quickly
- 📋 **Clipboard Integration**: Copy secret values directly to clipboard
- 🖥️ **Interactive CLI**: User-friendly interface with both arrow key navigation and command typing
- 🔒 **Local Storage**: All data is stored locally in `~/.myst/`

## Installation

```sh
go install github.com/Isaac-Fate/myst@latest
```

## Quick Start

1. Run MyST:

```sh
myst
```

2. First-time setup:
   - Set your master passphrase (min 8 characters)
   - This passphrase will encrypt/decrypt your secrets

3. Use the interactive menu to manage your secrets

## Commands

- `add`: Add a new secret
  ```
  Key: github-token
  Value: [hidden]
  Website (optional): github.com
  Notes (optional): Personal access token
  ```

- `find`: Search secrets
  - Search by key, website, or notes
  - For each secret:
    - Display in terminal
    - Copy to clipboard
    - Skip

- `list`: View all secrets
  - Shows all stored secrets
  - Select any to view/copy value

- `update`: Modify secrets
  - Select a secret
  - Update:
    - Value
    - Website
    - Notes

- `remove`: Delete secrets
  - Select a secret
  - Confirm deletion

- `help`: Show help

- `quit`: Exit

## Navigation

- Use ↑/↓ arrows to navigate
- Type commands directly
- Enter to select
- Ctrl+C to cancel

## Security

- AES-GCM encryption
- Master passphrase never stored
- Local SQLite database
- Separate search index

## Development

### Requirements

- Go 1.21+
- SQLite3

### Build

```sh
git clone https://github.com/Isaac-Fate/myst.git
cd myst
go build
```
