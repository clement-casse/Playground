package cyphergraphexporter // import "github.com/clement-casse/playground/otelcol-custom/exporter/cyphergraphexporter"

import (
	"errors"
	"net/url"

	"go.opentelemetry.io/collector/config/configopaque"
	"go.opentelemetry.io/collector/config/configretry"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

const (
	defaultDatabaseURI = "bolt://localhost:7687"
	defaultUserAgent   = "opentelemetrycollector.cyphergraphexporter"
)

var (
	// defaultResourcesMappers provides a selection of ResourceMappers following the
	// OpenTelemetry Resource Semantic [https://opentelemetry.io/docs/specs/semconv/resource/]
	defaultResourcesMappers = map[string]ResourceMapper{
		"service": {
			LabelID:     string(semconv.ServiceNameKey),
			OtherLabels: []string{string(semconv.ServiceVersionKey), string(semconv.ServiceInstanceIDKey), string(semconv.ServiceNamespaceKey)},
		},
		"container": {
			LabelID:     string(semconv.ContainerIDKey),
			OtherLabels: []string{string(semconv.ContainerNameKey)},
		},
		"container.image": {
			LabelID:     string(semconv.ContainerImageIDKey),
			OtherLabels: []string{string(semconv.ContainerImageNameKey)},
		},
		"host": {
			LabelID:     string(semconv.HostIDKey),
			OtherLabels: []string{string(semconv.HostNameKey), string(semconv.HostImageIDKey), string(semconv.HostImageNameKey), string(semconv.HostIPKey), string(semconv.HostTypeKey)},
		},
		"deployment": {
			LabelID: string(semconv.DeploymentEnvironmentKey),
		},
		"k8s.cluster": {
			LabelID:     string(semconv.K8SClusterUIDKey),
			OtherLabels: []string{string(semconv.K8SClusterNameKey)},
		},
		"k8s.node": {
			LabelID:     string(semconv.K8SNodeNameKey),
			OtherLabels: []string{string(semconv.K8SNodeUIDKey)},
		},
		"k8s.pod": {
			LabelID:     string(semconv.K8SPodUIDKey),
			OtherLabels: []string{string(semconv.K8SPodNameKey)},
		},
	}
)

// Config defines configuration for the Cypher Graph exporter.
type Config struct {
	configretry.BackOffConfig `mapstructure:"retry_on_failure"`

	// DatabaseURI is the target address of a Neo4j instance. The URI is the Bolt URI
	// thats uses the following scheme:
	//   - 'bolt', 'bolt+s' or 'bolt+ssc' allow to connects to a single instance database
	//   - 'neo4j', 'neo4j+s' or 'neo4j+ssc' allow to connect to a cluster of databases
	//
	// Refer to [neo4j.NewDriverWithContext] documentation to understand how this value
	// is used.
	DatabaseURI string `mapstructure:"db_uri,omitempty"`

	// Username corresponds to the user name used in the basic authentication process
	// to the database.
	Username string `mapstructure:"username,omitempty"`

	// Password corresponds to the secret used in the basic authentication process
	// to the database.
	Password configopaque.String `mapstructure:"password,omitempty"`

	// BearerToken corresponds to a base64-encoded string generated by a Single Sign-On Provider.
	BearerToken configopaque.String `mapstructure:"bearer_token,omitempty"`

	// KerberosTicket corresponds to a base64-encoded string representing the kerberos ticket.
	KerberosTicket configopaque.String `mapstructure:"kerberos_ticket,omitempty"`

	// UserAgent corresponds to the value of the User-Agent field that is used in
	// the "bolt" (websocket) connection.
	UserAgent string `mapstructure:"user_agent,omitempty"`

	// RessourceMappers allows to generate multiple nodes with different labels in the graph
	// for a single [go.opentelemetry.io/collector/pdata/pcommon.Resource] based on its attributes.
	// Graph node labels are provided as keys of the ResourceMappers field:
	//   resources:
	//     "k8s.pod":
	//       label_id: "k8s.pod.uid"
	//       other_labels: [ "k8s.pod.name" ]
	//     "k8s.node":
	//       label_id: "k8s.node.uid"
	//       other_labels: [ "k8s.node.name" ]
	// This configuration will generate for each spans processed by the receiver two different
	// resource nodes: one with label 'k8s.node' and one with label 'k8s.pod'. Nodes with the same
	// label sharing the same value of the attribute referenced in the 'label_id' field will be
	// merged together. This configuration will only generate one node for each single pods more
	// one for each single nodes in the kubernetes cluster.
	ResourceMappers map[string]ResourceMapper `mapstructure:"resources,omitempty"`
}

// ResourceMapper describes how [go.opentelemetry.io/collector/pdata/pcommon.Resource] attributes
// get parsed to create a node in the graph.
type ResourceMapper struct {
	// LabelID represents the attrribute Key of a resource whose value will be used as index for
	// for a given resource node.
	LabelID string `mapstructure:"label_id"`

	// OtherLabels represents additionnal labels whose value will be added as additional properties
	// for the nodes of the graph.
	OtherLabels []string `mapstructure:"other_labels"`
}

var (
	errMultipleAuthMethod = errors.New("cannot mix multiple authentication methods")
)

func (cfg *Config) Validate() error {
	if _, err := url.Parse(cfg.DatabaseURI); err != nil {
		return err
	}
	if cfg.Username != "" && (string(cfg.BearerToken) != "" || string(cfg.KerberosTicket) != "") {
		return errMultipleAuthMethod
	}
	if string(cfg.BearerToken) != "" && string(cfg.KerberosTicket) != "" {
		return errMultipleAuthMethod
	}
	return nil
}
