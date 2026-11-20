// Copyright (C) 2026 monoton-dev
//
// This software is provided 'as-is', without any express or implied
// warranty. In no event will the authors be held liable for any damages
// arising from the use of this software.
//
// Permission is granted to anyone to use this software for any purpose,
// including commercial applications, and to alter it and redistribute it
// freely, subject to the following restrictions:
//
// 1. The origin of this software must not be misrepresented; you must not
//    claim that you wrote the original software. If you use this software
//    in a product, an acknowledgment in the product documentation would be
//    appreciated but is not required.
// 2. Altered source versions must be plainly marked as such, and must not be
//    misrepresented as being the original software.
// 3. This notice may not be removed or altered from any source distribution.

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

const HELP_TEXT string = `usage: collect [options]
options:
  help             Displays this help message.
  version          Displays the build version.
  license          Displays the license for collect.`
const VERSION_TEXT string = "Version 0.1.0"
const LICENSE_TEXT string = `Copyright (C) 2026 monoton-dev

This software is provided 'as-is', without any express or implied
warranty.  In no event will the authors be held liable for any damages
arising from the use of this software.

Permission is granted to anyone to use this software for any purpose,
including commercial applications, and to alter it and redistribute it
freely, subject to the following restrictions:

1. The origin of this software must not be misrepresented; you must not
   claim that you wrote the original software. If you use this software
   in a product, an acknowledgment in the product documentation would be
   appreciated but is not required.
2. Altered source versions must be plainly marked as such, and must not be
   misrepresented as being the original software.
3. This notice may not be removed or altered from any source distribution.`

func startswith(str string, substr string) bool {
	if len(str) >= len(substr) {
		for i := 0; i < len(substr); i++ {
			if str[i] != substr[i] {
				return false
			}
		}
		return true
	}
	return false
}

func firstchar(str string) byte {
	for chr := range str {
		if str[chr] > ' ' {
			return str[chr]
		}
	}
	return 0
}

func endswith(str string, substr string) bool {
	if len(str) >= len(substr) {
		var j int = 0
		for i := len(str) - len(substr) - 1; i < len(substr); i++ {
			if str[i] != substr[j] {
				return false
			}
			j += 1
		}
		return true
	}
	return false
}

func is_valid_target(targets []Target, target Target) bool {
	if len(target.files) == 0 {
		fmt.Println("error: " + target.name + " has no files")
		return false
	}
	for it := range targets[0 : len(targets)-1] {
		if targets[it].name == target.name {
			fmt.Println("error: " + target.name + " already declared")
			return false
		}
	}
	return true
}

func find_target_files(target *Target) {
	for _, pattern := range target.files {
		files, err := filepath.Glob(pattern)
		if err != nil {
			fmt.Println(fmt.Errorf("%s", err.Error()))
			continue
		}
		target.found = append(target.found, files[:]...)
	}
	if len(target.found) == 0 {
		fmt.Println("error: could not resolve files for target " + target.name)
	}
}

type Target struct {
	name  string
	files []string
	found []string
}

func main() {
	args := os.Args[1:]

	if len(args) == 1 {
		if args[0] == "help" {
			fmt.Println(HELP_TEXT)
			return
		}
		if args[0] == "version" {
			fmt.Println(VERSION_TEXT)
			return
		}
		if args[0] == "license" {
			fmt.Println(LICENSE_TEXT)
			return
		}
	}
	if len(args) > 0 {
		fmt.Println("error: unknown argument(s)")
		fmt.Println(HELP_TEXT)
		os.Exit(1)
	}

	config, err := os.ReadFile(".collect")
	if err != nil {
		fmt.Println("error loading config:", err)
		os.Exit(1)
	}

	fmt.Println("> Reading config")

	var lines []string = strings.Split(string(config), "\n")

	var targets []Target

	for index, line := range lines {
		if len(line) == 0 || firstchar(line) == 0 {
			continue
		}
		if firstchar(line) == '#' {
			continue
		} else if startswith(line, "[") && endswith(line, "]") && len(line) > 2 {
			if len(targets) >= 1 {
				if !is_valid_target(targets, targets[len(targets)-1]) {
					os.Exit(1)
				}
			}
			targets = append(targets, Target{name: line[1 : len(line)-1], files: []string{}})
		} else if startswith(line, "    ") {
			targets[len(targets)-1].files = append(targets[len(targets)-1].files, line[4:])
		} else {
			fmt.Println("error in line ", index+1)
			os.Exit(1)
		}
	}

	fmt.Println("> Resolving targets")

	var make_args []string

	for _, target := range targets {
		find_target_files(&target)
		make_args = append(make_args, target.name+"=\""+strings.Join(target.found, " ")+"\" ")
	}

	var make_path, lookErr = exec.LookPath("make")

	if lookErr != nil {
		fmt.Println(fmt.Errorf("%v", lookErr))
		os.Exit(1)
	}

	fmt.Println("> Invoking make\n")

	syscall.Exec(make_path, make_args, os.Environ())
}
