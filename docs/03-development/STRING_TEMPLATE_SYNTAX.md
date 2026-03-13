# String Template Syntax Specification

This document describes the unified syntax for string templates with parameters in MDL (Mendix Definition Language). String templates allow embedding dynamic values into text using placeholders like `{1}`, `{2}`, etc.

## Overview

String templates are used in both microflows and pages to create dynamic text:
- **Microflows**: LOG statements, ShowMessage actions
- **Pages**: DynamicText content, ActionButton captions

Although the underlying Mendix metamodel uses different types (`Microflows$StringTemplate` vs `Forms$ClientTemplate`), MDL provides a **unified syntax** using the `WITH` clause.

## Basic Syntax

```sql
'Template text with {1} and {2}' WITH ({1} = expr1, {2} = expr2)
```

### Components

| Component | Description | Example |
|-----------|-------------|---------|
| Template text | String with `{n}` placeholders | `'Order {1} has {2} items'` |
| `WITH` clause | Maps placeholders to values | `WITH ({1} = $OrderNumber, {2} = $ItemCount)` |
| Placeholder | `{n}` where n is 1-based index | `{1}`, `{2}`, `{3}` |
| Expression | Value for the placeholder | `$Variable`, `'literal'`, `$widget.Attribute` |

## Expression Types

### Variables and Parameters

Use `$` prefix for variables and parameters:

```sql
-- Microflow parameter
'Processing order {1}' WITH ({1} = $OrderNumber)

-- Page parameter
'Welcome {1}' WITH ({1} = $Customer/Name)
```

### String Literals

Use single quotes for literal values:

```sql
'Status: {1}' WITH ({1} = 'Active')
```

### Data Source Attribute References (Pages Only)

In pages, widgets are nested within data containers (DataView, ListView, Gallery). To reference attributes from a specific data source, use the **widget name** as a qualifier:

```sql
$WidgetName.AttributeName
```

This explicitly identifies which data source provides the attribute value.

#### Example: Nested Data Sources

```sql
CREATE PAGE Sales.OrderDetailPage ($Order: Sales.Order)
BEGIN
  -- DataView named 'dvOrder' bound to page parameter
  DATAVIEW dvOrder DATASOURCE $Order
  BEGIN
    -- $dvOrder.OrderNumber refers to the dvOrder's context (Sales.Order)
    DYNAMICTEXT txtOrderNum (
      CONTENT 'Order #{1}' WITH ({1} = $dvOrder.OrderNumber)
    );

    -- Nested ListView named 'lvItems'
    LISTVIEW lvItems DATASOURCE $dvOrder/Sales.Order_OrderItem/Sales.OrderItem
    BEGIN
      -- Can reference BOTH parent (dvOrder) and current (lvItems) data sources
      DYNAMICTEXT txtItem (
        CONTENT 'Order {1}: {2}x {3} @ ${4}'
        WITH (
          {1} = $dvOrder.OrderNumber,   -- from parent DataView (Sales.Order)
          {2} = $lvItems.Quantity,       -- from current ListView (Sales.OrderItem)
          {3} = $lvItems.ProductName,    -- from current ListView (Sales.OrderItem)
          {4} = $lvItems.UnitPrice       -- from current ListView (Sales.OrderItem)
        )
      );

      -- Further nested DataView for product details
      DATAVIEW dvProduct DATASOURCE $lvItems/Sales.OrderItem_Product/Sales.Product
      BEGIN
        -- Three levels deep - can reference all in-scope data sources
        DYNAMICTEXT txtProductInfo (
          CONTENT 'Order {1} | Qty: {2} | Product: {3} (SKU: {4})'
          WITH (
            {1} = $dvOrder.OrderNumber,   -- grandparent (Sales.Order)
            {2} = $lvItems.Quantity,       -- parent (Sales.OrderItem)
            {3} = $dvProduct.Name,         -- current (Sales.Product)
            {4} = $dvProduct.SKU           -- current (Sales.Product)
          )
        );
      END;
    END;
  END;
END;
```

### Expressions

Any valid Mendix expression can be used:

```sql
'Total: {1}' WITH ({1} = toString($Total * 1.21))
'Items: {1}' WITH ({1} = toString(length($ItemList)))
```

## Usage in Microflows

### LOG Statement

