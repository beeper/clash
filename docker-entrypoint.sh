#!/bin/sh
if [[ ! -z "$CONFIG_FILE" ]]; then
	mkdir -p /root/.config/clash
	cat $CONFIG_FILE | envsubst > /root/.config/clash/config.yaml
fi
exec /clash
