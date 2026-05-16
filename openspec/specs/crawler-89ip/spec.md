## ADDED Requirements

### Requirement: Crawl 89ip.cn for proxies
The system SHALL crawl `https://www.89ip.cn/index_{page}.html` for pages 1 through 6 and extract proxy IP:port pairs.

#### Scenario: Successful crawl returns proxies
- **WHEN** the crawler fetches a valid page from 89ip.cn
- **THEN** the crawler SHALL extract IPs from `table tbody tr td:nth-child(1)` and ports from `table tbody tr td:nth-child(2)` and yield `http://{ip}:{port}` to ProxyChannel

#### Scenario: Page fetch fails
- **WHEN** the page returns non-200 status or connection error
- **THEN** the crawler SHALL skip that page and continue to the next page
