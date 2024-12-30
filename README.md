# README

This is just a little guy that allows you to quickly mute and unmute hosts in Datadog.

## Installation

```bash
go install github.com/hkrutzer/dd-mute-host@latest
```

## Configuration

The tool requires Datadog API credentials. Set the following environment variables:

- `DD_API_KEY`: Your Datadog API key
- `DD_APP_KEY`: Your Datadog application key

## Usage

### List Hosts
To see all hosts and their mute status:
```bash
dd-mute-host list
```

### Mute Hosts
To mute one or more hosts for a specified duration:
```bash
dd-mute-host mute [--duration minutes] host1 [host2 ...]
```

Options:
- `--duration`: How many minutes to mute the host for (default: 60)

Example:
```bash
# Mute a single host for 30 minutes
dd-mute-host mute --duration 30 webserver1

# Mute multiple hosts for the default duration (60 minutes)
dd-mute-host mute webserver1 webserver2 webserver3
```

### Unmute Hosts
To unmute one or more hosts:
```bash
dd-mute-host unmute host1 [host2 ...]
```

Example:
```bash
# Unmute a single host
dd-mute-host unmute webserver1

# Unmute multiple hosts
dd-mute-host unmute webserver1 webserver2 webserver3
```

