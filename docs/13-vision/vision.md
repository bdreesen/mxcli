---
title: Project Vision
status: complete
updated: 2024-11-09
category: project
tags: [vision, goals, strategy]
---

# Project Vision: Mendix Model Context Protocol (MCP-MX)

## Executive Summary

The Mendix Model Context Protocol (MCP-MX) project provides a bridge between Large Language Models (LLMs) and the Mendix platform through a standardised, text-based interface. The project enables AI-assisted development, migration, and maintenance of Mendix applications using the Mendix Model Definition Language (MDL).

## Current State

### Architecture Overview

The MCP-MX architecture separates concerns between protocol translation (MCP server) and model manipulation (REPL). The REPL serves as the central hub for all Mendix model operations, whilst the LLM uses specialised MCP services for external integrations.

```mermaid
graph TB
    subgraph LLMLayer ["LLM Integration Layer"]
        LLM["Large Language Model<br/>Claude/GPT/Copilot"]
    end
    
    subgraph MCPServices ["MCP Service Ecosystem"]
        MCPMX["MCP-MX<br/>Protocol Translation"]
        MCPOQL["MCP-OQL<br/>Runtime Queries"]
        MCPSQL["MCP-SQL<br/>Database Access"]
        MCPSPARQL["MCP-SPARQL<br/>Graph Databases"]
        MCPCAT["MCP Catalog<br/>API Discovery"]
    end
    
    subgraph REPLCore ["REPL - Central Hub"]
        REPL["Mendix REPL<br/>MDL Parser & Executor"]
        
        subgraph Backends ["REPL Backends"]
            SDK["Model SDK<br/>Service Mode"]
            FS["Direct Filesystem<br/>Fast Mode"]
            SP["Studio Pro Extension<br/>Live Mode"]
        end
        
        subgraph Validation ["Validation Layer"]
            MXBUILD["mxbuild<br/>Server Mode"]
        end
    end
    
    subgraph MendixEco ["Mendix Ecosystem"]
        PROJ["Mendix Project<br/>Filesystem"]
        STUDIO["Studio Pro<br/>IDE"]
        RUNTIME["Mendix Runtime<br/>Running App"]
    end
    
    LLM -->|MCP Protocol| MCPMX
    LLM -.->|OQL Validation| MCPOQL
    LLM -.->|Database Queries| MCPSQL
    LLM -.->|Graph Queries| MCPSPARQL
    LLM -.->|API Discovery| MCPCAT
    
    MCPMX -->|MDL Commands| REPL
    MCPOQL -->|OQL Queries| RUNTIME
    
    REPL -->|SDK API| SDK
    REPL -->|File I/O| FS
    REPL -->|Extension API| SP
    REPL -->|Compile & Validate| MXBUILD
    
    SDK -->|Read/Write| PROJ
    FS -->|Direct Access| PROJ
    SP <-->|Live Editing| STUDIO
    STUDIO -->|Project Files| PROJ
    MXBUILD -->|Build Errors| REPL
    
    style LLM fill:#e1f5fe
    style MCPMX fill:#f3e5f5
    style MCPOQL fill:#e8eaf6
    style MCPSQL fill:#e8eaf6
    style MCPSPARQL fill:#e8eaf6
    style MCPCAT fill:#e8eaf6
    style REPL fill:#e8f5e8
    style SDK fill:#fff3e0
    style FS fill:#fff8e1
    style SP fill:#f1f8e9
    style MXBUILD fill:#fce4ec
    style PROJ fill:#ffebee
    style STUDIO fill:#f3e5f5
    style RUNTIME fill:#e0f2f1
```

### Key Components

1. **MCP-MX Server**: Protocol translation layer - converts MCP protocol to MDL commands
2. **Mendix REPL**: Central hub for model manipulation with pluggable backends
3. **REPL Backends**: 
   - **Model SDK Service**: Remote SDK for production use
   - **Direct Filesystem**: Fast local file manipulation
   - **Studio Pro Extension**: Live editing in open IDE
4. **Validation Layer**: mxbuild in server mode for compile-time validation
5. **External MCP Services**: Specialized services for runtime queries and 3rd party integrations

### Current Capabilities

- **Model Inspection**: LLMs can query domain models, microflows, and pages
- **Model Modification**: Create and modify entities, attributes, and associations
- **SQL-like Syntax**: Familiar query interface for model operations
- **Safe Operations**: Sandboxed environment for AI-generated changes

## Future Vision

### Enhanced Architecture with Validation and External Services

