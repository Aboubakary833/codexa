# Codexa

**Codexa** is a terminal UI application designed to help developers quickly access concise, practical code snippets and common development patterns, without the verbosity of traditional documentation.

It is built for developers who want answers that go straight to the point.



## Demo video


<video src="https://github.com/user-attachments/assets/3547d373-68d2-48cc-97bb-00fa944697b7" width="600"></video>


---

## Why Codexa?

Documentation tools like `man` pages or online references are often too verbose for day-to-day development. While tools like `tldr` focus on command-line usage, Codexa targets a different need: **code and developer knowledge**.

Codexa focuses on:
- Short, focused code snippets
- Minimal explanations that highlight key ideas and pitfalls
- Developer patterns rather than full tutorials

The goal is not to teach from scratch, but to **refresh your memory quickly and efficiently**.

## Installation

### Linux / macOS

```bash
curl -sSf https://raw.githubusercontent.com/aboubakary833/codexa/main/scripts/install.sh | bash
```

### Windows (PowerShell)

```pwsh
iwr https://raw.githubusercontent.com/aboubakary833/codexa/main/scripts/install.ps1 -UseBasicParsing | iex
```

This will:

- Download the latest Codexa release
- Install the binary in a folder added to your PATH
- Install shell completions automatically for bash, zsh, fish, or PowerShell

## Design philosophy

- Concise over exhaustive
- Practical over theoretical
- Readable over clever
- Offline-first

<p>
Every snippet is intentionally short and focused on real-world usage.
To contribute to snippets visit: <a href="https://github.com/aboubakary833/cx-registry">Codexa snippets registry repository</a>
</p>

## Usage

### Running in browsing mode

To launch Codexa TUI run the following command:

```sh
codexa run
```

Or open a sipecific category of snippet, run:

```sh
codexa open js # [go|html|css]
```

### Downloading/Update snippets to local registry

To install/update for example "Go" set of snippets, run:

```sh
codexa sync go
```

You can install/update a specific snippet (context with timeout for example) by running:

```sh
codexa sync go context-timeout
```

## Status

Codexa is an experimental project built as a learning and exploration tool.  
The scope is intentionally limited to keep the project focused and maintainable.
