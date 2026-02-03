#!/usr/bin/env sh
set -e

gcloud auth activate-service-account --key-file=/app/scripts/service-account.json

pg_dump -v $URL > ./backup-$(date '+%d-%m-%Y').sql

gsutil cp ./backup-$(date '+%d-%m-%Y').sql gs://$BUCKET_NAME/

echo "Upload complete: gs://$BUCKET_NAME/backup-$(date '+%d-%m-%Y').sql"
