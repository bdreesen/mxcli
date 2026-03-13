// SPDX-License-Identifier: Apache-2.0

import { escapeHtml, escapeForTemplate } from './mxcliRunner';

/**
 * Returns the complete HTML content for the Mermaid-based webview panel.
 */
export function getMermaidWebviewContent(mermaidSource: string, title: string, theme: string = 'clean'): string {
	// Escape for embedding in JavaScript string
	const escapedSource = escapeForTemplate(mermaidSource);

	return `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>${escapeHtml(title)}</title>
	<style>
		body {
			margin: 0;
			padding: 0;
			background-color: var(--vscode-editor-background, #1e1e1e);
			color: var(--vscode-editor-foreground, #d4d4d4);
			font-family: var(--vscode-font-family, sans-serif);
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
		.toolbar button.active {
			background: var(--vscode-button-secondaryBackground, #3a3d41);
			outline: 1px solid var(--vscode-focusBorder, #007acc);
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
	</style>
</head>
<body>
	<div class="toolbar">
		<span class="title">${escapeHtml(title)}</span>
		<span class="zoom-info" id="zoom-info">100%</span>
		<button onclick="resetView()">Fit</button>
		<button onclick="zoomIn()">+</button>
		<button onclick="zoomOut()">&minus;</button>
		<button id="btn-direction" onclick="toggleDirection()" style="display:none" title="Toggle flow direction">&#x2B0D; LR</button>
		<button id="btn-details" onclick="toggleDetails()" style="display:none" title="Show/hide activity details">Details</button>
		<button onclick="toggleSource()">Source</button>
		<button onclick="copySource()">Copy</button>
	</div>
	<div class="legend" id="legend" style="display:none">
		<div class="legend-item"><div class="legend-swatch" style="background:#4a90d9"></div>Persistent</div>
		<div class="legend-item"><div class="legend-swatch" style="background:#e8a838"></div>Non-Persistent</div>
		<div class="legend-item"><div class="legend-swatch" style="background:#9b59b6"></div>External</div>
		<div class="legend-item"><div class="legend-swatch" style="background:#27ae60"></div>View</div>
	</div>
	<div id="diagram-viewport">
		<div id="diagram-canvas">
			<pre class="mermaid">${escapeHtml(mermaidSource)}</pre>
		</div>
	</div>
	<div id="node-popover"></div>
	<pre id="source-view" class="source-view">${escapeHtml(mermaidSource)}</pre>

	<script type="module">
		import mermaid from 'https://cdn.jsdelivr.net/npm/mermaid@11/dist/mermaid.esm.min.mjs';

		const isDark = document.body.classList.contains('vscode-dark') ||
			document.body.getAttribute('data-vscode-theme-kind') === 'vscode-dark' ||
			getComputedStyle(document.body).getPropertyValue('--vscode-editor-background').trim().match(/^#[0-3]/);

		mermaid.initialize({
			startOnLoad: false,
			theme: isDark ? 'dark' : 'default',
			securityLevel: 'loose',
			er: { useMaxWidth: false },
			flowchart: { useMaxWidth: false, htmlLabels: true },
		});

		// Expose mermaid for re-rendering (direction toggle)
		window.__mermaid = mermaid;

		// Render, then apply colors and fit view
		await mermaid.run();
		window.postRender();
	</script>
	<script>
		const mermaidSource = \`${escapedSource}\`;

		// --- Diagram metadata (emitted by Go backend as %% @key value) ---
		const diagramType = (mermaidSource.match(/%% @type (\\w+)/) || [])[1] || '';
		const diagramDirection = (mermaidSource.match(/%% @direction (\\w+)/) || [])[1] || '';
		let currentDirection = diagramDirection || 'LR';

		// Show direction toggle for flowcharts
		if (diagramType === 'flowchart') {
			document.getElementById('btn-direction').style.display = '';
			document.getElementById('btn-direction').textContent = '\\u2B0D ' + currentDirection;
		}

		// --- Node detail info (emitted by Go backend as %% @nodeinfo {...}) ---
		const nodeInfo = {};
		const nodeInfoMatch = mermaidSource.match(/%% @nodeinfo (\\{.+\\})/);
		if (nodeInfoMatch) {
			try {
				Object.assign(nodeInfo, JSON.parse(nodeInfoMatch[1]));
			} catch(e) {}
		}

		// Show details button if we have node info
		if (Object.keys(nodeInfo).length > 0) {
			document.getElementById('btn-details').style.display = '';
		}

		let detailsExpanded = false;

		// --- Entity color map (emitted by Go backend as %% @colors {...}) ---
		const colorMap = {};
		const colorMatch = mermaidSource.match(/%% @colors \\{(.+?)\\}/);
		if (colorMatch) {
			try {
				const parsed = JSON.parse('{' + colorMatch[1] + '}');
				Object.assign(colorMap, parsed);
			} catch(e) {}
		}

		const categoryColors = {
			persistent:    { fill: '#4a90d9', stroke: '#2c6fad', text: '#fff' },
			nonpersistent: { fill: '#e8a838', stroke: '#c48820', text: '#fff' },
			external:      { fill: '#9b59b6', stroke: '#7d3c98', text: '#fff' },
			view:          { fill: '#27ae60', stroke: '#1e8449', text: '#fff' },
		};
		const categoryColorsDark = categoryColors; // same palette works for both

		// --- Pan & Zoom state ---
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

		// Mouse wheel zoom
		viewport.addEventListener('wheel', (e) => {
			e.preventDefault();
			const rect = viewport.getBoundingClientRect();
			const mx = e.clientX - rect.left;
			const my = e.clientY - rect.top;

			const delta = e.deltaY > 0 ? 0.9 : 1.1;
			const newScale = Math.min(Math.max(scale * delta, 0.1), 10);

			// Zoom toward mouse position
			panX = mx - (mx - panX) * (newScale / scale);
			panY = my - (my - panY) * (newScale / scale);
			scale = newScale;
			applyTransform();
		}, { passive: false });

		// Mouse drag to pan
		viewport.addEventListener('mousedown', (e) => {
			if (e.button !== 0) return;
			isDragging = true;
			dragStartX = e.clientX;
			dragStartY = e.clientY;
			dragStartPanX = panX;
			dragStartPanY = panY;
			viewport.classList.add('dragging');
		});
		window.addEventListener('mousemove', (e) => {
			if (!isDragging) return;
			panX = dragStartPanX + (e.clientX - dragStartX);
			panY = dragStartPanY + (e.clientY - dragStartY);
			applyTransform();
		});
		window.addEventListener('mouseup', () => {
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

		// --- Post-render: apply entity colors + fit ---
		window.postRender = function() {
			const svg = canvas.querySelector('svg');
			if (!svg) return;

			// Remove max-width so we get the natural size for pan/zoom
			svg.style.maxWidth = 'none';

			// Show legend if we have color metadata
			if (Object.keys(colorMap).length > 0) {
				document.getElementById('legend').style.display = 'flex';
				applyEntityColors(svg);
			}

			// Fit diagram to viewport
			setTimeout(resetView, 50);
		};

		// Apply entity fill colors using inline styles (overrides Mermaid CSS)
		function applyEntityColors(svg) {
			// Strategy: Find text elements matching entity names in colorMap.
			// For each match, walk up the DOM to find the entity box rect(s)
			// and apply fill/stroke via inline styles (CSS overrides SVG attributes).
			const colored = new Set(); // track already-colored entity names

			// First pass: try matching by entity group IDs (g[id*="entity-Name"])
			for (const [name, cat] of Object.entries(colorMap)) {
				if (!categoryColors[cat]) continue;
				const groups = svg.querySelectorAll('g[id*="entity-' + name + '"]');
				groups.forEach((g) => {
					const rects = g.querySelectorAll('rect');
					rects.forEach((rect) => {
						rect.style.fill = categoryColors[cat].fill;
						rect.style.stroke = categoryColors[cat].stroke;
					});
					colored.add(name);
				});
			}

			// Second pass: match by text content for entities not yet colored
			const texts = svg.querySelectorAll('text');
			texts.forEach((t) => {
				const label = (t.textContent || '').trim();
				const cat = colorMap[label];
				if (!cat || !categoryColors[cat] || colored.has(label)) return;

				// Walk up to find the entity container group
				let g = t.closest('g');
				if (!g) return;

				// The text might be in a nested label group; walk up further
				// to find a group that contains rect elements
				let target = g;
				for (let i = 0; i < 3 && target; i++) {
					const rects = target.querySelectorAll(':scope > rect');
					if (rects.length > 0) {
						rects.forEach((rect) => {
							rect.style.fill = categoryColors[cat].fill;
							rect.style.stroke = categoryColors[cat].stroke;
						});
						// Color the entity name text for contrast
						const titleTexts = target.querySelectorAll('text');
						titleTexts.forEach((te) => {
							if ((te.textContent || '').trim() === label) {
								te.style.fill = categoryColors[cat].text;
							}
						});
						colored.add(label);
						break;
					}
					target = target.parentElement;
				}
			});
		}

		let showingSource = false;
		function toggleSource() {
			showingSource = !showingSource;
			document.getElementById('source-view').style.display = showingSource ? 'block' : 'none';
			document.getElementById('diagram-viewport').style.display = showingSource ? 'none' : 'flex';
			document.getElementById('legend').style.display = showingSource ? 'none' :
				(Object.keys(colorMap).length > 0 ? 'flex' : 'none');
		}

		function copySource() {
			navigator.clipboard.writeText(mermaidSource);
		}

		// --- Build current Mermaid source with active options applied ---
		function buildCurrentSource() {
			let src = mermaidSource;
			// Apply direction
			if (diagramType === 'flowchart') {
				src = src.replace(/flowchart (LR|TD)/, 'flowchart ' + currentDirection);
			}
			// Apply detail expansion: inject detail lines into node labels
			if (detailsExpanded && Object.keys(nodeInfo).length > 0) {
				const lines = src.split('\\n');
				for (let i = 0; i < lines.length; i++) {
					const line = lines[i];
					for (const [nodeId, details] of Object.entries(nodeInfo)) {
						// Match node definition: "    nodeId[" or "    nodeId[/"
						if (!line.trimStart().startsWith(nodeId)) continue;
						const after = line.substring(line.indexOf(nodeId) + nodeId.length);
						// Only expand box nodes ["..."] and trapezoid [/"..."/]
						if (after.startsWith('["') || after.startsWith('[/"')) {
							const detailHtml = details
								.map(d => '<br/><small>' + d.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '#quot;') + '</small>')
								.join('');
							// Insert detail HTML before the closing "]" or "/]"
							const closeBracket = after.endsWith('/]') ? '/]' : after.endsWith('"]') ? '"]' : ']';
							const insertPos = line.lastIndexOf(closeBracket);
							if (insertPos > 0) {
								lines[i] = line.substring(0, insertPos) + detailHtml + line.substring(insertPos);
							}
						}
						break;
					}
				}
				src = lines.join('\\n');
			}
			return src;
		}

		async function rerender() {
			const src = buildCurrentSource();
			const newPre = document.createElement('pre');
			newPre.className = 'mermaid';
			newPre.textContent = src;
			canvas.innerHTML = '';
			canvas.appendChild(newPre);
			await window.__mermaid.run();
			window.postRender();
		}

		async function toggleDirection() {
			currentDirection = currentDirection === 'LR' ? 'TD' : 'LR';
			document.getElementById('btn-direction').textContent = '\\u2B0D ' + currentDirection;
			await rerender();
		}

		async function toggleDetails() {
			detailsExpanded = !detailsExpanded;
			const btn = document.getElementById('btn-details');
			if (detailsExpanded) {
				btn.classList.add('active');
			} else {
				btn.classList.remove('active');
			}
			await rerender();
		}

		// --- Click-on-node popover ---
		const popover = document.getElementById('node-popover');

		// Close popover on click outside
		document.addEventListener('mousedown', (e) => {
			if (!popover.contains(e.target)) {
				popover.style.display = 'none';
			}
		});

		// After rendering, attach click handlers to flowchart nodes
		const origPostRender = window.postRender;
		window.postRender = function() {
			origPostRender();
			if (Object.keys(nodeInfo).length === 0) return;

			const svg = canvas.querySelector('svg');
			if (!svg) return;

			// Find all node groups and attach click handlers
			for (const [nodeId, details] of Object.entries(nodeInfo)) {
				// Mermaid generates nodes with id like "flowchart-nodeId-N" or just contains nodeId
				const nodeEls = svg.querySelectorAll('[id*="' + nodeId + '"]');
				nodeEls.forEach((el) => {
					el.style.cursor = 'pointer';
					el.addEventListener('click', (e) => {
						e.stopPropagation();

						// Find the node's label text
						const textEl = el.querySelector('span, text');
						const title = textEl ? (textEl.textContent || '').trim().split('\\n')[0] : nodeId;

						// Build popover content
						let html = '<div class="popover-title">' + escapeForHtml(title) + '</div>';
						details.forEach((line) => {
							html += '<div class="popover-line">' + escapeForHtml(line) + '</div>';
						});
						popover.innerHTML = html;

						// Position near the click
						const vpRect = viewport.getBoundingClientRect();
						let left = e.clientX - vpRect.left + 10;
						let top = e.clientY - vpRect.top + 10;
						// Keep within viewport
						popover.style.display = 'block';
						if (left + popover.offsetWidth > vpRect.width) {
							left = Math.max(0, e.clientX - vpRect.left - popover.offsetWidth - 10);
						}
						if (top + popover.offsetHeight > vpRect.height) {
							top = Math.max(0, e.clientY - vpRect.top - popover.offsetHeight - 10);
						}
						popover.style.left = (left + vpRect.left) + 'px';
						popover.style.top = (top + vpRect.top) + 'px';
					});
				});
			}
		};

		function escapeForHtml(s) {
			return s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
		}
	</script>
</body>
</html>`;
}
