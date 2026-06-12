package config

import (
	"os"
)

const (
	AppName     = "dncensor"
	ConfigDir   = "/etc/dncensor"
	SnapshotsDir = "/etc/dncensor/snapshots"
)

func Init() error {
	for _, dir := range []string{ConfigDir, SnapshotsDir} {
		if err := os.MkdirAll(dir, 0700); err != nil {
			return err
		}
	}
	return nil
}

func SnapshotsPath() string {
	return SnapshotsDir
}
