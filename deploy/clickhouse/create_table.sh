#!/bin/bash
clickhouse-client --host=$HOUSE_HOST --query="CREATE TABLE IF NOT EXISTS shkaff_stat (
  UserId UInt16,
  DbID   UInt16,
  TaskId UInt16,
  Service String,
  NewOperator UInt32,
  SuccessOperator UInt32,
  FailOperator UInt32,
  NewDump UInt32,
  SuccessDump UInt32,
  FailDump UInt32,
  NewRestore UInt32,
  SuccessRestore UInt32,
  FailRestore UInt32,
  Error String,
  CreateDate Date
 ) ENGINE = MergeTree( CreateDate, (UserId, CreateDate), 8192);"
