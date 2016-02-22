#!/usr/bin/env python
#-*- coding:utf-8 -*-

import os
import sys
import conf


def main():
    db_config = conf.get_postgres_config()
    print db_config

    base_dir, host, port, db_name, user, passwd = db_config

    base_dir = os.path.abspath(base_dir)
    
    data_dir = os.path.join(base_dir, "pg_data/")
    if os.path.exists(data_dir):
        print "%s is not empty!!!" % data_dir
        sys.exit(1)
    
    os.system('initdb -D %s' % data_dir)
    
    os.system('echo "shared_buffers = 512MB" >> %s/postgresql.conf' % data_dir)
    os.system('echo "temp_buffers = 32MB" >> %s/postgresql.conf' % data_dir)
    os.system('echo "work_mem = 32MB" >> %s/postgresql.conf' % data_dir)
    os.system('echo "maintenance_work_mem = 28MB" >> %s/postgresql.conf' % data_dir)
    os.system('echo "checkpoint_segments = 16" >> %s/postgresql.conf' % data_dir)
    
    os.system('pg_ctl -D %s -l %s/postgresql.log start -o "-p %d"' % (data_dir, data_dir, port))
    
    os.system('echo "wait 2 seconds..."')
    os.system('sleep 2s')
    
    os.system('createuser -s %s -h localhost -p %d' % (user, port))
    os.system('echo "create database %s encoding=\'utf-8\' template=template0;" | psql -U postgres -h localhost -p %d'
              % (db_name, port))

    os.system('psql %s -f ./conf/tables.sql -p %d -h localhost' % (db_name, port))


if __name__ == "__main__":
    main()
