## ADDED Requirements

### Requirement: Crawl zdaye.com for proxies
The system SHALL crawl `https://www.zdaye.com/dayProxy/{page}.html` for pages 1 through 4 with Cookie headers, then follow detail page links to extract proxy IP:port pairs.

#### Scenario: Successful crawl returns proxies
- **WHEN** the crawler fetches a valid list page from zdaye.com with Cookie headers
- **AND** finds detail page links matching `#J_posts_list .thread_item div div p a`
- **AND** fetches each detail page successfully
- **THEN** the crawler SHALL extract `ip:port` from text nodes in `.cont br` elements on detail pages, prefix with `http://`, and yield each to ProxyChannel

#### Scenario: List page fails
- **WHEN** the list page returns non-200 or connection error
- **THEN** the crawler SHALL skip that page and continue to the next page

#### Scenario: Detail page fails
- **WHEN** a detail page fetch fails
- **THEN** the crawler SHALL skip that detail page and continue to the next link
