package manifest

import "time"

type Manifest struct {
    Root      string      `json:"root"`
    Generated time.Time   `json:"generated_at"`
    Nodes     []*Node     `json:"nodes"`
}


