// SPDX-License-Identifier: Apache-2.0

/**
 * Returns inline JavaScript code for the page wireframe renderer.
 * Depends on shared state: wireframeColors, vscodeApi, and sketch helper functions.
 */
export function pageWireframeRendererJs(): string {
	return `
		// --- Page Wireframe renderer ---
		function renderPageWireframe(wfData) {
			var fontFamily = diagramTheme === 'clean'
				? "var(--vscode-font-family, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif)"
				: "'Architects Daughter', cursive";
			var canvasWidth = 800;
			var padding = 40;
			var gap = 6;
			var containerPad = 8;

			// Widget height defaults
			var heights = {
				textbox: 48, textarea: 80, datepicker: 48, checkbox: 32, radiobuttons: 32,
				combobox: 48, dropdown: 48, dropdownfilter: 48, textfilter: 48, numberfilter: 48, datefilter: 48,
				actionbutton: 36, title: 28, text: 22, dynamictext: 22, label: 22,
				snippetcall: 40, unknown: 32, navigationlistitem: 36,
			};

			function widgetCategory(w) {
				switch (w) {
					case 'dataview': case 'layoutgrid': case 'container': case 'groupbox': case 'tabcontainer': case 'footer': return 'structure';
					case 'textbox': case 'textarea': case 'datepicker': case 'checkbox': case 'radiobuttons': case 'combobox': case 'dropdown':
					case 'textfilter': case 'numberfilter': case 'dropdownfilter': case 'datefilter': return 'input';
					case 'datagrid': case 'gallery': case 'listview': case 'navigationlist': return 'data';
					case 'actionbutton': case 'linkbutton': return 'action';
					case 'title': case 'text': case 'dynamictext': case 'label': return 'text';
					case 'snippetcall': return 'snippet';
					default: return 'structure';
				}
			}

			// --- Pass 1: layout ---
			function layoutNode(node, x, y, availW) {
				node._x = x; node._y = y; node._w = availW;
				var w = node.widget;

				if (w === 'layoutgrid') {
					var totalH = 0;
					(node.rows || []).forEach(function(row) {
						var totalWeight = 0;
						row.columns.forEach(function(c) { totalWeight += c.weight || 1; });
						var rowMaxH = 0;
						var colX = x;
						row.columns.forEach(function(col) {
							var colW = Math.floor(availW * (col.weight || 1) / totalWeight);
							var colH = 0;
							(col.children || []).forEach(function(child, i) {
								if (i > 0) colH += gap;
								colH += layoutNode(child, colX + containerPad, y + totalH + colH + containerPad, colW - containerPad * 2);
							});
							colH += containerPad * 2;
							col._x = colX; col._y = y + totalH; col._w = colW; col._h = colH;
							if (colH > rowMaxH) rowMaxH = colH;
							colX += colW;
						});
						// Equalize column heights
						row.columns.forEach(function(col) { col._h = rowMaxH; });
						totalH += rowMaxH;
					});
					node._h = Math.max(totalH, 20);
					return node._h;
				}

				if (w === 'dataview' || w === 'groupbox') {
					var headerH = 24;
					var innerH = 0;
					(node.children || []).forEach(function(child, i) {
						if (i > 0) innerH += gap;
						innerH += layoutNode(child, x + containerPad, y + headerH + innerH + containerPad, availW - containerPad * 2);
					});
					node._h = headerH + innerH + containerPad * 2;
					return node._h;
				}

				if (w === 'tabcontainer') {
					var tabBarH = 30;
					var innerH = 0;
					// Layout first tab's children only
					var firstTab = (node.tabPages || [])[0];
					if (firstTab) {
						(firstTab.children || []).forEach(function(child, i) {
							if (i > 0) innerH += gap;
							innerH += layoutNode(child, x + containerPad, y + tabBarH + innerH + containerPad, availW - containerPad * 2);
						});
					}
					node._h = tabBarH + innerH + containerPad * 2;
					return node._h;
				}

				if (w === 'footer') {
					// Horizontal layout for buttons
					var btnH = 36;
					var btnX = x + containerPad;
					(node.children || []).forEach(function(child) {
						var bw = Math.max((child.caption || child.name || 'Button').length * 8 + 24, 80);
						child._x = btnX; child._y = y + 10; child._w = bw; child._h = btnH;
						btnX += bw + gap;
					});
					node._h = btnH + 20;
					return node._h;
				}

				if (w === 'listview') {
					var itemH = 40;
					var innerH = 0;
					// Show 3 placeholder rows
					for (var li = 0; li < 3; li++) {
						innerH += itemH + gap;
					}
					node._h = 28 + innerH;
					return node._h;
				}

				if (w === 'datagrid') {
					var headerRowH = 28;
					var dataRowH = 20;
					var numRows = 3;
					node._h = headerRowH + dataRowH * numRows + 8;
					return node._h;
				}

				if (w === 'gallery') {
					node._h = 130;
					return node._h;
				}

				if (w === 'navigationlist') {
					var nlH = 0;
					(node.children || []).forEach(function(child, i) {
						if (i > 0) nlH += 2;
						child._x = x + containerPad; child._y = y + nlH; child._w = availW - containerPad * 2; child._h = 36;
						nlH += 36;
					});
					node._h = nlH + containerPad;
					return node._h;
				}

				if (w === 'container') {
					var innerH = 0;
					(node.children || []).forEach(function(child, i) {
						if (i > 0) innerH += gap;
						innerH += layoutNode(child, x + containerPad, y + innerH + containerPad, availW - containerPad * 2);
					});
					node._h = innerH + containerPad * 2;
					return node._h;
				}

				if (w === 'actionbutton') {
					var bw = Math.max((node.caption || node.name || 'Button').length * 8 + 24, 80);
					node._w = bw;
					node._h = 36;
					return node._h;
				}

				// Leaf widgets
				node._h = heights[w] || heights.unknown;
				return node._h;
			}

			// Layout all root widgets
			var contentW = canvasWidth - padding * 2;
			var curY = 0;
			var titleBarH = 50;
			curY += titleBarH;
			(wfData.root || []).forEach(function(child, i) {
				if (i > 0) curY += gap;
				curY += layoutNode(child, padding, curY, contentW);
			});
			var totalHeight = curY + padding;

			// --- Pass 2: render SVG ---
			var svgWidth = canvasWidth;
			var svgHeight = totalHeight;
			var svg = '<svg xmlns="http://www.w3.org/2000/svg" width="' + svgWidth + '" height="' + svgHeight + '">';
			svg += '<defs>' + svgFilterDefs() + '</defs>';

			// Page title bar
			var pageName = wfData.name || 'Page';
			var pageTitle = wfData.title ? pageName + ' \\u2014 ' + wfData.title : pageName;
			var titleRng = makeRng(hashStr(pageName));
			svg += markerFill(padding, 4, contentW, 40, wireframeColors.structure.light, titleRng);
			var titleBorderRng = makeRng(hashStr(pageName) + 100);
			svg += '<path d="' + roughRoundedRect(padding, 4, contentW, 40, 6, titleBorderRng) + '" fill="none" stroke="' + wireframeColors.structure.base + '" stroke-width="2" stroke-linecap="round" filter="url(#pencil)"/>';
			svg += '<text x="' + (padding + 12) + '" y="30" font-size="14" fill="' + inkColor + '" font-weight="600" font-family="' + fontFamily + '">' + escHtml(pageTitle) + '</text>';

			// Layout info
			if (wfData.layout) {
				svg += '<text x="' + (padding + contentW - 8) + '" y="28" font-size="10" fill="' + secondaryColor + '" text-anchor="end" opacity="0.6" font-family="' + fontFamily + '">Layout: ' + escHtml(wfData.layout) + '</text>';
			}

			// PoC badge (sketch only)
			if (diagramTheme === 'sketch') {
				var badgeRng = makeRng(42);
				var badgeX = padding + contentW - 80;
				svg += '<path d="' + roughRoundedRect(badgeX, svgHeight - 22, 64, 18, 9, badgeRng) + '" fill="none" stroke="' + secondaryColor + '" stroke-width="1" filter="url(#pencil)"/>';
				svg += '<text x="' + (badgeX + 32) + '" y="' + (svgHeight - 9) + '" font-size="8" fill="' + secondaryColor + '" text-anchor="middle" font-family="' + fontFamily + '">PoC draft</text>';
			}

			function renderNode(node) {
				var x = node._x, y = node._y, w = node._w, h = node._h;
				var wt = node.widget;
				var cat = widgetCategory(wt);
				var colors = wireframeColors[cat] || wireframeColors.structure;
				var nodeRng = makeRng(hashStr(node.id));

				svg += '<g class="wf-node" data-node-id="' + escHtml(node.id) + '" style="cursor:pointer">';

				if (wt === 'layoutgrid') {
					// Render each column in each row
					(node.rows || []).forEach(function(row) {
						row.columns.forEach(function(col) {
							var colRng = makeRng(hashStr(node.id + '-col-' + col._x));
							svg += '<path d="' + roughRoundedRect(col._x, col._y, col._w, col._h, 2, colRng) + '" fill="none" stroke="' + colors.base + '" stroke-width="0.8" opacity="0.3" stroke-dasharray="4 3" filter="url(#pencil)"/>';
							(col.children || []).forEach(function(child) { renderNode(child); });
						});
					});
				} else if (wt === 'dataview') {
					// Light blue container with header
					if (diagramTheme === 'clean') {
						svg += '<rect x="' + x + '" y="' + y + '" width="' + w + '" height="' + h + '" rx="4" fill="var(--vscode-editor-background, #1e1e1e)" filter="url(#clean-shadow)"/>';
					}
					svg += markerFill(x, y, w, 24, colors.light, makeRng(hashStr(node.id) + 50));
					var dvBorderRng = makeRng(hashStr(node.id) + 100);
					svg += '<path d="' + roughRoundedRect(x, y, w, h, 4, dvBorderRng) + '" fill="none" stroke="' + colors.base + '" stroke-width="1.5" stroke-linecap="round" filter="url(#pencil)"/>';
					// Header line
					var headerLineRng = makeRng(hashStr(node.id) + 150);
					svg += '<path d="' + roughLine(x, y + 24, x + w, y + 24, headerLineRng, 0.8) + '" fill="none" stroke="' + colors.base + '" stroke-width="1" opacity="0.5" filter="url(#pencil)"/>';
					// Label
					var dvLabel = 'DataView';
					if (node.datasource) dvLabel += ': ' + node.datasource;
					svg += '<text x="' + (x + 8) + '" y="' + (y + 16) + '" font-size="10" fill="' + colors.base + '" font-weight="600" font-family="' + fontFamily + '">' + escHtml(dvLabel) + '</text>';
					(node.children || []).forEach(function(child) { renderNode(child); });

				} else if (wt === 'groupbox') {
					var gbBorderRng = makeRng(hashStr(node.id) + 100);
					svg += '<path d="' + roughRoundedRect(x, y, w, h, 4, gbBorderRng) + '" fill="none" stroke="' + colors.base + '" stroke-width="1.2" stroke-linecap="round" filter="url(#pencil)"/>';
					// Caption overlaid on top edge (fieldset style)
					if (node.caption) {
						var captionW = node.caption.length * 7 + 16;
						svg += '<rect x="' + (x + 10) + '" y="' + (y - 7) + '" width="' + captionW + '" height="14" fill="var(--vscode-editor-background, #1e1e1e)"/>';
						svg += '<text x="' + (x + 18) + '" y="' + (y + 4) + '" font-size="11" fill="' + colors.base + '" font-weight="600" font-family="' + fontFamily + '">' + escHtml(node.caption) + '</text>';
					}
					(node.children || []).forEach(function(child) { renderNode(child); });

				} else if (wt === 'tabcontainer') {
					var tcBorderRng = makeRng(hashStr(node.id) + 100);
					svg += '<path d="' + roughRoundedRect(x, y, w, h, 4, tcBorderRng) + '" fill="none" stroke="' + colors.base + '" stroke-width="1.2" stroke-linecap="round" filter="url(#pencil)"/>';
					// Tab bar
					var tabX = x;
					(node.tabPages || []).forEach(function(tp, tIdx) {
						var tabW = Math.max((tp.caption || 'Tab').length * 8 + 20, 60);
						var tabRng = makeRng(hashStr(node.id + '-tab-' + tIdx));
						if (tIdx === 0) {
							svg += markerFill(tabX, y, tabW, 28, colors.light, makeRng(hashStr(node.id) + 200 + tIdx));
						}
						svg += '<path d="' + roughRoundedRect(tabX, y, tabW, 28, 3, tabRng) + '" fill="none" stroke="' + colors.base + '" stroke-width="1" filter="url(#pencil)"/>';
						svg += '<text x="' + (tabX + tabW / 2) + '" y="' + (y + 18) + '" font-size="10" fill="' + inkColor + '" text-anchor="middle" font-family="' + fontFamily + '">' + escHtml(tp.caption || 'Tab') + '</text>';
						tabX += tabW + 2;
					});
					// Render first tab's children
					var firstTab = (node.tabPages || [])[0];
					if (firstTab) {
						(firstTab.children || []).forEach(function(child) { renderNode(child); });
					}

				} else if (wt === 'textbox' || wt === 'combobox' || wt === 'dropdown' || wt === 'datepicker' ||
				           wt === 'textfilter' || wt === 'numberfilter' || wt === 'dropdownfilter' || wt === 'datefilter') {
					// Label above input
					if (node.label) {
						svg += '<text x="' + x + '" y="' + (y + 12) + '" font-size="11" fill="' + inkColor + '" font-family="' + fontFamily + '">' + escHtml(node.label) + '</text>';
					}
					var inputY = node.label ? y + 18 : y + 4;
					var inputH = h - (node.label ? 22 : 8);
					var inputRng = makeRng(hashStr(node.id) + 100);
					svg += '<path d="' + roughRoundedRect(x, inputY, w, inputH, 3, inputRng) + '" fill="none" stroke="' + colors.base + '" stroke-width="1.2" stroke-linecap="round" filter="url(#pencil)"/>';
					// Placeholder text
					var placeholder = node.binding || '';
					if (placeholder) {
						svg += '<text x="' + (x + 8) + '" y="' + (inputY + inputH / 2 + 4) + '" font-size="10" fill="' + secondaryColor + '" opacity="0.5" font-family="' + fontFamily + '">[' + escHtml(placeholder) + ']</text>';
					}
					// Dropdown chevron
					if (wt === 'combobox' || wt === 'dropdown' || wt === 'dropdownfilter') {
						svg += '<text x="' + (x + w - 16) + '" y="' + (inputY + inputH / 2 + 4) + '" font-size="12" fill="' + secondaryColor + '" opacity="0.5">&#x25BE;</text>';
					}
					// Calendar icon for datepicker
					if (wt === 'datepicker' || wt === 'datefilter') {
						var cx = x + w - 20, cy = inputY + 4;
						var calRng = makeRng(hashStr(node.id) + 200);
						svg += '<path d="' + roughRoundedRect(cx, cy, 14, 14, 2, calRng) + '" fill="none" stroke="' + secondaryColor + '" stroke-width="0.8" opacity="0.4" filter="url(#pencil)"/>';
						svg += '<path d="' + roughLine(cx, cy + 5, cx + 14, cy + 5, calRng, 0.5) + '" fill="none" stroke="' + secondaryColor + '" stroke-width="0.6" opacity="0.4"/>';
					}

				} else if (wt === 'textarea') {
					if (node.label) {
						svg += '<text x="' + x + '" y="' + (y + 12) + '" font-size="11" fill="' + inkColor + '" font-family="' + fontFamily + '">' + escHtml(node.label) + '</text>';
					}
					var inputY = node.label ? y + 18 : y + 4;
					var inputH = h - (node.label ? 22 : 8);
					var taRng = makeRng(hashStr(node.id) + 100);
					svg += '<path d="' + roughRoundedRect(x, inputY, w, inputH, 3, taRng) + '" fill="none" stroke="' + colors.base + '" stroke-width="1.2" stroke-linecap="round" filter="url(#pencil)"/>';
					// Wavy placeholder lines
					for (var li = 0; li < 3; li++) {
						var lineY = inputY + 12 + li * 14;
						if (lineY + 6 > inputY + inputH) break;
						var lineW = w * (0.3 + Math.random() * 0.5);
						var lineRng = makeRng(hashStr(node.id) + 150 + li);
						svg += '<path d="' + roughLine(x + 8, lineY, x + 8 + lineW, lineY, lineRng, 0.6) + '" fill="none" stroke="' + secondaryColor + '" stroke-width="1.5" opacity="0.2" stroke-linecap="round"/>';
					}

				} else if (wt === 'checkbox') {
					var cbRng = makeRng(hashStr(node.id) + 100);
					svg += '<path d="' + roughRoundedRect(x, y + 6, 16, 16, 2, cbRng) + '" fill="none" stroke="' + colors.base + '" stroke-width="1.2" filter="url(#pencil)"/>';
					var lbl = node.label || node.binding || '';
					if (lbl) {
						svg += '<text x="' + (x + 22) + '" y="' + (y + 19) + '" font-size="11" fill="' + inkColor + '" font-family="' + fontFamily + '">' + escHtml(lbl) + '</text>';
					}

				} else if (wt === 'radiobuttons') {
					var rbLbl = node.label || node.binding || '';
					if (rbLbl) {
						svg += '<text x="' + x + '" y="' + (y + 12) + '" font-size="11" fill="' + inkColor + '" font-family="' + fontFamily + '">' + escHtml(rbLbl) + '</text>';
					}
					// Two radio circles
					for (var ri = 0; ri < 2; ri++) {
						var rcx = x + 8; var rcy = y + 22 + ri * 16;
						if (rcy + 8 > y + h) break;
						svg += '<circle cx="' + rcx + '" cy="' + rcy + '" r="6" fill="none" stroke="' + colors.base + '" stroke-width="1" filter="url(#pencil)"/>';
						svg += '<text x="' + (rcx + 12) + '" y="' + (rcy + 4) + '" font-size="10" fill="' + secondaryColor + '" font-family="' + fontFamily + '">Option ' + (ri + 1) + '</text>';
					}

				} else if (wt === 'actionbutton') {
					var isPrimary = (node.style || '').toLowerCase() === 'primary';
					if (isPrimary) {
						svg += markerFill(x, y, w, h, colors.light, makeRng(hashStr(node.id) + 50));
					}
					var btnRng = makeRng(hashStr(node.id) + 100);
					svg += '<path d="' + roughRoundedRect(x, y, w, h, 4, btnRng) + '" fill="none" stroke="' + colors.base + '" stroke-width="' + (isPrimary ? '2' : '1.2') + '" stroke-linecap="round" filter="url(#pencil)"/>';
					svg += '<text x="' + (x + w / 2) + '" y="' + (y + h / 2 + 4) + '" font-size="11" fill="' + (isPrimary ? inkColor : secondaryColor) + '" font-weight="' + (isPrimary ? '600' : '400') + '" text-anchor="middle" font-family="' + fontFamily + '">' + escHtml(node.caption || node.name || 'Button') + '</text>';

				} else if (wt === 'datagrid') {
					if (diagramTheme === 'clean') {
						svg += '<rect x="' + x + '" y="' + y + '" width="' + w + '" height="' + h + '" rx="4" fill="var(--vscode-editor-background, #1e1e1e)" filter="url(#clean-shadow)"/>';
					}
					var dgBorderRng = makeRng(hashStr(node.id) + 100);
					svg += '<path d="' + roughRoundedRect(x, y, w, h, 4, dgBorderRng) + '" fill="none" stroke="' + colors.base + '" stroke-width="1.5" stroke-linecap="round" filter="url(#pencil)"/>';
					// Header
					svg += markerFill(x + 1, y + 1, w - 2, 26, colors.light, makeRng(hashStr(node.id) + 50));
					var hdrLineRng = makeRng(hashStr(node.id) + 150);
					svg += '<path d="' + roughLine(x, y + 28, x + w, y + 28, hdrLineRng, 0.8) + '" fill="none" stroke="' + colors.base + '" stroke-width="1" filter="url(#pencil)"/>';
					// Column headers
					var cols = node.columns || [];
					if (cols.length > 0) {
						var colW = (w - 4) / cols.length;
						cols.forEach(function(col, ci) {
							var cx = x + 2 + ci * colW;
							svg += '<text x="' + (cx + 6) + '" y="' + (y + 18) + '" font-size="10" fill="' + inkColor + '" font-weight="600" font-family="' + fontFamily + '">' + escHtml(col.caption || col.binding || 'Col') + '</text>';
							if (ci > 0) {
								var sepRng = makeRng(hashStr(node.id) + 300 + ci);
								svg += '<path d="' + roughLine(cx, y + 2, cx, y + h - 2, sepRng, 0.5) + '" fill="none" stroke="' + colors.base + '" stroke-width="0.6" opacity="0.3"/>';
							}
						});
					}
					// Data rows (wavy placeholder lines)
					for (var dr = 0; dr < 3; dr++) {
						var rowY = y + 30 + dr * 20;
						if (rowY + 14 > y + h) break;
						if (cols.length > 0) {
							var colW = (w - 4) / cols.length;
							cols.forEach(function(col, ci) {
								var lineW = colW * (0.3 + Math.random() * 0.4);
								var dataRng = makeRng(hashStr(node.id) + 400 + dr * 10 + ci);
								svg += '<path d="' + roughLine(x + 2 + ci * colW + 6, rowY + 10, x + 2 + ci * colW + 6 + lineW, rowY + 10, dataRng, 0.5) + '" fill="none" stroke="' + secondaryColor + '" stroke-width="1.5" opacity="0.2" stroke-linecap="round"/>';
							});
						}
						// Row separator
						if (dr < 2) {
							var rowLineRng = makeRng(hashStr(node.id) + 500 + dr);
							svg += '<path d="' + roughLine(x + 2, rowY + 18, x + w - 2, rowY + 18, rowLineRng, 0.5) + '" fill="none" stroke="' + secondaryColor + '" stroke-width="0.5" opacity="0.15"/>';
						}
					}

				} else if (wt === 'listview') {
					if (diagramTheme === 'clean') {
						svg += '<rect x="' + x + '" y="' + y + '" width="' + w + '" height="' + h + '" rx="4" fill="var(--vscode-editor-background, #1e1e1e)" filter="url(#clean-shadow)"/>';
					}
					var lvBorderRng = makeRng(hashStr(node.id) + 100);
					svg += '<path d="' + roughRoundedRect(x, y, w, h, 4, lvBorderRng) + '" fill="none" stroke="' + colors.base + '" stroke-width="1.2" stroke-linecap="round" filter="url(#pencil)"/>';
					// Header
					svg += markerFill(x + 1, y + 1, w - 2, 26, colors.light, makeRng(hashStr(node.id) + 50));
					var dsLabel = node.datasource || 'ListView';
					svg += '<text x="' + (x + 8) + '" y="' + (y + 18) + '" font-size="10" fill="' + colors.base + '" font-weight="600" font-family="' + fontFamily + '">' + escHtml(dsLabel) + '</text>';
					// Placeholder rows
					for (var lr = 0; lr < 3; lr++) {
						var rowY = y + 30 + lr * 40;
						if (rowY + 30 > y + h) break;
						var rowRng = makeRng(hashStr(node.id) + 200 + lr);
						svg += '<path d="' + roughRoundedRect(x + 6, rowY, w - 12, 34, 3, rowRng) + '" fill="none" stroke="' + secondaryColor + '" stroke-width="0.8" opacity="0.25" stroke-dasharray="3 2" filter="url(#pencil)"/>';
						var lineW = (w - 24) * (0.3 + Math.random() * 0.4);
						var lineRng = makeRng(hashStr(node.id) + 250 + lr);
						svg += '<path d="' + roughLine(x + 14, rowY + 18, x + 14 + lineW, rowY + 18, lineRng, 0.5) + '" fill="none" stroke="' + secondaryColor + '" stroke-width="1.5" opacity="0.15" stroke-linecap="round"/>';
					}

				} else if (wt === 'gallery') {
					if (diagramTheme === 'clean') {
						svg += '<rect x="' + x + '" y="' + y + '" width="' + w + '" height="' + h + '" rx="4" fill="var(--vscode-editor-background, #1e1e1e)" filter="url(#clean-shadow)"/>';
					}
					var galBorderRng = makeRng(hashStr(node.id) + 100);
					svg += '<path d="' + roughRoundedRect(x, y, w, h, 4, galBorderRng) + '" fill="none" stroke="' + colors.base + '" stroke-width="1.2" stroke-linecap="round" filter="url(#pencil)"/>';
					// Header
					svg += markerFill(x + 1, y + 1, w - 2, 26, colors.light, makeRng(hashStr(node.id) + 50));
					var galLabel = node.datasource || 'Gallery';
					svg += '<text x="' + (x + 8) + '" y="' + (y + 18) + '" font-size="10" fill="' + colors.base + '" font-weight="600" font-family="' + fontFamily + '">' + escHtml(galLabel) + '</text>';
					// 2x2 card grid
					var cardW = (w - 20) / 2;
					var cardH = 44;
					for (var gr = 0; gr < 2; gr++) {
						for (var gc = 0; gc < 2; gc++) {
							var cx = x + 6 + gc * (cardW + 6);
							var cy = y + 32 + gr * (cardH + 6);
							var cardRng = makeRng(hashStr(node.id) + 200 + gr * 2 + gc);
							svg += '<path d="' + roughRoundedRect(cx, cy, cardW, cardH, 3, cardRng) + '" fill="none" stroke="' + secondaryColor + '" stroke-width="0.8" opacity="0.3" stroke-dasharray="3 2" filter="url(#pencil)"/>';
						}
					}

				} else if (wt === 'title') {
					svg += '<text x="' + x + '" y="' + (y + 20) + '" font-size="16" fill="' + inkColor + '" font-weight="600" font-family="' + fontFamily + '">' + escHtml(node.caption || node.content || '') + '</text>';

				} else if (wt === 'text' || wt === 'dynamictext' || wt === 'label') {
					svg += '<text x="' + x + '" y="' + (y + 15) + '" font-size="11" fill="' + secondaryColor + '" font-family="' + fontFamily + '">' + escHtml(node.content || node.caption || '') + '</text>';

				} else if (wt === 'snippetcall') {
					var snipRng = makeRng(hashStr(node.id) + 100);
					svg += '<path d="' + roughRoundedRect(x, y, w, h, 3, snipRng) + '" fill="none" stroke="' + colors.base + '" stroke-width="1" stroke-dasharray="6 3" stroke-linecap="round" filter="url(#pencil)"/>';
					svg += '<text x="' + (x + w / 2) + '" y="' + (y + h / 2 + 4) + '" font-size="10" fill="' + colors.base + '" text-anchor="middle" font-family="' + fontFamily + '">\\u27E8Snippet: ' + escHtml(node.content || '?') + '\\u27E9</text>';

				} else if (wt === 'footer') {
					// Separator line then children
					var footLineRng = makeRng(hashStr(node.id) + 100);
					svg += '<path d="' + roughLine(x, y + 4, x + w, y + 4, footLineRng, 0.8) + '" fill="none" stroke="' + secondaryColor + '" stroke-width="1" opacity="0.3" filter="url(#pencil)"/>';
					(node.children || []).forEach(function(child) { renderNode(child); });

				} else if (wt === 'navigationlist') {
					var nlBorderRng = makeRng(hashStr(node.id) + 100);
					svg += '<path d="' + roughRoundedRect(x, y, w, h, 4, nlBorderRng) + '" fill="none" stroke="' + colors.base + '" stroke-width="1.2" stroke-linecap="round" filter="url(#pencil)"/>';
					(node.children || []).forEach(function(child) { renderNode(child); });

				} else if (wt === 'navigationlistitem') {
					var nliRng = makeRng(hashStr(node.id) + 100);
					svg += '<path d="' + roughRoundedRect(x, y, w, h, 3, nliRng) + '" fill="none" stroke="' + colors.base + '" stroke-width="0.8" opacity="0.4" filter="url(#pencil)"/>';
					// Show children or a placeholder
					if ((node.children || []).length > 0) {
						(node.children || []).forEach(function(child) { renderNode(child); });
					} else {
						svg += '<text x="' + (x + 8) + '" y="' + (y + h / 2 + 4) + '" font-size="10" fill="' + secondaryColor + '" font-family="' + fontFamily + '">' + escHtml(node.action || 'Item') + '</text>';
					}

				} else if (wt === 'container') {
					// Transparent container, just render children
					(node.children || []).forEach(function(child) { renderNode(child); });

				} else {
					// Unknown widget — dashed rect with type name
					var unkRng = makeRng(hashStr(node.id) + 100);
					svg += '<path d="' + roughRoundedRect(x, y, w, h, 3, unkRng) + '" fill="none" stroke="' + secondaryColor + '" stroke-width="0.8" stroke-dasharray="4 3" opacity="0.4" filter="url(#pencil)"/>';
					svg += '<text x="' + (x + 8) + '" y="' + (y + h / 2 + 4) + '" font-size="9" fill="' + secondaryColor + '" opacity="0.5" font-family="' + fontFamily + '">' + escHtml(wt + (node.name ? ' ' + node.name : '')) + '</text>';
				}

				svg += '</g>';
			}

			// Render all root nodes
			(wfData.root || []).forEach(function(node) { renderNode(node); });

			// Footer (sketch only)
			if (diagramTheme === 'sketch') {
				svg += '<text x="10" y="' + (svgHeight - 6) + '" font-size="10" fill="' + secondaryColor + '" opacity="0.4" font-family="' + fontFamily + '">sketch — subject to change</text>';
			}
			svg += '</svg>';

			var canvas = document.getElementById('diagram-canvas');
			canvas.innerHTML = svg;

			// Click handlers for wireframe nodes
			canvas.querySelectorAll('.wf-node').forEach(function(g) {
				g.addEventListener('click', function(e) {
					e.stopPropagation();
					var nodeId = g.getAttribute('data-node-id');
					vscodeApi.postMessage({ type: 'nodeClicked', nodeId: nodeId });

					// Find node in data tree for popover
					var nd = findWireframeNode(wfData.root || [], nodeId);
					if (!nd) return;

					var popover = document.getElementById('node-popover');
					if (!popover) return;

					var html = '<div class="popover-title">' + escHtml(nd.widget + (nd.name ? ' ' + nd.name : '')) + '</div>';
					if (nd.label) html += '<div class="popover-line">Label: ' + escHtml(nd.label) + '</div>';
					if (nd.binding) html += '<div class="popover-line">Attribute: ' + escHtml(nd.binding) + '</div>';
					if (nd.datasource) html += '<div class="popover-line">Source: ' + escHtml(nd.datasource) + '</div>';
					if (nd.action) html += '<div class="popover-line">Action: ' + escHtml(nd.action) + '</div>';
					if (nd.caption) html += '<div class="popover-line">Caption: ' + escHtml(nd.caption) + '</div>';
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

		function findWireframeNode(nodes, id) {
			for (var i = 0; i < nodes.length; i++) {
				if (nodes[i].id === id) return nodes[i];
				var found = findWireframeNode(nodes[i].children || [], id);
				if (found) return found;
				// Search in rows
				if (nodes[i].rows) {
					for (var r = 0; r < nodes[i].rows.length; r++) {
						for (var c = 0; c < nodes[i].rows[r].columns.length; c++) {
							found = findWireframeNode(nodes[i].rows[r].columns[c].children || [], id);
							if (found) return found;
						}
					}
				}
				// Search in tabPages
				if (nodes[i].tabPages) {
					for (var t = 0; t < nodes[i].tabPages.length; t++) {
						found = findWireframeNode(nodes[i].tabPages[t].children || [], id);
						if (found) return found;
					}
				}
			}
			return null;
		}
	`;
}
