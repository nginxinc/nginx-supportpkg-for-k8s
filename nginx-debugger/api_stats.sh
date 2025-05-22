#!/usr/bin/env bash
for i in /api/8/processes /api/8/connections /api/8/slabs /api/8/http/requests /api/8/http/server_zones /api/8/http/location_zones /api/8/http/caches /api/8/http/upstreams /api/8/http/keyvals; do
  echo "**** $i ****" ; 
  curl -s "127.0.0.1:8080/$i" | jq .; 
  echo ""; 
done