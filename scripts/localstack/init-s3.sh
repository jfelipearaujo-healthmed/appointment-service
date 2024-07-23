#!/bin/sh

echo "Initializing Buckets S3..."

awslocal s3 mb s3://appointment-files

echo "Buckets S3 initialized!"