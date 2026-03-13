# MDL to BSON Mapping

This document describes how MDL constructs map to BSON structures in Mendix MPR files.

## Table of Contents

1. [MPR File Format Overview](#mpr-file-format-overview)
2. [BSON Structure Conventions](#bson-structure-conventions)
3. [Entity Mapping](#entity-mapping)
4. [Attribute Mapping](#attribute-mapping)
5. [Validation Rule Mapping](#validation-rule-mapping)
6. [Index Mapping](#index-mapping)
7. [Association Mapping](#association-mapping)
8. [Enumeration Mapping](#enumeration-mapping)
9. [Text and Localization](#text-and-localization)
10. [ID Generation](#id-generation)
11. [Complete Example](#complete-example)
12. [Implementing New Document Types](#implementing-new-document-types)
13. [Page Widget Mapping](#page-widget-mapping)

---

## MPR File Format Overview

Mendix projects are stored in `.mpr` files which contain:

### MPR v1 (Mendix < 10.18)
Single SQLite database file with:
- `Unit` table: Document metadata
- `UnitContents` table: BSON document contents

### MPR v2 (Mendix >= 10.18)
SQLite metadata file + separate content files:
- `.mpr` file: SQLite with `Unit` table (metadata only)
- `mprcontents/` folder: Individual `.mxunit` files containing BSON

### Unit Types

| UnitType | Document Type |
|----------|---------------|
| `DomainModels$DomainModel` | Domain model (entities, associations) |
| `DomainModels$ViewEntitySourceDocument` | OQL query for VIEW entities |
| `Microflows$Microflow` | Microflow definition |
| `Microflows$Nanoflow` | Nanoflow definition |
| `Pages$Page` | Page definition |
| `Pages$Layout` | Layout definition |
| `Pages$Snippet` | Snippet definition |
| `Pages$BuildingBlock` | Building block definition |
| `Enumerations$Enumeration` | Enumeration definition |
| `JavaActions$JavaAction` | Java action definition |
| `Security$ProjectSecurity` | Project security settings |
| `Security$ModuleSecurity` | Module security settings |
| `Navigation$NavigationDocument` | Navigation profile |
| `Settings$ProjectSettings` | Project settings |
| `BusinessEvents$BusinessEventService` | Business event service |
| `CustomWidgets$CustomWidget` | Custom widget definition |

---

## BSON Structure Conventions

### Standard Fields

Every BSON document contains:

| Field | Type | Description |
|-------|------|-------------|
| `$ID` | Binary (UUID) | Unique identifier |
| `$Type` | String | Fully qualified type name |

### ID Format

IDs are stored as BSON Binary subtype 0 (generic) containing UUID bytes:
```json
{
  "$ID": {
    "Subtype": 0,
    "Data": "base64-encoded-uuid"
  }
}
```

UUID string format: `xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx`

### Array Format

**CRITICAL**: Arrays in Mendix BSON have the **count of elements** as the first element:
```json
{
  "Items": [
    2,           // Count of items (2 items follow)
    { ... },     // First element
    { ... }      // Second element
  ]
}
```

**Important notes:**
- The count is an `int32` value representing the number of actual items that follow
- When writing arrays, you MUST include this count prefix
- When parsing arrays, skip the first element (the count) when iterating
- Missing or incorrect counts will cause Studio Pro to misinterpret the data

Example in Go:
```go
// Writing an array with count prefix
items := bson.A{int32(len(changes))} // Start with count
for _, change := range changes {
    items = append(items, serializeItem(change))
}
```

> **JSON Templates vs BSON**: This array format with count prefixes applies to **BSON serialization in Go code**. When editing JSON template files (like `sdk/widgets/templates/*.json`), use truly empty arrays `[]` - the version markers are added automatically during BSON serialization. Writing `[2]` in JSON creates an array containing the integer 2, not an empty array.

### Reference Types

The Mendix metamodel defines two types of references:

| Reference Type | Storage Format | Example Use |
|---------------|----------------|-------------|
| `BY_ID_REFERENCE` | Binary UUID | Index `AttributePointer` |
| `BY_NAME_REFERENCE` | Qualified name string | ValidationRule `Attribute` |

**BY_ID_REFERENCE** - Stored as BSON Binary containing UUID bytes:
```json
{
  "AttributePointer": {
    "Subtype": 0,
    "Data": "base64-uuid"
  }
}
```

**BY_NAME_REFERENCE** - Stored as qualified name string:
```json
{
  "Attribute": "MyModule.MyEntity.MyAttribute"
}
```

> **Critical**: Using the wrong reference format will cause Studio Pro to fail loading the model. The metamodel reflection data specifies which format each property uses via the `kind` field in `typeInfo`.

### Type Names: qualifiedName vs storageName

The metamodel defines two type identifiers for each element type:

| Field | Usage | Example |
|-------|-------|---------|
| `qualifiedName` | TypeScript SDK API, internal naming | `DomainModels$Index` |
| `storageName` | BSON `$Type` field value | `DomainModels$EntityIndex` |

**Critical**: The `$Type` field in BSON must use the `storageName`, not the `qualifiedName`. These are often identical, but not always:

```json
// From metamodel reflection data
"DomainModels$Index" : {
  "qualifiedName" : "DomainModels$Index",
  "storageName" : "DomainModels$EntityIndex",  // ← Use this for $Type!
  ...
}
```

Using the wrong type name causes Studio Pro to fail with:
```
TypeCacheUnknownTypeException: The type cache does not contain a type with qualified name DomainModels$Index
```

**Known differences** (Mendix 11.6):

| qualifiedName | storageName (use this) |
|---------------|------------------------|
| `DomainModels$Index` | `DomainModels$EntityIndex` |
| `DomainModels$Entity` | `DomainModels$EntityImpl` |

When adding support for new document types, always check the metamodel reflection data in `reference/mendixmodellib/reflection-data/<version>-structures.json` to find the correct `storageName`.

### Metamodel Reference Definition

From the metamodel reflection data (`*-structures.json`):

```json
{
  "DomainModels$ValidationRule": {
    "properties": {
      "attribute": {
        "storageName": "Attribute",
        "typeInfo": {
          "type": "ELEMENT",
          "elementType": "DomainModels$Attribute",
          "kind": "BY_NAME_REFERENCE"
        }
      }
    }
  },
  "DomainModels$IndexedAttribute": {
    "properties": {
      "attribute": {
        "storageName": "AttributePointer",
        "typeInfo": {
          "type": "ELEMENT",
          "elementType": "DomainModels$Attribute",
          "kind": "BY_ID_REFERENCE"
        }
      }
    }
  }
}
```

---

## Entity Mapping

### MDL Entity
```sql
/** Documentation text */
@Position(100, 200)
CREATE PERSISTENT ENTITY Module.EntityName (
  AttrName: String(200) NOT NULL
)
INDEX (AttrName);
/
```

### BSON Structure
```json
{
  "$ID": {"Subtype": 0, "Data": "<uuid>"},
  "$Type": "DomainModels$EntityImpl",
  "Name": "EntityName",
  "Documentation": "Documentation text",
  "Location": "100;200",
  "MaybeGeneralization": {
    "$ID": {"Subtype": 0, "Data": "<uuid>"},
    "$Type": "DomainModels$NoGeneralization",
    "Persistable": true
  },
  "Attributes": [
    3,
    { /* attribute BSON */ }
  ],
  "ValidationRules": [
    3,
    { /* validation rule BSON */ }
  ],
  "Indexes": [
    3,
    { /* index BSON */ }
  ],
  "AccessRules": [3],
  "Events": [3]
}
```

### Entity Type Mapping

| MDL | BSON Persistable | BSON Source |
|-----|------------------|-------------|
| `PERSISTENT` | `true` | `null` |
| `NON-PERSISTENT` | `false` | `null` |
| `VIEW` | `false` | `OqlViewEntitySource` |
| `EXTERNAL` | `false` | `ODataRemoteEntitySource` |

### VIEW Entity Structure

VIEW entities require two separate documents:

1. **ViewEntitySourceDocument** - Contains the OQL query (MODEL_UNIT)
2. **Entity with OqlViewEntitySource** - References the source document

#### ViewEntitySourceDocument BSON

This is a separate document (unit) that stores the OQL query:

```json
{
  "$ID": {"Subtype": 0, "Data": "<uuid>"},
  "$Type": "DomainModels$ViewEntitySourceDocument",
  "Name": "ActiveProducts",
  "Documentation": "Products that are currently active",
  "Excluded": false,
  "ExportLevel": "Hidden",
  "Oql": "SELECT p.Name AS Name, p.Price AS Price FROM Module.Product p WHERE p.IsActive = true"
}
```

#### Entity with OqlViewEntitySource

The entity's `Source` field references the ViewEntitySourceDocument by qualified name:

```json
{
  "$Type": "DomainModels$EntityImpl",
  "Name": "ActiveProducts",
  "MaybeGeneralization": {
    "$Type": "DomainModels$NoGeneralization",
    "Persistable": false
  },
  "Source": {
    "$ID": {"Subtype": 0, "Data": "<uuid>"},
    "$Type": "DomainModels$OqlViewEntitySource",
    "SourceDocument": "Module.ActiveProducts"
  },
  "Attributes": [...]
}
```

#### OqlViewValue for View Attributes

View entity attributes use `OqlViewValue` instead of `StoredValue`. The `Reference` field contains the OQL column alias:

```json
{
  "$Type": "DomainModels$Attribute",
  "Name": "Name",
  "NewType": {
    "$Type": "DomainModels$StringAttributeType",
    "Length": 0
  },
  "Value": {
    "$ID": {"Subtype": 0, "Data": "<uuid>"},
    "$Type": "DomainModels$OqlViewValue",
    "Reference": "Name"
  }
}
```

The `Reference` value must match the OQL column alias (e.g., `AS Name` in the SELECT clause).

### Location Format

Position is stored as semicolon-separated string:
```
"Location": "100;200"
```

Parsed from MDL:
```sql
@Position(100, 200)
```

---

## Attribute Mapping

### MDL Attribute
```sql
/** Attribute documentation */
AttrName: String(200) NOT NULL ERROR 'Required' UNIQUE ERROR 'Must be unique' DEFAULT 'value'
```

### BSON Structure
```json
{
  "$ID": {"Subtype": 0, "Data": "<uuid>"},
  "$Type": "DomainModels$Attribute",
  "Name": "AttrName",
  "Documentation": "Attribute documentation",
  "NewType": {
    "$ID": {"Subtype": 0, "Data": "<uuid>"},
    "$Type": "DomainModels$StringAttributeType",
    "Length": 200
  },
  "Value": {
    "$ID": {"Subtype": 0, "Data": "<uuid>"},
    "$Type": "DomainModels$StoredValue",
    "DefaultValue": "value"
  }
}
```

### Attribute Type Mapping

| MDL Type | BSON $Type | Additional Fields |
|----------|------------|-------------------|
| `String` | `DomainModels$StringAttributeType` | `Length: 200` (default) |
| `String(n)` | `DomainModels$StringAttributeType` | `Length: n` |
| `Integer` | `DomainModels$IntegerAttributeType` | - |
| `Long` | `DomainModels$LongAttributeType` | - |
| `Decimal` | `DomainModels$DecimalAttributeType` | - |
| `Boolean` | `DomainModels$BooleanAttributeType` | - |
| `DateTime` | `DomainModels$DateTimeAttributeType` | `LocalizeDate: false` |
| `AutoNumber` | `DomainModels$AutoNumberAttributeType` | - |
| `Binary` | `DomainModels$BinaryAttributeType` | - |
| `Enumeration(M.E)` | `DomainModels$EnumerationAttributeType` | `Enumeration: "Module.EnumName"` |
| `HashedString` | `DomainModels$HashedStringAttributeType` | - |

### Enumeration Attribute Type

The `Enumeration` field in `EnumerationAttributeType` uses a **BY_NAME_REFERENCE** (qualified name string), not a binary UUID:

```json
{
  "$ID": {"Subtype": 0, "Data": "<uuid>"},
  "$Type": "DomainModels$EnumerationAttributeType",
  "Enumeration": "MyModule.Status"
}
```

The qualified name references the enumeration document by its `Module.EnumName` path.

### Default Value Mapping

| MDL Default | BSON Structure |
|-------------|----------------|
| `DEFAULT 'text'` | `{$Type: "DomainModels$StoredValue", DefaultValue: "text"}` |
| `DEFAULT 123` | `{$Type: "DomainModels$StoredValue", DefaultValue: "123"}` |
| `DEFAULT TRUE` | `{$Type: "DomainModels$StoredValue", DefaultValue: "true"}` |
| `DEFAULT FALSE` | `{$Type: "DomainModels$StoredValue", DefaultValue: "false"}` |
| `DEFAULT 'EnumValue'` | `{$Type: "DomainModels$StoredValue", DefaultValue: "EnumValue"}` |
| (no default) | `Value` field absent or null |

---

## Validation Rule Mapping

Validation rules are stored separately from attributes in the entity's `ValidationRules` array.

### MDL Validation
```sql
AttrName: String NOT NULL ERROR 'Field is required' UNIQUE ERROR 'Must be unique'
```

### BSON Structure (Required Rule)

> **Important**: Field order matters. Studio Pro expects: `$ID`, `$Type`, `Attribute`, `Message`, `RuleInfo`

```json
{
  "$ID": {"Subtype": 0, "Data": "<uuid>"},
  "$Type": "DomainModels$ValidationRule",
  "Attribute": "Module.Entity.AttrName",
  "Message": {
    "$ID": {"Subtype": 0, "Data": "<uuid>"},
    "$Type": "Texts$Text",
    "Items": [
      3,
      {
        "$ID": {"Subtype": 0, "Data": "<uuid>"},
        "$Type": "Texts$Translation",
        "LanguageCode": "en_US",
        "Text": "Field is required"
      }
    ]
  },
  "RuleInfo": {
    "$ID": {"Subtype": 0, "Data": "<uuid>"},
    "$Type": "DomainModels$RequiredRuleInfo"
  }
}
```

### BSON Structure (Unique Rule)
```json
{
  "$ID": {"Subtype": 0, "Data": "<uuid>"},
  "$Type": "DomainModels$ValidationRule",
  "Attribute": "Module.Entity.AttrName",
  "Message": { /* same structure */ },
  "RuleInfo": {
    "$ID": {"Subtype": 0, "Data": "<uuid>"},
    "$Type": "DomainModels$UniqueRuleInfo"
  }
}
```

### Validation Rule Type Mapping

| MDL Constraint | BSON RuleInfo.$Type |
|----------------|---------------------|
| `NOT NULL` | `DomainModels$RequiredRuleInfo` |
| `UNIQUE` | `DomainModels$UniqueRuleInfo` |
| (future) `RANGE` | `DomainModels$RangeRuleInfo` |
| (future) `REGEX` | `DomainModels$RegexRuleInfo` |

### Attribute Reference (BY_NAME_REFERENCE)

The `Attribute` field in ValidationRule uses **BY_NAME_REFERENCE** and MUST be a qualified name string:

```json
{
  "Attribute": "Module.Entity.Attribute"
}
```

> **Critical**: Do NOT use binary UUID for this field. The metamodel specifies `"kind": "BY_NAME_REFERENCE"` for ValidationRule.attribute, which requires a qualified name string. Using binary UUID will cause `System.ArgumentNullException` in Studio Pro when it tries to resolve the attribute.

This is different from Index's `AttributePointer` which uses **BY_ID_REFERENCE** (binary UUID).

**Implementation Note**: When reading validation rules from BSON, the qualified name string is stored in the Go struct's `AttributeID` field. When re-serializing (e.g., when adding another entity to the domain model), the code must detect whether `AttributeID` contains a UUID (for newly created entities) or a qualified name string (for entities read from disk). If it contains dots, it's already a qualified name and can be used directly.

---

## Index Mapping

### MDL Index
```sql
INDEX (AttrName1, AttrName2 DESC)
```

### BSON Structure

> **Note**: Index uses `AttributePointer` with **BY_ID_REFERENCE** (binary UUID), unlike ValidationRule which uses qualified name strings.

```json
{
  "$ID": {"Subtype": 0, "Data": "<uuid>"},
  "$Type": "DomainModels$EntityIndex",
  "Attributes": [
    2,
    {
      "$ID": {"Subtype": 0, "Data": "<uuid>"},
      "$Type": "DomainModels$IndexedAttribute",
      "AttributePointer": {"Subtype": 0, "Data": "<attr1-uuid>"},
      "Ascending": true,
      "Type": "Normal"
    },
    {
      "$ID": {"Subtype": 0, "Data": "<uuid>"},
      "$Type": "DomainModels$IndexedAttribute",
      "AttributePointer": {"Subtype": 0, "Data": "<attr2-uuid>"},
      "Ascending": false,
      "Type": "Normal"
    }
  ]
}
```

### Attribute Reference (BY_ID_REFERENCE)

The `AttributePointer` field in IndexedAttribute uses **BY_ID_REFERENCE** and MUST be a binary UUID:

```json
{
  "AttributePointer": {"Subtype": 0, "Data": "<attribute-uuid>"}
}
```

This is different from ValidationRule's `Attribute` which uses **BY_NAME_REFERENCE** (qualified name string).

### Sort Order Mapping

| MDL | BSON Ascending |
|-----|----------------|
| `AttrName` (default) | `true` |
| `AttrName ASC` | `true` |
| `AttrName DESC` | `false` |

---

## Association Mapping

### MDL Association
```sql
-- Many Orders can reference one Customer (1-to-many from Customer perspective)
CREATE ASSOCIATION Module.Order_Customer
  FROM Module.Order       -- Entity holding the FK reference
  TO Module.Customer      -- Entity being referenced
  TYPE Reference
  OWNER Default           -- Creates 1-to-many cardinality
  DELETE_BEHAVIOR DELETE_BUT_KEEP_REFERENCES;
```

### BSON Structure
```json
{
  "$ID": {"Subtype": 0, "Data": "<uuid>"},
  "$Type": "DomainModels$Association",
  "Name": "Order_Customer",
  "Documentation": "",
  "ExportLevel": "Hidden",
  "GUID": {"Subtype": 0, "Data": "<uuid>"},
  "ParentPointer": {"Subtype": 0, "Data": "<order-entity-uuid>"},
  "ChildPointer": {"Subtype": 0, "Data": "<customer-entity-uuid>"},
  "Type": "Reference",
  "Owner": "Default",
  "ParentConnection": "0;50",
  "ChildConnection": "100;50",
  "StorageFormat": "Table",
  "Source": null,
  "DeleteBehavior": {
    "$ID": {"Subtype": 0, "Data": "<uuid>"},
    "$Type": "DomainModels$DeleteBehavior",
    "ChildDeleteBehavior": "DeleteMeButKeepReferences",
    "ChildErrorMessage": null,
    "ParentDeleteBehavior": "DeleteMeButKeepReferences",
    "ParentErrorMessage": null
  }
}
```

### Field Mapping

| MDL | BSON Field | Description |
|-----|------------|-------------|
| `FROM Entity` | `ParentPointer` | Entity holding the foreign key reference (BY_ID_REFERENCE) |
| `TO Entity` | `ChildPointer` | Entity being referenced (BY_ID_REFERENCE) |
| `DELETE_BEHAVIOR` | `DeleteBehavior.ChildDeleteBehavior` | Behavior when child (TO) entity is deleted |

### Association Type Mapping

| MDL Type | BSON Type |
|----------|-----------|
| `Reference` | `"Reference"` |
| `ReferenceSet` | `"ReferenceSet"` |

### Owner Mapping and Cardinality

The `Owner` setting determines relationship cardinality in Studio Pro:

| MDL Owner | BSON Owner | Cardinality | Use Case |
|-----------|------------|-------------|----------|
| `Default` | `"Default"` | **1-to-many** | Many Orders → One Customer |
| `Both` | `"Both"` | **1-to-1** | One Order ↔ One Customer |

**Important**: For many-to-one relationships, use `OWNER Default`. Using `OWNER Both` creates a one-to-one relationship.

```sql
-- Many-to-one (many orders to one customer)
CREATE ASSOCIATION Module.Order_Customer
FROM Module.Order TO Module.Customer
TYPE Reference
OWNER Default;  -- Creates 1-to-many: Customer has many Orders

-- One-to-one (bidirectional)
CREATE ASSOCIATION Module.Order_Customer
FROM Module.Order TO Module.Customer
TYPE Reference
OWNER Both;     -- Creates 1-to-1: Customer has one Order
```

### Delete Behavior Mapping

| MDL Behavior | BSON DeleteBehavior.Type |
|--------------|--------------------------|
| `DELETE_BUT_KEEP_REFERENCES` | `"DeleteMeButKeepReferences"` |
| `DELETE_CASCADE` | `"DeleteMeAndReferences"` |
| (default) | `"DeleteMeIfNoReferences"` |

---

## Enumeration Mapping

### MDL Enumeration
```sql
/** Status values */
CREATE ENUMERATION Module.Status (
  Active 'Active',
  Inactive 'Inactive',
  Pending 'Pending Review'
);
```

### BSON Structure
```json
{
  "$ID": {"Subtype": 0, "Data": "<uuid>"},
  "$Type": "Enumerations$Enumeration",
  "Name": "Status",
  "Documentation": "Status values",
  "Values": [
    3,
    {
      "$ID": {"Subtype": 0, "Data": "<uuid>"},
      "$Type": "Enumerations$EnumerationValue",
      "Name": "Active",
      "Caption": {
        "$Type": "Texts$Text",
        "Items": [
          3,
          {
            "$Type": "Texts$Translation",
            "LanguageCode": "en_US",
            "Text": "Active"
          }
        ]
      }
    },
    {
      "$ID": {"Subtype": 0, "Data": "<uuid>"},
      "$Type": "Enumerations$EnumerationValue",
      "Name": "Inactive",
      "Caption": { /* ... */ }
    },
    {
      "$ID": {"Subtype": 0, "Data": "<uuid>"},
      "$Type": "Enumerations$EnumerationValue",
      "Name": "Pending",
      "Caption": {
        "$Type": "Texts$Text",
        "Items": [
          3,
          {
            "$Type": "Texts$Translation",
            "LanguageCode": "en_US",
            "Text": "Pending Review"
          }
        ]
      }
    }
  ]
}
```

---

## Text and Localization

### Text Structure

All user-visible text uses the `Texts$Text` structure with translations:

```json
{
  "$ID": {"Subtype": 0, "Data": "<uuid>"},
  "$Type": "Texts$Text",
  "Items": [
    3,
    {
      "$ID": {"Subtype": 0, "Data": "<uuid>"},
      "$Type": "Texts$Translation",
      "LanguageCode": "en_US",
      "Text": "English text"
    },
    {
      "$ID": {"Subtype": 0, "Data": "<uuid>"},
      "$Type": "Texts$Translation",
      "LanguageCode": "nl_NL",
      "Text": "Dutch text"
    }
  ]
}
```

### Language Codes

Common language codes:
- `en_US` - English (United States)
- `en_GB` - English (United Kingdom)
- `nl_NL` - Dutch (Netherlands)
- `de_DE` - German (Germany)
- `fr_FR` - French (France)

### MDL Text Handling

Currently MDL uses the first available translation or `en_US` if available:

```sql
-- Error message uses en_US translation
AttrName: String NOT NULL ERROR 'This field is required'
```

Multi-language support is planned for future versions.

---

## ID Generation

When creating new elements, UUIDs must be generated for all `$ID` fields.

### UUID Format

Version 4 (random) UUIDs in standard format:
```
xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx
```

Where:
- `x` is any hexadecimal digit
- `4` indicates version 4
- `y` is one of: 8, 9, a, b

### Example

```go
// Go implementation
import "github.com/google/uuid"

id := uuid.New().String()
// "f47ac10b-58cc-4372-a567-0e02b2c3d479"
```

---

## Complete Example

### MDL Input
```sql
/** Customer entity */
@Position(100, 200)
CREATE PERSISTENT ENTITY Sales.Customer (
  /** Unique identifier */
  CustomerId: AutoNumber NOT NULL UNIQUE DEFAULT 1,
  /** Customer name */
  Name: String(200) NOT NULL ERROR 'Name is required',
  Email: String(200) UNIQUE ERROR 'Email must be unique'
)
INDEX (Name);
/
```

### BSON Output
```json
{
  "$ID": {"Subtype": 0, "Data": "..."},
  "$Type": "DomainModels$EntityImpl",
  "Name": "Customer",
  "Documentation": "Customer entity",
  "Location": "100;200",
  "MaybeGeneralization": {
    "$Type": "DomainModels$NoGeneralization",
    "Persistable": true
  },
  "Attributes": [
    3,
    {
      "$Type": "DomainModels$Attribute",
      "Name": "CustomerId",
      "Documentation": "Unique identifier",
      "NewType": {"$Type": "DomainModels$AutoNumberAttributeType"},
      "Value": {"$Type": "DomainModels$StoredValue", "DefaultValue": "1"}
    },
    {
      "$Type": "DomainModels$Attribute",
      "Name": "Name",
      "Documentation": "Customer name",
      "NewType": {"$Type": "DomainModels$StringAttributeType", "Length": 200}
    },
    {
      "$Type": "DomainModels$Attribute",
      "Name": "Email",
      "NewType": {"$Type": "DomainModels$StringAttributeType", "Length": 200}
    }
  ],
  "ValidationRules": [
    3,
    {
      "$ID": {"Subtype": 0, "Data": "..."},
      "$Type": "DomainModels$ValidationRule",
      "Attribute": "Sales.Customer.CustomerId",
      "RuleInfo": {"$ID": {...}, "$Type": "DomainModels$RequiredRuleInfo"}
    },
    {
      "$ID": {"Subtype": 0, "Data": "..."},
      "$Type": "DomainModels$ValidationRule",
      "Attribute": "Sales.Customer.CustomerId",
      "RuleInfo": {"$ID": {...}, "$Type": "DomainModels$UniqueRuleInfo"}
    },
    {
      "$ID": {"Subtype": 0, "Data": "..."},
      "$Type": "DomainModels$ValidationRule",
      "Attribute": "Sales.Customer.Name",
      "Message": {"$ID": {...}, "$Type": "Texts$Text", "Items": [3, {...}]},
      "RuleInfo": {"$ID": {...}, "$Type": "DomainModels$RequiredRuleInfo"}
    },
    {
      "$ID": {"Subtype": 0, "Data": "..."},
      "$Type": "DomainModels$ValidationRule",
      "Attribute": "Sales.Customer.Email",
      "Message": {"$ID": {...}, "$Type": "Texts$Text", "Items": [3, {...}]},
      "RuleInfo": {"$ID": {...}, "$Type": "DomainModels$UniqueRuleInfo"}
    }
  ],
  "Indexes": [
    3,
    {
      "$Type": "DomainModels$EntityIndex",
      "Attributes": [
        2,
        {"AttributePointer": "<name-attr-id>", "Ascending": true}
      ]
    }
  ],
  "AccessRules": [3],
  "Events": [3]
}
```

---

## Implementing New Document Types

When adding support for new Mendix metamodel types (microflows, pages, workflows, etc.), follow this checklist:

### 1. Find the Metamodel Definition

Locate the type in `reference/mendixmodellib/reflection-data/<version>-structures.json`:

```bash
# Search for a type
grep -A 20 '"Microflows\$Microflow"' reference/mendixmodellib/reflection-data/11.6.0-structures.json
```

### 2. Check storageName vs qualifiedName

**Always use `storageName` for the `$Type` field**:

```json
"Microflows$Microflow" : {
  "qualifiedName" : "Microflows$Microflow",
  "storageName" : "Microflows$Microflow",  // ← Use this
  ...
}
```

### 3. Check Property Reference Types

For each property that references another element, check the `kind` field:

```json
"properties": {
  "objectCollection": {
    "storageName": "ObjectCollection",
    "typeInfo": {
      "type": "ELEMENT",
      "elementType": "Microflows$MicroflowObjectCollection",
      "kind": "PART"  // Embedded object, serialize inline
    }
  },
  "returnType": {
    "storageName": "MicroflowReturnType",
    "typeInfo": {
      "type": "ELEMENT",
      "elementType": "DataTypes$DataType",
      "kind": "PART"
    }
  }
}
```

Reference kind mappings:

| Kind | Storage Format | Description |
|------|----------------|-------------|
| `PART` | Embedded BSON object | Child object serialized inline |
| `BY_ID_REFERENCE` | Binary UUID | Reference by ID (BSON Binary) |
| `BY_NAME_REFERENCE` | Qualified name string | Reference by name (e.g., "Module.Entity") |
| `LOOKUP` | Usually string | Named lookup in parent context |

### 4. Check Array Prefixes

Arrays have a type prefix as the first element. Common values:

| Prefix | Meaning |
|--------|---------|
| `2` | List/array type |
| `3` | Another array type (most common) |

Check existing BSON files to determine the correct prefix for each array property.

### 5. Check Default Values

The `defaultSettings` in the metamodel shows what fields can be omitted:

```json
"defaultSettings": {
  "documentation": "",
  "allowedModuleRoles": [],
  "markAsUsed": false
}
```

Fields matching their default value can often be omitted from BSON.

### 6. Test with Studio Pro

After implementing serialization:

1. Create a test MPR with Studio Pro
2. Use the SDK to modify it
3. Reopen in Studio Pro to verify no errors
4. Compare BSON output with original using debug tools

### Common Pitfalls

1. **Wrong $Type**: Using `qualifiedName` instead of `storageName`
2. **Wrong reference format**: Using UUID for BY_NAME_REFERENCE or vice versa
3. **Missing array prefix**: Forgetting the integer prefix in arrays
4. **Field ordering**: Some types require specific field order (use `bson.D` not `bson.M`)
5. **Missing $ID**: Every element needs a unique UUID in its `$ID` field

---

## Page Widget Mapping

### LayoutGridColumn

Layout grid columns have weight properties for responsive design:

| Field | Description | Values |
|-------|-------------|--------|
| `Weight` | Desktop column width | 1-12 for explicit, -1 for auto-fill |
| `PhoneWeight` | Phone column width | Usually -1 (auto) |
| `TabletWeight` | Tablet column width | Usually -1 (auto) |

**BSON Structure:**
```json
{
  "$ID": {"Subtype": 0, "Data": "<uuid>"},
  "$Type": "Forms$LayoutGridColumn",
  "Appearance": {...},
  "PhoneWeight": -1,
  "TabletWeight": -1,
  "Weight": 6,
  "Widgets": [3, ...]
}
```

> **Critical**: Use `Weight` (not `DesktopWeight`) for desktop column width. Using the wrong field name causes columns to display "Manual - 1" in Studio Pro.

### ActionButton with CaptionTemplate

ActionButton widgets use `CaptionTemplate` (not `Caption`) for parameterized button text with placeholders like `{1}`.

**BSON Structure:**
```json
{
  "$ID": {"Subtype": 0, "Data": "<uuid>"},
  "$Type": "Forms$ActionButton",
  "Action": {...},
  "Appearance": {...},
  "ButtonStyle": "Primary",
  "CaptionTemplate": {
    "$ID": {"Subtype": 0, "Data": "<uuid>"},
    "$Type": "Forms$ClientTemplate",
    "Fallback": {
      "$ID": {"Subtype": 0, "Data": "<uuid>"},
      "$Type": "Texts$Text",
      "Items": [3]
    },
    "Parameters": [
      2,
      {
        "$ID": {"Subtype": 0, "Data": "<uuid>"},
        "$Type": "Forms$ClientTemplateParameter",
        "AttributeRef": null,
        "Expression": "'Hello'",
        "FormattingInfo": {...},
        "SourceVariable": null
      }
    ],
    "Template": {
      "$ID": {"Subtype": 0, "Data": "<uuid>"},
      "$Type": "Texts$Text",
      "Items": [
        3,
        {
          "$ID": {"Subtype": 0, "Data": "<uuid>"},
          "$Type": "Texts$Translation",
          "LanguageCode": "en_US",
          "Text": "Save {1}"
        }
      ]
    }
  },
  "Name": "btnSave1",
  "RenderMode": "Button"
}
```

> **Critical Field Names and Structures:**
> 1. **`CaptionTemplate`** - Must be `CaptionTemplate`, NOT `Caption`. Using `Caption` causes the button text to not display in Studio Pro.
> 2. **`Fallback`** - Must be a `Texts$Text` object, NOT a string field like `FallbackValue`. Using `FallbackValue: ""` causes the template to fail.
> 3. **Array version markers differ by context:**
>    - `Parameters`: Use `[2, items...]` for non-empty, `[3]` for empty
>    - `Template.Items`: Use `[3, items...]` for non-empty, `[3]` for empty
>
> These differences were discovered by comparing SDK-generated BSON with Studio Pro-generated BSON. When in doubt, create a reference structure in Studio Pro and compare.

### ClientTemplate (DynamicText)

DynamicText widgets use `ClientTemplate` for parameterized content with placeholders like `{1}`.

**BSON Structure:**
```json
{
  "$ID": {"Subtype": 0, "Data": "<uuid>"},
  "$Type": "Forms$DynamicText",
  "Content": {
    "$ID": {"Subtype": 0, "Data": "<uuid>"},
    "$Type": "Forms$ClientTemplate",
    "Fallback": {
      "$ID": {"Subtype": 0, "Data": "<uuid>"},
      "$Type": "Texts$Text",
      "Items": [3]
    },
    "Parameters": [
      2,
      {
        "$ID": {"Subtype": 0, "Data": "<uuid>"},
        "$Type": "Forms$ClientTemplateParameter",
        "Expression": "'Hello World'"
      }
    ],
    "Template": {
      "$Type": "Texts$Text",
      "Items": [3, {"LanguageCode": "en_US", "Text": "{1}"}]
    }
  }
}
```

### ClientTemplateParameter Expression Format

The `Expression` field must contain a valid Mendix expression:

| Value Type | Expression Format | Example |
|------------|-------------------|---------|
| String literal | Single-quoted | `'Hello World'` |
| Variable | Dollar prefix | `$Parameter/Name` |
| Number | Unquoted | `42` |
| Boolean | Unquoted | `true` or `false` |
| Attribute path | Dollar + path | `$currentObject/Name` |

> **Critical**: String literals MUST be wrapped in single quotes. Using bare strings like `Hello` causes CE0117 "Error(s) in expression" errors in Studio Pro.

---

## Microflow Action Mapping

### CreateChangeAction (CREATE Object)

The CreateObjectAction uses `Microflows$CreateChangeAction` as its storageName.

**BSON Structure:**
```json
{
  "$ID": {"Subtype": 0, "Data": "<uuid>"},
  "$Type": "Microflows$CreateChangeAction",
  "Commit": "No",
  "Entity": "Module.EntityName",
  "ErrorHandlingType": "Rollback",
  "Items": [
    4,  // Count of items!
    {
      "$ID": {"Subtype": 0, "Data": "<uuid>"},
      "$Type": "Microflows$ChangeActionItem",
      "Association": "",
      "Attribute": "Module.Entity.AttributeName",
      "Type": "Set",
      "Value": "$ParameterName"
    },
    // ... more items
  ],
  "RefreshInClient": false,
  "VariableName": "NewProduct"
}
```

**Required fields for CreateChangeAction:**
| Field | Description | Required |
|-------|-------------|----------|
| `Commit` | "No", "Yes", or "YesWithoutEvents" | Yes |
| `Entity` | Qualified entity name (BY_NAME_REFERENCE) | Yes |
| `ErrorHandlingType` | "Rollback" or "Abort" | Yes |
| `Items` | Array of ChangeActionItem with count prefix | Yes |
| `RefreshInClient` | Boolean | Yes |
| `VariableName` | Output variable name (without $) | Yes |

### ChangeActionItem

Each item in the `Items` array represents an attribute assignment.

**Required fields:**
| Field | Description | Required |
|-------|-------------|----------|
| `$ID` | Unique UUID | Yes |
| `$Type` | `"Microflows$ChangeActionItem"` | Yes |
| `Association` | Empty string for attribute changes | Yes |
| `Attribute` | Qualified attribute name (BY_NAME_REFERENCE) | For attribute changes |
| `Type` | "Set", "Add", or "Remove" | Yes |
| `Value` | Expression string | Yes |

> **Critical**: The `Association` field MUST be present even when empty. Missing this field causes Studio Pro to fail silently, showing fewer items than expected.

### ChangeAction (CHANGE Object)

Similar to CreateChangeAction but uses `Microflows$ChangeAction`:

```json
{
  "$ID": {"Subtype": 0, "Data": "<uuid>"},
  "$Type": "Microflows$ChangeAction",
  "ChangeVariableName": "Product",
  "Commit": "No",
  "Items": [
    2,  // Count of items
    { /* ChangeActionItem */ },
    { /* ChangeActionItem */ }
  ],
  "RefreshInClient": false
}
```

---

## Debugging BSON Issues

When Studio Pro doesn't display data correctly (e.g., missing attributes, incorrect values), follow this debugging approach:

### 1. Create a Reference Copy

1. Create the expected structure manually in Studio Pro
2. Save and close Studio Pro
3. Use the SDK to extract and examine the BSON

### 2. Compare BSON Structures

Use this pattern to compare your generated BSON with Mendix-generated BSON:

```go
// In sdk/mpr/reader_units.go there's GetRawMicroflowByName for debugging
raw1, _ := reader.GetRawMicroflowByName("Module.BrokenMicroflow")
raw2, _ := reader.GetRawMicroflowByName("Module.WorkingMicroflow")

// Unmarshal and compare
var map1, map2 map[string]interface{}
bson.Unmarshal(raw1, &map1)
bson.Unmarshal(raw2, &map2)
// Compare field by field
```

### 3. Common BSON Issues

| Symptom | Likely Cause | Solution |
|---------|--------------|----------|
| Items missing in Studio Pro | Missing array count prefix | Add `int32(count)` as first element |
| First item missing | Count is 0 or missing | Verify count matches actual items |
| Fields showing default values | Missing required fields | Check metamodel for required fields |
| "Unknown type" error | Wrong $Type value | Use `storageName`, not `qualifiedName` |
| Silent failures | Missing optional-but-expected fields | Compare with Mendix-generated BSON |
| "(Empty caption)" on buttons | Wrong field name or structure | Use `CaptionTemplate`, not `Caption` |
| Template not displaying | Wrong Fallback structure | Use `Fallback: {Texts$Text}`, not `FallbackValue: ""` |

### 4. Key Patterns Discovered

1. **Array version markers vary by type**:
   - `Parameters` arrays use `[2, items...]` for non-empty
   - `Texts$Text.Items` arrays use `[3, items...]` for non-empty
   - Empty arrays typically use `[3]` alone
2. **Field names may differ from SDK types**: e.g., SDK uses `Caption` but BSON needs `CaptionTemplate`
3. **Object vs string fields**: e.g., `Fallback` must be a `Texts$Text` object, not a string
4. **Empty string vs null**: Some fields require empty string `""`, not null/omitted
5. **Required "optional" fields**: Some fields marked optional in metamodel are required by Studio Pro
6. **Field order**: Some elements require specific field ordering (use `bson.D`, not `bson.M`)

### 5. Debugging Checklist

- [ ] Array has count prefix as first element
- [ ] Count matches actual number of items
- [ ] All required fields are present (check metamodel)
- [ ] $Type uses `storageName` from reflection data
- [ ] References use correct format (BY_ID vs BY_NAME)
- [ ] Empty strings used where null might cause issues
- [ ] Field order matches Mendix-generated examples
