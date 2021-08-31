# Falco

[Falco](https://falco.org) is a *Cloud Native Runtime Security* tool designed to detect anomalous activity in your applications. You can use Falco to monitor runtime security of your Kubernetes applications and internal components.

## Introduction

This chart adds Falco to all nodes in your cluster using a DaemonSet.

Also provides a Deployment for generating Falco alerts. This is useful for testing purposes.

## Adding `falcosecurity` repository

Prior to install the chart, add the `falcosecurity` charts repository:

```bash
helm repo add falcosecurity https://falcosecurity.github.io/charts
helm repo update
```

## Installing the Chart

To install the chart with the release name `falco` run:

```bash
helm install falco falcosecurity/falco
```

After a few seconds, Falco should be running.

> **Tip**: List all releases using `helm list`, a release is a name used to track a specific deployment

## Uninstalling the Chart

To uninstall the `falco` deployment:

```bash
helm uninstall falco
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Configuration

The following table lists the configurable parameters of the Falco chart and their default values.

| Parameter                                       | Description                                                                                                        | Default                                                                                                                                   |
| ---                                             | ---                                                                                                                | ---                                                                                                                                       |
| `image.registry`                                | The image registry to pull from                                                                                    | `docker.io`                                                                                                                               |
| `image.repository`                              | The image repository to pull from                                                                                  | `falcosecurity/falco`                                                                                                                     |
| `image.tag`                                     | The image tag to pull                                                                                              | `0.27.0`                                                                                                                                  |
| `image.pullPolicy`                              | The image pull policy                                                                                              | `IfNotPresent`                                                                                                                            |
| `image.pullSecrets`                              | The image pull secretes                                                                                            | `[]`                                                                                                                                      |
| `containerd.enabled`                            | Enable ContainerD support                                                                                          | `true`                                                                                                                                    |
| `containerd.socket`                             | The path of the ContainerD socket                                                                                  | `/run/containerd/containerd.sock`                                                                                                         |
| `docker.enabled`                                | Enable Docker support                                                                                              | `true`                                                                                                                                    |
| `docker.socket`                                 | The path of the Docker daemon socket                                                                               | `/var/run/docker.sock`                                                                                                                    |
| `resources.requests.cpu`                        | CPU requested for being run in a node                                                                              | `100m`                                                                                                                                    |
| `resources.requests.memory`                     | Memory requested for being run in a node                                                                           | `512Mi`                                                                                                                                   |
| `resources.limits.cpu`                          | CPU limit                                                                                                          | `200m`                                                                                                                                    |
| `resources.limits.memory`                       | Memory limit                                                                                                       | `1024Mi`                                                                                                                                  |
| `extraArgs`                                     | Specify additional container args                                                                                  | `[]`                                                                                                                                      |
| `rbac.create`                                   | If true, create & use RBAC resources                                                                               | `true`                                                                                                                                    |
| `serviceAccount.create`                         | Create serviceAccount                                                                                              | `true`                                                                                                                                    |
| `serviceAccount.name`                           | Use this value as serviceAccountName                                                                               | ` `                                                                                                                                       |
| `fakeEventGenerator.enabled`                    | Run [falcosecurity/event-generator](https://github.com/falcosecurity/event-generator) for sample events            | `false`
| `fakeEventGenerator.args`                       | Arguments for `falcosecurity/event-generator`                                                                       | `run --loop ^syscall`                                                                                                                                   |
| `fakeEventGenerator.replicas`                   | How many replicas of `falcosecurity/event-generator` to run                                                                  | `1`                                                                                                                                       |
| `daemonset.updateStrategy.type`                 | The updateStrategy for updating the daemonset                                                                      | `RollingUpdate`                                                                                                                           |
| `daemonset.env`                                 | Extra environment variables passed to daemonset pods                                                               | `{}`                                                                                                                                      |
| `daemonset.podAnnotations`                      | Extra pod annotations to be added to pods created by the daemonset                                                 | `{}`                                                                                                                                      |
| `podSecurityPolicy.create`                      | If true, create & use podSecurityPolicy                                                                            | `false`                                                                                                                                   |
| `proxy.httpProxy`                               | Set the Proxy server if is behind a firewall                                                                       | ` `                                                                                                                                       |
| `proxy.httpsProxy`                              | Set the Proxy server if is behind a firewall                                                                       | ` `                                                                                                                                       |
| `proxy.noProxy`                                 | Set the Proxy server if is behind a firewall                                                                       | ` `                                                                                                                                       |
| `timezone`                                      | Set the daemonset's timezone                                                                                       | ` `                                                                                                                                       |
| `priorityClassName`                             | Set the daemonset's priorityClassName                                                                              | ` `                                                                                                                                       |
| `ebpf.enabled`                                  | Enable eBPF support for Falco instead of `falco-probe` kernel module                                               | `false`                                                                                                                                   |
| `ebpf.settings.hostNetwork`                     | Needed to enable eBPF JIT at runtime for performance reasons                                                       | `true`                                                                                                                                    |
| `auditLog.enabled`                              | Enable K8s audit log support for Falco                                                                             | `false`                                                                                                                                   |
| `auditLog.dynamicBackend.enabled`               | Deploy the Audit Sink where Falco listens for K8s audit log events                                                 | `false`                                                                                                                                   |
| `auditLog.dynamicBackend.url`                   | Define if Audit Sink client config should point to a fixed [url](https://kubernetes.io/docs/tasks/debug-application-cluster/audit/#url) (useful for development) instead of the default webserver service.                                                 | ``                                                                                                                                   |
| `falco.rulesFile`                               | The location of the rules files                                                                                    | `[/etc/falco/falco_rules.yaml, /etc/falco/falco_rules.local.yaml, /etc/falco/rules.available/application_rules.yaml, /etc/falco/rules.d]` |
| `falco.timeFormatISO8601`                       | Display times using ISO 8601 instead of local time zone                                                            | `false`                                                                                                                                   |
| `falco.jsonOutput`                              | Output events in json or text                                                                                      | `false`                                                                                                                                   |
| `falco.jsonIncludeOutputProperty`               | Include output property in json output                                                                             | `true`                                                                                                                                    |
| `falco.logStderr`                               | Send Falco debugging information logs to stderr                                                                    | `true`                                                                                                                                    |
| `falco.logSyslog`                               | Send Falco debugging information logs to syslog                                                                    | `true`                                                                                                                                    |
| `falco.logLevel`                                | The minimum level of Falco debugging information to include in logs                                                | `info`                                                                                                                                    |
| `falco.priority`                                | The minimum rule priority level to load and run                                                                    | `debug`                                                                                                                                   |
| `falco.bufferedOutputs`                         | Use buffered outputs to channels                                                                                   | `false`                                                                                                                                   |
| `falco.syscallEventDrops.actions`               | Actions to be taken when system calls were dropped from the circular buffer                                        | `[log, alert]`                                                                                                                            |
| `falco.syscallEventDrops.rate`                  | Rate at which log/alert messages are emitted                                                                       | `.03333`                                                                                                                                  |
| `falco.syscallEventDrops.maxBurst`              | Max burst of messages emitted                                                                                      | `10`                                                                                                                                      |
| `falco.outputs.output_timeout`                  | Duration in milliseconds to wait before considering the output timeout deadline exceed                             | `2000`                                                                                                                                    |
| `falco.outputs.rate`                            | Number of tokens gained per second                                                                                 | `1`                                                                                                                                       |
| `falco.outputs.maxBurst`                        | Maximum number of tokens outstanding                                                                               | `1000`                                                                                                                                    |
| `falco.syslogOutput.enabled`                    | Enable syslog output for security notifications                                                                    | `true`                                                                                                                                    |
| `falco.fileOutput.enabled`                      | Enable file output for security notifications                                                                      | `false`                                                                                                                                   |
| `falco.fileOutput.keepAlive`                    | Open file once or every time a new notification arrives                                                            | `false`                                                                                                                                   |
| `falco.fileOutput.filename`                     | The filename for logging notifications                                                                             | `./events.txt`                                                                                                                            |
| `falco.stdoutOutput.enabled`                    | Enable stdout output for security notifications                                                                    | `true`                                                                                                                                    |
| `falco.webserver.enabled`                       | Enable Falco embedded webserver to accept K8s audit events                                                         | `true`                                                                                                                                    |
| `falco.webserver.k8sAuditEndpoint`              | Endpoint where Falco embedded webserver accepts K8s audit events                                                   | `/k8s-audit`                                                                                                                              |
| `falco.webserver.listenPort`                    | Port where Falco embedded webserver listen to connections                                                          | `8765`                                                                                                                                    |
| `falco.webserver.nodePort`                      | Exposes the Falco embedded webserver through a NodePort                                                            | `false`                                                                                                                                   |
| `falco.webserver.sslEnabled`                    | Enable SSL on Falco embedded webserver                                                                             | `false`                                                                                                                                   |
| `falco.webserver.sslCertificate`                | Certificate bundle path for the Falco embedded webserver                                                           | `/etc/falco/certs/server.pem`                                                                                                             |
| `falco.programOutput.enabled`                   | Enable program output for security notifications                                                                   | `false`                                                                                                                                   |
| `falco.programOutput.keepAlive`                 | Start the program once or re-spawn when a notification arrives                                                     | `false`                                                                                                                                   |
| `falco.programOutput.program`                   | Command to execute for program output                                                                              | `mail -s "Falco Notification" someone@example.com`                                                                                        |
| `falco.httpOutput.enabled`                      | Enable http output for security notifications                                                                      | `false`                                                                                                                                   |
| `falco.httpOutput.url`                          | Url to notify using the http output when a notification arrives                                                    | `http://some.url`                                                                                                                         |
| `falco.grpc.enabled`                            | Enable the Falco gRPC server                                                                                       | `false`                                                                                                                                   |
| `falco.grpc.threadiness`                        | Number of threads (and context) the gRPC server will use, `0` by default, which means "auto"                                                           | `0`                                                                                                                                       |
| `falco.grpc.unixSocketPath`                     | Unix socket the gRPC server will create                                                                            | `unix:///var/run/falco/falco.sock`                                                                                                        |
| `falco.grpc.listenPort`                         | Port where Falco gRPC server listen to connections                                                                 | `5060`                                                                                                                                    |
| `falco.grpc.privateKey`                         | Key file path for the Falco gRPC server                                                                            | `/etc/falco/certs/server.key`                                                                                                             |
| `falco.grpc.certChain`                          | Cert file path for the Falco gRPC server                                                                           | `/etc/falco/certs/server.crt`                                                                                                             |
| `falco.grpc.rootCerts`                          | CA root file path for the Falco gRPC server                                                                        | `/etc/falco/certs/ca.crt`                                                                                                                 |
| `falco.grpcOutput.enabled`                      | Enable the gRPC output and events will be kept in memory until you read them with a gRPC client.                   | `false`                                                                                                                                   |
| `customRules`                                   | Third party rules enabled for Falco                                                                                | `{}`                                                                                                                                      |
| `certs.server.key`                              | Key used by gRPC and webserver                                                                                     | ` `                                                                                                                                       |
| `certs.server.crt`                              | Certificate used by gRPC and webserver                                                                             | ` `                                                                                                                                       |
| `certs.ca.crt`                                  | CA certificate used by gRPC, webserver and AuditSink validation                                                    | ` `                                                                                                                                       |
| `nodeSelector`                                  | The node selection constraint                                                                                      | `{}`                                                                                                                                      |
| `affinity`                                      | The affinity constraint                                                                                            | `{}`                                                                                                                                      |
| `tolerations`                                   | The tolerations for scheduling                                                                                     | `node-role.kubernetes.io/master:NoSchedule`                                                                                               |
| `scc.create`                                    | Create OpenShift's Security Context Constraint                                                                     | `true`                                                                                                                                    |
| `extraInitContainers`                           | A list of initContainers you want to add to the falco pod in the daemonset.                                        | `[]`                                                                                                                                      |
| `extraVolumes`                                  | A list of volumes you want to add to the falco daemonset.                                                          | `[]`                                                                                                                                      |
| `extraVolumeMounts`                             | A list of volumeMounts you want to add to the falco container in the falco daemonset.                              | `[]`                                                                                                                                      |
| `falcosidekick.enable`                             | Enable `falcosidekick` deployment                              | `false`                                                                                                                                      |

Specify each parameter using the `--set key=value[,key=value]` argument to `helm install`. For example,

```bash
helm install falco --set falco.jsonOutput=true falcosecurity/falco
```

Alternatively, a YAML file that specifies the values for the parameters can be provided while installing the chart. For example,

```bash
helm install falco -f values.yaml falcosecurity/falco
```

> **Tip**: You can use the default [values.yaml](values.yaml)

## Loading custom rules

Falco ships with a nice default ruleset. Is a good starting point but sooner or later we are going to need to add custom rules which fits our needs.

So the question is: How we can load custom rules in our Falco deployment?

We are going to create a file which contains custom rules so that we can keep it in a Git repository.

```bash
cat custom-rules.yaml
```

And the file looks like this one:

```yaml
customRules:
  rules-traefik.yaml: |-
    - macro: traefik_consider_syscalls
      condition: (evt.num < 0)

    - macro: app_traefik
      condition: container and container.image startswith "traefik"

    # Restricting listening ports to selected set

    - list: traefik_allowed_inbound_ports_tcp
      items: [443, 80, 8080]

    - rule: Unexpected inbound tcp connection traefik
      desc: Detect inbound traffic to traefik using tcp on a port outside of expected set
      condition: inbound and evt.rawres >= 0 and not fd.sport in (traefik_allowed_inbound_ports_tcp) and app_traefik
      output: Inbound network connection to traefik on unexpected port (command=%proc.cmdline pid=%proc.pid connection=%fd.name sport=%fd.sport user=%user.name %container.info image=%container.image)
      priority: NOTICE

    # Restricting spawned processes to selected set

    - list: traefik_allowed_processes
      items: ["traefik"]

    - rule: Unexpected spawned process traefik
      desc: Detect a process started in a traefik container outside of an expected set
      condition: spawned_process and not proc.name in (traefik_allowed_processes) and app_traefik
      output: Unexpected process spawned in traefik container (command=%proc.cmdline pid=%proc.pid user=%user.name %container.info image=%container.image)
      priority: NOTICE
```

So next step is to use the custom-rules.yaml file for installing the Falco Helm chart.

```bash
helm install falco -f custom-rules.yaml falcosecurity/falco
```

And we will see in our logs something like:

```bash
Tue Jun  5 15:08:57 2018: Loading rules from file /etc/falco/rules.d/rules-traefik.yaml:
```

And this means that our Falco installation has loaded the rules and is ready to help us.


## Enabling K8s audit event support

### Using scripts
This has been tested with Kops and Minikube. You will need the following components:

* A Kubernetes cluster greater than v1.13
* The apiserver must be configured with Dynamic Auditing feature, do it with the following flags:
  * `--audit-dynamic-configuration`
  * `--feature-gates=DynamicAuditing=true`
  * `--runtime-config=auditregistration.k8s.io/v1alpha1=true`

You can do it with the [scripts provided by Falco engineers](https://github.com/falcosecurity/evolution/tree/master/examples/k8s_audit_config)
just running:

```
cd examples/k8s_audit_config
bash enable-k8s-audit.sh minikube dynamic
```

Or in the case of Kops:

```
cd examples/k8s_audit_config
APISERVER_HOST=api.my-kops-cluster.com bash ./enable-k8s-audit.sh kops dynamic
```

Then you can install Falco chart enabling the enabling the `falco.webserver`
flag:

`helm install falco --set auditLog.enabled=true --set auditLog.dynamicBackend.enabled=true falcosecurity/falco`

And that's it, you will start to see the K8s audit log related alerts.

### Known validation failed error

Perhaps you may find the case where you receive an error like the following one:

```
helm install falco --set auditLog.enabled=true falcosecurity/falco
Error: validation failed: unable to recognize "": no matches for kind "AuditSink" in version "auditregistration.k8s.io/v1alpha1"
```

This means that the apiserver cannot recognize the `auditregistration.k8s.io`
resource, which means that the dynamic auditing feature hasn't been enabled
properly. You need to enable it or ensure that your using a Kubernetes version
greater than v1.13.

### Manual setup with NodePort on kOps

Using `kops edit cluster`, ensure these options are present, then run `kops update cluster` and `kops rolling-update cluster`:
```yaml
spec:
  kubeAPIServer:
    auditLogMaxBackups: 1
    auditLogMaxSize: 10
    auditLogPath: /var/log/k8s-audit.log
    auditPolicyFile: /srv/kubernetes/assets/audit-policy.yaml
    auditWebhookBatchMaxWait: 5s
    auditWebhookConfigFile: /srv/kubernetes/assets/webhook-config.yaml
  fileAssets:
  - content: |
      # content of the webserver CA certificate
      # remove this fileAsset and certificate-authority from webhook-config if using http
    name: audit-ca.pem
    roles:
    - Master
  - content: |
      apiVersion: v1
      kind: Config
      clusters:
      - name: falco
        cluster:
          # remove 'certificate-authority' when using 'http'
          certificate-authority: /srv/kubernetes/assets/audit-ca.pem
          server: https://localhost:32765/k8s-audit
      contexts:
      - context:
          cluster: falco
          user: ""
        name: default-context
      current-context: default-context
      preferences: {}
      users: []
    name: webhook-config.yaml
    roles:
    - Master
  - content: |
      # ... paste audit-policy.yaml here ...
      # https://raw.githubusercontent.com/falcosecurity/evolution/master/examples/k8s_audit_config/audit-policy.yaml
    name: audit-policy.yaml
    roles:
    - Master
```

Then you can install the Falco chart enabling these flags:

```shell
# without SSL (not recommended):
helm install falco --set auditLog.enabled=true --set falco.webserver.nodePort=32765 falcosecurity/falco

# with SSL:
helm install falco \
  --set auditLog.enabled=true \
  --set falco.webserver.sslEnabled=true \
  --set falco.webserver.nodePort=32765 \
  --set-file certs.server.key=/path/to/server.key \
  --set-file certs.server.crt=/path/to/server.crt \
  --set-file certs.ca.crt=/path/to/ca.crt \
  falcosecurity/falco
```

The webserver reuses the gRPC certificate setup, which is [documented here](https://falco.org/docs/grpc/#generate-valid-ca). Generating the client certificate isn't required.


## Enabling gRPC

The Falco gRPC server and the Falco gRPC Outputs APIs are not enabled by default.
Morover, Falco supports running a gRPC server with two main binding types:
- Over a local **unix socket** with no authentication 
- Over the **network** with mandatory mutual TLS authentication (mTLS)

> **Tip**: Once gRPC is enabled, you can deploy [falco-exporter](https://github.com/falcosecurity/falco-exporter) to export metrics to Prometheus.

### gRPC over unix socket (default)

The preferred way to use the gRPC is over a unix socket.

To install Falco with gRPC enabled over a **unix socket**, you have to:

```shell
helm install falco \
  --set falco.grpc.enabled=true \
  --set falco.grpcOutput.enabled=true \
  falcosecurity/falco
```

### gRPC over network

The gRPC server over the network can only be used with mutual authentication between the clients and the server using TLS certificates. 
How to generate the certificates is [documented here](https://falco.org/docs/grpc/#generate-valid-ca).

To install Falco with gRPC enabled over the **network**, you have to:

```shell
helm install falco \
  --set falco.grpc.enabled=true \
  --set falco.grpcOutput.enabled=true \
  --set falco.grpc.unixSocketPath="" \
  --set-file certs.server.key=/path/to/server.key \
  --set-file certs.server.crt=/path/to/server.crt \
  --set-file certs.ca.crt=/path/to/ca.crt \
  falcosecurity/falco
```
