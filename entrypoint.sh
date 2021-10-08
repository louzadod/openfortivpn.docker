#!/bin/bash
set -e
command=$1

/usr/sbin/pcscd --auto-exit

touch /vpn/config.cfg

case $command in
  reconfigure)
    exec vpnconfig --reconfigure /vpn/config.cfg;;
  edit)
    exec nano /vpn/config.cfg;;
  start)
    vpnconfig /vpn/config.cfg
    exec openfortivpn --seclevel-1 --config /vpn/config.cfg;;
  *)
    exit 1;;
esac
