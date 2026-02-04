package aws

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Client wraps the AWS S3 client with configuration
type Client struct {
	S3      *s3.Client
	Config  aws.Config
	Profile string
	Region  string
}

// NewClient creates a new AWS client with the specified profile
// Supports SSO profiles - user must run `aws sso login --profile <profile>` first
func NewClient(ctx context.Context, profile, region string) (*Client, error) {
	var opts []func(*config.LoadOptions) error

	if profile != "" {
		opts = append(opts, config.WithSharedConfigProfile(profile))
	}

	if region != "" {
		opts = append(opts, config.WithRegion(region))
	}

	cfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	s3Client := s3.NewFromConfig(cfg)

	return &Client{
		S3:      s3Client,
		Config:  cfg,
		Profile: profile,
		Region:  cfg.Region,
	}, nil
}

// WithRegion creates a new client with a different region
func (c *Client) WithRegion(ctx context.Context, region string) (*Client, error) {
	return NewClient(ctx, c.Profile, region)
}

// ProfileInfo contains information about an AWS profile
type ProfileInfo struct {
	Name       string
	Region     string
	SSOSession string
	AccountID  string
}

// ListProfiles returns a list of available AWS profiles from ~/.aws/config
func ListProfiles() ([]ProfileInfo, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, ".aws", "config")
	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open AWS config: %w", err)
	}
	defer file.Close()

	var profiles []ProfileInfo
	var currentProfile *ProfileInfo

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Check for section header
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			// Save previous profile if it exists and has SSO config
			if currentProfile != nil && currentProfile.SSOSession != "" {
				profiles = append(profiles, *currentProfile)
			}

			section := strings.TrimPrefix(strings.TrimSuffix(line, "]"), "[")

			// Skip sso-session sections, only get profiles
			if strings.HasPrefix(section, "sso-session ") {
				currentProfile = nil
				continue
			}

			// Extract profile name
			name := section
			if strings.HasPrefix(section, "profile ") {
				name = strings.TrimPrefix(section, "profile ")
			}

			currentProfile = &ProfileInfo{Name: name}
			continue
		}

		// Parse key-value pairs for current profile
		if currentProfile != nil && strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])

				switch key {
				case "region":
					currentProfile.Region = value
				case "sso_session":
					currentProfile.SSOSession = value
				case "sso_account_id":
					currentProfile.AccountID = value
				}
			}
		}
	}

	// Don't forget the last profile
	if currentProfile != nil && currentProfile.SSOSession != "" {
		profiles = append(profiles, *currentProfile)
	}

	return profiles, scanner.Err()
}
