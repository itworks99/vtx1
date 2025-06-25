#!/bin/sh
java -Xmx512M -cp "./tools/antlr-4.13.1-complete.jar:$CLASSPATH" org.antlr.v4.Tool "$@" 