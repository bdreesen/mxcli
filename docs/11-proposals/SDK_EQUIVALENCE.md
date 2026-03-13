# Mendix Model SDK Go - Equivalence Analysis

This document provides a comprehensive analysis of the original Mendix Model SDK (TypeScript) and the approach for creating a fully equivalent Go implementation.

## Table of Contents

1. [Overview](#overview)
2. [Original SDK Architecture](#original-sdk-architecture)
3. [Current Go Implementation](#current-go-implementation)
4. [Gap Analysis](#gap-analysis)
5. [Code Generation Strategy](#code-generation-strategy)
6. [Maintenance Process](#maintenance-process)
7. [Implementation Roadmap](#implementation-roadmap)

---

## Overview

The goal is to create a 100% functional equivalent of the Mendix Model SDK for local Mendix projects in Go. The original SDK consists of two main components:

| Component | Purpose | Go Equivalent |
|-----------|---------|---------------|
| **mendixmodelsdk** | High-level OO API for model manipulation | `modelsdk-go` packages |
| **mendixmodellib** | Low-level format handling, metamodel definitions | `mpr/` package + generated types |

### Key Differences

- **Original SDK**: Cloud-first, connects to Mendix Team Server, real-time collaboration
- **Go Implementation**: Local-first, works directly with `.mpr` files, no cloud dependency

---

## Original SDK Architecture

### mendixmodelsdk (v4.105.0)

The high-level SDK providing:

- **52 metamodel domains** with 820+ type definitions per version
- **Delta-based change system** for CRUD operations with undo/redo
- **Cloud connectivity** for Team Server (SVN/Git)
- **Real-time synchronization** via Server-Sent Events
- **Working copy management** with locking and collaboration

#### Key Domains

| Domain | Description | Types |
|--------|-------------|-------|
| `domainmodels` | Entities, attributes, associations | ~50 |
| `microflows` | Microflow activities and flows | 150+ |
| `pages` | UI widgets, layouts, snippets | 500+ |
| `workflows` | BPM workflow definitions | 200+ |
| `rest` | REST/OData services | ~40 |
| `webservices` | SOAP web services | ~30 |
| `security` | Access rules, user roles | ~25 |
| `navigation` | Menus and navigation | ~20 |
| `mappings` | Import/export mappings | ~30 |
| `expressions` | Expression language | ~100 |
| ... | 42 more domains | ... |

### mendixmodellib (v1.76.1)

The low-level library providing:

- **MPR file format handling** (SQLite + BSON)
- **Metamodel reflection data** for 88 Mendix versions (6.0.0 - 11.6.0)
- **Delta processing** for model transformations
- **GUID utilities** for binary encoding
- **Version compatibility** checking

#### Reflection Data Structure

Located in `libs/mendixmodellib/reflection-data/`:

```
reflection-data/
├── 6.0.0-structures.json      # Mendix 6.0.0 type definitions
├── 6.0.0-storageNames.json    # Storage name mappings
├── ...
├── 11.6.0-structures.json     # Latest version
└── 11.6.0-storageNames.json
```

Each `{version}-structures.json` contains:

```json
{
  "DomainModels$Entity": {
    "qualifiedName": "DomainModels$Entity",
    "storageName": "DomainModels$Entity",
    "superTypeName": "DomainModels$MaybeRemotableElement",
    "abstract": false,
    "type": "ELEMENT",
    "properties": {
      "name": {
        "name": "name",
        "storageName": "Name",
        "list": false,
        "typeInfo": {
          "type": "PRIMITIVE",
          "primitiveType": "STRING"
        }
      },
      "attributes": {
        "name": "attributes",
        "storageName": "Attributes",
        "list": true,
        "typeInfo": {
          "type": "ELEMENT",
          "elementType": "DomainModels$Attribute",
          "kind": "PART"
        }
      }
      // ... more properties
    },
    "defaultSettings": {
      "name": "",
      "documentation": ""
    }
  }
}
```

---

## Current Go Implementation

### Package Structure

```
modelsdk-go/
├── modelsdk.go           # Public API entry points
├── model/                # Core types (ID, Module, Project, etc.)
├── domainmodel/          # Entity, Attribute, Association types
├── microflows/           # Microflow, Nanoflow types
├── pages/                # Page, Layout, Snippet types
├── mpr/                  # MPR file reading/writing
│   ├── reader.go         # Read-only access
│   ├── writer.go         # Read-write access
│   ├── parser.go         # BSON parsing
│   └── utils.go          # UUID generation, etc.
└── examples/             # Usage examples
```

### Implemented Features

| Feature | Status | Notes |
|---------|--------|-------|
| MPR v1/v2 reading | ✅ Complete | SQLite + BSON parsing |
| MPR v1/v2 writing | ✅ Complete | Transaction support |
| Module CRUD | ✅ Complete | |
| Entity CRUD | ✅ Complete | |
| Attribute types | ✅ Complete | 9 types |
| Association CRUD | ✅ Complete | |
| Microflow basic | ⚠️ Partial | Basic structure only |
| Page basic | ⚠️ Partial | Basic structure only |
| JSON export | ✅ Complete | |

### Type Coverage Comparison

| Domain | TypeScript Lines | Go Lines | Coverage |
|--------|-----------------|----------|----------|
| domainmodels | 3,597 | ~478 | ~13% |
| microflows | 6,628 | ~200 | ~3% |
| pages | 34,062 | ~100 | <1% |
| workflows | 3,471 | 0 | 0% |
| rest | 2,453 | 0 | 0% |
| expressions | 7,591 | 0 | 0% |

---

## Gap Analysis

### Missing Metamodel Domains

48 of 52 domains are not implemented:

**High Priority:**
- `workflows` - Workflow/BPM definitions
- `rest` - REST service consumption/publishing
- `webservices` - SOAP web services
- `security` - Access rules, module roles
- `navigation` - Navigation profiles, menus
- `expressions` - Expression language types

**Medium Priority:**
- `mappings`, `importmappings`, `exportmappings`
- `jsonstructures`, `xmlschemas`
- `javaactions`, `javascriptactions`
- `customwidgets`
- `reports`, `datasets`
- `documenttemplates`

**Lower Priority:**
- `kafka`, `queues`, `businessevents`
- `appservices`
- `mlmappings`
- `databaseconnector`
- `exceldataimporter`
- ... and more

### Missing Type Details

Even in implemented domains, many subtypes are missing:

**Microflows (150+ activities needed):**
- `ActionActivity`, `LoopedActivity`
- `CreateObjectAction`, `ChangeObjectAction`, `DeleteAction`
- `RetrieveAction`, `AggregateListAction`
- `MicroflowCallAction`, `JavaActionCallAction`
- `ShowMessageAction`, `ShowPageAction`
- `LogMessageAction`, `ValidationFeedbackAction`
- ... and 140+ more

**Pages (500+ widgets needed):**
- Container widgets: `DataView`, `ListView`, `DataGrid`, `TemplateGrid`
- Input widgets: `TextBox`, `TextArea`, `DropDown`, `DatePicker`
- Button widgets: `ActionButton`, `LinkButton`
- Layout widgets: `Container`, `GroupBox`, `TabContainer`
- ... and 480+ more

### Missing Infrastructure

| Component | Description | Status |
|-----------|-------------|--------|
| Metamodel reflection | Runtime type inspection | ❌ Not implemented |
| Version compatibility | Check type availability per version | ❌ Not implemented |
| Delta system | Change tracking with undo/redo | ❌ Not implemented |
| Property system | Type-safe property access | ❌ Not implemented |
| Qualified names | Path-based element resolution | ⚠️ Basic only |

---

## Code Generation Strategy

### Recommended Approach: Generate from Reflection Data

The `mendixmodellib/reflection-data/` contains complete metamodel definitions in JSON format. We can generate Go code automatically from these definitions.

### Benefits

1. **Accuracy**: Generated types match the official metamodel exactly
2. **Maintainability**: New Mendix versions only require running the generator
3. **Completeness**: All 820+ types per version can be generated
4. **Version support**: Can generate version-specific types or union types

### Generator Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Code Generator                            │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐  │
│  │   Parser     │───▶│  Transformer │───▶│   Emitter    │  │
│  │              │    │              │    │              │  │
│  │ Read JSON    │    │ Build type   │    │ Generate Go  │  │
│  │ metamodel    │    │ hierarchy    │    │ source code  │  │
│  └──────────────┘    └──────────────┘    └──────────────┘  │
│         │                   │                   │           │
│         ▼                   ▼                   ▼           │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐  │
│  │ structures   │    │ Go type      │    │ .go files    │  │
│  │ .json        │    │ definitions  │    │ per domain   │  │
│  └──────────────┘    └──────────────┘    └──────────────┘  │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

### Generator Implementation

```go
// cmd/generate/main.go
package main

// Generator reads reflection data and produces Go types
type Generator struct {
    Version     string           // Target Mendix version
    Structures  map[string]*Structure
    OutputDir   string
}

// Structure from reflection-data JSON
type Structure struct {
    QualifiedName    string              `json:"qualifiedName"`
    StorageName      string              `json:"storageName"`
    SuperTypeName    string              `json:"superTypeName,omitempty"`
    Abstract         bool                `json:"abstract"`
    Type             string              `json:"type"` // ELEMENT, MODEL_UNIT, STRUCTURAL_UNIT
    Properties       map[string]*Property `json:"properties"`
    DefaultSettings  map[string]any      `json:"defaultSettings"`
}

// Property definition
type Property struct {
    Name        string    `json:"name"`
    StorageName string    `json:"storageName"`
    List        bool      `json:"list"`
    Public      bool      `json:"public"`
    TypeInfo    *TypeInfo `json:"typeInfo"`
}

// TypeInfo for property types
type TypeInfo struct {
    Type          string `json:"type"` // PRIMITIVE, ELEMENT, ENUMERATION, UNIT
    PrimitiveType string `json:"primitiveType,omitempty"`
    ElementType   string `json:"elementType,omitempty"`
    Kind          string `json:"kind,omitempty"` // PART, BY_ID_REFERENCE, BY_NAME_REFERENCE
}
```

### Generated Output Structure

```
generated/
├── domainmodels/
│   ├── entity.go
│   ├── attribute.go
│   ├── association.go
│   └── types.go
├── microflows/
│   ├── microflow.go
│   ├── activities.go      # 150+ activity types
│   ├── actions.go         # Action implementations
│   └── types.go
├── pages/
│   ├── page.go
│   ├── widgets.go         # 500+ widget types
│   ├── layouts.go
│   └── types.go
├── workflows/
│   ├── workflow.go
│   ├── activities.go
│   └── types.go
└── ... (48 more domains)
```

### Type Mapping

| JSON Type | Go Type |
|-----------|---------|
| `PRIMITIVE/STRING` | `string` |
| `PRIMITIVE/INTEGER` | `int64` |
| `PRIMITIVE/DOUBLE` | `float64` |
| `PRIMITIVE/BOOLEAN` | `bool` |
| `PRIMITIVE/DATE_TIME` | `time.Time` |
| `PRIMITIVE/GUID` | `model.ID` |
| `PRIMITIVE/POINT` | `model.Point` |
| `PRIMITIVE/SIZE` | `model.Size` |
| `PRIMITIVE/COLOR` | `model.Color` |
| `PRIMITIVE/BLOB` | `[]byte` |
| `ELEMENT` (single) | `*TypeName` |
| `ELEMENT` (list) | `[]*TypeName` |
| `BY_ID_REFERENCE` | `model.ID` |
| `BY_NAME_REFERENCE` | `model.QualifiedName` |
| `ENUMERATION` | Custom enum type |

### Example Generated Code

Input (`11.6.0-structures.json`):
```json
{
  "DomainModels$Entity": {
    "qualifiedName": "DomainModels$Entity",
    "superTypeName": "DomainModels$MaybeRemotableElement",
    "properties": {
      "name": {
        "storageName": "Name",
        "typeInfo": { "type": "PRIMITIVE", "primitiveType": "STRING" }
      },
      "attributes": {
        "storageName": "Attributes",
        "list": true,
        "typeInfo": { "type": "ELEMENT", "elementType": "DomainModels$Attribute", "kind": "PART" }
      },
      "generalization": {
        "storageName": "Generalization",
        "typeInfo": { "type": "ELEMENT", "elementType": "DomainModels$GeneralizationBase", "kind": "PART" }
      }
    }
  }
}
```

Output (`generated/domainmodels/entity.go`):
```go
// Code generated by modelsdk-generator. DO NOT EDIT.
// Source: 11.6.0-structures.json

package domainmodels

import "github.com/mendixlabs/mxcli/model"

// Entity represents a DomainModels$Entity element.
type Entity struct {
    model.BaseElement

    // Name is the entity name (storage: Name)
    Name string `json:"name" bson:"Name"`

    // Attributes contains the entity's attributes (storage: Attributes)
    Attributes []*Attribute `json:"attributes,omitempty" bson:"Attributes"`

    // Generalization defines inheritance (storage: Generalization)
    Generalization GeneralizationBase `json:"generalization,omitempty" bson:"Generalization"`
}

// GetName returns the entity's name.
func (e *Entity) GetName() string {
    return e.Name
}

// Ensure Entity implements the required interfaces.
var (
    _ model.Element      = (*Entity)(nil)
    _ model.NamedElement = (*Entity)(nil)
)
```

---

## Maintenance Process

### Adding Support for New Mendix Versions

When a new Mendix version is released:

```
┌─────────────────────────────────────────────────────────────┐
│                  Version Update Process                      │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  1. Obtain new reflection data                               │
│     └─▶ Get {version}-structures.json from mendixmodellib   │
│                                                              │
│  2. Run generator                                            │
│     └─▶ go run cmd/generate/main.go -version=11.7.0         │
│                                                              │
│  3. Review changes                                           │
│     └─▶ git diff generated/                                 │
│                                                              │
│  4. Update parser mappings                                   │
│     └─▶ Add new type handlers if needed                     │
│                                                              │
│  5. Run tests                                                │
│     └─▶ go test ./...                                       │
│                                                              │
│  6. Update version constants                                 │
│     └─▶ Add to supported versions list                      │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

### Step-by-Step Update Guide

#### 1. Obtain Reflection Data

```bash
# Option A: From npm package
npm pack mendixmodellib@latest
tar -xzf mendixmodellib-*.tgz
cp package/reflection-data/*.json libs/mendixmodellib/reflection-data/

# Option B: From existing installation
cp ~/Projects/mcpmxsdk/node_modules/mendixmodellib/reflection-data/*.json \
   libs/mendixmodellib/reflection-data/
```

#### 2. Run Generator

```bash
# Generate for specific version
go run cmd/generate/main.go -version=11.7.0 -output=generated/

# Generate for latest version
go run cmd/generate/main.go -latest -output=generated/

# Generate for all versions (creates version-specific packages)
go run cmd/generate/main.go -all -output=generated/
```

#### 3. Review Generated Changes

```bash
# See what changed
git diff generated/

# Check for breaking changes
go build ./...
go test ./...
```

#### 4. Update Version Support

```go
// version/versions.go
var SupportedVersions = []string{
    "10.0.0",
    "10.21.0",
    "11.0.0",
    "11.6.0",
    "11.7.0", // NEW
}

var LatestVersion = "11.7.0"
```

### Automation with CI/CD

```yaml
# .github/workflows/update-metamodel.yml
name: Update Metamodel

on:
  schedule:
    - cron: '0 0 * * 0'  # Weekly
  workflow_dispatch:
    inputs:
      version:
        description: 'Mendix version to add'
        required: true

jobs:
  update:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.24'

      - name: Setup Node
        uses: actions/setup-node@v4

      - name: Get latest mendixmodellib
        run: |
          npm pack mendixmodellib@latest
          tar -xzf mendixmodellib-*.tgz
          cp package/reflection-data/*.json libs/mendixmodellib/reflection-data/

      - name: Generate types
        run: go run cmd/generate/main.go -latest -output=generated/

      - name: Run tests
        run: go test ./...

      - name: Create PR
        uses: peter-evans/create-pull-request@v6
        with:
          title: "Update metamodel to latest version"
          branch: update-metamodel
```

---

## Implementation Roadmap

### Phase 1: Generator Foundation (Week 1-2)

- [ ] Create `cmd/generate/` package
- [ ] Implement JSON parser for reflection data
- [ ] Implement basic Go code emitter
- [ ] Generate `domainmodels` package as proof of concept
- [ ] Validate generated code compiles

### Phase 2: Core Domains (Week 3-4)

- [ ] Generate `microflows` with all 150+ activity types
- [ ] Generate `pages` with all 500+ widget types
- [ ] Generate `workflows` domain
- [ ] Update parser to use generated types

### Phase 3: Complete Domains (Week 5-6)

- [ ] Generate remaining 48 domains
- [ ] Implement enum types
- [ ] Add interface generation for polymorphic types
- [ ] Add BSON tags for serialization

### Phase 4: Infrastructure (Week 7-8)

- [ ] Implement metamodel reflection system
- [ ] Add version compatibility checking
- [ ] Implement delta change tracking
- [ ] Add qualified name resolution

### Phase 5: Testing & Documentation (Week 9-10)

- [ ] Add comprehensive unit tests
- [ ] Test against real Mendix projects
- [ ] Document API usage
- [ ] Create migration guide from TypeScript SDK

---

## Appendix: Metamodel Statistics

### Type Counts by Domain (v11.6.0)

| Domain | Types | Properties |
|--------|-------|------------|
| pages | 312 | 2,847 |
| microflows | 156 | 1,203 |
| expressions | 89 | 412 |
| workflows | 78 | 534 |
| domainmodels | 52 | 389 |
| rest | 45 | 312 |
| customwidgets | 38 | 267 |
| reports | 36 | 289 |
| webservices | 34 | 245 |
| datatypes | 28 | 156 |
| ... | ... | ... |
| **Total** | **~820** | **~7,500** |

### Supported Mendix Versions

```
6.x:  6.0.0 → 6.10.4  (11 versions)
7.x:  7.0.0 → 7.23.0  (24 versions)
8.x:  8.0.0 → 8.18.0  (19 versions)
9.x:  9.0.0 → 9.24.0  (26 versions)
10.x: 10.0.0 → 10.21.0 (22 versions)
11.x: 11.0.0 → 11.6.0  (7 versions)
────────────────────────────────
Total: 88 versions supported
```

---

## References

- [Mendix Model SDK Documentation](https://docs.mendix.com/apidocs-mxsdk/mxsdk/)
- [mendixmodelsdk npm package](https://www.npmjs.com/package/mendixmodelsdk)
- [mendixmodellib npm package](https://www.npmjs.com/package/mendixmodellib)
- [Mendix MPR File Format](https://docs.mendix.com/refguide/mpr-format/)
