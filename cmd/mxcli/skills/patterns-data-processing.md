# Data Processing Patterns

Patterns for loops, aggregates, batch processing, and data transformation.

## Loop Patterns

### Basic Loop

```mdl
/**
 * Process all items in a list
 */
CREATE MICROFLOW Module.ProcessItems (
  $Items: List of Module.Item
)
RETURNS Boolean
BEGIN
  DECLARE $ProcessedCount Integer = 0;

  LOOP $Item IN $Items
  BEGIN
    -- Process each item
    CHANGE $Item (ProcessedDate = [%CurrentDateTime%]);
    COMMIT $Item;
    SET $ProcessedCount = $ProcessedCount + 1;
  END LOOP;

  LOG INFO NODE 'Processing' 'Processed ' + $ProcessedCount + ' items';
  RETURN true;
END;
/
```

### Loop with Filtering

```mdl
/**
 * Process only active items
 */
CREATE MICROFLOW Module.ProcessActiveItems (
  $Items: List of Module.Item
)
RETURNS Integer
BEGIN
  DECLARE $Count Integer = 0;

  LOOP $Item IN $Items
  BEGIN
    IF $Item/IsActive THEN
      -- Process active item
      CHANGE $Item (LastProcessed = [%CurrentDateTime%]);
      COMMIT $Item;
      SET $Count = $Count + 1;
    END IF;
  END LOOP;

  RETURN $Count;
END;
/
```

### Loop with Accumulator

```mdl
/**
 * Calculate total value of all orders
 */
CREATE MICROFLOW Module.CalculateOrderTotal (
  $Orders: List of Module.Order
)
RETURNS Decimal
BEGIN
  DECLARE $Total Decimal = 0;

  LOOP $Order IN $Orders
  BEGIN
    SET $Total = $Total + $Order/Amount;
  END LOOP;

  RETURN $Total;
END;
/
```

## Aggregate Patterns

Aggregates use **function-call syntax** — there is no `AGGREGATE` keyword.

| Function | Syntax | Returns |
|----------|--------|---------|
| COUNT | `$n = COUNT($list)` | Integer |
| SUM | `$n = SUM($list.Attr)` | Decimal |
| AVERAGE | `$n = AVERAGE($list.Attr)` | Decimal |
| MINIMUM | `$n = MINIMUM($list.Attr)` | Same as attribute |
| MAXIMUM | `$n = MAXIMUM($list.Attr)` | Same as attribute |

**Important:** RETRIEVE implicitly declares its variable — do NOT add a separate DECLARE
before RETRIEVE, or you'll get CE0111 "Duplicate variable name".

### Count Items

```mdl
/**
 * Count active customers
 */
CREATE MICROFLOW Module.CountActiveCustomers ()
RETURNS Integer
BEGIN
  RETRIEVE $Customers FROM Module.Customer
    WHERE IsActive = true;

  $Count = COUNT($Customers);
  RETURN $Count;
END;
/
```

### Sum Values

```mdl
/**
 * Sum order amounts for a customer
 */
CREATE MICROFLOW Module.GetCustomerTotalOrders (
  $Customer: Module.Customer
)
RETURNS Decimal
BEGIN
  RETRIEVE $Orders FROM Module.Order
    WHERE Module.Order_Customer = $Customer;

  $Total = SUM($Orders.Amount);
  RETURN $Total;
END;
/
```

### Average Calculation

```mdl
/**
 * Calculate average order value
 */
CREATE MICROFLOW Module.GetAverageOrderValue ()
RETURNS Decimal
BEGIN
  RETRIEVE $Orders FROM Module.Order;

  $Average = AVERAGE($Orders.Amount);
  RETURN $Average;
END;
/
```

### Min/Max

```mdl
$MinPrice = MINIMUM($Products.Price);
$MaxPrice = MAXIMUM($Products.Price);
```

## List Operations

### Add to List

```mdl
/**
 * Collect matching items into a list
 */
CREATE MICROFLOW Module.CollectHighValueOrders (
  $Orders: List of Module.Order,
  $Threshold: Decimal
)
RETURNS List of Module.Order
BEGIN
  DECLARE $HighValue List of Module.Order = empty;

  LOOP $Order IN $Orders
  BEGIN
    IF $Order/Amount > $Threshold THEN
      ADD $Order TO $HighValue;
    END IF;
  END LOOP;

  RETURN $HighValue;
END;
/
```

