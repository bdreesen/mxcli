# MDL Language Specification

MDL (Mendix Definition Language) is a SQL-like domain-specific language for defining and manipulating Mendix application models. This specification describes the language syntax and its mapping to various backends.

## Documents

1. [Language Reference](./01-language-reference.md) - Complete MDL syntax and semantics
2. [Data Types](./02-data-types.md) - MDL data type system
3. [Domain Model](./03-domain-model.md) - Entities, attributes, associations
4. [BSON Mapping](./10-bson-mapping.md) - Mapping to MPR file format
5. [Model SDK Mapping](./11-model-sdk-mapping.md) - Mapping to modelsdk-go library

For the comprehensive quick reference of all MDL statements, see [MDL Quick Reference](../MDL_QUICK_REFERENCE.md).

## Quick Reference

```sql
-- Connection
CONNECT LOCAL '/path/to/app.mpr';
DISCONNECT;
STATUS;

-- Query
SHOW MODULES;
SHOW ENTITIES [IN ModuleName];
SHOW STRUCTURE [DEPTH 1|2|3] [IN Module] [ALL];
DESCRIBE ENTITY Module.EntityName;

-- Domain Model
CREATE PERSISTENT ENTITY Module.Name (
  AttrName: Type [NOT NULL] [UNIQUE] [DEFAULT value]
);
ALTER ENTITY Module.Name ADD (NewAttr: String(200));
DROP ENTITY Module.Name;

-- Microflows
CREATE MICROFLOW Module.Name BEGIN ... END;
DESCRIBE MICROFLOW Module.Name;

-- Pages
CREATE PAGE Module.Name (Title: 'Title', Layout: Module.Layout) { ... };
ALTER PAGE Module.Name { SET Caption = 'New' ON btnSave; };

-- Security
GRANT Module.Role ON Module.Entity (CREATE, DELETE, READ *, WRITE *);
GRANT EXECUTE ON MICROFLOW Module.Name TO Module.Role;

-- External SQL
SQL CONNECT postgres 'dsn' AS alias;
SQL alias SELECT * FROM table;

-- Navigation, Settings, Business Events, Java Actions
-- See MDL Quick Reference for full syntax
```

## Design Principles

1. **SQL-like syntax** - Familiar to developers with database experience
2. **Case-insensitive keywords** - `CREATE`, `create`, `Create` are equivalent
3. **Qualified names** - `Module.Element` format for cross-module references
4. **Statement terminators** - `;` or `/` to end statements
5. **Multi-line support** - Statements can span multiple lines
6. **Documentation comments** - `/** ... */` for element documentation
