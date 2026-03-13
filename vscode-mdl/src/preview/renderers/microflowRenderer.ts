// SPDX-License-Identifier: Apache-2.0

/**
 * Returns inline JavaScript code for the microflow renderer.
 * Depends on shared state: data, mfNodeMap, mfEdgeMap, mfCategoryColors,
 * collapsedActivities, vscodeApi, and sketch helper functions.
 */
export function microflowRendererJs(): string {
	return `
		// --- Microflow renderer ---
		function renderMicroflow(layout) {
			var maxX = 0, maxY = 0;
			layout.children.forEach(function(node) {
				var right = node.x + node.width;
				var bottom = node.y + node.height;
				if (right > maxX) maxX = right;
				if (bottom > maxY) maxY = bottom;
			});

			var padding = 60;
			var titleHeight = 40;
			var headerHeight = 28;
			var detailLineHeight = 16;
			var svgWidth = maxX + padding * 2;
			var svgHeight = maxY + padding * 2 + titleHeight + 20;
			var fontFamily = diagramTheme === 'clean'
				? "var(--vscode-font-family, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif)"
				: "'Architects Daughter', cursive";
			var errorColor = '#d9534f';

			var svg = '<svg xmlns="http://www.w3.org/2000/svg" width="' + svgWidth + '" height="' + svgHeight + '">';
			svg += '<defs>' + svgFilterDefs() + '</defs>';

			// Title
			var mfName = data.name || 'Microflow';
			svg += '<text x="' + padding + '" y="32" font-size="20" fill="' + inkColor + '" font-family="' + fontFamily + '">' + escHtml(mfName) + '</text>';

			// PoC badge (sketch only)
			if (diagramTheme === 'sketch') {
				var badgeRng = makeRng(42);
				var titleWidth = mfName.length * 11 + 10;
				var badgeX = padding + titleWidth;
				svg += '<path d="' + roughRoundedRect(badgeX, 17, 64, 20, 10, badgeRng) + '" fill="none" stroke="' + secondaryColor + '" stroke-width="1.2" filter="url(#pencil)"/>';
				svg += '<text x="' + (badgeX + 32) + '" y="31" font-size="9" fill="' + secondaryColor + '" text-anchor="middle" font-family="' + fontFamily + '">PoC draft</text>';
			}

			var offsetY = titleHeight;

			// Helper: render edges with offset (ox, oy are absolute offsets for section points)
			function renderMfEdges(edges, ox, oy) {
				if (!edges) return;
				edges.forEach(function(edge) {
					var edgeData = mfEdgeMap[edge.id];
					var isError = edgeData && edgeData.isErrorHandler;
					var edgeColor = isError ? errorColor : connectorColor;
					var dashAttr = isError ? ' stroke-dasharray="6 4"' : '';
					var edgeRng = makeRng(hashStr(edge.id) + 7);

					if (edge.sections) {
						edge.sections.forEach(function(section) {
							var pts = [{ x: section.startPoint.x + ox, y: section.startPoint.y + oy }];
							if (section.bendPoints) {
								section.bendPoints.forEach(function(bp) {
									pts.push({ x: bp.x + ox, y: bp.y + oy });
								});
							}
							pts.push({ x: section.endPoint.x + ox, y: section.endPoint.y + oy });

							var d = '';
							for (var i = 0; i < pts.length - 1; i++) {
								var seg = roughLine(pts[i].x, pts[i].y, pts[i + 1].x, pts[i + 1].y, edgeRng, 1.0);
								if (i === 0) { d = seg; } else { d += seg.replace(/^M [^ ]+ [^ ]+/, ''); }
							}
							svg += '<path d="' + d + '" fill="none" stroke="' + edgeColor + '" stroke-width="1.5" opacity="0.6" stroke-linecap="round"' + dashAttr + ' filter="url(#pencil)"/>';

							// Arrowhead
							var last = pts[pts.length - 1];
							var prev = pts[pts.length - 2];
							var angle = Math.atan2(last.y - prev.y, last.x - prev.x);
							svg += roughArrowhead(last.x, last.y, angle, edgeRng, edgeColor);

							// Edge label at midpoint
							if (edgeData && edgeData.label) {
								var allPts = [section.startPoint];
								if (section.bendPoints) allPts.push.apply(allPts, section.bendPoints);
								allPts.push(section.endPoint);
								var midIdx = Math.floor(allPts.length / 2);
								var lx = allPts[midIdx].x + ox;
								var ly = allPts[midIdx].y + oy - 6;
								svg += '<text x="' + lx + '" y="' + ly + '" font-size="10" fill="' + secondaryColor + '" opacity="0.8" text-anchor="middle" font-style="italic" font-family="' + fontFamily + '">' + escHtml(edgeData.label) + '</text>';
							}
						});
					}
				});
			}

			// Helper: render nodes with offset (ox, oy translate node positions to absolute SVG coords)
			function renderMfNodes(children, ox, oy) {
				if (!children) return;
				children.forEach(function(node) {
					var nd = mfNodeMap[node.id];
					if (!nd) return;
					var cat = nd.category || 'variable';
					var colors = mfCategoryColors[cat] || mfCategoryColors.variable;
					var x = node.x + ox;
					var y = node.y + oy;
					var w = node.width;
					var h = node.height;
					var nodeRng = makeRng(hashStr(node.id));

					var detailLines = nd.details || [];
					var hasDetails = detailLines.length > 0;
					var isCompoundLoop = nd.type === 'loop' && node.children && node.children.length > 0;

					svg += '<g class="mf-node" data-node-id="' + escHtml(node.id) + '" style="cursor:pointer">';

					if (nd.type === 'start' || nd.type === 'end' || nd.type === 'continue' || nd.type === 'break' || nd.type === 'error') {
						// Pill shape (rounded rect with r = h/2)
						var pillR = h / 2;
						if (diagramTheme === 'clean') {
							svg += '<rect x="' + x + '" y="' + y + '" width="' + w + '" height="' + h + '" rx="' + pillR + '" fill="var(--vscode-editor-background, #1e1e1e)" filter="url(#clean-shadow)"/>';
						}
						svg += markerFill(x, y, w, h, colors.light, makeRng(hashStr(node.id) + 50));
						var pillBorderRng = makeRng(hashStr(node.id) + 100);
						svg += '<path d="' + roughRoundedRect(x, y, w, h, pillR, pillBorderRng) + '" fill="none" stroke="' + colors.base + '" stroke-width="1.5" stroke-linecap="round" filter="url(#pencil)"/>';
						svg += '<text x="' + (x + w / 2) + '" y="' + (y + h / 2 + 5) + '" font-size="12" fill="' + inkColor + '" font-weight="600" text-anchor="middle" font-family="' + fontFamily + '">' + escHtml(nd.label) + '</text>';
					} else if (nd.type === 'split') {
						// Diamond shape
						var cx = x + w / 2;
						var cy = y + h / 2;
						var diamondRng = makeRng(hashStr(node.id) + 100);
						if (diagramTheme === 'clean') {
							// Clean diamond with solid fill
							var dd = 'M ' + cx + ' ' + y + ' L ' + (x + w) + ' ' + cy + ' L ' + cx + ' ' + (y + h) + ' L ' + x + ' ' + cy + ' Z';
							svg += '<path d="' + dd + '" fill="var(--vscode-editor-background, #1e1e1e)" filter="url(#clean-shadow)"/>';
						}
						var d1 = roughLine(cx, y, x + w, cy, diamondRng, 1.2);
						var d2 = roughLine(x + w, cy, cx, y + h, diamondRng, 1.2).replace(/^M [^ ]+ [^ ]+/, '');
						var d3 = roughLine(cx, y + h, x, cy, diamondRng, 1.2).replace(/^M [^ ]+ [^ ]+/, '');
						var d4 = roughLine(x, cy, cx, y, diamondRng, 1.2).replace(/^M [^ ]+ [^ ]+/, '');
						svg += markerFill(x + w * 0.2, y + h * 0.2, w * 0.6, h * 0.6, mfCategoryColors.controlflow.light, makeRng(hashStr(node.id) + 50));
						svg += '<path d="' + d1 + d2 + d3 + d4 + ' Z" fill="none" stroke="' + mfCategoryColors.controlflow.base + '" stroke-width="1.5" stroke-linecap="round" filter="url(#pencil)"/>';
						var splitLabel = nd.label;
						var maxChars = Math.floor(w / 7.5 - 4);
						if (splitLabel.length > maxChars && maxChars > 3) {
							splitLabel = splitLabel.substring(0, maxChars - 3) + '...';
						}
						svg += '<text x="' + cx + '" y="' + (cy + 4) + '" font-size="10" fill="' + inkColor + '" text-anchor="middle" font-family="' + fontFamily + '">' + escHtml(splitLabel) + '</text>';
					} else if (nd.type === 'merge') {
						// Small filled circle
						var mcx = x + w / 2;
						var mcy = y + h / 2;
						var mr = w / 2;
						svg += '<circle cx="' + mcx + '" cy="' + mcy + '" r="' + mr + '" fill="' + mfCategoryColors.controlflow.light + '" stroke="' + mfCategoryColors.controlflow.base + '" stroke-width="1.5" filter="url(#pencil)"/>';
					} else if (isCompoundLoop) {
						// Compound loop: double-border container with children rendered inside
						if (diagramTheme === 'clean') {
							svg += '<rect x="' + x + '" y="' + y + '" width="' + w + '" height="' + h + '" rx="6" fill="var(--vscode-editor-background, #1e1e1e)" filter="url(#clean-shadow)"/>';
						}
						svg += markerFill(x, y, w, headerHeight, colors.light, makeRng(hashStr(node.id) + 50));
						var loopOuterRng = makeRng(hashStr(node.id) + 100);
						svg += '<path d="' + roughRoundedRect(x, y, w, h, 6, loopOuterRng) + '" fill="none" stroke="' + colors.base + '" stroke-width="2" stroke-linecap="round" filter="url(#pencil)"/>';
						var loopInnerRng = makeRng(hashStr(node.id) + 150);
						svg += '<path d="' + roughRoundedRect(x + 3, y + 3, w - 6, h - 6, 4, loopInnerRng) + '" fill="none" stroke="' + colors.base + '" stroke-width="1" opacity="0.5" stroke-linecap="round" filter="url(#pencil)"/>';
						// Loop label at top
						var loopLabel = nd.label;
						if (hasDetails) {
							loopLabel += ' (' + detailLines.join(', ') + ')';
						}
						svg += '<text x="' + (x + w / 2) + '" y="' + (y + 18) + '" font-size="12" fill="' + inkColor + '" font-weight="600" text-anchor="middle" font-family="' + fontFamily + '">' + escHtml(loopLabel) + '</text>';
						// Close group for container, then render inner content
						svg += '</g>';
						renderMfEdges(node.edges, x, y);
						renderMfNodes(node.children, x, y);
						return; // skip the closing </g> below — already closed
					} else if (nd.type === 'loop') {
						// Simple loop without children (fallback)
						var isLoopCollapsed = hasDetails && !!collapsedActivities[node.id];
						if (diagramTheme === 'clean') {
							svg += '<rect x="' + x + '" y="' + y + '" width="' + w + '" height="' + h + '" rx="6" fill="var(--vscode-editor-background, #1e1e1e)" filter="url(#clean-shadow)"/>';
						}
						svg += markerFill(x, y, w, headerHeight, colors.light, makeRng(hashStr(node.id) + 50));
						var loopOuterRng2 = makeRng(hashStr(node.id) + 100);
						svg += '<path d="' + roughRoundedRect(x, y, w, h, 6, loopOuterRng2) + '" fill="none" stroke="' + colors.base + '" stroke-width="2" stroke-linecap="round" filter="url(#pencil)"/>';
						var loopInnerRng2 = makeRng(hashStr(node.id) + 150);
						svg += '<path d="' + roughRoundedRect(x + 3, y + 3, w - 6, h - 6, 4, loopInnerRng2) + '" fill="none" stroke="' + colors.base + '" stroke-width="1" opacity="0.5" stroke-linecap="round" filter="url(#pencil)"/>';
						svg += '<text x="' + (x + w / 2) + '" y="' + (y + 18) + '" font-size="12" fill="' + inkColor + '" font-weight="600" text-anchor="middle" font-family="' + fontFamily + '">' + escHtml(nd.label) + '</text>';
						if (hasDetails) {
							var loopToggleIcon = isLoopCollapsed ? '\\u25B6' : '\\u25BC';
							svg += '<text class="mf-collapse-toggle" data-node-id="' + escHtml(node.id) + '" x="' + (x + w - 16) + '" y="' + (y + 18) + '" font-size="10" fill="' + secondaryColor + '" text-anchor="middle" font-family="' + fontFamily + '" style="cursor:pointer">' + loopToggleIcon + '</text>';
						}
						if (!isLoopCollapsed) {
							detailLines.forEach(function(line, i) {
								var rowY = y + headerHeight + 12 + i * detailLineHeight;
								svg += '<text x="' + (x + 8) + '" y="' + rowY + '" font-size="10" fill="' + secondaryColor + '" font-family="' + fontFamily + '">' + escHtml(line) + '</text>';
							});
						} else {
							svg += '<text x="' + (x + w / 2) + '" y="' + (y + headerHeight + 12) + '" font-size="9" fill="' + secondaryColor + '" opacity="0.6" text-anchor="middle" font-family="' + fontFamily + '">' + detailLines.length + ' detail' + (detailLines.length !== 1 ? 's' : '') + '</text>';
						}
					} else {
						var isActionCollapsed = hasDetails && !!collapsedActivities[node.id];
						// Action: rounded rect with colored header + body
						if (diagramTheme === 'clean') {
							svg += '<rect x="' + x + '" y="' + y + '" width="' + w + '" height="' + h + '" rx="4" fill="var(--vscode-editor-background, #1e1e1e)" filter="url(#clean-shadow)"/>';
						}
						svg += markerFill(x, y, w, headerHeight, colors.light, makeRng(hashStr(node.id) + 50));
						var actionHeaderRng = makeRng(hashStr(node.id) + 100);
						svg += '<path d="' + roughRoundedRect(x, y, w, headerHeight, 4, actionHeaderRng) + '" fill="none" stroke="' + colors.base + '" stroke-width="1.5" stroke-linecap="round" filter="url(#pencil)"/>';
						svg += '<text x="' + (x + w / 2) + '" y="' + (y + 18) + '" font-size="12" fill="' + inkColor + '" font-weight="600" text-anchor="middle" font-family="' + fontFamily + '">' + escHtml(nd.label) + '</text>';
						if (hasDetails) {
							var actionToggleIcon = isActionCollapsed ? '\\u25B6' : '\\u25BC';
							svg += '<text class="mf-collapse-toggle" data-node-id="' + escHtml(node.id) + '" x="' + (x + w - 16) + '" y="' + (y + 18) + '" font-size="10" fill="' + secondaryColor + '" text-anchor="middle" font-family="' + fontFamily + '" style="cursor:pointer">' + actionToggleIcon + '</text>';
						}
						if (hasDetails && !isActionCollapsed) {
							var bodyY = y + headerHeight;
							var bodyH = h - headerHeight;
							if (bodyH > 2) {
								var bodyRng = makeRng(hashStr(node.id) + 200);
								svg += '<path d="' + roughRoundedRect(x, bodyY, w, bodyH, 4, bodyRng) + '" fill="none" stroke="' + colors.base + '" stroke-width="1" opacity="0.5" stroke-linecap="round" filter="url(#pencil)"/>';
							}
							detailLines.forEach(function(line, i) {
								var rowY = bodyY + 14 + i * detailLineHeight;
								var displayLine = line;
								var maxLineChars = Math.floor(w / 7 - 2);
								if (displayLine.length > maxLineChars && maxLineChars > 3) {
									displayLine = displayLine.substring(0, maxLineChars - 3) + '...';
								}
								svg += '<text x="' + (x + 8) + '" y="' + rowY + '" font-size="10" fill="' + secondaryColor + '" font-family="' + fontFamily + '">' + escHtml(displayLine) + '</text>';
							});
						} else if (hasDetails && isActionCollapsed) {
							svg += '<text x="' + (x + w / 2) + '" y="' + (y + headerHeight + 12) + '" font-size="9" fill="' + secondaryColor + '" opacity="0.6" text-anchor="middle" font-family="' + fontFamily + '">' + detailLines.length + ' detail' + (detailLines.length !== 1 ? 's' : '') + '</text>';
						}
					}

					svg += '</g>';
				});
			}

			// Render edges (under nodes)
			renderMfEdges(layout.edges, padding, padding + offsetY);

			// Render nodes
			renderMfNodes(layout.children, padding, padding + offsetY);

			// Footer (sketch only)
			if (diagramTheme === 'sketch') {
				svg += '<text x="10" y="' + (svgHeight - 8) + '" font-size="10" fill="' + secondaryColor + '" opacity="0.4" font-family="' + fontFamily + '">sketch — subject to change</text>';
			}
			svg += '</svg>';

			var canvas = document.getElementById('diagram-canvas');
			canvas.innerHTML = svg;

			// Click handlers for microflow collapse toggles
			canvas.querySelectorAll('.mf-collapse-toggle').forEach(function(el) {
				el.addEventListener('click', function(e) {
					e.stopPropagation();
					var nid = el.getAttribute('data-node-id');
					if (collapsedActivities[nid]) {
						delete collapsedActivities[nid];
					} else {
						collapsedActivities[nid] = true;
					}
					window.relayoutMicroflow();
				});
			});

			// Click handlers for node popover
			canvas.querySelectorAll('.mf-node').forEach(function(g) {
				g.addEventListener('click', function(e) {
					// Don't show popover if clicking the collapse toggle
					if (e.target.classList && e.target.classList.contains('mf-collapse-toggle')) return;
					e.stopPropagation();
					var nodeId = g.getAttribute('data-node-id');
					vscodeApi.postMessage({ type: 'nodeClicked', nodeId: nodeId });
					var nd = mfNodeMap[nodeId];
					if (!nd) return;

					var popover = document.getElementById('node-popover');
					if (!popover) return;

					var html = '<div class="popover-title">' + escHtml(nd.label) + '</div>';
					if (nd.category) {
						html += '<div class="popover-line" style="opacity:0.6;font-size:10px">' + escHtml(nd.type + ' \\u2022 ' + nd.category) + '</div>';
					}
					(nd.details || []).forEach(function(line) {
						html += '<div class="popover-line">' + escHtml(line) + '</div>';
					});
					popover.innerHTML = html;

					var vpRect = document.getElementById('diagram-viewport').getBoundingClientRect();
					var left = e.clientX - vpRect.left + 10;
					var top = e.clientY - vpRect.top + 10;
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
	`;
}
