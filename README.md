# go-mnemonics

[![Go](https://img.shields.io/badge/go-1.24+-blue)](https://go.dev/)
[![License](https://img.shields.io/github/license/nduyhai/go-mnemonics)](LICENSE)
[![Update Starred](https://github.com/nduyhai/go-mnemonics/actions/workflows/update-starred.yml/badge.svg)](https://github.com/nduyhai/go-mnemonics/actions/workflows/update-starred.yml)

A personal knowledge stash for Go developers. This repository automatically tracks and organizes GitHub starred repositories to create a curated collection of Go resources.

## Features

- ✅ Automatic weekly updates of starred GitHub repositories
- ✅ Generates both JSON and Markdown formats of starred repos
- ✅ GitHub Actions workflow for automation
- ✅ Containerization support with Docker
- ✅ Minimal dependencies (uses only Go standard library)

## How It Works

This project uses a GitHub Actions workflow to:

1. Run weekly (every Monday at 2AM UTC)
2. Fetch all repositories starred by the user
3. Organize repositories by programming language
4. Sort repositories alphabetically within each language
5. Generate structured data files:
   - `starred.json` - Complete data in JSON format
   - Language-specific Markdown files in the `stars/` directory:
     - `stars/go/starred.md` - Go repositories
     - `stars/java/starred.md` - Java repositories
     - And so on for each programming language
6. Automatically commit and push the updated files

## Getting Started

### Prerequisites

- GitHub account
- Personal access token with `repo` scope

### Setup

1. **Fork or clone this repository**

   ```bash
   git clone https://github.com/nduyhai/go-mnemonics
   cd go-mnemonics
   ```

2. **Update GitHub username**

   Edit `.github/workflows/update-starred.yml` and change the `GITHUB_USERNAME` value to your GitHub username.

3. **Add GitHub token**

   Add your GitHub personal access token as a repository secret named `GH_STAR_TOKEN`.

4. **Run manually or wait for scheduled run**

   You can manually trigger the workflow from the Actions tab in your repository, or wait for the scheduled run.

## Running Locally

To run the script locally:

```bash
export GITHUB_USERNAME=your-username
export GITHUB_TOKEN=your-token
go run .github/scripts/fetch_starred.go
```

## Docker Support

Build and run the container:

```bash
# Build the container
docker build -t go-mnemonics .

# Run the container
docker run -p 8080:8080 go-mnemonics
```

## Knowledge Structure

This repository includes a structured collection of Go knowledge organized into the following categories:

### Syntax
- [Variables](syntax/variables.md) - Variable declarations, zero values, constants, and iota
- [Functions](syntax/functions.md) - Function declarations, multiple returns, closures, and defer
- [Interfaces](syntax/interfaces.md) - Interface declarations, implementation, and type assertions

### Patterns
- [Goroutines](patterns/goroutines.md) - Lightweight concurrency, synchronization, and worker pools
- [Channels](patterns/channels.md) - Communication between goroutines, buffering, and select
- [Context](patterns/context.md) - Cancellation, timeouts, and request-scoped values

### Gotchas
- [Nil vs. Zero Values](gotchas/nil-vs-zero.md) - Understanding nil and zero values in Go
- [Slice Pitfalls](gotchas/slice-pitfalls.md) - Common mistakes with slices and underlying arrays

### Tips
- [Performance](tips/performance.md) - Optimization techniques and best practices
- [Testing](tips/testing.md) - Effective testing strategies and patterns

## Customization

You can customize this project in several ways:

- Modify the script to categorize repositories
- Change the update frequency in the workflow file
- Extend the script to include additional metadata
- Implement a web interface to browse your knowledge stash
- Add your own notes and examples to the knowledge structure

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

