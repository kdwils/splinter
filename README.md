# Splinter

Simple tool to split or flatten kubernetes manifests. 

It is often the case that I will use `helm template` where certain values are not supported and end up with a single massive manifest which is messy to work with. 

Instead, I would prefer to have multiple smaller manifests separated by resource kind.

# Usage

### Commands
| Command | Subcommand | Description |
| --- | ----------- | --- |
| splinter | split | split a single manifest into multiple manifests organized by resource kind |
| splinter | merge | merge individual manifests into a single manifest |


### Flags
| Flag | Required | Description | Command |
| --- | ----------- | --- | --- |
| --include, -i | false | optional flag to set what files or directories to include | all |
| --kustomize, -k | false | spit out a kustomization.yaml for the splintered or merged manifest(s)  | all |
| --output, -o | false | directory to write manifest(s) to | all |
| --delete, -d | false | delete all files used when merging manifests | merge |
| --std-out | false | print to stdout instead of writing to disk | merge |


# Examples

```
splinter split --input examples/metallb.yaml --output examples/split/
```

```
splinter split examples/metallb.yaml -o examples/split/
```

```
splinter merge examples/split/ -o examples/flatten/my-manifest.yaml
```

You can also pipe from stdin

```
helm template my-release sealed-secrets/sealed-secrets | splinter split -o my-dir/
helm template my-release sealed-secrets/sealed-secrets | splinter split -i anotherfile.yaml -o my-dir/
````

# Installation

Go:

`go install github.com/kdwils/splinter@latest`

Github:

check out the latest [release](https://github.com/kdwils/splinter/releases)