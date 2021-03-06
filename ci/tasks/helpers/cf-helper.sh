#!/usr/local/bin/dumb-init /bin/bash
cf_login() {
  local cf_app_url="api.$(pcf cf-info | grep system_domain | cut -d' ' -f 3)"
  local cf_app_user="$(pcf cf-info | grep admin_username | cut -d' ' -f 3)"
  local cf_app_password="$(pcf cf-info | grep admin_password | cut -d' ' -f 3)"
  cf login -a "$cf_app_url" -u "$cf_app_user" -p "$cf_app_password" --skip-ssl-validation -o system -s system
}

#cf.helper. wait for particular service status 
wait_service_status_change() {
  local instance_name=$1
  local status=$2
  local time=0
  echo "checking " $instance_name $status
  local verify_status=$(cf services | awk  '/'"$instance_name"'[ ].*'"$status"'/{print "'"$status"'"}')
  while [[ $verify_status == $status ]] && [[ $time -lt $INSTALL_TIMEOUT ]]; do
    echo "...${verify_status}"
    sleep 3m
    let "time=$time+3"
    verify_status=$(cf services | awk  '/'"$instance_name"'[ ].*'"$status"'/{print "'"$status"'"}')
  done
}

delete_service_app_if_exists() {
  local instance_name=$1
  local app_name=$2
  delete_bind $instance_name $app_name
  delete_application $app_name
  delete_service $instance_name $app_name
}

delete_bind() {
  local instance_name=$1
  local app_name=$2
  echo "check if $app_name binding exist"
  local binding=$(cf services | grep $instance_name | awk '/'"$app_name"'/{print "exist"}')
  if [[ $binding == "exist" ]]; then
      cf unbind-service $app_name $instance_name
      check_app_unbinding $instance_name $app_name
  fi
}

delete_service() {
  local instance_name=$1
  local app_name=$2
  local service=$(cf services | awk '/'"$instance_name"'[ $]/{print "exist"}')
  if [[ $service == "exist" ]]; then
    cf delete-service $instance_name -f
    wait_service_status_change $instance_name "delete in progress"
    service_status=$(cf services | awk  '/'"$instance_name"'[ $].*failed/{print "failed"}')
    if [[ $service_status == "failed" ]]; then
      cf purge-service-instance $instance_name -f
    fi
  fi
}

delete_application() {
  local app_name=$1
  echo "check if $app_name exists"
  local app=$(cf apps | awk '/'"$app_name"'[ $]/{print "exist"}')
  if [[ $app == "exist" ]]; then
    cf delete $app_name -f
  fi
}

create_service() {
  local instance_name=$1
  if [[ -z $MONGO_VERSION ]]; then
    cf create-service mongodb-odb "$SET_PLAN" $instance_name
  else
    cf create-service mongodb-odb "$SET_PLAN" $instance_name -c "{\"version\":\"$MONGO_VERSION\"}"
  fi
  
  wait_service_status_change $instance_name "create in progress"
  service_status=$(cf services | awk  '/'"$instance_name"'[ ].*succeeded/{print "succeeded"}')
  if [[ $service_status != "succeeded" ]]; then
    echo "FAILED! wrong status: $(cf service $instance_name)"
    cf logout
    exit 1
  fi
}

check_app_unbinding() {
  local instance_name=$1
  local app_name=$2
  local app_binding=$(cf services | grep "$instance_name " | awk '!/'"$app_name"'/{print "not bounded"}')
  local try=10
  until [[ $app_binding == "not bounded" ]]; do
    app_binding=$(cf services | grep "$instance_name " | awk '!/'"$app_name"'/{print "not bounded"}')
    if [[ $try -lt 0 ]]; then
      echo "ERROR: unbinding is getting too long"
      exit 1
    fi
    let "try--"
    echo "checking unbinding ($try)"
  done
}

check_app_started() {
  local app_name=$1
  local app=$(cf apps | grep "$app_name " | awk  '/started/{print "started"}')
  local try=10
  until [[ $app == "started" ]]; do
    app=$(cf apps | grep "$app_name " | awk  '/started/{print "started"}')
    if [[ $try -lt 0 ]]; then
      echo "ERROR: unbinding is getting too long"
      exit 1
    fi
    let "try--"
    echo "checking application status ($try)"
  done
  echo "Application started"
}