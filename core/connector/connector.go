package connector

import (
	"fmt"
	"net"
	"os"
	"os/exec"

	"path"

	"time"

	"strings"

	humanize "github.com/dustin/go-humanize"
	"github.com/ellcrys/crypto"
	"github.com/ellcrys/util"
	cutil "github.com/ncodes/cocoon-util"
	"github.com/ncodes/cocoon/core/common"
	"github.com/ncodes/cocoon/core/config"
	"github.com/ncodes/cocoon/core/connector/monitor"
	docker "github.com/ncodes/go-dockerclient"
	logging "github.com/op/go-logging"
)

var log = logging.MustGetLogger("connector")
var buildLog = logging.MustGetLogger("ccode.build")
var runLog = logging.MustGetLogger("ccode.run")
var configLog = logging.MustGetLogger("ccode.config")
var ccodeLog = logging.MustGetLogger("ccode")
var dckClient *docker.Client

var bridgeName = os.Getenv("BRIDGE_NAME")

func init() {
	runLog.SetBackend(config.MessageOnlyBackend)
}

// Connector defines a structure for starting and managing a cocoon (coode)
type Connector struct {
	waitCh            chan bool
	req               *Request
	connectorRPCAddr  string
	cocoonCodeRPCAddr string
	languages         []Language
	container         *docker.Container
	containerRunning  bool
	monitor           *monitor.Monitor
	healthCheck       *HealthChecker
}

// NewConnector creates a new connector
func NewConnector(req *Request, waitCh chan bool) *Connector {
	return &Connector{
		req:     req,
		waitCh:  waitCh,
		monitor: monitor.NewMonitor(),
	}
}

// Launch starts a cocoon code
func (cn *Connector) Launch(connectorRPCAddr, cocoonCodeRPCAddr string) {

	endpoint := "unix:///var/run/docker.sock"
	client, err := docker.NewClient(endpoint)
	if err != nil {
		log.Errorf("failed to create docker client. Is dockerd running locally?. %s", err)
		cn.Stop(true)
		return
	}

	dckClient = client
	cn.monitor.SetDockerClient(dckClient)
	cn.healthCheck = NewHealthChecker("127.0.0.1"+cn.cocoonCodeRPCAddr, cn.cocoonUnresponsive)

	// No need downloading, building and starting a cocoon code
	// if DEV_COCOON_RPC_ADDR has been specified. This means a dev cocoon code
	// is running at that address. Just start the connector's client.
	if devCocoonCodeRPCAddr := os.Getenv("DEV_COCOON_RPC_ADDR"); len(devCocoonCodeRPCAddr) > 0 {
		cn.cocoonCodeRPCAddr = devCocoonCodeRPCAddr
		log.Infof("[Dev] Will interact with cocoon code at %s", devCocoonCodeRPCAddr)
		cn.healthCheck.Start()
		return
	}

	log.Info("Ready to install cocoon code")
	log.Debugf("Found ccode url=%s and lang=%s", cn.req.URL, cn.req.Lang)

	lang := cn.GetLanguage(cn.req.Lang)
	if lang == nil {
		log.Errorf("cocoon code language (%s) not supported", cn.req.Lang)
		cn.Stop(true)
		return
	}

	newContainer, err := cn.prepareContainer(cn.req, lang)
	if err != nil {
		log.Error(err.Error())
		cn.Stop(true)
		return
	}

	lang.SetRunEnv(map[string]string{
		"COCOON_ID":          cn.req.ID,
		"CONNECTOR_RPC_ADDR": cn.connectorRPCAddr,
		"COCOON_RPC_ADDR":    cn.cocoonCodeRPCAddr, // cocoon code server will bind to the port of this address
		"COCOON_LINK":        cn.req.Link,          // the cocoon code id to link to natively
	})

	go cn.monitor.Monitor()

	go func() {
		if err = cn.run(newContainer, lang); err != nil {
			log.Error(err.Error())
			cn.Stop(true)
			return
		}
	}()
}

// cocoonUnresponsive is called when the cocoon code failed health check
func (cn *Connector) cocoonUnresponsive() {
	log.Info("Cocoon code has failed health check. Stopping cocoon code.")
	cn.Stop(true)
}

// SetAddrs sets the address of the connector and cocoon code RPC servers
func (cn *Connector) SetAddrs(connectorRPCAddr, cocoonCodeRPCAddr string) {
	cn.connectorRPCAddr = connectorRPCAddr
	cn.cocoonCodeRPCAddr = cocoonCodeRPCAddr
}

