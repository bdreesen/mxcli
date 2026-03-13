---
id: APP-001
category: App/Crud
tags: [entity, crud, pages, navigation]
timeout: 10m
---

# APP-001: Bookstore Inventory

## Prompt
Create an app to manage my bookstore inventory. I need to track books with title, author, ISBN, price, and stock quantity.

## Expected Outcome
Domain model with Book entity, CRUD pages (overview, detail, edit), navigation, and basic microflows for create/update/delete.

## Checks
- entity_exists: "*.Book"
- entity_has_attribute: "*.Book.Title String"
- entity_has_attribute: "*.Book.Author String"
- entity_has_attribute: "*.Book.ISBN String"
- entity_has_attribute: "*.Book.Price Decimal"
- entity_has_attribute: "*.Book.StockQuantity Integer"
- page_exists: "*Overview*"
- page_exists: "*Edit*"
- navigation_has_item: true
- mx_check_passes: true

## Acceptance Criteria
- Book entity has all specified attributes with appropriate types
- Overview page with data grid
- New/Edit page with form
- Delete confirmation
- Navigation menu item

## Iteration

### Prompt
Add a category field to the books, and let me filter the book list by category.

### Checks
- entity_has_attribute: "*.Book.Category"

### Acceptance Criteria
- Category attribute added to Book entity
- Book list can be filtered by category
