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

- jkstatus_balancer
    name
    type
    sticky_session
    sticky_session_force
    retries
    recover_time
    error_escalation_time
    max_reply_timeouts
    method
    lock
    member_count
    good
    degraded
    bad
    busy
    max_busy
    map_count
    time_to_maintenance_min
    time_to_maintenance_max
    last_reset_at
    last_reset_ago
- jkstatus_member
    name
    type
    host
    port
    address
    connection_pool_timeout
    ping_timeout
    connect_timeout
    prepost_timeout
    reply_timeout
    connection_ping_interval
    retries
    recovery_options
    max_packet_size
    activation
    lbfactor
    route
    redirect
    domain
    distance
    state
    lbmult
    lbvalue
    elected
    sessions
    errors
    client_errors
    reply_timeouts
    transferred
    read
    busy
    max_busy
    connected
    time_to_recover_min
    time_to_recover_max
    last_reset_at
    last_reset_ago

### Tags:

- jkstatus_balancer
  - name
- jkstatus_member
  - balancer
  - name

### Example Output:

```
jkstatus_balancer,name=balancer type=lb,sticky_session=True,sticky_session_force=False,retries=2,recover_time=60,error_escalation_time=30,max_reply_timeouts=0,method=Busyness,lock=Optimistic,member_count=2,good=2,degraded=0,bad=0,busy=7,max_busy=16,map_count=2,time_to_maintenance_min=26,time_to_maintenance_max=88,last_reset_at=1521759602,last_reset_ago=348 1522776987705587000
jkstatus_balancer,balancer=balancer,name=server1 type=ajp13,host=worker-box1,port=8009,address=10.200.0.2:8009,connection_pool_timeout=0,ping_timeout=10000,connect_timeout=10000,prepost_timeout=10000,reply_timeout=0,connection_ping_interval=100,retries=2,recovery_options=0,max_packet_size=16384,activation=ACT,lbfactor=10,route=workerbox1,redirect=,domain=,distance=0,state=OK,lbmult=1,lbvalue=4,elected=1180,sessions=1180,errors=0,client_errors=0,reply_timeouts=0,transferred=1812895,read=293257002,busy=4,max_busy=8,connected=348,time_to_recover_min=0,time_to_recover_max=0,last_reset_at=1521759602,last_reset_ago=348 1522776987705587000
jkstatus_balancer,balancer=balancer,name=server2 type=ajp13,host=worker-box2,port=8009,address=10.200.0.3:8009,connection_pool_timeout=0,ping_timeout=10000,connect_timeout=10000,prepost_timeout=10000,reply_timeout=0,connection_ping_interval=100,retries=2,recovery_options=0,max_packet_size=16384,activation=ACT,lbfactor=10,route=workerbox2,redirect=,domain=,distance=0,state=OK,lbmult=1,lbvalue=3,elected=1193,sessions=1193,errors=0,client_errors=0,reply_timeouts=0,transferred=1749155,read=265104501,busy=3,max_busy=8,connected=341,time_to_recover_min=0,time_to_recover_max=0,last_reset_at=1521759602,last_reset_ago=348 1522776987705587000
```
