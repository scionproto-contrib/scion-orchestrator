#!/bin/bash

# Quick and dirty script for testing the orchestrator's installation feature.
# Build the orchestrator with go build and run this in the repository's root.
# To clean up, delete the directories ./config and ./bin

set -euo pipefail

if [ ! -e config ] && [ ! -L config ]; then
ln -s release/leaf-as/config config
fi

if [ ! -d bin ]; then
mkdir bin
curl -fL -o bin/scion.tar.gz "https://github.com/scionproto/scion/releases/download/v0.12.0/scion_0.12.0_amd64_linux.tar.gz"
echo "65ca69fee123cb580640bd40396abc74706d84e51f01618c3a50bbdff18344ae  bin/scion.tar.gz" | sha256sum -c - || exit 1
tar -xzf bin/scion.tar.gz -C bin
rm bin/scion.tar.gz
fi

docker run -it --rm --detach --name orchestrator-test --network host -v "$PWD:/src:ro" ubuntu:24.04 bash
docker exec -i orchestrator-test bash << 'EOF'
printf '#!/bin/sh\nexit 0\n' > /usr/local/bin/systemctl && chmod +x /usr/local/bin/systemctl
mkdir -p /opt/bundle
cp /src/scion-orchestrator /opt/bundle
cp -ar /src/bin /opt/bundle
cp -Lr /src/config /opt/bundle
EOF

echo "Dropping you into a temporary container, test installation with one of:"
echo "cd /opt/bundle && ./scion-orchestrator install # CLI install"
echo "cd /opt/bundle && ./scion-orchestrator # Install Wizard at https://127.0.0.1:8843"
docker attach orchestrator-test
