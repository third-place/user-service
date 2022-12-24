set -e

# create dev user service database
psql -h localhost -p 54321 -U postgres -c 'create database user_service'

# create kafka topics
kafka-topics --bootstrap-server broker:9092 \
             --create \
             --topic images

kafka-topics --bootstrap-server broker:9092 \
             --create \
             --topic users
