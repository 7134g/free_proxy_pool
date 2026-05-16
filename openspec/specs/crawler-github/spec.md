## ADDED Requirements

### Requirement: Fetch proxies from GitHub proxy list
The system SHALL fetch proxy list from `https://raw.githubusercontent.com/TheSpeedX/PROXY-List/master/http.txt` and parse the plain-text response.

#### Scenario: Successful fetch returns proxies
- **WHEN** the crawler receives a valid response from GitHub
- **THEN** the crawler SHALL split the response by newlines, filter lines matching `\d+\.\d+\.\d+\.\d+:\d+` pattern, prefix with `http://`, and yield each to ProxyChannel

#### Scenario: Fetch fails
- **WHEN** the raw URL returns non-200 or connection error
- **THEN** the crawler SHALL report the error and terminate without yielding any proxies
