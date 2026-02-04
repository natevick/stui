# S3 TUI

A terminal user interface for browsing and downloading files from AWS S3, with full support for AWS IAM Identity Center (SSO).

![S3 TUI Demo](https://via.placeholder.com/800x400?text=S3+TUI+Screenshot)

## Features

- **Browse S3 buckets and prefixes** - Navigate your S3 storage like a file browser
- **AWS SSO support** - Works with IAM Identity Center profiles
- **Profile picker** - Select from available AWS profiles on startup
- **Multi-select** - Select multiple files/folders with spacebar
- **Download files** - Download individual files or entire prefixes
- **Sync folders** - Sync S3 prefixes to local directories (only downloads changed files)
- **Bookmarks** - Save frequently accessed locations
- **Demo mode** - Try the UI without AWS credentials

## Installation

### From Binary

Download the latest release for your platform from the [Releases](https://github.com/natevick/s3-tui/releases) page.

### From Source

```bash
go install github.com/natevick/s3-tui/cmd/s3tui@latest
```

Or clone and build:

```bash
git clone https://github.com/natevick/s3-tui.git
cd s3-tui
go build -o s3tui ./cmd/s3tui
```

## Usage

```bash
# Launch with profile picker
s3tui

# Launch with specific profile
s3tui --profile my-profile

# Launch directly into a bucket
s3tui --profile my-profile --bucket my-bucket

# Demo mode (no AWS credentials needed)
s3tui --demo
```

### AWS SSO Login

Before using with SSO profiles, authenticate:

```bash
aws sso login --profile my-profile
```

## Keyboard Shortcuts

### Navigation
| Key | Action |
|-----|--------|
| `↑/k`, `↓/j` | Move up/down |
| `Enter` | Open folder / Select |
| `Backspace` | Go back |
| `PgUp/PgDn` | Page up/down |

### Views
| Key | Action |
|-----|--------|
| `←/→` | Switch tabs |
| `Tab` | Next tab |
| `Shift+Tab` | Previous tab |
| `1/2/3` | Jump to tab |

### Actions
| Key | Action |
|-----|--------|
| `Space` | Select/deselect item |
| `d` | Download selected |
| `s` | Sync prefix to local |
| `b` | Add bookmark |
| `r` | Refresh |
| `/` | Filter list |

### General
| Key | Action |
|-----|--------|
| `?` | Toggle help |
| `Esc` | Cancel / Close |
| `q` | Quit |

## Configuration

S3 TUI uses your standard AWS configuration (`~/.aws/config` and `~/.aws/credentials`).

### Example SSO Profile

```ini
[sso-session my-sso]
sso_start_url = https://my-company.awsapps.com/start
sso_region = us-east-1
sso_registration_scopes = sso:account:access

[profile my-profile]
sso_session = my-sso
sso_account_id = 123456789012
sso_role_name = MyRole
region = us-west-2
```

## License

MIT License - see [LICENSE](LICENSE) for details.
