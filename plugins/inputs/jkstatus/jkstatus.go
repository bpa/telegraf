package jkstatus

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/internal"
	"github.com/influxdata/telegraf/plugins/inputs"
)

type JKStatusStatus struct {
	Server    JKStatusServer    `xml:"jk:server"`
	Balancers JKStatusBalancers `xml:"jk:balancers"`
}

type JKStatusServer struct {
	Name string `xml:"name,attr"`
	Port int    `xml:"port,attr"`
}

type JKStatusBalancers struct {
	Balancers []JKStatusBalancer
}

type JKStatusBalancer struct {
	Name                 string           `xml:"name,attr"`
	Type                 string           `xml:"type,attr"`
	StickySession        bool             `xml:"sticky_session,attr"`
	StickySessionForce   bool             `xml:"sticky_session_force,attr"`
	Retries              int              `xml:"retries,attr"`
	RecoverTime          int              `xml:"recover_time,attr"`
	ErrorEscalationTime  int              `xml:"error_escalation_time,attr"`
	MaxReplyTimeouts     int              `xml:"max_reply_timeouts,attr"`
	Method               string           `xml:"method,attr"`
	Lock                 string           `xml:"lock,attr"`
	MemberCount          int              `xml:"member_count,attr"`
	Good                 int              `xml:"good,attr"`
	Degraded             int              `xml:"degraded,attr"`
	Bad                  int              `xml:"bad,attr"`
	Busy                 int              `xml:"busy,attr"`
	MaxBusy              int              `xml:"max_busy,attr"`
	MapCount             int              `xml:"map_count,attr"`
	TimeToMaintenanceMin int64            `xml:"time_to_maintenance_min,attr"`
	TimeToMaintenanceMax int64            `xml:"time_to_maintenance_max,attr"`
	LastResetAt          int64            `xml:"last_reset_at,attr"`
	LastResetAgo         int64            `xml:"last_reset_ago,attr"`
	StatusMembers        []JKStatusMember `xml:"jk:member"`
	Maps                 []JKStatusMap    `xml:"jk:map"`
}

type JKStatusMember struct {
	Name                   string `xml:"name,attr"`
	Type                   string `xml:"type,attr"`
	Host                   string `xml:"host,attr"`
	Port                   int    `xml:"port,attr"`
	Address                string `xml:"address,attr"`
	ConnectionPoolTimeout  int    `xml:"connection_pool_timeout,attr"`
	PingTimeout            int    `xml:"ping_timeout,attr"`
	ConnectTimeout         int    `xml:"connect_timeout,attr"`
	PrepostTimeout         int    `xml:"prepost_timeout,attr"`
	ReplyTimeout           int    `xml:"reply_timeout,attr"`
	ConnectionPingInterval int    `xml:"connection_ping_interval,attr"`
	Retries                int    `xml:"retries,attr"`
	RecoveryOptions        int    `xml:"recovery_options,attr"`
	MaxPacketSize          int    `xml:"max_packet_size,attr"`
	Activation             string `xml:"activation,attr"`
	LbFactor               int    `xml:"lbfactor,attr"`
	Route                  string `xml:"route,attr"`
	Redirect               string `xml:"redirect,attr"`
	Domain                 string `xml:"domain,attr"`
	Distance               int    `xml:"distance,attr"`
	State                  string `xml:"state,attr"`
	LbMult                 int    `xml:"lbmult,attr"`
	LbValue                int    `xml:"lbvalue,attr"`
	Elected                int    `xml:"elected,attr"`
	Sessions               int    `xml:"sessions,attr"`
	Errors                 int    `xml:"errors,attr"`
	ClientErrors           int    `xml:"client_errors,attr"`
	ReplyTimeouts          int    `xml:"reply_timeouts,attr"`
	Transferred            int64  `xml:"transferred,attr"`
	Read                   int64  `xml:"read,attr"`
	Busy                   int    `xml:"busy,attr"`
	MaxBusy                int    `xml:"max_busy,attr"`
	Connected              int    `xml:"connected,attr"`
	TimeToRecoverMin       int64  `xml:"time_to_recover_min,attr"`
	TimeToRecoverMax       int64  `xml:"time_to_recover_max,attr"`
	LastResetAt            int64  `xml:"last_reset_at,attr"`
	LastResetAgo           int64  `xml:"last_reset_ago,attr"`
}

