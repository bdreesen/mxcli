# Proposal: Comprehensive MDL Syntax and Grammar Improvements (v2)

**Date:** 2026-01-20
**Author:** GitHub Copilot

## 1. Introduction

This document consolidates and expands upon previous proposals, offering a unified set of syntax and grammar improvements for all major Mendix Definition Language (MDL) document types: **Domain Models**, **Pages**, and **Microflows**.

The goal is to evolve MDL into a more modern, consistent, and expressive language that is:
-   **More Readable:** Intuitive for citizen developers and professionals alike.
-   **More Consistent:** Uniform rules for common operations across all document types.
-   **More Token-Efficient:** Compact syntax for faster processing by LLMs and other tools.
-   **More Expressive:** Powerful constructs for complex logic and UI definitions.

The analysis is based on the example files: `01-domain-model-examples.mdl`, `02-microflow-examples.mdl`, and `03-page-examples.mdl`.

---

## 2. Microflow Syntax Improvements

### 2.1. Unified Variable Declaration and Assignment

**Problem:** The language uses three different syntaxes for variable assignment (`DECLARE`, `SET`, and direct assignment with `=`), which is inconsistent and verbose.

**Proposal:** Adopt Go-style assignment operators.
-   **`:=` (Declare and Assign):** For first-time declaration and initialization.
-   **`=` (Assign):** For re-assigning a value to an existing variable.

**Example:**
```mdl
// Before
DECLARE $Counter Integer = 0;
SET $Counter = $Counter + 1;
$NewProduct = CREATE MfTest.Product(...);

// After
$Counter := 0;
$Counter = $Counter + 1;
$NewProduct := CREATE MfTest.Product(...);
```
**Benefit:** This makes the `DECLARE` and `SET` keywords obsolete, creating a single, clear rule for variable handling.

### 2.2. C-style Braces for Code Blocks

**Problem:** The `BEGIN...END`, `THEN...END IF`, and `BEGIN...END LOOP` constructs are verbose and less common in modern languages.

**Proposal:** Use curly braces `{}` to define all control flow blocks.

**Example:**
```mdl
// Before
IF $Product/IsActive THEN
  SET $ActiveCount = $ActiveCount + 1;
END IF;

// After
if ($Product/IsActive) {
  $ActiveCount = $ActiveCount + 1;
}
```
**Benefit:** Improves readability and token efficiency by aligning with a widely adopted standard.

### 2.3. Fluent APIs for List Operations

**Problem:** Chaining list operations (`FILTER`, `SORT`, `AVERAGE`) is clumsy, and the syntax is inconsistent.

**Proposal:** Introduce a fluent, pipelined syntax (method chaining) for all list and aggregate operations.

**Example:**
```mdl
// Before
$ActiveProducts = FILTER($ProductList, $IteratorProduct/IsActive = true);
$SortedProducts = SORT($ActiveProducts, Price DESC);
$AverageActivePrice = AVERAGE($SortedProducts.Price);

// After
$AverageActivePrice := $ProductList
  .filter($p -> $p/IsActive)
  .sort(Price DESC)
  .average($p -> $p/Price);
```
**Benefit:** Highly readable, expressive, and powerful for complex data manipulation.

### 2.4. Standardized Function and Action Calls

**Problem:** The `CALL` keyword for microflows and Java actions is verbose and inconsistent with built-in functions like `COUNT`.

**Proposal:** Remove the `CALL` keyword and treat all microflows and Java actions as standard, callable functions.

**Example:**
```mdl
// Before
$Result = CALL MICROFLOW MfTest.M003_StringOperations(FirstName = 'Hello', LastName = 'World!');

// After
$Result := MfTest.M003_StringOperations(FirstName: 'Hello', LastName: 'World!');
```
**Benefit:** Creates a single, unified syntax for all function-like calls.

---

## 3. Domain Model Syntax Improvements

### 3.1. Simplified Entity Definition

**Problem:** The `CREATE PERSISTENT ENTITY` syntax is verbose and separates attributes from their entity in a procedural style.

**Proposal:** Introduce a concise `entity` keyword with a `{}` block for attributes, similar to class definitions in other languages.

