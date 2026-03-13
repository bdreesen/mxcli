# REST Client Skill

Use this skill when working with Consumed REST Clients in MDL.

## Quick Reference

### Create a REST Client

```sql
CREATE REST CLIENT Module.ClientName
BASE URL 'https://api.example.com'
AUTHENTICATION NONE
BEGIN
  OPERATION GetData
    METHOD GET
    PATH '/data'
    RESPONSE JSON AS $Result;
END;
```

### Show/Describe REST Clients

```sql
-- List all REST clients
SHOW REST CLIENTS;

-- Show client source code
DESCRIBE REST CLIENT Module.ClientName;
```

## Key Components

### 1. Tokens (`src/parser/lexer/tokens.ts`)
- `Rest`, `Client`, `Clients`
- `Base`, `Url`
- `Authentication`, `Basic`, `None`
- `Operation`, `Method`, `Path`, `Timeout`
- `Header`, `Parameter`, `Query`, `Body`
- `Response`, `Json`, `Status`, `File`
- `Get`, `Post`, `Put`, `Patch`, `Delete`

### 2. Parser Rules (`src/parser/grammar/MendixParser.ts`)
- `createRestClient` - Main rule
- `restAuthentication` - BASIC or NONE
- `restOperation` - Operation definition
- `restHeader` - Header with value
- `restResponse` - Response handling

### 3. AST Types (`src/parser/ast/types.ts`)
- `CreateRestClientAST`
- `ShowRestClientsAST`
- `DescribeRestClientAST`
- `RestOperationDefinition`
- `RestAuthentication`

### 4. Handlers (`src/repl/handlers/rest-client-handlers.ts`)
- `ShowRestClientsHandler`
- `CreateRestClientHandler`
- `DescribeRestClientHandler`

### 5. Creator (`src/creators/RestClientCreator.ts`)
Creates `rest.ConsumedRestService` in the Model SDK.

### 6. Generator (`src/generators/RestClientGenerator.ts`)
Converts SDK objects back to MDL format.

## Testing

```bash
# Test parser
node /tmp/test-rest-parser.js

# Via REPL daemon
pnpm run repl:exec -- "SHOW REST CLIENTS;"
```

## SDK Classes Used

- `rest.ConsumedRestService` - Main client
- `rest.RestOperation` - Operations
- `rest.RestOperationMethod` / `RestOperationMethodWithBody` / `RestOperationMethodWithoutBody`
- `rest.ValueTemplate` - URL/path templates
- `rest.BasicAuthenticationScheme` - Basic auth
- `rest.StringValue` - Auth credentials
- `rest.HeaderWithValueTemplate` - Headers
- `services.HttpMethod` - GET/POST/PUT/PATCH/DELETE

## Common Patterns

### Add a new operation type
1. Add parser rule in `restOperation`
2. Add AST type property in `RestOperationDefinition`
3. Update visitor in `restOperation()`
4. Update creator in `RestClientCreator`
5. Update generator in `RestClientGenerator`

### Add authentication type
1. Add token if needed
2. Add alternative in `restAuthentication` rule
3. Add AST type (e.g., `RestOAuth2Auth`)
4. Update visitor
5. Update creator with SDK class

## Documentation

- Implementation: `/docs/02-features/REST_CLIENT_SYNTAX.md`
- Original Proposal: `/docs/06-future/rest-client-mdl-proposal.md`
