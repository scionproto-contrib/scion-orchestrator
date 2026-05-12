# Getting Started with scion-orchestrator

A guide for new contributors covering dev setup, running examples, integration testing, and debugging.

## Prerequisites

- **Go >= 1.23** — check with `go version`
- **git**
- **Docker + Docker Compose** — only needed for integration tests
- **Linux/macOS** recommended; Windows is supported via PowerShell scripts

---

## 1. Clone and build

```sh
git clone https://github.com/scionproto-contrib/scion-orchestrator.git
cd scion-orchestrator
```

The orchestrator needs SCION binaries (`daemon`, `dispatcher`, `router`, `control`, `scion`, `scion-pki`) alongside it at runtime. The `dev.sh` script clones SCION v0.12.0 into `dev/`, builds those binaries, and outputs them to `bin/`:

```sh
# Linux / macOS
./dev.sh

# Windows (PowerShell)
.\dev.ps1 -Repository https://github.com/scionproto/scion.git -Tag v0.12.0 -Path dev
```

After this completes you have:
- `./bin/` — SCION binaries
- `./scion-orchestrator` — the orchestrator binary

To rebuild just the orchestrator after making code changes:

```sh
CGO_ENABLED=0 go build
```

---

## 2. Project layout

```
main.go                 # Entry point — CLI arg parsing, command dispatch
standalone.go           # "run" command: starts SCION processes directly
install.go              # "install" command: registers SCION as system services
conf/                   # Config loading (TOML)
pkg/
  apiv1/                # REST API v1 (Gin handlers)
  bootstrap/            # Bootstrap server (serves topology to endhosts) and client
  scionca/              # Certificate Authority implementation
  metrics/              # Prometheus + status HTTP endpoints
  osutils/              # systemd (Linux) / launchd (macOS) integration
ui/                     # Web UI (Gin, HTML templates, static assets)
examples/               # Ready-to-use configs for common topologies
integration/            # Docker Compose environment for end-to-end tests
doc/                    # API and bootstrap documentation
```

The orchestrator reads a `./config/` directory at startup and a `./bin/` directory containing SCION binaries — both relative to the working directory where you run it.

---

## 3. Run an example locally

Pick an example from `examples/` and copy it to `./config/`:

```sh
# Full core AS with integrated CA (good starting point)
mkdir -p config
cp -R examples/core-as-ca/* ./config/

sudo ./scion-orchestrator run
```

> `sudo` is required because SCION components open raw sockets and bind to low ports.

Other example topologies:

| Directory | What it demonstrates |
|---|---|
| `core-as-ca` | Core AS running control service + border router + CA server |
| `leaf-as` | Leaf AS connecting to a parent core AS |
| `endhost` | Minimal endhost bootstrapping from an AS |
| `ca-only` | Standalone CA server without AS infrastructure |
| `core-as-pila` | Core AS with PILA integration |

Each example's `scion-orchestrator.toml` is the main config file. The key fields:

```toml
command = "service"       # "service" or "standalone"
isd_as  = "1-150"         # ISD-AS identifier for this node
mode    = "as"            # "as" or "endhost"

[bootstrap]
server = "127.0.0.1:8041" # AS: address to serve bootstrap on
                           # Endhost: address of AS bootstrap server

[ca]
server  = ":3000"
clients = ["123:client1.secret"]   # clientId:path-to-shared-secret

[api]
users = ["admin:admin.secret"]     # username:path-to-password-file

[metrics]
prometheus = "127.0.0.1:33401"
```

Once running you can reach:
- **Web UI**: `http://localhost:8180`
- **REST API**: `https://localhost:8843/api/v1/` (HTTP Basic Auth with credentials from config)
- **Prometheus metrics**: `http://localhost:33401/metrics`

---

## 4. Integration tests

Integration tests spin up two full ASes (1-150 and 1-151) and one endhost inside Docker containers and run the orchestrator binary inside each.

### Build for Linux/amd64 (required for the containers)

```sh
./dev_integration.sh
```

This builds SCION binaries with `GOOS=linux GOARCH=amd64`, copies them to `integration/bin/`, and copies the orchestrator binary to `integration/scion-orchestrator`.

### Start the environment

```sh
docker compose -f integration/docker-compose.yml up -d
```

Three containers start:

| Container | Role | IP |
|---|---|---|
| `as150full` | Core AS 1-150 | 10.150.0.254 |
| `as151full` | Leaf AS 1-151 | 10.150.0.251 |
| `host150full` | Endhost in AS 1-150 | 10.150.0.252 |

Each container mounts its config from `integration/AS150/`, `integration/AS151/`, or `integration/Host150/` and runs `./scion-orchestrator run`.

### Watch logs

```sh
docker logs -f as150full
docker logs -f host150full
```

### Tear down

```sh
docker compose -f integration/docker-compose.yml down
```

---

## 5. Making and verifying changes

After editing Go code:

```sh
# Rebuild
CGO_ENABLED=0 go build

# Run unit tests (no special setup needed)
go test ./...

# Rebuild for integration (cross-compile for containers)
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build && cp scion-orchestrator ./integration/

# Restart containers to pick up the new binary
docker compose -f integration/docker-compose.yml restart
```

---

## 6. Debugging

### Verbose logging

The orchestrator uses standard Go `log` output. Run it in the foreground and watch stdout/stderr:

```sh
sudo ./scion-orchestrator run 2>&1 | tee orchestrator.log
```

### Inspect running SCION processes

The orchestrator spawns SCION sub-processes. List them:

```sh
ps aux | grep -E 'daemon|dispatcher|router|control'
```

### REST API

The API is documented in [doc/api/Readme.md](doc/api/Readme.md). Quick test with curl:

```sh
# List CSRs (adjust host/port and credentials)
curl -u admin:$(cat config/admin.secret) https://localhost:8843/api/v1/cppki/csr -k
```

### Bootstrap server

The bootstrap server docs are in [doc/bootstrap/Readme.md](doc/bootstrap/Readme.md). You can verify it is serving topology info:

```sh
curl http://localhost:8041/topology
```

### Inside a container

```sh
docker exec -it as150full bash
# SCION binaries are at /root/AS150/bin/
/root/AS150/bin/scion --help
```

### Common issues

| Symptom | Likely cause |
|---|---|
| `permission denied` on startup when running install | Missing `sudo` — orchestrator needs root access to write systemd files |
| `bin/daemon: no such file` | `dev.sh` (or `dev_integration.sh`) hasn't been run yet |
| Bootstrap client fails to connect | Endhost `bootstrap.server` address doesn't match the AS IP |
| CA signing fails | Shared secret file path in `[ca] clients` is wrong or missing |

---

## 7. Further reading

- [README.md](README.md) — usage reference and platform support matrix
- [doc/api/Readme.md](doc/api/Readme.md) — REST API endpoints with curl examples
- [doc/bootstrap/Readme.md](doc/bootstrap/Readme.md) — bootstrap server configuration
- [status-quo-and-vision.md](status-quo-and-vision.md) — roadmap and architectural direction
- [SCION documentation](https://docs.scion.org) — background on the SCION network protocol
