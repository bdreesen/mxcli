# Case Study: Recreating MxGraphStudioDemo in MDL

## Overview

The **MxGraphStudioDemo** project (Mendix 11.6.3) is a small app that connects to a Graph Studio OData endpoint to display and edit customer/address data and visualize a hierarchical Bill of Materials (BOM) tree. This document analyzes what MDL can express today and provides scripts to recreate the supported portions.

## Project Structure

| Module | Role | Source |
|--------|------|--------|
| **OdataPlm** | App module — external entities from OData, overview/edit pages, BOM tree | Custom |
| **Main** | App module — home page, REST integration prototype, constants | Custom |
| Administration | Marketplace module — user management | v4.3.2 |
| Atlas_Core | Marketplace module — layouts and building blocks | v4.1.3 |
| Atlas_Web_Content | Marketplace module — page templates | v4.1.0 |
| DataWidgets | Marketplace module — DataGrid2, filters | v3.5.0 |
| FeedbackModule | Marketplace module — feedback widget | v4.0.3 |

## Architecture

```
Graph Studio OData API
    |
    v
[OdataPlm.MxPlmOdataApiClient] -- consumed OData service
    |
    v
External Entities: Customer, Address, Bom, Component, Sub_Component
    |
    v
Pages: Customer_Overview, Address_Overview, Customer_NewEdit, Address_NewEdit, BomTree
    |
    v
Navigation: Responsive -> Main.Home_Web -> buttons to OdataPlm pages
```

### Associations

- `OdataPlm.Address_2`: Customer -> Address (Reference)
- `OdataPlm.Components`: Bom -> Component (ReferenceSet)
- `OdataPlm.Sub_Components`: Component -> Sub_Component (ReferenceSet)

## MDL Coverage Analysis

### Fully Supported

| Feature | MDL Syntax | Status |
|---------|-----------|--------|
| External entities | `CREATE EXTERNAL ENTITY` | Supported |
| Persistent entities | `CREATE PERSISTENT ENTITY` | Supported |
| Non-persistent entities | `CREATE NON-PERSISTENT ENTITY` | Supported |
| Enumerations | `CREATE ENUMERATION` | Supported |
| Overview pages (DataGrid2 + filters) | `CREATE PAGE ... { DATAGRID ... }` | Supported |
| Edit pages (DataView + TextBox) | `CREATE PAGE ... { DATAVIEW ... }` | Supported |
| Snippets | `CREATE SNIPPET` | Supported |
| Navigation lists | `NAVIGATIONLIST` | Supported |
| Layout grids | `LAYOUTGRID / ROW / COLUMN` | Supported |
| Action buttons | `ACTIONBUTTON` | Supported |
| Text/Number filters | `TEXTFILTER / NUMBERFILTER` | Supported |
| Dynamic text / headings | `DYNAMICTEXT` | Supported |
| ComboBox | `COMBOBOX` | Supported |
| Navigation profiles | `CREATE OR REPLACE NAVIGATION` | Supported |
| Constants | `CREATE CONSTANT` | Supported |
| Page parameters | `Params: { $Var: Entity }` | Supported |
| Show page actions | `Action: SHOW_PAGE Module.Page` | Supported |
| Save/Cancel actions | `Action: SAVE_CHANGES / CANCEL_CHANGES` | Supported |
| Consumed OData clients | `CREATE ODATA CLIENT` | Supported |
| External entities with OData source | `CREATE EXTERNAL ENTITY ... FROM ODATA CLIENT` | Supported |

### Not Supported

| Feature | What's Missing | Workaround |
|---------|---------------|------------|
| **TreeNode widget** | Pluggable widget not in MDL grammar or SDK. `DESCRIBE PAGE` outputs `TREENODE treeNode1` but cannot configure it. | No workaround — requires Studio Pro to configure. Can create the page structure and add TreeNode manually. |
| **RestOperationCallAction** | Microflow activity not implemented in describe/create. Shows as `-- Unsupported action type`. | Must create REST call microflows in Studio Pro. |

### Partially Supported

| Feature | Notes |
|---------|-------|
| **BomTree page** | The page layout (LayoutGrid, DynamicText heading) is fully supported. Only the `TREENODE` widget inside it cannot be configured — MDL outputs just the widget name with no properties. |
| **Main.GetCustomers microflow** | The `RETRIEVE` and `RETURN` are supported, but the `RestOperationCallAction` at the start is emitted as a comment. |

