package backup

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/nathabonfim59/dncensor/internal/stack"
)

type Backup struct {
	Hash      string    `json:"hash"`
	Name      string    `json:"name"`
	StackType string    `json:"stack_type"`
	Content   []byte    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func Create(snapshotsDir, name string, s stack.Stack) (*Backup, error) {
	content, err := s.CaptureDNS()
	if err != nil {
		return nil, fmt.Errorf("capture DNS: %w", err)
	}

	h := sha256.Sum256(content)
	hash := hex.EncodeToString(h[:])

	b := &Backup{
		Hash:      hash,
		Name:      name,
		StackType: string(s.Type()),
		Content:   content,
		CreatedAt: time.Now(),
	}

	data, err := json.Marshal(b)
	if err != nil {
		return nil, fmt.Errorf("marshal backup: %w", err)
	}

	path := filepath.Join(snapshotsDir, hash+".json")
	if err := os.WriteFile(path, data, 0600); err != nil {
		return nil, fmt.Errorf("write backup: %w", err)
	}

	return b, nil
}

func List(snapshotsDir string) ([]*Backup, error) {
	entries, err := os.ReadDir(snapshotsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read snapshots dir: %w", err)
	}

	var backups []*Backup
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(snapshotsDir, entry.Name()))
		if err != nil {
			continue
		}
		var b Backup
		if err := json.Unmarshal(data, &b); err != nil {
			continue
		}
		backups = append(backups, &b)
	}

	sort.Slice(backups, func(i, j int) bool {
		return backups[i].CreatedAt.After(backups[j].CreatedAt)
	})

	return backups, nil
}

func Find(snapshotsDir, hashOrName string) (*Backup, error) {
	backups, err := List(snapshotsDir)
	if err != nil {
		return nil, err
	}

	var matches []*Backup
	for _, b := range backups {
		if b.Name == hashOrName {
			return b, nil
		}
		if strings.HasPrefix(b.Hash, hashOrName) {
			matches = append(matches, b)
		}
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("no backup found matching %q", hashOrName)
	}
	if len(matches) > 1 {
		var names []string
		for _, m := range matches {
			names = append(names, fmt.Sprintf("%s (%s)", m.Hash[:12], m.Name))
		}
		return nil, fmt.Errorf("multiple backups match %q: %s", hashOrName, strings.Join(names, ", "))
	}

	return matches[0], nil
}

func Delete(snapshotsDir, hashOrName string) error {
	b, err := Find(snapshotsDir, hashOrName)
	if err != nil {
		return err
	}

	path := filepath.Join(snapshotsDir, b.Hash+".json")
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("delete backup: %w", err)
	}

	return nil
}

func (b *Backup) Restore(s stack.Stack) error {
	if err := s.ApplyDNS(b.Content); err != nil {
		return fmt.Errorf("apply backup: %w", err)
	}
	return nil
}
