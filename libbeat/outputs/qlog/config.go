package qlog

import (
	"github.com/elastic/beats/libbeat/outputs/codec"
)

type Config struct {
	Codec       codec.Config `config:"codec"`
	Container   bool         `config:"container"     validate:"required"`
	IDC         string       `config:"idc"           validate:"required"`
	Level       int          `config:"level"`
	Hosts       []string     `config:"hosts"         validate:"required"`
	Topic       string       `config:"topic"`
	MaxRetries  int          `config:"max_retries"   validate:"min=-1,nonzero"`
	TimeOut     int          `config:"timeOut"`
	MinConnNum  int          `config:"minConnNum"    validate:"min=1"`
	ZkServer    []string     `config:"zkServer"`
	Cluster     string       `config:"cluster"       validate:"required"`
	Compression string       `config:"compression"`
	BulkMaxSize int          `config:"bulk_max_size"`
}

type Qlog struct {
	version byte
	unicode byte
	server  string
	topic   string
	level   uint32
	idc     string
	thread  string
	relay   string
	sec     uint32
	usec    uint32
	file    string
	line    uint32
	msg     []byte
	bin     []byte
}

var config = Config{}

var hostName string
var idc string

func defaultConfig() Config {
	return Config{
		Container:  true,
		MaxRetries: 5,
		TimeOut:    5,
		MinConnNum: 1,
		Level: 20000,
		BulkMaxSize: 2048,
	}
}
