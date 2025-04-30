# Ollama Down

Ollama Down is a command-line application for managing models in the Ollama ecosystem. It provides commands to fetch model information and install models efficiently.

## Why

Ollama is a great project. But it doesn't implement download properly. As a result, it is almost impossible to download a large model successfully in an average network.

With this tool, you can download the model data with more robust tools like `curl` or `aria2`, and then install the model easily with ollama.

## Quick Start

```bash
# List the links
./ollama_down get qwen3

# Download the files and save them to `./cache`
# Note that `sha256:xxx` should be renamed to `sha256-xxx`

# Copy the downloaded layers to ollama's storage
./ollama_down install qwen3

# Run ollama to verify the installation
ollama pull qwen3
```

You can also generate the data needed for download tools:

<details>
  <summary>For Aria2</summary>

```bash
# Generate a download script
./ollama_down get qwen3 --format aria2 --outFile download.txt

# Download layers
aria2c -c -x5 -m0 -k20m --no-conf -i download.txt

# Copy the downloaded layers to ollama's storage
./ollama_down install qwen3

# Run ollama to verify the installation
ollama pull qwen3
```

</details>

<details>
  <summary>For Curl</summary>

```bash
# Generate a download script
./ollama_down get qwen3 --format curl --outFile download.sh

# Download layers
sh download.sh

# Copy the downloaded layers to ollama's storage
./ollama_down install qwen3

# Run ollama to verify the installation
ollama pull qwen3
```

</details>

## Commands

### Get Command

Fetches information about a specific model.

```bash
./ollama-down get <model> [options]
```

Options:

- `--format <format>`: Output format, possible values are `links`, `curl`, `aria2` (default: `links`).
- `--outFile <outFile>`: Specify an output file to save the information.

### Install Command

Installs a specified model.

```bash
./ollama-down install <model> [options]
```

Options:

- `--root <root>`: Specify the Ollama root directory (default: `~/.ollama`).

## Build

```bash
go build cmd/ollama_down/ollama_down.go
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## License

This project is licensed under the MIT License. See the LICENSE file for more details.
