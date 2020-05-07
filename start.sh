#!/bin/sh

nohup ./cmd/workers >> logs/workers.log &
nohup ./cmd/api >> logs/api.log &
