---
name: container
description: Guide for using Apple's Container project - a macOS-native containerization platform that runs Linux containers in lightweight VMs. Use when working with the `container` CLI for building, running, and managing OCI containers on macOS. Covers container lifecycle, image management, networking, volumes, registry operations, and system configuration.
---

# Container

Apple's Container project provides macOS-native containerization using lightweight VMs (one per container) via the Virtualization framework. Each container has full VM isolation with minimal resource usage and fast boot times.

## Quick Reference

### System Setup

```bash
# Start services (installs kernel on first run)
container system start

# Stop services
container system stop

# Optional: Set up local DNS domain
sudo container system dns create test
container system property set dns.domain test
```

### Running Containers

```bash
# Run interactively
container run -it ubuntu:latest /bin/bash

# Run detached with port publishing
container run -d --name web -p 127.0.0.1:8080:80 nginx:latest

# Run with resources and auto-remove
container run --rm --cpus 4 --memory 8g my-image

# With volume mount
container run -v ~/data:/data my-image

# With SSH agent forwarding
container run -it --ssh alpine:latest
```

### Building Images

```bash
# Basic build
container build -t my-app:latest .

# Multi-platform build
container build --arch arm64 --arch amd64 -t my-app .

# Configure builder resources for large builds
container builder start --cpus 8 --memory 32g
```

### Container Management

```bash
container ls                      # List running
container ls -a                   # List all
container stop <id>               # Stop gracefully
container rm <id>                 # Delete
container exec -it <id> sh        # Execute command
container logs <id>               # View logs
container logs --boot <id>        # View boot logs
container stats <id>              # Resource usage
container inspect <id>            # Detailed info (JSON)
```

### Image Management

```bash
container image ls                # List images
container image pull <ref>        # Pull from registry
container image push <ref>        # Push to registry
container image tag <src> <dst>   # Tag image
container image rm <ref>          # Delete image
container image inspect <ref>     # Detailed info
```

### Networks (macOS 26+)

```bash
container network create foo                              # Create network
container network create foo --subnet 192.168.100.0/24   # With custom subnet
container run --network foo my-image                      # Use network
container network ls                                      # List networks
container network rm foo                                  # Delete network
```

### Volumes

```bash
container volume create mydata        # Create volume
container run -v mydata:/data image   # Use volume
container volume ls                   # List volumes
container volume rm mydata            # Delete volume
```

### Registry Auth

```bash
container registry login registry.example.com
container registry logout registry.example.com
```

## Key Concepts

**Architecture**: Unlike shared-VM approaches, Container runs each container in its own lightweight VM using the Virtualization framework. This provides VM-level isolation with container-like resource efficiency.

**Networking**: Containers attach to vmnet networks. The default network (192.168.64.0/24) is created at startup. Container-to-container communication works on macOS 26+. Host access uses the gateway IP (192.168.64.1).

**Builder**: Image builds use BuildKit in a utility container. Start with `container builder start` or let it auto-start on first build.

**DNS**: Optional local DNS (`sudo container system dns create <domain>`) enables hostname resolution like `my-container.test`.

## Detailed References

- **[Command Reference](references/command-reference.md)**: Complete CLI documentation with all options
- **[How-To Guide](references/how-to.md)**: Common tasks and workflows
- **[Technical Overview](references/technical-overview.md)**: Architecture and limitations

## Limitations

1. **Container-to-host networking**: No direct access to host's 127.0.0.1. Workaround: use `socat` to forward from gateway (192.168.64.1)
2. **Memory ballooning**: Memory freed by containers isn't returned to macOS. Restart containers to reclaim memory
3. **macOS 15**: Limited functionality - no container-to-container networking, no custom networks, potential IP address issues