// GetRequest returns the current cocoon launch request
func (cn *Connector) GetRequest() *Request {
	return cn.req
}

// GetCocoonCodeRPCAddr returns the RPC address of the cocoon code
func (cn *Connector) GetCocoonCodeRPCAddr() string {
	return cn.cocoonCodeRPCAddr
}

// prepareContainer fetches the cocoon code source, creates a container,
// moves the source in to it, builds the source within the container (if required)
// and configures default firewall.
func (cn *Connector) prepareContainer(req *Request, lang Language) (*docker.Container, error) {

	var containerID = util.Env("CONTAINER_ID", req.ID)

	_, err := cn.fetchSource(req, lang)
	if err != nil {
		return nil, err
	}

	// ensure cocoon code isn't already launched on a container
	c, err := cn.getContainer(containerID)
	if err != nil {
		return nil, fmt.Errorf("failed to check whether cocoon code is already active. %s ", err.Error())
	} else if c != nil {
		return nil, fmt.Errorf("cocoon code already exists on a container")
	}

	newContainer, err := cn.createContainer(req.ID, lang, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new container to run cocoon code. %s ", err.Error())
	}

	cn.container = newContainer
	cn.monitor.SetContainerID(cn.container.ID)
	cn.HookToMonitor(req)

	if lang.RequiresBuild() {
		var buildParams map[string]interface{}
		if len(req.BuildParams) > 0 {

			req.BuildParams, err = crypto.FromBase64(req.BuildParams)
			if err != nil {
				return nil, fmt.Errorf("failed to decode build parameter. Expects a base 64 encoded string. %s", err)
			}

			if err = util.FromJSON([]byte(req.BuildParams), &buildParams); err != nil {
				return nil, fmt.Errorf("failed to parse build parameter. Expects valid json string. %s", err)
			}
		}

		if err = lang.SetBuildParams(buildParams); err != nil {
			return nil, fmt.Errorf("failed to set and validate build parameter. %s", err)
		}

		err = cn.build(newContainer, lang)
		if err != nil {
			return nil, fmt.Errorf(err.Error())
		}

	} else {
		log.Info("Cocoon code does not require a build processing. Skipped.")
	}

	if err = cn.configFirewall(newContainer, req); err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	return newContainer, nil
}

// HookToMonitor is where all listeners to the monitor
// are attached.
func (cn *Connector) HookToMonitor(req *Request) {
	go func() {
		for evt := range cn.monitor.GetEmitter().On("monitor.report") {
			if cn.RestartIfDiskAllocExceeded(req, evt.Args[0].(monitor.Report).DiskUsage) {
				break
			}
		}
	}()
}

// RestartIfDiskAllocExceeded restarts the cocoon code is disk usages
// has exceeded its set limit.
func (cn *Connector) RestartIfDiskAllocExceeded(req *Request, curDiskSize int64) bool {
	if curDiskSize > req.DiskLimit {
		log.Errorf("cocoon code has used more than its allocated disk space (%s of %s)",
			humanize.Bytes(uint64(curDiskSize)),
			humanize.Bytes(uint64(req.DiskLimit)))
		if err := cn.restart(); err != nil {
			log.Error(err.Error())
			return false
		}
		return true
	}
	return false
}

// restart restarts the cocoon code. The running cocoon code is stopped
// and relaunched.
func (cn *Connector) restart() error {

	if dckClient == nil || cn.container == nil {
		return nil
	}

	if cn.monitor != nil {
		cn.monitor.Reset()
	}

	log.Info("Restarting cocoon code")

	cn.containerRunning = false

	err := dckClient.RemoveContainer(docker.RemoveContainerOptions{
		ID:            cn.container.ID,
		RemoveVolumes: true,
		Force:         true,
	})
	if err != nil {
		return fmt.Errorf("failed to remove container. %s", err)
	}

	newContainer, err := cn.prepareContainer(cn.req, cn.GetLanguage(cn.req.Lang))
	if err != nil {
		return fmt.Errorf("restart failed: %s", err)
	}

	go cn.monitor.Monitor()

	go func() {
		if err = cn.run(newContainer, cn.GetLanguage(cn.req.Lang)); err != nil {
			log.Errorf("restart failed: %s", err)
			cn.Stop(true)
		}
	}()

	return nil
}

// AddLanguage adds a new langauge to the launcher.
// Will return error if language is already added
func (cn *Connector) AddLanguage(lang Language) error {
	if cn.GetLanguage(lang.GetName()) != nil {
		return fmt.Errorf("language already exist")
	}
	cn.languages = append(cn.languages, lang)
	return nil
}

