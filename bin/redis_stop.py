#!/usr/bin/env python

import os

import conf


redis_config = conf.get_redis_config()
base_dir, host, port = redis_config

print "stop redis host:%s, port: %s" % (host,  port)
os.system('redis-cli -h %s -p %d shutdown' % (host, port))
