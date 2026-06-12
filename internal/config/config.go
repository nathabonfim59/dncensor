package config

import (
	"os"
)

const (
	AppName   = "dncensor"
	ConfigDir = "/etc/dncensor"
	BackupDir = "/etc/dncensor/backup"
)

func Init() error {
	for _, dir := range []string{ConfigDir, BackupDir} {
		if err := os.MkdirAll(dir, 0700); err != nil {
			return err
		}
	}
	return nil
}

func BackupPath() string {
	return BackupDir
}
