package data

// FlowLog defines the data schema for the flow-log message
// Note: normally snake case is not a Go convention. But mongodb decoder
// uses strict binding, i.e. src_app won't decode to SrcApp.
// For simplicity sake, this demo attempts to keep a consistent schema
// instead of creating a separate conventional struct.
type FlowLog struct {
	Src_App  string `json:"src_app" validate:"required"`
	Dest_App string `json:"dest_app" validate:"required"`
	Vpc_Id   string `json:"vpc_id" validate:"required"`
	Bytes_Tx int    `json:"bytes_tx" validate:"required"`
	Bytes_Rx int    `json:"bytes_rx" validate:"required"`
	Hour    int    `json:"hour" validate:"required"`
}
