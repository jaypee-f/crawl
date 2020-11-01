# Crawl
Crawl a domain. the crawler should visit each URL it finds on the same domain. 
It should print each URL visited, and a list of links found on that page

A timeout of 5 secs is given for retrieving page. 

Usage: go run cmd/main.go url

TODOs 
* Clean up command usage documentation.
* Configurable timeouts, concurrency.
* Observability on the network requests. 
* Sorted results.
* Finish when all consumers are idle