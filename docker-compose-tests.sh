#!/bin/bash
set -x
set -e

PROJ=$(basename $PWD)

docker-compose down --remove-orphans || echo new or partial install
docker rmi "${PROJ}-pgsectest" || echo image was not there
if [ ! -e config/tls/int_server/certs/localhost.pem ]; then
  rm -rf config/tls
  docker-compose up chainsmith --exit-code-from chainsmith
fi
docker-compose up -d postgres

for ((i=0;i<60;i++)); do
  docker-compose exec -u postgres postgres pg_isready && break
  echo $i
  sleep 1
done
docker-compose logs postgres

#export PGUSER=postgres
#export PGHOST=127.0.0.1
#export PGSSLCERT=config/tls/int_client/certs/postgres.pem
#export PGSSLKEY=config/tls/int_client/private/postgres.key.pem
#export PGSSLROOTCERT=config/tls/int_server/certs/ca-chain-bundle.cert.pem
docker-compose up pgsectest --exit-code-from pgsectest

echo "All is as expected"
