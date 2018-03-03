#!/bin/bash
git clone https://github.com/google/leveldb.git
cd leveldb/
make && make install
cd .. && rm -R leveldb
