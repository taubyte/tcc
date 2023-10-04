package parser

import (
	//lint:ignore ST1001 keeps defintion clean
	. "github.com/taubyte/tcc/schema"
)

func TaubyteAttributes(attrs ...*Attribute) []*Attribute {
	return append(
		attrs,
		String("id", IsCID()),
		String("name", IsVariableName()),
		String("description"),
		StringSlice("tags"),
	)
}

var TaubyteRessources = []*Node{
	DefineGroup("databases",
		DefineIter( //use name as "name"?
			TaubyteAttributes(
				String("match"),
				Bool("useRegex", Path("regex"), Compat("useRegex")),
				String("network-access", Path("access", "network"), InSet("all", "subnet", "host"), Default("all")),
				Int("replicas-min", Path("replicas", "min"), Default(1)),
				Int("replicas-max", Path("replicas", "max"), Default(3)),
				String("size", Path("storage", "size")),
				String("encryption-type", Path("encryption", "type")),
				String("encryption-key", Path("encryption", "key")),
			),
		)),
	DefineGroup("domains",
		DefineIter(
			TaubyteAttributes(
				String("fqdn", IsFqdn()),
				String("certificate", Path("certificate", "cert")),
				String("certificate-key", Path("certificate", "key")),
				String("certificate-type", Path("certificate", "type"), InSet("x509"), Default("x509")),
			),
		)),
	DefineGroup("functions",
		DefineIter(
			TaubyteAttributes(
				String("type", Path("trigger", "type"), InSet("http", "https", "pubsub", "p2p")),
				Bool("local", Path("trigger", "local")),
				String("pubsub-channel", Path("trigger", "channel")),
				String("p2p-protocol", Path("trigger", "protocol"), Compat("trigger", "service")),
				String("p2p-command", Path("trigger", "command")),
				String("http-method", Path("trigger", "method"), IsHttpMethod()),
				StringSlice("http-methods", Path("trigger", "methods")), // TO IMPLEMENT
				StringSlice("http-domains", Path("trigger", "domains"), Compat("domains")),
				StringSlice("http-paths", Path("trigger", "paths")),
				String("source"),
				String("timeout", Path("execution", "timeout")),
				String("memory", Path("execution", "memory")),
				String("call", Path("execution", "call")),
			),
		)),
	DefineGroup("libraries",
		DefineIter(
			TaubyteAttributes(
				String("path", Path("source", "path")),
				String("branch", Path("source", "branch")),
				String("git-provider", Path("source", Either("github")), Key()),
				String("github-id", Path("source", "github", "id")),
				String("github-fullname", Path("source", "github", "fullname")),
			),
		)),
	DefineGroup("messaging",
		DefineIter(
			TaubyteAttributes(
				Bool("local"),
				String("match", Path("channel", "match")),
				Bool("regex", Path("channel", "regex")),
				Bool("mqtt", Path("bridges", "mqtt", "enable")),
				Bool("websocket", Path("bridges", "websocket", "enable")),
			),
		)),
	DefineGroup("services",
		DefineIter(
			TaubyteAttributes(
				String("protocol"),
			),
		)),
	DefineGroup("smartops",
		DefineIter(
			TaubyteAttributes(
				String("source"),
				String("timeout", Path("execution", "timeout")),
				String("memory", Path("execution", "memory")),
				String("call", Path("execution", "call")),
			),
		)),
	DefineGroup("storages",
		DefineIter(
			TaubyteAttributes(
				String("match"),
				Bool("useRegex", Path("regex"), Compat("useRegex")),
				String("network-access", Path("access", "network"), InSet("all", "subnet", "host"), Default("all")),
				Bool("versioning", Path("object", "versioning")),
				String("ttl", Path("streaming", "ttl")),
				String("size", Path("streaming", Either("object", "streaming"), "ttl")),
			),
		)),
	DefineGroup("websites",
		DefineIter(
			TaubyteAttributes(
				StringSlice("domains", Path("domains")),
				StringSlice("paths", Path("paths"), Compat("source", "paths")),
				String("branch", Path("source", "branch")),
				String("git-provider", Path("source", Either("github")), Key()),
				String("github-id", Path("source", "github", "id")),
				String("github-fullname", Path("source", "github", "fullname")),
			),
		)),
}

var TaubyteProject = SchemaDefinition(
	Root(
		TaubyteAttributes(
			String("email", Path("notification", "email"), IsEmail()),
		),
		append(TaubyteRessources,
			DefineGroup("applications",
				DefineIterGroup(
					TaubyteAttributes(),
					TaubyteRessources...,
				),
			),
		)...,
	),
)
