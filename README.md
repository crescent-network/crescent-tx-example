# Crescent Transaction Example

This repository is a sample project to demonstrate how to wrap message(s) in a single transaction and sign to broadcast the transaction to the network. In `main.go`, you will see that there are some sample functions to create transactions messages.

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
