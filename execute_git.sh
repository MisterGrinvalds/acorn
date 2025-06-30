#!/bin/bash
cd /Users/mistergrinvalds/Repos/personal/bash-profile
export PATH=/usr/bin:/bin:/usr/local/bin:$PATH
git checkout -b fix/automation-framework-compatibility 2>&1
echo "Branch creation result: $?"
git status 2>&1