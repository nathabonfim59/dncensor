package stack

type StackType string

const (
	StackSystemdResolved StackType = "systemd-resolved"
	StackNetworkManager  StackType = "networkmanager"
	StackResolvConf      StackType = "resolvconf"
)

type Stack interface {
	Type() StackType
	Detect() bool
	CurrentDNS() (string, error)
	SetDNS(primary, secondary string) error
	SetDOH(endpoint string) error
	CaptureDNS() ([]byte, error)
	ApplyDNS(content []byte) error
	RequiresRoot() bool
}

func Detect() Stack {
	candidates := []Stack{
		&systemdResolvedStack{},
		&nmStack{},
		&resolvConfStack{},
	}
	for _, s := range candidates {
		if s.Detect() {
			return s
		}
	}
	return nil
}
