// SPDX-License-Identifier: Apache-2.0

package catalog

// createTables creates all catalog tables in the SQLite database.
func (c *Catalog) createTables() error {
	schemas := []string{
		// Catalog metadata table - for cache validation
		`CREATE TABLE IF NOT EXISTS catalog_meta (
			Key TEXT PRIMARY KEY,
			Value TEXT
		)`,

		// Projects table
		`CREATE TABLE IF NOT EXISTS projects (
			ProjectId TEXT PRIMARY KEY,
			ProjectName TEXT,
			MendixVersion TEXT,
			CreatedDate TEXT,
			LastSnapshotDate TEXT,
			SnapshotCount INTEGER DEFAULT 0
		)`,

		// Snapshots table
		`CREATE TABLE IF NOT EXISTS snapshots (
			SnapshotId TEXT PRIMARY KEY,
			SnapshotName TEXT,
			ProjectId TEXT,
			ProjectName TEXT,
			SnapshotDate TEXT,
			SnapshotSource TEXT,
			SourceId TEXT,
			SourceBranch TEXT,
			SourceRevision TEXT,
			ObjectCount INTEGER DEFAULT 0,
			IsActive INTEGER DEFAULT 0
		)`,

		// Modules table
		`CREATE TABLE IF NOT EXISTS modules (
			Id TEXT PRIMARY KEY,
			Name TEXT,
			QualifiedName TEXT,
			ModuleName TEXT,
			Folder TEXT,
			Description TEXT,
			IsSystemModule INTEGER DEFAULT 0,
			AppStoreVersion TEXT,
			AppStoreGuid TEXT,
			ProjectId TEXT,
			ProjectName TEXT,
			SnapshotId TEXT,
			SnapshotDate TEXT,
			SnapshotSource TEXT,
			SourceId TEXT,
			SourceBranch TEXT,
			SourceRevision TEXT
		)`,

		// Entities table
		`CREATE TABLE IF NOT EXISTS entities (
			Id TEXT PRIMARY KEY,
			Name TEXT,
			QualifiedName TEXT,
			ModuleName TEXT,
			Folder TEXT,
			EntityType TEXT,
			Description TEXT,
			Generalization TEXT,
			AttributeCount INTEGER DEFAULT 0,
			AssociationCount INTEGER DEFAULT 0,
			AccessRuleCount INTEGER DEFAULT 0,
			ValidationRuleCount INTEGER DEFAULT 0,
			HasEventHandlers INTEGER DEFAULT 0,
			IsExternal INTEGER DEFAULT 0,
			ExternalService TEXT,
			ProjectId TEXT,
			ProjectName TEXT,
			SnapshotId TEXT,
			SnapshotDate TEXT,
			SnapshotSource TEXT,
			SourceId TEXT,
			SourceBranch TEXT,
			SourceRevision TEXT
		)`,

		// Attributes table - stores entity attribute details
		`CREATE TABLE IF NOT EXISTS attributes (
			Id TEXT PRIMARY KEY,
			Name TEXT,
			EntityId TEXT,
			EntityQualifiedName TEXT,
			ModuleName TEXT,
			DataType TEXT,
			Length INTEGER,
			IsUnique INTEGER DEFAULT 0,
			IsRequired INTEGER DEFAULT 0,
			DefaultValue TEXT,
			IsCalculated INTEGER DEFAULT 0,
			Description TEXT,
			ProjectId TEXT,
			ProjectName TEXT,
			SnapshotId TEXT,
			SnapshotDate TEXT,
			SnapshotSource TEXT,
			SourceId TEXT,
			SourceBranch TEXT,
			SourceRevision TEXT
		)`,

		// Microflows table
		`CREATE TABLE IF NOT EXISTS microflows (
			Id TEXT PRIMARY KEY,
			Name TEXT,
			QualifiedName TEXT,
			ModuleName TEXT,
			Folder TEXT,
			MicroflowType TEXT,
			Description TEXT,
			ReturnType TEXT,
			ParameterCount INTEGER DEFAULT 0,
			ActivityCount INTEGER DEFAULT 0,
			Complexity INTEGER DEFAULT 1,
			ProjectId TEXT,
			ProjectName TEXT,
			SnapshotId TEXT,
			SnapshotDate TEXT,
			SnapshotSource TEXT,
			SourceId TEXT,
			SourceBranch TEXT,
			SourceRevision TEXT
		)`,

		// Nanoflows view (same structure, filtered)
		`CREATE VIEW IF NOT EXISTS nanoflows AS
			SELECT * FROM microflows WHERE MicroflowType = 'NANOFLOW'`,

		// Pages table
		`CREATE TABLE IF NOT EXISTS pages (
			Id TEXT PRIMARY KEY,
			Name TEXT,
			QualifiedName TEXT,
			ModuleName TEXT,
			Folder TEXT,
			Title TEXT,
			URL TEXT,
			LayoutRef TEXT,
			Description TEXT,
			ParameterCount INTEGER DEFAULT 0,
			WidgetCount INTEGER DEFAULT 0,
			ProjectId TEXT,
			ProjectName TEXT,
			SnapshotId TEXT,
			SnapshotDate TEXT,
			SnapshotSource TEXT,
			SourceId TEXT,
			SourceBranch TEXT,
			SourceRevision TEXT
		)`,

		// Snippets table
		`CREATE TABLE IF NOT EXISTS snippets (
			Id TEXT PRIMARY KEY,
			Name TEXT,
			QualifiedName TEXT,
			ModuleName TEXT,
			Folder TEXT,
			Description TEXT,
			ParameterCount INTEGER DEFAULT 0,
			WidgetCount INTEGER DEFAULT 0,
			ProjectId TEXT,
			ProjectName TEXT,
			SnapshotId TEXT,
			SnapshotDate TEXT,
			SnapshotSource TEXT,
			SourceId TEXT,
			SourceBranch TEXT,
			SourceRevision TEXT
		)`,

		// Layouts table
		`CREATE TABLE IF NOT EXISTS layouts (
			Id TEXT PRIMARY KEY,
			Name TEXT,
			QualifiedName TEXT,
			ModuleName TEXT,
			Folder TEXT,
			LayoutType TEXT,
			Description TEXT,
			ProjectId TEXT,
			ProjectName TEXT,
			SnapshotId TEXT,
			SnapshotDate TEXT,
			SnapshotSource TEXT,
			SourceId TEXT,
			SourceBranch TEXT,
			SourceRevision TEXT
		)`,

		// Enumerations table
		`CREATE TABLE IF NOT EXISTS enumerations (
			Id TEXT PRIMARY KEY,
			Name TEXT,
			QualifiedName TEXT,
			ModuleName TEXT,
			Folder TEXT,
			Description TEXT,
			ValueCount INTEGER DEFAULT 0,
			ProjectId TEXT,
			ProjectName TEXT,
			SnapshotId TEXT,
			SnapshotDate TEXT,
			SnapshotSource TEXT,
			SourceId TEXT,
			SourceBranch TEXT,
			SourceRevision TEXT
		)`,

		// Java Actions table
		`CREATE TABLE IF NOT EXISTS java_actions (
			Id TEXT PRIMARY KEY,
			Name TEXT,
			QualifiedName TEXT,
			ModuleName TEXT,
			Folder TEXT,
			Documentation TEXT,
			ExportLevel TEXT,
			ReturnType TEXT,
			ParameterCount INTEGER DEFAULT 0,
			ProjectId TEXT,
			ProjectName TEXT,
			SnapshotId TEXT,
			SnapshotDate TEXT,
			SnapshotSource TEXT,
			SourceId TEXT,
			SourceBranch TEXT,
			SourceRevision TEXT
		)`,

		// Activities table
		`CREATE TABLE IF NOT EXISTS activities (
			Id TEXT PRIMARY KEY,
			Name TEXT,
			Caption TEXT,
			ActivityType TEXT,
			Sequence INTEGER DEFAULT 0,
			MicroflowId TEXT,
			MicroflowQualifiedName TEXT,
			ModuleName TEXT,
			Folder TEXT,
			EntityRef TEXT,
			ActionType TEXT,
			Description TEXT,
			ProjectId TEXT,
			ProjectName TEXT,
			SnapshotId TEXT,
			SnapshotDate TEXT,
			SnapshotSource TEXT,
			SourceId TEXT,
			SourceBranch TEXT,
			SourceRevision TEXT
		)`,

		// Widgets table
		`CREATE TABLE IF NOT EXISTS widgets (
			Id TEXT PRIMARY KEY,
			Name TEXT,
			WidgetType TEXT,
			ContainerId TEXT,
			ContainerQualifiedName TEXT,
			ContainerType TEXT,
			ModuleName TEXT,
			Folder TEXT,
			EntityRef TEXT,
			AttributeRef TEXT,
			Description TEXT,
			ProjectId TEXT,
			ProjectName TEXT,
			SnapshotId TEXT,
			SnapshotDate TEXT,
			SnapshotSource TEXT,
			SourceId TEXT,
			SourceBranch TEXT,
			SourceRevision TEXT
		)`,

		// XPath expressions table
		`CREATE TABLE IF NOT EXISTS xpath_expressions (
			Id TEXT PRIMARY KEY,
			DocumentType TEXT,
			DocumentId TEXT,
			DocumentQualifiedName TEXT,
			ComponentType TEXT,
			ComponentId TEXT,
			ComponentName TEXT,
			XPathExpression TEXT,
			XPathAST TEXT,
			TargetEntity TEXT,
			ReferencedEntities TEXT,
			IsParameterized INTEGER DEFAULT 0,
			UsageType TEXT,
			ModuleName TEXT,
			Folder TEXT,
			Description TEXT,
			ProjectId TEXT,
			ProjectName TEXT,
			SnapshotId TEXT,
			SnapshotDate TEXT,
			SnapshotSource TEXT,
			SourceId TEXT,
			SourceBranch TEXT,
			SourceRevision TEXT
		)`,

		// OData Clients table (consumed OData services)
		`CREATE TABLE IF NOT EXISTS odata_clients (
			Id TEXT PRIMARY KEY,
			Name TEXT,
			QualifiedName TEXT,
			ModuleName TEXT,
			Version TEXT,
			ODataVersion TEXT,
			MetadataUrl TEXT,
			Validated INTEGER DEFAULT 0,
			ProjectId TEXT,
			ProjectName TEXT,
			SnapshotId TEXT,
			SnapshotDate TEXT,
			SnapshotSource TEXT,
			SourceId TEXT,
			SourceBranch TEXT,
			SourceRevision TEXT
		)`,

		// OData Services table (published OData services)
		`CREATE TABLE IF NOT EXISTS odata_services (
			Id TEXT PRIMARY KEY,
			Name TEXT,
			QualifiedName TEXT,
			ModuleName TEXT,
			Path TEXT,
			Version TEXT,
			ODataVersion TEXT,
			EntitySetCount INTEGER DEFAULT 0,
			AuthenticationTypes TEXT,
			ProjectId TEXT,
			ProjectName TEXT,
			SnapshotId TEXT,
			SnapshotDate TEXT,
			SnapshotSource TEXT,
			SourceId TEXT,
			SourceBranch TEXT,
			SourceRevision TEXT
		)`,

		// Workflows table
		`CREATE TABLE IF NOT EXISTS workflows (
			Id TEXT PRIMARY KEY,
			Name TEXT,
			QualifiedName TEXT,
			ModuleName TEXT,
			Folder TEXT,
			Description TEXT,
			ExportLevel TEXT,
			ParameterEntity TEXT,
			ActivityCount INTEGER DEFAULT 0,
			UserTaskCount INTEGER DEFAULT 0,
			MicroflowCallCount INTEGER DEFAULT 0,
			DecisionCount INTEGER DEFAULT 0,
			DueDate TEXT,
			ProjectId TEXT,
			ProjectName TEXT,
			SnapshotId TEXT,
			SnapshotDate TEXT,
			SnapshotSource TEXT,
			SourceId TEXT,
			SourceBranch TEXT,
			SourceRevision TEXT
		)`,

		// Business Event Services table
		`CREATE TABLE IF NOT EXISTS business_event_services (
			Id TEXT PRIMARY KEY,
			Name TEXT,
			QualifiedName TEXT,
			ModuleName TEXT,
			Documentation TEXT,
			ServiceName TEXT,
			EventNamePrefix TEXT,
			MessageCount INTEGER DEFAULT 0,
			PublishCount INTEGER DEFAULT 0,
			SubscribeCount INTEGER DEFAULT 0,
			ProjectId TEXT,
			ProjectName TEXT,
			SnapshotId TEXT,
			SnapshotDate TEXT,
			SnapshotSource TEXT,
			SourceId TEXT,
			SourceBranch TEXT,
			SourceRevision TEXT
		)`,

		// Navigation profiles table
		`CREATE TABLE IF NOT EXISTS navigation_profiles (
			ProfileName TEXT PRIMARY KEY,
			Kind TEXT,
			IsNative INTEGER DEFAULT 0,
			HomePage TEXT,
			HomePageType TEXT,
			LoginPage TEXT,
			NotFoundPage TEXT,
			MenuItemCount INTEGER DEFAULT 0,
			RoleBasedHomeCount INTEGER DEFAULT 0,
			OfflineEntityCount INTEGER DEFAULT 0,
			ProjectId TEXT,
			ProjectName TEXT,
			SnapshotId TEXT,
			SnapshotDate TEXT,
			SnapshotSource TEXT,
			SourceId TEXT,
			SourceBranch TEXT,
			SourceRevision TEXT
		)`,

		// Navigation menu items table
		`CREATE TABLE IF NOT EXISTS navigation_menu_items (
			Id INTEGER PRIMARY KEY AUTOINCREMENT,
			ProfileName TEXT NOT NULL,
			ItemPath TEXT NOT NULL,
			Depth INTEGER DEFAULT 0,
			Caption TEXT,
			ActionType TEXT,
			TargetPage TEXT,
			TargetMicroflow TEXT,
			SubItemCount INTEGER DEFAULT 0,
			ProjectId TEXT,
			SnapshotId TEXT
		)`,

		// Navigation role-based home pages table
		`CREATE TABLE IF NOT EXISTS navigation_role_homes (
			Id INTEGER PRIMARY KEY AUTOINCREMENT,
			ProfileName TEXT NOT NULL,
			UserRole TEXT NOT NULL,
			Page TEXT,
			Microflow TEXT,
			ProjectId TEXT,
			SnapshotId TEXT
		)`,

		`CREATE TABLE IF NOT EXISTS database_connections (
			Id TEXT PRIMARY KEY,
			Name TEXT,
			QualifiedName TEXT,
			ModuleName TEXT,
			Folder TEXT,
			DatabaseType TEXT,
			QueryCount INTEGER DEFAULT 0,
			ProjectId TEXT,
			ProjectName TEXT,
			SnapshotId TEXT,
			SnapshotDate TEXT,
			SnapshotSource TEXT,
			SourceId TEXT,
			SourceBranch TEXT,
			SourceRevision TEXT
		)`,

		// Role mappings table - user role to module role assignments
		`CREATE TABLE IF NOT EXISTS role_mappings (
			Id INTEGER PRIMARY KEY AUTOINCREMENT,
			UserRoleName TEXT NOT NULL,
			ModuleRoleName TEXT NOT NULL,
			ModuleName TEXT,
			ProjectId TEXT,
			SnapshotId TEXT
		)`,

		// References table - cross-references between objects
		// Enables queries like "find all callers of microflow X" or "find all usages of entity Y"
		`CREATE TABLE IF NOT EXISTS refs (
			Id INTEGER PRIMARY KEY AUTOINCREMENT,
			SourceType TEXT NOT NULL,
			SourceId TEXT NOT NULL,
			SourceName TEXT NOT NULL,
			TargetType TEXT NOT NULL,
			TargetId TEXT,
			TargetName TEXT NOT NULL,
			RefKind TEXT NOT NULL,
			ModuleName TEXT,
			ProjectId TEXT,
			SnapshotId TEXT
		)`,

		// Permissions table - queryable security permission matrix
		`CREATE TABLE IF NOT EXISTS permissions (
			Id INTEGER PRIMARY KEY AUTOINCREMENT,
			ModuleRoleName TEXT NOT NULL,
			ElementType TEXT NOT NULL,
			ElementName TEXT NOT NULL,
			MemberName TEXT,
			AccessType TEXT NOT NULL,
			XPathConstraint TEXT,
			ModuleName TEXT,
			ProjectId TEXT,
			SnapshotId TEXT
		)`,

		// Objects view - union of all object types
		`CREATE VIEW IF NOT EXISTS objects AS
			SELECT Id, 'MODULE' as ObjectType, Name, QualifiedName, '' as ModuleName, Folder, Description,
				ProjectId, ProjectName, SnapshotId, SnapshotDate, SnapshotSource
			FROM modules
			UNION ALL
			SELECT Id, 'ENTITY' as ObjectType, Name, QualifiedName, ModuleName, Folder, Description,
				ProjectId, ProjectName, SnapshotId, SnapshotDate, SnapshotSource
			FROM entities
			UNION ALL
			SELECT Id, 'MICROFLOW' as ObjectType, Name, QualifiedName, ModuleName, Folder, Description,
				ProjectId, ProjectName, SnapshotId, SnapshotDate, SnapshotSource
			FROM microflows WHERE MicroflowType = 'MICROFLOW'
			UNION ALL
			SELECT Id, 'NANOFLOW' as ObjectType, Name, QualifiedName, ModuleName, Folder, Description,
				ProjectId, ProjectName, SnapshotId, SnapshotDate, SnapshotSource
			FROM microflows WHERE MicroflowType = 'NANOFLOW'
			UNION ALL
			SELECT Id, 'PAGE' as ObjectType, Name, QualifiedName, ModuleName, Folder, Description,
				ProjectId, ProjectName, SnapshotId, SnapshotDate, SnapshotSource
			FROM pages
			UNION ALL
			SELECT Id, 'SNIPPET' as ObjectType, Name, QualifiedName, ModuleName, Folder, Description,
				ProjectId, ProjectName, SnapshotId, SnapshotDate, SnapshotSource
			FROM snippets
			UNION ALL
			SELECT Id, 'LAYOUT' as ObjectType, Name, QualifiedName, ModuleName, Folder, Description,
				ProjectId, ProjectName, SnapshotId, SnapshotDate, SnapshotSource
			FROM layouts
			UNION ALL
			SELECT Id, 'ENUMERATION' as ObjectType, Name, QualifiedName, ModuleName, Folder, Description,
				ProjectId, ProjectName, SnapshotId, SnapshotDate, SnapshotSource
			FROM enumerations
			UNION ALL
			SELECT Id, 'JAVA_ACTION' as ObjectType, Name, QualifiedName, ModuleName, Folder, Documentation as Description,
				ProjectId, ProjectName, SnapshotId, SnapshotDate, SnapshotSource
			FROM java_actions
			UNION ALL
			SELECT Id, 'ODATA_CLIENT' as ObjectType, Name, QualifiedName, ModuleName, '' as Folder, '' as Description,
				ProjectId, ProjectName, SnapshotId, SnapshotDate, SnapshotSource
			FROM odata_clients
			UNION ALL
			SELECT Id, 'ODATA_SERVICE' as ObjectType, Name, QualifiedName, ModuleName, '' as Folder, '' as Description,
				ProjectId, ProjectName, SnapshotId, SnapshotDate, SnapshotSource
			FROM odata_services
			UNION ALL
			SELECT Id, 'WORKFLOW' as ObjectType, Name, QualifiedName, ModuleName, Folder, Description,
				ProjectId, ProjectName, SnapshotId, SnapshotDate, SnapshotSource
			FROM workflows
			UNION ALL
			SELECT Id, 'BUSINESS_EVENT_SERVICE' as ObjectType, Name, QualifiedName, ModuleName, '' as Folder, Documentation as Description,
				ProjectId, ProjectName, SnapshotId, SnapshotDate, SnapshotSource
			FROM business_event_services
			UNION ALL
			SELECT Id, 'DATABASE_CONNECTION' as ObjectType, Name, QualifiedName, ModuleName, Folder, '' as Description,
				ProjectId, ProjectName, SnapshotId, SnapshotDate, SnapshotSource
			FROM database_connections`,

		// FTS5 virtual tables for full-text search
		`CREATE VIRTUAL TABLE IF NOT EXISTS strings USING fts5(
			QualifiedName,
			ObjectType,
			StringValue,
			StringContext,
			ModuleName
		)`,

		`CREATE VIRTUAL TABLE IF NOT EXISTS source USING fts5(
			QualifiedName,
			ObjectType,
			SourceText,
			ModuleName
		)`,

		// Indexes for common queries
		`CREATE INDEX IF NOT EXISTS idx_modules_name ON modules(Name)`,
		`CREATE INDEX IF NOT EXISTS idx_entities_name ON entities(Name)`,
		`CREATE INDEX IF NOT EXISTS idx_entities_module ON entities(ModuleName)`,
		`CREATE INDEX IF NOT EXISTS idx_microflows_name ON microflows(Name)`,
		`CREATE INDEX IF NOT EXISTS idx_microflows_module ON microflows(ModuleName)`,
		`CREATE INDEX IF NOT EXISTS idx_pages_name ON pages(Name)`,
		`CREATE INDEX IF NOT EXISTS idx_pages_module ON pages(ModuleName)`,
		`CREATE INDEX IF NOT EXISTS idx_layouts_name ON layouts(Name)`,
		`CREATE INDEX IF NOT EXISTS idx_layouts_module ON layouts(ModuleName)`,
		`CREATE INDEX IF NOT EXISTS idx_activities_microflow ON activities(MicroflowId)`,
		`CREATE INDEX IF NOT EXISTS idx_activities_type ON activities(ActivityType)`,
		`CREATE INDEX IF NOT EXISTS idx_widgets_container ON widgets(ContainerId)`,
		`CREATE INDEX IF NOT EXISTS idx_widgets_type ON widgets(WidgetType)`,
		`CREATE INDEX IF NOT EXISTS idx_xpath_document ON xpath_expressions(DocumentId)`,
		`CREATE INDEX IF NOT EXISTS idx_refs_source ON refs(SourceType, SourceName)`,
		`CREATE INDEX IF NOT EXISTS idx_refs_target ON refs(TargetType, TargetName)`,
		`CREATE INDEX IF NOT EXISTS idx_refs_kind ON refs(RefKind)`,
		`CREATE INDEX IF NOT EXISTS idx_attributes_entity ON attributes(EntityId)`,
		`CREATE INDEX IF NOT EXISTS idx_attributes_entity_qname ON attributes(EntityQualifiedName)`,
		`CREATE INDEX IF NOT EXISTS idx_java_actions_name ON java_actions(Name)`,
		`CREATE INDEX IF NOT EXISTS idx_java_actions_module ON java_actions(ModuleName)`,
		`CREATE INDEX IF NOT EXISTS idx_odata_clients_name ON odata_clients(Name)`,
		`CREATE INDEX IF NOT EXISTS idx_odata_clients_module ON odata_clients(ModuleName)`,
		`CREATE INDEX IF NOT EXISTS idx_odata_services_name ON odata_services(Name)`,
		`CREATE INDEX IF NOT EXISTS idx_odata_services_module ON odata_services(ModuleName)`,
		`CREATE INDEX IF NOT EXISTS idx_workflows_name ON workflows(QualifiedName)`,
		`CREATE INDEX IF NOT EXISTS idx_workflows_module ON workflows(ModuleName)`,
		`CREATE INDEX IF NOT EXISTS idx_be_services_name ON business_event_services(QualifiedName)`,
		`CREATE INDEX IF NOT EXISTS idx_be_services_module ON business_event_services(ModuleName)`,
		`CREATE INDEX IF NOT EXISTS idx_role_mappings_user_role ON role_mappings(UserRoleName)`,
		`CREATE INDEX IF NOT EXISTS idx_role_mappings_module_role ON role_mappings(ModuleRoleName)`,
		`CREATE INDEX IF NOT EXISTS idx_role_mappings_module ON role_mappings(ModuleName)`,
		`CREATE INDEX IF NOT EXISTS idx_permissions_role ON permissions(ModuleRoleName)`,
		`CREATE INDEX IF NOT EXISTS idx_permissions_element ON permissions(ElementType, ElementName)`,
		`CREATE INDEX IF NOT EXISTS idx_permissions_access ON permissions(AccessType)`,
		`CREATE INDEX IF NOT EXISTS idx_permissions_module ON permissions(ModuleName)`,
		`CREATE INDEX IF NOT EXISTS idx_nav_menu_items_profile ON navigation_menu_items(ProfileName)`,
		`CREATE INDEX IF NOT EXISTS idx_nav_menu_items_target_page ON navigation_menu_items(TargetPage)`,
		`CREATE INDEX IF NOT EXISTS idx_nav_menu_items_target_mf ON navigation_menu_items(TargetMicroflow)`,
		`CREATE INDEX IF NOT EXISTS idx_nav_role_homes_profile ON navigation_role_homes(ProfileName)`,
		`CREATE INDEX IF NOT EXISTS idx_nav_role_homes_role ON navigation_role_homes(UserRole)`,
	}

	for _, schema := range schemas {
		if _, err := c.db.Exec(schema); err != nil {
			return err
		}
	}

	return nil
}
