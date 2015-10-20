package service

import (
	"os"
	"testing"

	"github.com/gravitational/teleport/Godeps/_workspace/src/github.com/gravitational/configure"
	"github.com/gravitational/teleport/Godeps/_workspace/src/github.com/gravitational/log"
	. "github.com/gravitational/teleport/Godeps/_workspace/src/gopkg.in/check.v1"
	"github.com/gravitational/teleport/lib/utils"
)

func TestConfig(t *testing.T) { TestingT(t) }

type ConfigSuite struct {
}

var _ = Suite(&ConfigSuite{})

func (s *ConfigSuite) SetUpSuite(c *C) {
	log.Initialize("console", "INFO")
}

func (s *ConfigSuite) TestParseYAML(c *C) {
	var cfg Config
	err := configure.ParseYAML([]byte(configYAML), &cfg)
	c.Assert(err, IsNil)
	s.checkVariables(c, &cfg)
}

func (s *ConfigSuite) TestParseEnv(c *C) {
	vars := map[string]string{
		"TELEPORT_LOG_OUTPUT":                       "console",
		"TELEPORT_LOG_SEVERITY":                     "INFO",
		"TELEPORT_AUTH_SERVERS":                     "tcp://localhost:5000,unix:///var/run/auth.sock",
		"TELEPORT_DATA_DIR":                         "/tmp/data_dir",
		"TELEPORT_FQDN":                             "fqdn.example.com",
		"TELEPORT_AUTH_ENABLED":                     "true",
		"TELEPORT_AUTH_HTTP_ADDR":                   "tcp://localhost:4444",
		"TELEPORT_AUTH_SSH_ADDR":                    "tcp://localhost:5555",
		"TELEPORT_AUTH_DOMAIN":                      "a.fqdn.example.com",
		"TELEPORT_AUTH_TOKEN":                       "authtoken",
		"TELEPORT_AUTH_SECRET_KEY":                  "authsecret",
		"TELEPORT_AUTH_ALLOWED_TOKENS":              "node1.a.fqdn.example.com:token1,node2.a.fqdn.example.com:token2",
		"TELEPORT_AUTH_TRUSTED_USER_AUTHORITIES":    "a.example.com:cert1,b.example.com:cert2",
		"TELEPORT_AUTH_KEYS_BACKEND_TYPE":           "bolt",
		"TELEPORT_AUTH_KEYS_BACKEND_PARAMS":         "path:/keys",
		"TELEPORT_AUTH_KEYS_BACKEND_ADDITIONAL_KEY": "somekey",
		"TELEPORT_AUTH_EVENTS_BACKEND_TYPE":         "bolt",
		"TELEPORT_AUTH_EVENTS_BACKEND_PARAMS":       "path:/events",
		"TELEPORT_AUTH_RECORDS_BACKEND_TYPE":        "bolt",
		"TELEPORT_AUTH_RECORDS_BACKEND_PARAMS":      "path:/records",
		"TELEPORT_SSH_ENABLED":                      "true",
		"TELEPORT_SSH_TOKEN":                        "sshtoken",
		"TELEPORT_SSH_ADDR":                         "tcp://localhost:1234",
		"TELEPORT_SSH_SHELL":                        "/bin/bash",
		"TELEPORT_TUN_ENABLED":                      "true",
		"TELEPORT_TUN_TOKEN":                        "tuntoken",
		"TELEPORT_TUN_SERVER_ADDR":                  "tcp://telescope.example.com",
	}
	for k, v := range vars {
		c.Assert(os.Setenv(k, v), IsNil)
	}
	var cfg Config
	err := configure.ParseEnv(&cfg)
	c.Assert(err, IsNil)
	s.checkVariables(c, &cfg)
}

