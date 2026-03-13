// ─────────────────────────────────────────────
//  Mendix-styled document wrapper
//  Usage: typst compile document.typ output.pdf
// ─────────────────────────────────────────────

// ── Config ────────────────────────────────────
#let doc-title    = "OutTracker OutSystems Application"
#let doc-subtitle = "Functionality Assessment"
#let doc-date     = "March 2026"
#let doc-author   = ""

// ── Mendix brand colors ──────────────────────
#let mendix-teal    = rgb("#00B5B0")
#let mendix-dark    = rgb("#1B1534")
#let mendix-mid     = rgb("#697480")
#let mendix-light   = rgb("#E8F5F5")
#let mendix-divider = rgb("#D0DCDC")
#let body-text       = rgb("#1A1A1A")

// ── Page setup ────────────────────────────────
#set page(
  paper: "a4",
  margin: (top: 3cm, bottom: 2.8cm, left: 2.5cm, right: 2.5cm),

  header: context {
    if counter(page).get().first() > 1 [
      #set text(size: 8pt, fill: mendix-mid)
      #grid(
        columns: (1fr, auto),
        align: (left, right),
        doc-title,
        doc-date,
      )
      #line(length: 100%, stroke: 0.5pt + mendix-teal)
    ]
  },

  footer: context {
    set text(size: 8pt, fill: mendix-mid)
    line(length: 100%, stroke: 0.5pt + mendix-divider)
    v(-0.4em)
    grid(
      columns: (1fr, auto),
      align: (left, right),
      [Mendix — Confidential],
      [Page #counter(page).display("1 of 1", both: true)],
    )
  },
)

// ── Typography ────────────────────────────────
#set text(
  font: ("Noto Sans", "Helvetica Neue", "Helvetica", "Arial"),
  size: 10.5pt,
  fill: body-text,
  lang: "en",
)

#set par(
  justify: true,
  leading: 0.75em,
  spacing: 1.2em,
)

// ── Headings ──────────────────────────────────
#show heading.where(level: 1): it => {
  pagebreak(weak: true)
  v(1.5em)
  block[
    #set text(size: 20pt, weight: "bold", fill: mendix-dark)
    #it.body
    #v(-0.3em)
    #line(length: 100%, stroke: 2.5pt + mendix-teal)
  ]
  v(0.8em)
}

#show heading.where(level: 2): it => {
  v(1.2em)
  block[
    #set text(size: 14pt, weight: "bold", fill: mendix-dark)
    #it.body
    #v(-0.2em)
    #line(length: 40%, stroke: 1.5pt + mendix-teal)
  ]
  v(0.5em)
}

#show heading.where(level: 3): it => {
  v(0.8em)
  set text(size: 11.5pt, weight: "semibold", fill: mendix-teal)
  it
  v(0.3em)
}

#show heading.where(level: 4): it => {
  v(0.5em)
  set text(size: 10.5pt, weight: "semibold", fill: mendix-mid)
  it
}

// ── Code ──────────────────────────────────────
#show raw.where(block: false): it => {
  set text(font: ("JetBrains Mono", "Fira Mono", "Courier New"), size: 9.5pt)
  box(fill: mendix-light, inset: (x: 4pt, y: 2pt), radius: 3pt, baseline: 2pt, it)
}

#show raw.where(block: true): it => {
  set text(font: ("JetBrains Mono", "Fira Mono", "Courier New"), size: 9pt)
  block(
    width: 100%,
    fill: rgb("#F0F4F4"),
    stroke: (left: 3pt + mendix-teal, rest: 0.5pt + mendix-divider),
    radius: (right: 4pt),
    inset: (x: 1em, y: 0.8em),
    it,
  )
}

// ── Tables ────────────────────────────────────
#set table(
  stroke: (x, y) => (
    top: if y == 0 { 2pt + mendix-teal } else if y == 1 { 0.5pt + mendix-divider } else { none },
    bottom: 0.5pt + mendix-divider,
  ),
  fill: (x, y) => if y == 0 { mendix-dark } else if calc.odd(y) { white } else { rgb("#F7FAFA") },
  inset: (x: 0.8em, y: 0.6em),
)

#set table.header(repeat: true)

#show table.cell.where(y: 0): set text(fill: white, weight: "bold", size: 9.5pt)
#show table: set text(size: 9.5pt)

// ── Block quotes ─────────────────────────────
#show quote.where(block: true): it => {
  block(
    width: 100%,
    fill: mendix-light,
    stroke: (left: 4pt + mendix-teal),
    inset: (left: 1.5em, right: 1em, top: 0.8em, bottom: 0.8em),
    radius: (right: 4pt),
    [#set text(style: "italic", fill: mendix-mid); #it.body],
  )
}

// ── Lists ─────────────────────────────────────
#set list(marker: ([#text(fill: mendix-teal)[▸]], [#text(fill: mendix-mid)[–]]))
#set enum(numbering: n => text(fill: mendix-teal, weight: "bold")[#n.])

// ── Links ────────────────────────────────────
#show link: it => {
  set text(fill: mendix-teal)
  underline(it)
}

// ── Cover page ───────────────────────────────
#page(margin: 0pt)[
  // Teal top band
  #block(width: 100%, height: 12pt, fill: mendix-teal)

  #v(4cm)

  #block(inset: (left: 2.5cm, right: 2.5cm))[
    // Accent bar + title
    #line(length: 60pt, stroke: 4pt + mendix-teal)
    #v(0.5em)
    #text(size: 28pt, weight: "bold", fill: mendix-dark, doc-title)
    #v(0.4em)
    #text(size: 16pt, fill: mendix-teal, doc-subtitle)
    #v(3cm)

    #grid(
      columns: (auto, 1fr),
      gutter: 1em,
      // Meta block
      block(
        fill: mendix-light,
        stroke: (left: 3pt + mendix-teal),
        inset: (x: 1.2em, y: 1em),
        radius: (right: 4pt),
        [
          #set text(size: 9.5pt)
          #if doc-author != "" [
            *Author* #h(1fr) #doc-author \
          ]
          *Date* #h(1fr) #doc-date \
          *Classification* #h(1fr) Confidential
        ]
      ),
      [],
    )
  ]

  #place(bottom)[
    #block(width: 100%, fill: mendix-dark, inset: (x: 2.5cm, y: 1.2em))[
      #set text(fill: white, size: 10pt)
      #grid(
        columns: (1fr, auto),
        align: (left, right),
        [*Mendix*],
        text(fill: mendix-teal)[mendix.com],
      )
    ]
  ]
]

// ── Table of contents ─────────────────────────
#outline(
  title: [
    #text(size: 18pt, weight: "bold", fill: mendix-dark)[Contents]
    #v(0.3em)
    #line(length: 100%, stroke: 2pt + mendix-teal)
    #v(0.5em)
  ],
  indent: 1.5em,
  depth: 3,
)

#pagebreak()

// ── Document body ─────────────────────────────
#counter(page).update(1)

#import "@preview/cmarker:0.1.8": render
#render(
  read("OutTracker_Assessment.md"),
)
