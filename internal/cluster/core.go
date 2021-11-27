package cluster

import (
	"archive/tar"
	"context"
	"decept-defense/pkg/util"
	"fmt"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	appsV1 "k8s.io/api/apps/v1"
	apiV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	v1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	//metrics "k8s.io/client-go/tools/metrics"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/client-go/util/retry"
	_ "k8s.io/kubectl/pkg/cmd/cp"
	cmdUtil "k8s.io/kubectl/pkg/cmd/util"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

var (
	kubeConfig        = path.Join(util.WorkingPath(), "configs/.kube/config")
	deploymentsClient v1.DeploymentInterface
	client            *kubernetes.Clientset
	config            *rest.Config
)

func Setup() {
	var err error
	config, err = clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		zap.L().Fatal("BuildConfigFromFlags fail error:" + err.Error())
	}
	client, err = kubernetes.NewForConfig(config)
	if err != nil {
		zap.L().Fatal("BuildConfigFromFlags fail error:" + err.Error())
	}
	deploymentsClient = client.AppsV1().Deployments(apiV1.NamespaceDefault)
}

func CreateDeployment(podName, imageAddress string, containerPort int32) (*apiV1.Pod, error) {
	zap.L().Info(fmt.Sprintf("podName: %s imageAddress: %s, containerPort: %d", podName, imageAddress, containerPort))

	deployment := &appsV1.Deployment{
		ObjectMeta: metaV1.ObjectMeta{
			Name: podName,
		},
		Spec: appsV1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metaV1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "demo",
				},
			},
			Template: apiV1.PodTemplateSpec{
				ObjectMeta: metaV1.ObjectMeta{
					Labels: map[string]string{
						"app": "demo",
					},
				},
				Spec: apiV1.PodSpec{
					Containers: []apiV1.Container{
						{
							Name:  podName,
							Image: imageAddress,
							Ports: []apiV1.ContainerPort{
								{
									Name:          podName,
									Protocol:      apiV1.ProtocolTCP,
									ContainerPort: containerPort,
									//HostPort:      hostPort,
								},
							},
						},
					},
				},
			},
		},
	}

	if containerPort == 21 {

		var dataPort = apiV1.ContainerPort{
			Name:          "data",
			ContainerPort: 21100,
			HostPort:      21100,
		}

		var env1 = apiV1.EnvVar{
			Name:  "FTP_USER",
			Value: "admin",
		}

		var env2 = apiV1.EnvVar{
			Name:  "FTP_PASS",
			Value: "admin",
		}

		var env3 = apiV1.EnvVar{
			Name:  "PASV_ADDRESS",
			Value: "192.168.240.160",
		}

		var env4 = apiV1.EnvVar{
			Name:  "PASV_PROMISCUOUS",
			Value: "YES",
		}
		var env5 = apiV1.EnvVar{
			Name:  "PASV_MIN_PORT",
			Value: "21100",
		}
		var env6 = apiV1.EnvVar{
			Name:  "PASV_MAX_PORT",
			Value: "21100",
		}
		deployment.Spec.Template.Spec.Containers[0].Ports[0].HostPort = 21
		deployment.Spec.Template.Spec.Containers[0].Ports = append(deployment.Spec.Template.Spec.Containers[0].Ports, dataPort)
		deployment.Spec.Template.Spec.Containers[0].Env = append(deployment.Spec.Template.Spec.Containers[0].Env, env1)
		deployment.Spec.Template.Spec.Containers[0].Env = append(deployment.Spec.Template.Spec.Containers[0].Env, env2)
		deployment.Spec.Template.Spec.Containers[0].Env = append(deployment.Spec.Template.Spec.Containers[0].Env, env3)
		deployment.Spec.Template.Spec.Containers[0].Env = append(deployment.Spec.Template.Spec.Containers[0].Env, env4)
		deployment.Spec.Template.Spec.Containers[0].Env = append(deployment.Spec.Template.Spec.Containers[0].Env, env5)
		deployment.Spec.Template.Spec.Containers[0].Env = append(deployment.Spec.Template.Spec.Containers[0].Env, env6)

		//deployment.Spec.Template.Spec.Containers[0].Env[0].Name = "FTP_USER"
		//deployment.Spec.Template.Spec.Containers[0].Env[0].Value = "admin"
		//deployment.Spec.Template.Spec.Containers[0].Env[1].Name = "FTP_PASS"
		//deployment.Spec.Template.Spec.Containers[0].Env[1].Value = "admin"
		//deployment.Spec.Template.Spec.Containers[0].Env[2].Name = "PASV_ADDRESS"
		//deployment.Spec.Template.Spec.Containers[0].Env[2].Value = ""
		//deployment.Spec.Template.Spec.Containers[0].Env[3].Name = "PASV_PROMISCUOUS"
		//deployment.Spec.Template.Spec.Containers[0].Env[3].Value = "YES"
		//deployment.Spec.Template.Spec.Containers[0].Env[4].Name = "PASV_MIN_PORT"
		//deployment.Spec.Template.Spec.Containers[0].Env[4].Value = "21100"
		//deployment.Spec.Template.Spec.Containers[0].Env[5].Name = "PASV_MAX_PORT"
		//deployment.Spec.Template.Spec.Containers[0].Env[5].Value = "21100"
		zap.L().Info(fmt.Sprintf("ftp deployment: %v", deployment.Spec.Template.Spec.Containers[0].Ports[0]))
	}
	zap.L().Info("Creating deployment...")

	zap.L().Info(fmt.Sprintf("%v", deployment))

	_, err := deploymentsClient.Create(context.TODO(), deployment, metaV1.CreateOptions{})
	if err != nil {
		zap.L().Error("Creating deployment error : " + err.Error())
		return nil, err
	}
	//TODO remove wait logic
	err = waitForPodStartBuild(apiV1.NamespaceDefault, podName, time.Minute)
	if err != nil {
		zap.L().Error("waitForPodStartBuild error : " + err.Error())
		return nil, err
	}
	return GetPod(podName)
}

