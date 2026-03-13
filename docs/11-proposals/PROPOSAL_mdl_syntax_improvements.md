# Proposal: MDL Syntax and Grammar Improvements

**Date:** 2026-01-20
**Author:** GitHub Copilot

## 1. Introduction

This document proposes a set of improvements to the Mendix Definition Language (MDL) syntax and grammar. The goal is to enhance the language's **readability**, **consistency**, and **token efficiency**, making it more intuitive for citizen developers and more effective for processing by LLMs.

The analysis is based on the existing capabilities demonstrated in `mdl-examples/doctype-tests/02-microflow-examples.mdl`.

## 2. Key Areas for Improvement

### 2.1. Variable Declaration and Assignment

**Current Syntax:**
The language uses three different ways to handle variable assignments:
- `DECLARE $VarName Type = initial_value;` (Declaration)
- `SET $VarName = new_value;` (Re-assignment)
- `$VarName = CREATE ...;` or `$VarName = CALL ...;` (Assignment from expression)

**Problem:**
This inconsistency increases the learning curve and adds unnecessary verbosity (e.g., the `SET` keyword is redundant).

**Proposal: Introduce Go-style Assignment Operators**
Adopt a more concise and consistent approach:
- **`:=` (Declare and Assign):** Use for first-time declaration and initialization within a scope.
- **`=` (Assign):** Use for re-assigning a new value to an already declared variable.

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
This change would make the `DECLARE` and `SET` keywords obsolete, significantly cleaning up the syntax.

### 2.2. Control Flow Block Syntax

**Current Syntax:**
MDL uses `BEGIN...END` and `THEN...END IF` to delineate code blocks.
- `IF $condition THEN ... END IF;`
- `LOOP $item IN $list BEGIN ... END LOOP;`

**Problem:**
This syntax, while explicit, is verbose and less common in modern programming languages. It consumes more tokens and can be harder to read for developers accustomed to C-style syntax.

**Proposal: Adopt C-style Brace Syntax `{}`**
Use curly braces to define blocks for all control flow statements.

**Example:**
```mdl
// Before
IF $Product/IsActive THEN
  SET $ActiveCount = $ActiveCount + 1;
  LOG INFO NODE 'Test' 'Found active product';
END IF;

// After
IF ($Product/IsActive) {
  $ActiveCount = $ActiveCount + 1;
  LOG INFO NODE 'Test' 'Found active product';
}
```

### 2.3. Fluent APIs for List and Aggregate Operations

**Current Syntax:**
List operations are function-based, and aggregate functions have a unique syntax.
- `$ActiveProducts = FILTER($ProductList, $IteratorProduct/IsActive = true);`
- `$AveragePrice = AVERAGE($ProductList.Price);`

**Problem:**
Chaining multiple operations is clumsy and hard to read. The syntax for aggregation (`$List.Attribute`) is inconsistent with other function calls.

**Proposal: Introduce Fluent (Pipelined) Syntax**
Allow for method-chaining on list variables. This is a highly readable and expressive pattern.

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
This syntax is more intuitive, token-efficient, and powerful for complex data manipulation. It uses lambda-style expressions (`$p -> ...`) for clarity.

### 2.4. `CHANGE` Statement Readability

**Current Syntax:**
The `CHANGE` statement modifies multiple attributes in a flat list.
`CHANGE $Product (Name = $NewName, ModifiedDate = [%CurrentDateTime%]);`

**Problem:**
For objects with many attributes, this can become a long, hard-to-read line.

**Proposal: Introduce `WITH` block for `CHANGE`**
Allow a block syntax for grouping attribute changes, improving readability.

**Example:**
```mdl
// Before
CHANGE $Product (DailyAverage = $DailyAverage, LastCalculated = [%CurrentDateTime%]);

// After
CHANGE $Product WITH {
  DailyAverage = $DailyAverage,
  LastCalculated = [%CurrentDateTime%]
};
```

### 2.5. Unify Function and Action Call Syntax

**Current Syntax:**
- `CALL MICROFLOW MfTest.M001_HelloWorld()`
- `CALL JAVA ACTION CustomActivities.ExecuteOQLStatement(...)`
- `COUNT($ProductList)`

**Problem:**
The `CALL` keyword is verbose and inconsistent with built-in function calls like `COUNT`.

**Proposal: Standardize All Calls**
Remove the `CALL` keyword and treat microflows and Java actions as regular callable functions. The system can distinguish them by their path.

**Example:**
```mdl
// Before
$Result = CALL MICROFLOW MfTest.M003_StringOperations(FirstName = 'Hello', LastName = 'World!');
$OqlResult = CALL JAVA ACTION CustomActivities.ExecuteOQLStatement(...);

// After
$Result := MfTest.M003_StringOperations(FirstName: 'Hello', LastName: 'World!');
$OqlResult := CustomActivities.ExecuteOQLStatement(...);
```
Using named parameters with colons (`:`) could further improve clarity, distinguishing them from positional parameters.

## 3. Summary of Benefits

- **Improved Readability:** The proposed syntax is closer to modern programming languages, making it more familiar to a wider range of developers.
- **Increased Consistency:** Rules for variable assignment and function calls are unified, reducing cognitive load.
- **Enhanced Token Efficiency:** Removing redundant keywords (`DECLARE`, `SET`, `CALL`, `THEN`, `BEGIN`/`END`) and using braces makes the code more compact for LLM processing.
- **Greater Expressiveness:** Fluent APIs for list manipulation allow for more complex logic to be expressed clearly and concisely.

Adopting these changes would represent a significant evolution for MDL, making it a more powerful and user-friendly language for Mendix development.
