package apihandler

import (
	"fmt"

	devmodel "github.com/firgavin/eino-devops/model"
)

func flattenGraph(graph *devmodel.GraphSchema) (allNodes []*devmodel.Node, allEdges []*devmodel.Edge) {
	if graph == nil {
		return nil, nil
	}

	for _, node := range graph.Nodes {
		if node.GraphSchema != nil {
			subGraph := node.GraphSchema
			// 首先更新prefix node name，避免冲突
			for _, n := range subGraph.Nodes {
				n.Key = fmt.Sprintf("%s-%s", node.Key, n.Key)
			}
			for _, e := range subGraph.Edges {
				e.SourceNodeKey = fmt.Sprintf("%s-%s", node.Key, e.SourceNodeKey)
				e.TargetNodeKey = fmt.Sprintf("%s-%s", node.Key, e.TargetNodeKey)
			}

			subNodes, subEdges := flattenGraph(subGraph)

			// parent graph移除原始节点，并且更新edges
			var startNode, endNode *devmodel.Node
			for _, n := range subNodes {
				switch n.Type {
				case devmodel.NodeTypeOfStart:
					startNode = n
				case devmodel.NodeTypeOfEnd:
					endNode = n
				}
			}
			for _, edge := range graph.Edges {
				if edge.TargetNodeKey == node.Key {
					edge.TargetNodeKey = startNode.Key
				}
				if edge.SourceNodeKey == node.Key {
					edge.SourceNodeKey = endNode.Key
				}
			}
			allNodes = append(allNodes, subNodes...)
			allEdges = append(allEdges, subEdges...)
		} else {
			allNodes = append(allNodes, node)
		}
	}
	allEdges = append(allEdges, graph.Edges...)

	return deduplicate(allNodes, allEdges)
}

func deduplicate(nodes []*devmodel.Node, edges []*devmodel.Edge) ([]*devmodel.Node, []*devmodel.Edge) {
	nodeMap := make(map[string]*devmodel.Node)
	edgeMap := make(map[string]*devmodel.Edge)

	for _, n := range nodes {
		nodeMap[n.Key] = n
	}
	for _, e := range edges {
		key := fmt.Sprintf("%s->%s", e.SourceNodeKey, e.TargetNodeKey)
		edgeMap[key] = e
	}

	uniqueNodes := make([]*devmodel.Node, 0, len(nodeMap))
	for _, v := range nodeMap {
		uniqueNodes = append(uniqueNodes, v)
	}
	uniqueEdges := make([]*devmodel.Edge, 0, len(edgeMap))
	for _, v := range edgeMap {
		uniqueEdges = append(uniqueEdges, v)
	}

	return uniqueNodes, uniqueEdges
}
