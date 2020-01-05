package flexvolume

// Status is a status
type Status string

const (
	// StatusSuccess is a Status meaning success
	StatusSuccess Status = "Success"
	// StatusFailure is a Status meaning failure
	StatusFailure Status = "Failure"
	// StatusNotSupported is a Status meaning something isn't supported
	StatusNotSupported Status = "Not Supported"
)

// FlexVolume is
type FlexVolume interface {
	Capabilities() Capabilities
	Init() Response
	Attach(map[string]string) Response
	Detach(string) Response
	Mount(string, string, map[string]string) Response
	Unmount(string) Response
}

// Response is
type Response struct {
	Status  Status `json:"status"`
	Message string `json:"message"`
	Device  string `json:"device,omitempty"`
}

// Capabilities is a list of capabilities
type Capabilities struct {
	Attach bool `json:"attach"`
	Detach bool `json:"detach"`
}
