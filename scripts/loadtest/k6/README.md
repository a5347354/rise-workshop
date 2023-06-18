# How to use

## [install k6](https://k6.io/docs/get-started/installation/)

## Run
```shell
export RELAY_HOSTNAME=wss://relay-our-tide-385402-zzcjiygo2q-uc.a.run.app/ws
k6 run --vus 2 --duration 10s scripts/loadtest/k6/relay.js
```
