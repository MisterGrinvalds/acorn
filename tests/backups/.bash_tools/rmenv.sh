## Remove all environment variables matching pattern
rmenv()
{
	unset `env | grep -i "${1:-prox}" | egrep -o '^[^=]+'`
}

