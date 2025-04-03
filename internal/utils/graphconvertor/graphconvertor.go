package graphconvertor

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/firgavin/eino-devops/internal/utils/generic"
	"github.com/firgavin/eino-devops/model"
	"github.com/firgavin/eino-devops/model/d2"

	"oss.terrastruct.com/d2/d2graph"
	"oss.terrastruct.com/d2/d2layouts/d2elklayout"
	"oss.terrastruct.com/d2/d2lib"
	"oss.terrastruct.com/d2/d2renderers/d2svg"
	"oss.terrastruct.com/d2/d2themes/d2themescatalog"
	"oss.terrastruct.com/d2/lib/log"
	"oss.terrastruct.com/d2/lib/textmeasure"
	"oss.terrastruct.com/d2/lib/urlenc"
	"oss.terrastruct.com/util-go/go2"
)

// CanvasInfo is the top level structure from the input JSON
type CanvasInfo struct {
	Data struct {
		CanvasInfo d2.Graph `json:"canvas_info"`
	} `json:"data"`
}

// ConvertToD2 converts an Eino graph to D2 format
func ConvertToD2(r io.Reader, w io.Writer) error {
	// Read and parse JSON
	var canvasInfo CanvasInfo
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&canvasInfo); err != nil {
		return fmt.Errorf("error decoding JSON: %w", err)
	}

	// Write style classes
	if _, err := w.Write([]byte(styleClasses())); err != nil {
		return fmt.Errorf("error writing style classes: %w", err)
	}

	// Convert and write the graph
	if err := convertGraph(&canvasInfo.Data.CanvasInfo, "", w); err != nil {
		return fmt.Errorf("error converting graph: %w", err)
	}

	return nil
}

// ConvertModelToD2 converts a model.GraphSchema to D2 format
func ConvertModelToD2(graphSchema *model.GraphSchema, w io.Writer) error {
	// Write style classes
	if _, err := w.Write([]byte(styleClasses())); err != nil {
		return fmt.Errorf("error writing style classes: %w", err)
	}

	// Convert model.GraphSchema to our Graph structure
	graph := convertModelGraph(graphSchema)

	// Convert and write the graph
	if err := convertGraph(graph, "", w); err != nil {
		return fmt.Errorf("error converting graph: %w", err)
	}

	return nil
}

// convertModelGraph converts a model.GraphSchema to our Graph structure
func convertModelGraph(gs *model.GraphSchema) *d2.Graph {
	graph := &d2.Graph{
		Name:            gs.Name,
		Component:       string(gs.Component),
		NodeTriggerMode: string(gs.NodeTriggerMode),
	}

	// Convert nodes
	for _, modelNode := range gs.Nodes {
		node := d2.Node{
			Key:  modelNode.Key,
			Name: modelNode.Name,
			Type: string(modelNode.Type),
		}

		if modelNode.AllowOperate {
			node.AllowOperate = &modelNode.AllowOperate
		}

		if modelNode.GraphSchema != nil {
			node.GraphSchema = convertModelGraph(modelNode.GraphSchema)
		}

		graph.Nodes = append(graph.Nodes, node)
	}

	// Convert edges
	for _, modelEdge := range gs.Edges {
		edge := d2.Edge{
			ID:            modelEdge.ID,
			Name:          modelEdge.Name,
			SourceNodeKey: modelEdge.SourceNodeKey,
			TargetNodeKey: modelEdge.TargetNodeKey,
		}
		graph.Edges = append(graph.Edges, edge)
	}

	return graph
}

// convertGraph converts a Graph to D2 format and writes it to the writer
func convertGraph(graph *d2.Graph, parentID string, w io.Writer) error {
	indent := ""
	if parentID != "" {
		indent = "  "
		if _, err := fmt.Fprintf(w, "\"%s\": {\n", parentID); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(w, "%sdirection: down\n", indent); err != nil {
			return err
		}
	}

	// Process nodes
	for _, node := range graph.Nodes {
		if parentID != "" && (node.Key == "start" || node.Key == "end") {
			continue
		}

		name := node.Name
		if name == "" {
			name = node.Key
		}

		if _, err := fmt.Fprintf(w, "%s\"%s\": \"%s\" {class: %s}\n", indent, node.Key, name, node.Type); err != nil {
			return err
		}

		if node.GraphSchema != nil {
			if err := convertGraph(node.GraphSchema, node.Key, w); err != nil {
				return err
			}
		}
	}

	// Process edges
	for _, edge := range graph.Edges {
		if parentID != "" && (edge.SourceNodeKey == "start" || edge.TargetNodeKey == "end") {
			continue
		}

		if _, err := fmt.Fprintf(w, "%s\"%s\" -> \"%s\": {class: edge}\n",
			indent, edge.SourceNodeKey, edge.TargetNodeKey); err != nil {
			return err
		}
	}

	if parentID != "" {
		if _, err := fmt.Fprintln(w, "}"); err != nil {
			return err
		}
	}

	return nil
}

// styleClasses returns D2 style class definitions
func styleClasses() string {
	return `# Style definitions
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
direction: right
`
}

// GenerateD2FromJSON is a convenience function that takes a JSON string
// and returns a D2 diagram as a string
func GenerateD2FromJSON(jsonData string) (string, error) {
	r := strings.NewReader(jsonData)
	var b strings.Builder

	if err := ConvertToD2(r, &b); err != nil {
		return "", err
	}

	return b.String(), nil
}

// GenerateD2FromGraphSchema is a convenience function that takes a GraphSchema
// and returns a D2 diagram as a string
func GenerateD2FromGraphSchema(graphSchema *model.GraphSchema) (string, error) {
	var b strings.Builder

	if err := ConvertModelToD2(graphSchema, &b); err != nil {
		return "", err
	}

	return b.String(), nil
}

type SVGOption struct {
	leeching bool
}

type Option func(*SVGOption)

// Use this option will freeloads the api from d2 playground to generate the svg.
// Otherwise, we will use local generator to generate the svg.
func WithLeeching() Option {
	return func(o *SVGOption) {
		o.leeching = true
	}
}

var leechingSource = "aHR0cHM6Ly9hcGkuZDJsYW5nLmNvbS9yZW5kZXIvc3ZnP3NjcmlwdD0lcyZ0aGVtZT0wJnNrZXRjaD0wJmxheW91dD1lbGs="

func init() {
	b, _ := base64.StdEncoding.DecodeString(leechingSource)
	leechingSource = string(b)
}

func Dot2SVG(graphString string, opts ...Option) ([]byte, error) {
	var svgOption SVGOption
	for _, opt := range opts {
		opt(&svgOption)
	}

	if svgOption.leeching {
		script, err := urlenc.Encode(graphString)
		if err != nil {
			return nil, err
		}
		ls := fmt.Sprintf(leechingSource, script)
		resp, err := http.Get(ls)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		return io.ReadAll(resp.Body)
	}

	ruler, _ := textmeasure.NewRuler()
	layoutResolver := func(engine string) (d2graph.LayoutGraph, error) {
		return d2elklayout.DefaultLayout, nil
	}
	renderOpts := &d2svg.RenderOpts{
		Pad:     go2.Pointer(int64(5)),
		Center:  generic.PtrOf(true),
		ThemeID: &d2themescatalog.NeutralDefault.ID,
	}
	compileOpts := &d2lib.CompileOptions{
		LayoutResolver: layoutResolver,
		Ruler:          ruler,
	}
	ctx := log.WithDefault(context.Background())
	diagram, _, _ := d2lib.Compile(ctx, graphString, compileOpts, renderOpts)
	return d2svg.Render(diagram, renderOpts)
}
