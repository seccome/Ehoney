#!/bin/bash
#
DB_Port=3306
DB_User="root"
DB_Database="security#123456#"
DB_Password="12345"
REDIS_Port="6379"
REDIS_Password="123456"
Project_Dir=/home
Project_Port=8082
Project_Front_Port=8080
Local_Host="127.0.0.1"
File_Trace_port=5000

function install_soft() {
  echo -e "[ start install soft [ $1 ]"
  if command -v yum >/dev/null; then
    yum -q -y install $1
  elif command -v apt >/dev/null; then
    apt-get -qqy install $1
  elif command -v zypper >/dev/null; then
    zypper -q -n install $1
  elif command -v apk >/dev/null; then
    apk add -q $1
  else
    echo -e "[\033[31m ERROR \033[0m] Please install it first $1 "
    exit 1
  fi
}

function exit_if_process_error() {
  PROC_NAME=$1
  ProcNumber=$(ps -ef | grep -w $PROC_NAME | grep -v grep | wc -l)
}

function check_docker_service() {
  echo docker version
  if [ $? -ne 0 ]; then
    echo "docker service error, exit!!"
    exit
  else
    echo "docker service good"
  fi
}

function check_falco_exist() {
  exist=$(helm list | grep falco)
  if [ "${exist}" != "" ]; then
    echo "falco pod exit, skip install"

  else
    echo "${containerName} docker service start good"
  fi
}

function check_k3s_service() {
  echo k3s --version
  if [ $? -ne 0 ]; then
    echo "k3s service error, exit!!"
    exit
  else
    echo "k3s service good"
  fi
}

function check_docker_container_state() {
  sleep 2s
  # 查看进程是否存在
  containerName=$1
  exist=$(docker inspect --format '{{.State.Running}}' ${containerName})
  if [ "${exist}" != "true" ]; then
    echo "${containerName} docker service start error, exit!!"
    exit
  else
    echo "${containerName} docker service start good"
  fi
}

function check_file_trace_service() {
  sleep 7s
  curl http://${Local_Host}:5000/health | grep 'SUCCEED' &>/dev/null
  if [ $? -ne 0 ]; then
    echo "file_trace service error, exit!!"
    exit
  else
    echo "file_trace service good"
  fi
}

function check_decept_defense_service() {
  sleep 5s
  curl http://localhost:8082/api/public/health | grep 'code' &>/dev/null
  if [ $? -ne 0 ]; then
    echo "Cheat defense back end service health detection failed, please detect manually!"
  else
    echo "decept_defense service good"
  fi
}

function kill_if_process_exist() {
  PROC_NAME=$1
  echo "--------Start killing $PROC_NAME process and its child processes---------"
  ProcNumber=$(ps -ef | grep $PROC_NAME | grep -v "grep" | awk '{print $2}')
  if [ $ProcNumber ]; then
    echo "进程ID: $ProcNumber"
    ps --ppid $ProcNumber | awk '{if($1~/[0-9]+/) print $1}' | xargs kill -9
    kill -9 $ProcNumber
    echo "--------------------End of killing process---------------------------"
  fi
}

# 检查软件是否安装 curl wget zip go redis mysql docker kubectl;
function prepare_base_install() {
  for i in wget curl unzip kernel-devel-$(uname -r); do
    command -v $i &>/dev/null || install_soft $i
  done
  # yumRepoUpdate
}

function resetAgentJson() {
  dos2unix ${RelayDir}/agent/conf/agent.json
  echo "{
  \"honeyPublicIp\": \"\",
  \"strategyAddr\": \"${Local_Host}:${REDIS_Port}\",
  \"strategyPass\": \"${REDIS_Password}\",
  \"version\": \"1.0\",
  \"sshKeyUploadUrl\": \"http://${Local_Host}:${Project_Port}/api/public/protocol/key\"
  }" >${RelayDir}/agent/conf/agent.json

}