The enhanced architecture emphasizes the REPL as the central orchestrator with pluggable backends, integrated validation, and specialized MCP services for external integrations.

```mermaid
graph TB
    subgraph LLMWorkflow ["LLM-Driven Development Workflow"]
        LLM["Large Language Model<br/>AI Assistant"]
        
        subgraph MCPEcosystem ["MCP Service Ecosystem"]
            MCPMX["MCP-MX<br/>Model Operations"]
            MCPOQL["MCP-OQL<br/>Runtime Validation"]
            MCPSQL["MCP-SQL<br/>Database Integration"]
            MCPSPARQL["MCP-SPARQL<br/>Knowledge Graphs"]
            MCPCAT["MCP Catalog<br/>API Discovery"]
        end
    end
    
    subgraph REPLEngine ["REPL Engine - Central Orchestrator"]
        REPL["Mendix REPL<br/>MDL Parser & Executor"]
        
        subgraph BackendPlugins ["Backend Plugins"]
            SDK["Model SDK Service<br/>Production Mode"]
            FS["Direct Filesystem<br/>Fast Development"]
            SPEXT["Studio Pro Extension<br/>Live Collaboration"]
        end
        
        subgraph ValidationEngine ["Validation Engine"]
            MXBUILD["mxbuild Server<br/>Compile Validation"]
            PARSER["MDL Parser<br/>Syntax Validation"]
        end
    end
    
    subgraph MendixPlatform ["Mendix Platform"]
        PROJ["Project Files<br/>.mpr + Resources"]
        STUDIO["Studio Pro IDE<br/>Visual Development"]
        RUNTIME["Mendix Runtime<br/>Running Application"]
        CLI["CLI Tools<br/>mx, mxbuild"]
    end
    
    subgraph ExternalSystems ["External Systems"]
        DB[("Databases<br/>SQL/NoSQL")]
        GRAPH[("Knowledge Graphs<br/>RDF/SPARQL")]
        APIS["API Catalogs<br/>REST/GraphQL"]
        VCS["Version Control<br/>Git"]
    end
    
    %% LLM to MCP Services
    LLM -->|Model Changes| MCPMX
    LLM -.->|Query Validation| MCPOQL
    LLM -.->|Data Analysis| MCPSQL
    LLM -.->|Semantic Queries| MCPSPARQL
    LLM -.->|API Discovery| MCPCAT
    
    %% MCP-MX to REPL
    MCPMX -->|MDL Commands| REPL
    
    %% REPL to Backends
    REPL -->|SDK Calls| SDK
    REPL -->|File Operations| FS
    REPL -->|Extension API| SPEXT
    
    %% REPL Validation
    REPL -->|Compile Project| MXBUILD
    REPL -->|Parse MDL| PARSER
    MXBUILD -->|Build Errors| REPL
    PARSER -->|Syntax Errors| REPL
    
    %% Backends to Mendix
    SDK -->|Read/Write Model| PROJ
    FS -->|Direct File Access| PROJ
    SPEXT <-->|Live Sync| STUDIO
    STUDIO -->|Save Project| PROJ
    
    %% External MCP Services
    MCPOQL -->|OQL Queries| RUNTIME
    MCPSQL -->|SQL Queries| DB
    MCPSPARQL -->|SPARQL Queries| GRAPH
    MCPCAT -->|API Metadata| APIS
    
    %% Additional Integrations
    PROJ -.->|Version Control| VCS
    MXBUILD -.->|CLI Tools| CLI
    
    style LLM fill:#e1f5fe
    style MCPMX fill:#f3e5f5
    style MCPOQL fill:#e8eaf6
    style MCPSQL fill:#e8eaf6
    style MCPSPARQL fill:#e8eaf6
    style MCPCAT fill:#e8eaf6
    style REPL fill:#e8f5e8
    style SDK fill:#fff3e0
    style FS fill:#fff8e1
    style SPEXT fill:#f1f8e9
    style MXBUILD fill:#fce4ec
    style PARSER fill:#fce4ec
    style PROJ fill:#ffebee
    style STUDIO fill:#f3e5f5
    style RUNTIME fill:#e0f2f1
```

### Enhanced Capabilities Roadmap

#### 1. REPL Backend Plugins

The REPL supports multiple backend implementations, each optimized for different use cases:

**Model SDK Service Backend**
- Remote SDK service for production environments
- Full model API coverage
- Safe, transactional operations
- Currently implemented

