# Proposal: OData Clients, OData Services, and External Entities Support

## Overview

This proposal outlines the implementation plan for adding support for:

### Consumed (Client) OData Services
1. **OData Clients** - External API connections that consume remote OData services
2. **External Entities** - Entities backed by OData remote sources
3. **External Enumerations** - Enumerations from OData services
4. **External Actions** - Microflow activities that call OData functions/actions

### OData Services (Published)
5. **OData Services** - Expose Mendix entities as OData endpoints
6. **Published Entity Sets** - Configure which entities to expose and how
7. **Published Microflows** - Expose microflows as OData operations

## Current State

| Feature | SHOW | DESCRIBE | CREATE | ALTER | DROP |
|---------|------|----------|--------|-------|------|
| OData Client | No | No | No | No | No |
| External Entity | Partial* | Partial* | No | No | No |
| External Enumeration | No | No | No | No | No |
| External Action (microflow) | No | No | No | No | No |
| OData Service | No | No | No | No | No |
| Published Entity Set | No | No | No | No | No |
| Published Microflow | No | No | No | No | No |

*External entities appear in `SHOW ENTITIES` but without source information.

---

## Part 1: OData Clients (Consumed Services)

### BSON Storage Summary

#### OData Client
- **Document Type**: Module-level document
- **BSON $Type**: `Rest$ConsumedODataService`
- **Key Fields**: Name, ODataVersion, MetadataUrl, HTTPConfiguration, TimeoutExpression

#### External Entity
- **Document Type**: Regular entity with special `Source` field
- **BSON $Type**: `DomainModels$EntityImpl` (same as regular entities)
- **Source Type**: `Rest$ODataRemoteEntitySource`
- **Key Fields**: SourceDocument (service reference), EntitySet, RemoteName, Key, Capabilities

#### External Enumeration
- **Document Type**: Regular enumeration with special `Source` field
- **BSON $Type**: `Enumerations$Enumeration`
- **Source Type**: `Rest$ODataRemoteEnumerationSource`

#### External Action Call
- **Activity Type**: `Microflows$CallExternalAction`
- **Key Fields**: ConsumedODataService, ParameterMappings, VariableDataType

### Proposed MDL Syntax

#### 1.1 OData Client

##### SHOW
```sql
-- List all OData clients
SHOW ODATA CLIENTS;
SHOW ODATA CLIENTS IN MyModule;

-- Output columns: Module, Name, Version, ODataVersion, MetadataUrl, Validated
```

##### DESCRIBE
```sql
DESCRIBE ODATA CLIENT MyModule.SalesforceAPI;

-- Output:
-- CREATE ODATA CLIENT MyModule.SalesforceAPI (
--   Version: '1.0',
--   ODataVersion: OData4,
--   MetadataUrl: 'https://api.salesforce.com/odata/v4/$metadata',
--   Timeout: 300,
--   ProxyType: DefaultProxy
-- )
-- AUTHENTICATION CUSTOM HEADERS;
-- /
```

##### CREATE
```sql
CREATE ODATA CLIENT MyModule.SalesforceAPI (
  Version: '1.0',
  ODataVersion: OData4,
  MetadataUrl: 'https://api.salesforce.com/odata/v4/$metadata',
  Timeout: 300,
  ProxyType: DefaultProxy
)
AUTHENTICATION CUSTOM HEADERS
  'Authorization' = 'Bearer ' + @MyModule.APIToken;
```

##### ALTER
```sql
ALTER ODATA CLIENT MyModule.SalesforceAPI
  SET Timeout = 600;

ALTER ODATA CLIENT MyModule.SalesforceAPI
  SET MetadataUrl = 'https://new-api.salesforce.com/odata/v4/$metadata';

-- Refresh metadata from service
ALTER ODATA CLIENT MyModule.SalesforceAPI REFRESH METADATA;
```

