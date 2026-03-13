# Go Library API Reference

The `modelsdk-go` library provides programmatic access to Mendix projects from Go code. This is the underlying library that powers `mxcli`.

## Installation

```bash
go get github.com/mendixlabs/mxcli
```

## Quick Start

### Reading a Project

```go
package main

import (
    "fmt"
    "github.com/mendixlabs/mxcli"
)

func main() {
    // Open a Mendix project
    reader, err := modelsdk.Open("/path/to/MyApp.mpr")
    if err != nil {
        panic(err)
    }
    defer reader.Close()

    // List all modules
    modules, _ := reader.ListModules()
    for _, m := range modules {
        fmt.Printf("Module: %s\n", m.Name)
    }

    // Get domain model for a module
    dm, _ := reader.GetDomainModel(modules[0].ID)
    for _, entity := range dm.Entities {
        fmt.Printf("  Entity: %s\n", entity.Name)
        for _, attr := range entity.Attributes {
            fmt.Printf("    - %s: %s\n", attr.Name, attr.Type.GetTypeName())
        }
    }

    // List microflows
    microflows, _ := reader.ListMicroflows()
    fmt.Printf("Total microflows: %d\n", len(microflows))

    // List pages
    pages, _ := reader.ListPages()
    fmt.Printf("Total pages: %d\n", len(pages))
}
```

### Modifying a Project

```go
package main

import (
    "github.com/mendixlabs/mxcli"
)

func main() {
    // Open for writing
    writer, err := modelsdk.OpenForWriting("/path/to/MyApp.mpr")
    if err != nil {
        panic(err)
    }
    defer writer.Close()

    reader := writer.Reader()
    modules, _ := reader.ListModules()
    dm, _ := reader.GetDomainModel(modules[0].ID)

    // Create a new entity
    customer := modelsdk.NewEntity("Customer")
    writer.CreateEntity(dm.ID, customer)

    // Add attributes
    writer.AddAttribute(dm.ID, customer.ID, modelsdk.NewStringAttribute("Name", 200))
    writer.AddAttribute(dm.ID, customer.ID, modelsdk.NewStringAttribute("Email", 254))
    writer.AddAttribute(dm.ID, customer.ID, modelsdk.NewBooleanAttribute("IsActive"))
    writer.AddAttribute(dm.ID, customer.ID, modelsdk.NewDateTimeAttribute("CreatedDate", true))

    // Create another entity
    order := modelsdk.NewEntity("Order")
    writer.CreateEntity(dm.ID, order)

    // Create an association
    assoc := modelsdk.NewAssociation("Customer_Order", customer.ID, order.ID)
    writer.CreateAssociation(dm.ID, assoc)
}
```

### High-Level Fluent API

The `api/` package provides a simplified, fluent API inspired by the Mendix Web Extensibility Model API:

```go
package main

import (
    "github.com/mendixlabs/mxcli/api"
    "github.com/mendixlabs/mxcli/sdk/mpr"
)

func main() {
    writer, err := mpr.OpenForWriting("/path/to/MyApp.mpr")
    if err != nil {
        panic(err)
    }
    defer writer.Close()

    // Create the high-level API
    modelAPI := api.New(writer)

    // Set the current module context
    module, _ := modelAPI.Modules.GetModule("MyModule")
    modelAPI.SetModule(module)

    // Create entity with fluent builder
    customer, _ := modelAPI.DomainModels.CreateEntity("Customer").
        Persistent().
        WithStringAttribute("Name", 100).
        WithStringAttribute("Email", 254).
        WithIntegerAttribute("Age").
        WithBooleanAttribute("IsActive").
        WithDateTimeAttribute("CreatedDate", true).
        Build()

    // Create another entity
    order, _ := modelAPI.DomainModels.CreateEntity("Order").
        Persistent().
        WithDecimalAttribute("TotalAmount").
        WithDateTimeAttribute("OrderDate", true).
        Build()

    // Create association between entities
    _, _ = modelAPI.DomainModels.CreateAssociation("Customer_Orders").
        From("Customer").
        To("Order").
        OneToMany().
        Build()

    // Create enumeration
    _, _ = modelAPI.Enumerations.CreateEnumeration("OrderStatus").
        WithValue("Pending", "Pending").
        WithValue("Processing", "Processing").
        WithValue("Completed", "Completed").
        WithValue("Cancelled", "Cancelled").
        Build()

    // Create microflow
    _, _ = modelAPI.Microflows.CreateMicroflow("ACT_ProcessOrder").
        WithParameter("Order", "MyModule.Order").
        WithStringParameter("Message").
        ReturnsBoolean().
        Build()
}
```

