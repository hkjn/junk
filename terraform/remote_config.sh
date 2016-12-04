source tf_dockerized.sh
run_tf remote config \
       -backend=s3 \
       -backend-config="bucket=terraform-test-state" \
       -backend-config="key=junk-k8s/terraform.tfstate" \
       -backend-config="region=eu-west-1"