**Direct Filesystem Backend**
- Significantly faster than SDK for bulk operations
- Direct XML manipulation for .mpr files
- Access to all project resources (images, stylesheets, etc.)
- Planned implementation

**Studio Pro Extension Backend**
- Live collaboration with developers in IDE
- Real-time visual feedback
- Seamless integration with existing workflows
- Future implementation

```mermaid
sequenceDiagram
    participant LLM
    participant MCP as MCP-MX
    participant REPL as Mendix REPL
    participant Backend as Backend Plugin
    participant Val as mxbuild Server
    
    LLM->>MCP: "Create entity Customer"
    MCP->>REPL: CREATE ENTITY Customer...
    REPL->>REPL: Parse MDL
    REPL->>Backend: Execute via plugin
    Backend->>Backend: Write to project
    Backend-->>REPL: Success
    REPL->>Val: Compile project
    Val-->>REPL: Build errors (if any)
    REPL-->>MCP: Result + validation
    MCP-->>LLM: "Entity created successfully"
```

#### 2. Integrated Validation with mxbuild

**Benefits**:
- Compile-time validation of all model changes
- Immediate feedback on errors and warnings
- Ensures model consistency and correctness

The REPL integrates with mxbuild running in server mode to provide continuous validation:

```mermaid
sequenceDiagram
    participant REPL as Mendix REPL
    participant MXB as mxbuild Server
    participant PROJ as Project Files
    
    Note over MXB: mxbuild starts in server mode
    
    loop For each MDL command
        REPL->>PROJ: Apply changes
        REPL->>MXB: Compile project
        MXB->>PROJ: Read model
        MXB->>MXB: Validate & compile
        alt Compilation successful
            MXB-->>REPL: Success
        else Compilation failed
            MXB-->>REPL: Build errors & warnings
            REPL->>REPL: Report errors to user
        end
    end
```

#### 3. External MCP Services Integration

**Benefits**:
- Specialized services for different integration needs
- LLM can orchestrate multiple services
- Separation of concerns for better maintainability

**MCP-OQL** (Runtime Query Validation)
- Execute OQL queries against running Mendix runtime
- Validate OQL syntax and semantics
- Test view entity definitions with real data
- Performance analysis and optimization

**MCP-SQL** (Database Integration)
- Query external databases for migration scenarios
- Analyze existing database schemas
- Generate Mendix domain models from database structures

**MCP-SPARQL** (Knowledge Graph Integration)
- Query semantic data sources
- Integrate with ontologies and taxonomies
- Support for cultural heritage and scientific applications

**MCP Catalog** (API Discovery)
- Discover available REST/GraphQL APIs
- Generate Mendix integration modules
- Automatic service consumption setup

```mermaid
graph LR
    subgraph MCPOrchestration ["LLM Orchestrates Multiple Services"]
        LLM["Large Language Model"]
        
        LLM -->|1. Analyze DB| MCPSQL["MCP-SQL<br/>Schema Analysis"]
        LLM -->|2. Generate Model| MCPMX["MCP-MX<br/>Create Entities"]
        LLM -->|3. Validate OQL| MCPOQL["MCP-OQL<br/>Test Queries"]
        LLM -->|4. Discover APIs| MCPCAT["MCP Catalog<br/>Find Services"]
        
        MCPSQL -->|Schema Info| LLM
        MCPMX -->|Model Created| LLM
        MCPOQL -->|Query Results| LLM
        MCPCAT -->|API Specs| LLM
    end
```

## Strategic Applications

### 1. Database-First Application Generation

The LLM orchestrates multiple MCP services to generate Mendix applications from existing databases:

```mermaid
sequenceDiagram
    participant LLM as Large Language Model
    participant SQL as MCP-SQL
    participant MX as MCP-MX
    participant REPL as Mendix REPL
    participant OQL as MCP-OQL
    participant DB as Legacy Database
    participant RT as Mendix Runtime
    
    LLM->>SQL: Analyze database schema
    SQL->>DB: Query metadata
    DB-->>SQL: Tables, columns, relationships
    SQL-->>LLM: Schema structure
    
    LLM->>LLM: Design domain model
    
    LLM->>MX: CREATE ENTITY commands
    MX->>REPL: Execute MDL
    REPL->>REPL: Create entities & associations
    REPL-->>MX: Model created
    MX-->>LLM: Success
    
    LLM->>MX: CREATE VIEW ENTITY for queries
    MX->>REPL: Execute MDL
    REPL-->>MX: Views created
    
    LLM->>OQL: Test view queries
    OQL->>RT: Execute OQL
    RT-->>OQL: Query results
    OQL-->>LLM: Validation successful
    
    LLM-->>LLM: Application ready
```

