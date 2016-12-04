source aws_dockerized.sh

delete_instances() {
  local instances=$(run_aws ec2 describe-tags --filters \
                            "Name=key,Values=orchestration" \
                            "Name=value,Values=terraform" \
                            "Name=resource-type,Values=instance" \
                            --query 'Tags[*].ResourceId' \
                            --output text)
  for instance in $instances; do
    local i=$(echo $instance | tr -d '[:space:]')
    echo "[delete_resources.sh] Terminating EC2 instance '$i'..";
    run_aws ec2 terminate-instances --instance-ids $i
  done
}

delete_sgs() {
  local sgs=$(run_aws ec2 describe-tags --filters \
                      "Name=key,Values=orchestration" \
                      "Name=value,Values=terraform" \
                      "Name=resource-type,Values=security-group" \
                      --query 'Tags[*].ResourceId' \
                      --output text)
  for sg in $sgs; do
    local i=$(echo $sg | tr -d '[:space:]')
    echo "[delete_resources.sh] Deleting SG '$i'..";
    run_aws ec2 delete-security-group --group-id $i
  done
}

delete_subnets() {
  local subnets=$(run_aws ec2 describe-tags --filters \
                          "Name=key,Values=orchestration" \
                          "Name=value,Values=terraform" \
                          "Name=resource-type,Values=subnet" \
                      --query 'Tags[*].ResourceId' \
                      --output text)
  for subnet in $subnets; do
    local i=$(echo $subnet | tr -d '[:space:]')
    echo "[delete_resources.sh] Deleting subnet '$i'..";
    run_aws ec2 delete-subnet --subnet-id $i
  done
}

delete_route_tables() {
  local rts=$(run_aws ec2 describe-tags --filters \
                           "Name=key,Values=orchestration" \
                           "Name=value,Values=terraform" \
                           "Name=resource-type,Values=route-table" \
                           --query 'Tags[*].ResourceId' \
                          --output text)
  for rt in $rts; do
    local i=$(echo $rt | tr -d '[:space:]')
    echo "[delete_resources.sh] Deleting route table '$i'..";
    run_aws ec2 delete-route-table --route-table-id $i
  done
}


delete_vpc() {
  local gateways=$(run_aws ec2 describe-tags --filters \
                          "Name=key,Values=orchestration" \
                          "Name=value,Values=terraform" \
                          "Name=resource-type,Values=internet-gateway" \
                          --query 'Tags[*].ResourceId' \
                          --output text)
  local vpcs=$(run_aws ec2 describe-tags --filters \
                       "Name=key,Values=orchestration" \
                       "Name=value,Values=terraform" \
                       "Name=resource-type,Values=vpc" \
                       --query 'Tags[*].ResourceId' \
                       --output text)
  for gateway in $gateways; do
    local gw=$(echo $gateway | tr -d '[:space:]')
    vpc=$(run_aws ec2 describe-internet-gateways \
                  --internet-gateway-id $gw \
                  --query 'InternetGateways[*].Attachments[*].VpcId' \
                  --output text | tr -d '[:space:]')
    echo "[delete_resources.sh] Detaching IGW '$i' from VPC '$vpc'..";
    run_aws ec2 detach-internet-gateway --internet-gateway-id $gw --vpc-id $vpc
    echo "[delete_resources.sh] Deleting IGW '$i'..";
    run_aws ec2 delete-internet-gateway --internet-gateway-id $i
  done
  for vpc in $vpcs; do
    local i=$(echo $vpc | tr -d '[:space:]')
    echo "[delete_resources.sh] Deleting VPC '$i'..";
    run_aws ec2 delete-vpc --vpc-id $i
  done
}

delete_instances
delete_sgs
delete_subnets
delete_vpc
# TODO: This should delete the elastic-ip too, but since there's no
# way to tag those it would be difficult to script..
echo "Done deleting AWS resources tagged orchestration=terraform."