## MDL Scripts to Recreate the App

### Step 1: Create Modules

```sql
-- The OdataPlm and Main modules must exist first.
-- Marketplace modules (Atlas_Core, Administration, etc.) are assumed to be
-- installed via the Mendix Marketplace in Studio Pro.
CREATE MODULE OdataPlm;
CREATE MODULE Main;
```

### Step 2: Create Constants (needed by OData client)

```sql
CREATE CONSTANT Main.MxPlmGraphClient_graphmart AS String = 'http%3A%2F%2Fcambridgesemantics.com%2F...';
CREATE CONSTANT Main.MxPlmGraphClient_graphstudio_api AS String = 'https://graphstudio.mendixdemo.com/sp...';
CREATE CONSTANT Main.MxPlmGraphClient_password AS String = 'Welcome1!';
CREATE CONSTANT Main.MxPlmGraphClient_username AS String = 'Administrator';
CREATE CONSTANT OdataPlm.MxPlmOdataApiClient_Location AS String = 'https://graphstudio.mendixdemo.com/da...';
```

### Step 3: Create the Consumed OData Client

```sql
CREATE ODATA CLIENT OdataPlm.MxPlmOdataApiClient (
  ODataVersion: OData4,
  MetadataUrl: 'https://graphstudio.mendixdemo.com/dataondemand/Mx-PLM-example/MxPlmExample/$metadata',
  Timeout: 300,
  ServiceUrl: '@OdataPlm.MxPlmOdataApiClient_Location',
  UseAuthentication: Yes,
  HttpUsername: '@Main.MxPlmGraphClient_username',
  HttpPassword: '@Main.MxPlmGraphClient_password'
);
/
```

### Step 4: External Entities (OdataPlm)

```sql
CREATE EXTERNAL ENTITY OdataPlm.Customer
FROM ODATA CLIENT OdataPlm.MxPlmOdataApiClient
(
  EntitySet: 'Customer',
  RemoteName: 'Customer',
  Countable: Yes,
  Creatable: No,
  Deletable: No,
  Updatable: No
)
(
  customer_key: String,
  Industry: String,
  address_key: String,
  Name: String,
  _Id: String,
  Phone: String,
  Email: String,
  Contact_Person: String
);
/

CREATE EXTERNAL ENTITY OdataPlm.Address
FROM ODATA CLIENT OdataPlm.MxPlmOdataApiClient
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
  Country: String,
  State: String,
  Street: String,
  Post_Code: String,
  City: String,
  Zip_Code: Integer
);
/

CREATE EXTERNAL ENTITY OdataPlm.Bom
FROM ODATA CLIENT OdataPlm.MxPlmOdataApiClient
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
  Product_Version: Decimal,
  Product_Id: String,
  _Id: String,
  Name: String
);
/

CREATE EXTERNAL ENTITY OdataPlm.Component
FROM ODATA CLIENT OdataPlm.MxPlmOdataApiClient
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
  Level: Long,
  _Id: String,
  Unit_Of_Measure: String
);
/

CREATE EXTERNAL ENTITY OdataPlm.Sub_Component
FROM ODATA CLIENT OdataPlm.MxPlmOdataApiClient
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
  Unit_Of_Measure: String,
  Quantity_Required: Long,
  Component_Name: String,
  Level: Long,
  Component_Id: String
);
/
```

### Step 5: Test Entity and Enumeration (OdataPlm)

```sql
CREATE ENUMERATION OdataPlm.Test (
  wewe 'wewe'
);
/

CREATE PERSISTENT ENTITY OdataPlm.TestAbc (
  test: Enumeration(OdataPlm.Test)
);
/
```

### Step 6: Non-Persistent Entities (Main)

```sql
CREATE NON-PERSISTENT ENTITY Main.GetCustomerResponse ();
/

CREATE NON-PERSISTENT ENTITY Main.Customer (
  Customer__type: String,
  Customer_Value: String,
  CustomerId__type: String,
  CustomerId_Value: String,
  CustomerName__type: String,
  CustomerName_Value: String
);
/

CREATE NON-PERSISTENT ENTITY Main.Var (
  Value: String
);
/
```

### Step 7: Associations (Main)

