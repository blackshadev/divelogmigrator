#!/bin/sh
psql -U littledivelog littledivelog < "/docker-entrypoint-initdb.d/divelog.dump"
