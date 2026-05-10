# Proposal: Project Brain — Persistent Knowledge and Session Scaffolding for Long-Term AI Collaboration

## Problem Statement

mxcli is developed with significant AI involvement, yet the project's knowledge infrastructure is built for static human reading rather than persistent, connected memory. Three related friction points compound over time:

1. **Documentation goes stale without detection** — Proposals describe features as future work after they have already shipped (e.g., the registry implementation predating its own proposal by weeks). There is no mechanism to detect or prevent this drift.

2. **No orientation layer for new contributors or agent sessions** — Skills files cover specific tasks well, but there is no single place that answers: *what is the current state of the project, what is actively being worked on, and what is the path for adding something new end-to-end?* A new human contributor or an AI agent dropped into a fresh session must reconstruct this from scattered files, CLAUDE.md, and git history.

3. **Knowledge is not connected** — The internals documentation in `docs-site/src/internals/` covers the pipeline and key subsystems, but the pages are isolated documents rather than a traversable network. There is no way to navigate from "executor" to "backend interface" to "MPR backend" as a linked path, and no graph view for humans to explore topics.

These are not separate problems. They share the same root: the project has grown to a point where implicit conventions and scattered documents work for a small team in a single session but break down for parallel development, onboarding, and long-running AI collaboration across many sessions.

## Current State

The project already has several knowledge layers that are good individually but unconnected:

| Layer | Location | Audience | Character |
|---|---|---|---|
| Kitchen-sink AI context | `CLAUDE.md` | AI agents | Procedural, comprehensive, monolithic |
| Maintainer task skills | `.claude/skills/*.md` | AI agents (mxcli-dev) | Procedural ("how to do X") |
| User task skills | `.claude/skills/mendix/*.md` | AI agents (shipped to users) | Procedural ("how to do X") |
| Internals docs | `docs-site/src/internals/` | Curious users + contributors | Explanatory, isolated pages |
| Proposals | `docs/11-proposals/` | Contributors | Forward-looking, no lifecycle |
| Architecture docs | `docs/01-project/`, `docs/03-development/` | Contributors | Structural, sparse |

The gaps are: links between topics (graph), decision lifecycle (why things are the way they are), extension guides (how to add new things end-to-end), and a live "what's next" layer connected to GitHub data.

## Proposed System: Three Connected Layers

### Layer 1: The Wiki (Explanatory Knowledge)

A contributor-facing wiki in `wiki/` at the repo root. The character is **explanatory** — how things work and why — distinct from skills (which are procedural) and proposals (which are decisions). One concept per file, explicit links between topics, with a Mermaid index in `README.md` that renders as a graph and makes the clusters navigable.

The `docs-site/src/internals/` pages are good raw material and would be linked from the wiki or migrated into it; the distinction is audience and framing: docs-site explains internals to curious users, the wiki explains them to contributors who need to change things.

**Proposed structure:**

```
wiki/
  README.md               # visual index: Mermaid graph of clusters + entry points
  MAP.md                  # source path → wiki topic manifest (drives freshness hooks)

  pipeline/               # the journey from MDL text to MPR write
    overview.md           # full pipeline diagram, links to each step
    grammar.md            # ANTLR4, domain split, make grammar
    ast.md                # statement node hierarchy, how kinds map to types
    visitor.md            # ANTLR listener → AST construction
    executor.md           # registry, dispatch, ExecContext, register_stubs.go
    backend.md            # hexagonal boundary, interface design, why it exists
    mpr.md                # BSON, reader/writer, storage names, v1 vs v2

  design/                 # why things are the way they are (ADRs)
    mdl-syntax-principles.md      # design guidelines, anti-patterns, decision framework
    backend-boundary.md           # no sdk/mpr in executor, why enforced by checklist
    registry-explicit-wiring.md   # why init() was rejected, emptyRegistry() testing
    grammar-domain-split.md       # how the domain split reduces conflict surface

  subsystems/             # key components in isolation
    catalog.md            # SQLite catalog, what it indexes, query interface
    lsp.md                # LSP capabilities, how diagnostics flow, wiring
    widget-engine.md      # widget registry, .def.json, template loading
    version-awareness.md  # feature registry, checkFeature(), version gates
    repl.md               # REPL architecture, session state

  extending/              # end-to-end contributor guides
    new-command.md        # full path: grammar → AST → visitor → executor → backend → MPR
    new-document-type.md  # adding a new Mendix document type
    new-backend-method.md # interface → MPR implementation → mock stub
    new-lint-rule.md      # rule registration, Starlark vs Go rules
    new-widget.md         # .def.json, widget registry, template extraction
```

**Links to skills:** Each `extending/` guide links to the relevant maintainer skill for the step-by-step checklist. The wiki provides narrative and rationale; the skill provides the checklist and gotchas.

**ADR format for `design/`:** Each decision file follows a lightweight structure:
- **Context** — what problem prompted this decision
- **Decision** — what was chosen and what was explicitly rejected
- **Status** — `accepted` | `superseded` | `under review`
- **Consequences** — what the decision implies for contributors and agents

