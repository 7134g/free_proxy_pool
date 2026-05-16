## MODIFIED Requirements

### Requirement: Crawl ip3366.net for proxies
**FROM:** The system crawls `http://www.ip3366.net/?stype=1&page={page}` for pages 1 through 99 (99 pages).
**TO:** The system SHALL crawl `http://www.ip3366.net/?stype=1&page={page}` for pages 1 through 9 (9 pages).

#### Scenario: Successful crawl returns proxies
- **WHEN** the crawler fetches a valid page from ip3366.net
- **THEN** the crawler SHALL extract IPs from `#list > table > tbody > tr td:nth-child(1)` and ports from `td:nth-child(2)`, scheme from `td:nth-child(4)`, and yield `{scheme}://{ip}:{port}` to ProxyChannel

#### Scenario: Page fetch fails
- **WHEN** the page returns non-200 status or connection error
- **THEN** the crawler SHALL skip that page and continue to the next page
