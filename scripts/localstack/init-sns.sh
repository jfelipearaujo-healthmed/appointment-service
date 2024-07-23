#!/bin/sh

echo "Initializing SNS topics..."

awslocal sns create-topic \
    --name AppointmentTopic

awslocal sns create-topic \
    --name FeedbackTopic

echo "SNS topics initialized!"