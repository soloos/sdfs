#!/bin/bash
ps aux|grep sdfsd|grep -v fuse|grep bin|grep -v grep |awk '{print $2}'|xargs kill -SIGABRT
