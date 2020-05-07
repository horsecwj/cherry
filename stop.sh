#!/bin/sh
cat pids/workers.pid  | xargs kill -INT
cat pids/api.pid  | xargs kill -INT
