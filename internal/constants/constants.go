package constants

import (
	"fmt"
	"strings"
)

const (
	registry            = "registry.ollama.ai"
	ManifestURLPattern  = "https://" + registry + "/v2/library/{{model}}/manifests/{{tag}}"
	ManifestPathPattern = registry + "/library/{{model}}/{{tag}}"
	LayerURLPattern     = "https://" + registry + "/v2/library/{{model}}/blobs/{{digest}}"
	CacheDir            = "cache"
)

// InfoTemplateFunc defines the function signature for template formatters
type InfoTemplateFunc func(items []DownloadItem, values map[string]string) string

// DownloadItem represents a single item to be downloaded
type DownloadItem struct {
	URL      string
	FileName string
}

// InfoTemplates maps format names to their template functions
var InfoTemplates = map[string]InfoTemplateFunc{
	"links": func(items []DownloadItem, _ map[string]string) string {
		var result []string
		for _, item := range items {
			result = append(result, item.URL)
		}
		return strings.Join(result, "\n")
	},
	"curl": func(items []DownloadItem, values map[string]string) string {
		var result []string
		cacheDir := values["cacheDir"]
		for _, item := range items {
			result = append(result, fmt.Sprintf("curl -fLo %s%s %s", cacheDir, item.FileName, item.URL))
		}
		return strings.Join(result, "\n")
	},
	"aria2": func(items []DownloadItem, values map[string]string) string {
		var result []string
		cacheDir := values["cacheDir"]
		for _, item := range items {
			result = append(result, fmt.Sprintf("%s\n  out=%s%s", item.URL, cacheDir, item.FileName))
		}
		return strings.Join(result, "\n")
	},
}