func UpdateDeployment(deploymentName string, containerPort int32, imageAddress string) error {

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := deploymentsClient.Get(context.TODO(), deploymentName, metaV1.GetOptions{})
		if getErr != nil {
			return getErr
		}
		result.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort = containerPort
		result.Spec.Template.Spec.Containers[0].Image = imageAddress
		_, updateErr := deploymentsClient.Update(context.TODO(), result, metaV1.UpdateOptions{})
		return updateErr
	})
	if retryErr != nil {
		return retryErr
	}
	return nil
}

func DeleteDeployment(deploymentName string) error {
	deletePolicy := metaV1.DeletePropagationForeground
	if err := deploymentsClient.Delete(context.TODO(), deploymentName, metaV1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		return err
	}
	return nil
}

func DeploymentIsExist(deploymentName string) (bool, error) {
	pods, err := client.CoreV1().Pods(apiV1.NamespaceDefault).List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return false, err
	}
	for _, d := range pods.Items {
		data := strings.Split(d.Name, "-")
		if len(data) == 0 {
			continue
		}
		if data[0] == deploymentName {
			return true, nil
		}
	}
	return false, nil
}

func GetPod(podName string) (*apiV1.Pod, error) {
	pods, err := client.CoreV1().Pods(apiV1.NamespaceDefault).List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, d := range pods.Items {
		data := strings.Split(d.Name, "-")
		if len(data) == 0 {
			continue
		}
		if strings.HasPrefix(d.Name, podName) && strings.Contains(d.Name, podName) {
			pod, err := client.CoreV1().Pods(apiV1.NamespaceDefault).Get(context.TODO(), d.Name, metaV1.GetOptions{})
			if err != nil {
				return nil, err
			}
			return pod, nil
		}
	}
	return nil, err
}

func isPodRunning(podName, namespace string) wait.ConditionFunc {
	var containerName string = ""
	return func() (bool, error) {
		pods, _ := client.CoreV1().Pods(apiV1.NamespaceDefault).List(context.TODO(), metaV1.ListOptions{})
		for i, d := range pods.Items {
			data := strings.Split(d.Name, "-")
			if len(data) == 0 {
				continue
			}
			if strings.HasPrefix(d.Name, podName) && strings.Contains(d.Name, podName) {
				containerName = d.Name
				break
			} else if i == len(pods.Items)-1 {
				return false, nil
			} else {
				continue
			}
		}
		if containerName != "" {
			return true, nil
		}
		return false, nil
	}
}

