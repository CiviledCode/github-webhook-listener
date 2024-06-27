# Webhook Listener
Webhook listener will listen for an incoming request from a git webhook and process the action based on a json configuration file `config.json`. This is written with no dependencies in a single file in <100 lines of code. This software works under the assumption that local files cannot be written/read from publicly.

## Config Structure
- `ip`: The local IP that the http server will listen on.
- `port`: The local port that the http server will bind to.
- `webhooks`: A map string->endpoint which maps URLs to their respective endpoint configurations.

**Endpoint Config**
- `type`: The type of endpoint action. Only `command` has been implemeneted.
- `data`: Data specific to the action at hand. This is the command when the `command` action is specified.
- `secret_token`: Secret authority token that is validated before an action's interpreted.


### Example Config

```json
{
    "ip": "localhost",
    "port": 9001,
    "webhooks": {
      "/update_writeups": {
        "type": "cmd",
        "data": "/var/scripts/update_writeups.sh",
        "secret_token": "L1VKYopzL7SqBxtc8Grv9FEhwlb4fF0Q"
      },
      "/update_home": {
        "type": "command",
        "data": "/var/scripts/update_home.sh"
      }
    }
  }
```