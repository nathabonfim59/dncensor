package backup

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type BackupRecord struct {
	Timestamp   time.Time
	StackType   string
	BackupPath  string
	Description string
}

type BackupManager struct {
	BackupDir string
}

func New(backupDir string) *BackupManager {
	return &BackupManager{BackupDir: backupDir}
}

func (bm *BackupManager) Init() error {
	return os.MkdirAll(bm.BackupDir, 0700)
}

func (bm *BackupManager) List() ([]BackupRecord, error) {
	entries, err := os.ReadDir(bm.BackupDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var records []BackupRecord
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}

		stackType := detectStackType(e.Name())
		records = append(records, BackupRecord{
			Timestamp:   info.ModTime(),
			StackType:   string(stackType),
			BackupPath:  filepath.Join(bm.BackupDir, e.Name()),
			Description: e.Name(),
		})
	}

	sort.Slice(records, func(i, j int) bool {
		return records[i].Timestamp.After(records[j].Timestamp)
	})

	return records, nil
}

func detectStackType(name string) string {
	if strings.HasPrefix(name, "resolved") {
		return "systemd-resolved"
	}
	if strings.HasPrefix(name, "nm-") {
		return "networkmanager"
	}
	if strings.HasPrefix(name, "resolv.conf") {
		return "resolvconf"
	}
	return "unknown"
}

func (bm *BackupManager) Latest(stackType string) (*BackupRecord, error) {
	records, err := bm.List()
	if err != nil {
		return nil, err
	}

	for _, r := range records {
		if r.StackType == stackType {
			return &r, nil
		}
	}

	return nil, fmt.Errorf("no backup found for stack type %s", stackType)
}

func (bm *BackupManager) Exists(stackType string) bool {
	_, err := bm.Latest(stackType)
	return err == nil
}
