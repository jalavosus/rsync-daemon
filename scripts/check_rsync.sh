#!/bin/bash

exitcode=$("$(which rsync)" = "" && CODE=1 || CODE=0; echo $CODE)

exit exitcode