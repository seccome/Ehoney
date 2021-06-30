package k3s

import (
	"context"
	"decept-defense/models/honeycluster"
	"decept-defense/models/policyCenter"
	"decept-defense/models/redisCenter"
	"decept-defense/models/util"
	"decept-defense/models/util/comm"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"io"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"os"
	"strings"
	"time"
)

/**
根据podname 获取pod信息
*/
func GetPodinfo(podname string) *v1.Pod {
	kubeconfig := beego.AppConfig.String("kubeconfig")
	kubeconfig = util.GetCurrentPathString() + kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		logs.Error("[GetPodinfo] err:", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	pods, err := clientset.CoreV1().Pods("default").Get(context.TODO(), podname, metav1.GetOptions{})
	return pods

}

func FreshPods() {
	kubeconfig := beego.AppConfig.String("kubeconfig")
	kubeconfig = util.GetCurrentPathString() + kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		logs.Error("FreshPods err:", err)
	}
	clientset, err := kubernetes.NewForConfig(config)

	pods, err := clientset.CoreV1().Pods(v1.NamespaceDefault).List(context.TODO(), metav1.ListOptions{})

	for i := 0; i < len(pods.Items); i++ {
		CheckPodStatus(pods.Items[i].Name, v1.NamespaceDefault)

		//pstatus := 2
		//podsimage := ""
		//podsourcename := ""
		////var podport int32
		//podname := pods.Items[i].Name
		//podnamepaces := pods.Items[i].Namespace
		//poduid := pods.Items[i].UID
		//podip := pods.Items[i].Status.PodIP
		//if len(pods.Items[i].Spec.Containers) > 0 {
		//	podsourcename = pods.Items[i].Spec.Containers[0].Name
		//	podsimage = pods.Items[i].Spec.Containers[0].Image
		//	//if len(pods.Items[i].Spec.Containers[0].Ports) > 0 {
		//	//	podport = pods.Items[i].Spec.Containers[0].Ports[0].ContainerPort
		//	//} else {
		//	//	podport = 0
		//	//}
		//}
		//if len(pods.Items[i].Status.ContainerStatuses) > 0 {
		//	if pods.Items[i].Status.ContainerStatuses[0].State.Running != nil {
		//		pstatus = 1
		//	} else {
		//		pstatus = 2
		//	}
		//}
		//honeycluster.FreshPodInfo(podsourcename, podname, podnamepaces, podip, string(poduid), podsimage, pstatus)

	}
	time.Sleep(1 * time.Minute)
	go FreshPods()
}

// 新增蜜罐流量转发策略，insert数据库，下发策略到Redis
func CreateHoneyTransPolicyHandler(agentId string, honeyTypeId string, forwardPort int, honeyPotId string, honeyPotPort int, creator string, status int, path string, serverId string) {
	// 前端json转换为结构体
	//var transJson comm.TransparentTransponderJson
	// 前端json转换成策略结构体
	var transPolicyJson comm.TransparentTransponderPolicyJson
	taskid := util.GetUUID()
	currentTime := time.Now().Unix()
	policyCenter.InsertHoneyTransPolicy(taskid, agentId, forwardPort, honeyPotPort, honeyPotId, currentTime, creator, status, comm.AgentTypeRelay, path, serverId)
	transPolicyJson.ServerType = policyCenter.SelectHoneyPotType(honeyTypeId)
	transPolicyJson.HoneyIP = policyCenter.SelectHoneyPotIP(honeyPotId)
	transPolicyJson.TaskId = taskid
	transPolicyJson.HoneyPort = honeyPotPort
	transPolicyJson.ListenPort = forwardPort
	transPolicyJson.Status = status
	transPolicyJson.AgentId = agentId
	transPolicyJson.Type = comm.AgentTypeRelay
	transPolicyJson.Path = path
	// 策略结构体转换成策略json
	transPolicy, err2 := json.Marshal(transPolicyJson)
	if err2 != nil {
		logs.Error(fmt.Sprintf("[CreateHoneyTransPolicyHandler] Failed to marshal data: %v", err2))
		return
	}
	redisCenter.RedisSubProducerTransPolicy(string(transPolicy))
	return
}

// 删除蜜罐流量转发策略，insert数据库，下发策略到Redis
func DeleteHoneyTransPolicyHandler(taskId string, status int) {
	// 前端json转换为结构体
	//var transJson comm.TransOfflineJson
	currentTime := time.Now().Unix()
	policyCenter.UpdateHoneyTransPolicy(taskId, currentTime, status)
	policy := policyCenter.SelectHoneyPotInfoByTaskIdHoneyTrans(taskId)
	redisCenter.RedisSubProducerTransPolicy(policy)
	return
}

