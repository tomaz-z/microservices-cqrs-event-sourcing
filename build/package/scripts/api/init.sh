#!/bin/bash

go build -mod vendor -o ./bin/api ./services/*/cmd

./bin/api
