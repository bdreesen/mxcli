#!/bin/bash
# Convert Markdown to Professional PDF using Pandoc and Typst
#
# Features:
#   - Professional title page
#   - Table of contents
#   - Sans-serif fonts
#   - ASCII art diagrams converted to SVG (requires svgbob)
#
# Usage: ./md2pdf.sh input.md [output.pdf]
#
# Dependencies:
#   - pandoc (markdown processing)
#   - typst (PDF generation)
#   - svgbob (optional, for ASCII art to SVG conversion)

set -e

if [ -z "$1" ]; then
    echo "Usage: $0 input.md [output.pdf]"
    exit 1
fi

INPUT="$1"
BASENAME="${INPUT%.md}"
OUTPUT="${2:-$BASENAME.pdf}"
TYPST_FILE="$BASENAME.typ"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
DIAGRAMS_DIR="$SCRIPT_DIR/diagrams"

# Check dependencies
if ! command -v pandoc &> /dev/null; then
    echo "Error: pandoc is not installed"
    echo "  Install: brew install pandoc (macOS) or apt install pandoc (Linux)"
    exit 1
fi

if ! command -v typst &> /dev/null; then
    echo "Error: typst is not installed"
    echo "  Install: cargo install typst-cli or brew install typst"
    exit 1
fi

if [ ! -f "$INPUT" ]; then
    echo "Error: Input file '$INPUT' not found"
    exit 1
fi

# Check for mermaid-cli (preferred for diagrams)
HAS_MERMAID=false
if command -v mmdc &> /dev/null; then
    HAS_MERMAID=true
    echo "mermaid-cli found - will render Mermaid diagrams to SVG"
    mkdir -p "$DIAGRAMS_DIR"
else
    echo "Note: mermaid-cli (mmdc) not found"
    echo "  Install: npm install -g @mermaid-js/mermaid-cli"
    echo "  Mermaid diagrams will be skipped"
fi

# svgbob disabled - ASCII art renders better as styled code blocks in Typst
HAS_SVGBOB=false
# if command -v svgbob &> /dev/null; then
#     HAS_SVGBOB=true
#     echo "svgbob found - will convert ASCII art diagrams to SVG"
#     mkdir -p "$DIAGRAMS_DIR"
# fi

echo "Converting $INPUT to $OUTPUT..."

