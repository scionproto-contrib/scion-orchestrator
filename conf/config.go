package conf

import (
	"os"

	scionpila "github.com/netsys-lab/scion-pila"
	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Command       string
	Mode          string
	Bootstrap     Bootstrap `toml:"bootstrap,omitempty"`
	Metrics       Metrics
	IsdAs         string                    `toml:"isd_as,omitempty"`
	Ca            CA                        `toml:"ca,omitempty"`
	ServiceConfig ServiceConfig             `toml:"service_config,omitempty"`
	Api           Api                       `toml:"api,omitempty"`
	Pila          scionpila.SCIONPilaConfig `toml:"pila,omitempty"`
	UI            UI                        `toml:"ui,omitempty"`
}

type UI struct {
	Server string
}

type Api struct {
	Address string   `toml:"address,omitempty"`
	Users   []string `toml:"users,omitempty"`
}

type CA struct {
	CertValidityHours int `toml:"cert_validity_hours,omitempty"`
	Clients           []string
	Server            string
	PKCS11            *PKCS11 `toml:"pkcs11,omitempty"`
}

type PKCS11 struct {
	ModulePath string `toml:"module_path,omitempty"`
	TokenLabel string `toml:"token_label,omitempty"`
	Pin        string `toml:"token_pin,omitempty"`
	KeyLabel   string `toml:"key_label,omitempty"`
}

type Bootstrap struct {
	Server             string
	TopologyOverwrites []string `toml:"topology_overwrites,omitempty"`
	AllowedSubnets     []string `toml:"allowed_subnets,omitempty"`
	AllowClientFail    bool     `toml:"allow_client_fail,omitempty"`
}

type ServiceConfig struct {
	DisableDaemon          bool `toml:"disable_daemon,omitempty"`
	DisableDispatcher      bool `toml:"disable_dispatcher,omitempty"`
	DisableBootstrapServer bool `toml:"disable_bootstrap_server,omitempty"`
	DisablePilaServer      bool `toml:"disable_pila_server,omitempty"`
	DisableCertRenewal     bool `toml:"disable_cert_renewal,omitempty"`
}

type Metrics struct {
	Server     string
	Prometheus string
}

func NewConfig() *Config {
	return &Config{
		Command: "",
		Bootstrap: Bootstrap{
			Server: "",
		},
		Metrics: Metrics{
			Server:     "127.0.0.1:33401",
			Prometheus: "127.0.0.1:33402",
		},
		Ca: CA{
			CertValidityHours: 72,
			Server:            "127.0.0.1:3000",
			Clients:           []string{},
			PKCS11:            nil,
		},
		Api: Api{
			Address: "127.0.0.1:8843",
		},
		ServiceConfig: ServiceConfig{},
		Pila:          scionpila.SCIONPilaConfig{},
		UI: UI{
			Server: "127.0.0.1:8180",
		},
	}
}

func LoadConfig(path string) (*Config, error) {
	c := NewConfig()

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = toml.Unmarshal(content, c)
	if err != nil {
		return nil, err
	}

	// log.Println(c)

	return c, nil
}

func (c *Config) Save(path string) error {
	content, err := toml.Marshal(c)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, content, 0644)
	if err != nil {
		return err
	}

	return nil
}
