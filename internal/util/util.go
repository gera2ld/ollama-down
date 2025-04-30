package util

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "path/filepath"
    "regexp"
    "strings"

    "ollama_down/internal/constants"
)

type ModelInfo struct {
    Model string
    Tag   string
}

type Layer struct {
    Digest string `json:"digest"`
    Size   int64  `json:"size"`
}

type Manifest struct {
    Layers []Layer `json:"layers"`
}

// ParseModelName parses the model name into its components.
func ParseModelName(model string) ModelInfo {
    parts := strings.Split(model, ":")
    if len(parts) == 1 {
        return ModelInfo{
            Model: parts[0],
            Tag:   "latest",
        }
    }
    return ModelInfo{
        Model: parts[0],
        Tag:   parts[1],
    }
}

// FillTemplate fills a URL template with the provided values.
func FillTemplate(pattern string, values map[string]string) string {
    re := regexp.MustCompile(`\{\{(\w+)\}\}`)
    return re.ReplaceAllStringFunc(pattern, func(match string) string {
        key := match[2 : len(match)-2] // Remove {{ and }}
        if value, ok := values[key]; ok {
            return value
        }
        return match
    })
}

// FetchManifest fetches the manifest for a given model name.
func FetchManifest(modelInfo ModelInfo) (Manifest, error) {
    manifestURL := FillTemplate(constants.ManifestURLPattern, map[string]string{
        "model": modelInfo.Model,
        "tag":   modelInfo.Tag,
    })

    resp, err := http.Get(manifestURL)
    if err != nil {
        return Manifest{}, fmt.Errorf("failed to fetch manifest: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return Manifest{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return Manifest{}, fmt.Errorf("failed to read response body: %w", err)
    }

    var manifest Manifest
    if err := json.Unmarshal(body, &manifest); err != nil {
        return Manifest{}, fmt.Errorf("failed to parse manifest: %w", err)
    }

    return manifest, nil
}

// GetLayerFileName returns the file name for a given layer digest.
func GetLayerFileName(digest string) string {
    return strings.ReplaceAll(digest, ":", "-")
}

// EnsureDir ensures that the directory exists, creating it if necessary.
func EnsureDir(dir string) error {
    if err := os.MkdirAll(dir, os.ModePerm); err != nil {
        return fmt.Errorf("failed to create directory %s: %w", dir, err)
    }
    return nil
}

// JoinPaths joins multiple path elements into a single path.
func JoinPaths(elements ...string) string {
    return filepath.Join(elements...)
}
