package cst

const (
	// Docker

	HttpdPortFile    = "/etc/apache2/port.conf"
	HttpdConfFile    = "/etc/apache2/apache2.conf"
	HttpdVhostsDir   = "/etc/apache2/sites-available/"
	HttpdHost        = "http.front.com"
	HttpdPort        = "80"
	DockerConfigFile = "/usr/go-apache-test/config.json"
	DockerVhostsDir  = "/usr/go-apache-test/vhosts/"
	DockerTestsDir   = "/usr/go-apache-test/tests/"
	DockerHostsFile  = "/etc/hosts"
	BackendHost      = "backend.front.com"
	BackendPort      = "3011"

	// local

	// HttpdPortFile    = "tmp/apache2/port.conf"
	// HttpdConfFile    = "tmp/apache2/apache2.conf"
	// HttpdVhostsDir   = "tmp/apache2/sites-available/"
	// HttpdHost        = "http.front.com"
	// HttpdPort        = "80"
	// DockerConfigFile = "tmp/config.json"
	// DockerVhostsDir  = "tmp/vhosts/"
	// DockerTestsDir   = "tmp/tests/"
	// DockerHostsFile  = "tmp/hosts"
	// BackendHost      = "backend.front.com"
	// BackendPort      = "3011"
)
