package store

// Flow represents a single data point of network flow
type Flow struct {
	// Src is the source application name
	Src string `json:"src_app"`
	// Dst is the destination application name
	Dst     string `json:"dst_app"`
	VpcID   string `json:"vpc_id"`
	BytesTx int    `json:"bytes_tx"`
	BytsRx  int    `json:"bytes_rx"`
	Hour    int    `json:"hour"`
}
