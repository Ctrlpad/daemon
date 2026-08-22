package desktop

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/ini.v1"
)

type Entry struct {
	Path     string
	Binary   string
	Argv     []string
	Name     string `ini:"Name"`
	Exec     string `ini:"Exec"`
	Terminal bool   `ini:"Terminal"`
}

func Find(target string) (Entry, error) {
	entryPath, err := lookupEntry(target)
	if err != nil {
		return Entry{}, err
	}
	entry, err := loadEntry(entryPath)
	if err != nil {
		return Entry{}, err
	}
	return entry, nil
}

func lookupEntry(target string) (string, error) {
	for _, dir := range searchDirs() {
		path := filepath.Join(dir, "applications", target+".desktop")
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", nil
}

func searchDirs() []string {
	return filepath.SplitList(os.Getenv("XDG_DATA_DIRS"))
}

func loadEntry(entryPath string) (Entry, error) {
	entryData, err := ini.Load(entryPath)
	if err != nil {
		return Entry{}, err
	}

	var entry Entry
	err = entryData.Section("Desktop Entry").MapTo(&entry)
	if err != nil {
		return Entry{}, err
	}
	entry.Path = entryPath
	entry.Argv = parseExec(entry.Exec)
	entry.Binary, err = resolveBinary(entry)
	if err != nil {
		return Entry{}, err
	}

	return entry, nil
}

func parseExec(value string) []string {
	fields := strings.Fields(value)
	argv := make([]string, 0, len(fields))

	for _, field := range fields {
		switch {
		case field == "%%":
			argv = append(argv, "%")
		case len(field) == 2 && strings.HasPrefix(field, "%"):
			continue
		default:
			argv = append(argv, field)
		}
	}
	return argv
}

func resolveBinary(entry Entry) (string, error) {
	if len(entry.Argv) == 0 {
		return "", fmt.Errorf("desktop entry %s has no Exec", entry.Path)
	}
	binary := entry.Argv[0]

	if strings.Contains(binary, "/") {
		if !isExecutable(binary) {
			return "", fmt.Errorf("binary %s is not executable", binary)
		}
		return binary, nil
	}

	if prefix, _, found := strings.Cut(entry.Path, "/share/applications/"); found {
		if candidate := filepath.Join(prefix, "bin", binary); isExecutable(candidate) {
			return candidate, nil
		}
	}

	binaryPath, err := exec.LookPath(binary)
	if err != nil {
		return "", fmt.Errorf("binary %s not found: %w", binary, err)
	}
	return binaryPath, nil
}

func isExecutable(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir() && info.Mode()&0o111 != 0
}