##### DROP
```sql
DROP ODATA CLIENT MyModule.SalesforceAPI;
```

#### 1.2 External Entities

##### SHOW
```sql
-- Show external entities (filter from SHOW ENTITIES)
SHOW EXTERNAL ENTITIES;
SHOW EXTERNAL ENTITIES IN MyModule;

-- Output columns: Module, Name, Service, EntitySet, RemoteName, Countable, Creatable, Deletable
```

##### DESCRIBE
```sql
DESCRIBE EXTERNAL ENTITY MyModule.RemoteAccount;

-- Output:
-- CREATE EXTERNAL ENTITY MyModule.RemoteAccount
-- FROM ODATA CLIENT MyModule.SalesforceAPI
-- (
--   EntitySet: 'Accounts',
--   RemoteName: 'Account',
--   Countable: Yes,
--   Creatable: Yes,
--   Deletable: No
-- )
-- KEY (Id MAPS TO 'AccountId' AS Edm.String)
-- (
--   Id: String(200),
--   Name: String(255) MAPS TO 'AccountName' (Creatable, Updatable, Filterable),
--   Industry: String(100) MAPS TO 'Industry' (Filterable, Sortable),
--   CreatedDate: DateTime MAPS TO 'CreatedDate' (ReadOnly)
-- );
-- /
```

##### CREATE
```sql
CREATE EXTERNAL ENTITY MyModule.RemoteAccount
FROM ODATA CLIENT MyModule.SalesforceAPI
(
  EntitySet: 'Accounts',
  RemoteName: 'Account',
  Countable: Yes,
  Creatable: Yes,
  Deletable: No
)
KEY (Id MAPS TO 'AccountId' AS Edm.String)
(
  Id: String(200),
  Name: String(255) MAPS TO 'AccountName' (Creatable, Updatable, Filterable),
  Industry: String(100) MAPS TO 'Industry' (Filterable, Sortable),
  CreatedDate: DateTime MAPS TO 'CreatedDate' (ReadOnly)
);
```

##### ALTER
```sql
-- Add attribute
ALTER EXTERNAL ENTITY MyModule.RemoteAccount
  ADD Website: String(500) MAPS TO 'Website' (Creatable, Updatable);

-- Modify capabilities
ALTER EXTERNAL ENTITY MyModule.RemoteAccount
  SET Creatable = No;
```

##### DROP
```sql
DROP EXTERNAL ENTITY MyModule.RemoteAccount;
```

#### 1.3 External Enumerations

##### SHOW
```sql
SHOW EXTERNAL ENUMERATIONS;
SHOW EXTERNAL ENUMERATIONS IN MyModule;
```

##### DESCRIBE
```sql
DESCRIBE EXTERNAL ENUMERATION MyModule.AccountStatus;

-- Output:
-- CREATE EXTERNAL ENUMERATION MyModule.AccountStatus
-- FROM ODATA CLIENT MyModule.SalesforceAPI
-- MAPS TO 'AccountStatusEnum'
-- (
--   Active MAPS TO 'ACTIVE' CAPTION 'Active',
--   Inactive MAPS TO 'INACTIVE' CAPTION 'Inactive',
--   Pending MAPS TO 'PENDING' CAPTION 'Pending Review'
-- );
-- /
```

##### CREATE
```sql
CREATE EXTERNAL ENUMERATION MyModule.AccountStatus
FROM ODATA CLIENT MyModule.SalesforceAPI
MAPS TO 'AccountStatusEnum'
(
  Active MAPS TO 'ACTIVE' CAPTION 'Active',
  Inactive MAPS TO 'INACTIVE' CAPTION 'Inactive',
  Pending MAPS TO 'PENDING' CAPTION 'Pending Review'
);
```

#### 1.4 External Actions in Microflows

