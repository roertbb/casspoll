#!/bin/bash

cqlsh -f /schema/down.cql
cqlsh -f /schema/up.cql