FROM localstack/localstack:1.0.0

# S3初期データ
RUN mkdir /init-data-s3
COPY init-data-s3/* /init-data-s3

# 初期設定スクリプト
COPY init.sh /docker-entrypoint-initaws.d
