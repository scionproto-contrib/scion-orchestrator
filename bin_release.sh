#!/bin/bash
set -e

# Check if dev directory exists
if [ -d "dev" ]; then
    echo "Directory dev exists."
else 
    mkdir dev
    cd dev
    git clone https://github.com/scionproto/scion.git
    cd scion
    git checkout v0.12.0
    cd ../../
fi

build_and_copy_binaries() {
    cd dev/scion/
    local bin_dir="$1"

    # Determine the binary suffix based on the target OS
    local suffix=""
    if [ "$GOOS" == "windows" ]; then
        suffix=".exe"
    fi

    # Ensure bin directory exists
    mkdir -p "../../$bin_dir"

    # Array of relative paths for each command
    declare -a commands=(
        "scion/cmd/scion"
        "scion-pki/cmd/scion-pki"
        "router/cmd/router"
        "control/cmd/control"
        "daemon/cmd/daemon"
        "dispatcher/cmd/dispatcher"
    )

    echo "Changing to directory: $(pwd)"

    # Build each command with the specified environment variables
    for cmd_path in "${commands[@]}"; do
        local bin_name=$(basename "$cmd_path")
        echo "Building ${bin_name}${suffix}..."
        CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -o "${cmd_path}/${bin_name}${suffix}" "./${cmd_path}"
    done

    # Copy binaries to bin directory
    echo "Copying binaries to ../../$bin_dir/"
    cp "scion/cmd/scion/scion${suffix}" "../../$bin_dir/"
    cp "scion-pki/cmd/scion-pki/scion-pki${suffix}" "../../$bin_dir/"
    cp "router/cmd/router/router${suffix}" "../../$bin_dir/"
    cp "control/cmd/control/control${suffix}" "../../$bin_dir/"
    cp "daemon/cmd/daemon/daemon${suffix}" "../../$bin_dir/"
    cp "dispatcher/cmd/dispatcher/dispatcher${suffix}" "../../$bin_dir/"

    cd ../../

    # Final build for scion-orchestrator with environment variables, if needed
    echo "Building scion-orchestrator${suffix}..."
    CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -o "scion-orchestrator${suffix}"

    # Copy final binary to bin directory
    if [ -f "scion-orchestrator${suffix}" ]; then
        echo "Copying scion-orchestrator${suffix} to $bin_dir/"
        cp "scion-orchestrator${suffix}" "$bin_dir/"
    else 
        echo "Warning: scion-orchestrator${suffix} not found after build."
    fi
}

# --- Build for various platforms by uncommenting the desired lines ---

 GOOS=linux GOARCH=amd64 build_and_copy_binaries "./bin_release/linux_amd64"
# GOOS=linux GOARCH=arm64 build_and_copy_binaries "./bin_release/linux_arm64"
# GOOS=darwin GOARCH=amd64 build_and_copy_binaries "./bin_release/darwin_amd64"
# GOOS=darwin GOARCH=arm64 build_and_copy_binaries "./bin_release/darwin_arm64"
#GOOS=windows GOARCH=amd64 build_and_copy_binaries "./bin_release/windows_amd64"