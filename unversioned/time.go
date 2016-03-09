package unversioned

import (
	"time"
)
// Time is a wrapper around time.Time which supports correct
// marshaling to YAML and JSON.  Wrappers are provided for many
// of the factory methods that the time package offers.
//
// +protobuf.options.marshal=false
type Time struct {
	time.Time `protobuf:"Timestamp,1,req,name=time"`
}