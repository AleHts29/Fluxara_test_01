package domain

type Server struct {
	Host string        `mapstructure:"ip"`
	Port string        `mapstructure:"port"`
	TLS  TLSServerConf `mapstructure:"tls"`
}

type TLSServerConf struct {
	Enable              bool   `mapstructure:"enable"`
	HttpsPort           string `mapstructure:"https_port"`
	CertFile            string `mapstructure:"cert_file"`
	KeyFile             string `mapstructure:"key_file"`
	RedirectHTTPToHTTPS bool   `mapstructure:"redirect_http_to_https"`
}
