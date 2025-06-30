function venv() { 
	if [ -z "$1" ]
	then
			source .venv/bin/activate
	else
		export VENV=$1 && source $ENVS_LOCATION/$VENV/bin/activate ; 
	fi
}