##### Syntax
```sql
-- Call external action (function/action from OData service)
$Result = CALL EXTERNAL ACTION MyModule.SalesforceAPI.CreateAccount (
  accountName = $CompanyName,
  accountType = 'Business'
);

-- With error handling
$Result = CALL EXTERNAL ACTION MyModule.SalesforceAPI.CreateAccount (
  accountName = $CompanyName
) ON ERROR ROLLBACK;

-- Void action (no return value)
CALL EXTERNAL ACTION MyModule.SalesforceAPI.SendNotification (
  message = $NotificationText
);
```

##### DESCRIBE Microflow Output
When describing a microflow containing external action calls:
```sql
CREATE MICROFLOW MyModule.CreateSalesforceAccount (
  CompanyName: String
) RETURNS MyModule.RemoteAccount
BEGIN
  $Result = CALL EXTERNAL ACTION MyModule.SalesforceAPI.CreateAccount (
    accountName = $CompanyName,
    accountType = 'Business'
  );
  RETURN $Result;
END;
```

---

## Part 2: OData Services (Published)

### BSON Storage Summary

#### PublishedODataService2
- **Document Type**: Module-level document
- **BSON $Type**: `ODataPublish$PublishedODataService2`
- **Key Fields**:
  - `Name`: Service name
  - `Path`: URL path for the service endpoint
  - `Namespace`: OData namespace (default: "DefaultNamespace")
  - `ServiceName`: Display name for the service
  - `Version`: Service version (default: "1.0.0")
  - `ODataVersion`: OData version ("OData4")
  - `EntitySets[]`: Array of published entity sets
  - `EntityTypes[]`: Array of entity type definitions
  - `Enumerations[]`: Published enumerations
  - `Microflows[]`: Published microflow operations
  - `AllowedModuleRoles[]`: Security roles (BY_NAME_REFERENCE)
  - `AuthenticationMicroflow`: Custom auth microflow (BY_NAME_REFERENCE)
  - `AuthenticationTypes[]`: ["Basic", "Guest", "Microflow", "Session"]
  - `Description`, `Summary`, `Documentation`: API documentation
  - `PublishAssociations`: Whether to expose navigation properties
  - `SupportsGraphQL`: Enable GraphQL support
  - `UseGeneralization`: Expose entity inheritance

#### EntitySet
- **BSON $Type**: `ODataPublish$EntitySet`
- **Key Fields**:
  - `ExposedName`: OData entity set name
  - `EntityTypePointer`: Reference to EntityType (BY_ID_REFERENCE)
  - `ReadMode`: How to read data (`ReadSource` | `CallMicroflowToRead`)
  - `InsertMode`: How to create data (`ChangeSource` | `ChangeNotSupported` | `CallMicroflowToChange`)
  - `UpdateMode`: How to update data (same options as InsertMode)
  - `DeleteMode`: How to delete data (same options as InsertMode)
  - `QueryOptions`: OData query capabilities
  - `UsePaging`: Enable server-side paging
  - `PageSize`: Page size for paging (default: 10000)

#### EntityType
- **BSON $Type**: `ODataPublish$EntityType`
- **Key Fields**:
  - `Entity`: Reference to domain model entity (BY_NAME_REFERENCE)
  - `ExposedName`: Name exposed in OData schema
  - `ChildMembers[]`: Array of published attributes and associations
  - `Description`, `Summary`: Documentation

#### PublishedMember (Abstract)
Concrete types:
- `ODataPublish$PublishedAttribute` - Exposes entity attributes
- `ODataPublish$PublishedAssociationEnd` - Exposes navigation properties
- `ODataPublish$PublishedId` - Exposes the entity key

#### Change/Read Modes
- `ODataPublish$ReadSource` - Read directly from database
- `ODataPublish$CallMicroflowToRead` - Call microflow to read (custom logic)
- `ODataPublish$ChangeSource` - Write directly to database
- `ODataPublish$ChangeNotSupported` - Operation not allowed
- `ODataPublish$CallMicroflowToChange` - Call microflow for changes (custom logic)

