# Collect

Version 0.1.0

Collect is a simple tool to generate file lists for makefile targets. That's it.

## Install

Prequesits: go 1.21.0

You can either install collect by cloning this repository and running `go build` or by running `go install github.com/monoton-dev/collect`.

## Usage

Add a `.collect` file in the root directory of your repository / project. Collect searches for the .collect file in the current working directory (CWD).

Inside the `.collect` file you can:

- use comments by starting a line with a hash / pound sign (#),
- specify target names in brackets ([]). Example `[main]` would specify the target `main`,
- add file paths and glob patterns by indenting them with four spaces under a target.

```
usage: collect [options]
options:
  help             Displays this help message.
  version          Displays the build version.
  license          Displays the license for collect.
```

By running `collect` without any arguments, collect tries to find all files specified for the targets and invokes make.

## Use Cases

Ideal for projects using make that require wildcards or other path expanding tools. Collect offers a straightforward way to declare dependencies for your make targets.

This project is licensed under the Zlib license.
