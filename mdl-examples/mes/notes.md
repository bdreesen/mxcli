Here's the summary of issues found and fixes needed:

mxcli Skill File Improvements
Issue 1: EXTENDS syntax position not documented clearly
Problem: AI generated EXTENDS after the closing ) because mxcli syntax entity shows it that way. The actual parser requires EXTENDS BEFORE the opening (.

Files to update:

generate-domain-model.md — Add a dedicated "Entity Generalization (EXTENDS)" section with correct/wrong examples, and a note that mxcli syntax entity output is misleading
mdl-entities.md — Add a critical note to the existing Generalization section showing correct vs wrong placement
CLAUDE.md — Add a short EXTENDS quick-reference in the syntax section
Correct syntax:


-- ✅ EXTENDS before opening parenthesis
CREATE PERSISTENT ENTITY Mod.Photo EXTENDS System.Image (
  PhotoCaption: String(200)
);

-- ❌ EXTENDS after closing parenthesis = parse error
CREATE PERSISTENT ENTITY Mod.Photo (
  PhotoCaption: String(200)
) EXTENDS System.Image;
Also consider fixing the mxcli syntax entity help output itself — it currently shows EXTENDS after ) which is misleading.

Issue 2: Reserved keywords list is incomplete
Problem: Caption is a reserved keyword but wasn't listed in any of the "common reserved keywords" lists. AI used Caption as an attribute name and hit a parse error.

Files to update:

check-syntax.md — Expand the table from 9 to ~16 entries. Add: Caption, Content, Label, Range, Source, Status, Title
generate-domain-model.md — Same expansion in the "Common Reserved Keywords" bullet list
CLAUDE.md — Same expansion in the quick-reference list
All files should mention ./mxcli syntax keywords for the full 320+ list
Issue 3: system-module.md examples use reserved keyword
Problem: The System.Image usage examples use Caption: String(200) as an attribute name, which is a reserved keyword. When AI copies these examples, it gets parse errors.

File to update:

system-module.md — Rename Caption to PhotoCaption in all ProductPhoto examples (appears ~3 times), add a comment noting the rename reason
That covers the three things that tripped me up. The EXTENDS position was the biggest time sink.