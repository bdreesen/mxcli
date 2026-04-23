# Mendix Nanoflow Skill

This skill provides comprehensive guidance for writing Mendix nanoflows in MDL (Mendix Definition Language) syntax.

## When to Use This Skill

Use this skill when:
- Writing CREATE NANOFLOW statements
- Debugging nanoflow syntax errors
- Converting Studio Pro nanoflows to MDL
- Understanding nanoflow vs microflow differences

## Nanoflow vs Microflow

Nanoflows execute **client-side** (browser or native app). They share the same flow structure as microflows but have important restrictions:

| Aspect | Nanoflow | Microflow |
|--------|----------|-----------|
| Execution | Client-side (browser/native) | Server-side |
| Database | Client-side offline DB | Server-side DB |
| Error handling | `$latestError` (String only) | Full error events |
| ErrorEvent | **Forbidden** | Allowed |
| Return types | No `Binary`, no `Float` | All types |
| Concurrency | N/A | `AllowConcurrentExecution` |
| Entity access | N/A | `ApplyEntityAccess` |
| Calling nanoflows | `call nanoflow` | `call microflow` (not nanoflow) |
| JavaScript actions | Allowed | Not allowed |

## Nanoflow Structure

**CRITICAL: All nanoflows MUST have JavaDoc-style documentation**

```mdl
/**
 * Nanoflow description explaining what it does
 *
 * @param $Parameter1 Description of first parameter
 * @param $Parameter2 Description of second parameter
 * @returns Description of return value
 * @since 1.0.0
 * @author Team Name
 */
create nanoflow Module.NanoflowName (
  $Parameter1: type,
  $Parameter2: type
)
returns ReturnType
[folder 'FolderPath']
begin
  -- Nanoflow logic here
  return $value;
end;
/
```

### Key Differences from Microflow Syntax

- Use `create nanoflow` (not `create microflow`)
- No `as $ReturnVariable` — nanoflows do not support `ReturnVariableName`
- No `AllowConcurrentExecution` option
- Otherwise identical syntax: same parameters, return types, body, folder, comment

### Parameter Types

Same as microflows:

```mdl
$Name: string
$count: integer
$Amount: decimal
$IsActive: boolean
$date: datetime
$Customer: Module.Entity
$ProductList: list of Module.Product
$status: enum Module.OrderStatus
```

### Allowed Return Types

```
boolean, integer, decimal, string, datetime, enumeration, object, list, void
```

**NOT allowed**: `Binary`, `Float`

## Allowed Actions

### Actions available in nanoflows

| Action | MDL Syntax | Notes |
|--------|-----------|-------|
| CreateVariable | `declare $var type = value;` | Same as microflow |
| ChangeVariable | `set $var = value;` | Same as microflow |
| CreateObject | `$var = create Module.Entity (...)` | Same as microflow |
| ChangeObject | `change $var (...)` | Same as microflow |
| CommitObject | `commit $var;` | Same as microflow |
| DeleteObject | `delete $var;` | Same as microflow |
| RollbackObject | `rollback $var;` | Same as microflow |
| RetrieveAction | `retrieve $var from ...` | From client DB |
| CreateList | `declare $list list of ... = empty;` | Same as microflow |
| ChangeList | `change list ...` | Same as microflow |
| ListOperation | `list operation ...` | Same as microflow |
| AggregateList | `aggregate list ...` | Same as microflow |
| ShowPage | `show page Module.Page(...)` | Same as microflow |
| ClosePage | `close page;` | Same as microflow |
| ShowMessage | `show message ...` | Same as microflow |
| ValidationFeedback | `validation feedback ...` | Same as microflow |
| LogMessage | `log info/warning/error ...` | Same as microflow |
| CastAction | `cast ...` | Same as microflow |
| CallMicroflow | `call microflow Module.Name(...)` | Calls server-side |
| **CallNanoflow** | `call nanoflow Module.Name(...)` | **Nanoflow-only** |
| IF/ELSE | `if ... then ... end if;` | Same as microflow |
| LOOP | `loop $var in $list begin ... end loop;` | Same as microflow |
| WHILE | `while condition begin ... end while;` | Same as microflow |

### Actions NOT available in nanoflows

These will cause validation errors if used:

| Action | Reason |
|--------|--------|
| `call java action` | Server-side JVM execution |
| `rest call` / `send rest request` | Server-side HTTP |
| `call external` | Server-side external calls |
| `download file` | Server-side file streaming |
| `generate document` | Server-side document generation |
| `import xml` / `export xml` | Server-side XML processing |
| `show home page` | Server-side navigation |
| All workflow actions | Server-side workflow engine |
| All metrics actions | Server-side telemetry |
| `send email` | Server-side email |
| `push to client` | Server-side push (nanoflows ARE client-side) |
| `execute database query` | Server-side SQL |
| `transform json` | Server-side JSON transform |

### ErrorEvent is Forbidden

Nanoflows cannot use `ErrorEvent`. Error handling uses `on error continue` or `on error { ... }` blocks on individual activities, with `$latestError` (String) as the only predefined error variable.

## Calling Nanoflows

### CALL NANOFLOW

```mdl
-- Call with result
$Result = call nanoflow Module.ValidateForm(Customer = $Customer);

-- Call without result (void nanoflow)
call nanoflow Module.RefreshUI(Page = $CurrentPage);

-- Call with error handling
$Result = call nanoflow Module.ProcessLocally(data = $data) on error continue;
```

