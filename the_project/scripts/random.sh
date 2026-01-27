#!/bin/bash

echo "DATABASE_URL: $DATABASE_URL"

url=$(wget --spider -S https://en.wikipedia.org/wiki/Special:Random 2>&1 \
  | grep -i "Location:" \
  | awk '{print $2}' \
  | tr -d '\r')

/app/project random "$url"