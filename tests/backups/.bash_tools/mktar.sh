# Creates a *.tar.gz archive of a file or folder
mktar() 
{ 
	tar cvzf "${1%%/}.tar.gz"  "${1%%/}/"
} 