# Function to extract and convert diagrams (Mermaid and ASCII art)
convert_diagrams() {
    local input_file="$1"
    local output_file="$2"
    local diagram_count=0
    local in_diagram=false
    local in_mermaid=false
    local diagram_content=""
    local line_buffer=""

    # Create a preprocessed markdown file
    > "$output_file"

    while IFS= read -r line || [[ -n "$line" ]]; do
        # Detect Mermaid code blocks
        if [[ "$line" =~ ^\`\`\`mermaid$ ]]; then
            in_mermaid=true
            diagram_content=""
            continue
        fi

        if [ "$in_mermaid" = true ]; then
            if [[ "$line" =~ ^\`\`\`$ ]]; then
                # End of mermaid block
                ((diagram_count++))
                local svg_file="$DIAGRAMS_DIR/mermaid_${diagram_count}.svg"

                if [ "$HAS_MERMAID" = true ]; then
                    # Create temp file for mermaid content
                    local mmd_file=$(mktemp).mmd
                    echo "$diagram_content" > "$mmd_file"

                    # Output as PNG for better compatibility (foreignObject in SVG not supported by Typst)
                    local png_file="$DIAGRAMS_DIR/mermaid_${diagram_count}.png"
                    mmdc -i "$mmd_file" -o "$png_file" -b white -s 3 2>/dev/null && {
                        echo "" >> "$output_file"
                        # Use relative path for the image
                        echo "![Diagram $diagram_count](diagrams/mermaid_${diagram_count}.png)" >> "$output_file"
                        echo "" >> "$output_file"
                    } || {
                        # Fallback: output as code block if mermaid fails
                        echo '```mermaid' >> "$output_file"
                        echo -n "$diagram_content" >> "$output_file"
                        echo '```' >> "$output_file"
                    }
                    rm -f "$mmd_file"
                else
                    # No mermaid, output as code block
                    echo '```mermaid' >> "$output_file"
                    echo -n "$diagram_content" >> "$output_file"
                    echo '```' >> "$output_file"
                fi
                in_mermaid=false
                diagram_content=""
            else
                diagram_content+="$line"$'\n'
            fi
            continue
        fi

        # Detect ASCII art diagram blocks (code blocks starting with box-drawing characters)
        if [[ "$line" =~ ^\`\`\`$ ]] && [ "$in_diagram" = false ]; then
            # Start of a code block - peek ahead to see if it's ASCII art
            line_buffer="$line"
            in_diagram="pending"
            diagram_content=""
            continue
        fi

        if [ "$in_diagram" = "pending" ]; then
            # Check if this looks like ASCII art (contains box-drawing characters)
            if [[ "$line" =~ [┌┐└┘│─├┤┬┴┼╔╗╚╝║═╠╣╦╩╬▲▼◄►●○■□] ]] || [[ "$line" =~ ^[[:space:]]*[\+\-\|] ]]; then
                in_diagram=true
                diagram_content="$line"$'\n'
            else
                # Not ASCII art, output the buffered line and continue
                echo "$line_buffer" >> "$output_file"
                echo "$line" >> "$output_file"
                in_diagram=false
                line_buffer=""
            fi
            continue
        fi

        if [ "$in_diagram" = true ]; then
            if [[ "$line" =~ ^\`\`\`$ ]]; then
                # End of diagram block
                ((diagram_count++))
                local svg_file="$DIAGRAMS_DIR/ascii_${diagram_count}.svg"

                if [ "$HAS_SVGBOB" = true ]; then
                    # Convert Unicode box-drawing to ASCII equivalents for svgbob
                    ascii_diagram=$(echo "$diagram_content" | sed \
                        -e 's/┌/+/g' -e 's/┐/+/g' -e 's/└/+/g' -e 's/┘/+/g' \
                        -e 's/├/+/g' -e 's/┤/+/g' -e 's/┬/+/g' -e 's/┴/+/g' -e 's/┼/+/g' \
                        -e 's/─/-/g' -e 's/│/|/g' \
                        -e 's/╔/+/g' -e 's/╗/+/g' -e 's/╚/+/g' -e 's/╝/+/g' \
                        -e 's/╠/+/g' -e 's/╣/+/g' -e 's/╦/+/g' -e 's/╩/+/g' -e 's/╬/+/g' \
                        -e 's/═/-/g' -e 's/║/|/g' \
                        -e 's/▲/^/g' -e 's/▼/v/g' -e 's/◄/</g' -e 's/►/>/g' \
                        -e 's/●/*/g' -e 's/○/o/g' -e 's/■/#/g' -e 's/□/o/g' \
                    )

                    # Convert ASCII art to SVG
                    echo "$ascii_diagram" | svgbob --output "$svg_file" 2>/dev/null || {
                        # Fallback: output as code block if svgbob fails
                        echo '```' >> "$output_file"
                        echo "$diagram_content" >> "$output_file"
                        echo '```' >> "$output_file"
                        in_diagram=false
                        continue
                    }
                    echo "" >> "$output_file"
                    # Use relative path for the image
                    echo "![Diagram $diagram_count](diagrams/ascii_${diagram_count}.svg)" >> "$output_file"
                    echo "" >> "$output_file"
                else
                    # No svgbob, output as code block
                    echo '```' >> "$output_file"
                    echo -n "$diagram_content" >> "$output_file"
                    echo '```' >> "$output_file"
                fi
                in_diagram=false
                diagram_content=""
            else
                diagram_content+="$line"$'\n'
            fi
            continue
        fi

        echo "$line" >> "$output_file"
    done < "$input_file"
}

# Preprocess markdown if mermaid or svgbob is available
PROCESSED_INPUT="$INPUT"
if [ "$HAS_MERMAID" = true ] || [ "$HAS_SVGBOB" = true ]; then
    PROCESSED_INPUT="${BASENAME}_processed.md"
    convert_diagrams "$INPUT" "$PROCESSED_INPUT"
fi

# Convert markdown to typst using pandoc
pandoc "$PROCESSED_INPUT" -o "$TYPST_FILE" \
    --to=typst \
    --wrap=none \
    --standalone \
    -V documentclass=article

# Extract metadata from markdown for title page
TITLE=$(grep -m1 '^# ' "$INPUT" | sed 's/^# //')
# Extract status line
STATUS=$(grep -m1 '^\*\*Status\*\*:' "$INPUT" | sed 's/.*: //' | sed 's/\*//g' || echo "")
# Extract date from status or use current date
DOC_DATE=$(echo "$STATUS" | grep -oE '(January|February|March|April|May|June|July|August|September|October|November|December) [0-9]{4}' || date +"%B %Y")

# Create enhanced typst file with professional formatting
TEMP_FILE=$(mktemp)
cat > "$TEMP_FILE" << 'TYPST_HEADER'
// ============================================================================
// Professional Document Template
// ============================================================================

// Page setup
#set page(
  paper: "a4",
  margin: (x: 2.5cm, y: 2.5cm),
  header: context {
    if counter(page).get().first() > 1 [
      #set text(size: 9pt, fill: gray)
      #h(1fr)
      _Mendix for Agentic IDEs_
    ]
  },
  footer: context {
    if counter(page).get().first() > 1 [
      #set text(size: 9pt, fill: gray)
      #h(1fr)
      #counter(page).display("1")
    ]
  }
)

// Typography - Sans-serif fonts (with cross-platform fallbacks)
#set text(
  font: ("Helvetica Neue", "Helvetica", "Arial"),
  size: 10.5pt,
  hyphenate: true
)