### Proposed MDL Syntax

#### 2.1 OData Service

##### SHOW
```sql
-- List all OData services
SHOW ODATA SERVICES;
SHOW ODATA SERVICES IN MyModule;

-- Output columns: Module, Name, Path, Version, ODataVersion, EntitySets, AuthTypes
```

##### DESCRIBE
```sql
DESCRIBE ODATA SERVICE MyModule.CustomerAPI;

-- Output:
-- CREATE ODATA SERVICE MyModule.CustomerAPI (
--   Path: '/odata/customers',
--   Version: '1.0.0',
--   ODataVersion: OData4,
--   Namespace: 'MyApp.Customers',
--   ServiceName: 'Customer Service',
--   Summary: 'API for managing customers',
--   PublishAssociations: Yes
-- )
-- AUTHENTICATION Basic, Session
-- {
--   PUBLISH ENTITY MyModule.Customer AS 'Customers' (
--     ReadMode: SOURCE,
--     InsertMode: SOURCE,
--     UpdateMode: SOURCE,
--     DeleteMode: NOT_SUPPORTED,
--     UsePaging: Yes,
--     PageSize: 100
--   )
--   EXPOSE (
--     Id AS 'customerId',
--     Name AS 'customerName' (Filterable, Sortable),
--     Email AS 'email',
--     CreatedDate AS 'createdAt' (ReadOnly)
--   );
--
--   PUBLISH ENTITY MyModule.Order AS 'Orders' (
--     ReadMode: MICROFLOW MyModule.GetOrdersForOData,
--     InsertMode: MICROFLOW MyModule.CreateOrderViaOData,
--     UpdateMode: NOT_SUPPORTED,
--     DeleteMode: NOT_SUPPORTED
--   )
--   EXPOSE (*);
-- }
-- /
```

##### CREATE
```sql
CREATE ODATA SERVICE MyModule.CustomerAPI (
  Path: '/odata/customers',
  Version: '1.0.0',
  ODataVersion: OData4,
  Namespace: 'MyApp.Customers',
  ServiceName: 'Customer Service',
  Summary: 'API for managing customers',
  PublishAssociations: Yes
)
AUTHENTICATION Basic, Session
{
  PUBLISH ENTITY MyModule.Customer AS 'Customers' (
    ReadMode: SOURCE,
    InsertMode: SOURCE,
    UpdateMode: SOURCE,
    DeleteMode: NOT_SUPPORTED,
    UsePaging: Yes,
    PageSize: 100
  )
  EXPOSE (
    Id,
    Name (Filterable, Sortable),
    Email,
    CreatedDate (ReadOnly)
  );
};
```

##### ALTER
```sql
-- Add authentication type
ALTER ODATA SERVICE MyModule.CustomerAPI
  ADD AUTHENTICATION Microflow MyModule.ValidateAPIKey;

-- Change version
ALTER ODATA SERVICE MyModule.CustomerAPI
  SET Version = '2.0.0';

-- Add entity to service
ALTER ODATA SERVICE MyModule.CustomerAPI
  ADD ENTITY MyModule.Invoice AS 'Invoices' (ReadMode: SOURCE);
```

##### DROP
```sql
DROP ODATA SERVICE MyModule.CustomerAPI;
```

#### 2.2 Published Microflows

```sql
-- Within ODATA SERVICE block:
{
  PUBLISH MICROFLOW MyModule.CalculateDiscount AS 'CalculateDiscount' (
    $CustomerId: Integer,
    $Amount: Decimal
  ) RETURNS Decimal;
}
```

#### 2.3 Security Access

```sql
-- Grant access to OData service
GRANT ACCESS ON ODATA SERVICE MyModule.CustomerAPI TO MyModule.Admin, MyModule.User;

-- Revoke access
REVOKE ACCESS ON ODATA SERVICE MyModule.CustomerAPI FROM MyModule.Guest;
```

---