### Layer 2: Live Project State (Connected to GitHub)

The wiki provides static knowledge. The current state of the project — what is in progress, what is ready to pick up, what is blocked — lives in GitHub (issues, milestones, project board) and should not be duplicated in the wiki. It would go stale immediately.

Instead, the connection is a **synthesized briefing** generated on demand from GitHub data and interpreted against the wiki's context layer. This is the piece that feeds contributor session planning and personal second-brain systems (e.g., meowary, which already has GitHub CLI integration).

A scheduled Claude Code agent (using the `/schedule` skill) could generate a `wiki/CURRENT.md` file periodically — not raw GitHub data, but a synthesized "given open issues, active PRs, and recent decisions, here is what the project needs next and why." This gives Meowary and other contributor systems a pull-able feed, and gives AI agents dropping into a fresh session a starting point without querying the GitHub API themselves.

GitHub Discussions with a "Decisions" category provides a free Atom/RSS feed for significant architectural decisions — subscribable by contributor second-brain systems without any additional tooling.

### Layer 3: Freshness Enforcement (Hooks and Commands)

Documentation that is not enforced degrades. Three mechanisms keep the wiki current:

**`MAP.md` — the source-to-topic manifest:**
A file that maps source paths to wiki topics:
```
mdl/executor/registry.go        → wiki/pipeline/executor.md
mdl/grammar/domains/            → wiki/pipeline/grammar.md
mdl/backend/                    → wiki/pipeline/backend.md
sdk/mpr/                        → wiki/pipeline/mpr.md
```
This manifest is what makes automated freshness checking possible.

**Hook — flag on edit:**
A `PostToolUse` hook on writes to mapped source paths looks up `MAP.md` and surfaces which wiki topics may be affected. It does not generate updates — it makes the connection visible so it is not forgotten during the session.

**Command — draft the update:**
`/mxcli-dev:update-wiki` reads the current diff, identifies affected topics via `MAP.md`, and drafts updates to the relevant pages. Invoked intentionally at the end of a feature, as part of the definition of done.

**Review integration:**
The existing `/mxcli-dev:review` command gains a wiki freshness check: for each source file changed in the PR, verify that the corresponding wiki topic has been touched or explicitly noted as unaffected. This is a review prompt, not a hard block.

A CI check (`TestWikiFreshness` or a Makefile target) could compare the modification dates of source files against their mapped wiki topics and warn when a source file is newer than its documentation by more than a threshold. This makes staleness visible at PR time rather than months later.

## Skill Level Mapping

The two-level skill structure maps cleanly to the wiki:

| | Maintainer (mxcli-dev) | User (mendix/) |
|---|---|---|
| **Procedural skills** | `.claude/skills/*.md` | `.claude/skills/mendix/*.md` |
| **Explanatory wiki** | `wiki/` (this proposal) | `docs-site/src/` (already exists) |
| **Live state** | `wiki/CURRENT.md` (generated) | GitHub Discussions feed |

User-level explanatory docs (`docs-site/`) already exist and are published. The maintainer wiki is the missing half.

## Relationship to Meowary

Retran's meowary system (https://github.com/retran/meowary) is a contributor personal second-brain with GitHub CLI integration and session-planning scaffolding. The project brain is the *supply side* that feeds meowary and equivalent systems:

- **Project brain publishes:** current state (`wiki/CURRENT.md`), decisions (GitHub Discussions Atom feed), architecture context (`wiki/`)
- **Meowary pulls:** GitHub data + project brain content → session briefing for that contributor
- **AI agents read:** wiki + CURRENT.md at session start → oriented without querying GitHub

The project brain does not replace meowary or personal second-brain systems. It provides the project-level context that personal systems lack.

## Summary and Priority

| Component | Problem solved | Dependencies | Effort |
|---|---|---|---|
| `wiki/pipeline/` | Orients contributors to how the system works | Migrate/link docs-site internals | Low |
| `wiki/extending/` | End-to-end contributor guides | Pipeline wiki | Medium |
| `wiki/design/` (ADRs) | Captures why, prevents decision re-litigation | None | Low per decision |
| `MAP.md` manifest | Enables freshness hooks and CI checks | Pipeline wiki | Low |
| Hook + `/mxcli-dev:update-wiki` | Makes wiki updates part of workflow | MAP.md | Low |
| `wiki/CURRENT.md` (scheduled agent) | Live session briefing, feeds personal second brains | GitHub API access | Medium |
| GitHub Discussions feed | RSS/Atom for decisions, subscribable by meowary | None | Low |

**Recommended order:** Start with `wiki/pipeline/` and `wiki/extending/` — these have the highest immediate value for contributor onboarding and AI sessions, and the raw material already exists in `docs-site/src/internals/`. Add `MAP.md` and the hook immediately after, since they are low effort and establish the freshness mechanism before content accumulates. ADRs can be added incrementally as decisions are made or revisited. The scheduled briefing agent and Discussions feed are independent and can be done in any order.