type JKStatusMap struct {
	ID              int    `xml:"id,attr"`
	Server          string `xml:"server,attr"`
	Uri             string `xml:"uri,attr"`
	Type            string `xml:"type,attr"`
	Source          string `xml:"source,attr"`
	ReplyTimeout    int    `xml:"reply_timeout,attr"`
	StickyIgnore    int    `xml:"sticky_ignore,attr"`
	Stateless       int    `xml:"stateless,attr"`
	FailOnStatus    int    `xml:"fail_on_status,attr"`
	Active          int    `xml:"active,attr"`
	Disabled        int    `xml:"disabled,attr"`
	Stopped         int    `xml:"stopped,attr"`
	UseServerErrors int    `xml:"use_server_errors,attr"`
}

type JKStatus struct {
	URL      string
	Username string
	Password string
	Timeout  internal.Duration

	SSLCA              string `toml:"ssl_ca"`
	SSLCert            string `toml:"ssl_cert"`
	SSLKey             string `toml:"ssl_key"`
	InsecureSkipVerify bool

	client  *http.Client
	request *http.Request
}

var sampleconfig = `
  ## URL of the JKStatusStatus server status
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
`

func (s *JKStatus) Description() string {
	return "Gather metrics from the JKStatus page."
}

func (s *JKStatus) SampleConfig() string {
	return sampleconfig
}

func (s *JKStatus) Gather(acc telegraf.Accumulator) error {
	if s.client == nil {
		client, err := s.createHttpClient()
		if err != nil {
			return err
		}
		s.client = client
	}

	if s.request == nil {
		_, err := url.Parse(s.URL)
		if err != nil {
			return err
		}
		request, err := http.NewRequest("GET", s.URL, nil)
		if err != nil {
			return err
		}
		request.SetBasicAuth(s.Username, s.Password)
		s.request = request
	}

	resp, err := s.client.Do(s.request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received HTTP status code %d from %q; expected 200",
			resp.StatusCode, s.URL)
	}

	var status JKStatusStatus
	xml.NewDecoder(resp.Body).Decode(&status)
	fmt.Printf("%#v\n", status)

	jkss := map[string]interface{}{
		"name": status.Server.Name,
		"port": status.Server.Port,
	}
	acc.AddFields("jkstatus", jkss, nil)

	//	// add jkstatus_jvm_memorypool measurements
	//	for _, mp := range status.JKStatusStatusJvm.JvmMemoryPools {
	//		tcmpTags := map[string]string{
	//			"name": mp.Name,
	//			"type": mp.Type,
	//		}
	//
	//		tcmpFields := map[string]interface{}{
	//			"init":      mp.UsageInit,
	//			"committed": mp.UsageCommitted,
	//			"max":       mp.UsageMax,
	//			"used":      mp.UsageUsed,
	//		}
	//
	//		acc.AddFields("jkstatus_jvm_memorypool", tcmpFields, tcmpTags)
	//	}
	//
	//	// add jkstatus_connector measurements
	//	for _, c := range status.JKStatusStatusConnectors {
	//		name, err := strconv.Unquote(c.Name)
	//		if err != nil {
	//			name = c.Name
	//		}
	//
	//		tccTags := map[string]string{
	//			"name": name,
	//		}
	//
	//		tccFields := map[string]interface{}{
	//			"max_threads":          c.ThreadInfo.MaxThreads,
	//			"current_thread_count": c.ThreadInfo.CurrentThreadCount,
	//			"current_threads_busy": c.ThreadInfo.CurrentThreadsBusy,
	//			"max_time":             c.RequestInfo.MaxTime,
	//			"processing_time":      c.RequestInfo.ProcessingTime,
	//			"request_count":        c.RequestInfo.RequestCount,
	//			"error_count":          c.RequestInfo.ErrorCount,
	//			"bytes_received":       c.RequestInfo.BytesReceived,
	//			"bytes_sent":           c.RequestInfo.BytesSent,
	//		}
	//
	//		acc.AddFields("jkstatus_connector", tccFields, tccTags)
	//	}

	return nil
}

func (s *JKStatus) createHttpClient() (*http.Client, error) {
	tlsConfig, err := internal.GetTLSConfig(
		s.SSLCert, s.SSLKey, s.SSLCA, s.InsecureSkipVerify)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
		Timeout: s.Timeout.Duration,
	}

	return client, nil
}

func init() {
	inputs.Add("jkstatus", func() telegraf.Input {
		return &JKStatus{
			URL:      "http://127.0.0.1:8080/jkstatus?mime=xml",
			Username: "jkstatus",
			Password: "s3cret",
			Timeout:  internal.Duration{Duration: 5 * time.Second},
		}
	})
}