```sql
-- Auto-created for OdataPlm associations via consumed OData service.
-- Main module associations:
CREATE ASSOCIATION Main.Customer_GetCustomerResponse (
  Main.Customer -> Main.GetCustomerResponse
);
/

CREATE ASSOCIATION Main.Var_GetCustomerResponse (
  Main.Var -> Main.GetCustomerResponse
);
/
```

### Step 8: Snippet (Entity_Menu)

```sql
CREATE SNIPPET OdataPlm.Entity_Menu {
  DYNAMICTEXT text1 (Content: 'Entities', RenderMode: H2)
  NAVIGATIONLIST navigationList1 {
    ITEM (Action: SHOW_PAGE 'OdataPlm.Address_Overview') {
      DYNAMICTEXT text2 (Content: 'Address')
    }
    ITEM (Action: SHOW_PAGE 'OdataPlm.Customer_Overview') {
      DYNAMICTEXT text3 (Content: 'Customer')
    }
  }
}
```

### Step 9: Pages (OdataPlm)

#### Customer Overview

```sql
CREATE PAGE OdataPlm.Customer_Overview
(Title: 'Customer Overview', Layout: Atlas_Core.Atlas_Default)
{
  LAYOUTGRID layoutGrid1 {
    ROW row1 {
      COLUMN col1 (DesktopWidth: AutoFill) {
        SNIPPETCALL snippetCall1 (Snippet: OdataPlm.Entity_Menu)
      }
      COLUMN col2 (DesktopWidth: AutoFill) {
        DYNAMICTEXT text1 (Content: 'Customer', RenderMode: H2)
        DATAGRID dataGrid2_1 (DataSource: DATABASE OdataPlm.Customer) {
          COLUMN col1 (Attribute: customer_key, Caption: 'customer key') {
            TEXTFILTER textFilter2
          }
          COLUMN col2 (Attribute: Industry, Caption: 'Industry') {
            TEXTFILTER textFilter3
          }
          COLUMN col3 (Attribute: address_key, Caption: 'address key') {
            TEXTFILTER textFilter4
          }
          COLUMN col4 (Attribute: Name, Caption: 'Name') {
            TEXTFILTER textFilter5
          }
          COLUMN col5 (Attribute: _Id, Caption: 'Id') {
            TEXTFILTER textFilter6
          }
          COLUMN col6 (Attribute: Phone, Caption: 'Phone') {
            TEXTFILTER textFilter7
          }
          COLUMN col7 (Attribute: Email, Caption: 'Email') {
            TEXTFILTER textFilter8
          }
          COLUMN col8 (Attribute: Contact_Person, Caption: 'Contact Person') {
            TEXTFILTER textFilter1
          }
          COLUMN col9 (Attribute: customer_key, ShowContentAs: customContent) {
            ACTIONBUTTON actionButton1 (Action: SHOW_PAGE OdataPlm.Customer_NewEdit, Style: Primary)
          }
        }
      }
    }
  }
}
```

#### Customer Edit (Popup)

```sql
CREATE PAGE OdataPlm.Customer_NewEdit
(Title: 'Edit Customer', Layout: Atlas_Core.PopupLayout, Params: { $Customer: OdataPlm.Customer })
{
  LAYOUTGRID layoutGrid1 {
    ROW row1 {
      COLUMN col1 (DesktopWidth: AutoFill) {
        DATAVIEW dataView1 (DataSource: $Customer) {
          TEXTBOX textBox1 (Label: 'customer key', Attribute: OdataPlm.Customer.customer_key)
          TEXTBOX textBox2 (Label: 'Industry', Attribute: OdataPlm.Customer.Industry)
          TEXTBOX textBox3 (Label: 'address key', Attribute: OdataPlm.Customer.address_key)
          TEXTBOX textBox4 (Label: 'Name', Attribute: OdataPlm.Customer.Name)
          TEXTBOX textBox5 (Label: 'Id', Attribute: OdataPlm.Customer._Id)
          TEXTBOX textBox6 (Label: 'Phone', Attribute: OdataPlm.Customer.Phone)
          TEXTBOX textBox7 (Label: 'Email', Attribute: OdataPlm.Customer.Email)
          TEXTBOX textBox8 (Label: 'Contact Person', Attribute: OdataPlm.Customer.Contact_Person)
          COMBOBOX comboBox1 (Attribute: OdataPlm.Address.address_key)
          FOOTER footer1 {
            ACTIONBUTTON actionButton1 (Caption: 'Save', Action: SAVE_CHANGES CLOSE_PAGE, Style: Success)
            ACTIONBUTTON actionButton2 (Caption: 'Cancel', Action: CANCEL_CHANGES CLOSE_PAGE)
          }
        }
      }
    }
  }
}
```