type PodDetail struct {
	Status string
	HostIP string
	PodIP  string
}

func GetPodDetailInfo(podName string) PodDetail {
	pod, _ := client.CoreV1().Pods(apiV1.NamespaceDefault).Get(context.TODO(), podName, metaV1.GetOptions{})

	status := ""
	if pod.Status.Phase == apiV1.PodRunning {
		status = "Running"
	} else if pod.Status.Phase == apiV1.PodFailed {
		status = "Failed"
	} else {
		status = "Unknown"
	}

	return PodDetail{
		Status: status,
		HostIP: pod.Status.HostIP,
		PodIP:  pod.Status.PodIP,
	}
}

func waitForPodStartBuild(namespace, podName string, timeout time.Duration) error {
	return wait.PollImmediate(time.Second, timeout, isPodRunning(podName, namespace))
}

func RemoveFromPod(podName, containerName, srcPath string) error {
	if err := checkDestinationIsFile(podName, containerName, srcPath); err == nil {
		return exec(podName, containerName, []string{"rm", "-rf", srcPath})
	}
	return nil
}

func CopyToPod(podName, containerName, srcPath, destPath string) error {

	zap.L().Error(fmt.Sprintf("podName[%s] containerName[%s]  srcPath[%s] destPath[%s]", podName, containerName, srcPath, destPath))

	reader, writer := io.Pipe()
	if destPath != "/" && strings.HasSuffix(string(destPath[len(destPath)-1]), "/") {
		destPath = destPath[:len(destPath)-1]
	}
	createDir(podName, containerName, destPath)
	destPath = destPath + "/" + path.Base(srcPath)
	go func() {
		defer writer.Close()
		err := makeTar(srcPath, destPath, writer)
		zap.L().Error(fmt.Sprintf("makeTar err: %v", err))
		cmdUtil.CheckErr(err)
	}()
	var cmdArr []string

	cmdArr = []string{"tar", "-xf", "-"}
	destDir := path.Dir(destPath)
	if len(destDir) > 0 {
		cmdArr = append(cmdArr, "-C", destDir)
	}
	req := client.CoreV1().RESTClient().
		Post().
		Namespace(apiV1.NamespaceDefault).
		Resource("pods").
		Name(podName).
		SubResource("exec").
		VersionedParams(&apiV1.PodExecOptions{
			Container: containerName,
			Command:   cmdArr,
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       false,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		zap.L().Error(fmt.Sprintf("CopyToPod NewSPDYExecutor err: %v", err))
		return err
	}
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  reader,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Tty:    false,
	})
	if err != nil {
		zap.L().Error(fmt.Sprintf("exec stream err: %v", err))
		return err
	}
	return nil
}

func makeTar(srcPath, destPath string, writer io.Writer) error {
	tarWriter := tar.NewWriter(writer)
	defer tarWriter.Close()

	srcPath = path.Clean(srcPath)
	destPath = path.Clean(destPath)
	return recursiveTar(path.Dir(srcPath), path.Base(srcPath), path.Dir(destPath), path.Base(destPath), tarWriter)
}

func recursiveTar(srcBase, srcFile, destBase, destFile string, tw *tar.Writer) error {
	srcPath := path.Join(srcBase, srcFile)
	matchedPaths, err := filepath.Glob(srcPath)
	if err != nil {
		return err
	}
	for _, fPath := range matchedPaths {
		stat, err := os.Lstat(fPath)
		if err != nil {
			return err
		}
		if stat.IsDir() {
			files, err := ioutil.ReadDir(fPath)
			if err != nil {
				return err
			}
			if len(files) == 0 {
				//case empty directory
				hdr, _ := tar.FileInfoHeader(stat, fPath)
				hdr.Name = destFile
				if err := tw.WriteHeader(hdr); err != nil {
					return err
				}
			}
			for _, f := range files {
				if err := recursiveTar(srcBase, path.Join(srcFile, f.Name()), destBase, path.Join(destFile, f.Name()), tw); err != nil {
					return err
				}
			}
			return nil
		} else if stat.Mode()&os.ModeSymlink != 0 {
			//case soft link
			hdr, _ := tar.FileInfoHeader(stat, fPath)
			target, err := os.Readlink(fPath)
			if err != nil {
				return err
			}

			hdr.Linkname = target
			hdr.Name = destFile
			if err := tw.WriteHeader(hdr); err != nil {
				return err
			}
		} else {
			//case regular file or other file type like pipe
			hdr, err := tar.FileInfoHeader(stat, fPath)
			if err != nil {
				return err
			}
			hdr.Name = destFile

			if err := tw.WriteHeader(hdr); err != nil {
				return err
			}

			f, err := os.Open(fPath)
			if err != nil {
				return err
			}
			defer f.Close()

			if _, err := io.Copy(tw, f); err != nil {
				return err
			}
			return f.Close()
		}
	}
	return nil
}

