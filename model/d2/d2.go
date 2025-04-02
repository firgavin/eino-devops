package d2

// Edge represents a connection between nodes in the graph
type Edge struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	SourceNodeKey string `json:"source_node_key"`
	TargetNodeKey string `json:"target_node_key"`
}

// Node represents a node in the graph
type Node struct {
	Key          string `json:"key"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	AllowOperate *bool  `json:"allow_operate,omitempty"`
	GraphSchema  *Graph `json:"graph_schema,omitempty"`
}

// Graph represents a graph structure with nodes and edges
type Graph struct {
	Nodes           []Node `json:"nodes"`
	Edges           []Edge `json:"edges"`
	Name            string `json:"name,omitempty"`
	Component       string `json:"component,omitempty"`
	NodeTriggerMode string `json:"node_trigger_mode,omitempty"`
}