```sql
-- Simple template
LOG INFO NODE 'OrderService' 'Processing order: {1}' WITH ({1} = $OrderNumber);

-- Multiple parameters
LOG INFO NODE 'OrderService' 'Order {1} for {2} totaling {3}' WITH (
  {1} = $OrderNumber,
  {2} = $CustomerName,
  {3} = toString($TotalAmount)
);

-- Without template (simple concatenation still works)
LOG INFO NODE 'OrderService' 'Processing order: ' + $OrderNumber;
```

### ShowMessage Action (Future)

```sql
SHOW MESSAGE INFO 'Order {1} created successfully' WITH ({1} = $OrderNumber);
```

## Usage in Pages

### DynamicText Widget

```sql
-- Simple template with literal
DYNAMICTEXT txtWelcome (
  CONTENT 'Welcome {1}!' WITH ({1} = 'User'),
  RENDERMODE 'H3'
);

-- Template with data source attribute
DYNAMICTEXT txtOrderInfo (
  CONTENT 'Order #{1} - {2}' WITH (
    {1} = $dvOrder.OrderNumber,
    {2} = $dvOrder.Status
  ),
  RENDERMODE 'Paragraph'
);
```

### ActionButton Caption

```sql
-- Button with dynamic caption
ACTIONBUTTON btnConfirm 'Confirm Order #{1}'
  WITH ({1} = $dvOrder.OrderNumber)
  ACTION CALL_MICROFLOW 'Sales.ConfirmOrder';

-- Multiple placeholders
ACTIONBUTTON btnProcess 'Process {1} items for ${2}'
  WITH ({1} = $dvOrder.ItemCount, {2} = $dvOrder.TotalAmount)
  ACTION SAVE_CHANGES
  STYLE Primary;
```

## Mendix Internal Mapping

### Microflows: `Microflows$StringTemplate`

The MDL syntax maps to the Mendix internal structure:

```
MDL:
  'Order {1} for {2}' WITH ({1} = $OrderNumber, {2} = $CustomerName)

BSON:
  {
    "$Type": "Microflows$StringTemplate",
    "Text": "'Order {1} for {2}'",
    "Parameters": [
      {"$Type": "Microflows$TemplateParameter", "Expression": "$OrderNumber"},
      {"$Type": "Microflows$TemplateParameter", "Expression": "$CustomerName"}
    ]
  }
```

### Pages: `Forms$ClientTemplate`

When referencing a page parameter's attribute (e.g., `$Product.Name` where `$Product` is a page parameter):

```
MDL:
  Content: '$Product.Name'
  -- or explicit --
  Content: 'Product: {1}', ContentParams: [$Product.Name]

BSON:
  {
    "$Type": "Forms$ClientTemplate",
    "Template": {"$Type": "Texts$Text", "Items": [{"Text": "{1}"}]},
    "Parameters": [
      {
        "$Type": "Forms$ClientTemplateParameter",
        "AttributeRef": {"Attribute": "Sales.Product.Name"},
        "Expression": "",
        "SourceVariable": {
          "$Type": "Forms$PageVariable",
          "PageParameter": "Product",
          "UseAllPages": false,
          "Widget": ""
        }
      }
    ]
  }
```

**Important**: When `SourceVariable` is set, it indicates the parameter binding. The `AttributeRef` still contains the full entity path for type resolution, but the `SourceVariable.PageParameter` preserves which page parameter is being referenced. This allows distinguishing between multiple parameters of the same entity type.

When referencing a widget's data source attribute (e.g., `$dvOrder.OrderNumber` where `dvOrder` is a DataView):

```
MDL:
  'Order #{1}' WITH ({1} = $dvOrder.OrderNumber)

BSON:
  {
    "$Type": "Forms$ClientTemplate",
    "Template": {"$Type": "Texts$Text", "Items": [{"Text": "Order #{1}"}]},
    "Parameters": [
      {
        "$Type": "Forms$ClientTemplateParameter",
        "AttributeRef": {"Attribute": "Sales.Order.OrderNumber"},
        "Expression": "",
        "SourceVariable": {
          "$Type": "Forms$PageVariable",
          "PageParameter": "",
          "UseAllPages": false,
          "Widget": "dvOrder"
        }
      }
    ]
  }
```

When using a simple expression (not a data source attribute):