## Implementation Plan

### Phase 1: Read Support (SHOW/DESCRIBE)

#### 1.1 SDK Layer (`sdk/`)
- [ ] Add `ConsumedODataService` struct in new `sdk/rest/` package
- [ ] Add `ODataService` struct in `sdk/odatapublish/` package
- [ ] Add `ODataRemoteEntitySource` struct
- [ ] Add `ODataRemoteEnumerationSource` struct
- [ ] Add parser support in `sdk/mpr/parser.go`
- [ ] Add `ListConsumedODataServices()` to reader interface
- [ ] Add `ListODataServices()` to reader interface

#### 1.2 Grammar (`mdl/grammar/`)
- [ ] Add tokens: `ODATA`, `CLIENT`, `EXTERNAL`, `MAPS`, `AUTHENTICATION`, `PUBLISHED`, `PUBLISH`, `EXPOSE`, `SOURCE`
- [ ] Add `showODataClientsStatement` rule
- [ ] Add `describeODataClientStatement` rule
- [ ] Add `showExternalEntitiesStatement` rule
- [ ] Add `describeExternalEntityStatement` rule
- [ ] Add `showODataServicesStatement` rule
- [ ] Add `describeODataServiceStatement` rule

#### 1.3 AST (`mdl/ast/`)
- [ ] Add `ShowODataClients` to `ShowObjectType`
- [ ] Add `DescribeODataClient` to `DescribeObjectType`
- [ ] Add `ShowExternalEntities` to `ShowObjectType`
- [ ] Add `DescribeExternalEntity` to `DescribeObjectType`
- [ ] Add `ShowODataServices` to `ShowObjectType`
- [ ] Add `DescribeODataService` to `DescribeObjectType`

#### 1.4 Visitor (`mdl/visitor/`)
- [ ] Handle new SHOW/DESCRIBE statement types

#### 1.5 Executor (`mdl/executor/`)
- [ ] Implement `showODataClients()`
- [ ] Implement `describeODataClient()`
- [ ] Modify `showEntities()` to indicate external entities
- [ ] Implement `describeExternalEntity()`
- [ ] Implement `showODataServices()`
- [ ] Implement `describeODataService()`

#### 1.6 Catalog (`mdl/catalog/`)
- [ ] Add `CATALOG.ODATA_CLIENTS` table
- [ ] Add `CATALOG.ODATA_SERVICES` table
- [ ] Add `IsExternal` column to `CATALOG.ENTITIES`

### Phase 2: Write Support (CREATE/ALTER/DROP)

#### 2.1 SDK Layer
- [ ] Add `ConsumedODataService` to writer
- [ ] Add `ODataService` to writer
- [ ] Add BSON serialization for `Rest$ConsumedODataService`
- [ ] Add BSON serialization for `ODataPublish$PublishedODataService2`
- [ ] Add BSON serialization for `Rest$ODataRemoteEntitySource`

#### 2.2 Grammar
- [ ] Add `createODataClientStatement` rule
- [ ] Add `alterODataClientStatement` rule
- [ ] Add `dropODataClientStatement` rule
- [ ] Add `createExternalEntityStatement` rule
- [ ] Add `createODataServiceStatement` rule
- [ ] Add `alterODataServiceStatement` rule
- [ ] Add `dropODataServiceStatement` rule
- [ ] Add authentication clause rules

#### 2.3 AST
- [ ] Add `CreateODataClientStmt` struct
- [ ] Add `AlterODataClientStmt` struct
- [ ] Add `DropODataClientStmt` struct
- [ ] Add `CreateExternalEntityStmt` struct
- [ ] Add `ODataAttributeMapping` struct
- [ ] Add `CreateODataServiceStmt` struct
- [ ] Add `AlterODataServiceStmt` struct
- [ ] Add `DropODataServiceStmt` struct