function is_port_bind() {
  processor=$(netstat -lnpt | grep $1 | awk '{print $2}')
  if [ "$processor" != "" ]; then
    echo "true"
  else
    echo "false"
  fi
}

function stop_docker_container_if_exist() {
  docker_container_name=$1
  processor=$(docker ps | grep docker_container_name)
  echo "stop docker container if exist [$docker_container_name]..."
  if [ "$processor" != "" ]; then
    docker stop $docker_container_name
    sleep 2s
    docker rm -f $docker_container_name
    sleep 2s
  fi
}

function ports_check() {
  command -v lsof &>/dev/null || install_soft lsof
  echo "Start to detect whether the necessary ports required by the service are occupied!!"
  shouldreturn="false"
  httpfrontport=$(is_port_bind $Project_Front_Port)
  httpwebport=$(is_port_bind $Project_Port)
  redisport=$(is_port_bind $REDIS_Port)
  mysqlport=$(is_port_bind $DB_Port)
  filetraceport=$(is_port_bind $File_Trace_port)
  if [ "$httpfrontport" = "true" ]; then
    echo "$Project_Front_Port Occupied!!"
    shouldreturn="true"
  fi
  if [ "$httpwebport" = "true" ]; then
    echo "$Project_Port Occupied!!"
    shouldreturn="true"
  fi
  if [ "$redisport" = "true" ]; then
    echo "$REDIS_Port Occupied!!"
    shouldreturn="true"
  fi
  if [ "$mysqlport" = "true" ]; then
    echo "$DB_Port Occupied!!"
    shouldreturn="true"
  fi
  if [ "$filetraceport" = "true" ]; then
    echo "$File_Trace_port Occupied!!"
    shouldreturn="true"
  fi

  if [ "$shouldreturn" = "true" ]; then
    echo "Make sure the above ports are available!!"
    echo "The program started to exit!!"
    exit
  fi
}

function setup_iptables() {
  echo "--------------------开始安装Iptables-----------------------------"
  yum -y install iptables-services
  echo "*filter
:INPUT ACCEPT [0:0]
:FORWARD ACCEPT [0:0]
:OUTPUT ACCEPT [0:0]
COMMIT" >/etc/sysconfig/iptables
  systemctl start iptables.service
  systemctl restart iptables.service
  # iptables -t nat -A POSTROUTING -s 172.17.0.0/24 -j SNAT --to-source $Local_Host
  # iptables -A OUTPUT -p all -m state --state RELATED,ESTABLISHED -j ACCEPT
  # iptables -A OUTPUT -p udp --dport 53 -m state --state NEW,RELATED,ESTABLISHED -j ACCEPT
  # iptables -A OUTPUT -p all -d $Local_Host -m state --state NEW,RELATED,ESTABLISHED -j ACCEPT
  # iptables-save
  echo "--------------------Iptables安装完毕-----------------------------"
}

function prepare_check() {
  Project_Dir=$(
    cd $(dirname $0)
    pwd
  )
  echo $(
    cd $(dirname $0)
    pwd
  )
  isRoot=$(id -u -n | grep root | wc -l)
  if [ "x$isRoot" != "x1" ]; then
    echo -e "[\033[31m ERROR \033[0m] Please use root to execute the installation script (请用 root 用户执行安装脚本)"
    exit 1
  fi
  processor=$(cat /proc/cpuinfo | grep "processor" | wc -l)
  if [ $processor -lt 2 ]; then
    echo -e "[\033[31m ERROR \033[0m] The CPU is less than 2 cores (CPU 小于 2核，机器的 CPU 需要至少 2核)"
    exit 1
  fi
  memTotal=$(cat /proc/meminfo | grep MemTotal | awk '{print $2}')
  if [ $memTotal -lt 3750000 ]; then
    echo -e "[\033[31m ERROR \033[0m] Memory less than 4G (内存小于 4G，机器的内存需要至少 4G)"
    exit 1
  fi
}

