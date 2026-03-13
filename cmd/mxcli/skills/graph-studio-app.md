# Building Mendix Apps on Graph Studio Data

This skill covers how to build a Mendix application that consumes data from Altair RapidMiner Graph Studio via OData. Graph Studio exposes knowledge graph data through "Data on Demand" OData V4 endpoints, which Mendix can consume as external entities.

## When to Use This Skill

- User wants to build a Mendix app that uses Graph Studio / RapidMiner graph data
- User asks about connecting to a Graph Studio OData endpoint
- User needs to create external entities from a graphmart Data on Demand endpoint
- User wants to visualize or interact with knowledge graph data in Mendix pages

## Architecture

```
┌──────────────────────────────────────────────┐
│  ALTAIR GRAPH STUDIO                         │
│                                              │
│  Graphmart  ──▶  Data on Demand Endpoint     │
│  (knowledge       (OData V4 REST API)        │
│   graph data)                                │
└───────────────────┬──────────────────────────┘
                    │ HTTP/OData4 + Basic Auth
┌───────────────────▼──────────────────────────┐
│  MENDIX APP                                  │
│                                              │
│  Constants (credentials, URL)                │
│           ▼                                  │
│  OData Client (consumed OData service)       │
│           ▼                                  │
│  External Entities (mapped to graph classes) │
│           ▼                                  │
│  Pages: Overview grids, Edit popups,         │
│         Navigation, Detail views             │
└──────────────────────────────────────────────┘
```

## Graph Studio Concepts

| Concept | Description |
|---------|-------------|
| **Graphmart** | A curated dataset in Graph Studio — the knowledge graph you query |
| **Data on Demand** | Graph Studio's OData V4 endpoint feature that exposes graphmart data |
| **Endpoint URL** | Format: `https://<host>/dataondemand/<EndpointName>/<EndpointName>` |
| **Metadata URL** | Append `/$metadata` to the endpoint URL |
| **Classes** | Graph entities exposed as OData entity sets (map to Mendix external entities) |
| **Properties** | Graph attributes exposed as OData properties (map to Mendix entity attributes) |
| **Authentication** | Basic auth with Graph Studio username/password |

## Step-by-Step: Build an App on Graph Studio Data

### Step 1: Gather Connection Info

From Graph Studio, navigate to the Data on Demand screen for your endpoint and note:

1. **Service Root URL** — e.g., `https://graphstudio.example.com/dataondemand/MyGraphmart/MyGraphmart`
2. **Metadata URL** — Service Root URL + `/$metadata`
3. **Username** — Graph Studio user with access to the graphmart
4. **Password** — Graph Studio password
5. **Entity sets** — The classes exposed by the endpoint (e.g., Customer, Product, Order)

### Step 2: Create Module and Constants

Create a dedicated module for the Graph Studio integration and store credentials in constants:

```sql
CREATE MODULE GraphData;

-- Service URL constant (environment-specific, overrideable per deployment)
CREATE CONSTANT GraphData.ServiceLocation
  TYPE String
  DEFAULT 'https://graphstudio.example.com/dataondemand/MyGraphmart/MyGraphmart';

-- Credentials (set actual values per environment)
CREATE CONSTANT GraphData.Username
  TYPE String
  DEFAULT 'api-user';

CREATE CONSTANT GraphData.Password
  TYPE String
  DEFAULT 'changeme';
```

### Step 3: Create the Consumed OData Client

```sql
CREATE ODATA CLIENT GraphData.GraphStudioClient (
  ODataVersion: OData4,
  MetadataUrl: 'https://graphstudio.example.com/dataondemand/MyGraphmart/MyGraphmart/$metadata',
  Timeout: 300,
  ServiceUrl: '@GraphData.ServiceLocation',
  UseAuthentication: Yes,
  HttpUsername: '@GraphData.Username',
  HttpPassword: '@GraphData.Password'
);
/
```

**Key properties:**
- `ODataVersion: OData4` — Graph Studio uses OData V4
- `Timeout: 300` — Graph queries can be slow for large datasets; allow 5 minutes
- `ServiceUrl: '@GraphData.ServiceLocation'` — References the constant (prefix `@` for constant references)
- `UseAuthentication: Yes` — Graph Studio requires authentication
- `HttpUsername/HttpPassword` — Reference constants with `@` prefix for environment-specific credentials

### Step 4: Create External Entities

For each graph class you want to use, create an external entity mapped to the OData client:

