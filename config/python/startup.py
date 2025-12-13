from os import getenv, path
from sys import version

runtime_version = version.split(" ")[0]
print("Running Startup Scripts for Python {}".format(runtime_version[0]))
if runtime_version[0] == "2":
    exec(open(path.join(getenv("XDG_CONFIG_HOME"), "dotfiles/config/python/startup2.py")).read())
elif runtime_version[0] == "3":
    exec(open(path.join(getenv("XDG_CONFIG_HOME"), "dotfiles/config/python/startup3.py")).read())
else:
    print("\n WARNING: Python version not detected. History writing to $HOME")
