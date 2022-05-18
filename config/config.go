package config

type Config struct {
	Mysql      *Mysql      `json:"mysql" mapstructure:"mysql"`
	HttpServer *HttpServer `json:"http_server" mapstructure:"http_server"`
}
type Mysql struct {
	Address  string `json:"address" mapstructure:"address"`
	Username string `json:"username" mapstructure:"username"`
	Password string `json:"password" mapstructure:"password"`
	Database string `json:"database" mapstructure:"database"`
}
type HttpServer struct {
	ListenIp   string `json:"listen_ip" mapstructure:"listen_ip"`
	ListenPort int    `json:"listen_port" mapstructure:"listen_port"`
}