```sql
CREATE EXTERNAL ENTITY GraphData.Customer
FROM ODATA CLIENT GraphData.GraphStudioClient
(
  EntitySet: 'Customer',
  RemoteName: 'Customer',
  Countable: Yes,
  Creatable: No,
  Deletable: No,
  Updatable: No
)
(
  _Id: String,
  customer_key: String,
  Name: String,
  Email: String,
  Phone: String,
  Industry: String,
  Contact_Person: String
);
/

CREATE EXTERNAL ENTITY GraphData.Address
FROM ODATA CLIENT GraphData.GraphStudioClient
(
  EntitySet: 'Address',
  RemoteName: 'Address',
  Countable: Yes,
  Creatable: No,
  Deletable: No,
  Updatable: No
)
(
  address_key: String,
  Street: String,
  City: String,
  State: String,
  Country: String,
  Post_Code: String,
  Zip_Code: Integer
);
/
```

**Important notes on external entities:**
- `EntitySet` and `RemoteName` must match the OData entity set name exactly (case-sensitive)
- `Countable: Yes` enables `$count` queries (pagination support)
- For read-only graph data, set `Creatable/Deletable/Updatable: No`
- Attribute names must match the OData property names from the `$metadata`
- Graph Studio may expose properties with underscores (e.g., `customer_key`, `Contact_Person`)
- If a graph property label exceeds 128 characters, it may need a shorter name configured in the endpoint

### Step 5: Create Navigation Snippet

Build a sidebar menu snippet for navigating between entity overview pages:

```sql
CREATE SNIPPET GraphData.Entity_Menu {
  DYNAMICTEXT txtTitle (Content: 'Graph Data', RenderMode: H2)
  NAVIGATIONLIST navMenu {
    ITEM (Action: SHOW_PAGE 'GraphData.Customer_Overview') {
      DYNAMICTEXT txt1 (Content: 'Customers')
    }
    ITEM (Action: SHOW_PAGE 'GraphData.Address_Overview') {
      DYNAMICTEXT txt2 (Content: 'Addresses')
    }
  }
}
```

### Step 6: Create Overview Pages

Overview pages display graph data in a filterable data grid with the sidebar menu:

```sql
CREATE PAGE GraphData.Customer_Overview
(Title: 'Customer Overview', Layout: Atlas_Core.Atlas_Default)
{
  LAYOUTGRID layoutGrid1 {
    ROW row1 {
      COLUMN col1 (DesktopWidth: AutoFill) {
        SNIPPETCALL snippetCall1 (Snippet: GraphData.Entity_Menu)
      }
      COLUMN col2 (DesktopWidth: AutoFill) {
        DYNAMICTEXT txtHeading (Content: 'Customers', RenderMode: H2)
        DATAGRID dgCustomers (DataSource: DATABASE GraphData.Customer) {
          COLUMN col1 (Binds: Name, Caption: 'Name') {
            TEXTFILTER tfName
          }
          COLUMN col2 (Binds: Email, Caption: 'Email') {
            TEXTFILTER tfEmail
          }
          COLUMN col3 (Binds: Phone, Caption: 'Phone') {
            TEXTFILTER tfPhone
          }
          COLUMN col4 (Binds: Industry, Caption: 'Industry') {
            TEXTFILTER tfIndustry
          }
          COLUMN col5 (Binds: Contact_Person, Caption: 'Contact Person') {
            TEXTFILTER tfContact
          }
          COLUMN colActions (Binds: _Id, ShowContentAs: customContent) {
            ACTIONBUTTON btnEdit (Action: SHOW_PAGE GraphData.Customer_Detail, Style: Primary)
          }
        }
      }
    }
  }
}
```

### Step 7: Create Detail/Edit Pages

For read-only graph data, create detail popup pages:

```sql
CREATE PAGE GraphData.Customer_Detail
(
  Title: 'Customer Details',
  Layout: Atlas_Core.PopupLayout,
  Params: { $Customer: GraphData.Customer }
)
{
  LAYOUTGRID layoutGrid1 {
    ROW row1 {
      COLUMN col1 (DesktopWidth: AutoFill) {
        DATAVIEW dvCustomer (DataSource: $Customer) {
          TEXTBOX txtName (Label: 'Name', Binds: GraphData.Customer.Name, Editable: Never)
          TEXTBOX txtEmail (Label: 'Email', Binds: GraphData.Customer.Email, Editable: Never)
          TEXTBOX txtPhone (Label: 'Phone', Binds: GraphData.Customer.Phone, Editable: Never)
          TEXTBOX txtIndustry (Label: 'Industry', Binds: GraphData.Customer.Industry, Editable: Never)
          TEXTBOX txtContact (Label: 'Contact Person', Binds: GraphData.Customer.Contact_Person, Editable: Never)
          FOOTER footer1 {
            ACTIONBUTTON btnClose (Caption: 'Close', Action: CLOSE_PAGE)
          }
        }
      }
    }
  }
}
```

