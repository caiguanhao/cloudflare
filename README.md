# cloudflare

Simple cloudflare command line tool.

```
cloudflare ls | jq -r '.[].name'
cloudflare records example.com | jq -r '.[] | "\(.id)\t\(.type)\t\(.content)\t\(.name)"'
cloudflare addrecord example.com test TXT helloworld
cloudflare delrecord example.com 73c939efa2eb425e95471c96ec67bc88
```
