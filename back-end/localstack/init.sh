#!/bin/bash

BUCKET_NAME=my-circle-bucket
ACCOUNTS_DIR_NAME=accounts
INIT_S3_DATA_DIR_PATH=/init-data-s3

# S3バケットの作成
awslocal s3api create-bucket --bucket $BUCKET_NAME
# 初期データの投入
awslocal s3api put-object --bucket $BUCKET_NAME --key "${ACCOUNTS_DIR_NAME}/no1.png" --body "${INIT_S3_DATA_DIR_PATH}/no1.png"
awslocal s3api put-object --bucket $BUCKET_NAME --key "${ACCOUNTS_DIR_NAME}/no2.png" --body "${INIT_S3_DATA_DIR_PATH}/no2.png"
