package backup

import (
	"backup-tool/config"
	"backup-tool/internal/logging"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Backuper struct {
	config config.PostgresConfig
	backup config.BackupConfig
	logger *logging.Logger
}

func NewBackuper(pgConfig config.PostgresConfig, backupConfig config.BackupConfig, logger *logging.Logger) *Backuper {
	return &Backuper{
		config: pgConfig,
		backup: backupConfig,
		logger: logger,
	}
}

func (b *Backuper) CreateBackup() error {
	timestamp := time.Now().Format("20060102_150405")
	backupDir := filepath.Join(b.backup.OutputDir, fmt.Sprintf("backup_%s", timestamp))

	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %v", err)
	}

	b.logger.Info("Creating backup in directory: %s", backupDir)

	cmdArgs := []string{
		"-h", b.config.Host,
		"-p", fmt.Sprintf("%d", b.config.Port),
		"-U", b.config.User,
		"-d", b.config.Database,
		"-F", "d",
		"-f", backupDir,
		"-j", "4",
		"-v",
	}

	env := append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", b.config.Password))

	cmd := exec.Command("pg_dump", cmdArgs...)
	cmd.Env = env

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	b.logger.Info("Starting pg_dump...")

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("pg_dump failed: %v, stderr: %s", err, stderr.String())
	}

	b.logger.Info("pg_dump completed successfully")

	if err := b.splitLargeFiles(backupDir); err != nil {
		b.logger.Warn("Failed to split files: %v", err)
	}

	b.cleanOldBackups()

	return nil
}

func (b *Backuper) cleanOldBackups() {
	cutoff := time.Now().AddDate(0, 0, -b.backup.RetentionDays)

	entries, err := os.ReadDir(b.backup.OutputDir)
	if err != nil {
		b.logger.Error("Failed to read backup directory: %v", err)
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() || !strings.HasPrefix(entry.Name(), "backup_") {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		if info.ModTime().Before(cutoff) {
			path := filepath.Join(b.backup.OutputDir, entry.Name())
			b.logger.Info("Removing old backup: %s", path)
			os.RemoveAll(path)
		}
	}
}

func (b *Backuper) splitLargeFiles(backupDir string) error {
	splitter := NewFileSplitter(b.backup.SplitSizeMB, b.logger)

	return filepath.Walk(backupDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		maxBytes := int64(b.backup.SplitSizeMB) * 1024 * 1024
		if info.Size() > maxBytes {
			b.logger.Info("Splitting large file: %s (%.2f MB)", path, float64(info.Size())/1024/1024)
			return splitter.SplitFile(path)
		}

		return nil
	})
}
