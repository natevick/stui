package download

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/natevick/s3-tui/internal/aws"
)

// SyncResult contains the result of a sync operation
type SyncResult struct {
	ToDownload []aws.S3Object // Files that need to be downloaded
	Unchanged  []aws.S3Object // Files that are already up to date
	TotalBytes int64          // Total bytes to download
}

// SyncManager handles sync operations
type SyncManager struct {
	client *aws.Client
}

// NewSyncManager creates a new sync manager
func NewSyncManager(client *aws.Client) *SyncManager {
	return &SyncManager{client: client}
}

// CompareFiles compares S3 objects with local files and returns sync plan
func (s *SyncManager) CompareFiles(ctx context.Context, bucket, prefix, localDir string) (*SyncResult, error) {
	// List all S3 objects
	objects, err := s.client.ListAllObjects(ctx, bucket, prefix)
	if err != nil {
		return nil, fmt.Errorf("failed to list S3 objects: %w", err)
	}

	// Build local file map
	localFiles, err := s.buildLocalFileMap(localDir, prefix)
	if err != nil {
		return nil, fmt.Errorf("failed to scan local directory: %w", err)
	}

	result := &SyncResult{}

	for _, obj := range objects {
		relPath := strings.TrimPrefix(obj.Key, prefix)
		localPath := filepath.Join(localDir, relPath)

		localInfo, exists := localFiles[relPath]
		if !exists {
			// File doesn't exist locally
			result.ToDownload = append(result.ToDownload, obj)
			result.TotalBytes += obj.Size
			continue
		}

		// Quick check: size comparison
		if localInfo.Size() != obj.Size {
			result.ToDownload = append(result.ToDownload, obj)
			result.TotalBytes += obj.Size
			continue
		}

		// Detailed check: ETag comparison
		// Note: For multipart uploads, ETag is not MD5, so we skip hash check for those
		if !strings.Contains(obj.ETag, "-") {
			localHash, err := computeFileMD5(localPath)
			if err != nil {
				// If we can't compute hash, download to be safe
				result.ToDownload = append(result.ToDownload, obj)
				result.TotalBytes += obj.Size
				continue
			}

			if localHash != obj.ETag {
				result.ToDownload = append(result.ToDownload, obj)
				result.TotalBytes += obj.Size
				continue
			}
		}

		// File matches
		result.Unchanged = append(result.Unchanged, obj)
	}

	return result, nil
}

// localFileInfo wraps os.FileInfo for our needs
type localFileInfo struct {
	os.FileInfo
	path string
}

// buildLocalFileMap builds a map of relative path -> file info
func (s *SyncManager) buildLocalFileMap(localDir, prefix string) (map[string]os.FileInfo, error) {
	files := make(map[string]os.FileInfo)

	// If directory doesn't exist, return empty map
	if _, err := os.Stat(localDir); os.IsNotExist(err) {
		return files, nil
	}

	err := filepath.Walk(localDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Get relative path
		relPath, err := filepath.Rel(localDir, path)
		if err != nil {
			return err
		}

		// Normalize path separators
		relPath = filepath.ToSlash(relPath)
		files[relPath] = info

		return nil
	})

	return files, err
}

// computeFileMD5 computes the MD5 hash of a file
func computeFileMD5(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// Sync performs a sync operation, downloading only changed/new files
func (s *SyncManager) Sync(ctx context.Context, bucket, prefix, localDir string, manager *Manager) error {
	// Compare files
	result, err := s.CompareFiles(ctx, bucket, prefix, localDir)
	if err != nil {
		return err
	}

	if len(result.ToDownload) == 0 {
		return nil // Nothing to download
	}

	// Initialize progress for sync
	files := make(map[string]*FileProgress)
	for _, obj := range result.ToDownload {
		relPath := strings.TrimPrefix(obj.Key, prefix)
		localPath := filepath.Join(localDir, relPath)
		files[obj.Key] = &FileProgress{
			Key:       obj.Key,
			LocalPath: localPath,
			Size:      obj.Size,
			Status:    StatusPending,
		}
	}

	manager.progressMu.Lock()
	manager.progress = Progress{
		TotalFiles: len(result.ToDownload),
		TotalBytes: result.TotalBytes,
		Files:      files,
		Status:     StatusInProgress,
	}
	manager.progressMu.Unlock()

	// Download the files
	return manager.downloadWithWorkers(ctx, bucket, result.ToDownload, prefix, localDir)
}
