# cloudflare-ddns

A Go daemon for updating CloudFlare DNS automatically.

## Usage

Simply invoke the binary, for example `TOKEN=<CF Token> TYPE=AAAA DOMAIN="nenad.dev" ./cloudflare-ddns`.

## Configuration

The daemon must have few parameters configured from the environment.

Here is a table of required and non-required parameters:

| Variable      | Explanation |  Required |
| ------------- | ------------- | ------ |
| TOKEN  | CloudFlare API token, must allow Zone.Zone, Zone.DNS permissions| Yes |
| DOMAIN  | Domain to be updated  | Yes |
| TYPE  | IP version, allowed A for IPv4 or AAAA for IPv6  | Yes |
| TIMEOUT  | Timeout for HTTP calls to CloudFlare  | No |
| PROXIED  | If record should be proxied to CloudFlare, default true  | No |

## Periodic tasks

_TODO_

- Add systemd service file

## TODOs

- Allow configuration of ipify domain name
- Add developer environment
- Add compiled builds
