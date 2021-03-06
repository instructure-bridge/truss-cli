package truss

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

// KubectlCmd wrapper for kubectl
type KubectlCmd struct {
	kubeconfig     string
	portForwardCmd *exec.Cmd
}

// Kubectl wrapper for kubectl
func Kubectl(kubeconfig string) *KubectlCmd {
	return &KubectlCmd{
		kubeconfig: kubeconfig,
	}
}

// PortForward kubectl port-forward
func (kubectl *KubectlCmd) PortForward(port, listen, namespace, target string, timeoutSeconds int) error {
	log.Debugln("Opening connection port forward for", port)
	argsWithCmd := []string{"port-forward", "-n=" + namespace, target, listen + ":" + port}
	kubectl.portForwardCmd = exec.Command("kubectl", argsWithCmd...)

	if kubectl.kubeconfig != "" {
		_, err := os.Stat(kubectl.kubeconfig)
		if err != nil {
			return err
		}
		kubectl.portForwardCmd.Env = append(os.Environ(), "KUBECONFIG="+kubectl.kubeconfig)
	}

	if err := kubectl.portForwardCmd.Start(); err != nil {
		return fmt.Errorf("Failed to port forward: %v", string(err.(*exec.ExitError).Stderr))
	}

	return waitForPort(listen, timeoutSeconds)
}

// ClosePortForward sigterm kubectl port-forward
func (kubectl *KubectlCmd) ClosePortForward() error {
	return kubectl.portForwardCmd.Process.Signal(syscall.SIGTERM)
}

// Run kubectl
func (kubectl *KubectlCmd) Run(arg ...string) ([]byte, error) {
	cmd := exec.Command("kubectl", arg...)

	if kubectl.kubeconfig != "" {
		cmd.Env = append(os.Environ(), "KUBECONFIG="+kubectl.kubeconfig)
	}

	bytes, err := cmd.Output()
	if err != nil {
		return nil, errors.New(string(err.(*exec.ExitError).Stderr))
	}

	return bytes, nil
}

func waitForPort(port string, timeoutSeconds int) error {
	log.Debugln("Waiting for port", port)
	for i := 0; i < timeoutSeconds; i++ {
		conn, err := net.Dial("tcp", ":"+port)
		if conn != nil {
			defer conn.Close()
		}
		if err == nil {
			return nil
		}
		time.Sleep(time.Second)
	}

	return fmt.Errorf("Could not reach port: %v", port)
}
