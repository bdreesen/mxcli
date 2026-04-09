# OData and Integration Statements

Statements for managing OData services, external entities, and browsing service contracts.

Mendix supports consuming and publishing OData services. Consumed services (OData clients) connect to remote OData endpoints and make external entity types available for use in your application. Published services expose your domain model entities as OData endpoints for other applications to consume.

## OData Client Statements

| Statement | Description |
|-----------|-------------|
| [CREATE ODATA CLIENT](create-odata-client.md) | Create a consumed OData service (auto-fetches $metadata) |
| [ALTER ODATA CLIENT](alter-odata-client.md) | Modify OData client properties |
| [DROP ODATA CLIENT](drop-odata-client.md) | Remove a consumed OData service |

## OData Service Statements (Published)

| Statement | Description |
|-----------|-------------|
| [CREATE ODATA SERVICE](create-odata-service.md) | Publish entities as an OData endpoint |
| [ALTER ODATA SERVICE](alter-odata-service.md) | Modify published service properties |
| [DROP ODATA SERVICE](drop-odata-service.md) | Remove a published OData service |

## External Entity Statements

| Statement | Description |
|-----------|-------------|
| [CREATE EXTERNAL ENTITY](create-external-entity.md) | Import an entity from a consumed OData service |
| [CREATE EXTERNAL ENTITIES](#create-external-entities) | Bulk-create external entities from cached $metadata |

## CREATE EXTERNAL ENTITIES

Bulk-create external entities from a consumed OData service's cached `$metadata`.

```sql
-- Create all entity types from the contract
CREATE EXTERNAL ENTITIES FROM Module.Service;

-- Create into a specific module
CREATE EXTERNAL ENTITIES FROM Module.Service INTO TargetModule;

-- Filter to specific entities
CREATE EXTERNAL ENTITIES FROM Module.Service ENTITIES (Customer, Order);

-- Idempotent — update existing entities
CREATE OR MODIFY EXTERNAL ENTITIES FROM Module.Service;
```

## Contract Browsing Statements

Browse available assets from cached service contracts without network access.

| Statement | Description |
|-----------|-------------|
| `SHOW CONTRACT ENTITIES FROM Module.Service` | List entity types from cached $metadata |
| `SHOW CONTRACT ACTIONS FROM Module.Service` | List actions/functions from cached $metadata |
| `DESCRIBE CONTRACT ENTITY Module.Service.Entity` | Show entity properties, types, keys |
| `DESCRIBE CONTRACT ENTITY Module.Service.Entity FORMAT mdl` | Generate CREATE EXTERNAL ENTITY |
| `DESCRIBE CONTRACT ACTION Module.Service.Action` | Show action parameters and return type |
| `SHOW CONTRACT CHANNELS FROM Module.Service` | List channels from cached AsyncAPI |
| `SHOW CONTRACT MESSAGES FROM Module.Service` | List messages from cached AsyncAPI |
| `DESCRIBE CONTRACT MESSAGE Module.Service.Message` | Show message payload properties |

## Related Show/Describe Statements

| Statement | Syntax |
|-----------|--------|
| Show OData clients | `SHOW ODATA CLIENTS [IN module]` |
| Show OData services | `SHOW ODATA SERVICES [IN module]` |
| Show external entities | `SHOW EXTERNAL ENTITIES [IN module]` |
| Show external actions | `SHOW EXTERNAL ACTIONS [IN module]` |
| Describe OData client | `DESCRIBE ODATA CLIENT Module.Name` |
| Describe OData service | `DESCRIBE ODATA SERVICE Module.Name` |
| Describe external entity | `DESCRIBE EXTERNAL ENTITY Module.Name` |

## Catalog Tables

After `REFRESH CATALOG`, the following tables are available for SQL queries:

| Table | Contents |
|-------|----------|
| `CATALOG.ODATA_CLIENTS` | Consumed OData services |
| `CATALOG.ODATA_SERVICES` | Published OData services |
| `CATALOG.EXTERNAL_ENTITIES` | Imported external entities |
| `CATALOG.EXTERNAL_ACTIONS` | External actions used in microflows |
| `CATALOG.CONTRACT_ENTITIES` | Entity types from cached $metadata |
| `CATALOG.CONTRACT_ACTIONS` | Actions/functions from cached $metadata |
| `CATALOG.CONTRACT_MESSAGES` | Messages from cached AsyncAPI |

### Example: Find available entities not yet imported

```sql
SELECT ce.EntityName, ce.ServiceQualifiedName, ce.PropertyCount
FROM CATALOG.CONTRACT_ENTITIES ce
LEFT JOIN CATALOG.EXTERNAL_ENTITIES ee
  ON ce.ServiceQualifiedName = ee.ServiceName
  AND ce.EntityName = ee.RemoteName
WHERE ee.Id IS NULL;
```
