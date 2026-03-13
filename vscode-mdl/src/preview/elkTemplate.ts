// SPDX-License-Identifier: Apache-2.0

import { escapeHtml, escapeForTemplate } from './mxcliRunner';
import { sketchHelpersJs } from './renderers/sketchHelpers';
import { microflowRendererJs } from './renderers/microflowRenderer';
import { domainModelRendererJs } from './renderers/domainModelRenderer';
import { moduleOverviewRendererJs } from './renderers/moduleOverviewRenderer';
import { queryPlanRendererJs } from './renderers/queryPlanRenderer';
import { pageWireframeRendererJs } from './renderers/pageWireframeRenderer';

/**
 * Returns the complete HTML content for the ELK.js-based webview panel.
 * Composes all renderer JS code strings into a single HTML document.
 */
export function getElkWebviewContent(jsonData: string, title: string, theme: string = 'clean'): string {
	const escapedJson = escapeForTemplate(jsonData);

	return `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>${escapeHtml(title)}</title>
	${theme === 'sketch' ? '<link href="https://fonts.googleapis.com/css2?family=Architects+Daughter&display=swap" rel="stylesheet">' : ''}
	<style>
		body {
			margin: 0;
			padding: 0;
			background-color: var(--vscode-editor-background, #1e1e1e);
			color: var(--vscode-editor-foreground, #d4d4d4);
			font-family: ${theme === 'sketch' ? "'Architects Daughter', cursive" : "var(--vscode-font-family, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif)"};
			overflow: hidden;
			height: 100vh;
			display: flex;
			flex-direction: column;
		}
		.toolbar {
			display: flex;
			gap: 8px;
			padding: 8px 16px;
			align-items: center;
			border-bottom: 1px solid var(--vscode-panel-border, #333);
			flex-shrink: 0;
		}
		.toolbar button {
			background: var(--vscode-button-background, #0e639c);
			color: var(--vscode-button-foreground, #fff);
			border: none;
			padding: 4px 12px;
			cursor: pointer;
			font-size: 12px;
			border-radius: 2px;
			font-family: ${theme === 'sketch' ? "'Architects Daughter', cursive" : "inherit"};
		}
		.toolbar button:hover {
			background: var(--vscode-button-hoverBackground, #1177bb);
		}
		.toolbar .title {
			flex: 1;
			font-size: 13px;
			font-weight: 600;
		}
		.toolbar .zoom-info {
			font-size: 11px;
			opacity: 0.7;
			min-width: 40px;
			text-align: center;
		}
		.toolbar button.active {
			background: var(--vscode-button-secondaryBackground, #3a3d41);
			outline: 1px solid var(--vscode-focusBorder, #007acc);
		}
		.legend {
			display: flex;
			gap: 12px;
			padding: 6px 16px;
			font-size: 11px;
			border-bottom: 1px solid var(--vscode-panel-border, #333);
			flex-shrink: 0;
		}
		.legend-item {
			display: flex;
			align-items: center;
			gap: 4px;
		}
		.legend-swatch {
			width: 12px;
			height: 12px;
			border-radius: 2px;
		}
		#diagram-viewport {
			flex: 1;
			overflow: hidden;
			cursor: grab;
			position: relative;
		}
		#diagram-viewport.dragging {
			cursor: grabbing;
		}
		#diagram-canvas {
			transform-origin: 0 0;
			position: absolute;
			top: 0;
			left: 0;
		}
		#diagram-canvas svg {
			display: block;
		}
		.source-view {
			display: none;
			flex: 1;
			background: var(--vscode-textCodeBlock-background, #2d2d2d);
			padding: 12px 16px;
			font-family: var(--vscode-editor-fontFamily, monospace);
			font-size: 12px;
			white-space: pre-wrap;
			overflow: auto;
		}
		#loading {
			display: flex;
			align-items: center;
			justify-content: center;
			flex: 1;
			font-size: 14px;
			opacity: 0.6;
		}
		#node-popover {
			display: none;
			position: fixed;
			background: var(--vscode-editorWidget-background, #252526);
			border: 1px solid var(--vscode-editorWidget-border, #454545);
			border-radius: 4px;
			padding: 8px 12px;
			font-size: 12px;
			font-family: var(--vscode-editor-fontFamily, monospace);
			line-height: 1.5;
			max-width: 400px;
			z-index: 100;
			box-shadow: 0 2px 8px rgba(0,0,0,0.4);
			white-space: pre-wrap;
			pointer-events: auto;
		}
		#node-popover .popover-title {
			font-weight: 600;
			margin-bottom: 4px;
			padding-bottom: 4px;
			border-bottom: 1px solid var(--vscode-panel-border, #333);
		}
		#node-popover .popover-line {
			opacity: 0.85;
		}
		.mf-node.highlighted path, .entity-node.highlighted path {
			stroke-width: 3 !important;
			opacity: 1 !important;
		}
		.mf-node.highlighted circle {
			stroke-width: 3 !important;
			opacity: 1 !important;
		}
		.mf-node.dimmed, .entity-node.dimmed, .wf-node.dimmed {
			opacity: 0.3;
			transition: opacity 0.15s;
		}
		.wf-node.highlighted path {
			stroke-width: 3 !important;
			opacity: 1 !important;
		}
	</style>
</head>
<body>
	<div class="toolbar">
		<span class="title">${escapeHtml(title)}</span>
		<span class="zoom-info" id="zoom-info">100%</span>
		<button onclick="resetView()">Fit</button>
		<button onclick="zoomIn()">+</button>
		<button onclick="zoomOut()">&minus;</button>
		<button onclick="toggleSource()">Source</button>
		<button onclick="copySource()">Copy</button>
		<button id="btn-collapse" onclick="toggleCollapseAll()" style="display:none">Collapse All</button>
	</div>
	<div class="legend" id="legend"></div>
	<div id="loading">Loading ELK layout...</div>
	<div id="diagram-viewport" style="display:none">
		<div id="diagram-canvas"></div>
	</div>
	<div id="node-popover"></div>
	<pre id="source-view" class="source-view"></pre>

	<script>
		const vscodeApi = acquireVsCodeApi();
		const rawJson = \`${escapedJson}\`;
		const data = JSON.parse(rawJson);
		const diagramType = data.type || 'module-overview';
		window.diagramData = data;

		// Populate source view
		document.getElementById('source-view').textContent = JSON.stringify(data, null, 2);

		// Set up legend based on diagram type
		var legendEl = document.getElementById('legend');
		if (diagramType === 'microflow') {
			legendEl.innerHTML =
				'<div class="legend-item"><div class="legend-swatch" style="background:#5ba55b; border: 1px solid #3d7a3d"></div>Event</div>' +
				'<div class="legend-item"><div class="legend-swatch" style="background:#4a90d9; border: 1px solid #2c6fad"></div>Object</div>' +
				'<div class="legend-item"><div class="legend-swatch" style="background:#3db8b8; border: 1px solid #2a9494"></div>Retrieve</div>' +
				'<div class="legend-item"><div class="legend-swatch" style="background:#8e6cbf; border: 1px solid #6b4f94"></div>Call</div>' +
				'<div class="legend-item"><div class="legend-swatch" style="background:#e89b3e; border: 1px solid #c07e2a"></div>Navigation</div>' +
				'<div class="legend-item"><div class="legend-swatch" style="background:#7f8c8d; border: 1px solid #636e70"></div>Variable</div>' +
				'<div class="legend-item"><div class="legend-swatch" style="background:#e8c93e; border: 1px solid #c4a82a"></div>Decision</div>' +
				'<div class="legend-item"><div class="legend-swatch" style="background:#d9534f; border: 1px solid #b52b27"></div>Validation</div>';
		} else if (diagramType === 'oql-queryplan') {
			legendEl.innerHTML =
				'<div class="legend-item"><div class="legend-swatch" style="background:#5ba55b; border: 1px solid #3d7a3d"></div>Table</div>' +
				'<div class="legend-item"><div class="legend-swatch" style="background:#8e6cbf; border: 1px solid #6b4f94"></div>Subquery</div>' +
				'<div class="legend-item"><div class="legend-swatch" style="background:#e89b3e; border: 1px solid #c07e2a"></div>Join</div>' +
				'<div class="legend-item"><div class="legend-swatch" style="background:#d9534f; border: 1px solid #b52b27"></div>Filter</div>' +
				'<div class="legend-item"><div class="legend-swatch" style="background:#4a90d9; border: 1px solid #2c6fad"></div>Aggregate</div>';
		} else if (diagramType === 'domainmodel') {
			legendEl.innerHTML =
				'<div class="legend-item"><div class="legend-swatch" style="background:#4a90d9; border: 1px solid #2c6fad"></div>Persistent</div>' +
				'<div class="legend-item"><div class="legend-swatch" style="background:#e89b3e; border: 1px solid #c07e2a"></div>Non-persistent</div>' +
				'<div class="legend-item"><div class="legend-swatch" style="background:#8e6cbf; border: 1px solid #6b4f94"></div>External</div>' +
				'<div class="legend-item"><div class="legend-swatch" style="background:#5ba55b; border: 1px solid #3d7a3d"></div>View</div>';
		} else if (diagramType === 'page') {
			legendEl.innerHTML =
				'<div class="legend-item"><div class="legend-swatch" style="background:#4a90d9; border: 1px solid #2c6fad"></div>Structure</div>' +
				'<div class="legend-item"><div class="legend-swatch" style="background:#5ba55b; border: 1px solid #3d7a3d"></div>Input</div>' +
				'<div class="legend-item"><div class="legend-swatch" style="background:#e89b3e; border: 1px solid #c07e2a"></div>Data</div>' +
				'<div class="legend-item"><div class="legend-swatch" style="background:#8e6cbf; border: 1px solid #6b4f94"></div>Action</div>';
		} else {
			legendEl.innerHTML =
				'<div class="legend-item"><div class="legend-swatch" style="background:#4a90d9; border: 1px solid #2c6fad"></div>App Module</div>' +
				'<div class="legend-item"><div class="legend-swatch" style="background:#8e6cbf; border: 1px solid #6b4f94"></div>System Module</div>';
		}

		const diagramTheme = '${theme}';
		${sketchHelpersJs()}
	</script>
	<script type="module">
		const loading = document.getElementById('loading');
		const d = window.diagramData;

		if (d.type === 'page') {
			// Pages use direct wireframe rendering — no ELK needed
			loading.textContent = 'Rendering wireframe...';
			window.renderSvg(null);
			loading.style.display = 'none';
			document.getElementById('diagram-viewport').style.display = '';
			setTimeout(window.resetView, 50);
		} else {
		loading.textContent = 'Importing ELK.js...';

		try {
			const elkModule = await import('https://cdn.jsdelivr.net/npm/elkjs@0.9.3/+esm');
			loading.textContent = 'Computing layout...';
			const ELK = elkModule.default || elkModule;

			var elkGraph;
			if (d.type === 'oql-queryplan') {
				// OQL Query Plan: tables + joins + result node, flowing left to right
				var tables = d.tables || [];
				var joins = d.joins || [];
				var children = [];
				var edges = [];

				// Table nodes (left layer)
				tables.forEach(function(t) {
					children.push({ id: t.id, width: t.width, height: t.height, labels: [{ text: t.entity }] });
				});

				// Join nodes (middle layer)
				joins.forEach(function(j) {
					children.push({ id: j.id, width: j.width, height: j.height, labels: [{ text: j.joinType }] });
				});

				// Result node (right side): the view entity itself
				var resultName = d.entityName || 'Result';
				var resultCols = (d.columns || []).length;
				var resultWidth = Math.max(resultName.length * 7.5 + 24, 100);
				var resultHeight = 28 + Math.max(resultCols, 1) * 18;
				children.push({ id: 'result', width: resultWidth, height: resultHeight, labels: [{ text: resultName }] });

				// Separate scalar subquery tables from regular tables
				var scalarTables = tables.filter(function(t) { return t.joinType === 'scalar'; });
				var regularTables = tables.filter(function(t) { return t.joinType !== 'scalar'; });

				if (joins.length > 0) {
					// Edges: each table pair connects to its join node, join connects to next join or result
					joins.forEach(function(j) {
						edges.push({ id: j.id + '-left', sources: [j.leftId], targets: [j.id] });
						edges.push({ id: j.id + '-right', sources: [j.rightId], targets: [j.id] });
					});
					// Last join connects to result
					edges.push({ id: 'join-result', sources: [joins[joins.length - 1].id], targets: ['result'] });
				} else if (regularTables.length > 0) {
					// No joins: single table flows directly to result
					edges.push({ id: 'direct', sources: [regularTables[0].id], targets: ['result'] });
				}

				// Scalar subquery tables connect directly to result
				scalarTables.forEach(function(st, idx) {
					edges.push({ id: 'scalar-' + idx, sources: [st.id], targets: ['result'] });
				});

				elkGraph = {
					id: 'root',
					layoutOptions: {
						'elk.algorithm': 'layered',
						'elk.direction': 'RIGHT',
						'elk.spacing.nodeNode': '40',
						'elk.layered.spacing.nodeNodeBetweenLayers': '80',
						'elk.edgeRouting': 'ORTHOGONAL',
						'elk.layered.crossingMinimization.strategy': 'LAYER_SWEEP',
					},
					children: children,
					edges: edges,
				};
			} else if (d.type === 'microflow') {
				// Microflow: nodes as ELK children, flows as edges
				var mfNodes = d.nodes || [];
				var mfEdges = d.edges || [];

				// Recursively build ELK node, handling compound nodes (loops)
				function buildElkNode(n) {
					var elkNode = { id: n.id, labels: [{ text: n.label }] };
					if (n.children && n.children.length > 0) {
						// Compound node (loop) — let ELK compute size from children
						elkNode.children = n.children.map(buildElkNode);
						elkNode.edges = (n.edges || []).map(function(e) {
							return { id: e.id, sources: [e.sourceId], targets: [e.targetId] };
						});
						elkNode.layoutOptions = {
							'elk.algorithm': 'layered',
							'elk.direction': 'RIGHT',
							'elk.spacing.nodeNode': '20',
							'elk.layered.spacing.nodeNodeBetweenLayers': '40',
							'elk.edgeRouting': 'ORTHOGONAL',
							'elk.padding': '[top=40,left=12,bottom=12,right=12]',
						};
					} else {
						elkNode.width = n.width;
						elkNode.height = n.height;
					}
					return elkNode;
				}

				elkGraph = {
					id: 'root',
					layoutOptions: {
						'elk.algorithm': 'layered',
						'elk.direction': 'RIGHT',
						'elk.spacing.nodeNode': '30',
						'elk.layered.spacing.nodeNodeBetweenLayers': '60',
						'elk.edgeRouting': 'ORTHOGONAL',
						'elk.layered.crossingMinimization.strategy': 'LAYER_SWEEP',
					},
					children: mfNodes.map(buildElkNode),
					edges: mfEdges.map(function(e) {
						return { id: e.id, sources: [e.sourceId], targets: [e.targetId] };
					}),
				};
			} else if (d.type === 'domainmodel') {
				// Domain model: entities as nodes, associations + generalizations as edges
				var entities = d.entities || [];
				var assocs = d.associations || [];
				var gens = d.generalizations || [];

				elkGraph = {
					id: 'root',
					layoutOptions: {
						'elk.algorithm': 'layered',
						'elk.direction': 'DOWN',
						'elk.spacing.nodeNode': '30',
						'elk.layered.spacing.nodeNodeBetweenLayers': '60',
						'elk.edgeRouting': 'ORTHOGONAL',
						'elk.layered.crossingMinimization.strategy': 'LAYER_SWEEP',
					},
					children: entities.map(function(e) {
						return { id: e.id, width: e.width, height: e.height, labels: [{ text: e.name }] };
					}),
					edges: assocs.map(function(a) {
						return { id: a.id, sources: [a.sourceId], targets: [a.targetId] };
					}).concat(gens.map(function(g, i) {
						return { id: 'gen-' + i, sources: [g.childId], targets: [g.parentId] };
					})),
				};
			} else {
				// Module overview
				d.modules = d.modules || [];
				d.edges = d.edges || [];
				elkGraph = {
					id: 'root',
					layoutOptions: {
						'elk.algorithm': 'layered',
						'elk.direction': 'DOWN',
						'elk.spacing.nodeNode': '40',
						'elk.layered.spacing.nodeNodeBetweenLayers': '100',
						'elk.edgeRouting': 'ORTHOGONAL',
					},
					children: d.modules.map(function(m) {
						return { id: m.id, width: calcNodeWidth(m), height: nodeHeight, labels: [{ text: m.name }] };
					}),
					edges: d.edges.map(function(e, i) {
						return { id: 'e' + i, sources: [e.source], targets: [e.target] };
					}),
				};
			}

			const elk = new ELK();
			window.elkInstance = elk;
			window.elkGraph = elkGraph;
			const layoutResult = await elk.layout(elkGraph);
			window.renderSvg(layoutResult);
			loading.style.display = 'none';
			document.getElementById('diagram-viewport').style.display = '';
			if (diagramType === 'domainmodel' || diagramType === 'microflow') {
				document.getElementById('btn-collapse').style.display = '';
			}
			setTimeout(window.resetView, 50);
		} catch(err) {
			loading.textContent = 'ELK error: ' + err.message;
		}
		} // end of else (non-page ELK path)
	</script>
	<script>

		// --- Build lookups ---
		var moduleMap = {};
		if (data.modules) {
			data.modules.forEach(function(m) { moduleMap[m.id] = m; });
		}
		var entityMap = {};
		if (data.entities) {
			data.entities.forEach(function(e) { entityMap[e.id] = e; });
		}
		var collapsedEntities = {};
		var collapsedActivities = {};
		var allCollapsed = false;
		var mfNodeMap = {};
		function indexMfNodes(nodes) {
			if (!nodes) return;
			nodes.forEach(function(n) {
				mfNodeMap[n.id] = n;
				if (n.children) indexMfNodes(n.children);
			});
		}
		indexMfNodes(data.nodes);
		var mfEdgeMap = {};
		function indexMfEdges(nodes, edges) {
			if (edges) edges.forEach(function(e) { mfEdgeMap[e.id] = e; });
			if (nodes) nodes.forEach(function(n) {
				if (n.edges) n.edges.forEach(function(e) { mfEdgeMap[e.id] = e; });
				if (n.children) indexMfEdges(n.children, null);
			});
		}
		indexMfEdges(data.nodes, data.edges);

		// --- Render SVG (exposed on window for module script) ---
		window.renderSvg = renderSvg;
		function renderSvg(layout) {
			if (diagramType === 'page') {
				renderPageWireframe(data);
			} else if (diagramType === 'microflow') {
				renderMicroflow(layout);
			} else if (diagramType === 'oql-queryplan') {
				renderQueryPlan(layout);
			} else if (diagramType === 'domainmodel') {
				renderDomainModel(layout);
			} else {
				renderModuleOverview(layout);
			}
		}

		// --- Domain model collapse/expand ---
		window.relayoutDomainModel = async function() {
			if (diagramType !== 'domainmodel' || !window.elkInstance || !window.elkGraph) return;
			var headerH = 28;
			var attrLineH = 18;
			var entities = data.entities || [];
			var graph = JSON.parse(JSON.stringify(window.elkGraph));
			graph.children.forEach(function(child) {
				var ent = entityMap[child.id];
				if (ent && ent.attributes && ent.attributes.length > 0 && collapsedEntities[child.id]) {
					child.height = headerH + attrLineH;
				}
			});
			var result = await window.elkInstance.layout(graph);
			renderDomainModel(result);
			applyTransform();
		};

		// --- Microflow collapse/expand ---
		window.relayoutMicroflow = async function() {
			if (diagramType !== 'microflow' || !window.elkInstance || !window.elkGraph) return;
			var headerH = 28;
			var graph = JSON.parse(JSON.stringify(window.elkGraph));
			function collapseChildren(children) {
				if (!children) return;
				children.forEach(function(child) {
					var nd = mfNodeMap[child.id];
					// Skip compound nodes (loops with children) — their size is computed by ELK
					if (child.children && child.children.length > 0) {
						collapseChildren(child.children);
						return;
					}
					if (nd && nd.details && nd.details.length > 0 && collapsedActivities[child.id]) {
						child.height = headerH + 8;
					}
				});
			}
			collapseChildren(graph.children);
			var result = await window.elkInstance.layout(graph);
			renderMicroflow(result);
			applyTransform();
		};

		function toggleCollapseAll() {
			allCollapsed = !allCollapsed;
			var btn = document.getElementById('btn-collapse');
			btn.textContent = allCollapsed ? 'Expand All' : 'Collapse All';

			if (diagramType === 'domainmodel') {
				var entities = data.entities || [];
				entities.forEach(function(e) {
					if (e.attributes && e.attributes.length > 0) {
						if (allCollapsed) {
							collapsedEntities[e.id] = true;
						} else {
							delete collapsedEntities[e.id];
						}
					}
				});
				window.relayoutDomainModel();
			} else if (diagramType === 'microflow') {
				var nodes = data.nodes || [];
				nodes.forEach(function(n) {
					if (n.details && n.details.length > 0) {
						if (allCollapsed) {
							collapsedActivities[n.id] = true;
						} else {
							delete collapsedActivities[n.id];
						}
					}
				});
				window.relayoutMicroflow();
			}
		}

		${microflowRendererJs()}
		${moduleOverviewRendererJs()}
		${queryPlanRendererJs()}
		${domainModelRendererJs()}
		${pageWireframeRendererJs()}

		// --- Pan & Zoom ---
		let scale = 1;
		let panX = 0;
		let panY = 0;
		let isDragging = false;
		let dragStartX = 0;
		let dragStartY = 0;
		let dragStartPanX = 0;
		let dragStartPanY = 0;

		const viewport = document.getElementById('diagram-viewport');
		const canvas = document.getElementById('diagram-canvas');
		const zoomInfo = document.getElementById('zoom-info');

		function applyTransform() {
			canvas.style.transform = 'translate(' + panX + 'px,' + panY + 'px) scale(' + scale + ')';
			zoomInfo.textContent = Math.round(scale * 100) + '%';
		}

		viewport.addEventListener('wheel', function(e) {
			e.preventDefault();
			const rect = viewport.getBoundingClientRect();
			const mx = e.clientX - rect.left;
			const my = e.clientY - rect.top;
			const delta = e.deltaY > 0 ? 0.9 : 1.1;
			const newScale = Math.min(Math.max(scale * delta, 0.1), 10);
			panX = mx - (mx - panX) * (newScale / scale);
			panY = my - (my - panY) * (newScale / scale);
			scale = newScale;
			applyTransform();
		}, { passive: false });

		viewport.addEventListener('mousedown', function(e) {
			if (e.button !== 0) return;
			isDragging = true;
			dragStartX = e.clientX;
			dragStartY = e.clientY;
			dragStartPanX = panX;
			dragStartPanY = panY;
			viewport.classList.add('dragging');
		});
		window.addEventListener('mousemove', function(e) {
			if (!isDragging) return;
			panX = dragStartPanX + (e.clientX - dragStartX);
			panY = dragStartPanY + (e.clientY - dragStartY);
			applyTransform();
		});
		window.addEventListener('mouseup', function() {
			isDragging = false;
			viewport.classList.remove('dragging');
		});

		function zoomIn() {
			const rect = viewport.getBoundingClientRect();
			const cx = rect.width / 2;
			const cy = rect.height / 2;
			const newScale = Math.min(scale * 1.25, 10);
			panX = cx - (cx - panX) * (newScale / scale);
			panY = cy - (cy - panY) * (newScale / scale);
			scale = newScale;
			applyTransform();
		}
		function zoomOut() {
			const rect = viewport.getBoundingClientRect();
			const cx = rect.width / 2;
			const cy = rect.height / 2;
			const newScale = Math.max(scale * 0.8, 0.1);
			panX = cx - (cx - panX) * (newScale / scale);
			panY = cy - (cy - panY) * (newScale / scale);
			scale = newScale;
			applyTransform();
		}
		window.resetView = resetView;
		function resetView() {
			const svg = canvas.querySelector('svg');
			if (!svg) { scale = 1; panX = 0; panY = 0; applyTransform(); return; }
			const vw = viewport.clientWidth;
			const vh = viewport.clientHeight;
			const sw = svg.getBoundingClientRect().width / scale;
			const sh = svg.getBoundingClientRect().height / scale;
			scale = Math.min(vw / sw, vh / sh, 2) * 0.95;
			panX = (vw - sw * scale) / 2;
			panY = (vh - sh * scale) / 2;
			applyTransform();
		}

		// Close popover on click outside
		document.addEventListener('mousedown', function(e) {
			var popover = document.getElementById('node-popover');
			if (popover && !popover.contains(e.target)) {
				popover.style.display = 'none';
			}
		});

		// --- Highlight/Dim message handler for diagram-with-source sync ---
		window.addEventListener('message', function(event) {
			var msg = event.data;
			if (msg.type === 'highlightNode') {
				var allNodes = document.querySelectorAll('.mf-node, .entity-node, .wf-node');
				if (!msg.nodeId) {
					// Clear all highlights
					allNodes.forEach(function(g) {
						g.classList.remove('highlighted', 'dimmed');
					});
					return;
				}
				allNodes.forEach(function(g) {
					var id = g.getAttribute('data-node-id') || g.getAttribute('data-entity');
					if (id === msg.nodeId) {
						g.classList.add('highlighted');
						g.classList.remove('dimmed');
					} else {
						g.classList.remove('highlighted');
						g.classList.add('dimmed');
					}
				});
			}
		});

		let showingSource = false;
		function toggleSource() {
			showingSource = !showingSource;
			document.getElementById('source-view').style.display = showingSource ? 'block' : 'none';
			document.getElementById('diagram-viewport').style.display = showingSource ? 'none' : '';
			document.getElementById('legend').style.display = showingSource ? 'none' : 'flex';
		}

		function copySource() {
			navigator.clipboard.writeText(JSON.stringify(data, null, 2));
		}
	</script>
</body>
</html>`;
}
