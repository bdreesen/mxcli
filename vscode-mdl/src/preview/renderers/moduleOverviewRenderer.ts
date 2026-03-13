// SPDX-License-Identifier: Apache-2.0

/**
 * Returns inline JavaScript code for the module overview renderer.
 * Depends on shared state: data, moduleMap, and sketch helper functions.
 */
export function moduleOverviewRendererJs(): string {
	return `
		// --- Module Overview renderer ---
		function renderModuleOverview(layout) {
			var maxX = 0, maxY = 0;
			layout.children.forEach(function(node) {
				var right = node.x + node.width;
				var bottom = node.y + node.height;
				if (right > maxX) maxX = right;
				if (bottom > maxY) maxY = bottom;
			});

			var padding = 60;
			var titleHeight = 40;
			var svgWidth = maxX + padding * 2;
			var svgHeight = maxY + padding * 2 + titleHeight + 20;
			var fontFamily = diagramTheme === 'clean'
				? "var(--vscode-font-family, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif)"
				: "'Architects Daughter', cursive";

			var svg = '<svg xmlns="http://www.w3.org/2000/svg" width="' + svgWidth + '" height="' + svgHeight + '">';
			svg += '<defs>' + svgFilterDefs() + '</defs>';

			svg += '<text x="' + padding + '" y="32" font-size="20" fill="' + inkColor + '" font-family="' + fontFamily + '">System Overview</text>';

			// PoC badge (sketch only)
			if (diagramTheme === 'sketch') {
				var badgeRng = makeRng(42);
				var badgeX = padding + 195;
				svg += '<path d="' + roughRoundedRect(badgeX, 17, 64, 20, 10, badgeRng) + '" fill="none" stroke="' + secondaryColor + '" stroke-width="1.2" filter="url(#pencil)"/>';
				svg += '<text x="' + (badgeX + 32) + '" y="31" font-size="9" fill="' + secondaryColor + '" text-anchor="middle" font-family="' + fontFamily + '">PoC draft</text>';
			}

			var offsetY = titleHeight;

			if (layout.edges) {
				layout.edges.forEach(function(edge, idx) {
					var edgeData = data.edges[idx];
					var thickness = edgeData ? Math.min(Math.max(Math.ceil(edgeData.count / 3), 1), 4) : 1;
					var edgeRng = makeRng(idx * 1000 + 7);

					if (edge.sections) {
						edge.sections.forEach(function(section) {
							var pts = [{ x: section.startPoint.x + padding, y: section.startPoint.y + padding + offsetY }];
							if (section.bendPoints) {
								section.bendPoints.forEach(function(bp) {
									pts.push({ x: bp.x + padding, y: bp.y + padding + offsetY });
								});
							}
							pts.push({ x: section.endPoint.x + padding, y: section.endPoint.y + padding + offsetY });

							var d = '';
							for (var i = 0; i < pts.length - 1; i++) {
								var seg = roughLine(pts[i].x, pts[i].y, pts[i + 1].x, pts[i + 1].y, edgeRng, 1.0);
								if (i === 0) { d = seg; } else { d += seg.replace(/^M [^ ]+ [^ ]+/, ''); }
							}
							svg += '<path d="' + d + '" fill="none" stroke="' + connectorColor + '" stroke-width="' + thickness + '" opacity="0.5" stroke-linecap="round" filter="url(#pencil)"/>';

							var last = pts[pts.length - 1];
							var prev = pts[pts.length - 2];
							var angle = Math.atan2(last.y - prev.y, last.x - prev.x);
							svg += roughArrowhead(last.x, last.y, angle, edgeRng);
						});
					}

					if (edgeData && edgeData.count > 1 && edge.sections && edge.sections.length > 0) {
						var section = edge.sections[0];
						var points = [section.startPoint];
						if (section.bendPoints) points.push.apply(points, section.bendPoints);
						points.push(section.endPoint);
						var mid = Math.floor(points.length / 2);
						var lx = points[mid].x + padding;
						var ly = points[mid].y + padding + offsetY - 6;
						svg += '<text x="' + lx + '" y="' + ly + '" font-size="10" fill="' + secondaryColor + '" opacity="0.6" text-anchor="middle" font-family="' + fontFamily + '">' + edgeData.count + '</text>';
					}
				});
			}

			layout.children.forEach(function(node) {
				var m = moduleMap[node.id];
				var isSystem = m && m.isSystem;
				var base = isSystem ? sysBase : appBase;
				var light = isSystem ? sysLight : appLight;
				var x = node.x + padding;
				var y = node.y + padding + offsetY;
				var w = node.width;
				var h = node.height;
				var nodeRng = makeRng(hashStr(node.id));

				svg += '<g class="module-node" data-module="' + escHtml(node.id) + '" style="cursor:pointer">';
				// Drop shadow for clean theme
				if (diagramTheme === 'clean') {
					svg += '<rect x="' + x + '" y="' + y + '" width="' + w + '" height="' + h + '" rx="6" fill="var(--vscode-editor-background, #1e1e1e)" filter="url(#clean-shadow)"/>';
				}
				svg += markerFill(x, y, w, h, light, nodeRng);
				var borderRng = makeRng(hashStr(node.id) + 100);
				svg += '<path d="' + roughRoundedRect(x, y, w, h, 6, borderRng) + '" fill="none" stroke="' + base + '" stroke-width="1.5" stroke-linecap="round" filter="url(#pencil)"/>';
				svg += '<text x="' + (x + w / 2) + '" y="' + (y + 22) + '" font-size="13" fill="' + inkColor + '" text-anchor="middle" font-family="' + fontFamily + '">' + escHtml(node.id) + '</text>';
				if (m) {
					var stats = m.entityCount + 'E  ' + m.microflowCount + 'MF  ' + m.pageCount + 'P';
					svg += '<text x="' + (x + w / 2) + '" y="' + (y + 40) + '" font-size="10" fill="' + secondaryColor + '" opacity="0.7" text-anchor="middle" font-family="' + fontFamily + '">' + stats + '</text>';
				}
				svg += '</g>';
			});

			// Footer (sketch only)
			if (diagramTheme === 'sketch') {
				svg += '<text x="10" y="' + (svgHeight - 8) + '" font-size="10" fill="' + secondaryColor + '" opacity="0.4" font-family="' + fontFamily + '">sketch — subject to change</text>';
			}
			svg += '</svg>';

			var canvas = document.getElementById('diagram-canvas');
			canvas.innerHTML = svg;

			canvas.querySelectorAll('.module-node').forEach(function(g) {
				g.addEventListener('click', function(e) {
					e.stopPropagation();
					var moduleName = g.getAttribute('data-module');
					if (moduleName) {
						vscodeApi.postMessage({ type: 'openModule', moduleName: moduleName });
					}
				});
			});
		}
	`;
}