// Headings
#set heading(numbering: "1.1.1")
#show heading.where(level: 1): it => {
  pagebreak(weak: true)
  v(1.5em)
  text(size: 18pt, weight: "bold", fill: rgb("#1a365d"), it)
  v(1em)
}
#show heading.where(level: 2): it => {
  v(1.2em)
  text(size: 14pt, weight: "bold", fill: rgb("#2c5282"), it)
  v(0.6em)
}
#show heading.where(level: 3): it => {
  v(0.8em)
  text(size: 12pt, weight: "bold", fill: rgb("#2d3748"), it)
  v(0.4em)
}
#show heading.where(level: 4): it => {
  v(0.6em)
  text(size: 11pt, weight: "bold", style: "italic", it)
  v(0.3em)
}

// Code blocks - monospace with nice styling
// Detect ASCII art diagrams by looking for box-drawing characters
#show raw.where(block: true): it => {
  let content = it.text
  let is_diagram = content.contains("┌") or content.contains("└") or content.contains("│") or content.contains("─") or content.contains("╔") or content.contains("║") or content.contains("├") or content.contains("┬")

  if is_diagram {
    // Diagram styling - centered, no background tint, larger font
    set text(font: ("Menlo", "Monaco", "Courier New", "Courier"), size: 7pt)
    block(
      fill: white,
      stroke: 1pt + rgb("#e2e8f0"),
      inset: 15pt,
      radius: 6pt,
      width: 100%,
      align(center, it)
    )
  } else {
    // Regular code block styling
    set text(font: ("Menlo", "Monaco", "Courier New", "Courier"), size: 8.5pt)
    block(
      fill: rgb("#f7fafc"),
      stroke: (left: 3pt + rgb("#4299e1")),
      inset: (left: 12pt, right: 10pt, top: 10pt, bottom: 10pt),
      radius: (right: 4pt),
      width: 100%,
      it
    )
  }
}

// Inline code
#show raw.where(block: false): it => {
  set text(font: ("Menlo", "Monaco", "Courier New", "Courier"), size: 9.5pt)
  box(fill: rgb("#edf2f7"), inset: (x: 4pt, y: 2pt), radius: 3pt, it)
}

