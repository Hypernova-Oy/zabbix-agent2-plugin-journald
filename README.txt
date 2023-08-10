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

