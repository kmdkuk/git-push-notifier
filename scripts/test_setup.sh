#!/usr/bin/env bash

rm -rf tmp/test
mkdir -p tmp/test

mkdir -p tmp/test/nongit
touch tmp/test/nongit/hogehoge.txt

mkdir -p tmp/test/nocommitpushdir
mkdir -p tmp/test/plaingitdir

git init tmp/test/nocommitpushdir
git init tmp/test/plaingitdir

touch tmp/test/nocommitpushdir/hogehoge.txt