### Remove from List

```mdl
/**
 * Remove inactive items from list
 */
CREATE MICROFLOW Module.FilterActiveItems (
  $Items: List of Module.Item
)
RETURNS List of Module.Item
BEGIN
  DECLARE $ToRemove List of Module.Item = empty;

  -- Collect items to remove
  LOOP $Item IN $Items
  BEGIN
    IF NOT $Item/IsActive THEN
      ADD $Item TO $ToRemove;
    END IF;
  END LOOP;

  -- Remove collected items
  LOOP $Item IN $ToRemove
  BEGIN
    REMOVE $Item FROM $Items;
  END LOOP;

  RETURN $Items;
END;
/
```

## Batch Processing

### Process in Batches

```mdl
/**
 * Process large dataset in batches
 * Commits after each batch to avoid memory issues
 */
CREATE MICROFLOW Module.BatchProcess (
  $Items: List of Module.Item,
  $BatchSize: Integer
)
RETURNS Integer
BEGIN
  DECLARE $Processed Integer = 0;
  DECLARE $BatchCount Integer = 0;

  LOOP $Item IN $Items
  BEGIN
    -- Process item
    CHANGE $Item (Status = 'Processed');

    SET $BatchCount = $BatchCount + 1;
    SET $Processed = $Processed + 1;

    -- Commit batch
    IF $BatchCount >= $BatchSize THEN
      COMMIT $Item;
      SET $BatchCount = 0;
      LOG INFO NODE 'Batch' 'Processed ' + $Processed + ' items';
    END IF;
  END LOOP;

  -- Final commit for remaining items
  IF $BatchCount > 0 THEN
    LOG INFO NODE 'Batch' 'Final batch: ' + $Processed + ' total';
  END IF;

  RETURN $Processed;
END;
/
```

## Data Transformation

### Copy Entity

```mdl
/**
 * Create a copy of an order
 */
CREATE MICROFLOW Module.CopyOrder (
  $Source: Module.Order
)
RETURNS Module.Order
BEGIN
  DECLARE $Copy AS Module.Order;

  $Copy = CREATE Module.Order (
    OrderNumber = 'COPY-' + $Source/OrderNumber,
    Amount = $Source/Amount,
    Status = 'Draft',
    CreatedDate = [%CurrentDateTime%]
  );

  -- Copy association
  SET $Copy/Module.Order_Customer = $Source/Module.Order_Customer;

  COMMIT $Copy;
  RETURN $Copy;
END;
/
```

### Transform List

```mdl
/**
 * Create summary records from detail records
 */
CREATE MICROFLOW Module.CreateOrderSummaries (
  $Orders: List of Module.Order
)
RETURNS List of Module.OrderSummary
BEGIN
  DECLARE $Summaries List of Module.OrderSummary = empty;
  DECLARE $Summary AS Module.OrderSummary;

  LOOP $Order IN $Orders
  BEGIN
    $Summary = CREATE Module.OrderSummary (
      OrderNumber = $Order/OrderNumber,
      TotalAmount = $Order/Amount,
      CustomerName = $Order/Module.Order_Customer/Name
    );
    ADD $Summary TO $Summaries;
  END LOOP;

  RETURN $Summaries;
END;
/
```

## Error Handling in Loops

### Continue on Error

```mdl
/**
 * Process items, log errors but continue
 */
CREATE MICROFLOW Module.ProcessWithErrorHandling (
  $Items: List of Module.Item
)
RETURNS Integer
BEGIN
  DECLARE $Processed Integer = 0;
  DECLARE $Errors Integer = 0;

  LOOP $Item IN $Items
  BEGIN
    IF $Item/Data = empty THEN
      LOG WARNING NODE 'Process' 'Skipping item with empty data: ' + $Item/Code;
      SET $Errors = $Errors + 1;
    ELSE
      CHANGE $Item (Status = 'Processed');
      COMMIT $Item;
      SET $Processed = $Processed + 1;
    END IF;
  END LOOP;

  LOG INFO NODE 'Process' 'Completed: ' + $Processed + ' processed, ' + $Errors + ' errors';
  RETURN $Processed;
END;
/
```

## Delta Merge / List Matching

The correct Mendix pattern for merging imported data into existing records.

### ✅ CORRECT: Single Loop + RETRIEVE

