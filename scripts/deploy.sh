#!/bin/bash

echo 'Starting to Deploy...'
ssh -o StrictHostKeyChecking=no "${REMOTE_USER}"@"${REMOTE_HOST}" -i private_key "
        cd ${TARGET}
        git checkout ${BRANCH}
        git fetch --all
        git reset --hard origin/${BRANCH}
        git pull origin ${BRANCH} &&
        cd ./docker &&
        docker-compose -p ${ENVIRONMENT} -f docker-compose.${ENVIRONMENT}.yaml down --volumes
        docker-compose -p ${ENVIRONMENT} -f docker-compose.${ENVIRONMENT}.yaml pull
        docker-compose -p ${ENVIRONMENT} -f docker-compose.${ENVIRONMENT}.yaml up -d
        "
echo 'Deployment completed successfully'