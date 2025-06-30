# Creates a *.zip archive of a file or folder
mkzip() 
{ 
    zip -r "${1%%/}.zip" "${1%%/}/"
} 