```mdl
/**
 * Merge flagged items into quotation details by line number.
 * Uses RETRIEVE for O(N) lookup instead of nested loops.
 */
CREATE MICROFLOW Module.SUB_ApplyDelta (
  $FlaggedItems: List of Module.FlaggedItem,
  $Details: List of Module.QuotationDetail
)
RETURNS Integer
BEGIN
  DECLARE $Updated Integer = 0;

  LOOP $FlaggedItem IN $FlaggedItems
  BEGIN
    -- Find matching detail by key (acts as List Operation: Find)
    RETRIEVE $Match FROM Module.QuotationDetail
      WHERE LineNumber = $FlaggedItem/LineNumber
      LIMIT 1;

    IF $Match != empty THEN
      -- Append logic: don't overwrite, concatenate
      IF $FlaggedItem/ReviewReason != empty THEN
        CHANGE $Match (
          ReviewReason = $Match/ReviewReason + '\n' + $FlaggedItem/ReviewReason);
      END IF;
      COMMIT $Match;
      SET $Updated = $Updated + 1;
    END IF;
  END LOOP;

  RETURN $Updated;
END;
/
```

### ❌ WRONG: Nested Loops (O(N^2))

```mdl
-- ANTI-PATTERN: Creates empty list, then uses nested loops for matching
DECLARE $FlaggedItems List of Module.FlaggedItem = empty;  -- Ghost list!
LOOP $FlaggedItem IN $FlaggedItems  -- Loops over nothing!
BEGIN
  LOOP $Detail IN $Details  -- O(N^2) matching
  BEGIN
    IF $Detail/LineNumber = $FlaggedItem/LineNumber THEN
      CHANGE $Detail (ReviewReason = $FlaggedItem/ReviewReason);
    END IF;
  END LOOP;
END LOOP;
```

**Why this is wrong:**
1. `DECLARE ... = empty` creates an empty list — the loop never executes
2. Nested loops are O(N^2) — use `RETRIEVE ... LIMIT 1` for O(N) matching
3. `mxcli check` will flag both issues automatically

## Anti-Patterns

### NEVER create empty lists as loop sources

If you need to process imported data, pass it as a **microflow parameter**:

```mdl
-- WRONG: Empty list variable
DECLARE $Items List of Module.Item = empty;
LOOP $Item IN $Items BEGIN ... END LOOP;  -- Never executes!

-- CORRECT: Accept as parameter
CREATE MICROFLOW Module.ProcessItems ($Items: List of Module.Item) ...
```

### NEVER use nested LOOPs for list matching

Use `RETRIEVE ... FROM $List WHERE ... LIMIT 1` inside a single loop:

```mdl
-- WRONG: Nested loops
LOOP $Source IN $SourceList BEGIN
  LOOP $Target IN $TargetList BEGIN
    IF $Source/Key = $Target/Key THEN ...
  END LOOP;
END LOOP;

-- CORRECT: Single loop with RETRIEVE
LOOP $Source IN $SourceList BEGIN
  RETRIEVE $Match FROM Module.Target WHERE Key = $Source/Key LIMIT 1;
  IF $Match != empty THEN ...
END LOOP;
```

### NEVER overwrite when merging — use append logic

When merging notes, comments, or reasons from multiple sources:

```mdl
-- WRONG: Overwrites existing value
CHANGE $Detail (ReviewReason = $FlaggedItem/ReviewReason);

-- CORRECT: Append with separator
IF $FlaggedItem/ReviewReason != empty THEN
  CHANGE $Detail (
    ReviewReason = $Detail/ReviewReason + '\n' + $FlaggedItem/ReviewReason);
END IF;
```

## Best Practices

1. **Commit inside loops carefully**: Can cause performance issues on large sets
2. **Use batch commits**: Commit every N records for large datasets
3. **Log progress**: Add logging for long-running operations
4. **Handle errors gracefully**: Don't let one bad record stop the whole process
5. **Return counts**: Help callers know what was processed
6. **Use meaningful variable names**: `$ProcessedCount` not `$c`
7. **Pass data as parameters**: Never create empty list variables to iterate over
8. **Use RETRIEVE for matching**: Single loop + RETRIEVE, never nested loops
9. **Validate with `mxcli check`**: Run `mxcli check script.mdl` to catch anti-patterns before execution
