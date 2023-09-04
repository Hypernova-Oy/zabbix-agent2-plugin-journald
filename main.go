package main

import (
	"fmt"
	"os"
	"os/exec"

	"git.zabbix.com/ap/plugin-support/plugin"
	"git.zabbix.com/ap/plugin-support/plugin/container"
)

var VERSION = "0.0.1"

// PluginOptions option from config file
type PluginOptions struct {
	//The name of the systemctl service whose journal entries to look for
	ServiceName string `conf:""`

	// since is the amount of seconds to look for log messages. This should be the same as the interval of running the attached Item.
	// See man journalctl, parameter --since
	// Defaults to -60seconds
	Since int `conf:"optional,default=-60seconds"`

	// Which log-level to look for? Your application might not log in a format journald understands, thus you need to use Grep to filter entries.
	// See man journalctl, parameter --priority
	// The log levels are the usual syslog log levels as documented in syslog(3), i.e.  "emerg" (0), "alert" (1), "crit" (2), "err" (3), "warning" (4), "notice" (5), "info" (6), "debug" (7).
	// Defaults to 'err'
	LogLevel string `conf:"optional,default='err'"`

	// Arbitrary regexp to further filter the log entries. To use this as a generic regexp search tool, just place . as the LogLevel parameter.
	// If your application doesn't log in the journald-format, use something like \b(FATAL|ERROR|WARN)\b to detect application-specific log levels.
	// See man journalctl, parameter --grep
	// Defaults to none
	Grep string `conf:"optional"`
}

type Plugin struct {
	plugin.Base
	options PluginOptions
}

var impl Plugin

var plugin_description = "Zabbix Agent Journald plugin v" + VERSION

func (p *Plugin) Export(key string, params []string, ctx plugin.ContextProvider) (result interface{}, err error) {
	p.Infof("received request to handle %s key with %d parameters", key, len(params))

	var serviceName string
	if len(params) > 0 {
		serviceName = params[0]
	} else {
		return nil, fmt.Errorf("Missing parameter 1 ServiceName")
	}

	var since string
	if len(params) > 1 {
		since = params[1]
	} else {
		since = "-60seconds"
	}

	var logLevel string
	if len(params) > 2 {
		logLevel = params[2]
	} else {
		logLevel = "err"
	}

	var grep string
	if len(params) > 3 {
		grep = params[3]
	}

	commandArgs := make([]string, 0, 10)
	commandArgs = append(commandArgs, "-q", "-u", serviceName, "--since", since, "--case-sensitive")
	if logLevel != "" {
		commandArgs = append(commandArgs, "--priority", logLevel)
	}
	if grep != "" {
		commandArgs = append(commandArgs, "--grep", grep)
	}

	cmd := exec.Command("journalctl", commandArgs...)

	var logEntries []byte
	logEntries, err = cmd.CombinedOutput()
	if exiterr, ok := err.(*exec.ExitError); ok {
		// journalctl return 1 when there are no log entries matching the regexp filter. Which is not an error, but actually
		// a good thing.
		if len(logEntries) == 0 && exiterr.ExitCode() == 1 {
			return string(""), nil
		}
	}
	if err != nil {
		return nil, fmt.Errorf("Executing the command '%s' failed with '%+v'. %s", cmd, err, logEntries)
	}

	return string(logEntries), nil
}

func init() {
	plugin.RegisterMetrics(&impl, "ZBX_Journald", "zbx_journald", "Read log messages from journald.")
}

func main() {
	if len(os.Args) >= 2 && (os.Args[1] == "--version" || os.Args[1] == "-v" || os.Args[1] == "--help" || os.Args[1] == "-h") {
		fmt.Printf(plugin_description + "\n")
		os.Exit(0)
	}

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