## API Reference

### Core Types

| Type | Description |
|------|-------------|
| `modelsdk.ID` | Unique identifier for model elements (UUID) |
| `modelsdk.Module` | Represents a Mendix module |
| `modelsdk.Project` | Represents a Mendix project |
| `modelsdk.DomainModel` | Contains entities and associations |
| `modelsdk.Entity` | An entity in the domain model |
| `modelsdk.Attribute` | An attribute of an entity |
| `modelsdk.Association` | A relationship between entities |
| `modelsdk.Microflow` | A microflow (server-side logic) |
| `modelsdk.Nanoflow` | A nanoflow (client-side logic) |
| `modelsdk.Page` | A page in the UI |
| `modelsdk.Layout` | A page layout template |

### Reader Methods

```go
// Open a project
reader, _ := modelsdk.Open("path/to/project.mpr")
defer reader.Close()

// Metadata
reader.Path()                    // Get file path
reader.Version()                 // Get MPR version (1 or 2)
reader.GetMendixVersion()        // Get Mendix Studio Pro version

// Modules
reader.ListModules()             // List all modules
reader.GetModule(id)             // Get module by ID
reader.GetModuleByName(name)     // Get module by name

// Domain Models
reader.ListDomainModels()        // List all domain models
reader.GetDomainModel(moduleID)  // Get domain model for module

// Microflows & Nanoflows
reader.ListMicroflows()          // List all microflows
reader.GetMicroflow(id)          // Get microflow by ID
reader.ListNanoflows()           // List all nanoflows
reader.GetNanoflow(id)           // Get nanoflow by ID

// Pages & Layouts
reader.ListPages()               // List all pages
reader.GetPage(id)               // Get page by ID
reader.ListLayouts()             // List all layouts
reader.GetLayout(id)             // Get layout by ID

// Other
reader.ListEnumerations()        // List all enumerations
reader.ListConstants()           // List all constants
reader.ListScheduledEvents()     // List all scheduled events
reader.ExportJSON()              // Export entire model as JSON
```

### Writer Methods

```go
// Open for writing
writer, _ := modelsdk.OpenForWriting("path/to/project.mpr")
defer writer.Close()

// Access the reader
reader := writer.Reader()

// Modules
writer.CreateModule(module)
writer.UpdateModule(module)
writer.DeleteModule(id)

// Entities
writer.CreateEntity(domainModelID, entity)
writer.UpdateEntity(domainModelID, entity)
writer.DeleteEntity(domainModelID, entityID)

// Attributes
writer.AddAttribute(domainModelID, entityID, attribute)

// Associations
writer.CreateAssociation(domainModelID, association)
writer.DeleteAssociation(domainModelID, associationID)

// Microflows & Nanoflows
writer.CreateMicroflow(microflow)
writer.UpdateMicroflow(microflow)
writer.DeleteMicroflow(id)
writer.CreateNanoflow(nanoflow)
writer.UpdateNanoflow(nanoflow)
writer.DeleteNanoflow(id)

// Pages & Layouts
writer.CreatePage(page)
writer.UpdatePage(page)
writer.DeletePage(id)
writer.CreateLayout(layout)
writer.UpdateLayout(layout)
writer.DeleteLayout(id)

// Other
writer.CreateEnumeration(enumeration)
writer.CreateConstant(constant)
```

### Helper Functions

```go
// Create attributes
modelsdk.NewStringAttribute(name, length)
modelsdk.NewIntegerAttribute(name)
modelsdk.NewDecimalAttribute(name)
modelsdk.NewBooleanAttribute(name)
modelsdk.NewDateTimeAttribute(name, localize)
modelsdk.NewEnumerationAttribute(name, enumID)

// Create entities
modelsdk.NewEntity(name)                 // Persistable entity
modelsdk.NewNonPersistableEntity(name)   // Non-persistable entity

// Create associations
modelsdk.NewAssociation(name, parentID, childID)      // Reference (1:N)
modelsdk.NewReferenceSetAssociation(name, p, c)       // Reference set (M:N)

// Create flows
modelsdk.NewMicroflow(name)
modelsdk.NewNanoflow(name)

// Create pages
modelsdk.NewPage(name)

// Generate IDs
modelsdk.GenerateID()
```

