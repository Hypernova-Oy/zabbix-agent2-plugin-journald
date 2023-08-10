package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"git.zabbix.com/ap/plugin-support/plugin"
	"git.zabbix.com/ap/plugin-support/plugin/container"
)

// PluginOptions option from config file
type PluginOptions struct {
	//The name of the systemctl service whose journal entries to look for
	ServiceName string `conf:""`

	// SinceS is the amount of seconds to look for log messages. This should be the same as the interval of running the attached Item.
	// Defaults to 60s
	SinceS int `conf:"optional,default=60"`

	// Which log-level to look for?
	// Defaults to 'ERROR'
	LogLevel string `conf:"optional,default='ERROR'"`
}

type Plugin struct {
	plugin.Base
	options PluginOptions
}

var impl Plugin

func (p *Plugin) Export(key string, params []string, ctx plugin.ContextProvider) (result interface{}, err error) {
	p.Infof("received request to handle %s key with %d parameters", key, len(params))

	p.Errf("Params '%s'", params)

	var serviceName string
	if len(params) > 0 {
		serviceName = params[0]
	} else {
		return nil, fmt.Errorf("Missing parameter 1 ServiceName")
	}
	p.Errf("serviceName '%s'", serviceName)

	var sinceS string
	if len(params) > 1 {
		_, err = strconv.Atoi(params[1])
		if err != nil {
			return nil, fmt.Errorf("Bad parameter 2 SinceS '%s'", err)
		}
		sinceS = params[1]
	} else {
		sinceS = "60"
	}
	p.Errf("sinceS '%d'", sinceS)

	var logLevel string
	if len(params) > 2 {
		logLevel = params[2]
	} else {
		logLevel = "ERROR"
	}
	p.Errf("logLevel '%s'", logLevel)

	//cmd := exec.Command("ls", "/usr/")
	cmd := exec.Command("journalctl", "-u ", serviceName, "--since", "-"+sinceS+"seconds", "-g", "["+logLevel+"]", "--case-sensitive", "-q")

	var logEntries []byte
	logEntries, err = cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("Executing the command '%s' failed with '%s'", cmd, err)
	}
	p.Errf("logEntries '%s'", logEntries)

	return string(logEntries), nil
}

func init() {
	fmt.Printf("init()%s", os.Args)
	plugin.RegisterMetrics(&impl, "Myip", "myip", "Return the external IP address of the host where agent is running.")
}

func main() {
	fmt.Printf("main()%s", os.Args)
	h, err := container.NewHandler(impl.Name())
	if err != nil {
		panic(fmt.Sprintf("failed to create plugin handler %s", err.Error()))
	}
	impl.Logger = &h

	err = h.Execute()
	if err != nil {
		panic(fmt.Sprintf("failed to execute plugin handler %s", err.Error()))
	}
}
