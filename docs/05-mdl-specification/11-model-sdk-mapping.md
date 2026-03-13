# MDL to Model SDK Mapping

This document describes how MDL constructs map to the modelsdk-go library types and API.

## Table of Contents

1. [Package Structure](#package-structure)
2. [Entity Mapping](#entity-mapping)
3. [Attribute Mapping](#attribute-mapping)
4. [Validation Rule Mapping](#validation-rule-mapping)
5. [Index Mapping](#index-mapping)
6. [Association Mapping](#association-mapping)
7. [Enumeration Mapping](#enumeration-mapping)
8. [API Usage Examples](#api-usage-examples)

---

## Package Structure

The modelsdk-go library is organized into packages:

| Package | Description |
|---------|-------------|
| `modelsdk` | Main API: `Open()`, `OpenForWriting()`, helpers |
| `model` | Core types: `ID`, `Module`, `Element`, `Point` |
| `api` | High-level fluent API: `ModelAPI`, builders for entities, microflows, pages |
| `sdk/domainmodel` | Domain model types: `Entity`, `Attribute`, `Association` |
| `sdk/microflows` | Microflow types (60+ activity types) |
| `sdk/pages` | Page and widget types (50+ widgets) |
| `sdk/widgets` | Embedded widget templates for pluggable widgets |
| `sdk/mpr` | MPR file reading/writing, BSON parsing |
| `sql` | External database connectivity (PostgreSQL, Oracle, SQL Server) |
| `mdl/executor` | MDL statement execution engine |
| `mdl/catalog` | SQLite-based catalog for cross-reference queries |
| `mdl/linter` | Linting framework with built-in and Starlark rules |

---

## Entity Mapping

### MDL Entity
```sql
/** Customer entity */
@Position(100, 200)
CREATE PERSISTENT ENTITY Sales.Customer (
  Name: String(200) NOT NULL
);
```

### Go SDK Types

```go
import (
    "github.com/mendixlabs/mxcli/sdk/domainmodel"
    "github.com/mendixlabs/mxcli/model"
)

entity := &domainmodel.Entity{
    BaseElement: model.BaseElement{
        ID:       model.ID("generated-uuid"),
        TypeName: "DomainModels$EntityImpl",
    },
    Name:          "Customer",
    Documentation: "Customer entity",
    Location:      model.Point{X: 100, Y: 200},
    Persistable:   true,
    Attributes:    []*domainmodel.Attribute{...},
}
```

### Entity Type Mapping

| MDL | Go Field | Value |
|-----|----------|-------|
| `PERSISTENT` | `Persistable` | `true` |
| `NON-PERSISTENT` | `Persistable` | `false` |
| `VIEW` | `Persistable`, `Source` | `false`, `"OqlView"` |
| `@Position(x,y)` | `Location` | `model.Point{X: x, Y: y}` |
| `/** doc */` | `Documentation` | `"doc"` |

### Entity Struct Definition

```go
// domainmodel/domainmodel.go
type Entity struct {
    model.BaseElement
    ContainerID      model.ID    // Parent domain model ID
    Name             string
    Documentation    string
    Location         model.Point
    Persistable      bool

    // Entity members
    Attributes       []*Attribute
    Indexes          []*Index
    ValidationRules  []*ValidationRule
    AccessRules      []*AccessRule
    EventHandlers    []*EventHandler

    // Generalization
    Generalization    Generalization
    GeneralizationID  model.ID
    GeneralizationRef string      // e.g., "System.User"

    // External/View entity
    Source            string      // "OqlViewEntitySource", etc.
    RemoteSource      string
    OqlQuery          string      // For view entities
}
```

---

## Attribute Mapping

### MDL Attribute
```sql
/** Customer name */
Name: String(200) NOT NULL ERROR 'Required' DEFAULT 'Unknown'
```

### Go SDK Types

```go
attr := &domainmodel.Attribute{
    BaseElement: model.BaseElement{
        ID:       model.ID("generated-uuid"),
        TypeName: "DomainModels$Attribute",
    },
    Name:          "Name",
    Documentation: "Customer name",
    Type: &domainmodel.StringAttributeType{
        Length: 200,
    },
    Value: &domainmodel.AttributeValue{
        Type:         "StoredValue",
        DefaultValue: "Unknown",
    },
}
```

### Attribute Type Mapping

| MDL Type | Go Type |
|----------|---------|
| `String` | `*StringAttributeType{Length: 200}` |
| `String(n)` | `*StringAttributeType{Length: n}` |
| `Integer` | `*IntegerAttributeType{}` |
| `Long` | `*LongAttributeType{}` |
| `Decimal` | `*DecimalAttributeType{}` |
| `Boolean` | `*BooleanAttributeType{}` |
| `DateTime` | `*DateTimeAttributeType{}` |
| `AutoNumber` | `*AutoNumberAttributeType{}` |
| `Binary` | `*BinaryAttributeType{}` |
| `Enumeration(M.E)` | `*EnumerationAttributeType{EnumerationID: id}` |

### Attribute Type Interface

```go
// domainmodel/domainmodel.go
type AttributeType interface {
    GetTypeName() string
}

type StringAttributeType struct {
    model.BaseElement
    Length int
}

func (t *StringAttributeType) GetTypeName() string {
    return "String"
}

type IntegerAttributeType struct {
    model.BaseElement
}

func (t *IntegerAttributeType) GetTypeName() string {
    return "Integer"
}

// ... similar for other types
```

### Attribute Value

```go
type AttributeValue struct {
    model.BaseElement
    Type         string    // "StoredValue" or "CalculatedValue"
    DefaultValue string    // String representation of default
    MicroflowID  model.ID  // For calculated values
}
```

---

## Validation Rule Mapping

### MDL Validation
```sql
Name: String NOT NULL ERROR 'Name is required' UNIQUE ERROR 'Name must be unique'
```

### Go SDK Types

```go
// Required validation
requiredRule := &domainmodel.ValidationRule{
    BaseElement: model.BaseElement{
        ID: model.ID("generated-uuid"),
    },
    AttributeID: attrID,  // or qualified name like "Module.Entity.Attr"
    Type:        "Required",
    ErrorMessage: &model.Text{
        Translations: map[string]string{
            "en_US": "Name is required",
        },
    },
}

// Unique validation
uniqueRule := &domainmodel.ValidationRule{
    BaseElement: model.BaseElement{
        ID: model.ID("generated-uuid"),
    },
    AttributeID:  attrID,
    Type:         "Unique",
    ErrorMessage: &model.Text{
        Translations: map[string]string{
            "en_US": "Name must be unique",
        },
    },
}
```

### Validation Rule Type Mapping

| MDL Constraint | Go Type Field |
|----------------|---------------|
| `NOT NULL` | `Type: "Required"` |
| `UNIQUE` | `Type: "Unique"` |
| `NOT NULL ERROR 'msg'` | `Type: "Required"`, `ErrorMessage: {...}` |
| `UNIQUE ERROR 'msg'` | `Type: "Unique"`, `ErrorMessage: {...}` |

### ValidationRule Struct

```go
type ValidationRule struct {
    model.BaseElement
    ContainerID  model.ID     // Parent entity ID
    AttributeID  model.ID     // Can be UUID or qualified name
    Type         string       // "Required", "Unique", "Range", "Regex"
    ErrorMessage *model.Text  // Localized error message
    Rule         ValidationRuleInfo  // Additional rule details
}
```

---

## Index Mapping

### MDL Index
```sql
INDEX (Name, CreatedAt DESC)
```

### Go SDK Types

```go
index := &domainmodel.Index{
    BaseElement: model.BaseElement{
        ID: model.ID("generated-uuid"),
    },
    Attributes: []*domainmodel.IndexAttribute{
        {
            AttributeID: nameAttrID,
            Ascending:   true,
        },
        {
            AttributeID: createdAtAttrID,
            Ascending:   false,
        },
    },
}
```

### Index Struct

```go
type Index struct {
    model.BaseElement
    ContainerID  model.ID           // Parent entity ID
    Name         string             // Optional index name
    Attributes   []*IndexAttribute  // Indexed columns
    AttributeIDs []model.ID         // Alternative: just IDs
}

type IndexAttribute struct {
    model.BaseElement
    AttributeID model.ID
    Ascending   bool
}
```

### Sort Order Mapping

| MDL | Go Ascending |
|-----|--------------|
| `AttrName` | `true` |
| `AttrName ASC` | `true` |
| `AttrName DESC` | `false` |

---

## Association Mapping

### MDL Association
```sql
CREATE ASSOCIATION Sales.Order_Customer
  FROM Sales.Customer
  TO Sales.Order
  TYPE Reference
  OWNER Default
  DELETE_BEHAVIOR DELETE_BUT_KEEP_REFERENCES;
```

### Go SDK Types

```go
assoc := &domainmodel.Association{
    BaseElement: model.BaseElement{
        ID:       model.ID("generated-uuid"),
        TypeName: "DomainModels$Association",
    },
    Name:     "Order_Customer",
    ParentID: customerEntityID,
    ChildID:  orderEntityID,
    Type:     domainmodel.AssociationTypeReference,
    Owner:    domainmodel.AssociationOwnerDefault,
    ParentDeleteBehavior: &domainmodel.DeleteBehavior{
        Type: domainmodel.DeleteBehaviorTypeDeleteMeButKeepReferences,
    },
}
```

### Association Type Mapping

| MDL | Go Constant |
|-----|-------------|
| `Reference` | `AssociationTypeReference` |
| `ReferenceSet` | `AssociationTypeReferenceSet` |

### Owner Mapping

| MDL | Go Constant |
|-----|-------------|
| `Default` | `AssociationOwnerDefault` |
| `Both` | `AssociationOwnerBoth` |
| `Parent` | (not yet defined) |
| `Child` | (not yet defined) |

### Delete Behavior Mapping

| MDL | Go Constant |
|-----|-------------|
| `DELETE_BUT_KEEP_REFERENCES` | `DeleteBehaviorTypeDeleteMeButKeepReferences` |
| `DELETE_CASCADE` | `DeleteBehaviorTypeDeleteMeAndReferences` |

### Association Struct

```go
type Association struct {
    model.BaseElement
    ContainerID          model.ID
    Name                 string
    Documentation        string
    ParentID             model.ID
    ChildID              model.ID
    Type                 AssociationType
    Owner                AssociationOwner
    ParentConnection     model.Point
    ChildConnection      model.Point
    ParentDeleteBehavior *DeleteBehavior
    ChildDeleteBehavior  *DeleteBehavior
}

type AssociationType string

const (
    AssociationTypeReference    AssociationType = "Reference"
    AssociationTypeReferenceSet AssociationType = "ReferenceSet"
)

type AssociationOwner string

const (
    AssociationOwnerDefault AssociationOwner = "Default"
    AssociationOwnerBoth    AssociationOwner = "Both"
)
```

---

## Enumeration Mapping

### MDL Enumeration
```sql
CREATE ENUMERATION Sales.OrderStatus (
  Draft 'Draft Order',
  Pending 'Pending Approval',
  Approved 'Approved'
);
```

### Go SDK Types

```go
enum := &model.Enumeration{
    BaseElement: model.BaseElement{
        ID:       model.ID("generated-uuid"),
        TypeName: "Enumerations$Enumeration",
    },
    Name: "OrderStatus",
    Values: []*model.EnumerationValue{
        {
            Name: "Draft",
            Caption: &model.Text{
                Translations: map[string]string{
                    "en_US": "Draft Order",
                },
            },
        },
        {
            Name: "Pending",
            Caption: &model.Text{
                Translations: map[string]string{
                    "en_US": "Pending Approval",
                },
            },
        },
        {
            Name: "Approved",
            Caption: &model.Text{
                Translations: map[string]string{
                    "en_US": "Approved",
                },
            },
        },
    },
}
```

---

## API Usage Examples

### Reading Entities

```go
import (
    modelsdk "github.com/mendixlabs/mxcli"
)

// Open project read-only
reader, err := modelsdk.Open("/path/to/project.mpr")
if err != nil {
    log.Fatal(err)
}
defer reader.Close()

// List modules
modules, err := reader.ListModules()
for _, m := range modules {
    fmt.Printf("Module: %s\n", m.Name)
}

// Get domain model
dm, err := reader.GetDomainModel(moduleID)
for _, entity := range dm.Entities {
    fmt.Printf("Entity: %s (persistable: %v)\n",
        entity.Name, entity.Persistable)

    for _, attr := range entity.Attributes {
        fmt.Printf("  - %s: %s\n",
            attr.Name, attr.Type.GetTypeName())
    }
}
```

### Creating Entities

```go
// Open project for writing
writer, err := modelsdk.OpenForWriting("/path/to/project.mpr")
if err != nil {
    log.Fatal(err)
}
defer writer.Close()

// Create entity using helper
entity := modelsdk.NewEntity("Customer")
entity.Documentation = "Customer master data"
entity.Location = model.Point{X: 100, Y: 200}

// Add attributes using helpers
entity.Attributes = append(entity.Attributes,
    modelsdk.NewAutoNumberAttribute("CustomerId"),
    modelsdk.NewStringAttribute("Name", 200),
    modelsdk.NewStringAttribute("Email", 200),
)

// Create in domain model
err = writer.CreateEntity(domainModelID, entity)
if err != nil {
    log.Fatal(err)
}
```

### Helper Functions

```go
// modelsdk.go - Public helper functions

// NewEntity creates a new persistent entity
func NewEntity(name string) *domainmodel.Entity {
    return &domainmodel.Entity{
        BaseElement: model.BaseElement{
            ID: model.ID(GenerateID()),
        },
        Name:        name,
        Persistable: true,
    }
}

// NewNonPersistableEntity creates a non-persistent entity
func NewNonPersistableEntity(name string) *domainmodel.Entity {
    entity := NewEntity(name)
    entity.Persistable = false
    return entity
}

// NewStringAttribute creates a string attribute
func NewStringAttribute(name string, length int) *domainmodel.Attribute {
    return &domainmodel.Attribute{
        BaseElement: model.BaseElement{
            ID: model.ID(GenerateID()),
        },
        Name: name,
        Type: &domainmodel.StringAttributeType{Length: length},
    }
}

// NewIntegerAttribute creates an integer attribute
func NewIntegerAttribute(name string) *domainmodel.Attribute {
    return &domainmodel.Attribute{
        BaseElement: model.BaseElement{
            ID: model.ID(GenerateID()),
        },
        Name: name,
        Type: &domainmodel.IntegerAttributeType{},
    }
}

// ... similar helpers for other types
```

### MDL Executor Integration

The MDL executor (`mdl/executor` package) translates MDL AST to SDK calls:

```go
// executor/cmd_entities.go
func (e *Executor) execCreateEntity(s *ast.CreateEntityStmt) error {
    // Find module
    module, err := e.findModule(s.Name.Module)
    if err != nil {
        return err
    }

    // Get domain model
    dm, err := e.reader.GetDomainModel(module.ID)
    if err != nil {
        return err
    }

    // Create entity
    entity := &domainmodel.Entity{
        Name:        s.Name.Name,
        Documentation: s.Documentation,
        Location:    model.Point{X: s.Position.X, Y: s.Position.Y},
        Persistable: s.Kind != ast.EntityNonPersistent,
    }

    // Convert attributes
    for _, a := range s.Attributes {
        attr := &domainmodel.Attribute{
            Name:          a.Name,
            Documentation: a.Documentation,
            Type:          convertDataType(a.Type),
        }
        if a.HasDefault {
            attr.Value = &domainmodel.AttributeValue{
                DefaultValue: fmt.Sprintf("%v", a.DefaultValue),
            }
        }
        entity.Attributes = append(entity.Attributes, attr)
    }

    // Write to project
    return e.writer.CreateEntity(dm.ID, entity)
}
```

---

## Type Conversion Functions

### MDL AST to SDK Types

```go
// executor/executor.go
func convertDataType(dt ast.DataType) domainmodel.AttributeType {
    switch dt.Kind {
    case ast.TypeString:
        return &domainmodel.StringAttributeType{Length: dt.Length}
    case ast.TypeInteger:
        return &domainmodel.IntegerAttributeType{}
    case ast.TypeLong:
        return &domainmodel.LongAttributeType{}
    case ast.TypeDecimal:
        return &domainmodel.DecimalAttributeType{}
    case ast.TypeBoolean:
        return &domainmodel.BooleanAttributeType{}
    case ast.TypeDateTime:
        return &domainmodel.DateTimeAttributeType{}
    case ast.TypeAutoNumber:
        return &domainmodel.AutoNumberAttributeType{}
    case ast.TypeBinary:
        return &domainmodel.BinaryAttributeType{}
    case ast.TypeEnumeration:
        return &domainmodel.EnumerationAttributeType{
            // EnumerationID resolved from dt.EnumRef
        }
    default:
        return &domainmodel.StringAttributeType{Length: 200}
    }
}
```

### SDK Types to MDL Output

```go
// executor/executor.go
func getAttributeTypeName(at domainmodel.AttributeType) string {
    switch t := at.(type) {
    case *domainmodel.StringAttributeType:
        if t.Length > 0 {
            return fmt.Sprintf("String(%d)", t.Length)
        }
        return "String"
    case *domainmodel.IntegerAttributeType:
        return "Integer"
    case *domainmodel.LongAttributeType:
        return "Long"
    case *domainmodel.DecimalAttributeType:
        return "Decimal"
    case *domainmodel.BooleanAttributeType:
        return "Boolean"
    case *domainmodel.DateTimeAttributeType:
        return "DateTime"
    case *domainmodel.AutoNumberAttributeType:
        return "AutoNumber"
    case *domainmodel.BinaryAttributeType:
        return "Binary"
    case *domainmodel.EnumerationAttributeType:
        if t.EnumerationID != "" {
            return fmt.Sprintf("Enumeration(%s)", t.EnumerationID)
        }
        return "Enumeration"
    default:
        return "Unknown"
    }
}
```

---

## High-Level Fluent API

The `api/` package provides a simplified builder API as an alternative to direct SDK type construction.

### Entity Builder

```go
import "github.com/mendixlabs/mxcli/api"

modelAPI := api.New(writer)
module, _ := modelAPI.Modules.GetModule("Sales")
modelAPI.SetModule(module)

entity, _ := modelAPI.DomainModels.CreateEntity("Customer").
    Persistent().
    WithStringAttribute("Name", 200).
    WithIntegerAttribute("Age").
    WithEnumerationAttribute("Status", "Sales.CustomerStatus").
    Build()
```

### Microflow Builder

```go
mf, _ := modelAPI.Microflows.CreateMicroflow("ACT_ProcessOrder").
    WithParameter("Order", "Sales.Order").
    WithStringParameter("Note").
    ReturnsBoolean().
    Build()
```

### Enumeration Builder

```go
enum, _ := modelAPI.Enumerations.CreateEnumeration("OrderStatus").
    WithValue("Draft", "Draft").
    WithValue("Active", "Active").
    WithValue("Closed", "Closed").
    Build()
```

### MDL to API Mapping

| MDL Statement | Fluent API Method |
|---------------|-------------------|
| `CREATE PERSISTENT ENTITY` | `DomainModels.CreateEntity().Persistent().Build()` |
| `CREATE NON-PERSISTENT ENTITY` | `DomainModels.CreateEntity().NonPersistent().Build()` |
| `CREATE ASSOCIATION` | `DomainModels.CreateAssociation().Build()` |
| `CREATE ENUMERATION` | `Enumerations.CreateEnumeration().Build()` |
| `CREATE MICROFLOW` | `Microflows.CreateMicroflow().Build()` |
| `CREATE PAGE` | `Pages.CreatePage().Build()` |
