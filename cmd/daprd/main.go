/*
Copyright 2021 The Dapr Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"github.com/liuxd6825/components-contrib/liuxd/applog"
	applog_loader "github.com/liuxd6825/dapr/pkg/components/liuxd/applogger"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/valyala/fasthttp"
	"go.uber.org/automaxprocs/maxprocs"

	"github.com/liuxd6825/dapr/pkg/runtime"
	"github.com/dapr/kit/logger"

	// Included components in compiled daprd.

	// Secret stores.
	"github.com/liuxd6825/components-contrib/secretstores"
	alicloud_paramstore "github.com/liuxd6825/components-contrib/secretstores/alicloud/parameterstore"
	"github.com/liuxd6825/components-contrib/secretstores/aws/parameterstore"
	"github.com/liuxd6825/components-contrib/secretstores/aws/secretmanager"
	"github.com/liuxd6825/components-contrib/secretstores/azure/keyvault"
	gcp_secretmanager "github.com/liuxd6825/components-contrib/secretstores/gcp/secretmanager"
	"github.com/liuxd6825/components-contrib/secretstores/hashicorp/vault"
	secretstore_kubernetes "github.com/liuxd6825/components-contrib/secretstores/kubernetes"
	secretstore_env "github.com/liuxd6825/components-contrib/secretstores/local/env"
	secretstore_file "github.com/liuxd6825/components-contrib/secretstores/local/file"

	secretstores_loader "github.com/liuxd6825/dapr/pkg/components/secretstores"

	// State Stores.
	"github.com/liuxd6825/components-contrib/state"
	"github.com/liuxd6825/components-contrib/state/aerospike"
	state_dynamodb "github.com/liuxd6825/components-contrib/state/aws/dynamodb"
	state_azure_blobstorage "github.com/liuxd6825/components-contrib/state/azure/blobstorage"
	state_cosmosdb "github.com/liuxd6825/components-contrib/state/azure/cosmosdb"
	state_azure_tablestorage "github.com/liuxd6825/components-contrib/state/azure/tablestorage"
	"github.com/liuxd6825/components-contrib/state/cassandra"
	"github.com/liuxd6825/components-contrib/state/cockroachdb"
	"github.com/liuxd6825/components-contrib/state/couchbase"
	"github.com/liuxd6825/components-contrib/state/gcp/firestore"
	"github.com/liuxd6825/components-contrib/state/hashicorp/consul"
	"github.com/liuxd6825/components-contrib/state/hazelcast"
	state_jetstream "github.com/liuxd6825/components-contrib/state/jetstream"
	"github.com/liuxd6825/components-contrib/state/memcached"
	"github.com/liuxd6825/components-contrib/state/mongodb"
	state_mysql "github.com/liuxd6825/components-contrib/state/mysql"
	state_oci_objectstorage "github.com/liuxd6825/components-contrib/state/oci/objectstorage"
	state_oracledatabase "github.com/liuxd6825/components-contrib/state/oracledatabase"
	"github.com/liuxd6825/components-contrib/state/postgresql"
	state_redis "github.com/liuxd6825/components-contrib/state/redis"
	"github.com/liuxd6825/components-contrib/state/rethinkdb"
	"github.com/liuxd6825/components-contrib/state/sqlserver"
	"github.com/liuxd6825/components-contrib/state/zookeeper"

	state_loader "github.com/liuxd6825/dapr/pkg/components/state"

	// Pub/Sub.
	pubs "github.com/liuxd6825/components-contrib/pubsub"
	pubsub_snssqs "github.com/liuxd6825/components-contrib/pubsub/aws/snssqs"
	pubsub_eventhubs "github.com/liuxd6825/components-contrib/pubsub/azure/eventhubs"
	"github.com/liuxd6825/components-contrib/pubsub/azure/servicebus"
	pubsub_gcp "github.com/liuxd6825/components-contrib/pubsub/gcp/pubsub"
	pubsub_hazelcast "github.com/liuxd6825/components-contrib/pubsub/hazelcast"
	pubsub_inmemory "github.com/liuxd6825/components-contrib/pubsub/in-memory"
	pubsub_jetstream "github.com/liuxd6825/components-contrib/pubsub/jetstream"
	pubsub_kafka "github.com/liuxd6825/components-contrib/pubsub/kafka"
	pubsub_mqtt "github.com/liuxd6825/components-contrib/pubsub/mqtt"
	"github.com/liuxd6825/components-contrib/pubsub/natsstreaming"
	pubsub_pulsar "github.com/liuxd6825/components-contrib/pubsub/pulsar"
	"github.com/liuxd6825/components-contrib/pubsub/rabbitmq"
	pubsub_redis "github.com/liuxd6825/components-contrib/pubsub/redis"
	configuration_loader "github.com/liuxd6825/dapr/pkg/components/configuration"
	pubsub_loader "github.com/liuxd6825/dapr/pkg/components/pubsub"

	// Name resolutions.
	nr "github.com/liuxd6825/components-contrib/nameresolution"
	nr_consul "github.com/liuxd6825/components-contrib/nameresolution/consul"
	nr_kubernetes "github.com/liuxd6825/components-contrib/nameresolution/kubernetes"
	nr_mdns "github.com/liuxd6825/components-contrib/nameresolution/mdns"

	nr_loader "github.com/liuxd6825/dapr/pkg/components/nameresolution"

	// Bindings.
	"github.com/liuxd6825/components-contrib/bindings"
	dingtalk_webhook "github.com/liuxd6825/components-contrib/bindings/alicloud/dingtalk/webhook"
	"github.com/liuxd6825/components-contrib/bindings/alicloud/oss"
	"github.com/liuxd6825/components-contrib/bindings/alicloud/tablestore"
	"github.com/liuxd6825/components-contrib/bindings/apns"
	"github.com/liuxd6825/components-contrib/bindings/aws/dynamodb"
	"github.com/liuxd6825/components-contrib/bindings/aws/kinesis"
	"github.com/liuxd6825/components-contrib/bindings/aws/s3"
	"github.com/liuxd6825/components-contrib/bindings/aws/ses"
	"github.com/liuxd6825/components-contrib/bindings/aws/sns"
	"github.com/liuxd6825/components-contrib/bindings/aws/sqs"
	"github.com/liuxd6825/components-contrib/bindings/azure/blobstorage"
	bindings_cosmosdb "github.com/liuxd6825/components-contrib/bindings/azure/cosmosdb"
	bindings_cosmosdbgremlinapi "github.com/liuxd6825/components-contrib/bindings/azure/cosmosdbgremlinapi"
	"github.com/liuxd6825/components-contrib/bindings/azure/eventgrid"
	"github.com/liuxd6825/components-contrib/bindings/azure/eventhubs"
	"github.com/liuxd6825/components-contrib/bindings/azure/servicebusqueues"
	"github.com/liuxd6825/components-contrib/bindings/azure/signalr"
	"github.com/liuxd6825/components-contrib/bindings/azure/storagequeues"
	"github.com/liuxd6825/components-contrib/bindings/cron"
	"github.com/liuxd6825/components-contrib/bindings/gcp/bucket"
	"github.com/liuxd6825/components-contrib/bindings/gcp/pubsub"
	"github.com/liuxd6825/components-contrib/bindings/graphql"
	"github.com/liuxd6825/components-contrib/bindings/http"
	"github.com/liuxd6825/components-contrib/bindings/influx"
	"github.com/liuxd6825/components-contrib/bindings/kafka"
	"github.com/liuxd6825/components-contrib/bindings/kubernetes"
	"github.com/liuxd6825/components-contrib/bindings/localstorage"
	"github.com/liuxd6825/components-contrib/bindings/mqtt"
	"github.com/liuxd6825/components-contrib/bindings/mysql"
	"github.com/liuxd6825/components-contrib/bindings/postgres"
	"github.com/liuxd6825/components-contrib/bindings/postmark"
	bindings_rabbitmq "github.com/liuxd6825/components-contrib/bindings/rabbitmq"
	"github.com/liuxd6825/components-contrib/bindings/redis"
	"github.com/liuxd6825/components-contrib/bindings/rethinkdb/statechange"
	"github.com/liuxd6825/components-contrib/bindings/smtp"
	"github.com/liuxd6825/components-contrib/bindings/twilio/sendgrid"
	"github.com/liuxd6825/components-contrib/bindings/twilio/sms"
	"github.com/liuxd6825/components-contrib/bindings/twitter"
	bindings_zeebe_command "github.com/liuxd6825/components-contrib/bindings/zeebe/command"
	bindings_zeebe_jobworker "github.com/liuxd6825/components-contrib/bindings/zeebe/jobworker"

	bindings_loader "github.com/liuxd6825/dapr/pkg/components/bindings"

	// HTTP Middleware.

	middleware "github.com/liuxd6825/components-contrib/middleware"
	"github.com/liuxd6825/components-contrib/middleware/http/bearer"
	"github.com/liuxd6825/components-contrib/middleware/http/oauth2"
	"github.com/liuxd6825/components-contrib/middleware/http/oauth2clientcredentials"
	"github.com/liuxd6825/components-contrib/middleware/http/opa"
	"github.com/liuxd6825/components-contrib/middleware/http/ratelimit"
	"github.com/liuxd6825/components-contrib/middleware/http/sentinel"

	http_middleware_loader "github.com/liuxd6825/dapr/pkg/components/middleware/http"
	http_middleware "github.com/liuxd6825/dapr/pkg/middleware/http"

	"github.com/liuxd6825/components-contrib/configuration"
	configuration_redis "github.com/liuxd6825/components-contrib/configuration/redis"

	applog_mongo "github.com/liuxd6825/components-contrib/liuxd/applog/mongo"
	eventstorage "github.com/liuxd6825/components-contrib/liuxd/eventstorage"
	eventstorage_mongo "github.com/liuxd6825/components-contrib/liuxd/eventstorage/es_mongo"
	eventstorage_loader "github.com/liuxd6825/dapr/pkg/components/liuxd/eventstorage"
)

var (
	log        = logger.NewLogger("dapr.runtime")
	logContrib = logger.NewLogger("dapr.contrib")
)

func main() {
	// set GOMAXPROCS
	_, _ = maxprocs.Set()

	rt, err := runtime.FromFlags()
	if err != nil {
		log.Fatal(err)
	}

	err = rt.Run(
		runtime.WithEventStorage(
			// lxd 注册mongodb事件溯源组件
			eventstorage_loader.New("mongodb", func() eventstorage.EventStorage {
				return eventstorage_mongo.NewMongoEventSourcing(logContrib)
			}),
		),
		runtime.WithApplog(
			// lxd 注册mongodb事件溯源组件
			applog_loader.New("mongodb", func() applog.Logger {
				return applog_mongo.NewLogger(logContrib)
			}),
		),
		runtime.WithSecretStores(
			secretstores_loader.New("kubernetes", func() secretstores.SecretStore {
				return secretstore_kubernetes.NewKubernetesSecretStore(logContrib)
			}),
			secretstores_loader.New("azure.keyvault", func() secretstores.SecretStore {
				return keyvault.NewAzureKeyvaultSecretStore(logContrib)
			}),
			secretstores_loader.New("hashicorp.vault", func() secretstores.SecretStore {
				return vault.NewHashiCorpVaultSecretStore(logContrib)
			}),
			secretstores_loader.New("aws.secretmanager", func() secretstores.SecretStore {
				return secretmanager.NewSecretManager(logContrib)
			}),
			secretstores_loader.New("aws.parameterstore", func() secretstores.SecretStore {
				return parameterstore.NewParameterStore(logContrib)
			}),
			secretstores_loader.New("gcp.secretmanager", func() secretstores.SecretStore {
				return gcp_secretmanager.NewSecreteManager(logContrib)
			}),
			secretstores_loader.New("local.file", func() secretstores.SecretStore {
				return secretstore_file.NewLocalSecretStore(logContrib)
			}),
			secretstores_loader.New("local.env", func() secretstores.SecretStore {
				return secretstore_env.NewEnvSecretStore(logContrib)
			}),
			secretstores_loader.New("alicloud.parameterstore", func() secretstores.SecretStore {
				return alicloud_paramstore.NewParameterStore(logContrib)
			}),
		),
		runtime.WithStates(
			state_loader.New("redis", func() state.Store {
				return state_redis.NewRedisStateStore(logContrib)
			}),
			state_loader.New("consul", func() state.Store {
				return consul.NewConsulStateStore(logContrib)
			}),
			state_loader.New("azure.blobstorage", func() state.Store {
				return state_azure_blobstorage.NewAzureBlobStorageStore(logContrib)
			}),
			state_loader.New("azure.cosmosdb", func() state.Store {
				return state_cosmosdb.NewCosmosDBStateStore(logContrib)
			}),
			state_loader.New("azure.tablestorage", func() state.Store {
				return state_azure_tablestorage.NewAzureTablesStateStore(logContrib)
			}),
			state_loader.New("cassandra", func() state.Store {
				return cassandra.NewCassandraStateStore(logContrib)
			}),
			state_loader.New("memcached", func() state.Store {
				return memcached.NewMemCacheStateStore(logContrib)
			}),
			state_loader.New("mongodb", func() state.Store {
				return mongodb.NewMongoDB(logContrib)
			}),
			state_loader.New("zookeeper", func() state.Store {
				return zookeeper.NewZookeeperStateStore(logContrib)
			}),
			state_loader.New("gcp.firestore", func() state.Store {
				return firestore.NewFirestoreStateStore(logContrib)
			}),
			state_loader.New("postgresql", func() state.Store {
				return postgresql.NewPostgreSQLStateStore(logContrib)
			}),
			state_loader.New("sqlserver", func() state.Store {
				return sqlserver.NewSQLServerStateStore(logContrib)
			}),
			state_loader.New("hazelcast", func() state.Store {
				return hazelcast.NewHazelcastStore(logContrib)
			}),
			state_loader.New("couchbase", func() state.Store {
				return couchbase.NewCouchbaseStateStore(logContrib)
			}),
			state_loader.New("aerospike", func() state.Store {
				return aerospike.NewAerospikeStateStore(logContrib)
			}),
			state_loader.New("rethinkdb", func() state.Store {
				return rethinkdb.NewRethinkDBStateStore(logContrib)
			}),
			state_loader.New("aws.dynamodb", state_dynamodb.NewDynamoDBStateStore),
			state_loader.New("mysql", func() state.Store {
				return state_mysql.NewMySQLStateStore(logContrib)
			}),
			state_loader.New("oci.objectstorage", func() state.Store {
				return state_oci_objectstorage.NewOCIObjectStorageStore(logContrib)
			}),
			state_loader.New("jetstream", func() state.Store {
				return state_jetstream.NewJetstreamStateStore(logContrib)
			}),
			state_loader.New("oracledatabase", func() state.Store {
				return state_oracledatabase.NewOracleDatabaseStateStore(logContrib)
			}),
			state_loader.New("cockroachdb", func() state.Store {
				return cockroachdb.New(logContrib)
			}),
		),
		runtime.WithConfigurations(
			configuration_loader.New("redis", func() configuration.Store {
				return configuration_redis.NewRedisConfigurationStore(logContrib)
			}),
		),
		runtime.WithPubSubs(
			pubsub_loader.New("azure.eventhubs", func() pubs.PubSub {
				return pubsub_eventhubs.NewAzureEventHubs(logContrib)
			}),
			pubsub_loader.New("azure.servicebus", func() pubs.PubSub {
				return servicebus.NewAzureServiceBus(logContrib)
			}),
			pubsub_loader.New("gcp.pubsub", func() pubs.PubSub {
				return pubsub_gcp.NewGCPPubSub(logContrib)
			}),
			pubsub_loader.New("hazelcast", func() pubs.PubSub {
				return pubsub_hazelcast.NewHazelcastPubSub(logContrib)
			}),
			pubsub_loader.New("jetstream", func() pubs.PubSub {
				return pubsub_jetstream.NewJetStream(logContrib)
			}),
			pubsub_loader.New("kafka", func() pubs.PubSub {
				return pubsub_kafka.NewKafka(logContrib)
			}),
			pubsub_loader.New("mqtt", func() pubs.PubSub {
				return pubsub_mqtt.NewMQTTPubSub(logContrib)
			}),
			pubsub_loader.New("natsstreaming", func() pubs.PubSub {
				return natsstreaming.NewNATSStreamingPubSub(logContrib)
			}),
			pubsub_loader.New("pulsar", func() pubs.PubSub {
				return pubsub_pulsar.NewPulsar(logContrib)
			}),
			pubsub_loader.New("rabbitmq", func() pubs.PubSub {
				return rabbitmq.NewRabbitMQ(logContrib)
			}),
			pubsub_loader.New("redis", func() pubs.PubSub {
				return pubsub_redis.NewRedisStreams(logContrib)
			}),
			pubsub_loader.New("snssqs", func() pubs.PubSub {
				return pubsub_snssqs.NewSnsSqs(logContrib)
			}),
			pubsub_loader.New("in-memory", func() pubs.PubSub {
				return pubsub_inmemory.New(logContrib)
			}),
		),
		runtime.WithNameResolutions(
			nr_loader.New("mdns", func() nr.Resolver {
				return nr_mdns.NewResolver(logContrib)
			}),
			nr_loader.New("kubernetes", func() nr.Resolver {
				return nr_kubernetes.NewResolver(logContrib)
			}),
			nr_loader.New("consul", func() nr.Resolver {
				return nr_consul.NewResolver(logContrib)
			}),
		),
		runtime.WithInputBindings(
			bindings_loader.NewInput("aws.sqs", func() bindings.InputBinding {
				return sqs.NewAWSSQS(logContrib)
			}),
			bindings_loader.NewInput("aws.kinesis", func() bindings.InputBinding {
				return kinesis.NewAWSKinesis(logContrib)
			}),
			bindings_loader.NewInput("azure.eventgrid", func() bindings.InputBinding {
				return eventgrid.NewAzureEventGrid(logContrib)
			}),
			bindings_loader.NewInput("azure.eventhubs", func() bindings.InputBinding {
				return eventhubs.NewAzureEventHubs(logContrib)
			}),
			bindings_loader.NewInput("azure.servicebusqueues", func() bindings.InputBinding {
				return servicebusqueues.NewAzureServiceBusQueues(logContrib)
			}),
			bindings_loader.NewInput("azure.storagequeues", func() bindings.InputBinding {
				return storagequeues.NewAzureStorageQueues(logContrib)
			}),
			bindings_loader.NewInput("cron", func() bindings.InputBinding {
				return cron.NewCron(logContrib)
			}),
			bindings_loader.NewInput("dingtalk.webhook", func() bindings.InputBinding {
				return dingtalk_webhook.NewDingTalkWebhook(logContrib)
			}),
			bindings_loader.NewInput("gcp.pubsub", func() bindings.InputBinding {
				return pubsub.NewGCPPubSub(logContrib)
			}),
			bindings_loader.NewInput("kafka", func() bindings.InputBinding {
				return kafka.NewKafka(logContrib)
			}),
			bindings_loader.NewInput("kubernetes", func() bindings.InputBinding {
				return kubernetes.NewKubernetes(logContrib)
			}),
			bindings_loader.NewInput("mqtt", func() bindings.InputBinding {
				return mqtt.NewMQTT(logContrib)
			}),
			bindings_loader.NewInput("rabbitmq", func() bindings.InputBinding {
				return bindings_rabbitmq.NewRabbitMQ(logContrib)
			}),
			bindings_loader.NewInput("rethinkdb.statechange", func() bindings.InputBinding {
				return statechange.NewRethinkDBStateChangeBinding(logContrib)
			}),
			bindings_loader.NewInput("twitter", func() bindings.InputBinding {
				return twitter.NewTwitter(logContrib)
			}),
			bindings_loader.NewInput("zeebe.jobworker", func() bindings.InputBinding {
				return bindings_zeebe_jobworker.NewZeebeJobWorker(logContrib)
			}),
		),
		runtime.WithOutputBindings(
			bindings_loader.NewOutput("alicloud.oss", func() bindings.OutputBinding {
				return oss.NewAliCloudOSS(logContrib)
			}),
			bindings_loader.NewOutput("alicloud.tablestore", func() bindings.OutputBinding {
				return tablestore.NewAliCloudTableStore(log)
			}),
			bindings_loader.NewOutput("apns", func() bindings.OutputBinding {
				return apns.NewAPNS(logContrib)
			}),
			bindings_loader.NewOutput("aws.s3", func() bindings.OutputBinding {
				return s3.NewAWSS3(logContrib)
			}),
			bindings_loader.NewOutput("aws.ses", func() bindings.OutputBinding {
				return ses.NewAWSSES(logContrib)
			}),
			bindings_loader.NewOutput("aws.sqs", func() bindings.OutputBinding {
				return sqs.NewAWSSQS(logContrib)
			}),
			bindings_loader.NewOutput("aws.sns", func() bindings.OutputBinding {
				return sns.NewAWSSNS(logContrib)
			}),
			bindings_loader.NewOutput("aws.kinesis", func() bindings.OutputBinding {
				return kinesis.NewAWSKinesis(logContrib)
			}),
			bindings_loader.NewOutput("aws.dynamodb", func() bindings.OutputBinding {
				return dynamodb.NewDynamoDB(logContrib)
			}),
			bindings_loader.NewOutput("azure.blobstorage", func() bindings.OutputBinding {
				return blobstorage.NewAzureBlobStorage(logContrib)
			}),
			bindings_loader.NewOutput("azure.cosmosdb", func() bindings.OutputBinding {
				return bindings_cosmosdb.NewCosmosDB(logContrib)
			}),
			bindings_loader.NewOutput("azure.cosmosdb.gremlinapi", func() bindings.OutputBinding {
				return bindings_cosmosdbgremlinapi.NewCosmosDBGremlinAPI(logContrib)
			}),
			bindings_loader.NewOutput("azure.eventgrid", func() bindings.OutputBinding {
				return eventgrid.NewAzureEventGrid(logContrib)
			}),
			bindings_loader.NewOutput("azure.eventhubs", func() bindings.OutputBinding {
				return eventhubs.NewAzureEventHubs(logContrib)
			}),
			bindings_loader.NewOutput("azure.servicebusqueues", func() bindings.OutputBinding {
				return servicebusqueues.NewAzureServiceBusQueues(logContrib)
			}),
			bindings_loader.NewOutput("azure.signalr", func() bindings.OutputBinding {
				return signalr.NewSignalR(logContrib)
			}),
			bindings_loader.NewOutput("azure.storagequeues", func() bindings.OutputBinding {
				return storagequeues.NewAzureStorageQueues(logContrib)
			}),
			bindings_loader.NewOutput("cron", func() bindings.OutputBinding {
				return cron.NewCron(logContrib)
			}),
			bindings_loader.NewOutput("dingtalk.webhook", func() bindings.OutputBinding {
				return dingtalk_webhook.NewDingTalkWebhook(logContrib)
			}),
			bindings_loader.NewOutput("gcp.bucket", func() bindings.OutputBinding {
				return bucket.NewGCPStorage(logContrib)
			}),
			bindings_loader.NewOutput("gcp.pubsub", func() bindings.OutputBinding {
				return pubsub.NewGCPPubSub(logContrib)
			}),
			bindings_loader.NewOutput("http", func() bindings.OutputBinding {
				return http.NewHTTP(logContrib)
			}),
			bindings_loader.NewOutput("influx", func() bindings.OutputBinding {
				return influx.NewInflux(logContrib)
			}),
			bindings_loader.NewOutput("kafka", func() bindings.OutputBinding {
				return kafka.NewKafka(logContrib)
			}),
			bindings_loader.NewOutput("localstorage", func() bindings.OutputBinding {
				return localstorage.NewLocalStorage(logContrib)
			}),
			bindings_loader.NewOutput("mqtt", func() bindings.OutputBinding {
				return mqtt.NewMQTT(logContrib)
			}),
			bindings_loader.NewOutput("mysql", func() bindings.OutputBinding {
				return mysql.NewMysql(logContrib)
			}),
			bindings_loader.NewOutput("postgres", func() bindings.OutputBinding {
				return postgres.NewPostgres(logContrib)
			}),
			bindings_loader.NewOutput("postmark", func() bindings.OutputBinding {
				return postmark.NewPostmark(logContrib)
			}),
			bindings_loader.NewOutput("rabbitmq", func() bindings.OutputBinding {
				return bindings_rabbitmq.NewRabbitMQ(logContrib)
			}),
			bindings_loader.NewOutput("redis", func() bindings.OutputBinding {
				return redis.NewRedis(logContrib)
			}),
			bindings_loader.NewOutput("smtp", func() bindings.OutputBinding {
				return smtp.NewSMTP(logContrib)
			}),
			bindings_loader.NewOutput("twilio.sms", func() bindings.OutputBinding {
				return sms.NewSMS(logContrib)
			}),
			bindings_loader.NewOutput("twilio.sendgrid", func() bindings.OutputBinding {
				return sendgrid.NewSendGrid(logContrib)
			}),
			bindings_loader.NewOutput("twitter", func() bindings.OutputBinding {
				return twitter.NewTwitter(logContrib)
			}),
			bindings_loader.NewOutput("zeebe.command", func() bindings.OutputBinding {
				return bindings_zeebe_command.NewZeebeCommand(logContrib)
			}),
			bindings_loader.NewOutput("graphql", func() bindings.OutputBinding {
				return graphql.NewGraphQL(logContrib)
			}),
		),
		runtime.WithHTTPMiddleware(
			http_middleware_loader.New("uppercase", func(metadata middleware.Metadata) (http_middleware.Middleware, error) {
				return func(h fasthttp.RequestHandler) fasthttp.RequestHandler {
					return func(ctx *fasthttp.RequestCtx) {
						body := string(ctx.PostBody())
						ctx.Request.SetBody([]byte(strings.ToUpper(body)))
						h(ctx)
					}
				}, nil
			}),
			http_middleware_loader.New("oauth2", func(metadata middleware.Metadata) (http_middleware.Middleware, error) {
				return oauth2.NewOAuth2Middleware().GetHandler(metadata)
			}),
			http_middleware_loader.New("oauth2clientcredentials", func(metadata middleware.Metadata) (http_middleware.Middleware, error) {
				return oauth2clientcredentials.NewOAuth2ClientCredentialsMiddleware(log).GetHandler(metadata)
			}),
			http_middleware_loader.New("ratelimit", func(metadata middleware.Metadata) (http_middleware.Middleware, error) {
				return ratelimit.NewRateLimitMiddleware(log).GetHandler(metadata)
			}),
			http_middleware_loader.New("bearer", func(metadata middleware.Metadata) (http_middleware.Middleware, error) {
				return bearer.NewBearerMiddleware(log).GetHandler(metadata)
			}),
			http_middleware_loader.New("opa", func(metadata middleware.Metadata) (http_middleware.Middleware, error) {
				return opa.NewMiddleware(log).GetHandler(metadata)
			}),
			http_middleware_loader.New("sentinel", func(metadata middleware.Metadata) (http_middleware.Middleware, error) {
				return sentinel.NewMiddleware(log).GetHandler(metadata)
			}),
		),
	)
	if err != nil {
		log.Fatalf("fatal error from runtime: %s", err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, os.Interrupt)
	<-stop
	rt.ShutdownWithWait()
}
