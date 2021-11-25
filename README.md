# go-apache-test

Apache2.4 test CLI - Test your vhosts

Purpose is mainly to test some rewrite mod directives.
If you want to go further, open an issue.

## Install

Requirements: docker

```
docker pull jbustin1/go-apache-test
```

## Usage

```sh
docker run --rm \
-v ${PWD}/tests:/usr/go-apache-test/tests \
-v ${PWD}/vhosts:/usr/go-apache-test/vhosts \
-v ${PWD}/config.json:/usr/go-apache-test/config.json \
jbustin1/go-apache-test
```

- `vhosts` directory should contain all conf apache files (named as [vhostfilename].conf)
- `tests` directory should contain all json test files (named as [vhostfilename].json)
- `config.json` is a configuration file

### Precision about vhosts

Your vhosts should work without any external resources (like SSL certificates, docroot, ...).

Or you should share this resources with apache2 by mounting them in docker command.

Example:

```sh
docker run --rm \
-v ${PWD}/tests:/usr/go-apache-test/tests \
-v ${PWD}/vhosts:/usr/go-apache-test/vhosts \
-v ${PWD}/config.json:/usr/go-apache-test/config.json \
-v ${PWD}/certificates:/tmp/certificates \
-v /path/to/docroot:/tmp/docroot \
jbustin1/go-apache-test
```

## Json

```json
{
  "test": {
    "fetchTimeoutMs": 3000
  },
  "debug": {
    "logLevel": "info"
  },
  "httpd": {
    "host": "httpd.front.com",
    "port": "80",
    "backends": ["unique.backend.com"],
    "backendsPort": "3011"
  }
}
```

- `fetchTimeoutMs` time in ms before test fails
- `logLevel` string ("error", "info", "debug")
- `httpd.host` apache hostname (no need to modify)
- `httpd.port` apache default active port (no need to modify)
- `httpd.backends` array of string
- `httpd.backendsPort` unique port for all backends inside go-apache-test (no need to modify)

Most important key is `backends`.

You should list all the backends expected by the vhost proxies.
For example, if I set a proxy balancer to route my trafic to a specific backend, I need to declare it in configuration.

```
  <Proxy balancer://nodejscluster>
    ProxySet lbmethod=bybusyness
    BalancerMember http://proxy.new-stack.com:3011 route=new-stack
    Proxyset stickysession=JSESSIONID
  </Proxy>
```

```json
  "httpd": {
    "host": "httpd.front.com",
    "port": "80",
    "backends": ["proxy.new-stack.com"],
    "backendsPort": "3011"
  }
```
