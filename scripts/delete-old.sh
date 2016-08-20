#!/bin/bash

# Used for debugging ElasticSearch integration
# Use with caution!

curl -XGET http://127.0.0.1:19200/_cat/indices
curl -XDELETE http://127.0.0.1:19200/nginx-*
curl -XDELETE http://127.0.0.1:19200/php-*

# curl -XDELETE http://127.0.0.1:19200/.kibana