// Tables
#set table(
  stroke: (x, y) => (
    top: if y == 0 { 1.5pt + rgb("#2d3748") } else { 0.5pt + rgb("#e2e8f0") },
    bottom: if y == 0 { 1pt + rgb("#4a5568") } else { 0.5pt + rgb("#e2e8f0") },
    left: 0pt,
    right: 0pt,
  ),
  inset: 8pt,
)
#show table.cell.where(y: 0): strong

// Links
#show link: it => {
  set text(fill: rgb("#3182ce"))
  underline(it)
}

// Block quotes
#show quote: it => {
  block(
    fill: rgb("#ebf8ff"),
    stroke: (left: 4pt + rgb("#4299e1")),
    inset: (left: 16pt, right: 12pt, top: 12pt, bottom: 12pt),
    radius: (right: 4pt),
    it
  )
}

// Horizontal rule
#let horizontalrule = {
  v(1em)
  line(length: 100%, stroke: 1pt + rgb("#e2e8f0"))
  v(1em)
}

// ============================================================================
// Title Page
// ============================================================================

#page(header: none, footer: none)[
  #set align(center)

  #v(3cm)

  // Logo placeholder (optional)
  #block(
    width: 100%,
    height: 2cm,
    // Add logo here if available
  )

  #v(2cm)

  // Title
  #text(size: 32pt, weight: "bold", fill: rgb("#1a365d"))[
    Mendix for Agentic IDEs
  ]

  #v(0.5cm)

  #text(size: 18pt, fill: rgb("#4a5568"))[
    Vision & Architecture
  ]

  #v(3cm)

  // Metadata box
  #block(
    width: 70%,
    stroke: 1pt + rgb("#e2e8f0"),
    radius: 8pt,
    inset: 20pt,
  )[
    #set align(left)
    #set text(size: 11pt)

    #grid(
      columns: (1fr, 2fr),
      row-gutter: 12pt,
      [*Status:*], [Vision Document],
      [*Updated:*], [February 2026],
      [*Audience:*], [Product Strategy, Architecture, Engineering Leadership],
    )
  ]

  #v(1fr)

  // Footer
  #set text(size: 10pt, fill: rgb("#718096"))
  _This document outlines how Mendix can become the preferred target for agentic code generation of business applications._

  #v(2cm)
]

// ============================================================================
// Table of Contents
// ============================================================================

#page(header: none)[
  #heading(outlined: false, numbering: none)[Table of Contents]
  #v(1em)
  #outline(
    title: none,
    indent: 1.5em,
    depth: 4,
  )
]

// ============================================================================
// Document Content
// ============================================================================

TYPST_HEADER

# Process the pandoc output to remove title and promote all headings by one level
# Skip pandoc's preamble and extract content starting from first heading
# Remove level-1 title, then promote: == → =, === → ==, ==== → ===, etc.
# Using awk for single-pass transformation (avoids multiple promotions)
sed -n '/^=/,$p' "$TYPST_FILE" | \
    sed '/^= [^=]/d' | \
    awk '{
      if (/^====== /) { sub(/^====== /, "===== "); }
      else if (/^===== /) { sub(/^===== /, "==== "); }
      else if (/^==== /) { sub(/^==== /, "=== "); }
      else if (/^=== /) { sub(/^=== /, "== "); }
      else if (/^== /) { sub(/^== /, "= "); }
      print
    }' >> "$TEMP_FILE"

mv "$TEMP_FILE" "$TYPST_FILE"

# Compile typst to PDF
echo "Compiling PDF..."
typst compile "$TYPST_FILE" "$OUTPUT" 2>&1 | head -20

# Cleanup
if [ "$HAS_SVGBOB" = true ] && [ -f "$PROCESSED_INPUT" ]; then
    rm -f "$PROCESSED_INPUT"
fi

if [ -f "$OUTPUT" ]; then
    echo ""
    echo "✓ Created: $OUTPUT"
    echo "  Typst source: $TYPST_FILE"
    if [ -d "$DIAGRAMS_DIR" ] && [ "$(ls -A $DIAGRAMS_DIR 2>/dev/null)" ]; then
        echo "  Diagrams: $DIAGRAMS_DIR/"
    fi
else
    echo "✗ Failed to create PDF"
    exit 1
fi
