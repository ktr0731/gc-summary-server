# Groove Coaster Summary
[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy)  
Receive your last result in Groove Coaster 3.  
===

## Description  
Groove Coaster Summary generates summary of last result in Groove Coaster 3.  
Also GC-Summary listening in HTTPS.  
If accessed to the server, it returns response with summary plain text.  

## Equipments
- Heroku (or any PaaS, servers)
- Redis

## Installation
If use Heroku, please click Heroku deploy button!  

Manual.  
``` go
$ go get github.com/lycoris0731/gc-summary-server
```

## Usage
If you want to deploy except Heroku, you need to set environment variables.  
``` sh
export NESICA_CARD_ID=""  # Your NESiCA ID
export NESICA_PASSWORD="" # Your NESiCA ID password
export PORT=""            # Port for listening
export REDIS_URL=""       # Redis server URL
```

Start server.  
``` go
$ ./gc-summary-server
```

## Also
[Linka (Slack bot)](https://github.com/lycoris0731/linka)

## License
Please see [LICENSE](./LICENSE).
