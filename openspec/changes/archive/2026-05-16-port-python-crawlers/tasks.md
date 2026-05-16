## 1. 修改现有爬虫

- [x] 1.1 修改 crawl_ip3366.go：genSeek 中循环改为 `page < 10`（与 Python 版 9 页对齐）
- [x] 1.2 更新 crawl_ip3366 测试，验证页数范围正确

## 2. 新增 89ip 爬虫

- [x] 2.1 创建 crawl_89ip.go：实现 spider 接口，genSeek 生成 6 页 URL，parse 用 goquery 解析 table tbody tr
- [x] 2.2 在 crawl_test.go 中添加 TestCrawler_crawl89ip

## 3. 新增 66daili.com 爬虫

- [x] 3.1 创建 crawl_66daili.go：genSeek 生成 10 页 URL，parse 用 regexp 从 `<li>` 提取 IP 和端口
- [x] 3.2 在 crawl_test.go 中添加 TestCrawler_crawl66DaiLi

## 4. 新增 ProxyScrape 爬虫

- [x] 4.1 创建 crawl_proxyscrape.go：genSeek 设置 API URL，parse 按行分割 + 正则过滤 ip:port
- [x] 4.2 在 crawl_test.go 中添加 TestCrawler_crawlProxyScrape

## 5. 新增 GitHub 代理列表爬虫

- [x] 5.1 创建 crawl_github.go：genSeek 设置 raw URL，parse 按行分割 + 正则过滤 ip:port
- [x] 5.2 在 crawl_test.go 中添加 TestCrawler_crawlGitHub

## 6. 新增 zdaye.com 爬虫

- [x] 6.1 创建 crawl_zdaye.go：genSeek 生成 4 页 URL，parse 列表页提取详情链接 + 抓取详情页解析 IP:PORT
- [x] 6.2 在 crawl_test.go 中添加 TestCrawler_crawlZdaye

## 7. 注册新爬虫

- [x] 7.1 在 register.go 的 Crawler() 中注册所有新爬虫
- [x] 7.2 运行 `go build .` 确认编译通过
- [x] 7.3 运行 `go test ./crawler/cell/` 确认所有爬虫测试通过
