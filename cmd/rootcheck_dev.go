//go:build dev

package cmd

// requireRoot is a no-op in dev builds so the UI can be tested without root.
func requireRoot(string) {}

// initConfig is a no-op in dev builds.
func initConfig() error { return nil }
