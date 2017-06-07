#!/bin/bash

sudo stop ecs
sudo stop weave

echo ""
echo "Joining ECS Cluster: ${ecs_cluster_name}!"
echo ""
echo ECS_CLUSTER=${ecs_cluster_name} >> /etc/ecs/ecs.config

echo ""
echo "Upgrading ECS"
echo ""
sudo yum update -y ecs-init

echo ""
echo "Tagging instance"
echo ""
INSTANCE_ID=$(curl http://169.254.169.254/latest/meta-data/instance-id)
export AWS_DEFAULT_REGION=eu-west-1
/usr/local/bin/aws ec2 create-tags --resources "$INSTANCE_ID" --tags Key="weave:peerGroupName",Value="privacy-negotiator"

# echo ""
# echo "Upgrading Weave Scope"
# echo ""
# sudo curl -L git.io/scope -o /usr/local/bin/scope
# sudo chmod a+x /usr/local/bin/scope
# sudo stop scope
# sudo start scope

echo ""
echo "Disabling Weave Scope"
echo ""
rm /etc/init/scope.conf

echo ""
echo "Upgrading Weave Net"
echo ""
sudo curl -L git.io/weave -o /usr/local/bin/weave
sudo chmod a+x /usr/local/bin/weave
sudo stop weave
sudo service docker restart
sudo start weave
sudo start ecs