If the endpoint supports writes (`Creatable: Yes`, `Updatable: Yes`), use editable fields with Save/Cancel:

```sql
CREATE PAGE GraphData.Customer_NewEdit
(
  Title: 'Edit Customer',
  Layout: Atlas_Core.PopupLayout,
  Params: { $Customer: GraphData.Customer }
)
{
  LAYOUTGRID layoutGrid1 {
    ROW row1 {
      COLUMN col1 (DesktopWidth: AutoFill) {
        DATAVIEW dvCustomer (DataSource: $Customer) {
          TEXTBOX txtName (Label: 'Name', Binds: GraphData.Customer.Name)
          TEXTBOX txtEmail (Label: 'Email', Binds: GraphData.Customer.Email)
          TEXTBOX txtPhone (Label: 'Phone', Binds: GraphData.Customer.Phone)
          TEXTBOX txtIndustry (Label: 'Industry', Binds: GraphData.Customer.Industry)
          TEXTBOX txtContact (Label: 'Contact Person', Binds: GraphData.Customer.Contact_Person)
          FOOTER footer1 {
            ACTIONBUTTON btnSave (Caption: 'Save', Action: SAVE_CHANGES CLOSE_PAGE, Style: Success)
            ACTIONBUTTON btnCancel (Caption: 'Cancel', Action: CANCEL_CHANGES CLOSE_PAGE)
          }
        }
      }
    }
  }
}
```

### Step 8: Create Home Page and Navigation

```sql
CREATE PAGE GraphData.Home
(Title: 'Graph Studio Data', Layout: Atlas_Core.Atlas_Default)
{
  LAYOUTGRID layoutGrid1 {
    ROW row1 {
      COLUMN col1 (DesktopWidth: AutoFill) {
        DYNAMICTEXT txtTitle (Content: 'Graph Studio Data Explorer', RenderMode: H1)
      }
    }
    ROW row2 {
      COLUMN col1 (DesktopWidth: AutoFill) {
        ACTIONBUTTON btnCustomers (Caption: 'Customer Overview', Action: SHOW_PAGE GraphData.Customer_Overview)
        ACTIONBUTTON btnAddresses (Caption: 'Address Overview', Action: SHOW_PAGE GraphData.Address_Overview)
      }
    }
  }
}

CREATE OR REPLACE NAVIGATION Responsive
  HOME PAGE GraphData.Home
  MENU (
    MENU ITEM 'Home' PAGE GraphData.Home;
    MENU 'Graph Data' (
      MENU ITEM 'Customers' PAGE GraphData.Customer_Overview;
      MENU ITEM 'Addresses' PAGE GraphData.Address_Overview;
    );
  )
;
```

## Handling Hierarchical Data (Bill of Materials)

Graph Studio often exposes hierarchical data — e.g., a BOM (Bill of Materials) with parent-child relationships via ReferenceSet associations.

```sql
-- Parent entity
CREATE EXTERNAL ENTITY GraphData.Bom
FROM ODATA CLIENT GraphData.GraphStudioClient
(
  EntitySet: 'Bom',
  RemoteName: 'Bom',
  Countable: Yes,
  Creatable: No,
  Deletable: No,
  Updatable: No
)
(
  bom_key: String,
  Name: String,
  Product_Id: String,
  Product_Version: Decimal
);
/

-- Child entities
CREATE EXTERNAL ENTITY GraphData.Component
FROM ODATA CLIENT GraphData.GraphStudioClient
(
  EntitySet: 'Component',
  RemoteName: 'Component',
  Countable: Yes,
  Creatable: No,
  Deletable: No,
  Updatable: No
)
(
  component_key: String,
  Name: String,
  Quantity_Required: Long,
  Unit_Of_Measure: String,
  Level: Long
);
/

CREATE EXTERNAL ENTITY GraphData.Sub_Component
FROM ODATA CLIENT GraphData.GraphStudioClient
(
  EntitySet: 'Sub_Component',
  RemoteName: 'Sub_Component',
  Countable: Yes,
  Creatable: No,
  Deletable: No,
  Updatable: No
)
(
  sub_component_key: String,
  Component_Name: String,
  Quantity_Required: Long,
  Unit_Of_Measure: String,
  Level: Long
);
/
```