function prepare_conf() {
  dos2unix $Project_Dir/configs/configs.toml
  DB_Port=$(readConfValue $Project_Dir/configs/configs.toml dbport)
  DB_User=$(readConfValue $Project_Dir/configs/configs.toml dbuser)
  DB_Database=$(readConfValue $Project_Dir/configs/configs.toml dbname)
  DB_Password=$(readConfValue $Project_Dir/configs/configs.toml dbpassword)
  REDIS_Port=$(readConfValue $Project_Dir/configs/configs.toml redisport)
  REDIS_Password=$(readConfValue $Project_Dir/configs/configs.toml redispassword)
  echo -e "DB_Port: "$DB_Port
  echo -e "DB_User: "$DB_User
  echo -e "DB_Database: "$DB_Database
  echo -e "DB_Password: "$DB_Password
  echo -e "REDIS_Port: "$REDIS_Port
  echo -e "REDIS_Password: "$REDIS_Password
}

function component_installer() {
  setupDocker # 安装Docker
  setupK3s    # 安装K3S
  setupFalco
  setupRedis       # 安装Redis
  setupMysqlDocker # 安装MySQL 5.6 以及sec_pot database 的镜像
  setupFileTrace
  setupDeceptDefence
  setupRelayAgent
}
function setupFalco() {
  echo "--------------------Start deploying Falco----------------------------"
  Project_Dir=$(
    cd $(dirname $0)
    pwd
  )
  echo "start setup falco"
  helmFile=/usr/bin/helm
  if [ ! -f "${helmFile}" ]; then
    echo "helm file not found, install falco!"
    falco_install
  else
    exist=$(helm list | grep falco)
    if [ "${exist}" != "" ]; then
      echo "falco pod exit, skip install"
    else
      echo "falco pod not found, install falco!"
      falco_install
    fi
  fi
  echo "--------------------End of Falco installation-----------------------------"
}

function falco_install() {
  cd $Project_Dir/tool/helm
  cp helm /usr/bin/helm
  chmod +x /usr/bin/helm
  findResult=$(echo $PATH | grep '/usr/bin')
  echo "find usr/bin result: ${findResult}"
  if [ "${findResult}" != "" ]; then
    echo \  >>/etc/profile
    echo PATH=$PATH:/usr/bin:/usr/local/bin >>/etc/profile
    source /etc/profile
  fi
  helm repo add stable http://mirror.azure.cn/kubernetes/charts/
  helm repo add falcosecurity https://falcosecurity.github.io/charts
  helm repo update
  helm delete falco
  swpConfValue $Project_Dir/tool/falco/values.yaml 259 false true
  swpConfValue $Project_Dir/tool/falco/values.yaml 260 127.0.0.1 $Local_Host
  cd $Project_Dir/tool/falco
  helm install falco . -n default
}

