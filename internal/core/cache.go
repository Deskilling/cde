package core

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"

	"cde/internal/editor/model"
)

type CacheMap struct {
	Hash      string          `json:"hash"`
	Workspace model.Workspace `json:"workspace"`
}

type CacheJson struct {
	Workspaces map[string]CacheMap `json:"cachemap"`
}

var Cache *CacheJson
var LatestHash string

func GetCacheDir() string {
	xdgCacheHome := os.Getenv("XDG_CACHE_HOME")
	if xdgCacheHome != "" {
		return xdgCacheHome
	}

	return filepath.Join(os.Getenv("HOME"), ".cache")
}

func GetCacheLocation() string {
	locations := map[string]string{
		"darwin": filepath.Join(GetCacheDir(), "cde-cache.json"),
		"linux":  filepath.Join(GetCacheDir(), "cde-cache.json"),
	}

	path, ok := locations[runtime.GOOS]
	if ok {
		return path
	}

	return filepath.Join(GetCacheDir(), "cde-cache.json")
}

func LoadCache() (err error) {
	content, err := os.ReadFile(GetCacheLocation())
	if err != nil {
		slog.Error("hmm", "err", err)
		_, err := os.Create(GetCacheLocation())
		if err != nil {
			return fmt.Errorf("failed creating cache file: %w", err)
		}
		Cache = &CacheJson{
			map[string]CacheMap{},
		}
	}

	err = json.Unmarshal(content, &Cache)
	if err != nil {
		slog.Error("hmm", "err", err)
		Cache = &CacheJson{
			map[string]CacheMap{},
		}
	}

	return nil
}

func CalculateHash(path string) (hashSha1 string, err error) {
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return "", err
	}
	defer func() {
		err = file.Close()
	}()

	shaHash := sha256.New()

	_, err = io.Copy(io.MultiWriter(shaHash), file)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(shaHash.Sum(nil)), err
}

func ValidCache(name string, path string) bool {
	entry, ok := Cache.Workspaces[name]

	hash, err := CalculateHash(path)
	if err != nil {
		return false
	}

	slog.Debug("got sha", "path", path, "sha", hash)
	slog.Debug("entry hash", "hash", entry.Hash)

	if ok && hash == entry.Hash {
		return true
	} else {
		entry.Hash = hash
		Cache.Workspaces[name] = entry
	}

	LatestHash = entry.Hash
	slog.Debug(entry.Hash)

	return false
}

func WriteCache(name string, hash string, ws model.Workspace) error {
	entry := Cache.Workspaces[name]

	entry.Hash = hash
	entry.Workspace = ws

	Cache.Workspaces[name] = entry

	content, err := json.Marshal(Cache)
	if err != nil {
		return fmt.Errorf("failed to marshal cache: %w", err)
	}

	err = os.WriteFile(GetCacheLocation(), content, 0600)
	if err != nil {
		return fmt.Errorf("failed to save Cache: %w", err)
	}

	return nil
}