#### Address Overview

```sql
CREATE PAGE OdataPlm.Address_Overview
(Title: 'Address Overview', Layout: Atlas_Core.Atlas_Default)
{
  LAYOUTGRID layoutGrid1 {
    ROW row1 {
      COLUMN col1 (DesktopWidth: AutoFill) {
        SNIPPETCALL snippetCall1 (Snippet: OdataPlm.Entity_Menu)
      }
      COLUMN col2 (DesktopWidth: AutoFill) {
        DYNAMICTEXT text1 (Content: 'Address', RenderMode: H2)
        DATAGRID dataGrid2_1 (DataSource: DATABASE OdataPlm.Address) {
          COLUMN col1 (Attribute: address_key, Caption: 'address key') {
            TEXTFILTER textFilter2
          }
          COLUMN col2 (Attribute: Country, Caption: 'Country') {
            TEXTFILTER textFilter3
          }
          COLUMN col3 (Attribute: State, Caption: 'State') {
            TEXTFILTER textFilter4
          }
          COLUMN col4 (Attribute: Street, Caption: 'Street') {
            TEXTFILTER textFilter5
          }
          COLUMN col5 (Attribute: Post_Code, Caption: 'Post Code') {
            TEXTFILTER textFilter6
          }
          COLUMN col6 (Attribute: City, Caption: 'City') {
            TEXTFILTER textFilter1
          }
          COLUMN col7 (Attribute: Zip_Code, Caption: 'Zip Code') {
            NUMBERFILTER numberFilter1 (FilterType: equal)
          }
          COLUMN col8 (Attribute: address_key, ShowContentAs: customContent) {
            ACTIONBUTTON actionButton1 (Action: SHOW_PAGE OdataPlm.Address_NewEdit, Style: Primary)
          }
        }
      }
    }
  }
}
```

#### Address Edit (Popup)

```sql
CREATE PAGE OdataPlm.Address_NewEdit
(Title: 'Edit Address', Layout: Atlas_Core.PopupLayout, Params: { $Address: OdataPlm.Address })
{
  LAYOUTGRID layoutGrid1 {
    ROW row1 {
      COLUMN col1 (DesktopWidth: AutoFill) {
        DATAVIEW dataView1 (DataSource: $Address) {
          TEXTBOX textBox1 (Label: 'address key', Attribute: OdataPlm.Address.address_key)
          TEXTBOX textBox2 (Label: 'Country', Attribute: OdataPlm.Address.Country)
          TEXTBOX textBox3 (Label: 'State', Attribute: OdataPlm.Address.State)
          TEXTBOX textBox4 (Label: 'Street', Attribute: OdataPlm.Address.Street)
          TEXTBOX textBox5 (Label: 'Post Code', Attribute: OdataPlm.Address.Post_Code)
          TEXTBOX textBox6 (Label: 'City', Attribute: OdataPlm.Address.City)
          TEXTBOX textBox7 (Label: 'Zip Code', Attribute: OdataPlm.Address.Zip_Code)
          FOOTER footer1 {
            ACTIONBUTTON actionButton1 (Caption: 'Save', Action: SAVE_CHANGES CLOSE_PAGE, Style: Success)
            ACTIONBUTTON actionButton2 (Caption: 'Cancel', Action: CANCEL_CHANGES CLOSE_PAGE)
          }
        }
      }
    }
  }
}
```

#### BOM Tree (Partial — TreeNode widget not configurable)

```sql
-- NOTE: The TREENODE widget cannot be fully configured via MDL.
-- This creates the page structure. Add the TreeNode widget configuration in Studio Pro.
CREATE PAGE OdataPlm.BomTree
(Title: 'Bom tree', Layout: Atlas_Core.Atlas_Default)
{
  LAYOUTGRID layoutGrid1 {
    ROW row1 {
      COLUMN col1 (DesktopWidth: AutoFill) {
        DYNAMICTEXT text1 (Content: 'BOM', RenderMode: H1)
      }
    }
    ROW row2 {
      COLUMN col1 (DesktopWidth: AutoFill) {
        -- TREENODE widget goes here (configure in Studio Pro)
        -- It displays Bom -> Component -> Sub_Component hierarchy
      }
    }
  }
}
```

