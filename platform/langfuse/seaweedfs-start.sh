#!/bin/sh
# Starts SeaweedFS (master + volume + filer + S3 gateway in one process) for Langfuse,
# with credentials from S3_ACCESS_KEY / S3_SECRET_KEY, and creates the "langfuse" bucket.
set -eu
cat > /tmp/s3.json <<JSON
{"identities":[{"name":"langfuse","credentials":[{"accessKey":"${S3_ACCESS_KEY}","secretKey":"${S3_SECRET_KEY}"}],"actions":["Admin","Read","Write","List","Tagging"]}]}
JSON
weed server -dir=/data -ip.bind=0.0.0.0 -s3 -s3.port=8333 -s3.config=/tmp/s3.json -master.volumeSizeLimitMB=1024 -volume.max=0 &
pid=$!
# The filer registers a few seconds after start; retry until the bucket is listed.
until echo "s3.bucket.list" | weed shell -master=localhost:9333 2>/dev/null | grep -q "langfuse"; do
  echo "s3.bucket.create -name langfuse" | weed shell -master=localhost:9333 >/dev/null 2>&1 || true
  sleep 2
done
touch /tmp/bucket-ready
echo "seaweedfs: S3 ready on :8333, bucket langfuse"
wait "$pid"
