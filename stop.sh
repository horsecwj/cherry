#!/bin/sh
cat pids/api.pid  | xargs kill -INT
cat pids/workers.pid  | xargs kill -INT
