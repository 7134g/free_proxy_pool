## Why

现有的 Go 代理池爬虫数量较少（仅 4 个），且部分爬虫（如 ip3366）的抓取页数与 Python 版不一致。将 Python 版已验证的代理源移植到 Go 版，可增加代理来源多样性，提高可用代理数量和质量。

## What Changes

- **新增 5 个爬虫**：89ip、66daili.com（非 66ip.cn）、ProxyScrape API、GitHub 代理列表、zdaye.com
- **修改 ip3366**：抓取页数从 99 页改为 9 页，与 Python 版一致
- **调整 register.go**：注册全部新爬虫
- **新增爬虫测试**：为每个新爬虫添加单元测试

## Capabilities

### New Capabilities

- `crawler-89ip`: 从 89ip.cn 抓取免费代理（HTML 解析）
- `crawler-66daili`: 从 66daili.com 抓取免费代理（正则解析）
- `crawler-proxyscrape`: 从 proxyscrape API 获取代理列表（纯文本解析）
- `crawler-github`: 从 GitHub 原始代理列表获取代理（纯文本解析）
- `crawler-zdaye`: 从 zdaye.com 抓取免费代理（需要 Cookie + 代理链，两阶段爬取）

### Modified Capabilities

- `crawler-ip3366`: 抓取页数从 99 页改为 9 页，与 Python 版对齐

## Impact

- `crawler/cell/` — 新增 5 个 `crawl_*.go` 文件，修改 `crawl_ip3366.go`、`register.go`
- `crawler/cell/crawl_test.go` — 新增对应测试
- 无外部依赖变更，新爬虫均复用现有 `xhttp` 和 `goquery` 库
