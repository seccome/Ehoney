# Falcosidekick

![falcosidekick](https://github.com/falcosecurity/falcosidekick/raw/master/imgs/falcosidekick_color.png)

![release](https://flat.badgen.net/github/release/falcosecurity/falcosidekick/latest?color=green) ![last commit](https://flat.badgen.net/github/last-commit/falcosecurity/falcosidekick) ![licence](https://flat.badgen.net/badge/license/MIT/blue) ![docker pulls](https://flat.badgen.net/docker/pulls/falcosecurity/falcosidekick?icon=docker)

## Description

A simple daemon for enhancing available outputs for [Falco](https://sysdig.com/opensource/falco/). It takes a falco's event and forwards it to different outputs.

It works as a single endpoint for as many as you want `falco` instances :

![falco_with_falcosidekick](https://github.com/falcosecurity/falcosidekick/raw/master/imgs/falco_with_falcosidekick.png)

## Outputs

Currently available outputs are :

* [**Slack**](https://slack.com)
* [**Rocketchat**](https://rocket.chat/)
* [**Mattermost**](https://mattermost.com/)
* [**Teams**](https://products.office.com/en-us/microsoft-teams/group-chat-software)
* [**Datadog**](https://www.datadoghq.com/)
* [**Discord**](https://www.discord.com/)
* [**AlertManager**](https://prometheus.io/docs/alerting/alertmanager/)
* [**Elasticsearch**](https://www.elastic.co/)
* [**Loki**](https://grafana.com/oss/loki)
* [**NATS**](https://nats.io/)
* [**STAN (NATS Streaming)**](https://docs.nats.io/nats-streaming-concepts/intro)
* [**Influxdb**](https://www.influxdata.com/products/influxdb-overview/)
* [**AWS Lambda**](https://aws.amazon.com/lambda/features/)
* [**AWS SQS**](https://aws.amazon.com/sqs/features/)
* [**AWS SNS**](https://aws.amazon.com/sns/features/)
* [**AWS CloudWatchLogs**](https://aws.amazon.com/cloudwatch/features/)
* **SMTP** (email)
* [**Opsgenie**](https://www.opsgenie.com/)
* [**StatsD**](https://github.com/statsd/statsd) (for monitoring of `falcosidekick`)
* [**DogStatsD**](https://docs.datadoghq.com/developers/dogstatsd/?tab=go) (for monitoring of `falcosidekick`)
* **Webhook**
* [**Azure Event Hubs**](https://azure.microsoft.com/en-in/services/event-hubs/)
* [**Prometheus**](https://prometheus.io/) (for both events and monitoring of `falcosidekick`)
* [**GCP PubSub**](https://cloud.google.com/pubsub)
* [**Google Chat**](https://workspace.google.com/products/chat/)
* [**Apache Kafka**](https://kafka.apache.org/)
* [**PagerDuty**](https://pagerduty.com/)
* [**Kubeless**](https://https://kubeless.io/)

## Adding `falcosecurity` repository

Prior to install the chart, add the `falcosecurity` charts repository:

```bash
helm repo add falcosecurity https://falcosecurity.github.io/charts
helm repo update
```

## Installing the Chart

To install the chart with the release name `falcosidekick` run:

```bash
helm install falcosidekick falcosecurity/falcosidekick
```

After a few seconds, Falcosidekick should be running.

> **Tip**: List all releases using `helm list`, a release is a name used to track a specific deployment

## Minumiun Kubernetes version

The minimum Kubernetes version required is 1.17.x

## Uninstalling the Chart

To uninstall the `falcosidekick` deployment:

```bash
helm uninstall falcosidekick
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Configuration

The following table lists the configurable parameters of the Falcosidekick chart and their default values.

| Parameter                                   | Description                                                                                                                                                                            | Default                                                                                           |
| ------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------- |
| `replicaCount`                              | number of running pods                                                                                                                                                                 | `1`                                                                                               |
| `listenport`                                | port to listen for daemon                                                                                                                                                              | `2801`                                                                                            |
| `config.debug`                              | if *true* all outputs will print in stdout the payload they send                                                                                                                       | `false`                                                                                           |
| `config.customfields`                       | a list of comma separated custom fields to add to falco events, syntax is "key:value,key:value"                                                                                        |                                                                                                   |
| `config.checkcert`                          | check if ssl certificate of the output is valid                                                                                                                                        | `true`                                                                                            |
| `config.slack.webhookurl`                   | Slack Webhook URL (ex: https://hooks.slack.com/services/XXXX/YYYY/ZZZZ), if not `empty`, Slack output is *enabled*                                                                     |                                                                                                   |
| `config.slack.footer`                       | Slack Footer                                                                                                                                                                           | https://github.com/falcosecurity/falcosidekick                                                    |
| `config.slack.icon`                         | Slack icon (avatar)                                                                                                                                                                    | https://raw.githubusercontent.com/falcosecurity/falcosidekick/master/imgs/falcosidekick_color.png |
| `config.slack.username`                     | Slack username                                                                                                                                                                         | `falcosidekick`                                                                                   |
| `config.slack.outputformat`                 | `all` (default), `text` (only text is displayed in Slack), `fields` (only fields are displayed in Slack)                                                                               | `all`                                                                                             |
| `config.slack.minimumpriority`              | minimum priority of event for using use this output, order is `emergency|alert|critical|error|warning|notice|informational|debug or ""`                                                | `debug`                                                                                           |
| `config.slack.messageformat`                | a Go template to format Slack Text above Attachment, displayed in addition to the output from `slack.outputformat`. If empty, no Text is displayed before Attachment                   |                                                                                                   |
| `config.rocketchat.webhookurl`              | Rocketchat Webhook URL (ex: https://XXXX/hooks/YYYY), if not `empty`, Rocketchat output is *enabled*                                                                                   |                                                                                                   |
| `config.rocketchat.icon`                    | Rocketchat icon (avatar)                                                                                                                                                               | https://raw.githubusercontent.com/falcosecurity/falcosidekick/master/imgs/falcosidekick_color.png |
| `config.rocketchat.username`                | Rocketchat username                                                                                                                                                                    | `falcosidekick`                                                                                   |
| `config.rocketchat.outputformat`            | `all` (default), `text` (only text is displayed in Rocketcaht), `fields` (only fields are displayed in Rocketchat)                                                                     | `all`                                                                                             |
| `config.rocketchat.minimumpriority`         | minimum priority of event for using use this output, order is `emergency|alert|critical|error|warning|notice|informational|debug or ""`                                                | `debug`                                                                                           |
| `config.rockerchat.messageformat`           | a Go template to format Rocketchat Text above Attachment, displayed in addition to the output from `slack.outputformat`. If empty, no Text is displayed before Attachment              |                                                                                                   |
| `config.mattermost.webhookurl`              | Mattermost Webhook URL (ex: https://XXXX/hooks/YYYY), if not `empty`, Mattermost output is *enabled*                                                                                   |                                                                                                   |
| `config.mattermost.footer`                  | Mattermost Footer                                                                                                                                                                      | https://github.com/falcosecurity/falcosidekick                                                    |
| `config.mattermost.icon`                    | Mattermost icon (avatar)                                                                                                                                                               | https://raw.githubusercontent.com/falcosecurity/falcosidekick/master/imgs/falcosidekick_color.png |
| `config.mattermost.username`                | Mattermost username                                                                                                                                                                    | `falcosidekick`                                                                                   |
| `config.mattermost.outputformat`            | `all` (default), `text` (only text is displayed in Slack), `fields` (only fields are displayed in Mattermost)                                                                          | `all`                                                                                             |
| `config.mattermost.minimumpriority`         | minimum priority of event for using use this output, order is `emergency|alert|critical|error|warning|notice|informational|debug or ""`                                                | `debug`                                                                                           |
| `config.mattermost.messageformat`           | a Go template to format Mattermost Text above Attachment, displayed in addition to the output from `slack.outputformat`. If empty, no Text is displayed before Attachment              |                                                                                                   |
| `config.teams.webhookurl`                   | Teams Webhook URL (ex: https://outlook.office.com/webhook/XXXXXX/IncomingWebhook/YYYYYY"), if not `empty`, Teams output is *enabled*                                                   |                                                                                                   |
| `config.teams.activityimage`                | Teams section image                                                                                                                                                                    | https://raw.githubusercontent.com/falcosecurity/falcosidekick/master/imgs/falcosidekick_color.png |
| `config.teams.outputformat`                 | `all` (default), `text` (only text is displayed in Teams), `facts` (only facts are displayed in Teams)                                                                                 | `all`                                                                                             |
| `config.teams.minimumpriority`              | minimum priority of event for using use this output, order is `emergency|alert|critical|error|warning|notice|informational|debug or ""`                                                | `debug`                                                                                           |
| `config.datadog.apikey`                     | Datadog API Key, if not `empty`, Datadog output is *enabled*                                                                                                                           |                                                                                                   |
| `config.datadog.host`                       | Datadog host. Override if you are on the Datadog EU site. Defaults to american site with "https://api.datadoghq.com"                                                                   | https://api.datadoghq.com                                                                         |
| `config.datadog.minimumpriority`            | minimum priority of event for using use this output, order is `emergency|alert|critical|error|warning|notice|informational|debug or ""`                                                | `debug`                                                                                           |
| `config.discord.webhookurl`                 | Discord WebhookURL (ex: https://discord.com/api/webhooks/xxxxxxxxxx...), if not empty, Discord output is enabled                                                                       |                                                                                                   |
| `config.discord.icon`                       | Discord icon (avatar)                                                                                                                                                                  | https://raw.githubusercontent.com/falcosecurity/falcosidekick/master/imgs/falcosidekick_color.png |
| `config.discord.minimumpriority`            | minimum priority of event for using use this output, order is `emergency|alert|critical|error|warning|notice|informational|debug or ""`                                                | `debug`                                                                                           |
| `config.alertmanager.hostport`              | AlertManager http://host:port, if not `empty`, AlertManager is *enabled*                                                                                                               |                                                                                                   |
| `config.alertmanager.minimumpriority`       | minimum priority of event for using use this output, order is `emergency|alert|critical|error|warning|notice|informational|debug or ""`                                                | `debug`                                                                                           |
| `config.elasticsearch.hostport`             | Elasticsearch http://host:port, if not `empty`, Elasticsearch is *enabled*                                                                                                             |                                                                                                   |
| `config.elasticsearch.index`                | Elasticsearch index                                                                                                                                                                    | `falco`                                                                                           |
| `config.elasticsearch.type`                 | Elasticsearch document type                                                                                                                                                            | `event`                                                                                           |
| `config.elasticsearch.suffix`               | date suffix for index rotation : `daily`, `monthly`, `annually`, `none`                                                                                                                | `daily`                                                                                           |
| `config.elasticsearch.minimumpriority`      | minimum priority of event for using use this output, order is `emergency|alert|critical|error|warning|notice|informational|debug or ""`                                                | `debug`                                                                                           |
| `config.influxdb.hostport`                  | Influxdb http://host:port, if not `empty`, Influxdb is *enabled*                                                                                                                       |                                                                                                   |
| `config.influxdb.database`                  | Influxdb database                                                                                                                                                                      | `falco`                                                                                           |
| `config.influxdb.user`                      | User to use if auth is enabled in Influxdb                                                                                                                                             |                                                                                                   |
| `config.influxdb.password`                  | Password to use if auth is enabled in Influxdb                                                                                                                                         |                                                                                                   |
| `config.influxdb.minimumpriority`           | minimum priority of event for using use this output, order is `emergency|alert|critical|error|warning|notice|informational|debug or ""`                                                | `debug`                                                                                           |
| `config.loki.hostport`                      | Loki http://host:port, if not `empty`, Loki is *enabled*                                                                                                                               |                                                                                                   |
| `config.loki.minimumpriority`               | minimum priority of event for using use this output, order is `emergency|alert|critical|error|warning|notice|informational|debug or ""`                                                | `debug`                                                                                           |
| `config.nats.hostport`                      | NATS "nats://host:port", if not `empty`, NATS is *enabled*                                                                                                                             |                                                                                                   |
| `config.nats.minimumpriority`               | minimum priority of event for using use this output, order is `emergency|alert|critical|error|warning|notice|informational|debug or ""`                                                | `debug`                                                                                           |
| `config.stan.hostport`                      | Stan nats://{domain or ip}:{port}, if not empty, STAN output is *enabled*                                                                                                              |                                                                                                   |
| `config.stan.clusterid`                     | Cluster name, if not empty, STAN output is *enabled*                                                                                                                                   | `debug`                                                                                           |
| `config.stan.clientid`                      | Client ID, if not empty, STAN output is *enabled*                                                                                                                                      |                                                                                                   |
| `config.stan.minimumpriority`               | minimum priority of event for using use this output, order is `emergency|alert|critical|error|warning|notice|informational|debug or ""`                                                | `debug`                                                                                           |
| `config.aws.accesskeyid`                    | AWS Access Key Id (optionnal if you use EC2 Instance Profile)                                                                                                                          |                                                                                                   |
| `config.aws.secretaccesskey`                | AWS Secret Access Key (optionnal if you use EC2 Instance Profile)                                                                                                                      |                                                                                                   |
| `config.aws.region`                         | AWS Region (optionnal if you use EC2 Instance Profile)                                                                                                                                 |                                                                                                   |
| `config.aws.cloudwatchlogs.loggroup`        | AWS CloudWatch Logs Group name, if not empty, CloudWatch Logs output is enabled                                                                                                        |                                                                                                   |
| `config.aws.cloudwatchlogs.logstream`       | AWS CloudWatch Logs Stream name, if empty, Falcosidekick will try to create a log stream                                                                                               | `debug`                                                                                           |
| `config.aws.cloudwatchlogs.minimumpriority` | minimum priority of event for using use this output, order is `emergency|alert|critical|error|warning|notice|informational|debug or ""`                                                | `debug`                                                                                           |
| `config.aws.lambda.functionname`            | AWS Lambda Function Name, if not empty, AWS Lambda output is enabled                                                                                                                   |                                                                                                   |
| `config.aws.lambda.minimumpriority`         | minimum priority of event for using use this output, order is `emergency|alert|critical|error|warning|notice|informational|debug or ""`                                                | `debug`                                                                                           |
| `config.aws.sns.topicarn`                   | AWS SNS TopicARN, if not empty, AWS SNS output is enabled                                                                                                                              |                                                                                                   |
| `config.aws.sns.rawjson`                    | Send RawJSON from `falco` or parse it to AWS SNS                                                                                                                                       |                                                                                                   |
| `config.aws.sns.minimumpriority`            | minimum priority of event for using use this output, order is `emergency|alert|critical|error|warning|notice|informational|debug or ""`                                                | `debug`                                                                                           |
| `config.aws.sqs.url`                        | AWS SQS Queue URL, if not empty, AWS SQS output is enabled                                                                                                                             |                                                                                                   |
| `config.aws.sqs.minimumpriority`            | minimum priority of event for using use this output, order is `emergency|alert|critical|error|warning|notice|informational|debug or ""`                                                | `debug`                                                                                           |
| `config.smtp.hostport`                      | "host:port" address of SMTP server, if not empty, SMTP output is enabled                                                                                                               |                                                                                                   |
| `config.smtp.user`                          | user to access SMTP server                                                                                                                                                             |                                                                                                   |
| `config.smtp.password`                      | password to access SMTP server                                                                                                                                                         |                                                                                                   |
| `config.smtp.from`                          | Sender address (mandatory if SMTP output is enabled)                                                                                                                                   |                                                                                                   |
| `config.smtp.to`                            | comma-separated list of Recipident addresses, can't be empty (mandatory if SMTP output is enabled)                                                                                     |                                                                                                   |
| `config.smtp.outputformat`                  | html, text                                                                                                                                                                             | `html`                                                                                            |
| `config.smtp.minimumpriority`               | minimum priority of event for using use this output, order is `emergency|alert|critical|error|warning|notice|informational|debug or ""`                                                | `debug`                                                                                           |
| `config.opsgenie.apikey`                    | Opsgenie API Key, if not empty, Opsgenie output is enabled                                                                                                                             |                                                                                                   |
| `config.opsgenie.region`                    | (`us` or `eu`) region of your domain                                                                                                                                                   | `us`                                                                                              |
| `config.opsgenie.minimumpriority`           | minimum priority of event for using use this output, order is `emergency|alert|critical|error|warning|notice|informational|debug or ""`                                                | `debug`                                                                                           |
| `config.statsd.forwarder`                   | The address for the StatsD forwarder, in the form http://host:port, if not empty StatsD is enabled                                                                                     |                                                                                                   |
| `config.statsd.namespace`                   | A prefix for all metrics                                                                                                                                                               | `falcosidekick`                                                                                   |
| `config.dogstatsd.forwarder`                | The address for the DogStatsD forwarder, in the form http://host:port, if not empty DogStatsD is enabled                                                                               |                                                                                                   |
| `config.dogstatsd.namespace`                | A prefix for all metrics                                                                                                                                                               | `falcosidekick`                                                                                   |
| `config.dogstatsd.tags`                     | A comma-separated list of tags to add to all metrics                                                                                                                                   |                                                                                                   |
| `config.webhook.address`                    | Webhook address, if not empty, Webhook output is enabled                                                                                                                               |                                                                                                   |
| `config.webhook.customHeaders`              | a list of comma separated custom headers to add, syntax is "key:value,key:value"                                                                                                       |                                                                                                   |
| `config.webhook.minimumpriority`            | minimum priority of event for using use this output, order is `emergency|alert|critical|error|warning|notice|informational|debug or ""`                                                | `debug`                                                                                           |
| `config.azure.eventhub.name`                | Name of the Hub, if not empty, EventHub is *enabled*                                                                                                                                   |                                                                                                   |
| `config.azure.eventhub.namespace`           | Name of the space the Hub is in                                                                                                                                                        |                                                                                                   |
| `config.azure.eventhub.minimumpriority`     | minimum priority of event for using use this output, order is `emergency|alert|critical|error|warning|notice|informational|debug or ""`                                                | `debug`                                                                                           |
| `config.gcp.credentials`                    | Base64 encoded JSON key file for the GCP service account                                                                                                                               |                                                                                                   |
| `config.gcp.pubsub.projectid`               | ID of the GCP project                                                                                                                                                                  |                                                                                                   |
| `config.gcp.pubsub.topic`                   | Name of the Pub/Sub topic                                                                                                                                                              |                                                                                                   |
| `config.gcp.eventhub.minimumpriority`       | minimum priority of event for using use this output, order is `emergency|alert|critical|error|warning|notice|informational|debug or ""`                                                | `debug`                                                                                           |
| `config.googlechat.webhookurl`              | Google Chat Webhook URL (ex: https://chat.googleapis.com/v1/spaces/XXXXXX/YYYYYY), if not `empty`, Google Chat output is *enabled*                                                     |                                                                                                   |
| `config.googlechat.outputformat`            | `all` (default), `text` (only text is displayed in Google chat)                                                                                                                        | `all`                                                                                             |
| `config.googlechat.minimumpriority`         | minimum priority of event for using use this output, order is `emergency|alert|critical|error|warning|notice|informational|debug or ""`                                                | `debug`                                                                                           |
| `config.googlechat.messageformat`           | a Go template to format Google Chat Text above Attachment, displayed in addition to the output from `config.googlechat.outputformat`. If empty, no Text is displayed before Attachment |                                                                                                   |
| `config.kafka.url`                          | Apache Kafka URL (ex: http://kafka). Defaults to port 9092 if no port is specified after the domain, if not empty, Kafka output is *enabled*                                           |                                                                                                   |
| `config.kafka.topic`                        | `all` (default), `text` (only text is displayed in Google chat)                                                                                                                        | `all`                                                                                             |
| `config.kafka.partition`                    | a Go template to format Google Chat Text above Attachment, displayed in addition to the output from `config.googlechat.outputformat`. If empty, no Text is displayed before Attachment |                                                                                                   |
| `config.kafka.minimumpriority`              | minimum priority of event for using use this output, order is `emergency|alert|critical|error|warning|notice|informational|debug or ""`                                                | `debug`                                                                                           |
| `config.pagerduty.apiKey`                   | Pagerduty API Key, if not empty, Pagerduty output is enabled                                                                                                                           |                                                                                                   |
| `config.pagerduty.service`                  | Service to create an incident (mandatory)                                                                                                                                              |                                                                                                   |
| `config.pagerduty.assignee`                 | A list of comma separated users to assign. Cannot be provided if pagerduty.escalationpolicy is already specified                                                                       |                                                                                                   |
| `config.pagerduty.escalationpolicy`         | Escalation policy to assign. Cannot be provided if pagerduty.escalationpolicy is already specified                                                                                     |                                                                                                   |
| `config.pagerduty.minimumpriority`          | minimum priority of event for using use this output, order is `emergency|alert|critical|error|warning|notice|informational|debug or ""`                                                | `debug`                                                                                           |
| `kubeless.function`                         | Name of Kubeless function, if not empty, EventHub is *enabled*                                                                                                                         |                                                                                                   |
| `kubeless.namespace`                        | Namespace of Kubeless function (mandatory)                                                                                                                                             |                                                                                                   |
| `kubeless.port`                             | Port of service of Kubeless function. Default is `8080`.                                                                                                                               |                                                                                                   |
| `kubeless.minimumpriority`                  | minimum priority of event for using use this output, order is `emergency|alert|critical|error|warning|notice|informational|debug or ""`                                                | `debug`                                                                                           |

Specify each parameter using the `--set key=value[,key=value]` argument to `helm install`. For example,

```bash
helm install falcosidekick --set debug=true falcosecurity/falcosidekick
```

Alternatively, a YAML file that specifies the values for the parameters can be provided while installing the chart. For example:

```bash
helm install falcosidekick -f values.yaml falcosecurity/falcosidekick
```

> **Tip**: You can use the default [values.yaml](values.yaml)

## Metrics

A `prometheus` endpoint can be scrapped at `/metrics`.