### 2. Version Migration Assistant

Leverage MDL versioning to create an intelligent migration tool for Mendix platform upgrades.

```mermaid
flowchart TD
    OLD[Mendix 9.x Project] -->|MDL Analysis| SCAN[Migration Scanner]
    SCAN -->|Compatibility Check| ISSUES[Issue Detection<br/>- Deprecated APIs<br/>- Breaking Changes<br/>- Performance Issues]
    ISSUES -->|AI-Powered Fixes| FIX[Automated Fixes<br/>- API Updates<br/>- Pattern Modernization<br/>- Performance Optimization]
    FIX -->|Validation| TEST[Automated Testing<br/>- Model Validation<br/>- Runtime Testing<br/>- Performance Benchmarks]
    TEST -->|Success| NEW[Mendix 10.x Project]
    TEST -->|Issues Found| ISSUES
```

### 3. Cross-Platform Migration Engine

Enable migration from Mendix to other platforms through MDL standardization.

```mermaid
graph TB
    subgraph PlatMig ["Platform Migration"]
        MX["Mendix Project<br/>.mpr Format"]
        MDL["MDL Representation<br/>Platform-Agnostic"]
        
        subgraph TargetPlatforms ["Target Platforms"]
            SPRING["Spring Boot<br/>+ React"]
            NET[".NET Core<br/>+ Angular"]
            DJANGO["Django<br/>+ Vue.js"]
            CUSTOM["Custom Platform<br/>Configurable Output"]
        end
        
        MX -->|Extract & Analyze| MDL
        MDL -->|Java Generator| SPRING
        MDL -->|C# Generator| NET
        MDL -->|Python Generator| DJANGO
        MDL -->|Template Engine| CUSTOM
    end
```

## Technology Evolution

### MDL 2.0 Specification

The next generation of Mendix Model Definition Language will include:

1. **Version Awareness**: Built-in support for Mendix platform versions
2. **Extended Coverage**: Support for all Mendix constructs (styling, workflows, etc.)
3. **Semantic Annotations**: AI-friendly metadata for better understanding
4. **Validation Rules**: Built-in constraints and business rules
5. **Dependency Tracking**: Automatic impact analysis for changes

### Performance Considerations

Different backend implementations offer varying performance characteristics:

| Operation Type | SDK Backend | Filesystem Backend |
|----------------|-------------|-------------------|
| Entity Creation | 2-5 seconds | Sub-second |
| Domain Model Query | 1-3 seconds | Sub-second |
| Bulk Operations | 30-120 seconds | 1-5 seconds |
| Project Analysis | 5-30 minutes | 30-180 seconds |

## Risk Analysis and Mitigation

### Technical Risks

1. **Platform Migration Risk**: 
   - **Risk**: MDL makes Mendix projects more portable to competitors
   - **Mitigation**: Focus on Mendix-specific optimizations and ecosystem integration
   - **Opportunity**: Position as industry standard for low-code AI integration

2. **Performance Complexity**:
   - **Risk**: Multiple execution modes increase maintenance overhead
   - **Mitigation**: Shared core with mode-specific adapters
   - **Strategy**: Gradual rollout with fallback to proven SDK mode

3. **IDE Integration Challenges**:
   - **Risk**: Studio Pro API limitations or breaking changes
   - **Mitigation**: Abstract interface layer with version compatibility
   - **Backup**: Continue with external project file manipulation

### Business Opportunities

1. **AI Development Platform**: Integration of AI capabilities into low-code development
2. **Enterprise Migration Services**: Automated legacy system modernisation
3. **Developer Productivity**: Significant improvements in development speed for common tasks
4. **Quality Assurance**: AI-assisted code review and optimisation

## Conclusion

MCP-MX provides a practical approach to integrating AI capabilities with low-code platforms. By bridging natural language AI and visual development tools, the project enables:

- **Developer Productivity**: AI-assisted development workflows
- **Platform Migration**: Automated legacy system modernisation  
- **Quality Assurance**: AI-assisted validation and optimisation
- **Rapid Development**: Faster prototyping and iteration

The modular architecture provides flexibility whilst maintaining the safety and reliability required for enterprise development. The project aims to establish practical patterns for AI integration in low-code development platforms.

---

*This document represents the strategic vision for MCP-MX. Technical specifications and detailed implementation plans are available in the accompanying architecture documents.*