**Example:**
```mdl
// Before
@Position(500, 100)
CREATE PERSISTENT ENTITY DmTest.Customer (
  FirstName: String(100) NOT NULL,
  Email: String(200) UNIQUE ERROR 'Email must be unique',
  IsActive: Boolean DEFAULT true
);

// After
@Position(500, 100)
entity DmTest.Customer {
  FirstName: String(100) NOT NULL;
  Email: String(200) UNIQUE ERROR 'Email must be unique';
  IsActive: Boolean DEFAULT true;
}
```
**Benefit:** More declarative, readable, and aligned with modern object-oriented syntax.

### 3.2. Modernized Association Syntax

**Problem:** The `CREATE ASSOCIATION ... FROM ... TO ...` syntax is lengthy and separates the relationship definition from the entities involved.

**Proposal:** Define associations directly within the entity definition.

**Example:**
```mdl
// Before
CREATE ASSOCIATION DmTest.Order_Customer
FROM DmTest.SalesOrder TO DmTest.Customer
TYPE Reference;

// After
entity DmTest.SalesOrder {
  // ... other attributes
  association Order_Customer -> DmTest.Customer; // Many-to-one
}

entity DmTest.Project1 {
  // ... other attributes
  association Project_Employees -> List<DmTest.Employee>; // Many-to-many
}
```
**Benefit:** Co-locates the relationship with the entity, making the domain model much easier to understand at a glance.

---

## 4. Page Syntax Improvements

### 4.1. Declarative, Hierarchical Page Structure

**Problem:** The current page syntax is flat and procedural (`LAYOUTGRID ... ROW ... COLUMN ... WIDGET`), making it hard to visualize the UI's nested structure.

**Proposal:** Adopt a declarative, hierarchical syntax similar to modern UI frameworks like React (JSX) or SwiftUI.

**Example:**
```mdl
// Before
CREATE PAGE PgTest.P012_Product_Manage
TITLE 'Product Management'
BEGIN
  DATAGRID ProductGrid
    SOURCE DATABASE PgTest.Product
  BEGIN
    HEADER
      ACTIONBUTTON btnNew1 'New' ACTION CREATE_OBJECT PgTest.Product SHOW_PAGE 'PgTest.P012_Product_Manage_Edit';
    COLUMN 'Actions'
    BEGIN
      ACTIONBUTTON btnEdit 'Edit' ACTION SHOW_PAGE 'PgTest.P012_Product_Manage_Edit' PASSING ($Product = $currentObject);
      ACTIONBUTTON btnDelete 'Delete' ACTION DELETE_ACTION;
    END;
  END;
END;

// After
page PgTest.P012_Product_Manage {
  title: 'Product Management';

  DataGrid(source: database(PgTest.Product)) {
    header {
      ActionButton(caption: 'New') {
        action create(PgTest.Product) show_page('PgTest.P012_Product_Manage_Edit');
      }
    }
    Column(caption: 'Actions') {
      ActionButton(caption: 'Edit') {
        action show_page('PgTest.P012_Product_Manage_Edit', passing: $currentObject);
      }
      ActionButton(caption: 'Delete', style: Danger) {
        action delete();
      }
    }
  }
}
```
**Benefit:** The code structure directly mirrors the UI component tree, making it vastly more intuitive.

### 4.2. Simplified Widget Definitions

**Problem:** Widget definitions are verbose and mix identity, properties, and actions in a flat structure.

**Proposal:** Use a concise, function-call-like syntax for widgets, with named parameters for properties and nested blocks for content or actions.

**Example:**
```mdl
// Before
ACTIONBUTTON btnSave 'Save {1}' WITH ({1} = 'abc')
  ACTION SAVE_CHANGES
  STYLE Primary;

// After
ActionButton(caption: 'Save {1}', style: Primary) with ({1} = 'abc') {
  action save_changes();
}
```
**Benefit:** Cleaner, more structured, and clearly separates configuration from actions.

## 5. Summary of Benefits

-   **Unified Experience:** The proposed changes create a consistent feel across all document types.
-   **Improved Readability:** Syntax is more declarative and familiar to developers from other ecosystems.
-   **Reduced Verbosity:** Removing redundant keywords (`CREATE`, `SET`, `CALL`, `BEGIN`/`END`) makes the code more compact.
-   **Enhanced Expressiveness:** Fluent APIs and hierarchical UI definitions allow complex ideas to be expressed simply.