func checkDestinationIsDir(podName string, containerName string, destPath string) error {
	return exec(podName, containerName, []string{"test", "-d", destPath})
}

func checkDestinationIsFile(podName string, containerName string, destPath string) error {
	return exec(podName, containerName, []string{"test", "-e", destPath})
}

func createDir(podName string, containerName string, destPath string) error {
	return exec(podName, containerName, []string{"mkdir", "-p", destPath})
}

func exec(podName string, containerName string, cmd []string) error {
	req := client.CoreV1().RESTClient().
		Post().
		Namespace(apiV1.NamespaceDefault).
		Resource("pods").
		Name(podName).
		SubResource("exec").
		VersionedParams(&apiV1.PodExecOptions{
			Container: containerName,
			Command:   cmd,
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       false,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return err
	}
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  strings.NewReader(""),
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Tty:    false,
	})
	if err != nil {
		return err
	}
	return nil
}

func CopyFromPod(podName, containerName, srcPath string, destPath string) error {
	reader, outStream := io.Pipe()
	cmdArr := []string{"tar", "cf", "-", srcPath}
	req := client.CoreV1().RESTClient().
		Get().
		Namespace(apiV1.NamespaceDefault).
		Resource("pods").
		Name(podName).
		SubResource("exec").
		VersionedParams(&apiV1.PodExecOptions{
			Container: containerName,
			Command:   cmdArr,
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       false,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return err
	}
	go func() {
		defer outStream.Close()
		err = exec.Stream(remotecommand.StreamOptions{
			Stdin:  os.Stdin,
			Stdout: outStream,
			Stderr: os.Stderr,
			Tty:    false,
		})
	}()
	prefix := getPrefix(srcPath)
	prefix = path.Clean(prefix)
	destPath = path.Join(destPath, path.Base(prefix))
	err = unTarAll(reader, destPath, prefix)
	return err
}

func unTarAll(reader io.Reader, destDir, prefix string) error {
	tarReader := tar.NewReader(reader)
	for {
		header, err := tarReader.Next()
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}

		if !strings.HasPrefix(header.Name, prefix) {
			return fmt.Errorf("tar contents corrupted")
		}

		mode := header.FileInfo().Mode()
		destFileName := filepath.Join(destDir, header.Name[len(prefix):])

		baseName := filepath.Dir(destFileName)
		if err := os.MkdirAll(baseName, 0755); err != nil {
			return err
		}
		if header.FileInfo().IsDir() {
			if err := os.MkdirAll(destFileName, 0755); err != nil {
				return err
			}
			continue
		}

		path, err := filepath.EvalSymlinks(baseName)
		if err != nil {
			return err
		}

		if mode&os.ModeSymlink != 0 {
			linkname := header.Linkname

			if !filepath.IsAbs(linkname) {
				_ = filepath.Join(path, linkname)
			}

			if err := os.Symlink(linkname, destFileName); err != nil {
				return err
			}
		} else {
			outFile, err := os.Create(destFileName)
			if err != nil {
				return err
			}
			defer outFile.Close()
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return err
			}
			if err := outFile.Close(); err != nil {
				return err
			}
		}
	}

	return nil
}

func getPrefix(file string) string {
	return strings.TrimLeft(file, "/")
}

func int32Ptr(i int32) *int32 { return &i }
