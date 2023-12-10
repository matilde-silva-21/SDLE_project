#! /bin/bash

trap 'pkill -P $$' TERM

if [[ $# != 1 ]]; then

    echo "Error: Wrong number of arguments. Usage: $0 <BACKEND_PORT>"
    exit 1

else
    
BACKEND_DIR="$PWD/backend"
FRONTEND_DIR="$PWD/frontend"

cd $BACKEND_DIR && go run main.go $1

cd $FRONTEND_DIR && npm i && npm run dev

fi