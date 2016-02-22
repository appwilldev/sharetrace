#!/usr/bin/env python

import os

import conf

redis_config = conf.get_redis_config()
base_dir, host, port = redis_config

base_dir = os.path.abspath(base_dir)

name = "default"

data_dir = os.path.join(base_dir, "redis_data/" + name)
pid_path = os.path.join(base_dir, "redis_data/redis_pids/%s.pid" % name)
log_path = os.path.join(base_dir, "redis_data/redis_logs/%s.log" % name)
conf_file = os.path.join(base_dir, "redis_data/redis_confs/%s.conf" % name)

os.system("mkdir -p %s" % data_dir)
os.system("mkdir -p %s" % os.path.join(base_dir, "redis_data/redis_pids"))
os.system("mkdir -p %s" % os.path.join(base_dir, "redis_data/redis_logs"))
os.system("mkdir -p %s" % os.path.join(base_dir, "redis_data/redis_confs"))

if host == 'localhost':
    host = '127.0.0.1'

os.system('sed -e "s|__PORT__|%s|" \
            -e "s|__DATA_DIR__|%s|" \
            -e "s|__PID_PATH__|%s|" \
            -e "s|__LOG_PATH__|%s|" \
            -e "s|__BIND__|%s|" \
            ./conf/redis.conf > %s \
    ' % (port, data_dir, pid_path, log_path, host, conf_file))

print "start %s at: %s port: %s" % (name, data_dir, port)
os.system('redis-server %s' % conf_file)