function readConfValue() {
  configfile=$1
  key=$2
  ww='"'
  while read line; do
    k=${line%=*}
    v=${line#*=}
    v=${v// /}
    k=${k// /}
    if [ "$key" == "$k" ]; then
      v=${v//"${ww}"/}
      echo ${v}
    fi
  done <$configfile
}

function swpConfValue() {
  configfile=$1
  old=$3
  value=$4
  echo "$2s/$old/$value/" $configfile
  sed -i "$2s/$old/$value/" $configfile
}

function setupK3s() {
  echo "--------------------Start deploying K3S-----------------------------"
  # curl -sfL https://get.k3s.io | sh -
  # K3S 配置修改
  echo "start setup k3s service"
  dos2unix ${Project_Dir}/tool/k3s/k3s.sh
  sh ${Project_Dir}/tool/k3s/k3s.sh
  sed -i '28c ExecStart=/usr/local/bin/k3s server --docker --no-deploy traefik' /etc/systemd/system/k3s.service
  sed -i '29c #' /etc/systemd/system/k3s.service
  sudo systemctl daemon-reload
  sudo systemctl restart k3s
  swpConfValue /etc/rancher/k3s/k3s.yaml 5 127.0.0.1 $Local_Host
  echo \ export KUBECONFIG=/etc/rancher/k3s/k3s.yaml >>/etc/profile
  # 覆盖项目中的k3s 配置
  yes | cp -rf /etc/rancher/k3s/k3s.yaml ${Project_Dir}/configs/.kube/config
  source /etc/profile
  sleep 1s
  echo "--------------------End of K3S installation-----------------------------"
  check_k3s_service
}

function setupRedis() {
  echo "--------------------Start deploying Redis-----------------------------"
  rm -rf /etc/decept-defense/data
  mkdir -p /etc/decept-defense/data
  docker pull redis:5.0.6
  docker run -p ${REDIS_Port}:6379 -v /etc/decept-defense/data:/data --name decept-redis -d redis:5.0.6 redis-server --requirepass ""${REDIS_Password}""
  echo "docker run -p ${REDIS_Port}:6379 -v /etc/decept-defense/data:/data --name decept-redis -d redis:5.0.6 redis-server --requirepass ${REDIS_Password}"
  check_docker_container_state decept-redis
  echo "--------------------End of redis installation------------------------------"
}

function setupDocker() {
  echo docker version
  if [ $? -ne 0 ]; then
    echo "docker is uninstalled, start install!"
    install_Docker
  else
    if [ $(grep -c "47.96.71.197:90" /usr/lib/systemd/system/docker.service) -eq '1' ]; then
      echo "docker.service is configured, skip!"
    else
      echo "docker is unconfigured, restart install!"
      sudo systemctl stop docker
      install_Docker
    fi
  fi
  sudo systemctl restart docker
  sleep 1s
  check_docker_service
}

function install_Docker() {
  echo "--------------------Start deploying Docker-----------------------------"
  #  安装依赖包
  sudo yum install -y yum-utils \
    device-mapper-persistent-data \
    lvm2
  # centos8 需要
  # yum install -y https://download.docker.com/linux/fedora/30/x86_64/stable/Packages/containerd.io-1.2.6-3.3.fc30.x86_64.rpm
  # 添加源，使用了阿里云镜像
  sudo yum-config-manager \
    --add-repo \
    http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo
  # 配置缓存
  sudo yum makecache fast
  # 安装最新稳定版本的docker
  sudo yum install -y docker-ce

  # 配置docker文件
  sed -i "13c ExecStart=/usr/bin/dockerd --insecure-registry=47.96.71.197:90" /usr/lib/systemd/system/docker.service

  #  sudo tee /etc/docker/daemon.json <<-'EOF'
  #  {
  #   "registry-mirrors": ["https://docker.mirrors.ustc.edu.cn"],
  #   "fixed-cidr":"172.17.0.0/24"
  #  }
  #EOF

  # 启动docker引擎并设置开机启动
  sudo systemctl daemon-reload
  sudo systemctl start docker
  sudo systemctl enable docker
  # 配置当前用户对docker的执行权限

  sudo groupadd docker
  sudo gpasswd -a ${USER} docker

  #sudo systemctl daemon-reload
}

function setupRelayAgent() {
  echo "--------------------Start deploying protocol agent-----------------------------"
  RelayDir=/home/relay

  kill_if_process_exist decept-agent

  #创建下载路径文件夹
  if [ ! -d "${RelayDir}" ]; then
    echo "${RelayDir} not exist, start create"
    mkdir -p ${RelayDir}
    if [ $? -ne 0 ]; then
      echo "Failed to create folder ${RelayDir}"
      exit 1
    else
      sudo chmod -R 755 ${RelayDir}
      echo "Create folder ${RelayDir} successfully"
    fi
  fi

  cp ${Project_Dir}/agent/decept-agent.tar.gz ${RelayDir}/

  # TODO 替换redis 配置

  cd ${RelayDir}/

  tar zxvf ${RelayDir}/decept-agent.tar.gz

  cd ${RelayDir}/agent

  resetAgentJson

  nohup ./decept-agent -mode RELAY >/dev/null &

  cd $Project_Dir
  # TODO cp 协议代理 到指定目录

  echo "--------------------Protocol agent deployment complete-----------------------------"
  setupProxyFile
}

function setupProxyFile() {
  echo "--------------------Start deploying protocol file-----------------------------"
  ProtocolPath=/home/ehoney_proxy
  #创建下载路径文件夹
  if [ -d "${ProtocolPath}" ]; then
    rm -rf ${ProtocolPath}
  fi
  mkdir -p ${ProtocolPath}
  if [ $? -ne 0 ]; then
    echo "Failed to create folder ${ProtocolPath}"
    exit 1
  else
    echo "Create folder ${ProtocolPath} successfully"
  fi

  cp -r ${Project_Dir}/tool/protocol/* ${ProtocolPath}
  sudo chmod -R 755 /home/ehoney_proxy/
  echo "--------------------Protocol file deployment complete-----------------------------"
}

function setupMysqlDocker() {
  echo "--------------------Start installing the database container-----------------------------"
  dos2unix ${Project_Dir}/tool/mysql-docker/setup.sh
  dos2unix ${Project_Dir}/tool/mysql-docker/privileges.sql
  sed -i "3c update user set authentication_string = password('${DB_Password}') where user = 'root';" $Project_Dir/tool/mysql-docker/privileges.sql
  sed -i "5c GRANT ALL PRIVILEGES ON *.* TO 'root'@'%' IDENTIFIED BY '${DB_Password}' WITH GRANT OPTION;" $Project_Dir/tool/mysql-docker/privileges.sql
  cd $Project_Dir/tool/mysql-docker
  docker build -t ehoney-mysql .
  db_data_dir=/var/lib/ehoney-db-data
  if [ ! -d "${db_data_dir}" ]; then
    echo "start mkdir $db_data_dir"
    mkdir $db_data_dir
    sudo chmod -R 777 $db_data_dir
  fi
  docker run -d -p $DB_Port:3306 -v $db_data_dir:/var/lib/mysql -e TZ=Asia/Shanghai --name ehoney-mysql ehoney-mysql:latest
  check_docker_container_state ehoney-mysql
  echo "--------------------End of database container installation-----------------------------"
  cd $Project_Dir
}

function setupDeceptDefence() {
  echo "--------------------Start install DeceptDefence-------------------------"
  cd $Project_Dir
  chmod +x $Project_Dir/dockerStart.sh
  dos2unix $Project_Dir/dockerStart.sh
  chmod -R 755 $Project_Dir/tool/file_token_trace
  # 覆盖项目中的k3s 配置
  yes | cp -rf /etc/rancher/k3s/k3s.yaml ${Project_Dir}/configs/.kube/config
  docker build -t decept-defense .
  docker run -d -v $Project_Dir/configs/:/go/src/configs/ -v $Project_Dir/upload/:/go/src/upload/ --network host -e TZ=Asia/Shanghai --name decept-defense-web -e CONFIGS="host:${Local_Host};" decept-defense:latest
  echo "--------------------End of DeceptDefence installation-------------------------"
  check_decept_defense_service
}

function setupFileTrace() {
  echo "--------------------Start install FileTrace---------------------------"
  chmod +x $Project_Dir/tool/filetrace/filetracemsg
  cd $Project_Dir/tool/filetrace
  echo "filetrace param: -dbuser ${DB_User} -dbpassword ${DB_Password} -dbhost ${Local_Host} -dbname ${DB_Database} -dbport ${DB_Port} "
  nohup ./filetracemsg -dbuser ${DB_User} -dbpassword ${DB_Password} -dbhost ${Local_Host} -dbname ${DB_Database} -dbport ${DB_Port} >/dev/null &
  check_file_trace_service
  echo "--------------------End of  FileTrace installation---------------------------"
}

function setupFileTraceDocker() {
  echo "--------------------Start install FileTrace---------------------------"
  cd $Project_Dir/tool/filetrace
  docker build -t filetracemsg .
  docker run -itd -p 5000:5000 --name filetrace-msg-service -v /home/filetrace/infile:/mnt/infile -v /home/filetrace/outfile:/mnt/outfile -e SQLALCHEMY_DATABASE_URI="mysql+pymysql://${DB_User}:${DB_Password}@${Local_Host}:${DB_Port}/${DB_Database}?charset=utf8" filetracemsg:latest
  check_docker_container_state filetrace-msg-service
  echo "--------------------End of FileTrace installation---------------------------"
}

function main() {
  echo "Notice [If there is a coding problem during startup, Please install dos2unix and execute dos2unix quick-start.sh]"
  setUpIp
  Project_Dir=$(
    cd $(dirname $0)
    pwd
  )
  yum install -y dos2unix
  prepare_conf
  ports_check
  #setup_iptables
  prepare_base_install
  component_installer
  echo "----------------------------------------------------------"
  echo "all the services are ready and happy to use!!!"
  echo "Please visit url: http://${Local_Host}:${Project_Front_Port}/decept-defense"
  echo "----------------------------------------------------------"
}
function yumRepoUpdate() {
  wget http://mirrors.163.com/.help/CentOS7-Base-163.repo
  mv /etc/yum.repos.d/CentOS-Base.repo /etc/yum.repos.d/CentOS-Base.repo.backup
  mv CentOS7-Base-163.repo /etc/yum.repos.d/CentOS-Base.repo
  yum clean all
  yum makecache
}

function getIpAddr() {
  # 获取IP命令
  ipaddr=$(ifconfig -a | grep inet | grep -v 127.0.0.1 | grep -v inet6 | awk '{print $2}' | tr -d "addr:"​)
  array=($(echo $ipaddr | tr '\n' ' ')) # IP地址分割，区分是否多网卡
  num=${#array[@]}                      #获取数组元素的个数

  # 选择安装的IP地址
  if [ $num -eq 1 ]; then
    #echo "*单网卡"
    Local_Host=${array[*]}
  elif [ $num -gt 1 ]; then
    echo -e "\033[036m---------------------------------------------------------\033[0m"
    echo -e "\033[036m----Please select the IP address used by this machine---\033[0m"
    for i in "${!array[@]}"; do
      echo -e "\033[032m*    "$i" "${array[$i]}"	\033[0m"
    done
    #选择需要安装的服务类型
    input="1"
    while :; do
      read -r -p "Please select the IP address (serial number) to use: " input
      if [ "${input}" != "" ]; then
        break
      fi
    done
    Local_Host=${array[$input]}
  else
    echo -e "IP of network card is not set, please check the server environment!"
    exit 1
  fi
}

# 校验IP地址合法性
function isValidIp() {
  local ip=$1
  local ret=1
  if [[ $ip =~ ^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$ ]]; then
    ip=(${ip//\./ }) # 按.分割，转成数组，方便下面的判断
    [[ ${ip[0]} -le 255 && ${ip[1]} -le 255 && ${ip[2]} -le 255 && ${ip[3]} -le 255 ]]
    ret=$?
  fi
  return $ret
}

function setUpIp() {
  getIpAddr               #自动获取IP
  isValidIp ${Local_Host} # IP校验
  if [ $? -ne 0 ]; then
    echo -e "\033[31m*The IP address obtained automatically is invalid, please try again! \033[0m"
    exit 1
  fi
  echo "The IP used by this machine is set to: ${Local_Host}"
  #
  #  if [ ! -f "/home/ehoney_ip.txt" ]; then
  #    touch /home/ehoney_ip.txt
  #  fi
  #  echo ${Local_Host} >/home/ehoney_ip.txt
}

main