func ExecPodCmd(honeypotid string, cmdstr []string) error {
	honeypotinfo := honeycluster.SelectHoneyInfoById(honeypotid)
	kubeconfig := beego.AppConfig.String("kubeconfig")
	kubeconfig = util.GetCurrentPathString() + "/" + kubeconfig
	errinfo := fmt.Errorf("%s", "")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		errinfo = fmt.Errorf("%s", "ExecPodCmd clientcmd.BuildConfigFromFlags ERR:", err)
		logs.Error("ExecPodCmd clientcmd.BuildConfigFromFlags ERR:", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		errinfo = fmt.Errorf("ExecPodCmd kubernetes.NewForConfig ERR:", err)
		logs.Error("ExecPodCmd kubernetes.NewForConfig ERR:", err)
	}
	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(util.Strval(honeypotinfo[0]["podname"])).
		Namespace("default").
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Command: cmdstr,
			Stdin:   true,
			Stdout:  true,
			Stderr:  true,
			TTY:     false,
		}, scheme.ParameterCodec)
	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		errinfo = fmt.Errorf("cmd err:", err)
		logs.Error("cmd err:", err)
	}
	screen := struct {
		io.Reader
		io.Writer
	}{os.Stdin, os.Stdout}

	if err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  screen,
		Stdout: screen,
		Stderr: screen,
		Tty:    false,
	}); err != nil {
		errinfo = fmt.Errorf("[ExecPodCmd] exec Error:", err)
		logs.Error("[ExecPodCmd] exec Error:", err)
	}
	if errinfo.Error() != "" {
		return errinfo
	} else {
		return nil
	}
}

/**
Response:
*v1.Pod
bool 是否成功
error
int 0是进行中，1是成功 2是失败
int
*/
func CheckPodStatus(podname string, namespace string) (*v1.Pod, bool, error, int, int) {
	issuccess := false
	createstatus := 2
	errcode := 2
	kubeconfig := beego.AppConfig.String("kubeconfig")
	kubeconfig = util.GetCurrentPathString() + kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		fmt.Println("[CheckPodStatus] BuildConfigFromFlags err:", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	pods, err := clientset.CoreV1().Pods(namespace).Get(context.TODO(), podname, metav1.GetOptions{})
	if err != nil {
		logs.Error("[CheckPodStatus] pod get err", err)
		issuccess = false
		createstatus = 2
		errcode = 0
	} else {
		if pods != nil {

			//if pods.Status.Phase == "Pending" {
			//	err = fmt.Errorf("%s", podname+"创建异常："+" 状态 Pending 可能 docker服务关闭")
			//	honeycluster.FreshPodInfo(pods.Spec.Containers[0].Name, pods.Name, pods.Namespace, pods.Status.PodIP, string(pods.UID), pods.Spec.Containers[0].Image, 2)
			//	logs.Error(podname + "创建异常：" + " 状态 Pending 可能 docker服务关闭")
			//	issuccess = false
			//	createstatus = 2
			//	errcode = 1
			//	//return pods, false,err,2, 1
			//} else
			if len(pods.Status.ContainerStatuses) >0 {
				if pods.Status.ContainerStatuses[0].State.Waiting != nil && pods.Status.ContainerStatuses[0].State.Waiting.Reason == "ContainerCreating" {
					issuccess = true
					createstatus = 0
					errcode = 2
				} else if pods.Status.ContainerStatuses[0].State.Terminated != nil && pods.Status.ContainerStatuses[0].State.Terminated.Reason == "Completed" {
					issuccess = true
					createstatus = 0
					errcode = 2
				} else if pods.Status.ContainerStatuses[0].State.Waiting != nil && pods.Status.ContainerStatuses[0].State.Waiting.Reason == "ErrImagePull" {
					err = fmt.Errorf("%s", podname+"创建异常："+" 状态 ErrImagePull 镜像拉取失败")
					logs.Error(podname + "创建异常：" + " 状态 ErrImagePull 镜像拉取失败")
					honeycluster.FreshPodInfo(pods.Spec.Containers[0].Name, pods.Name, pods.Namespace, pods.Status.PodIP, string(pods.UID), pods.Spec.Containers[0].Image, 2)
					issuccess = false
					createstatus = 2
					errcode = 3
					//return pods, false,err,2, 3
				} else if pods.Status.ContainerStatuses[0].State.Waiting != nil && pods.Status.ContainerStatuses[0].State.Waiting.Reason == "CrashLoopBackOff" {
					err = fmt.Errorf("%s", podname+"创建异常："+" 状态 CrashLoopBackOff 容器创建失败")
					logs.Error(podname + "创建异常：" + " 状态 CrashLoopBackOff 容器创建失败")
					honeycluster.FreshPodInfo(pods.Spec.Containers[0].Name, pods.Name, pods.Namespace, pods.Status.PodIP, string(pods.UID), pods.Spec.Containers[0].Image, 2)
					issuccess = false
					createstatus = 2
					errcode = 4
					//return pods, false,nil,2, 4
				} else if pods.Status.ContainerStatuses[0].State.Running != nil && pods.Status.ContainerStatuses[0].State.Running != nil {
					honeycluster.FreshPodInfo(pods.Spec.Containers[0].Name, pods.Name, pods.Namespace, pods.Status.PodIP, string(pods.UID), pods.Spec.Containers[0].Image, 1)
					err = fmt.Errorf("%s", pods.Spec.Containers[0].Name+" 创建成功！")
					issuccess = true
					createstatus = 1
					errcode = 5
					//return pods, true,nil,1, 5
				}
			}
		}
	}
	return pods, issuccess, err, createstatus, errcode
}

func IsPodRunningV3(podName, namespace string) (*v1.Pod, bool, error, int, int) {
	var containerName string = ""
	kubeconfig := beego.AppConfig.String("kubeconfig")
	kubeconfig = util.GetCurrentPathString() + kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		logs.Error("kubernetes err:", err)
		return nil, false, err, 2, 10
	}
	pods, err := client.CoreV1().Pods(v1.NamespaceDefault).List(context.TODO(), metav1.ListOptions{})
	for i, d := range pods.Items {

		if strings.Contains(d.Name, podName+"-") {
			containerName = d.Name
			break
		} else if i == len(pods.Items)-1 {
			return nil, false, nil, 2, 9
		} else {
			continue
		}
	}

	return CheckPodStatus(containerName, namespace)

	//pod, err := client.CoreV1().Pods(namespace).Get(context.TODO(), containerName, metav1.GetOptions{})
	//if err != nil {
	//	return nil, false, err, 0
	//}
	//if len(pod.Status.ContainerStatuses) > 0 {
	//	if pod.Status.ContainerStatuses[0].State.Running != nil {
	//		return pod, true, nil, 1
	//	} else if pod.Status.ContainerStatuses[0].State.Waiting != nil {
	//		err = fmt.Errorf("%s", pod.Status.ContainerStatuses[0].State.Waiting.Reason)
	//		return pod, false, err, 2
	//	} else if pod.Status.ContainerStatuses[0].State.Terminated != nil {
	//		err = fmt.Errorf("%s", pod.Status.ContainerStatuses[0].State.Terminated.Reason)
	//		return pod, false, err, 1
	//	} else {
	//		return pod, false, nil, 2
	//	}
	//}
	//return pod, false, nil, 0
}

