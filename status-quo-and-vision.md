# SCION Orchestrator: Status Quo & Vision

## 1. Executive Summary: Core Vision
The ultimate goal of the `scion-orchestrator` is to provide a **"zero-effort AS operator experience."** Running a SCION Autonomous System (AS) should be as simple and frictionless as possible. 
*   **Short-term goal:** Flawless operation of single-node ASes.
*   **Long-term goal:** Seamless multi-node setups and intra-AS centralized management.

To achieve this, the orchestrator must shift its focus heavily toward robust infrastructure management, reliable service observability, and automated lifecycle handling (updates, backups, certificate rotations) while explicitly deprecating non-core features like end-host logic.

## 2. Status Quo (Current Capabilities)
Currently, `scion-orchestrator` is a cross-platform Go tool (MacOS, Linux, Windows) that functions as a host connectivity bootstrap and AS lifecycle manager. 

**Core Features Implemented:**
*   **Execution Modes:** `run` (standalone processes), `install` (system services), and `shutdown`.
*   **Certificate Authority (CA):** Built-in dedicated CA server to issue SCION AS certificates with JWT authentication for Control Services.
*   **Automated Cert Renewal:** Background cron-like process for renewing AS certificates when approaching expiration.
*   **Bootstrapping & Metrics:** Built-in bootstrapping server (`:8041`) and telemetry integration.

**Current Architectural Bottlenecks (Identified via Issue Tracker):**
*   Service execution requires root privileges instead of a dedicated `scion` user (Issue #11).
*   Certificate management lacks resilience (e.g., handling future-valid certs, multi-CA fallbacks) (Issues #14, #15, #16).
*   Deployment is manual (building from source/binaries) instead of using standard package managers (Issue #17).
*   Scope creep: End-host logic is polluting the infrastructure focus (Issue #23).

## 3. Strategic Roadmap & Future Work

### 3.1. Infrastructure & Deployment (Stability First)
*   **Deprecate End-Host Logic:** Remove SCION end-host logic (Issue #23) to keep the focus strictly on AS infrastructure.
*   **Package Management:** Move away from manual binary shipping. Integrate with standard package managers (e.g., Debian packages at `netsec-ethz`) (Issue #17).
*   **Service Hardening:** Run systemd units as the `scion` user, not root (Issue #11). Improve service functionalities (observability, automated restarts, config backups/restores).
*   **Gateway Integration:** Add support for SCION IP Gateway and/or SCITRA Tun (Issue #22).

### 3.2. Advanced Certificate Management
*   **Multi-CA & Fallback:** Support certificate renewal to multiple core ASes to ensure high availability (Issue #16).
*   **Future Certificates:** Enable the storage and loading of certificates that become valid in the future (Issue #15).
*   **Offline Renewal:** Provide a mechanism for offline renewals (using standard IP to connect to a central instance) when SCION connectivity to the CA drops.
*   **Key Management:** Allow operators to upload their own AS keys via the API/UI (Issue #10).

### 3.3. Telemetry, Monitoring & Alerting
*   **Minimalist Alerting System:** Implement a highly pragmatic monitoring approach. Admins provide an email address and receive alerts *only* for critical state changes: Link down, Service down, Certificate expired (Issue #20).
*   **Metrics Integration:** Deploy `node_exporter` alongside the orchestrator to expose standard infrastructure metrics (Issue #5).

### 3.4. Intra-AS Topology & Node Management
*   **Topology Synchronization:** Implement sync mechanisms between AS-local orchestrator instances (Issue #18).
*   **Intra-AS Central Management:** Extend the API for centralized management of multiple orchestrator nodes *within* the same AS. (Note: No global inter-AS central instance is desired). Stick to local UI for single-node AS setups.

### 3.5. User Interface (UX)
*   **Complete UI Workflow:** Finish open UI work (Issue #19) and provide a dedicated Settings page (Issue #21).
*   **Bugfixes:** Fix connectivity page in `run` mode (Issue #9), improve placeholder UX (Issue #12), and verify SCION interface editing (Issue #13).

