#!/bin/bash
# Convert Markdown to PDF using Pandoc and Typst
#
# Usage: ./md2pdf.sh input.md [output.pdf]
#
# If output.pdf is not specified, uses input filename with .pdf extension

set -e

if [ -z "$1" ]; then
    echo "Usage: $0 input.md [output.pdf]"
    exit 1
fi

INPUT="$1"
BASENAME="${INPUT%.md}"
OUTPUT="${2:-$BASENAME.pdf}"
TYPST_FILE="$BASENAME.typ"

# Check dependencies
if ! command -v pandoc &> /dev/null; then
    echo "Error: pandoc is not installed"
    exit 1
fi

if ! command -v typst &> /dev/null; then
    echo "Error: typst is not installed"
    exit 1
fi

if [ ! -f "$INPUT" ]; then
    echo "Error: Input file '$INPUT' not found"
    exit 1
fi

echo "Converting $INPUT to $OUTPUT..."

# Convert markdown to typst using pandoc
pandoc "$INPUT" -o "$TYPST_FILE" \
    --to=typst \
    --wrap=none \
    --standalone \
    -V documentclass=article

# Add custom typst preamble for better formatting
TEMP_FILE=$(mktemp)
cat > "$TEMP_FILE" << 'EOF'
#set document(title: none)
#set page(paper: "a4", margin: (x: 2cm, y: 2.5cm))
#set text(size: 11pt)
#set heading(numbering: "1.1")
#show heading.where(level: 1): it => {
  v(1em)
  text(size: 16pt, weight: "bold", it)
  v(0.5em)
}
#show heading.where(level: 2): it => {
  v(0.8em)
  text(size: 13pt, weight: "bold", it)
  v(0.3em)
}
#show heading.where(level: 3): it => {
  v(0.5em)
  text(size: 11pt, weight: "bold", it)
  v(0.2em)
}
#show raw.where(block: true): it => {
  set text(size: 9pt)
  block(fill: luma(245), inset: 10pt, radius: 4pt, width: 100%, it)
}
#show raw.where(block: false): it => {
  box(fill: luma(240), inset: (x: 3pt, y: 0pt), radius: 2pt, it)
}

// Define pandoc compatibility functions
#let horizontalrule = line(length: 100%, stroke: 0.5pt + gray)

EOF

# Combine preamble with pandoc output (skip pandoc's default preamble)
# Extract content after pandoc's #set commands
sed -n '/^=/,$p' "$TYPST_FILE" >> "$TEMP_FILE"
mv "$TEMP_FILE" "$TYPST_FILE"

# Compile typst to PDF
typst compile "$TYPST_FILE" "$OUTPUT"

echo "Created: $OUTPUT"
echo "Typst source: $TYPST_FILE"
