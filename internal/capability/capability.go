package capability

// Capabilities holds all the requested permissions, categorized by area.
type Capabilities struct {
	Filesystem []FileSystemCapability `yaml:"filesystem,omitempty"`
	Network    *NetworkCapability     `yaml:"network,omitempty"`
	Process    []ProcessCapability    `yaml:"process,omitempty"`
}

// FileSystemCapability defines permissions for a specific path.
type FileSystemCapability struct {
	Path        string   `yaml:"path"`
	Permissions []string `yaml:"permissions"`
	Recursive   bool     `yaml:"recursive,omitempty"`
}

// NetworkCapability now uses two distinct slice types
type NetworkCapability struct {
	Egress  []EgressRule  `yaml:"egress,omitempty"`
	Ingress []IngressRule `yaml:"ingress,omitempty"`
}

// EgressRule includes the Hosts field
type EgressRule struct {
	Protocol string   `yaml:"protocol"`
	Hosts    []string `yaml:"hosts"`
	Ports    []int    `yaml:"ports"`
}

// IngressRule correctly omits the Hosts field
type IngressRule struct {
	Protocol       string   `yaml:"protocol"`
	Ports          []int    `yaml:"ports"`
	AllowedOrigins []string `yaml:"allowed_origins,omitempty"`
}

// ProcessCapability defines a process-related permission.
// Using pointers allows us to easily determine which type of rule it is.
type ProcessCapability struct {
	Exec   *ExecRule `yaml:"exec,omitempty"`
	Kill   []string  `yaml:"kill,omitempty"`
	List   []string  `yaml:"list,omitempty"`
	Signal []string  `yaml:"signal,omitempty"`
}

// ExecRule defines the constraints for executing a command.
type ExecRule struct {
	Command string   `yaml:"command"`
	Args    []string `yaml:"args,omitempty"`
}
