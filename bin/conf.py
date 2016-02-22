#!/usr/bin/env python
#-*- coding:utf-8 -*-

import os
import ConfigParser

cf = ConfigParser.ConfigParser()
cf.read("./conf/config.ini")


def readTest():
    s = cf.sections()
    print 'section:', s

    o = cf.options("postgres")
    print 'options:', o

    v = cf.items("postgres")
    print 'postgres:', v
    

def get_postgres_config():
    base_dir = cf.get("postgres", "base_dir")
    host = cf.get("postgres", "host")
    port = cf.getint("postgres", "port")
    user = cf.get("postgres", "user")
    passwd = cf.get("postgres", "passwd")
    db_name = cf.get("postgres", "db_name")
    return base_dir, host, port, db_name, user, passwd



def get_redis_config():
    base_dir = cf.get("postgres", "base_dir")
    host = cf.get("redis", "host")
    port = cf.getint("redis", "port")
    return base_dir, host, port

#readTest()
