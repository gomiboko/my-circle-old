FROM localstack/localstack:1.2.0

# S3初期データ
RUN mkdir /init-data-s3
COPY init-data-s3/* /init-data-s3

# 初期設定スクリプト
ENV READY_DIR_PATH /etc/localstack/init/ready.d
RUN mkdir ${READY_DIR_PATH}
COPY init.sh ${READY_DIR_PATH}
RUN chmod +x ${READY_DIR_PATH}/init.sh
