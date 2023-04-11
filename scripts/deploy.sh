#!/bin/bash

echo 'Starting to Deploy...'
ssh -o StrictHostKeyChecking=no "${REMOTE_USER}"@"${REMOTE_HOST}" -i private_key "
        cd ${TARGET}
        git checkout main
        git fetch --all
        git reset --hard origin/main
        git pull origin main &&
        cd ./docker &&
        sudo docker-compose -f docker-compose.${ENVIRONMENT}.yaml down -v
        sudo docker-compose -f docker-compose.${ENVIRONMENT}.yaml up -d
        "
echo 'Deployment completed successfully'