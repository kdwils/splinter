# Splinter

Simple tool to split or flatten kubernetes yaml resources. It is often the case that I will use `helm template` and end up with a single massive manifest which is messy to work with. 

Instead, I would prefer to have multiple yaml files where each resource kind lives in its own file.

# Usage
| Command | Subcommand | Description | Flag | Description | Flag | Description | Flag | Description |
| --- | ----------- | --- | --- | --- | --- | --- | --- | --- |
| splinter | - | splinter a single manifest into multiple files by resource kind | --input, -i | individual file or directory to splinter, or both  | --output, -o | directory to output manifests to | --kustomize, -k | spit out a simple kustomization.yaml for split resources |
| splinter | flatten | flatten multiple resource files into a single manifest | --input, -i | list of files or directory to flatten, or both | --output, -i | path & filename to output manifest to |

```
splinter --input examples/metallb.yaml --output examples/split/
```

```
splinter -i examples/ -o examples/split/
```

```
splinter flatten -i examples/split/ -o examples/flatten/my-manifest.yaml
```

You can also pipe from stdin

```
helm template my-release sealed-secrets/sealed-secrets | splinter -i some-file.yaml -o my-dir/
````

# Installation

`go get github.com/kdwils/splinter`