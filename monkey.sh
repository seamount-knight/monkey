#!/usr/bin/env sh

sync
chmod +x /*.sh /monkey/monkey

sync
supervisord --nodaemon --configuration /etc/supervisord.conf
