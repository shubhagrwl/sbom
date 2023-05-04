# Define the Go command and flags
GO := go
GOFLAGS := -ldflags="-s -w"

# Define the output binary name
BINARY := myprogram

# Default build target
default: all

# Build for all platforms
all: windows linux macos

# Build for Windows
windows:
	GOOS=windows GOARCH=amd64 $(GO) build $(GOFLAGS) -o $(BINARY)-windows-amd64.exe

# Build for Linux
linux:
	GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) -o $(BINARY)-linux-amd64

# Build for macOS
macos:
	GOOS=darwin GOARCH=amd64 $(GO) build $(GOFLAGS) -o $(BINARY)-macos-amd64

# Clean build artifacts
clean:
	rm -f $(BINARY)-windows-amd64.exe
	rm -f $(BINARY)-linux-amd64
	rm -f $(BINARY)-macos-amd64
