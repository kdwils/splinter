# Splinter

Simple tool to split or flatten kubernetes yaml resources. It is often the case that I will use `helm template` where certain values are not supported and end up with a single massive manifest which is messy to work with. 

Instead, I would prefer to have multiple yaml files where each resource kind lives in its own file.

# Usage

### Commands
| Command | Subcommand | Description |
| --- | ----------- | --- |
| splinter | - | splinter a single manifest into multiple manifests organize by resource kind |

### Flags
| Flag | Required | Description |
| --- | ----------- | --- |
| --include, -i | false | optional flag to set what files or directories to include
| --kustomize, -k | false | spit out a kustomization.yaml for the splintered or merged manifest(s)
| --output, -o | false | directory to write manifest(s) to
| --merge | false | merge a list of manifests into a single manifest


# Examples

```
splinter --input examples/metallb.yaml --output examples/split/
```

```
splinter examples/metallb.yaml -o examples/split/
```

```
splinter --merge examples/split/ -o examples/flatten/my-manifest.yaml
```

You can also pipe from stdin

```
helm template my-release sealed-secrets/sealed-secrets | splinter -o my-dir/
helm template my-release sealed-secrets/sealed-secrets | splinter -i anotherfile.yaml -o my-dir/
````

# Installation

Go:

`go install github.com/kdwils/splinter@latest`

Github:

check out the latest [release](https://github.com/kdwils/splinter/releases)