```
MDL:
  'Hello {1}' WITH ({1} = 'World')

BSON:
  {
    "$Type": "Forms$ClientTemplate",
    "Parameters": [
      {
        "$Type": "Forms$ClientTemplateParameter",
        "AttributeRef": null,
        "Expression": "'World'"
      }
    ]
  }
```

## Syntax Reference

### Complete Grammar

```
templateString
    : STRING_LITERAL (WITH LPAREN templateParamList RPAREN)?
    ;

templateParamList
    : templateParam (COMMA templateParam)*
    ;

templateParam
    : LBRACE NUMBER RBRACE EQUALS templateValue
    ;

templateValue
    : expression                    // Any Mendix expression
    | dataSourceAttributeRef        // $WidgetName.Attribute (pages only)
    ;

dataSourceAttributeRef
    : VARIABLE DOT IDENTIFIER       // e.g., $dvOrder.OrderNumber
    ;
```

### Summary Table

| Context | Syntax | Example |
|---------|--------|---------|
| Microflow variable | `$VarName` | `$OrderNumber` |
| Microflow attribute | `$Var/Attr` | `$Order/OrderNumber` |
| Page parameter | `$ParamName` | `$Customer` |
| Page parameter attribute | `$ParamName.Attr` | `$Product.Name` |
| Page data source attribute | `$WidgetName.Attr` | `$dvOrder.OrderNumber` |
| String literal | `'text'` | `'Hello'` |
| Expression | any expression | `toString($Total)` |

## Migration from PARAMETERS Syntax

The previous `PARAMETERS [...]` syntax is deprecated. Migrate as follows:

```sql
-- Old syntax (deprecated)
ACTIONBUTTON btn 'Save {1}' PARAMETERS ['Hello'];
DYNAMICTEXT txt (CONTENT 'Value: {1}' PARAMETERS ['test']);

-- New unified syntax
ACTIONBUTTON btn 'Save {1}' WITH ({1} = 'Hello');
DYNAMICTEXT txt (CONTENT 'Value: {1}' WITH ({1} = 'test'));
```

Benefits of `WITH` syntax:
1. **Explicit mapping** - Clear which placeholder gets which value
2. **Consistent** - Same syntax for microflows and pages
3. **Flexible ordering** - Can use `{1}` and `{3}` without `{2}`
4. **Self-documenting** - Placeholder purpose is clear from the mapping

## Best Practices

1. **Use meaningful data source names** - Name widgets descriptively (`dvOrder`, `lvItems`, `galleryProducts`)

2. **Be explicit about data sources** - Always qualify attributes with the widget name in nested contexts

3. **Keep templates readable** - For complex templates, format the `WITH` clause on multiple lines

4. **Prefer templates over concatenation** - `'Order {1}' WITH ({1} = $num)` is clearer than `'Order ' + $num`

5. **Use expressions for formatting** - Apply formatting in the expression: `WITH ({1} = formatDecimal($price, 2))`

## Related: Unified Parameter Syntax for CALL Statements

MDL uses a consistent parameter passing syntax across different call types. Parameter **names** do not use the `$` prefix (that's reserved for variable **values**).

### Microflow Calls

```sql
-- Parameter names without $ prefix
CALL MICROFLOW Module.ProcessOrder (OrderId = $Id, CustomerName = 'John');

-- With result variable
$Result = CALL MICROFLOW Module.Calculate (Value = 100, Multiplier = 2);
```

### Java Action Calls

```sql
-- Same syntax as microflows
CALL JAVA ACTION CustomActivities.ExecuteOQL (OqlStatement = 'SELECT...');

-- With result
$Count = CALL JAVA ACTION CustomActivities.CountRecords (EntityName = 'Order');
```

### Nanoflow Calls

```sql
CALL NANOFLOW Module.ValidateInput (Input = $FormData, Strict = true);
```

### Comparison with String Templates

| Context | Parameter Syntax | Example |
|---------|------------------|---------|
| String template | `{n} = expr` | `WITH ({1} = $Name)` |
| Microflow call | `Name = expr` | `(FirstName = 'John')` |
| Java action call | `Name = expr` | `(Statement = $Query)` |

The key difference:
- **Templates** use positional `{1}`, `{2}` placeholders for text interpolation
- **Calls** use named parameters for function invocation
