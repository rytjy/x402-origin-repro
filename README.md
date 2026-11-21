# x402 SIWX — configured-origin canonicalisation, offline repro

Reproduces the behaviour difference behind x402-foundation/x402 PR **#4328**
(`fix(go): canonicalize SIWX configured origin`, commit `6d160f9`) **without**
the upstream module: no `go.mod` toolchain download, no network, standard
library only.

## Run it

```bash
go test ./...        # must PASS  (~4 ms)
bash must_fail.sh    # must print MUTANT KILLED, exit 0
```

`must_fail.sh` copies the module to a temp dir, flips `var normalize` from
`normalizeFixed` to `normalizeBuggy` (the pre-fix `parsed.Host` behaviour) and
requires the same table to go RED. If the mutant survives, the script exits 1.

## Result (2026-11-21, go1.22.2, linux/amd64)

```
ok  x402origin  0.004s                       # fixed
--- FAIL: .../default_https_port_dropped     # buggy
--- FAIL: .../host_lowercased
--- FAIL: .../default_http_port_dropped
--- FAIL: .../ipv6_default_port
```

| input | pre-fix (`parsed.Host`) | canonical |
|---|---|---|
| `https://api.example.com:443` | `api.example.com:443` | `api.example.com` |
| `https://API.Example.com` | `API.Example.com` | `api.example.com` |
| `http://api.example.com:80` | `api.example.com:80` | `api.example.com` |
| `https://api.example.com:8443` | `api.example.com:8443` | `api.example.com:8443` |
| `https://[::1]:443` | `[::1]:443` | `[::1]` |

## 证明什么 / 不证明什么

**证明**

- 在给定的七个输入上，`parsed.Host` 与「小写 + 去默认端口 + IPv6 保括号」的规范化结果不同；差异是可枚举、可断言的。
- 反向对照成立：把实现换回 `parsed.Host`，同一张表立刻变红（`must_fail.sh`）。
- 这些差异与 PR #4328 的改动方向一致（该 PR 新增 `canonicalOriginHost` 与 +48 行 `validate_test.go`）。

**不证明**

- **不证明上游整包测试会红。** 本仓库是独立复刻，不是上游代码；它跑的是 `normalizeBuggy` / `normalizeFixed` 两个复刻函数。
- 不证明线上任何部署受影响。影响面是「配置里的 origin 字符串与请求比较时不相等」，是否可被利用取决于调用点。
- 不覆盖 IDN（`bücher.example` → punycode）：那需要 `golang.org/x/net/idna`，会破坏「仅标准库」这一条。
- 不替代 PR 里的测试：上游的 `validate_test.go` 才是交付物，这里是让人 3 秒看懂差异的旁证。

## 为什么不用上游测试

`x402` 的 `go.mod` 要求 **go1.26.10**；本机为 go1.22.2，`go test ./extensions/signinwithx/`
会触发工具链下载（无网络时挂起）。2026-11-19 已因此放弃本地跑整包，改交 CI。
本模块绕开这个问题：`go 1.22` + 零依赖 ⇒ 离线可跑。
