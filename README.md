# Splinter

A command-line tool to split or merge Kubernetes manifests.

## Use case

Splinter simplifies working with large Kubernetes manifests. 

It helps in scenarios where you need to:
- Break down large manifests into manageable, individual files organized by resource type
- Generate a `kustomization.yaml` to glue the individual files together

Or...

- Recombine the individual files back into a single manifest

## Installation

### Using Go

```bash
go install github.com/kdwils/splinter@latest
```

### From Releases

Download the latest binary from [GitHub Releases](https://github.com/kdwils/splinter/releases)

## Usage

### Commands

| Command | Description |
|---------|-------------|
| `split` | Split a single manifest into multiple files organized by resource kind |
| `merge` | Merge multiple manifest files into a single output (prints to stdout by default) |

### Global Flags

| Flag | Short | Required | Description |
|------|--------|----------|-------------|
| `--include` | `-i` | No | Files or directories to include |
| `--output` | `-o` | No | Output directory/file path |


## Examples

### Splitting Manifests

![split gif](vhs/split.gif)

With the include flag
```bash
splinter split -i examples/merged/merged.yaml -o examples/split/
```

With an argument
```bash
splinter split examples/merged/merged.yaml -o examples/split/
```

Split and generate a Kustomization file:
```bash
splinter split -k -i examples/merged/merged.yaml -o examples/split/
```

### Merging Manifests

![merge gif](vhs/merge.gif)

Merge multiple files into stdout:
```bash
splinter merge -i examples/split/
```

Merge into a specific file:
```bash
splinter merge -i examples/split/ -o examples/flatten/my-manifest.yaml
```

### Working with Pipes

Split Helm output:
```bash
helm template my-release sealed-secrets/sealed-secrets | splinter split -o my-dir/
```

Split Helm output and include additional files:
```bash
helm template my-release sealed-secrets/sealed-secrets | splinter split -i existing.yaml -o my-dir/
```

## Development

This project uses [Nix](https://nixos.org/) for development environment consistency. 

Assuming a working Nix installation:

#### Start the nix develop shell

```shell
nix develop .
```
#### Run tests
```shell
go test ./...
```

#### Build binary
```shell
go build
```

#### Run a command
```shell
go run main.go merge -i examples/split/
```

## Acknowledgements

* [vhs](https://github.com/charmbracelet/vhs) to make the gifs in this readme
