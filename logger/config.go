package logger

import (
	"errors"
	"fmt"
	"github.com/xiaoyu830411/properties"
	"os"
	"strconv"
	"strings"
)

var conf config

type config struct {
	logs              map[string]LogConfig
	appenders         map[string]AppenderConfig
	managementService ManagementServiceConfig
}

type LogConfig struct {
	Id           string
	Level        Level
	AppenderList []string
}

type AppenderConfig struct {
	Id     string
	Type   string
	Params map[string]string
}

type ManagementServiceConfig struct {
	Port    int
	Service bool
}

func initConfig() error {
	workDir, err := getCurrentWorkDirectory()
	if err != nil {
		return err
	}

	c, err := loadConfig(workDir + "log4g.properties")
	if err != nil {
		return err
	}

	conf = *c
	return nil
}

func loadConfig(path string) (*config, error) {
	var conf config
	p, err := properties.Load(path)
	if err != nil {
		return nil, err
	}

	conf = config{
		logs:      make(map[string]LogConfig),
		appenders: make(map[string]AppenderConfig),
	}

	if appenders, ok := p.GetSection("appender"); ok {
		elements := make(map[string](map[string]string))
		for id, value := range appenders.Elements() {
			if strings.Index(id, ".") == -1 {
				return nil, errors.New("Appender[" + id + "] has no attributes!")
			}

			keys := strings.Split(id, ".")
			key := keys[0]
			if _, ok = elements[key]; !ok {
				elements[key] = make(map[string]string)
			}

			elements[key][keys[1]] = value
		}
		for id, parameters := range elements {
			t, ok := parameters["type"]
			if !ok {
				return nil, errors.New("No type be specified in appender[" + id + "]")
			} else {
				delete(parameters, "type")
			}

			fmt.Printf("T is %v in appender", t)
			appenderConfig := AppenderConfig{
				Id:     id,
				Params: parameters,
				Type:   t,
			}

			conf.appenders[id] = appenderConfig
		}
	}

	if logs, ok := p.GetSection("logger"); ok {
		for id, value := range logs.Elements() {
			values := strings.Split(value, ",")
			logConfig := LogConfig{}

			l := len(values)

			if l == 0 {
				continue
			}

			if l > 0 {
				level, err := GetLevelByName(strings.TrimSpace(values[0]))
				if err != nil {
					return nil, errors.New(fmt.Sprintf("Invalid logger config[%v]. The logger config : packageName = level, appender1, appender2...", id))
				}

				logConfig.Level = *level
			}

			if l > 1 {
				appenders := values[1:l]
				logConfig.AppenderList = appenders
			}

			conf.logs[id] = logConfig
		}
	}

	if management, ok := p.GetSection("management"); ok {
		managementServiceConfig := ManagementServiceConfig{
			Port:    18080,
			Service: true,
		}

		if port, ok := management.Get("port"); ok {
			i, err := strconv.Atoi(port)
			if err != nil || i < 0 || i > 65535 {
				return nil, errors.New(fmt.Sprintf("Invalid port[%v] in management. It should be a number[1-65535]!", port))
			}

			managementServiceConfig.Port = i
		}

		if service, ok := management.Get("service"); ok {
			b, err := strconv.ParseBool(service)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Invalid service[%v] in management. It should be false or true!", b))
			}

			managementServiceConfig.Service = b
		}

		conf.managementService = managementServiceConfig
	}

	return &conf, nil
}

func GetLogConfig(pkg string) LogConfig {
	for {
		if v, ok := conf.logs[pkg]; ok {
			return v
		} else {
			if pos := strings.LastIndex(pkg, "/"); pos != -1 {
				pkg = pkg[:pos]
			} else {
				break
			}
		}
	}

	return conf.logs["__root__"]
}

func getCurrentWorkDirectory() (string, error) {
	return os.Getwd()
}
