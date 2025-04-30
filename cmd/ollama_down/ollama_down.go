package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"ollama_down/internal/constants"
	"ollama_down/internal/util"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "ollama-down",
		Short: "Ollama model manager",
	}

	var getCmd = &cobra.Command{
		Use:   "get <model>",
		Short: "Display information about a model",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			format, _ := cmd.Flags().GetString("format")
			outFile, _ := cmd.Flags().GetString("outFile")

			template, ok := constants.InfoTemplates[format]
			if !ok {
				fmt.Printf("Invalid format: %s\n", format)
				os.Exit(1)
			}

			modelInfo := util.ParseModelName(args[0])
			manifest, err := util.FetchManifest(modelInfo)
			if err != nil {
				fmt.Printf("Error fetching manifest: %v\n", err)
				os.Exit(1)
			}

			var items []constants.DownloadItem
			for _, layer := range manifest.Layers {
				url := util.FillTemplate(constants.LayerURLPattern, map[string]string{
					"model":  modelInfo.Model,
					"tag":    modelInfo.Tag,
					"digest": layer.Digest,
				})
				items = append(items, constants.DownloadItem{
					URL:      url,
					FileName: util.GetLayerFileName(layer.Digest),
				})
			}

			cacheDirPrefix := constants.CacheDir
			if cacheDirPrefix != "" && !strings.HasSuffix(cacheDirPrefix, "/") {
				cacheDirPrefix += "/"
			}

			output := template(items, map[string]string{
				"cacheDir": cacheDirPrefix,
			})

			if outFile != "" {
				if err := os.WriteFile(outFile, []byte(output), 0644); err != nil {
					fmt.Printf("Error writing to file: %v\n", err)
					os.Exit(1)
				}
			} else {
				fmt.Println(output)
			}
		},
	}

	var installCmd = &cobra.Command{
		Use:   "install <model>",
		Short: "Install a model",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			root, _ := cmd.Flags().GetString("root")
			root = os.ExpandEnv(strings.Replace(root, "~", "$HOME", 1))
			blobsRoot := filepath.Join(root, "models/blobs")

			modelInfo := util.ParseModelName(args[0])
			manifest, err := util.FetchManifest(modelInfo)
			if err != nil {
				fmt.Printf("Error fetching manifest: %v\n", err)
				os.Exit(1)
			}

			for _, layer := range manifest.Layers {
				fileName := util.GetLayerFileName(layer.Digest)
				source := filepath.Join("cache", fileName)
				target := filepath.Join(blobsRoot, fileName)

				if _, err := os.Stat(target); err == nil {
					fmt.Printf("Layer exists: %s\n", fileName)
					continue
				}

				fmt.Printf("Installing layer: %s\n", fileName)
				if err := os.Rename(source, target); err != nil {
					fmt.Printf("Error installing layer: %v\n", err)
					os.Exit(1)
				}
			}

			fmt.Println("Layers installed. Run command to finish:")
			fmt.Printf("ollama pull %s\n", args[0])
		},
	}

	// Add flags
	getCmd.Flags().String("format", "links", "Output format (links, curl, aria2)")
	getCmd.Flags().String("outFile", "", "Output file")
	installCmd.Flags().String("root", "~/.ollama", "Ollama root directory")

	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(installCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