// GetLanguage will return a langauges or nil if not found
func (cn *Connector) GetLanguage(name string) Language {
	for _, l := range cn.languages {
		if l.GetName() == name {
			return l
		}
	}
	return nil
}

// GetLanguages returns all languages added to the launcher
func (cn *Connector) GetLanguages() []Language {
	return cn.languages
}

// fetchSource fetches the cocoon code source from
// a remote address
func (cn *Connector) fetchSource(req *Request, lang Language) (string, error) {

	if !cutil.IsGithubRepoURL(req.URL) {
		return "", fmt.Errorf("only public source code hosted on github is supported") // TODO: support zip files
	}

	return cn.fetchFromGit(req, lang)
}

// findLaunch looks for a previous stored launch/Redeployment by id
// TODO: needs implementation
func (cn *Connector) findLaunch(id string) interface{} {
	return nil
}

// fetchFromGit fetchs cocoon code from git repo.
// and returns the download directory.
func (cn *Connector) fetchFromGit(req *Request, lang Language) (string, error) {

	var repoTarURL, downloadDst string
	var err error

	// checks if job was previously deployed. find a job by the job name.
	if cn.findLaunch(req.ID) != nil {
		return "", fmt.Errorf("cocoon code was previously launched") // TODO: fetch last launch tag and use it
	}

	repoTarURL, err = cutil.GetGithubRepoRelease(req.URL, req.Tag)
	if err != nil {
		return "", fmt.Errorf("Failed to fetch release from github repo. %s", err)
	}

	// set tag to latest if not provided
	tagStr := req.Tag
	if tagStr == "" {
		tagStr = "latest"
	}

	// determine download directory
	downloadDst = lang.GetDownloadDestination()

	// delete download directory if it exists
	if _, err := os.Stat(downloadDst); err == nil {
		log.Info("Download destination is not empty. Deleting content")
		if err = os.RemoveAll(downloadDst); err != nil {
			return "", fmt.Errorf("failed to delete contents of download directory")
		}
		log.Info("Download directory has been deleted")
	}

	// create the download directory
	if err = os.MkdirAll(downloadDst, os.ModePerm); err != nil {
		return "", fmt.Errorf("Failed to create download directory. %s", err)
	}

	log.Infof("Downloading cocoon repository with tag=%s, dst=%s", tagStr, downloadDst)
	filePath := path.Join(downloadDst, fmt.Sprintf("%s.tar.gz", req.ID))
	err = cutil.DownloadFile(repoTarURL, filePath, func(buf []byte) {})
	if err != nil {
		return "", err
	}

	log.Info("Successfully downloaded cocoon code")
	log.Debugf("Unpacking cocoon code to %s", filePath)

	// unpack tarball
	cmd := "tar"
	args := []string{"-xf", filePath, "-C", downloadDst, "--strip-components", "1"}
	if err = exec.Command(cmd, args...).Run(); err != nil {
		return "", fmt.Errorf("Failed to unpack cocoon code repo tarball. %s", err)
	}

	log.Infof("Successfully unpacked cocoon code to %s", downloadDst)

	os.Remove(filePath)
	log.Info("Deleted the cocoon code tarball")

	return downloadDst, nil
}

// getContainer returns a container with a
// matching name or nil if not found.
func (cn *Connector) getContainer(name string) (*docker.APIContainers, error) {
	apiContainers, err := dckClient.ListContainers(docker.ListContainersOptions{All: true})
	if err != nil {
		return nil, err
	}

	for _, c := range apiContainers {
		if util.InStringSlice(c.Names, "/"+name) {
			return &c, nil
		}
	}

	return nil, nil
}

