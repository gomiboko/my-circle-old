#!/bin/bash

readonly BUCKET_NAME="my-circle-bucket"
readonly ACCOUNTS_DIR_NAME="accounts"
readonly CIRCLES_DIR_NAME="circles"
readonly INIT_S3_DATA_DIR_PATH="/init-data-s3"

# S3バケットの作成
awslocal s3api create-bucket --bucket $BUCKET_NAME
# 初期データの投入
awslocal s3api put-object --bucket $BUCKET_NAME --key "${ACCOUNTS_DIR_NAME}/no1.png" --body "${INIT_S3_DATA_DIR_PATH}/no1.png"
awslocal s3api put-object --bucket $BUCKET_NAME --key "${ACCOUNTS_DIR_NAME}/no2.png" --body "${INIT_S3_DATA_DIR_PATH}/no2.png"
# サークルアイコンの初期データ投入
awslocal s3api put-object --bucket $BUCKET_NAME --key "${CIRCLES_DIR_NAME}/54a2e0a21c246a49c4b2f3057ea78da4a38952dbbfa450bc120bde5d99f0a7eb" --body "${INIT_S3_DATA_DIR_PATH}/circle1.png"
awslocal s3api put-object --bucket $BUCKET_NAME --key "${CIRCLES_DIR_NAME}/793b589f61a6426a9e6f1891f9ad9db4dfa10b3cea192fe4fa736100e0c02976" --body "${INIT_S3_DATA_DIR_PATH}/circle2.png"
