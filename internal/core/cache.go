package core

import (
	"crypto/sha1"
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

func LoadCache() {
	content, err := os.ReadFile(GetCacheLocation())
	if err != nil {
		slog.Error("hmm", "err", err)
		os.Create(GetCacheLocation())
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
}

func CalculateHash(filepath string) (hashSha1 string,err error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	sha1Hash := sha1.New()

	_, err = io.Copy(io.MultiWriter(sha1Hash), file)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(sha1Hash.Sum(nil)), nil
}

func ValidCache(name string, path string) bool{
	entry, ok := Cache.Workspaces[name]

	hash, err := CalculateHash(path)
	if err != nil {
		return false
	}

	slog.Debug("got sha1", "path", path, "sha1", hash)
	slog.Debug("entry hash", "hash", entry.Hash)

	if ok && hash == entry.Hash{
		return true
	} else {
		entry.Hash = hash
		Cache.Workspaces[name] = entry
	}

	LatestHash = entry.Hash
	slog.Debug(entry.Hash)

	return false
}

func WriteCache(name string, hash string, ws model.Workspace) error{
	entry := Cache.Workspaces[name]

	entry.Hash = hash
	entry.Workspace = ws

	Cache.Workspaces[name] = entry

	content, err := json.Marshal(Cache)
	if err != nil {
		return fmt.Errorf("failed to marshal cache: %w", err)
	}

	err = os.WriteFile(GetCacheLocation(), content, 0777)
	if err != nil {
		return fmt.Errorf("failed to save Cache: %w", err)
	}

	return nil
}
