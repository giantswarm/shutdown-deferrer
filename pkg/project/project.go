package project

var (
	description = "The shutdown deferrer is a service monitoring DrainerConfig status for POD that it's running within. It serves as a gatekeeper for POD termination"
	gitSHA      = "n/a"
	name        = "shutdown-deferrer"
	source      = "https://github.com/giantswarm/shutdown-deferrer"
	version     = "0.1.0-dev"
)

func Description() string {
	return description
}

func GitSHA() string {
	return gitSHA
}

func Name() string {
	return name
}

func Source() string {
	return source
}

func Version() string {
	return version
}
