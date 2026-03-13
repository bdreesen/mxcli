// SPDX-License-Identifier: Apache-2.0

/**
 * Returns inline JavaScript code for sketch-style drawing helpers.
 * Defines color constants, PRNG, rough drawing primitives, and escHtml.
 * All drawing primitives branch on the global `diagramTheme` variable
 * to produce either sketch (hand-drawn) or clean (professional) output.
 */
export function sketchHelpersJs(): string {
	return `
		// --- Sketch style colors ---
		const appBase = '#4a90d9';
		const appLight = '#d0e1f9';
		const sysBase = '#8e6cbf';
		const sysLight = '#e2d5f0';
		const inkColor = '#2c3e50';
		const secondaryColor = '#5a5a5a';
		const connectorColor = '#6b7b8d';

		// Domain model category colors
		const categoryColors = {
			persistent:    { base: '#4a90d9', light: '#d0e1f9' },
			nonpersistent: { base: '#e89b3e', light: '#fce4c0' },
			external:      { base: '#8e6cbf', light: '#e2d5f0' },
			view:          { base: '#5ba55b', light: '#d4edda' },
		};

		// Microflow category colors
		var mfCategoryColors = {
			event:       { base: '#5ba55b', light: '#d4edda' },
			object:      { base: '#4a90d9', light: '#d0e1f9' },
			retrieve:    { base: '#3db8b8', light: '#cceded' },
			call:        { base: '#8e6cbf', light: '#e2d5f0' },
			navigation:  { base: '#e89b3e', light: '#fce4c0' },
			variable:    { base: '#7f8c8d', light: '#dce0e0' },
			controlflow: { base: '#e8c93e', light: '#fdf3cc' },
			validation:  { base: '#d9534f', light: '#f5d0cf' },
			loop:        { base: '#3db8b8', light: '#cceded' },
			log:         { base: '#7f8c8d', light: '#dce0e0' },
		};

		// Page wireframe category colors
		var wireframeColors = {
			structure: { base: '#4a90d9', light: '#d0e1f9' },
			input:     { base: '#5ba55b', light: '#d4edda' },
			data:      { base: '#e89b3e', light: '#fce4c0' },
			action:    { base: '#8e6cbf', light: '#e2d5f0' },
			text:      { base: '#7f8c8d', light: '#dce0e0' },
			snippet:   { base: '#3db8b8', light: '#cceded' },
		};

		// --- Seeded PRNG ---
		function makeRng(seed) {
			var s = seed || 1;
			return function() {
				s = (s * 16807) % 2147483647;
				return (s - 1) / 2147483646;
			};
		}
		function hashStr(str) {
			var h = 0;
			for (var i = 0; i < str.length; i++) {
				h = ((h << 5) - h) + str.charCodeAt(i);
				h |= 0;
			}
			return Math.abs(h) || 1;
		}

		// --- SVG filter definitions (theme-aware) ---
		function svgFilterDefs() {
			if (diagramTheme === 'clean') {
				// No-op pencil filter (identity) so existing filter="url(#pencil)" refs are harmless
				// Plus a subtle drop-shadow for node boxes
				return '<filter id="pencil"><feOffset dx="0" dy="0"/></filter>' +
					'<filter id="marker-texture"><feOffset dx="0" dy="0"/></filter>' +
					'<filter id="clean-shadow" x="-4%" y="-4%" width="112%" height="112%"><feDropShadow dx="1" dy="2" stdDeviation="2" flood-color="#000" flood-opacity="0.15"/></filter>';
			}
			return '<filter id="pencil"><feTurbulence type="turbulence" baseFrequency="0.03" numOctaves="4" result="noise"/><feDisplacementMap in="SourceGraphic" in2="noise" scale="1.5"/></filter>' +
				'<filter id="marker-texture"><feTurbulence type="fractalNoise" baseFrequency="0.04 0.15" numOctaves="3" result="noise"/><feDisplacementMap in="SourceGraphic" in2="noise" scale="2"/><feGaussianBlur stdDeviation="0.3"/></filter>';
		}

		// --- Rough drawing helpers (theme-aware) ---
		function roughLine(x1, y1, x2, y2, rng, jitter) {
			if (diagramTheme === 'clean') {
				return 'M ' + x1.toFixed(1) + ' ' + y1.toFixed(1) + ' L ' + x2.toFixed(1) + ' ' + y2.toFixed(1);
			}
			jitter = jitter || 1.5;
			var dx = x2 - x1, dy = y2 - y1;
			var len = Math.sqrt(dx * dx + dy * dy);
			if (len < 1) return 'M ' + x1 + ' ' + y1 + ' L ' + x2 + ' ' + y2;
			var segments = Math.max(Math.ceil(len / 18), 2);
			var px = -dy / len, py = dx / len;
			var d = 'M ' + x1.toFixed(1) + ' ' + y1.toFixed(1);
			for (var i = 1; i <= segments; i++) {
				var t = i / segments;
				var x = x1 + dx * t;
				var y = y1 + dy * t;
				if (i < segments) {
					var j = (rng() - 0.5) * 2 * jitter;
					x += px * j;
					y += py * j;
				}
				d += ' L ' + x.toFixed(1) + ' ' + y.toFixed(1);
			}
			return d;
		}

		function roughRoundedRect(x, y, w, h, r, rng) {
			if (diagramTheme === 'clean') {
				// Precise rounded rectangle using arcs
				return 'M ' + (x + r) + ' ' + y +
					' L ' + (x + w - r) + ' ' + y +
					' A ' + r + ' ' + r + ' 0 0 1 ' + (x + w) + ' ' + (y + r) +
					' L ' + (x + w) + ' ' + (y + h - r) +
					' A ' + r + ' ' + r + ' 0 0 1 ' + (x + w - r) + ' ' + (y + h) +
					' L ' + (x + r) + ' ' + (y + h) +
					' A ' + r + ' ' + r + ' 0 0 1 ' + x + ' ' + (y + h - r) +
					' L ' + x + ' ' + (y + r) +
					' A ' + r + ' ' + r + ' 0 0 1 ' + (x + r) + ' ' + y +
					' Z';
			}
			var d = 'M ' + (x + r) + ' ' + y;
			// top edge
			d += roughLine(x + r, y, x + w - r, y, rng, 1.2).replace(/^M [^ ]+ [^ ]+/, '');
			d += ' Q ' + (x + w + (rng() - 0.5) * 0.8).toFixed(1) + ' ' + (y + (rng() - 0.5) * 0.8).toFixed(1) + ' ' + (x + w) + ' ' + (y + r);
			// right edge
			d += roughLine(x + w, y + r, x + w, y + h - r, rng, 1.2).replace(/^M [^ ]+ [^ ]+/, '');
			d += ' Q ' + (x + w + (rng() - 0.5) * 0.8).toFixed(1) + ' ' + (y + h + (rng() - 0.5) * 0.8).toFixed(1) + ' ' + (x + w - r) + ' ' + (y + h);
			// bottom edge
			d += roughLine(x + w - r, y + h, x + r, y + h, rng, 1.2).replace(/^M [^ ]+ [^ ]+/, '');
			d += ' Q ' + (x + (rng() - 0.5) * 0.8).toFixed(1) + ' ' + (y + h + (rng() - 0.5) * 0.8).toFixed(1) + ' ' + x + ' ' + (y + h - r);
			// left edge
			d += roughLine(x, y + h - r, x, y + r, rng, 1.2).replace(/^M [^ ]+ [^ ]+/, '');
			d += ' Q ' + (x + (rng() - 0.5) * 0.8).toFixed(1) + ' ' + (y + (rng() - 0.5) * 0.8).toFixed(1) + ' ' + (x + r) + ' ' + y;
			d += ' Z';
			return d;
		}

		function markerFill(x, y, w, h, lightColor, rng) {
			if (diagramTheme === 'clean') {
				return '<rect x="' + x + '" y="' + y + '" width="' + w + '" height="' + h + '" rx="4" fill="' + lightColor + '" opacity="0.6"/>';
			}
			var inset = 3;
			var spacing = 4.5;
			var paths = '';
			for (var pass = 0; pass < 2; pass++) {
				var sw = pass === 0 ? 4 : 3.5;
				var op = pass === 0 ? 0.65 : 0.4;
				var offsetY = pass * 1;
				for (var ly = y + inset + offsetY; ly < y + h - inset; ly += spacing) {
					var jx1 = (rng() - 0.5) * 2;
					var jx2 = (rng() - 0.5) * 2;
					var jy = (rng() - 0.5) * 4;
					var mx = x + w / 2 + jx1;
					var my = ly + jy;
					paths += '<path d="M ' + (x + inset) + ' ' + ly +
						' Q ' + mx.toFixed(1) + ' ' + my.toFixed(1) + ' ' + (x + w - inset) + ' ' + (ly + jx2).toFixed(1) +
						'" fill="none" stroke="' + lightColor + '" stroke-width="' + sw +
						'" opacity="' + op + '" stroke-linecap="round" filter="url(#marker-texture)"/>';
				}
			}
			return paths;
		}

		function roughArrowhead(tipX, tipY, angle, rng, color) {
			color = color || connectorColor;
			if (diagramTheme === 'clean') {
				var len = 10;
				var spread = 0.35;
				var x1 = tipX - len * Math.cos(angle - spread);
				var y1 = tipY - len * Math.sin(angle - spread);
				var x2 = tipX - len * Math.cos(angle + spread);
				var y2 = tipY - len * Math.sin(angle + spread);
				return '<path d="M ' + x1.toFixed(1) + ' ' + y1.toFixed(1) +
					' L ' + tipX.toFixed(1) + ' ' + tipY.toFixed(1) +
					' L ' + x2.toFixed(1) + ' ' + y2.toFixed(1) +
					'" fill="' + color + '" stroke="' + color + '" stroke-width="1" stroke-linejoin="round"/>';
			}
			var len = 10;
			var spread = 0.35;
			var j = (rng() - 0.5) * 1.5;
			var x1 = tipX - len * Math.cos(angle - spread) + j;
			var y1 = tipY - len * Math.sin(angle - spread) + j;
			var x2 = tipX - len * Math.cos(angle + spread) - j;
			var y2 = tipY - len * Math.sin(angle + spread) - j;
			return '<path d="M ' + x1.toFixed(1) + ' ' + y1.toFixed(1) +
				' L ' + tipX.toFixed(1) + ' ' + tipY.toFixed(1) +
				' L ' + x2.toFixed(1) + ' ' + y2.toFixed(1) +
				'" fill="none" stroke="' + color + '" stroke-width="1.5" stroke-linecap="round" filter="url(#pencil)"/>';
		}

		// Hollow triangle arrowhead for generalizations
		function roughTriangleArrow(tipX, tipY, angle, rng) {
			if (diagramTheme === 'clean') {
				var len = 12;
				var spread = 0.4;
				var x1 = tipX - len * Math.cos(angle - spread);
				var y1 = tipY - len * Math.sin(angle - spread);
				var x2 = tipX - len * Math.cos(angle + spread);
				var y2 = tipY - len * Math.sin(angle + spread);
				return '<path d="M ' + x1.toFixed(1) + ' ' + y1.toFixed(1) +
					' L ' + tipX.toFixed(1) + ' ' + tipY.toFixed(1) +
					' L ' + x2.toFixed(1) + ' ' + y2.toFixed(1) +
					' Z" fill="var(--vscode-editor-background, #1e1e1e)" stroke="' + connectorColor + '" stroke-width="1.5" stroke-linejoin="round"/>';
			}
			var len = 12;
			var spread = 0.4;
			var j = (rng() - 0.5) * 1;
			var x1 = tipX - len * Math.cos(angle - spread) + j;
			var y1 = tipY - len * Math.sin(angle - spread) + j;
			var x2 = tipX - len * Math.cos(angle + spread) - j;
			var y2 = tipY - len * Math.sin(angle + spread) - j;
			return '<path d="M ' + x1.toFixed(1) + ' ' + y1.toFixed(1) +
				' L ' + tipX.toFixed(1) + ' ' + tipY.toFixed(1) +
				' L ' + x2.toFixed(1) + ' ' + y2.toFixed(1) +
				' Z" fill="var(--vscode-editor-background, #1e1e1e)" stroke="' + connectorColor + '" stroke-width="1.5" stroke-linecap="round" filter="url(#pencil)"/>';
		}

		// --- Calculate node dimensions for module overview ---
		function calcNodeWidth(m) {
			var nameLen = m.name.length;
			var statsText = m.entityCount + 'E ' + m.microflowCount + 'MF ' + m.pageCount + 'P';
			var maxChars = Math.max(nameLen, statsText.length);
			return Math.max(maxChars * 9 + 30, 130);
		}
		const nodeHeight = 54;

		function escHtml(s) {
			return s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');
		}
	`;
}
