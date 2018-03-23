package jkstatus

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/influxdata/telegraf/testutil"

	"github.com/stretchr/testify/require"
)

var jkstatus = `<?xml version="1.0" encoding="UTF-8"?>
<jk:status xmlns:jk="http://jkstatus.apache.org">
  <jk:server name="my-server" port="8084"/>
  <jk:time datetime="20180322162127" tz="MDT" unix="1521757287"/>
  <jk:software web_server="Apache/2.4.7 (Ubuntu) mod_jk/1.2.37" jk_version="mod_jk/1.2.37"/>
  <jk:balancers count="1">
    <jk:balancer name="balancer" type="lb" sticky_session="True" sticky_session_force="False" retries="2" recover_time="60" error_escalation_time="30" max_reply_timeouts="0" method="Busyness" lock="Optimistic" member_count="2" good="2" degraded="0" bad="0" busy="1" max_busy="18" map_count="2" time_to_maintenance_min="47" time_to_maintenance_max="109" last_reset_at="1521756001" last_reset_ago="1286">
      <jk:member name="myserver1" type="ajp13" host="my-server1" port="8009" address="10.1.1.1:8009" connection_pool_timeout="0" ping_timeout="10000" connect_timeout="10000" prepost_timeout="10000" reply_timeout="0" connection_ping_interval="100" retries="2" recovery_options="0" max_packet_size="16384" activation="ACT" lbfactor="10" route="myserver1" redirect="" domain="" distance="0" state="OK" lbmult="1" lbvalue="0" elected="4534" sessions="4534" errors="0" client_errors="0" reply_timeouts="0" transferred="3237735" read="92505649" busy="0" max_busy="9" connected="342" time_to_recover_min="0" time_to_recover_max="0" last_reset_at="1521756001" last_reset_ago="1286"/>
      <jk:member name="myserver2" type="ajp13" host="my-server2" port="8009" address="10.1.1.2:8009" connection_pool_timeout="0" ping_timeout="10000" connect_timeout="10000" prepost_timeout="10000" reply_timeout="0" connection_ping_interval="100" retries="2" recovery_options="0" max_packet_size="16384" activation="ACT" lbfactor="10" route="myserver2" redirect="" domain="" distance="0" state="OK" lbmult="1" lbvalue="1" elected="4651" sessions="4651" errors="0" client_errors="0" reply_timeouts="0" transferred="3482618" read="130672970" busy="1" max_busy="9" connected="336" time_to_recover_min="0" time_to_recover_max="0" last_reset_at="1521756001" last_reset_ago="1286"/>
      <jk:map id="1" server="my-server.domain.com" uri="/*" type="Wildchar" source="JkMount" reply_timeout="-1" sticky_ignore="0" stateless="0" fail_on_status="" active="" disabled="" stopped="" use_server_errors="0"/>
      <jk:map id="2" server="my-server.domain.com [*:80]" uri="/*" type="Wildchar" source="JkMount" reply_timeout="-1" sticky_ignore="0" stateless="0" fail_on_status="" active="" disabled="" stopped="" use_server_errors="0"/>
    </jk:balancer>
  </jk:balancers>
  <jk:result type="OK" message="Action finished"/>
</jk:status>`

func TestHTTPJKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, jkstatus)
	}))
	defer ts.Close()

	tc := JKStatus{
		URL:      ts.URL,
		Username: "jkstatus",
		Password: "s3cret",
	}

	var acc testutil.Accumulator
	err := tc.Gather(&acc)
	require.NoError(t, err)

	// jkstatus_server
	jkServerFields := map[string]interface{}{
		"name": "my-server",
		"port": 8084,
	}
	acc.AssertContainsFields(t, "jkstatus", jkServerFields)

	//	// jkstatus_connector
	//	connectorFields := map[string]interface{}{
	//		"max_threads":          int64(200),
	//		"current_thread_count": int64(5),
	//		"current_threads_busy": int64(1),
	//		"max_time":             int(68),
	//		"processing_time":      int(88),
	//		"request_count":        int(2),
	//		"error_count":          int(1),
	//		"bytes_received":       int64(0),
	//		"bytes_sent":           int64(9286),
	//	}
	//	connectorTags := map[string]string{
	//		"name": "http-apr-8080",
	//	}
	//	acc.AssertContainsTaggedFields(t, "jkstatus_connector", connectorFields, connectorTags)
}