### Fluent API Namespaces

| Namespace | Description |
|-----------|-------------|
| `modelAPI.DomainModels` | Create/modify entities, attributes, associations |
| `modelAPI.Enumerations` | Create/modify enumerations and values |
| `modelAPI.Microflows` | Create microflows with parameters and return types |
| `modelAPI.Pages` | Create pages with widgets (DataView, TextBox, etc.) |
| `modelAPI.Modules` | List and retrieve modules |

## Package Structure

```
github.com/mendixlabs/mxcli/
├── modelsdk.go          # Main package with convenience functions
├── model/               # Core model types (ID, Module, Project, etc.)
├── api/                 # High-level fluent API (builders)
│   ├── api.go           # ModelAPI entry point
│   ├── domainmodels.go  # EntityBuilder, AssociationBuilder
│   ├── enumerations.go  # EnumerationBuilder
│   ├── microflows.go    # MicroflowBuilder
│   ├── pages.go         # PageBuilder, widget builders
│   └── modules.go       # ModulesAPI
├── sdk/
│   ├── domainmodel/     # Domain model types (Entity, Attribute, Association)
│   ├── microflows/      # Microflow and Nanoflow types
│   ├── pages/           # Page, Layout, and Widget types
│   └── mpr/             # MPR file reader and writer
└── examples/            # Example applications
```

## MPR File Format

Mendix projects are stored in `.mpr` files which are SQLite databases containing BSON-encoded model elements.

### MPR v1 (Mendix < 10.18)
- Single `.mpr` file containing all model data
- Documents stored as BSON blobs in SQLite

### MPR v2 (Mendix >= 10.18)
- `.mpr` file contains references and metadata
- `mprcontents/` folder contains individual document files
- Better for Git versioning and large projects

The library automatically detects and handles both formats.

## Model Structure

```
Project
├── Modules
│   ├── Domain Model
│   │   ├── Entities
│   │   │   ├── Attributes
│   │   │   ├── Indexes
│   │   │   ├── Access Rules
│   │   │   ├── Validation Rules
│   │   │   └── Event Handlers
│   │   ├── Associations
│   │   └── Annotations
│   ├── Microflows
│   │   ├── Parameters
│   │   └── Activities & Flows
│   ├── Nanoflows
│   ├── Pages
│   │   ├── Widgets
│   │   └── Data Sources
│   ├── Layouts
│   ├── Snippets
│   ├── Enumerations
│   ├── Constants
│   ├── Scheduled Events
│   └── Java Actions
└── Project Documents
```

## Examples

### Read Project Information

```bash
cd examples/read_project
go run main.go /path/to/MyApp.mpr
```

### Modify Project

```bash
cd examples/modify_project
go run main.go /path/to/MyApp.mpr
```

**Warning**: Always backup your `.mpr` file before modifying it!

## Comparison with Official SDK

| Feature | Mendix Model SDK (TypeScript) | modelsdk-go |
|---------|-------------------------------|-------------|
| Language | TypeScript/JavaScript | Go |
| Runtime | Node.js | Native binary |
| Cloud Required | Yes (Platform API) | No |
| Local Files | No | Yes |
| Real-time Collaboration | Yes | No |
| Read Operations | Yes | Yes |
| Write Operations | Yes | Yes |
| Type Safety | Yes (TypeScript) | Yes (Go) |
| CLI Tool | No | Yes (mxcli) |
| SQL-like DSL | No | Yes (MDL) |

## Resources

- [Mendix Model SDK Documentation](https://docs.mendix.com/apidocs-mxsdk/mxsdk/)
- [Mendix Metamodel Documentation](https://docs.mendix.com/apidocs-mxsdk/mxsdk/mendix-metamodel/)
- [MPR File Format Discussion](https://community.mendix.com/link/space/studio-pro/questions/86892)