func (s *ConfigSuite) checkVariables(c *C, cfg *Config) {

	// check logs section
	c.Assert(cfg.Log.Output, Equals, "console")
	c.Assert(cfg.Log.Severity, Equals, "INFO")

	// check common section
	c.Assert(cfg.DataDir, Equals, "/tmp/data_dir")
	c.Assert(cfg.FQDN, Equals, "fqdn.example.com")
	c.Assert(cfg.AuthServers, DeepEquals, NetAddrSlice{
		{Network: "tcp", Addr: "localhost:5000"},
		{Network: "unix", Addr: "/var/run/auth.sock"},
	})

	// auth section
	c.Assert(cfg.Auth.Enabled, Equals, true)
	c.Assert(cfg.Auth.HTTPAddr, Equals,
		utils.NetAddr{Network: "tcp", Addr: "localhost:4444"})
	c.Assert(cfg.Auth.SSHAddr, Equals,
		utils.NetAddr{Network: "tcp", Addr: "localhost:5555"})
	c.Assert(cfg.Auth.Domain, Equals, "a.fqdn.example.com")
	c.Assert(cfg.Auth.Token, Equals, "authtoken")
	c.Assert(cfg.Auth.SecretKey, Equals, "authsecret")

	c.Assert(cfg.Auth.AllowedTokens, DeepEquals,
		KeyVal{
			"node1.a.fqdn.example.com": "token1",
			"node2.a.fqdn.example.com": "token2",
		})

	c.Assert(cfg.Auth.TrustedUserAuthorities, DeepEquals,
		KeyVal{
			"a.example.com": "cert1",
			"b.example.com": "cert2",
		})

	c.Assert(cfg.Auth.KeysBackend.Type, Equals, "bolt")
	c.Assert(cfg.Auth.KeysBackend.Params,
		DeepEquals, KeyVal{"path": "/keys"})
	c.Assert(cfg.Auth.KeysBackend.AdditionalKey, Equals, "somekey")

	c.Assert(cfg.Auth.EventsBackend.Type, Equals, "bolt")
	c.Assert(cfg.Auth.EventsBackend.Params,
		DeepEquals, KeyVal{"path": "/events"})

	c.Assert(cfg.Auth.RecordsBackend.Type, Equals, "bolt")
	c.Assert(cfg.Auth.RecordsBackend.Params,
		DeepEquals, KeyVal{"path": "/records"})

	// SSH section
	c.Assert(cfg.SSH.Enabled, Equals, true)
	c.Assert(cfg.SSH.Addr, Equals,
		utils.NetAddr{Network: "tcp", Addr: "localhost:1234"})
	c.Assert(cfg.SSH.Token, Equals, "sshtoken")
	c.Assert(cfg.SSH.Shell, Equals, "/bin/bash")

	// Tun section
	c.Assert(cfg.Tun.Enabled, Equals, true)
	c.Assert(cfg.Tun.ServerAddr, Equals,
		utils.NetAddr{Network: "tcp", Addr: "telescope.example.com"})
	c.Assert(cfg.Tun.Token, Equals, "tuntoken")
}

const configYAML = `
log:
  output: console
  severity: INFO

data_dir: /tmp/data_dir
fqdn: fqdn.example.com
auth_servers: ['tcp://localhost:5000', 'unix:///var/run/auth.sock']

auth:
  enabled: true
  http_addr: 'tcp://localhost:4444'
  ssh_addr: 'tcp://localhost:5555'
  domain: a.fqdn.example.com
  token: authtoken
  secret_key: authsecret
  allowed_tokens: 
    node1.a.fqdn.example.com: token1
    node2.a.fqdn.example.com: token2

  trusted_user_authorities: 
    a.example.com: cert1
    b.example.com: cert2

  keys_backend:
    type: bolt
    params: {path: "/keys"}
    additional_key: somekey

  events_backend:
    type: bolt
    params: {path: "/events"}

  records_backend:
    type: bolt
    params: {path: "/records"}

ssh:
  enabled: true
  token: sshtoken
  addr: 'tcp://localhost:1234'
  shell: /bin/bash

tun:
  enabled: true
  token: tuntoken
  server_addr: 'tcp://telescope.example.com'
`
