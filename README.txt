# zabbix-agent2-plugin-journald

Fetches messages from Journald.
Very experimental.

Supposed to use the DBUS for fetchin messages, but journald doesnt support DBUS due to architectural reasons.
Separate API and Golang impl exists, but currently no time to work on this more.

This is still better than running direct shell commands from Zabbix Server.

And a nice introduction on using the new Zabbix Agent2 and Golang plugins.

# Deploy

Fetch the release from Github and direct the Zabbix Agent2's plugin config to load it from somewhere:

/etc/zabbix/zabbix_agent2.d/plugins.d/zap_journald.conf:

  Plugins.zap_journald.System.Path=/usr/local/zabbix/go/plugins/zap_journald

Make sure the zabbix-user is in the group systemd-journal, so it can access the journald logs.

# Configure

In Zabbix, configure an Item with the following settings:

  zbx_journald[<ServiceName>,<Since>,<LogLevel>,<Grep>]

  <ServiceName>: (mandatory)
  The name of the systemd service to grep for log entries, eg. networking

  <Since>: (default -60seconds)
  See man journalctl, parameter --since
  How many seconds to look in the past from now, for messages matching the given filters

  <LogLevel>: (default err)
  Which log-level to look for? Your application might not log in a format journald understands, thus you need to use Grep to filter entries.
  See man journalctl, parameter --priority
  The log levels are the usual syslog log levels as documented in syslog(3), i.e.  "emerg" (0), "alert" (1), "crit" (2), "err" (3), "warning" (4), "notice" (5), "info" (6), "debug" (7).

  <Grep>:
  Arbitrary regexp to further filter the log entries. To use this as a generic regexp search tool, just place . as the LogLevel parameter.
  If your application doesn't log in the journald-format, use something like \b(FATAL|ERROR|WARN)\b to detect application-specific log levels.
  See man journalctl, parameter --grep

  Examples:
  zbx_journald[ssauthenticator,-60seconds,7,\b(FATAL|ERROR|WARN)\b]
  zbx_journald[ssauthenticator,-5minutes,err,Auth|auth]
  zbx_journald[ssauthenticator]

# Version

You can check the installed version by calling the plugin directly with the --version -flag.
