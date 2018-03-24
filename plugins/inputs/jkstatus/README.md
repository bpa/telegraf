# JKStatus Input Plugin

The JKStatus plugin collects statistics available from the jkstatus page from the `http://<host>/jkstatus?mime=xml URL.` (`mime=xml` will return only xml data).

### Configuration:

```toml
# Gather metrics from the jkstatus page.
[[inputs.jkstatus]]
  ## URL of the jkstatus page
  # url = "http://127.0.0.1:8080/jkstatus?mime=xml"

  ## HTTP Basic Auth Credentials
  # username = "jkstatus"
  # password = "s3cret"

  ## Request timeout
  # timeout = "5s"

  ## Optional SSL Config
  # ssl_ca = "/etc/telegraf/ca.pem"
  # ssl_cert = "/etc/telegraf/cert.pem"
  # ssl_key = "/etc/telegraf/key.pem"
  ## Use SSL but skip chain & host verification
  # insecure_skip_verify = false
```

### Measurements & Fields:

### Tags:

### Example Output:

```
```
