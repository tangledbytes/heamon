# Heamon - A minimal Monitoring System

![GitHub issues](https://img.shields.io/github/issues/utkarsh-pro/heamon)
![GitHub](https://img.shields.io/github/license/utkarsh-pro/heamon)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/utkarsh-pro/heamon)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/utkarsh-pro/heamon)

<div style="display:flex;margin: 2rem 0">
    <div style="max-width:40%">
        <image src="./docs/images/heamon_ui.png">
    </div>
    <div style="padding:0 1rem;">
        <div>Heamon is minimal monitoring system written in Go.</div>
        <ol>
            <li>Hot Configuration Reloading</li>
            <li>Configurable Email alerts on service degradation/failure</li>
            <li>Minimal UI</li>
        </ol>
        <div>
            Heamon has an extensible event driven architecture which allows writing plugins for it extremely easy. Email Notifications is implemented leveraging the events generated within Heamon.
        </div>
    </div>
</div>

## Sample Heamon Configuration

```yaml
title: Heamon Application # Title will be used on the Heamon UI
port: 3000 # Heamon will start on this port
authentication:
  username: ${{.HEAMON_USER}} # Will be replaced with env var $HEAMON_USER
  password: ${{.HEAMON_PASS}} # Will be replaced with env var $HEAMON_PASS
monitor:
  interval: 5 # Interval in minutes at which heamon probe-bots will query the service
  services:
    - name: Service 1 # Name of the service
      host: <service-host> # Host name, for Sagacious it will be sagacious.dev
      health_check_endpoint: / # The endpoint which will be queried for health check
      failure: 90 # Service state will be marked as "FAIL" if more than 90% requests fails
      degraded: 80 # Service state will be marked as "DEGRADED" if more than 80% requests fail
      initial_down_time: 5 # Expected initial down time for the service. The service will be marked in unknown state for this duration
    - name: Sagacious
      host: sagacious.dev
      health_check_endpoint: /
      failure: 90
      degraded: 80
      initial_down_time: 5
    # Any number of services can be placed here...

# Plugins are optional features of heamon and can be turned off by removing them from the
# the config file
plugins:
  alert:
    email:
      smtp:
        host: smtp.sendgrid.net
        port: 25
        username: ${{.SENDGRID_USERNAME}}
        password: ${{.SENDGRID_API}}
      from: <alert@sagacious.dev>
      to:
        - utkarsh@sagacious.dev
        - example@example.com
      duration: 60 # Duration to wait before triggering another mail. This will not be respected if the application
                   # state changes
      
```

## Roadmap

- v0.2.0 - Add Auto Service Discovery for Docker Swarm
- v0.3.0 - Add Auto Service Discovery for Kubernetes
- v0.4.0 - Add support for either WASM or Go Plugins (Hashicorp Go Plugins)
- v0.5.0 - Robust Test Suite for Heamon