package tui

import (
	"github.com/natevick/s3-tui/internal/aws"
	"github.com/natevick/s3-tui/internal/bookmarks"
	"github.com/natevick/s3-tui/internal/download"
)

// ViewType represents the current active view
type ViewType int

const (
	ViewProfiles ViewType = iota
	ViewBuckets
	ViewBrowser
	ViewDownload
	ViewBookmarks
	ViewHelp
)

// Message types for inter-component communication

// BucketsLoadedMsg is sent when buckets are loaded
type BucketsLoadedMsg struct {
	Buckets []aws.Bucket
	Err     error
}

// BucketSelectedMsg is sent when a bucket is selected
type BucketSelectedMsg struct {
	Bucket string
	Region string
}

// ObjectsLoadedMsg is sent when objects are loaded
type ObjectsLoadedMsg struct {
	Objects []aws.S3Object
	Prefix  string
	Err     error
}

// NavigatePrefixMsg is sent when navigating to a prefix
type NavigatePrefixMsg struct {
	Prefix string
}

// NavigateBackMsg is sent when navigating back
type NavigateBackMsg struct{}

// SwitchViewMsg is sent to switch to a different view
type SwitchViewMsg struct {
	View ViewType
}

// StartDownloadMsg initiates a download
type StartDownloadMsg struct {
	Bucket    string
	Key       string
	LocalPath string
	IsPrefix  bool // true if downloading all files under a prefix
}

// DownloadProgressMsg reports download progress
type DownloadProgressMsg struct {
	Progress download.Progress
}

// DownloadCompleteMsg is sent when a download completes
type DownloadCompleteMsg struct {
	Key string
	Err error
}

// AllDownloadsCompleteMsg is sent when all downloads are done
type AllDownloadsCompleteMsg struct {
	TotalFiles int
	TotalBytes int64
	Failed     int
}

// CancelDownloadMsg cancels the current download
type CancelDownloadMsg struct{}

// BookmarksLoadedMsg is sent when bookmarks are loaded
type BookmarksLoadedMsg struct {
	Bookmarks []bookmarks.Bookmark
	Err       error
}

// AddBookmarkMsg adds a bookmark
type AddBookmarkMsg struct {
	Bucket string
	Prefix string
	Name   string
}

// BookmarkAddedMsg confirms bookmark was added
type BookmarkAddedMsg struct {
	Bookmark bookmarks.Bookmark
	Err      error
}

// RemoveBookmarkMsg removes a bookmark
type RemoveBookmarkMsg struct {
	ID string
}

// BookmarkRemovedMsg confirms bookmark was removed
type BookmarkRemovedMsg struct {
	ID  string
	Err error
}

// SelectBookmarkMsg selects a bookmark to navigate to
type SelectBookmarkMsg struct {
	Bookmark bookmarks.Bookmark
}

// ErrorMsg reports an error
type ErrorMsg struct {
	Err error
}

// StatusMsg updates the status bar
type StatusMsg struct {
	Message string
}

// WindowSizeMsg is sent when the terminal window is resized
type WindowSizeMsg struct {
	Width  int
	Height int
}

// PromptInputMsg is used for text input prompts
type PromptInputMsg struct {
	Prompt       string
	DefaultValue string
	Callback     func(string) // Called with the input value
}

// PromptConfirmMsg is used for confirmation prompts
type PromptConfirmMsg struct {
	Prompt   string
	Callback func(bool) // Called with true for yes, false for no
}

// ClosePromptMsg closes any open prompt
type ClosePromptMsg struct{}

// RefreshMsg requests a refresh of the current view
type RefreshMsg struct{}

// TickMsg is sent for periodic updates
type TickMsg struct{}
