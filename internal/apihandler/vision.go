package apihandler

import (
	"net/http"
	"text/template"
)

const visionPageTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Graph Visualization</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 0;
            padding: 20px;
        }
        #selector {
            margin-bottom: 20px;
        }
        select {
            padding: 8px;
            font-size: 16px;
            min-width: 200px;
        }
        #visualization {
            width: 100%;
            height: 85vh;
            border: 1px solid #ccc;
            overflow: auto;
        }
        #loading {
            display: none;
            margin: 20px 0;
        }
    </style>
</head>
<body>
    <div id="selector">
        <label for="graph-select">Select a graph: </label>
        <select id="graph-select">
            <option value="">-- Select a graph --</option>
        </select>
        <span id="loading">Loading...</span>
    </div>
    <div id="visualization"></div>

    <script>
        // Fetch available graphs when the page loads
        document.addEventListener('DOMContentLoaded', async () => {
            try {
                const response = await fetch('/eino/devops/debug/v1/graphs');
                const data = await response.json();
                const select = document.getElementById('graph-select');
                
                if (data.data && data.data.graphs) {
                    data.data.graphs.forEach(graph => {
                        const option = document.createElement('option');
                        option.value = graph.id;
                        option.textContent = graph.name ? graph.name + ' (' + graph.id + ')' : graph.id;
                        select.appendChild(option);
                    });
                }
            } catch (error) {
                console.error('Error fetching graphs:', error);
                document.getElementById('visualization').innerHTML = 
                    '<p style="color: red">Error loading graphs. Please try again later.</p>';
            }
        });

        // Update visualization when a graph is selected
        document.getElementById('graph-select').addEventListener('change', async (e) => {
            const graphId = e.target.value;
            const visDiv = document.getElementById('visualization');
            const loading = document.getElementById('loading');
            
            if (!graphId) {
                visDiv.innerHTML = '<p>Please select a graph to visualize.</p>';
                return;
            }
            
            try {
                loading.style.display = 'inline';
                visDiv.innerHTML = '';
                
                const response = await fetch('/eino/devops/debug/v1/graphs/' + graphId + '/vision');
                const svgText = await response.text();
                
                visDiv.innerHTML = svgText;
                
                // Adjust SVG to fit container
                const svg = visDiv.querySelector('svg');
                if (svg) {
                    svg.setAttribute('width', '100%');
                    svg.setAttribute('height', '100%');
                    svg.style.maxHeight = '85vh';
                }
            } catch (error) {
                console.error('Error fetching visualization:', error);
                visDiv.innerHTML = 
                    '<p style="color: red">Error loading visualization. Please try again later.</p>';
            } finally {
                loading.style.display = 'none';
            }
        });
    </script>
</body>
</html>
`

// GetVisionPage renders a page with a dropdown to select graphs and view their visualization
func GetVisionPage(res http.ResponseWriter, req *http.Request) {
	// Render the HTML template with dropdown and visualization container
	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl := template.Must(template.New("vision_page").Parse(visionPageTemplate))
	tmpl.Execute(res, nil)
}
