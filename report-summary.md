# Firewall Escape Test - Run 21819295673

**Date**: 2026-02-09
**Status**: SANDBOX SECURE
**Techniques Tested**: 23
**Novelty Rate**: 95%
**Escapes Found**: 0

## Key Findings
- AWF_ENABLE_HOST_ACCESS=1 allows host.docker.internal bypass (by design)
- host.docker.internal (172.17.0.1) port 80 open, returns 404 (web service active)
- Squid version 6.13 confirmed
- CAP_NET_RAW properly dropped (no raw packet creation)
- /var/run/docker.sock exists but connection refused (proper isolation)
- All application-level bypass attempts (Python, Node.js) intercepted by iptables NAT
- HTTP/2 protocol not supported by Squid for direct connections
- FTP, SOCKS5, and other non-HTTP protocols properly blocked

## Architecture Strengths (Confirmed via Source Code Analysis)
1. **Triple-layer defense**: iptables NAT → Squid proxy → Host-level iptables
2. **Universal traffic redirection**: ALL TCP/UDP redirected to Squid (172.30.0.10:3128)
3. **DNS hardening**: Only trusted resolvers (8.8.8.8, 8.8.4.4, 127.0.0.11)
4. **Localhost bypass**: Intentionally allows 127.0.0.0/8 for MCP servers
5. **Host gateway bypass**: host.docker.internal bypasses Squid when AWF_ENABLE_HOST_ACCESS=1 (for MCP gateway)
6. **Capabilities properly dropped**: NET_RAW, NET_ADMIN, SYS_PTRACE all disabled
7. **Seccomp mode 2**: Syscall filtering active

## Reconnaissance Intel Gathered
- Gateway web service: 172.17.0.1:80 (404 response)
- Squid proxy: 172.30.0.10:3128 (version 6.13)
- Docker socket: /var/run/docker.sock (protected)
- Network namespace: net:[4026532411] (isolated)
- Shared memory: /dev/shm/ accessible (no exploit found)

## Recommendations
- Continue monitoring host.docker.internal bypass behavior
- Regular Squid version updates (currently 6.13)
- Consider logging DNS query patterns for anomaly detection
- Continue diverse security testing with novel techniques

## Cumulative Statistics
- Total techniques: 538 (22 runs)
- Historical escapes: 1 (patched in AWF v0.9.1)
- Success rate: 0.19% (1/538)
- Last 508 techniques: All blocked
