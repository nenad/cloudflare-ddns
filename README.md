# cloudflare-ddns

A Go daemon for updating CloudFlare DNS automatically.

## Usage

Simply invoke the binary, for example `./cloudflare-ddns` and you will be provided with a list of all options.

## Building

To build the binary, run `go build -o cloudflare-ddns cmd/cloudflare-ddns/main.go`

## Configuration

The daemon must have few parameters configured before it can be useful.

Here is a table of required and non-required parameters:

| Variable      | Explanation |  Required | Default |
| ------------- | ------------- | ------ | ----- |
| -token  | CloudFlare API token, must allow Zone.Zone, Zone.DNS permissions| Yes | |
| -domain  | Domain to be updated  | Yes | |
| -type  | IP version, allowed A for IPv4 or AAAA for IPv6  | No | A |
| -timeout  | Timeout for HTTP calls to CloudFlare  | No | 10s | 
| -ttl  | TTL for CloudFlare record  | No | 1 |
| -proxied  | If record should be proxied to CloudFlare, default true  | No | true | 
| -interface  | Network interface name, if provided will be used to retrieve IP address | No | |
| -cache  | Should the last record from CloudFlare be cached on disk | No | false | 

## Periodic tasks

_TODO_

- Add systemd service file

## TODOs

- Allow configuration of ipify domain name
- Add developer environment
- Add compiled builds
