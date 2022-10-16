#!/bin/bash
#

DB_User="root"
DB_Database="sec_ehoneypot"
DB_Password="Ehoney2021"
Project_Dir=/home
Project_Port=8082
DB_Port=3306
Local_Host="127.0.0.1"

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

function kill_if_process_exist2() {
  PROC_NAME=$1
  echo "--------Start killing $PROC_NAME process and its child processes---------"
  ProcNumber=$(ps -ef | grep $PROC_NAME | grep -v "grep" | awk '{print $2}')
  if [ $ProcNumber ]; then
    echo "进程ID: $ProcNumber"
    kill -9 $ProcNumber
    echo "--------------------End of killing process---------------------------"
  fi
}

# 检查软件是否安装 curl wget zip go redis mysql docker kubectl;
function prepare_base_install() {
  for i in wget curl unzip gcc gcc-c++ kernel-devel-$(uname -r); do
    command -v $i &>/dev/null || install_soft $i
  done
  # yumRepoUpdate
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
  echo $docker_container_name
  processor=$(docker ps -a | grep $docker_container_name)
  if [ "$processor" != "" ]; then
    echo "stop and rm docker container [$docker_container_name]..."
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
  httpwebport=$(is_port_bind $Project_Port)
  mysqlport=$(is_port_bind $DB_Port)

  if [ "$httpwebport" = "true" ]; then
    echo "$Project_Port Occupied!!"
    shouldreturn="true"
  fi

  if [ "$mysqlport" = "true" ]; then
    echo "$DB_Port Occupied!!"
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
  echo -e "DB_Port: "$DB_Port
  echo -e "DB_User: "$DB_User
  echo -e "DB_Database: "$DB_Database
  echo -e "DB_Password: "$DB_Password
}

function component_installer() {
  setupDocker # 安装Docker
  setupK3s    # 安装K3S
  setupFalco
  setupMysqlDocker # 安装MySQL 5.6 以及sec_pot database 的镜像
  setupDeceptDefence
}
function setupFalco() {
  echo "start deploying falco >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
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
  echo "end of falco installation >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
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

function install_k3s() {
  # K3S 配置修改
  echo "start setup k3s service"
  dos2unix ${Project_Dir}/tool/k3s/k3s.sh
  sh ${Project_Dir}/tool/k3s/k3s.sh
  sed -i '28c ExecStart=/usr/local/bin/k3s server --docker --no-deploy traefik' /etc/systemd/system/k3s.service
  sed -i '29c #' /etc/systemd/system/k3s.service
  sudo systemctl daemon-reload
  sudo systemctl restart k3s
  swpConfValue /etc/rancher/k3s/k3s.yaml 5 127.0.0.1 $Local_Host
  if [ $(grep -c "KUBECONFIG" /etc/profile) -ge '1' ]; then
    echo "KUBECONFIG is configured, skip!"
  else
    echo \  >>/etc/profile
    echo export KUBECONFIG=/etc/rancher/k3s/k3s.yaml >>/etc/profile
  fi
  # 覆盖项目中的k3s 配置
  yes | cp -rf /etc/rancher/k3s/k3s.yaml ${Project_Dir}/configs/.kube/config
  source /etc/profile
  sleep 1s
  check_k3s_service
}

function setupK3s() {
  echo "start deploying k3s >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
  if [ -e "/etc/rancher/k3s/k3s.yaml" ]; then
    echo "k3s is installed, skip! "
  else
    echo "k3s is uninstalled, start install!"
    install_k3s
  fi
  echo "end of k3s installation >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
}

function setupDocker() {
  echo docker version
  if [ $? -ne 0 ]; then
    echo "docker is uninstalled, start install!"
    install_Docker
  else
    if [ $(grep -c "dockerd" /usr/lib/systemd/system/docker.service) -eq '1' ]; then
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
  echo "start deploying docker >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
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

  sudo tee /etc/docker/daemon.json <<-'EOF'
    {
     "registry-mirrors": ["https://docker.mirrors.ustc.edu.cn"]
    }
EOF

  # 启动docker引擎并设置开机启动
  sudo systemctl daemon-reload
  sudo systemctl start docker
  sudo systemctl enable docker
  # 配置当前用户对docker的执行权限
  sudo groupadd docker
  sudo gpasswd -a ${USER} docker
}

function setupMysqlDocker() {
  echo "start installing the database container >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
  stop_docker_container_if_exist ehoney-mysql
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
  echo "end of database container installation >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
  cd $Project_Dir
}

function install_go() {
  sudo cd $Project_Dir
  sudo cp -r tool/go /usr/local/go
  sudo chmod -R 755 /usr/local/go
  echo \  >>/etc/profile
  echo PATH=$PATH:/usr/local/go/bin >>/etc/profile
  echo \  >>/etc/profile
  echo export GO111MODULE=on >>/etc/profile
  echo \  >>/etc/profile
  echo export GOPROXY=https://goproxy.cn >>/etc/profile
  source /etc/profile
  sleep 1s
  go version
}

function setupDeceptDefence() {
  echo "start install ehoney server >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
  kill_if_process_exist decept-defense
  if [ ! -d "/usr/local/go" ]; then
    echo "golang is uninstalled, start install!"
    install_go
  else
    echo "golang is installed, skip! "
  fi
  cd $Project_Dir
  yes | cp -rf /etc/rancher/k3s/k3s.yaml ${Project_Dir}/configs/.kube/config
  mkdir -p /var/decept-agent/ssh/
  sudo chmod -R 755 $Project_Dir/tool/protocol
  go build .
  nohup ./decept-defense --ip ${Local_Host} >/dev/null &
  check_decept_defense_service
  echo "end of ehoney server installation >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
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
  #ports_check
  prepare_base_install
  component_installer
  echo "--------------------------------------------------------------"
  echo "all the services are ready and happy to use!!!"
  echo "please set the correct system time zone!!!!!!!!"
  echo "please visit url: [ http://${Local_Host}:8082/decept-defense ]"
  echo "--------------------------------------------------------------"
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
    echo -e "\033[036m-本次更新涉及数据库结构改变, 如果由老版本升级推荐删除文件夹/var/lib/ehoney-db-data \033[0m"
    echo -e "\033[036m-本次更新涉及数据库结构改变, 如果由老版本升级推荐删除文件夹/var/lib/ehoney-db-data \033[0m"
    echo -e "\033[036m-本次更新涉及数据库结构改变, 如果由老版本升级推荐删除文件夹/var/lib/ehoney-db-data \033[0m"
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
}

main
