# Butler Git Hooks

## Usage

1. Create a folder `$GIT_DIR/git_hooks` in the root directory of your template.
2. Create a git hook specific file in `git_hooks` e.g `$GIT_DIR/git_hooks/pre-commit`
3. Run Butler and run `Install Git Hooks`

**pre-commit**
```sh
#!/bin/sh

echo "hook executed!"
```

_Git Hooks are installed automatically when a new project template is created._

## Run hooks in different languages

Node.js
```sh
#!/usr/bin/env node

console.log("hook executed!")
```

Python
```py
#!/usr/bin/python

print "hook executed!"
```

Powershell
```ps
powershell -ExecutionPolicy RemoteSigned -Command .\.git_hooks\scripts\build.ps1
```

## Details
Hooks are programs you can place in the `$GIT_DIR/git_hooks` directory to
trigger actions at certain points in git's execution.
 
Before Git invokes a hook, it changes its working directory to either
the root of the working tree in a non-bare repository, or to the
$GIT_DIR in a bare repository.
 
Hooks can get their arguments via the environment, command-line
arguments, and stdin. See the documentation for each hook below for
details.
 
The currently supported hooks are described [here](https://git-scm.com/docs/githooks).

