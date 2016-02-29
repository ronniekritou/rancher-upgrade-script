# Rancher Auto Upgrade Script

The needed ENV variables are : 

- RANCHER_LAUNCH_CONFIG_JSON
- RANCHER_AUTH_SID
- RANCHER_AUTH_KEY
- RANCHER_BASE_URL


### ex : 

export RANCHER_LAUNCH_CONFIG_JSON='{"inServiceStrategy": { "launchConfig": { "imageUuid":"docker:MY_DOCKER_IMAGE","ports":["80:80/tcp"] }},"toServiceStrategy":null}' RANCHER_AUTH_SID=MY_RANCHER_AUTH_SID RANCHER_AUTH_KEY=MY_RANCHER_AUTH_KEY RANCHER_BASE_URL=MY_MANCHER_HOST_URL/v1/projects/MY_RANCHER_PROJECT_ID/services/MY_RANCHER_SERVICE_ID && rancher-upgrade-script 