func IsPodRunningV2(podName, namespace string) (*v1.Pod, bool, error, int) {
	var containerName string = ""
	kubeconfig := beego.AppConfig.String("kubeconfig")
	kubeconfig = util.GetCurrentPathString() + kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		logs.Error("kubernetes err:", err)
		return nil, false, nil, 0
	}
	pods, err := client.CoreV1().Pods(v1.NamespaceDefault).List(context.TODO(), metav1.ListOptions{})
	for i, d := range pods.Items {

		if strings.Contains(d.Name, podName+"-") {
			containerName = d.Name
			break
		} else if i == len(pods.Items)-1 {
			return nil, false, nil, 0
		} else {
			continue
		}
	}
	pod, err := client.CoreV1().Pods(namespace).Get(context.TODO(), containerName, metav1.GetOptions{})
	if err != nil {
		return nil, false, err, 0
	}
	if len(pod.Status.ContainerStatuses) > 0 {
		if pod.Status.ContainerStatuses[0].State.Running != nil {
			return pod, true, nil, 1
		} else if pod.Status.ContainerStatuses[0].State.Waiting != nil {
			err = fmt.Errorf("%s", pod.Status.ContainerStatuses[0].State.Waiting.Reason)
			return pod, false, err, 2
		} else if pod.Status.ContainerStatuses[0].State.Terminated != nil {
			err = fmt.Errorf("%s", pod.Status.ContainerStatuses[0].State.Terminated.Reason)
			return pod, false, err, 1
		} else {
			return pod, false, nil, 2
		}
	}
	return pod, false, nil, 0
}

func IsPodRunning(podName, namespace string) wait.ConditionFunc {
	var containerName string = ""
	return func() (bool, error) {
		kubeconfig := beego.AppConfig.String("kubeconfig")
		kubeconfig = util.GetCurrentPathString() + kubeconfig
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		client, err := kubernetes.NewForConfig(config)
		if err != nil {
			logs.Error("kubernetes err:", err)
		}
		pods, err := client.CoreV1().Pods(v1.NamespaceDefault).List(context.TODO(), metav1.ListOptions{})
		for i, d := range pods.Items {

			if strings.Contains(d.Name, podName) {
				containerName = d.Name
				break
			} else if i == len(pods.Items)-1 {
				return false, nil
			} else {
				continue
			}
		}
		pod, err := client.CoreV1().Pods(namespace).Get(context.TODO(), containerName, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		if len(pod.Status.ContainerStatuses) > 0 {
			if pod.Status.ContainerStatuses[0].State.Running != nil {
				return true, nil
			} else {
				err = fmt.Errorf("%s", pod.Status.ContainerStatuses[0].State.Waiting.Reason)
				return false, err
			}
		}
		return false, nil
	}
}

func WaitForPodRunning(namespace, podName string, timeout time.Duration) error {
	return wait.PollImmediate(time.Second, timeout, IsPodRunning(podName, namespace))
}
