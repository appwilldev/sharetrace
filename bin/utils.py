#!/usr/bin/env python
#-*- coding:utf-8 -*-

import commands


def do_shell_cmd(cmd):
    print 'begin to do: ' + cmd

    status, res = commands.getstatusoutput(cmd)
    if status and status != 0:
        print 'command execute failed: %d' % status
        return status