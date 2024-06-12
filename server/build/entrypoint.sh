#!/bin/bash

chown mattermost:mattermost /mattermost/client/plugins
chown mattermost:mattermost /mattermost/config
chown mattermost:mattermost /mattermost/plugins

if [ "${1:0:1}" = '-' ]; then
    set -- mattermost "$@"
fi

exec "$@"