**Important**: Same parameter matching rules as microflows — parameter names must exactly match the target nanoflow's signature (without `$` prefix). Use `describe nanoflow Module.Name` to verify.

### Calling Microflows from Nanoflows

Nanoflows can call microflows (triggers server round-trip):

```mdl
$ServerResult = call microflow Module.FetchFromServer(query = $query);
```

### Calling Nanoflows from Microflows

Microflows **cannot** call nanoflows directly. Use `call nanoflow` only inside nanoflow bodies.

## Security: GRANT/REVOKE

Control which module roles can execute a nanoflow:

```mdl
-- Grant execution permission
grant execute on nanoflow Module.NanoflowName to Module.RoleName;

-- Grant to multiple roles
grant execute on nanoflow Module.NanoflowName to Module.Role1, Module.Role2;

-- Revoke permission
revoke execute on nanoflow Module.NanoflowName from Module.RoleName;
```

**Note**: Nanoflow security is design-time only (AllowedModuleRoles). Unlike microflows, nanoflows do not have `ApplyEntityAccess`.

## Error Handling

### Predefined Variables

Nanoflows have only one predefined error variable:
- `$latestError` — String (not an object like in microflows)

### Error Handling Patterns

```mdl
-- On error continue
call microflow Module.ServerAction() on error continue;
if $latestError != empty then
  show message error 'Server call failed: ' + $latestError;
end if;

-- Custom error handler
$Result = call nanoflow Module.RiskyOperation() on error {
  log warning node 'NanoflowError' 'Operation failed: ' + $latestError;
  return $DefaultValue;
};
```

## Complete Example

```mdl
/**
 * Validates a customer form before saving
 *
 * Runs client-side for immediate feedback. Calls server
 * microflow only if local validation passes.
 *
 * @param $Customer The customer object to validate
 * @returns true if validation passes
 * @since 1.2.0
 * @author SPAM Team
 */
create nanoflow Shop.NFV_ValidateCustomerForm (
  $Customer: Shop.Customer
)
returns boolean
folder 'Customers/Validation'
begin
  -- Validate required fields
  if $Customer/Name = empty or $Customer/Name = '' then
    validation feedback $Customer attribute Name message 'Name is required';
    return false;
  end if;

  if $Customer/Email = empty or $Customer/Email = '' then
    validation feedback $Customer attribute Email message 'Email is required';
    return false;
  end if;

  -- Server-side uniqueness check
  $IsUnique = call microflow Shop.ACT_CheckEmailUnique(Email = $Customer/Email)
    on error continue;

  if $latestError != empty then
    show message warning 'Could not verify email uniqueness. Please try again.';
    return false;
  end if;

  if not $IsUnique then
    validation feedback $Customer attribute Email message 'Email already exists';
    return false;
  end if;

  return true;
end;
/
```

## Naming Conventions

Follow the same conventions as microflows with nanoflow-specific prefixes:

| Prefix | Purpose | Example |
|--------|---------|---------|
| `NFV_` | Validation nanoflow | `NFV_ValidateOrder` |
| `NFA_` | Action nanoflow | `NFA_ProcessLocally` |
| `NFS_` | Sub-nanoflow (helper) | `NFS_FormatAddress` |
| `DS_` | Data source nanoflow | `DS_GetActiveProducts` |
| `ON_` | On-change handler | `ON_StatusChanged` |

## Validation Checklist

Before executing a nanoflow script, verify:

- [ ] Uses `create nanoflow` (not `create microflow`)
- [ ] No `as $ReturnVariable` in return declaration
- [ ] Return type is not `Binary` or `Float`
- [ ] No microflow-only actions (Java, REST, workflow, import/export, etc.)
- [ ] No `ErrorEvent` in flow body
- [ ] All `call nanoflow` parameter names match target signature
- [ ] Every flow path ends with `return`
- [ ] No code after `return` statements
- [ ] All entity/association names are fully qualified
- [ ] Nanoflow ends with `/` separator

## Common Errors

| Error | Message | Fix |
|-------|---------|-----|
| CE0125 | Not supported in nanoflows | Remove microflow-only action |
| CE6051 | Web and native activities mixed | Use only web OR native actions |
| CW0701 | Deprecated list parameter | Set `UseListParameterByReference` to true |
| Parse error | Binary/Float return type | Use allowed return type |

## Quick Reference

### Nanoflow Declaration
```mdl
create nanoflow Module.Name ($Param: type) returns ReturnType
folder 'Path' begin ... end; /
```

### Call Nanoflow (inside nanoflow body)
```mdl
$result = call nanoflow Module.Name(Param = $value);
call nanoflow Module.Name(Param = $value) on error continue;
```

### Security
```mdl
grant execute on nanoflow Module.Name to Module.Role;
revoke execute on nanoflow Module.Name from Module.Role;
```

### Error Handling
```mdl
call nanoflow ... on error continue;
call nanoflow ... on error { log ...; return ...; };
```

## Related Documentation

- [Write Microflows Skill](write-microflows.md) — Server-side microflow syntax
- [MDL Syntax Guide](../../docs/02-features/mdl-syntax.md)
- [Mendix Nanoflow Documentation](https://docs.mendix.com/refguide/nanoflows/)
