# cloudflare-ddns

A Go daemon for updating CloudFlare DNS automatically.

## Usage

Simply invoke the binary, for example `./cloudflare-ddns` and you will be provided with a list of all options.

## How it works

By default, the daemon will make an API call to https://api.ipify.org for getting your external IPv4 address, or https://api6.ipify.org for your IPv6 address, and based on the parameters provided will send an API request to CloudFlare to update your DNS entry.

The parameter `-interface <name>` can be used if the IP you want the DNS entry to point to the unicast address of the interface instead of making an API call to ipify.org. That means, for IPv4 it will be (most likely) a private IP, and for IPv6 it will be a global unicast address.

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

The service can be run as a cron task on every hour by simply modifying the crontab and adding:
```
0 */1 * * * /path/to/cloudflare-ddns -token <token here> -domain <domain here> > /dev/null
```

## TODOs

- Allow configuration of ipify domain name
- Add developer environment
