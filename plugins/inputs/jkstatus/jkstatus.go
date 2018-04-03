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

type Status struct {
	XMLName   xml.Name  `xml:"http://jkstatus.apache.org status">`
	Balancers Balancers `xml:"balancers"`
}

type Balancers struct {
	Count     int        `xml:"count,attr"`
	Balancers []Balancer `xml:"balancer"`
}

type Balancer struct {
	Name                 string   `xml:"name,attr"`
	Type                 string   `xml:"type,attr"`
	StickySession        bool     `xml:"sticky_session,attr"`
	StickySessionForce   bool     `xml:"sticky_session_force,attr"`
	Retries              int      `xml:"retries,attr"`
	RecoverTime          int      `xml:"recover_time,attr"`
	ErrorEscalationTime  int      `xml:"error_escalation_time,attr"`
	MaxReplyTimeouts     int      `xml:"max_reply_timeouts,attr"`
	Method               string   `xml:"method,attr"`
	Lock                 string   `xml:"lock,attr"`
	MemberCount          int      `xml:"member_count,attr"`
	Good                 int      `xml:"good,attr"`
	Degraded             int      `xml:"degraded,attr"`
	Bad                  int      `xml:"bad,attr"`
	Busy                 int      `xml:"busy,attr"`
	MaxBusy              int      `xml:"max_busy,attr"`
	MapCount             int      `xml:"map_count,attr"`
	TimeToMaintenanceMin int      `xml:"time_to_maintenance_min,attr"`
	TimeToMaintenanceMax int      `xml:"time_to_maintenance_max,attr"`
	LastResetAt          int64    `xml:"last_reset_at,attr"`
	LastResetAgo         int64    `xml:"last_reset_ago,attr"`
	Members              []Member `xml:"member"`
}

type Member struct {
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
	TimeToRecoverMin       int    `xml:"time_to_recover_min,attr"`
	TimeToRecoverMax       int    `xml:"time_to_recover_max,attr"`
	LastResetAt            int64  `xml:"last_reset_at,attr"`
	LastResetAgo           int64  `xml:"last_reset_ago,attr"`
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

	var status Status
	xml.NewDecoder(resp.Body).Decode(&status)

	// add jkstatus_balancer measurements
	for _, b := range status.Balancers.Balancers {
		balancerTags := map[string]string{
			"name": b.Name,
		}

		balancerFields := map[string]interface{}{
			"type":                    b.Type,
			"sticky_session":          b.StickySession,
			"sticky_session_force":    b.StickySessionForce,
			"retries":                 b.Retries,
			"recover_time":            b.RecoverTime,
			"error_escalation_time":   b.ErrorEscalationTime,
			"max_reply_timeouts":      b.MaxReplyTimeouts,
			"method":                  b.Method,
			"lock":                    b.Lock,
			"member_count":            b.MemberCount,
			"good":                    b.Good,
			"degraded":                b.Degraded,
			"bad":                     b.Bad,
			"busy":                    b.Busy,
			"max_busy":                b.MaxBusy,
			"map_count":               b.MapCount,
			"time_to_maintenance_min": b.TimeToMaintenanceMin,
			"time_to_maintenance_max": b.TimeToMaintenanceMax,
			"last_reset_at":           b.LastResetAt,
			"last_reset_ago":          b.LastResetAgo,
		}

		acc.AddFields("jkstatus_balancer", balancerFields, balancerTags)

		for _, m := range b.Members {
			memberTags := map[string]string{
				"balancer": b.Name,
				"name":     m.Name,
			}

			memberFields := map[string]interface{}{
				"type":                     m.Type,
				"host":                     m.Host,
				"port":                     m.Port,
				"address":                  m.Address,
				"connection_pool_timeout":  m.ConnectionPoolTimeout,
				"ping_timeout":             m.PingTimeout,
				"connect_timeout":          m.ConnectTimeout,
				"prepost_timeout":          m.PrepostTimeout,
				"reply_timeout":            m.ReplyTimeout,
				"connection_ping_interval": m.ConnectionPingInterval,
				"retries":                  m.Retries,
				"recovery_options":         m.RecoveryOptions,
				"max_packet_size":          m.MaxPacketSize,
				"activation":               m.Activation,
				"lbfactor":                 m.LbFactor,
				"route":                    m.Route,
				"redirect":                 m.Redirect,
				"domain":                   m.Domain,
				"distance":                 m.Distance,
				"state":                    m.State,
				"lbmult":                   m.LbMult,
				"lbvalue":                  m.LbValue,
				"elected":                  m.Elected,
				"sessions":                 m.Sessions,
				"errors":                   m.Errors,
				"client_errors":            m.ClientErrors,
				"reply_timeouts":           m.ReplyTimeouts,
				"transferred":              m.Transferred,
				"read":                     m.Read,
				"busy":                     m.Busy,
				"max_busy":                 m.MaxBusy,
				"connected":                m.Connected,
				"time_to_recover_min":      m.TimeToRecoverMin,
				"time_to_recover_max":      m.TimeToRecoverMax,
				"last_reset_at":            m.LastResetAt,
				"last_reset_ago":           m.LastResetAgo,
			}

			acc.AddFields("jkstatus_member", memberFields, memberTags)
		}
	}

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
