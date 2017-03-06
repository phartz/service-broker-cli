package sbcli

const (
	UsageText = "NAME:\n" +
		"   sb - A command line tool to interact with a9s Service Broker\n" +
		"USAGE:\n" +
		"   sb [global options] command [arguments...] [command options]\n"
)

func GetHelpText(command string) string {
	switch command {
	case "create-service", "cs":
		return "NAME:\n" +
			"   create-service - Create a service instance\n" +
			"\n" +
			"USAGE:\n" +
			"   sb create-service SERVICE PLAN SERVICE_INSTANCE [-c PARAMETERS_AS_JSON] [-t TAGS]\n" +
			"\n" +
			"   Optionally provide service-specific configuration parameters in a valid JSON object in-line:\n" +
			"\n" +
			"   sb create-service SERVICE PLAN SERVICE_INSTANCE -c '{\"name\":\"value\",\"name\":\"value\"}'\n" +
			"\n" +
			"   Optionally provide a file containing service-specific configuration parameters in a valid JSON object.\n" +
			"   The path to the parameters file can be an absolute or relative path to a file:\n" +
			"\n" +
			"   sb create-service SERVICE PLAN SERVICE_INSTANCE -c PATH_TO_FILE\n" +
			"\n" +
			"ALIAS:\n" +
			"   cs\n" +
			"\n" +
			"OPTIONS:\n" +
			"   -c      Valid JSON object containing service-specific configuration parameters, provided either in-line or in a file. For a list of supported configuration parameters, see documentation for the particular service offering.\n" +
			"   -t      User provided tags\n"
	case "update-service":
		return "NAME:\n" +
			"   update-service - Update a service instance\n" +
			"\n" +
			"USAGE:\n" +
			"   sb update-service SERVICE_INSTANCE [-p NEW_PLAN] [-c PARAMETERS_AS_JSON] [-t TAGS]\n" +
			"\n" +
			"   Optionally provide service-specific configuration parameters in a valid JSON object in-line.\n" +
			"   sb update-service -c '{\"name\":\"value\",\"name\":\"value\"}'\n" +
			"\n" +
			"   Optionally provide a file containing service-specific configuration parameters in a valid JSON object. \n" +
			"   The path to the parameters file can be an absolute or relative path to a file.\n" +
			"   sb update-service -c PATH_TO_FILE\n" +
			"\n" +
			"   Optionally provide a list of comma-delimited tags that will be written to the VCAP_SERVICES environment variable for any bound applications.\n" +
			"\n" +
			"OPTIONS:\n" +
			"   -c      Valid JSON object containing service-specific configuration parameters, provided either in-line or in a file. For a list of supported configuration parameters, see documentation for the particular service offering.\n" +
			"   -p      Change service plan for a service instance\n" +
			"   -t      User provided tags\n"
	case "delete-service", "ds":
		return "NAME:\n" +
			"   delete-service - Delete a service instance\n" +
			"\n" +
			"USAGE:\n" +
			"   sb delete-service SERVICE_INSTANCE [-f]\n" +
			"\n" +
			"ALIAS:\n" +
			"   ds\n" +
			"\n" +
			"OPTIONS:\n" +
			"   -f      Force deletion without confirmation\n"
	case "create-service-key", "csk":
		return "NAME:\n" +
			"   create-service-key - Create key for a service instance\n" +
			"\n" +
			"USAGE:\n" +
			"   sb create-service-key SERVICE_INSTANCE SERVICE_KEY [-c PARAMETERS_AS_JSON]\n" +
			"\n" +
			"   Optionally provide service-specific configuration parameters in a valid JSON object in-line.\n" +
			"   sb create-service-key SERVICE_INSTANCE SERVICE_KEY -c '{\"name\":\"value\",\"name\":\"value\"}'\n" +
			"\n" +
			"   Optionally provide a file containing service-specific configuration parameters in a valid JSON object. The path to the parameters file can be an absolute or relative path to a file.\n" +
			"   sb create-service-key SERVICE_INSTANCE SERVICE_KEY -c PATH_TO_FILE\n" +
			"\n" +
			"ALIAS:\n" +
			"   csk\n" +
			"\n" +
			"OPTIONS:\n" +
			"   -c      Valid JSON object containing service-specific configuration parameters, provided either in-line or in a file. For a list of supported configuration parameters, see documentation for the particular service offering.\n" +
			"   -j      Print only JSON output\n"
	case "service-keys", "sk":
		return "NAME:\n" +
			"   service-keys - List keys for a service instance\n" +
			"\n" +
			"USAGE:\n" +
			"   sb service-keys SERVICE_INSTANCE\n" +
			"\n" +
			"ALIAS:\n" +
			"   sk\n"
	case "delete-service-key", "dsk":
		return "NAME:\n" +
			"   delete-service-key - Delete a service key\n" +
			"\n" +
			"USAGE:\n" +
			"   sb delete-service-key SERVICE_INSTANCE SERVICE_KEY [-f]\n" +
			"\n" +
			"ALIAS:\n" +
			"   dsk\n" +
			"\n" +
			"OPTIONS:\n" +
			"   -f      Force deletion without confirmation\n"
	case "auth":
		return "NAME:\n" +
			"   auth - Authenticate user non-interactively\n" +
			"\n" +
			"USAGE:\n" +
			"   sb auth USERNAME PASSWORD\n"
	case "version", "-v":
		return "NAME:\n" +
			"   version - Print the version\n" +
			"\n" +
			"USAGE:\n" +
			"   sb version\n" +
			"   sb -v\n"
	case "login", "l":
		return "NAME:\n" +
			"   login - Log user in\n" +
			"\n" +
			"USAGE:\n" +
			"   sb login\n" +
			"\n" +
			"ALIAS:\n" +
			"   l\n"
	case "logout", "lo":
		return "NAME:\n" +
			"   logout - Log user out\n" +
			"\n" +
			"USAGE:\n" +
			"   sb logout\n" +
			"\n" +
			"ALIAS:\n" +
			"   lo\n"
	case "target", "t":
		return "NAME:\n" +
			"   target - Set or view the target\n" +
			"\n" +
			"USAGE:\n" +
			"   sb target [-o ORG] [-s SPACE]\n" +
			"\n" +
			"ALIAS:\n" +
			"   t\n"
	case "marketplace", "m":
		return "NAME:\n" +
			"   marketplace - List available offerings in the marketplace\n" +
			"\n" +
			"USAGE:\n" +
			"   sb marketplace\n" +
			"\n" +
			"ALIAS:\n" +
			"   m\n"
	case "services", "s":
		return "NAME:\n" +
			"   services - List all service instances in the target space\n" +
			"\n" +
			"USAGE:\n" +
			"   sb services\n" +
			"\n" +
			"ALIAS:\n" +
			"   s\n"
	case "service":
		return "NAME:\n" +
			"   service - Show service instance info\n" +
			"\n" +
			"USAGE:\n" +
			"   sb service SERVICE_INSTANCE\n"
	case "help", "h":
		return "NAME:\n" +
			"   help - Show help\n" +
			"\n" +
			"USAGE:\n" +
			"   sb help [COMMAND]\n" +
			"\n" +
			"ALIAS:\n" +
			"   h\n"
	default:
		return "Sorry! No help available..."
	}
}
