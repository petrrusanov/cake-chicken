#!/bin/sh
echo "prepare environment"

chown -R app:app /srv /var 2>/dev/null

echo "start backend server"
/sbin/su-exec app /srv/backend $@
