#!/bin/sh
set -e
(
   ROOTDIR=$(dirname $PWD)/..
   cd $ROOTDIR/db/postgres/migrations
   $ROOTDIR/.tools/bin/go-bindata -o ./migrations.go -pkg migrations -nomemcopy ./*.sql
)