#### 2.4 Executor
- [ ] Implement `createODataClient()`
- [ ] Implement `alterODataClient()`
- [ ] Implement `dropODataClient()`
- [ ] Implement `createExternalEntity()`
- [ ] Implement `createODataService()`
- [ ] Implement `alterODataService()`
- [ ] Implement `dropODataService()`

### Phase 3: Microflow External Actions

#### 3.1 SDK Layer
- [ ] Add `CallExternalAction` activity struct
- [ ] Add `ExternalActionParameterMapping` struct
- [ ] Add parser support for `Microflows$CallExternalAction`

#### 3.2 Grammar
- [ ] Add `CALL EXTERNAL ACTION` to microflow statements

#### 3.3 AST
- [ ] Add `CallExternalActionStmt` struct

#### 3.4 Executor
- [ ] Implement external action call execution
- [ ] Add to microflow DESCRIBE output

---

## BSON Field Mappings

### ConsumedODataService (OData Client)

| MDL Property | BSON Field | Type | Default |
|--------------|------------|------|---------|
| Name | Name | STRING | (required) |
| Version | Version | STRING | "1.0" |
| ODataVersion | ODataVersion | ENUM | "OData4" |
| MetadataUrl | MetadataUrl | STRING | (required) |
| Timeout | TimeoutExpression | STRING | "300" |
| ProxyType | ProxyType | ENUM | "DefaultProxy" |
| Description | Description | STRING | "" |

### ODataRemoteEntitySource

| MDL Property | BSON Field | Type |
|--------------|------------|------|
| EntitySet | EntitySet | STRING |
| RemoteName | RemoteName | STRING |
| Countable | Countable | BOOLEAN |
| Creatable | Creatable | BOOLEAN |
| Deletable | Deletable | BOOLEAN |
| SkipSupported | SkipSupported | BOOLEAN |
| TopSupported | TopSupported | BOOLEAN |

### ODataMappedValue (Attribute)

| MDL Modifier | BSON Field | Type |
|--------------|------------|------|
| MAPS TO 'name' | RemoteName | STRING |
| AS Edm.Type | RemoteType | STRING |
| Creatable | Creatable | BOOLEAN |
| Updatable | Updatable | BOOLEAN |
| Filterable | Filterable | BOOLEAN |
| Sortable | Sortable | BOOLEAN |
| ReadOnly | Creatable=false, Updatable=false | - |

### PublishedODataService2

| MDL Property | BSON Field | Type | Default |
|--------------|------------|------|---------|
| Name | Name | STRING | (required) |
| Path | Path | STRING | "" |
| Version | Version | STRING | "1.0.0" |
| ODataVersion | ODataVersion | ENUM | "OData4" |
| Namespace | Namespace | STRING | "DefaultNamespace" |
| ServiceName | ServiceName | STRING | "" |
| Summary | Summary | STRING | "" |
| Description | Description | STRING | "" |
| PublishAssociations | PublishAssociations | BOOLEAN | true |
| SupportsGraphQL | SupportsGraphQL | BOOLEAN | false |
| UseGeneralization | UseGeneralization | BOOLEAN | false |
| AuthenticationTypes | AuthenticationTypes | ENUM[] | [] |
| AllowedModuleRoles | AllowedModuleRoles | BY_NAME_REFERENCE[] | [] |
| AuthenticationMicroflow | AuthenticationMicroflow | BY_NAME_REFERENCE | null |

### EntitySet (Published)

| MDL Property | BSON Field | Type | Default |
|--------------|------------|------|---------|
| ExposedName | ExposedName | STRING | "" |
| EntityType | EntityTypePointer | BY_ID_REFERENCE | null |
| ReadMode | ReadMode | PART | null |
| InsertMode | InsertMode | PART | null |
| UpdateMode | UpdateMode | PART | null |
| DeleteMode | DeleteMode | PART | null |
| UsePaging | UsePaging | BOOLEAN | false |
| PageSize | PageSize | INTEGER | 10000 |

