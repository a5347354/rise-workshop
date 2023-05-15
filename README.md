# Distributed Systems Project 

## How to run this Project on local
### Install tools
```bash
# auto watch file change and hot reload server
go get -u github.com/silenceper/gowatch

# install mockgen
go get -u github.com/golang/mock/mockgen

# install govulncheck
go install golang.org/x/vuln/cmd/govulncheck@latest
```

### View coverReport
```bash
go tool cover -html=coverage.out
```

## Phase 1
Building a simple nostr client

### Question Exercises
- What are some of the challenges you faced while working on Phase 1?
  - I didn't know Nostra very well, so it took me some time to understand it.
  - I spent some time debugging why there was no ACK response until I found out that our classmates had already discussed it on Discord. 


- What kind of failures do you expect to a project such as DISTRISE to encounter?
  - Misconfiguration caused DISTRISE to be unable to serve on time during a release change. 
  - High traffic has led to the service going down.