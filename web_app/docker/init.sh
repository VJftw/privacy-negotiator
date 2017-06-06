#!/usr/bin/with-contenv sh

echo "Using API: ${API_ENDPOINT}"

sed -i -e "s#API_ENDPOINT#${API_ENDPOINT}#g" /usr/share/nginx/html/*.js