### EntityType (Published)

| MDL Property | BSON Field | Type |
|--------------|------------|------|
| Entity | Entity | BY_NAME_REFERENCE |
| ExposedName | ExposedName | STRING |
| ChildMembers | ChildMembers | PART[] |
| Description | Description | STRING |
| Summary | Summary | STRING |

---

## Verification Commands

```bash
# Phase 1 verification - OData Clients
./bin/mxcli -p app.mpr -c "SHOW ODATA CLIENTS"
./bin/mxcli -p app.mpr -c "DESCRIBE ODATA CLIENT MyModule.API"
./bin/mxcli -p app.mpr -c "SHOW EXTERNAL ENTITIES"
./bin/mxcli -p app.mpr -c "DESCRIBE EXTERNAL ENTITY MyModule.RemoteEntity"

# Phase 1 verification - OData Services
./bin/mxcli -p app.mpr -c "SHOW ODATA SERVICES"
./bin/mxcli -p app.mpr -c "DESCRIBE ODATA SERVICE MyModule.CustomerAPI"

# Phase 2 verification (round-trip)
./bin/mxcli -p app.mpr -c "DESCRIBE ODATA CLIENT MyModule.API" > /tmp/client.mdl
./bin/mxcli check /tmp/client.mdl
./bin/mxcli -p app.mpr exec /tmp/client.mdl

./bin/mxcli -p app.mpr -c "DESCRIBE ODATA SERVICE MyModule.CustomerAPI" > /tmp/published.mdl
./bin/mxcli check /tmp/published.mdl
./bin/mxcli -p app.mpr exec /tmp/published.mdl

# Phase 3 verification
./bin/mxcli -p app.mpr -c "DESCRIBE MICROFLOW MyModule.CallExternalMF"
# Should show CALL EXTERNAL ACTION statements
```

---

## Dependencies

- Requires understanding of `Microflows$HttpConfiguration` structure for authentication
- May need constant references for proxy settings
- HTTPConfiguration is a required embedded PART - need to generate valid defaults
- OData services have complex embedded structures (EntitySet, EntityType, etc.)

---

## Open Questions

1. **Metadata Refresh**: Should `ALTER ODATA CLIENT ... REFRESH METADATA` actually fetch from the URL, or just clear cached metadata?

2. **Authentication**: How detailed should authentication configuration be? Options:
   - Basic: Just reference constants for credentials
   - Full: Support inline header expressions, OAuth flows, etc.

3. **Import vs Create**: Should we support importing from a metadata URL directly?
   ```sql
   IMPORT ODATA CLIENT FROM 'https://api.example.com/$metadata'
     INTO MyModule.ExampleAPI;
   ```

4. **Attribute Auto-Generation**: For external entities, should CREATE support auto-generating attributes from service metadata?

5. **Service Wizard**: Should there be a shorthand for common patterns?
   ```sql
   -- Quick publish with defaults
   PUBLISH ENTITY MyModule.Customer TO ODATA SERVICE MyModule.API;
   ```

---

## References

- `reference/mendixmodellib/reflection-data/11.6.0-structures.json` - Type definitions
- `mx-test-projects/QueryDemoApp-main/QueryDemoApp.mpr` - Working examples
- Types: `Rest$ConsumedODataService`, `ODataPublish$PublishedODataService2`, `ODataPublish$EntitySet`, `ODataPublish$EntityType`

---

## Timeline Estimate

| Phase | Effort | Dependencies |
|-------|--------|--------------|
| Phase 1a (OData Client Read) | 2-3 days | None |
| Phase 1b (OData Service Read) | 2-3 days | None |
| Phase 2a (OData Client Write) | 3-4 days | Phase 1a |
| Phase 2b (OData Service Write) | 4-5 days | Phase 1b |
| Phase 3 (Microflow Actions) | 2-3 days | Phase 1a |

Total: ~13-18 days