// createContainer creates a brand new container,
// and copies the cocoon source code to it.
func (cn *Connector) createContainer(name string, lang Language, env []string) (*docker.Container, error) {
	_, cocoonCodePort, _ := net.SplitHostPort(cn.cocoonCodeRPCAddr)

	container, err := dckClient.CreateContainer(docker.CreateContainerOptions{
		Name: name,
		Config: &docker.Config{
			Image: lang.GetImage(),
			Labels: map[string]string{
				"name": name,
				"type": "cocoon_code",
			},
			WorkingDir: lang.GetSourceRootDir(),
			Tty:        true,
			ExposedPorts: map[docker.Port]struct{}{
				docker.Port(fmt.Sprintf("%s/tcp", cocoonCodePort)): struct{}{},
			},
			Cmd:       []string{"bash"},
			Env:       env,
			Memory:    common.MBToByte(int64(cn.req.Memory)),
			CPUShares: cn.req.CPUShares,
		},
		HostConfig: &docker.HostConfig{
			Memory:    common.MBToByte(int64(cn.req.Memory)),
			CPUShares: cn.req.CPUShares,
			PortBindings: map[docker.Port][]docker.PortBinding{
				docker.Port(fmt.Sprintf("%s/tcp", cocoonCodePort)): []docker.PortBinding{
					docker.PortBinding{HostIP: "127.0.0.1", HostPort: cocoonCodePort},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	// no matter what happens, remove download directory
	defer func() {
		os.RemoveAll(lang.GetDownloadDestination())
		log.Info("Removed download directory")
	}()

	// copy source directory to the container's source directory
	cmd := "docker"
	args := []string{"cp", lang.GetDownloadDestination(), fmt.Sprintf("%s:%s", container.ID, lang.GetCopyDestination())}
	if err = exec.Command(cmd, args...).Run(); err != nil {
		return nil, fmt.Errorf("failed to copy cocoon code source to cocoon. %s", err)
	}

	log.Info("Copied cocoon code source to cocoon")

	return container, nil
}

// stopContainer stop container. Kill it if it doesn't
// end after 5 seconds.
func (cn *Connector) stopContainer(id string) error {
	if err := dckClient.StopContainer(id, uint((5 * time.Second).Seconds())); err != nil {
		return err
	}
	cn.containerRunning = false
	return nil
}

// Executes is a general purpose function
// to execute a command in a running container. If container is not running, it starts it.
// It accepts the container, a unique name for the execution
// and a callback function that is passed a lifecycle status and a value.
// If priviledged is set to true, command will attain root powers.
// Supported statuses are before (before command is executed), after (after command is executed)
// and end (when command exits).
func (cn *Connector) execInContainer(container *docker.Container, name string, command []string, priviledged bool, logger *logging.Logger, cb func(string, interface{}) error) error {

	containerStatus, err := dckClient.InspectContainer(container.ID)
	if err != nil {
		return fmt.Errorf("failed to inspect container before executing command [%s]. %s", name, err)
	}

	if !containerStatus.State.Running {
		err := dckClient.StartContainer(container.ID, nil)
		if err != nil {
			return fmt.Errorf("failed start container for exec [%s]. %s", name, err.Error())
		}
		cn.containerRunning = true
	}

	exec, err := dckClient.CreateExec(docker.CreateExecOptions{
		Container:    container.ID,
		AttachStderr: true,
		AttachStdout: true,
		Cmd:          command,
		Privileged:   priviledged,
	})

	if err != nil {
		return fmt.Errorf("failed to create exec [%s] object. %s", name, err)
	}

	if err = cb("before", nil); err != nil {
		return err
	}

	outStream := NewLogStreamer()
	outStream.SetLogger(logger)

	go func() {
		err = dckClient.StartExec(exec.ID, docker.StartExecOptions{
			OutputStream: outStream.GetWriter(),
			ErrorStream:  outStream.GetWriter(),
		})
		if err != nil {
			log.Infof("failed to start exec [%s] command. %s", name, err)
		}
	}()

	go func() {
		err := outStream.Start()
		if err != nil {
			log.Errorf("failed to start exec [%s] output stream logger. %s", name, err)
		}
	}()

	execExitCode := 0
	time.Sleep(1 * time.Second)

	if err = cb("after", nil); err != nil {
		outStream.Stop()
		return err
	}

	for cn.containerRunning {
		execIns, err := dckClient.InspectExec(exec.ID)
		if err != nil {
			outStream.Stop()
			return err
		}

		if execIns.Running {
			time.Sleep(500 * time.Millisecond)
			continue
		}

		execExitCode = execIns.ExitCode
		break
	}

	outStream.Stop()

	if err = cb("end", execExitCode); err != nil {
		return err
	}

	if execExitCode != 0 {
		return fmt.Errorf("Exec [%s] exited with code=%d", name, execExitCode)
	}

	return nil
}

// deleteBridge deletes the bridge the connector's docker daemon is attached to
func (cn *Connector) deleteBridge(bridgeName string) error {
	if bridgeName == "" {
		return nil
	}

	args := []string{"-c", `
		ip link set dev ` + bridgeName + ` down &&
		brctl delbr ` + bridgeName + `
	`}

	_, err := exec.Command("bash", args...).Output()
	if err != nil {
		return fmt.Errorf("failed to delete bridge. %s", err)
	}
	return nil
}

// build starts up the container and builds the cocoon code
// according to the build script provided by the languaged.
func (cn *Connector) build(container *docker.Container, lang Language) error {
	cmd := []string{"bash", "-c", lang.GetBuildScript()}
	return cn.execInContainer(container, "BUILD", cmd, false, buildLog, func(state string, val interface{}) error {
		switch state {
		case "before":
			log.Info("Building cocoon code...")
		case "end":
			if val.(int) == 0 {
				log.Info("Build succeeded!")
			} else {
				return fmt.Errorf("Build has failed with exit code=%d", val.(int))
			}
		}
		return nil
	})
}

// Run the cocoon code. First it gets the IP address of the container and sets
// the language environment.
func (cn *Connector) run(container *docker.Container, lang Language) error {
	return cn.execInContainer(container, "RUN", lang.GetRunScript(), false, runLog, func(state string, val interface{}) error {
		switch state {
		case "before":
			log.Info("Starting cocoon code")
		case "after":
			cn.healthCheck.Start()
			return nil
		case "end":
			if val.(int) == 0 {
				log.Info("Cocoon code successfully stop")
				return nil
			}
		}
		return nil
	})
}

// getDefaultFirewall returns the default firewall rules
// for a cocoon container.
func (cn *Connector) getDefaultFirewall() string {

	_, cocoonCodeRPCPort, _ := net.SplitHostPort(cn.cocoonCodeRPCAddr)
	connectorRPCIP, connectorRPCPort, _ := net.SplitHostPort(cn.connectorRPCAddr)

	return strings.TrimSpace(`iptables -F && 
			iptables -P INPUT DROP && 
			iptables -P FORWARD DROP &&
			iptables -P OUTPUT DROP &&
			iptables -A OUTPUT -m conntrack --ctstate RELATED,ESTABLISHED -j ACCEPT &&
			iptables -A OUTPUT -p tcp -d ` + connectorRPCIP + ` --dport ` + connectorRPCPort + ` -j ACCEPT
			iptables -A INPUT -m conntrack --ctstate RELATED,ESTABLISHED -j ACCEPT &&
			iptables -A INPUT -p tcp --dport ` + cocoonCodeRPCPort + ` -j ACCEPT 
			dnsIPs="$(cat /etc/resolv.conf | grep 'nameserver' | cut -c12-)" &&
			for ip in $dnsIPs;
			do 
				iptables -A OUTPUT -m state --state NEW,ESTABLISHED -d ${ip} -p udp --dport 53 -j ACCEPT;
				iptables -A OUTPUT -m state --state ESTABLISHED -p udp -s ${ip} --sport 53 -j ACCEPT;
				iptables -A OUTPUT -m state --state NEW,ESTABLISHED -d ${ip} -p tcp --dport 53 -j ACCEPT;
				iptables -A OUTPUT -m state --state ESTABLISHED -p tcp -s ${ip} --sport 53 -j ACCEPT;
			done`)
}

// configFirewall configures the container firewall.
func (cn *Connector) configFirewall(container *docker.Container, req *Request) error {
	cmd := []string{"bash", "-c", cn.getDefaultFirewall()}
	return cn.execInContainer(container, "CONFIG-FIREWALL", cmd, true, configLog, func(state string, val interface{}) error {
		switch state {
		case "before":
			log.Info("Configuring firewall for cocoon")
		case "end":
			if val.(int) == 0 {
				log.Info("Firewall configured for cocoon")
			}
		}
		return nil
	})
}

// Stop closes the client, stops the container if it is still running
// and deletes the container. This will effectively bring the launcher
// to a halt. Set failed parameter to true to set a positve exit code or
// false for 0 exit code.
func (cn *Connector) Stop(failed bool) error {

	defer func() {
		cn.waitCh <- failed
	}()

	if dckClient == nil || cn.container == nil {
		return nil
	}

	if cn.monitor != nil {
		cn.monitor.Stop()
	}

	if cn.healthCheck != nil {
		cn.healthCheck.Stop()
	}

	cn.containerRunning = false

	err := dckClient.RemoveContainer(docker.RemoveContainerOptions{
		ID:            cn.container.ID,
		RemoveVolumes: true,
		Force:         true,
	})
	if err != nil {
		return fmt.Errorf("failed to remove container. %s", err)
	}

	if err = cn.deleteBridge(bridgeName); err != nil {
		return fmt.Errorf("failed to delete bridge")
	}

	return nil
}