// SPDX-License-Identifier: Apache-2.0

/**
 * Returns inline JavaScript code for the OQL query plan renderer.
 * Depends on shared state: data, and sketch helper functions.
 */
export function queryPlanRendererJs(): string {
	return `
		// --- OQL Query Plan renderer ---
		function renderQueryPlan(layout) {
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
			var attrLineHeight = 18;
			var svgWidth = maxX + padding * 2;
			var svgHeight = maxY + padding * 2 + titleHeight + 20;
			var fontFamily = diagramTheme === 'clean'
				? "var(--vscode-font-family, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif)"
				: "'Architects Daughter', cursive";

			// Colors
			var tableBase = '#5ba55b';
			var tableLight = '#d4edda';
			var subqueryBase = '#8e6cbf';
			var subqueryLight = '#e2d5f0';
			var joinBase = '#e89b3e';
			var joinLight = '#fce4c0';
			var filterColor = '#d9534f';
			var aggColor = '#4a90d9';
			var resultBase = '#5ba55b';
			var resultLight = '#d4edda';

			var svg = '<svg xmlns="http://www.w3.org/2000/svg" width="' + svgWidth + '" height="' + svgHeight + '">';
			svg += '<defs>' + svgFilterDefs() + '</defs>';

			// Title
			var entityName = data.entityName || 'Query Plan';
			svg += '<text x="' + padding + '" y="32" font-size="20" fill="' + inkColor + '" font-family="' + fontFamily + '">' + escHtml(entityName) + ' \\u2014 OQL Query Plan</text>';

			// PoC badge (sketch only)
			if (diagramTheme === 'sketch') {
				var badgeRng = makeRng(42);
				var titleWidth = (entityName.length + 18) * 11 + 10;
				var badgeX = padding + titleWidth;
				svg += '<path d="' + roughRoundedRect(badgeX, 17, 64, 20, 10, badgeRng) + '" fill="none" stroke="' + secondaryColor + '" stroke-width="1.2" filter="url(#pencil)"/>';
				svg += '<text x="' + (badgeX + 32) + '" y="31" font-size="9" fill="' + secondaryColor + '" text-anchor="middle" font-family="' + fontFamily + '">PoC draft</text>';
			}

			var offsetY = titleHeight;

			// Build lookups
			var tableMap = {};
			(data.tables || []).forEach(function(t) { tableMap[t.id] = t; });
			var joinMap = {};
			(data.joins || []).forEach(function(j) { joinMap[j.id] = j; });

			// Render edges
			if (layout.edges) {
				layout.edges.forEach(function(edge, idx) {
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
							svg += '<path d="' + d + '" fill="none" stroke="' + connectorColor + '" stroke-width="1.5" opacity="0.6" stroke-linecap="round" filter="url(#pencil)"/>';

							// Arrowhead at target end
							var last = pts[pts.length - 1];
							var prev = pts[pts.length - 2];
							var angle = Math.atan2(last.y - prev.y, last.x - prev.x);
							svg += roughArrowhead(last.x, last.y, angle, edgeRng);
						});
					}
				});
			}

			// Render nodes
			layout.children.forEach(function(node) {
				var x = node.x + padding;
				var y = node.y + padding + offsetY;
				var w = node.width;
				var h = node.height;
				var nodeRng = makeRng(hashStr(node.id));

				if (tableMap[node.id]) {
					// --- Table node ---
					var t = tableMap[node.id];
					var isSubquery = t.joinType === 'scalar' || (t.entity && t.entity.indexOf('(') === 0);
					var tBase = isSubquery ? subqueryBase : tableBase;
					var tLight = isSubquery ? subqueryLight : tableLight;

					// Drop shadow for clean theme
					if (diagramTheme === 'clean') {
						svg += '<rect x="' + x + '" y="' + y + '" width="' + w + '" height="' + h + '" rx="4" fill="var(--vscode-editor-background, #1e1e1e)" filter="url(#clean-shadow)"/>';
					}

					// Header fill
					svg += markerFill(x, y, w, headerHeight, tLight, makeRng(hashStr(node.id) + 50));

					// Header border
					var headerRng = makeRng(hashStr(node.id) + 100);
					svg += '<path d="' + roughRoundedRect(x, y, w, headerHeight, 4, headerRng) + '" fill="none" stroke="' + tBase + '" stroke-width="1.5" stroke-linecap="round" filter="url(#pencil)"/>';

					// Entity name + alias
					var label = t.entity;
					if (t.alias) label += ' (' + t.alias + ')';
					svg += '<text x="' + (x + w / 2) + '" y="' + (y + 18) + '" font-size="12" fill="' + inkColor + '" font-weight="600" text-anchor="middle" font-family="' + fontFamily + '">' + escHtml(label) + '</text>';

					// Body
					var bodyY = y + headerHeight;
					var bodyH = h - headerHeight;
					if (bodyH > 2) {
						var bodyRng = makeRng(hashStr(node.id) + 200);
						svg += '<path d="' + roughRoundedRect(x, bodyY, w, bodyH, 4, bodyRng) + '" fill="none" stroke="' + tBase + '" stroke-width="1" opacity="0.5" stroke-linecap="round" filter="url(#pencil)"/>';
					}

					// Attribute rows
					var rowIdx = 0;
					(t.attributes || []).forEach(function(attr) {
						var rowY = bodyY + 14 + rowIdx * attrLineHeight;
						var displayName = attr.alias || attr.name;
						if (attr.isAggregate) {
							// Aggregate: blue indicator
							svg += '<text x="' + (x + 8) + '" y="' + rowY + '" font-size="10" fill="' + aggColor + '" font-family="' + fontFamily + '">';
							svg += escHtml(attr.expression || displayName);
							svg += '</text>';
							if (attr.alias) {
								svg += '<text x="' + (x + w - 8) + '" y="' + rowY + '" font-size="9" fill="' + secondaryColor + '" text-anchor="end" font-style="italic" font-family="' + fontFamily + '">AS ' + escHtml(attr.alias) + '</text>';
							}
						} else {
							svg += '<text x="' + (x + 8) + '" y="' + rowY + '" font-size="10" fill="' + secondaryColor + '" font-family="' + fontFamily + '">';
							svg += escHtml(displayName);
							svg += '</text>';
						}
						rowIdx++;
					});

					// Filter rows
					(t.filters || []).forEach(function(f) {
						var rowY = bodyY + 14 + rowIdx * attrLineHeight;
						var display = f;
						if (display.length > 40) display = display.substring(0, 37) + '...';
						svg += '<text x="' + (x + 8) + '" y="' + rowY + '" font-size="9" fill="' + filterColor + '" font-family="' + fontFamily + '">';
						svg += '\\u2A2F ' + escHtml(display);
						svg += '</text>';
						rowIdx++;
					});
				} else if (joinMap[node.id]) {
					// --- Join node ---
					var j = joinMap[node.id];

					// Drop shadow for clean theme
					if (diagramTheme === 'clean') {
						svg += '<rect x="' + x + '" y="' + y + '" width="' + w + '" height="' + h + '" rx="6" fill="var(--vscode-editor-background, #1e1e1e)" filter="url(#clean-shadow)"/>';
					}

					// Fill
					svg += markerFill(x, y, w, h, joinLight, makeRng(hashStr(node.id) + 50));

					// Border
					var joinBorderRng = makeRng(hashStr(node.id) + 100);
					svg += '<path d="' + roughRoundedRect(x, y, w, h, 6, joinBorderRng) + '" fill="none" stroke="' + joinBase + '" stroke-width="1.5" stroke-linecap="round" filter="url(#pencil)"/>';

					// Join type label
					var joinLabel = (j.joinType || 'JOIN').toUpperCase();
					svg += '<text x="' + (x + w / 2) + '" y="' + (y + 20) + '" font-size="12" fill="' + inkColor + '" font-weight="600" text-anchor="middle" font-family="' + fontFamily + '">' + escHtml(joinLabel) + '</text>';

					// Condition text
					if (j.condition) {
						var condDisplay = j.condition;
						if (condDisplay.length > 35) condDisplay = condDisplay.substring(0, 32) + '...';
						svg += '<text x="' + (x + w / 2) + '" y="' + (y + 38) + '" font-size="9" fill="' + secondaryColor + '" opacity="0.8" text-anchor="middle" font-style="italic" font-family="' + fontFamily + '">' + escHtml(condDisplay) + '</text>';
					}
				} else if (node.id === 'result') {
					// --- Result node (view entity) ---
					// Drop shadow for clean theme
					if (diagramTheme === 'clean') {
						svg += '<rect x="' + x + '" y="' + y + '" width="' + w + '" height="' + h + '" rx="4" fill="var(--vscode-editor-background, #1e1e1e)" filter="url(#clean-shadow)"/>';
					}

					// Header fill
					svg += markerFill(x, y, w, headerHeight, resultLight, makeRng(hashStr(node.id) + 50));

					// Header border
					var resultHeaderRng = makeRng(hashStr(node.id) + 100);
					svg += '<path d="' + roughRoundedRect(x, y, w, headerHeight, 4, resultHeaderRng) + '" fill="none" stroke="' + resultBase + '" stroke-width="2" stroke-linecap="round" filter="url(#pencil)"/>';

					// Entity name
					svg += '<text x="' + (x + w / 2) + '" y="' + (y + 18) + '" font-size="13" fill="' + inkColor + '" font-weight="600" text-anchor="middle" font-family="' + fontFamily + '">' + escHtml(entityName) + '</text>';

					// Output columns
					var colBodyY = y + headerHeight;
					var colBodyH = h - headerHeight;
					if (colBodyH > 2) {
						var colBodyRng = makeRng(hashStr(node.id) + 200);
						svg += '<path d="' + roughRoundedRect(x, colBodyY, w, colBodyH, 4, colBodyRng) + '" fill="none" stroke="' + resultBase + '" stroke-width="1" opacity="0.5" stroke-linecap="round" filter="url(#pencil)"/>';
					}

					var columns = data.columns || [];
					columns.forEach(function(col, ci) {
						var colY = colBodyY + 14 + ci * attrLineHeight;
						var colLabel = col.alias || col.expression;
						svg += '<text x="' + (x + 8) + '" y="' + colY + '" font-size="10" fill="' + secondaryColor + '" font-family="' + fontFamily + '">' + escHtml(colLabel) + '</text>';
					});

					// GROUP BY badge
					if (data.groupBy) {
						var gbY = y + h + 6;
						var gbRng = makeRng(hashStr(node.id) + 300);
						var gbText = 'GROUP BY: ' + data.groupBy;
						if (gbText.length > 40) gbText = gbText.substring(0, 37) + '...';
						var gbWidth = gbText.length * 6.5 + 16;
						svg += '<path d="' + roughRoundedRect(x, gbY, gbWidth, 18, 8, gbRng) + '" fill="none" stroke="' + aggColor + '" stroke-width="1" opacity="0.7" filter="url(#pencil)"/>';
						svg += '<text x="' + (x + gbWidth / 2) + '" y="' + (gbY + 13) + '" font-size="9" fill="' + aggColor + '" text-anchor="middle" font-family="' + fontFamily + '">' + escHtml(gbText) + '</text>';
					}
				}
			});

			// Footer (sketch only)
			if (diagramTheme === 'sketch') {
				svg += '<text x="10" y="' + (svgHeight - 8) + '" font-size="10" fill="' + secondaryColor + '" opacity="0.4" font-family="' + fontFamily + '">sketch — subject to change</text>';
			}
			svg += '</svg>';

			var canvas = document.getElementById('diagram-canvas');
			canvas.innerHTML = svg;
		}
	`;
}
