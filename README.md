# Splinter

Simple tool to split or flatten kubernetes yaml resources. It is often the case that I will use `helm template` and end up with a single massive manifest which is messy to work with. 

Instead, I would prefer to have multiple yaml files where each resource kind lives in its own file.

# Usage
| Command | Subcommand | Description | Flag | Description | Flag | Description |  
| --- | ----------- | --- | --- | --- | --- | --- |
| splinter | - | splinter a single manifest into multiple files by resource kind | --include, -i | individual file or directory to splinter, or both  | --output, -o | directory to output manifests to |
| splinter | flatten | flatten multiple resource files into a single manifest | --include, -i | list of files or directory to flatten, or both | --output, -i | path & filename to output manifest to |

`splinter --include examples/metallb.yaml --output examples/split/`

`splinter -i examples/ -o examples/split/`

`splinter flatten -i examples/split/ -o examples/flatten/my-manifest.yaml`

# Installation

`go get github.com/kdwils/splinter`