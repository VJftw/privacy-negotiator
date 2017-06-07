#!/usr/bin/with-contenv sh

echo "Using API: ${API_ENDPOINT}"
sed -i -e "s#API_ENDPOINT#${API_ENDPOINT}#g" /usr/share/nginx/html/*.js

echo "Using Facebook App: ${FB_APP_ID}"
sed -i -e "s#FB_APP_ID#${FB_APP_ID}#g" /usr/share/nginx/html/*.js
