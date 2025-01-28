# DotBack

DotBack is a command-line tool for backing up and restoring your dotfiles and application configurations using GitHub as storage.

## Features

- GitHub authentication via Personal Access Token (PAT)
- Secure token storage using system keyring
- Scan system for dotfiles and application configurations
- Interactive backup wizard
- Interactive restore wizard with symlink support
- Machine-specific configuration management

## Installation

```bash
go install github.com/amroessam/dotback@latest
```

## Usage

### Authentication

DotBack uses GitHub Personal Access Tokens (PATs) for authentication. The token should have the following scopes:
- `repo` (Full control of private repositories)
- `read:user` (Read all user profile data)
- `user:email` (Access user email addresses)

You can log in using one of these methods:

1. Using environment variable:
```bash
export GITHUB_TOKEN=your_token_here
dotback login
```

2. Interactive prompt:
```bash
dotback login
# You will be prompted to enter your token
```

To log out and remove the stored token:
```bash
dotback logout
```

### Scan for Configuration Files
```bash
dotback scan
dotback scan --verbose  # For detailed output
```

### Backup Your Configuration
```bash
dotback backup
```

### Restore Your Configuration
```bash
dotback restore
```

## Requirements

- Go 1.22 or later
- GitHub account
- GitHub Personal Access Token with appropriate permissions
- System keyring support (macOS Keychain, Linux Secret Service, or Windows Credential Manager)

## Security

DotBack takes security seriously:
- Your GitHub token is stored securely in your system's keyring
- No sensitive information is stored in plain text
- All GitHub communication is done over HTTPS
- Tokens are validated before being stored

## License

MIT License 