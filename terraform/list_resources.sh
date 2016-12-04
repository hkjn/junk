source aws_dockerized.sh
run_aws ec2 describe-tags --filters "Name=key,Values=orchestration" "Name=value,Values=terraform"
