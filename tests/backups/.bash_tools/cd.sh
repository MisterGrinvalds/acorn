## Always `ll` after a `cd`
cd()
{ 
	builtin cd "$@" && ll 
} 

