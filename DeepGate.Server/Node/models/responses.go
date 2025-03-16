package models

// NodeSortedHostResponse represents a sorted host response with IP and name
type NodeSortedHostResponse struct {
	HostIP   string `json:"hostIp"`
	HostName string `json:"hostName"`
}
