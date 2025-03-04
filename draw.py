import json
from typing import Optional, List
from dataclasses import dataclass


@dataclass
class Edge:
    id: str
    name: str
    source_node_key: str
    target_node_key: str


@dataclass
class Node:
    key: str
    name: str
    type: str
    allow_operate: Optional[bool] = None
    graph_schema: Optional["Graph"] = None


@dataclass
class Graph:
    nodes: List[Node]
    edges: List[Edge]
    name: str = ""
    component: str = ""
    node_trigger_mode: str = ""

    @classmethod
    def from_dict(cls, data: dict):
        # Convert edges
        edges = [
            Edge(
                id=e["id"],
                name=e["name"],
                source_node_key=e["source_node_key"],
                target_node_key=e["target_node_key"],
            )
            for e in data.get("edges", [])
        ]

        # Convert nodes
        nodes = []
        for n in data.get("nodes", []):
            node = Node(
                key=n["key"],
                name=n["name"],
                type=n["type"],
                allow_operate=n.get("allow_operate"),
            )
            if n.get("graph_schema"):
                node.graph_schema = cls.from_dict(n["graph_schema"])
            nodes.append(node)

        # Create graph
        return Graph(
            name=data.get("name", ""),
            component=data.get("component", ""),
            nodes=nodes,
            edges=edges,
            node_trigger_mode=data.get("node_trigger_mode", ""),
        )


def print_style_classes():
    """Print D2 style class definitions"""
    print(
        """# Style definitions
classes: {
  Lambda: {
    shape: hexagon  # 六边形代表函数/计算单元
    style: {
      fill: "#8B4513"  # 咖啡色，代表函数执行
      stroke: "#654321"
      font-color: "#ffffff"
      shadow: true
      stroke-width: 2
      border-radius: 4  # 直角一些，表示计算性质
      font-size: 14
      3d: true  # 表示实际执行的节点
    }
  }
  start: {
    shape: circle  # 圆形是流程图中标准的开始节点形状
    style: {
      fill: "#00B140"  # 绿色，表示起点
      stroke: "#008631"
      font-color: "#ffffff"
      shadow: true
      stroke-width: 2
      border-radius: 20  # 圆形，表示起始点
      font-size: 14
    }
  }
  end: {
    shape: circle  # 圆形是流程图中标准的结束节点形状
    style: {
      fill: "#DC143C"  # 深红色，表示终点
      stroke: "#B22222"
      font-color: "#ffffff"
      shadow: true
      stroke-width: 2
      border-radius: 20  # 圆形，表示结束点
      font-size: 14
    }
  }
  ChatModel: {
    shape: cloud  # 云形状代表AI服务/模型
    style: {
      fill: "#1E90FF"  # 天蓝色，表示AI模型
      stroke: "#0066CC"
      font-color: "#ffffff"
      shadow: true
      stroke-width: 3  # 加粗边框，强调模型的重要性
      border-radius: 10
      font-size: 14
      3d: true  # 表示实际执行的节点
    }
  }
  ChatTemplate: {
    shape: document  # 文档形状代表模板/文本处理
    style: {
      fill: "#FFB6C1"  # 浅粉色，表示模板
      stroke: "#FF69B4"
      font-color: "#4A4A4A"  # 深色文字，提高可读性
      shadow: true
      stroke-width: 1  # 细边框，表示轻量级处理
      border-radius: 8
      font-size: 14
    }
  }
  Passthrough: {
    shape: queue  # 队列形状表示数据传递
    style: {
      fill: "#E6E6FA"  # 淡紫色，表示数据传递
      stroke: "#9370DB"
      font-color: "#4A4A4A"
      shadow: true
      stroke-width: 1
      border-radius: 10
      font-size: 14
      opacity: 0.85  # 半透明，表示传递性质
    }
  }
  Branch: {
    shape: diamond  # 菱形是流程图中标准的判断节点形状
    style: {
      fill: "#FFD700"  # 金色，表示决策点
      stroke: "#DAA520"
      font-color: "#2c3e50"
      shadow: true
      stroke-width: 2
      border-radius: 2  # 尖锐的边角，强调判断特性
      font-size: 14
    }
  }
  Parallel: {
    shape: rectangle  # 矩形代表并行处理块
    style: {
      fill: "#9370DB"  # 紫色，表示并行处理
      stroke: "#7B68EE"
      font-color: "#ffffff"
      shadow: true
      stroke-width: 2
      border-radius: 2
      font-size: 14
      multiple: true  # 表示多重执行
      opacity: 0.9
    }
  }
  edge: {
    style: {
      animated: true  # 动画效果表示流动
      stroke: "#666666"  # 中性灰色
      stroke-width: 1.5
      font-size: 12
      opacity: 0.9
      stroke-dash: 0
    }
  }
  container: {
    shape: package  # 包装形状代表容器/组合
    style: {
      fill: "#F8F8FF"  # 非常浅的背景色
      stroke: "#E6E6FA"
      font-color: "#4A4A4A"
      shadow: true
      stroke-width: 1  # 细边框，不喧宾夺主
      border-radius: 12
      font-size: 14
      opacity: 0.95
    }
  }
}
"""
    )


def convert_d2(graph: Graph, parent_id: str = ""):
    indent = ""
    if parent_id:
        indent = "  "
        print(f'"{parent_id}": {{')
        print(f"{indent}direction: down")

    # First process nodes
    for node in graph.nodes:
        if parent_id and (node.type == "start" or node.type == "end"):
            continue
        name = node.name or node.key
        key = node.key

        print(f'{indent}"{key}": "{name}" {{class: {node.type}}}')
        if node.graph_schema:
            convert_d2(node.graph_schema, node.key)

    # Then process edges
    for edge in graph.edges:
        if parent_id and (
            edge.source_node_key == "start" or edge.target_node_key == "end"
        ):
            continue
        print(
            f'{indent}"{edge.source_node_key}" -> "{edge.target_node_key}": '
            + '{class: "edge"}'
        )

    if parent_id:
        print("}")


if __name__ == "__main__":
    # Read and parse JSON
    with open("data3.json", "r") as f:
        o = json.load(f)
        o = o["data"]["canvas_info"]

    graph = Graph.from_dict(o)

    print_style_classes()
    convert_d2(graph)