# GameBox

This is a game server that provides various casual games.
It is highly extensible and easy to add new games.

## Getting Started

### Prerequisites

- Golang
- Docker
- Google Cloud SDK

### Installation

```shell
# Prepare tools
make install-tools
docker-compose up -d

# Prepare application config
export GAMEBOX_CONFIG_FILEPATH=$(pwd)/config/default.json
# Optional: Use custom config file
cp config/default.json config/{custom_config_name}.json  
export GAMEBOX_CONFIG_FILEPATH=$(pwd)/config/{custom_config_name}.json
```
