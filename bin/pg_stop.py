#!/usr/bin/env python
#-*- coding:utf-8 -*-

import os
import conf


def main():
    db_config = conf.get_postgres_config()
    base_dir, host, port, db_name, user, passwd = db_config
    base_dir = os.path.abspath(base_dir)
    data_dir = os.path.join(base_dir, "pg_data/")
    os.system('pg_ctl stop -D %s -p %s -m fast' % (data_dir, port))


if __name__ == "__main__":
    main()