### Step 10: Home Page (Main)


```sql
CREATE PAGE Main.Home_Web
(Title: 'Homepage', Layout: Atlas_Core.Atlas_TopBar)
{
  LAYOUTGRID layoutGrid3 {
    ROW row1 {
      COLUMN col1 (DesktopWidth: AutoFill) {
        DYNAMICTEXT text1 (Content: 'Tekst', RenderMode: H1)
      }
    }
    ROW row2 {
      COLUMN col1 (DesktopWidth: AutoFill) {
        ACTIONBUTTON actionButton3 (Caption: 'Customer Overview', Action: SHOW_PAGE OdataPlm.Customer_Overview)
        ACTIONBUTTON actionButton4 (Caption: 'Bom tree', Action: SHOW_PAGE OdataPlm.BomTree)
        DATAGRID dataGrid2_1 {
          COLUMN col1 (Attribute: Customer__type, Caption: 'Customer type') {
            TEXTFILTER textFilter1
          }
          COLUMN col2 (Attribute: Customer_Value, Caption: 'Customer Value') {
            TEXTFILTER textFilter2
          }
          COLUMN col3 (Attribute: CustomerId__type, Caption: 'Customer id type') {
            TEXTFILTER textFilter3
          }
          COLUMN col4 (Attribute: CustomerId_Value, Caption: 'Customer id Value') {
            TEXTFILTER textFilter4
          }
          COLUMN col5 (Attribute: CustomerName__type, Caption: 'Customer name type') {
            TEXTFILTER textFilter5
          }
          COLUMN col6 (Attribute: CustomerName_Value, Caption: 'Customer name Value') {
            TEXTFILTER textFilter6
          }
          COLUMN col7 (Attribute: Customer__type, ShowContentAs: customContent) {
            ACTIONBUTTON actionButton1 (Action: SHOW_PAGE Main.Customer_View, Style: Primary)
            ACTIONBUTTON actionButton2 (Action: DELETE_OBJECT, Style: Primary)
          }
        }
      }
    }
  }
}
```

### Step 11: Navigation

```sql
CREATE OR REPLACE NAVIGATION Responsive
  HOME PAGE Main.Home_Web
  MENU (
    MENU ITEM 'Home' PAGE Main.Home_Web;
  )
;
```

### Step 12: Microflow (Partial)

```sql
-- NOTE: The RestOperationCallAction is not supported in MDL.
-- This microflow must be completed in Studio Pro by adding the REST call action.
CREATE MICROFLOW Main.GetCustomers ()
RETURNS List of Main.Customer AS $CustomerList
BEGIN
  -- TODO: Add REST operation call to get customerResponse (requires Studio Pro)
  RETRIEVE $BindingList FROM ASSOCIATION $customerResponse/Main.Customer_GetCustomerResponse;
  RETURN $BindingList;
END;
/
```

## Summary

### What MDL Can Do (~95% of this app)

- Create consumed OData clients with authentication
- Create all external entities linked to the OData client
- Create all domain model entities (external, persistent, non-persistent)
- Create all associations
- Create all overview pages with DataGrid2, text filters, number filters
- Create all edit popup pages with DataView, TextBox, ComboBox, action buttons
- Create the sidebar navigation snippet with NavigationList
- Create the home page with layout grid, buttons, and data grid
- Set up navigation profiles with home pages and menus
- Define constants
- Define enumerations

### What Requires Studio Pro (~5%)

| Gap | Impact | Notes |
|-----|--------|-------|
| **TreeNode widget** | Medium — only affects the BOM tree page | Page structure can be created in MDL; widget needs Studio Pro |
| **RestOperationCallAction** | Low — only used in 1 microflow | The Main.GetCustomers microflow prototype; main app uses OData not REST |

### Recommended Workflow

1. Create the Mendix project in Studio Pro
2. Install marketplace modules (Atlas_Core, Administration, DataWidgets)
3. **With MDL**: Run the scripts above (Steps 2-12) to create the entire app
4. **In Studio Pro**: Add TreeNode widget configuration to BomTree page
5. **In Studio Pro**: Complete the GetCustomers microflow with REST call action (if needed)
