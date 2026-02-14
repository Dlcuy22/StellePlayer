# Makefile
# Build and development automation for Endmi.
#
# Targets:
#   - run: runs the application locally
#   - build-amd64: builds the binary for windows amd64
#   - build-all: builds binaries for Linux, Windows, and macOS (multiple architectures)

BINARY_NAME=StellePlayer
VERSION=1.0.0

# Detect OS
ifeq ($(OS),Windows_NT)
	RM = powershell -Command "if (Test-Path bin) { Remove-Item -Recurse -Force bin }; if (Test-Path $(BINARY_NAME).exe) { Remove-Item -Force $(BINARY_NAME).exe }"
	MKDIR = powershell -Command "if (!(Test-Path bin)) { New-Item -ItemType Directory bin }"
	MKDIR_BUILD = powershell -Command "if (!(Test-Path build)) { New-Item -ItemType Directory build }"
	RM_BUILD = powershell -Command "if (Test-Path build) { Remove-Item -Recurse -Force build }"
	MV_INSTALLER = powershell -Command "Move-Item -Path Scripts/$(BINARY_NAME)_Setup_$(VERSION).exe -Destination bin/$(BINARY_NAME)_Setup_$(VERSION).exe -Force"
	MAKENSIS = makensis
	GO_BUILD = powershell -Command "$$env:GOOS='$(1)'; $$env:GOARCH='$(2)'; go build -buildvcs=false -ldflags '-X main.Version=$(VERSION)' -o $(3) ."
	NSIS_STEP = powershell -Command "Write-Output '[Windows] Compiling NSIS installer...'; & makensis Scripts/build.nsi; if ($$LASTEXITCODE -eq 0) { Write-Output '[Windows] Moving installer to bin/...'; Move-Item -Path Scripts/$(BINARY_NAME)_Setup_$(VERSION).exe -Destination bin/$(BINARY_NAME)_Setup_$(VERSION).exe -Force } else { Write-Output '[WARNING] NSIS compilation failed. Setup exe will not be created.' }"
else
	RM = rm -rf bin/ $(BINARY_NAME).exe
	MKDIR = mkdir -p bin
	MKDIR_BUILD = mkdir -p build
	RM_BUILD = rm -rf build/
	MV_INSTALLER = mv Scripts/$(BINARY_NAME)_Setup_$(VERSION).exe bin/
	MAKENSIS = makensis
	GO_BUILD = GOOS=$(1) GOARCH=$(2) go build -ldflags "-X main.Version=$(VERSION)" -o $(3) .
	NSIS_STEP = \
		echo "[Windows] Compiling NSIS installer..."; \
		if $(MAKENSIS) Scripts/build.nsi; then \
			echo "[Windows] Moving installer to bin/..."; \
			$(MV_INSTALLER); \
		else \
			echo "[WARNING] NSIS compilation failed. Setup exe will not be created."; \
		fi
endif

run:
	@echo Running StellePlayer...
	go run .

build-amd64:
	@echo Building local binary...
	go build -o $(BINARY_NAME).exe .

build-all: 
	@echo Starting full cross-platform build...
	@$(MAKE) build-windows
	@$(MAKE) build-linux
	@$(MAKE) build-macos

	@echo Build complete. Check the bin/ folder.


build-windows:
	@echo "[Windows] Preparing directory..."
	@$(MKDIR)
	@$(MKDIR_BUILD)
	@echo "[Windows] Building standalone binary (amd64)..."
	@$(call GO_BUILD,windows,amd64,bin/$(BINARY_NAME)-windows-amd64.exe)
	@echo "[Windows] Building binary for installer..."
	@$(call GO_BUILD,windows,amd64,build/$(BINARY_NAME).exe)
	-@$(NSIS_STEP)
	@$(RM_BUILD)

build-linux:
	@echo "[Linux] Preparing directory..."
	@$(MKDIR)
	@echo "[Linux] Building amd64..."
	@$(call GO_BUILD,linux,amd64,bin/$(BINARY_NAME)-linux-amd64)
	@echo "[Linux] Building arm64..."
	@$(call GO_BUILD,linux,arm64,bin/$(BINARY_NAME)-linux-arm64)
	@echo "[Linux] Building arm..."
	@$(call GO_BUILD,linux,arm,bin/$(BINARY_NAME)-linux-arm)

build-macos:
	@echo "[MacOS] Preparing directory..."
	@$(MKDIR)
	@echo "[MacOS] Building amd64..."
	@$(call GO_BUILD,darwin,amd64,bin/$(BINARY_NAME)-darwin-amd64)
	@echo "[MacOS] Building arm64..."
	@$(call GO_BUILD,darwin,arm64,bin/$(BINARY_NAME)-darwin-arm64)

clean:
	@echo Cleaning up build artifacts...
	@$(RM)
	@$(RM_BUILD)
