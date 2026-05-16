## ADDED Requirements

### Requirement: Crawl 66daili.com for proxies
The system SHALL crawl `http://www.66daili.com/?page={page}` for pages 1 through 10 and extract proxy IP:port pairs using regex.

#### Scenario: Successful crawl returns proxies
- **WHEN** the crawler fetches a valid page from 66daili.com
- **THEN** the crawler SHALL extract IPs by matching `<li>` tags with `min-width: 120px` and ports from `<li>` tags with `min-width: 60px`, and yield `http://{ip}:{port}` to ProxyChannel

#### Scenario: Page fetch fails
- **WHEN** the page returns non-200 status or connection error
- **THEN** the crawler SHALL skip that page and continue to the next page
