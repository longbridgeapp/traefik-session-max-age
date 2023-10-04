This repository includes traefik plugin which sets sticky session cookie's max age.

### Configuration

For traefik plugin, static configuration must define the module name (as is usual for Go packages).

The following declaration (given here in YAML) defines a plugin:

```yaml
# Static configuration

experimental:
  plugins:
    session-max-age:
      moduleName: github.com/longbridgeapp/traefik-session-max-age
      version: v0.1.0
```

Here is an example of a file provider dynamic configuration (given here in YAML), where the interesting part is the `http.middlewares` section:

```yaml
# Dynamic configuration

http:
  routers:
    my-router:
      rule: host(`demo.localhost`)
      service: service-foo
      entryPoints:
        - web
      middlewares:
        - session-max-age

  services:
    service-foo:
      loadBalancer:
        servers:
          - url: http://127.0.0.1:5000

  middlewares:
    session-max-age:
      plugin:
        session-max-age:
          cookieName: traefik_cookie
          maxAge: 100000
```
