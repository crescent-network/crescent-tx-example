# Crescent Transaction Example

This repository is a sample project to demonstrate how to wrap message(s) in a single transaction and signs the transaction to broadcast it the transaction to the network.

In `main.go`, you will find that it demonstrates how to send `MsgMMOrder` for market makers. There are also other sample functions available. Feel free to add/update/delete them depending on your usage.

## Dependency

| Dependency    | Version |
| ------------- | ------- |
| Go            | 1.18    |
| Crescent Core | 3.0.x   |

## Configuration

A configuration file called `config.toml` is required to run the program. Make sure to copy `example.toml` to create `config.toml` file and change values for you need. The config source code can be found in [this config file](/config/config.go).

## Usage

```bash
# See the code in main.go file to understand how this program is developed
go run main.go
```

## Resources

- [Network Configurations for Crescent Mainnet and Testnet(s)](https://docs.crescent.network/other-information/network-configurations)
