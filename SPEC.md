# TakeDoc
Its a service which enable x402 payments api for paid documents. the price is in USD  using tempo blockchain for it. 

## Tech stack
- echo
- typespec
- flag (for reading command line argument)
- sqllite
- grpc (to get details of available docs)

## Features
- Use Hexagonal Architecture (Ports & Adapters)
- Smoke testing
- health api (for rest service, db )
- save the log of who accessed the files to db such as blockchain address, ip, or other relavent details
- oberservability

## Rules
- use AskUserQuestion tool for any unsure features / task
- make phase-wise gated plan
- save the plan in same repo as md file
- use relavent skills 