#!/bin/bash
ps aux|grep solofsd|grep -v fuse|grep bin|grep -v grep |awk '{print $2}'|xargs kill -SIGABRT 2>/dev/null