For hierarchical visualization (tree views), the page layout can be created in MDL but the TreeNode pluggable widget must be configured in Studio Pro:

```sql
CREATE PAGE GraphData.BomTree
(Title: 'Bill of Materials', Layout: Atlas_Core.Atlas_Default)
{
  LAYOUTGRID layoutGrid1 {
    ROW row1 {
      COLUMN col1 (DesktopWidth: AutoFill) {
        DYNAMICTEXT txtTitle (Content: 'Bill of Materials', RenderMode: H1)
      }
    }
    ROW row2 {
      COLUMN col1 (DesktopWidth: AutoFill) {
        -- TreeNode widget for hierarchical Bom -> Component -> Sub_Component display
        -- Must be configured in Studio Pro (pluggable widget not fully supported in MDL)
        DATAGRID dgBom (DataSource: DATABASE GraphData.Bom) {
          COLUMN col1 (Binds: Name, Caption: 'BOM Name') {
            TEXTFILTER tf1
          }
          COLUMN col2 (Binds: Product_Id, Caption: 'Product ID') {
            TEXTFILTER tf2
          }
          COLUMN col3 (Binds: Product_Version, Caption: 'Version') {
            NUMBERFILTER nf1 (FilterType: equal)
          }
        }
      }
    }
  }
}
```

## Performance Considerations

Graph Studio Data on Demand endpoints have specific performance characteristics to keep in mind:

1. **Timeout** — Set `Timeout: 300` (5 minutes) or higher on the OData client. Graph queries on large graphmarts can be slow.
2. **No server-side joins** — Graph Studio generates separate OData queries for joins and combines results in memory. Query one entity at a time and use Mendix pages for relationships.
3. **Pagination** — Set `Countable: Yes` on external entities to enable pagination in data grids.
4. **Filters** — Use `TEXTFILTER` and `NUMBERFILTER` widgets to push filter criteria to the OData endpoint, reducing data transfer.
5. **Read-only by default** — Most Graph Studio endpoints expose read-only data. Set `Creatable/Updatable/Deletable: No` unless the endpoint explicitly supports writes.

## Data Type Mapping

Common Graph Studio property types and their Mendix equivalents:

| Graph Studio Type | Mendix Type | Notes |
|-------------------|-------------|-------|
| String / IRI | String | Most graph properties are strings |
| Integer | Integer | 32-bit integers |
| Long | Long | 64-bit integers |
| Decimal / Double | Decimal | Floating point values |
| Boolean | Boolean | true/false |
| DateTime | DateTime | ISO 8601 date-time values |
| Date | String | May need String if no native Date mapping |

## Exploration Commands

Use these commands to inspect an existing Graph Studio integration:

```sql
-- List all consumed OData services
SHOW ODATA CLIENTS;
SHOW ODATA CLIENTS IN GraphData;

-- Inspect a specific client
DESCRIBE ODATA CLIENT GraphData.GraphStudioClient;

-- List external entities
SHOW ENTITIES IN GraphData;

-- Inspect a specific external entity with OData source details
DESCRIBE EXTERNAL ENTITY GraphData.Customer;

-- List all associations
SHOW ASSOCIATIONS IN GraphData;
```

## Checklist

Before deploying:
- [ ] Constants created for service URL, username, password
- [ ] Constants have environment-specific values (dev vs prod Graph Studio instances)
- [ ] OData client created with `ODataVersion: OData4` and `UseAuthentication: Yes`
- [ ] OData client `Timeout` set high enough for graph queries (300+)
- [ ] External entities match the entity set names from `$metadata` exactly
- [ ] External entities have correct CRUD flags (usually read-only for graph data)
- [ ] Attribute names match OData property names exactly (case-sensitive)
- [ ] Overview pages use `TEXTFILTER` / `NUMBERFILTER` for server-side filtering
- [ ] Navigation configured with home page and menu items

## MDL Limitations

| What | Status | Notes |
|------|--------|-------|
| `CREATE ODATA CLIENT` | Supported | Full support for OData V4 + Basic Auth |
| `CREATE EXTERNAL ENTITY ... FROM ODATA CLIENT` | Supported | Maps to Graph Studio entity sets |
| Overview pages with DataGrid2 | Supported | Filters push to OData endpoint |
| Detail/Edit pages | Supported | DataView with TextBox, ComboBox, etc. |
| Navigation, snippets | Supported | Full sidebar menu pattern |
| TreeNode widget | Not configurable | Page structure only; configure widget in Studio Pro |
| Hierarchical tree views | Partial | Use DataGrid as fallback; TreeNode needs Studio Pro |
