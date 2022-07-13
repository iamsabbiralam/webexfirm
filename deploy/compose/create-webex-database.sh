export PGUSER=postgres
for db in hrmdb; do
	psql <<-EOSQL
		    	 CREATE USER ${db} with password 'secret';
		    	 CREATE DATABASE ${db};
		    	 GRANT ALL PRIVILEGES ON DATABASE $db TO ${db};
	EOSQL

	export PGDATABASE=${db}

	for e in ${CREATE_EXTENSION}; do
		psql <<-EOSQL
			     CREATE EXTENSION IF NOT EXISTS "${e}";
		EOSQL
	done
done
