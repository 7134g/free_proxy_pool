## ADDED Requirements

### Requirement: Fetch proxies from ProxyScrape API
The system SHALL fetch proxy list from `https://api.proxyscrape.com/v2/?request=displayproxies&protocol=http&timeout=10000&country=all&ssl=all&anonymity=all` and parse the plain-text response.

#### Scenario: Successful API fetch returns proxies
- **WHEN** the crawler receives a valid response from ProxyScrape API
- **THEN** the crawler SHALL split the response by newlines, filter lines matching `\d+\.\d+\.\d+\.\d+:\d+` pattern, prefix with `http://`, and yield each to ProxyChannel

#### Scenario: API fetch fails
- **WHEN** the API returns non-200 or connection error
- **THEN** the crawler SHALL report the error and terminate without yielding any proxies
