#!/bin/bash

echo ""
echo "Joining ECS Cluster: ${ecs_cluster_name}!"
echo ""
echo ECS_CLUSTER=${ecs_cluster_name} >> /etc/ecs/ecs.config

echo ""
echo "Upgrading ECS"
echo ""
sudo stop ecs
sudo yum update -y ecs-init
sudo stop ecs
sudo service docker restart

# echo ""
# echo "Upgrading Weave Scope"
# echo ""
# sudo curl -L git.io/scope -o /usr/local/bin/scope
# sudo chmod a+x /usr/local/bin/scope
# sudo stop scope
# sudo start scope

# echo ""
# echo "Disabling Weave Scope"
# echo ""
# rm /etc/init/scope.conf

echo ""
echo "Installing Weave Net"
echo ""
sudo curl -L git.io/weave -o /usr/local/bin/weave
sudo chmod a+x /usr/local/bin/weave

sudo curl -L https://raw.githubusercontent.com/weaveworks/integrations/master/aws/ecs/packer/to-upload/weave.conf -o /etc/init/weave.conf
sudo curl -L https://raw.githubusercontent.com/weaveworks/integrations/master/aws/ecs/packer/to-upload/ecs.override -o /etc/init/ecs.override

sudo mkdir /etc/weave

sudo curl -L https://raw.githubusercontent.com/weaveworks/integrations/master/aws/ecs/packer/to-upload/peers.sh -o /etc/weave/peers.sh
sudo chmod +x /etc/weave/peers.sh
sudo curl -L https://raw.githubusercontent.com/weaveworks/integrations/master/aws/ecs/packer/to-upload/run.sh -o /etc/weave/run.sh
sudo chmod +x /etc/weave/run.sh

sudo start weave
