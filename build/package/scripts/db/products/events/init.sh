#!/bin/bash

for ((i = 1 ; i < 10 ; i++)); do
    echo "Pinging db ${i}..."
	curl http://db-products-events:8000 > /dev/null 2>&1
    if [ 0 -eq $? ]
    then
        echo "DB ping successful!"
        break
    fi
    sleep 2
done


aws dynamodb create-table \
    --table-name Events \
    --attribute-definitions \
        AttributeName=ID,AttributeType=S \
    --key-schema \
        AttributeName=ID,KeyType=HASH \
    --provisioned-throughput \
        ReadCapacityUnits=10,WriteCapacityUnits=5 \
    --endpoint-url http://db-products-events:8000/
