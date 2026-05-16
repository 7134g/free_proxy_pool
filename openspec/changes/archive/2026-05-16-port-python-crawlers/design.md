## Context

Go 代理池项目中爬虫位于 `crawler/cell/` 包下。每个爬虫实现 `spider` 接口（`name()`、`genSeek()`、`run()`、`parse()`），通过 `register.go` 的 `Crawler()` 函数注册启动。现有爬虫使用 `xhttp.Get()` 发起 HTTP 请求，使用 `goquery` 解析 HTML 返回 `ip:port`。

Python 版已验证的 5 个新代理源中，有 4 个可复用现有模式（HTTP GET + 解析），1 个（zdaye）需要 Cookie 和代理链支持。

## Goals / Non-Goals

**Goals:**
- 新增 5 个爬虫：89ip、66daili.com、ProxyScrape、GitHub 代理列表、zdaye.com
- ip3366 抓取页数修正为 9 页
- 所有新爬虫注册到 Crawler() 启动流程
- 每个爬虫有对应的单元测试

**Non-Goals:**
- 不修改现有爬虫的已有逻辑（除 ip3366 页数外）
- 不新增外部依赖
- 不修改 `xhttp` 或 `store` 等底层模块
- 不涉及代理验证逻辑

## Decisions

1. **89ip 使用 goquery HTML 解析** — 与现有爬虫一致，`td:nth-child(1)` 取 IP、`td:nth-child(2)` 取端口。URL 格式 `https://www.89ip.cn/index_{page}.html`，爬取 6 页

2. **66daili.com 使用正则解析** — Python 版使用 `re.findall` 从 `<li>` 提取 IP 和端口。由于 66daili.com 页面结构不规整（非标准 table），goquery 难以准确定位，保留正则方式。在 Go 中用 `regexp` 包处理 `[]byte`

3. **ProxyScrape 和 GitHub 使用行解析 + 正则** — 返回纯文本格式，每行一个 `ip:port`。直接用 `bytes.Split` 按换行分割 + 正则过滤。无需 goquery

4. **zdaye 需要 Cookie + 代理链** — Python 版在请求 zdaye 时带 Cookie 头 + 使用代理池中的代理。Go 端：
   - 使用 `xhttp.GetHeader()` 传 Cookie
   - 代理链方面：zdaye 本身需要先有一个可用代理才能抓取。实现上可以在 `genSeek()` 中从 `Store` 读取可用代理，通过 `xhttp.SetLocalProxy()` 设置，但这需要爬虫能访问 store。另一种简化方案：允许 zdaye 不强制使用代理，先直接请求，失败则跳过

5. **ip3366 页数调整** — Python `range(1,10)` = 9 页，直接将 Go 中 `page < 100` 改为 `page < 10`

## Risks / Trade-offs

- **zdaye 依赖代理**：zdaye 可能反爬，如果没有可用代理会失败。备选方案：直接请求不设代理，能抓到最好，抓不到不强求
- **66daili 正则脆弱性**：网站改版可能导致正则失效。但 Python 版已在生产验证，且正则模式匹配的是 `min-width` 样式属性，相对稳定
- **ProxyScrape/GitHub API 变化**：外部 API URL 可能变更。这是所有爬虫的共同风险
