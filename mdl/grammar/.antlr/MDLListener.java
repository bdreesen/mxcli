// Generated from /workspaces/ModelSDKGo/pkg/mdl/grammar/MDL.g4 by ANTLR 4.13.1
import org.antlr.v4.runtime.tree.ParseTreeListener;

/**
 * This interface defines a complete listener for a parse tree produced by
 * {@link MDLParser}.
 */
public interface MDLListener extends ParseTreeListener {
	/**
	 * Enter a parse tree produced by {@link MDLParser#program}.
	 * @param ctx the parse tree
	 */
	void enterProgram(MDLParser.ProgramContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#program}.
	 * @param ctx the parse tree
	 */
	void exitProgram(MDLParser.ProgramContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#statement}.
	 * @param ctx the parse tree
	 */
	void enterStatement(MDLParser.StatementContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#statement}.
	 * @param ctx the parse tree
	 */
	void exitStatement(MDLParser.StatementContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#terminator}.
	 * @param ctx the parse tree
	 */
	void enterTerminator(MDLParser.TerminatorContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#terminator}.
	 * @param ctx the parse tree
	 */
	void exitTerminator(MDLParser.TerminatorContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#connectionStatement}.
	 * @param ctx the parse tree
	 */
	void enterConnectionStatement(MDLParser.ConnectionStatementContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#connectionStatement}.
	 * @param ctx the parse tree
	 */
	void exitConnectionStatement(MDLParser.ConnectionStatementContext ctx);
	/**
	 * Enter a parse tree produced by the {@code ConnectLocal}
	 * labeled alternative in {@link MDLParser#connectStatement}.
	 * @param ctx the parse tree
	 */
	void enterConnectLocal(MDLParser.ConnectLocalContext ctx);
	/**
	 * Exit a parse tree produced by the {@code ConnectLocal}
	 * labeled alternative in {@link MDLParser#connectStatement}.
	 * @param ctx the parse tree
	 */
	void exitConnectLocal(MDLParser.ConnectLocalContext ctx);
	/**
	 * Enter a parse tree produced by the {@code ConnectFilesystem}
	 * labeled alternative in {@link MDLParser#connectStatement}.
	 * @param ctx the parse tree
	 */
	void enterConnectFilesystem(MDLParser.ConnectFilesystemContext ctx);
	/**
	 * Exit a parse tree produced by the {@code ConnectFilesystem}
	 * labeled alternative in {@link MDLParser#connectStatement}.
	 * @param ctx the parse tree
	 */
	void exitConnectFilesystem(MDLParser.ConnectFilesystemContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#disconnectStatement}.
	 * @param ctx the parse tree
	 */
	void enterDisconnectStatement(MDLParser.DisconnectStatementContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#disconnectStatement}.
	 * @param ctx the parse tree
	 */
	void exitDisconnectStatement(MDLParser.DisconnectStatementContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#statusStatement}.
	 * @param ctx the parse tree
	 */
	void enterStatusStatement(MDLParser.StatusStatementContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#statusStatement}.
	 * @param ctx the parse tree
	 */
	void exitStatusStatement(MDLParser.StatusStatementContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#ddlStatement}.
	 * @param ctx the parse tree
	 */
	void enterDdlStatement(MDLParser.DdlStatementContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#ddlStatement}.
	 * @param ctx the parse tree
	 */
	void exitDdlStatement(MDLParser.DdlStatementContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#createStatement}.
	 * @param ctx the parse tree
	 */
	void enterCreateStatement(MDLParser.CreateStatementContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#createStatement}.
	 * @param ctx the parse tree
	 */
	void exitCreateStatement(MDLParser.CreateStatementContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#createModuleStatement}.
	 * @param ctx the parse tree
	 */
	void enterCreateModuleStatement(MDLParser.CreateModuleStatementContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#createModuleStatement}.
	 * @param ctx the parse tree
	 */
	void exitCreateModuleStatement(MDLParser.CreateModuleStatementContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#alterStatement}.
	 * @param ctx the parse tree
	 */
	void enterAlterStatement(MDLParser.AlterStatementContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#alterStatement}.
	 * @param ctx the parse tree
	 */
	void exitAlterStatement(MDLParser.AlterStatementContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#dropStatement}.
	 * @param ctx the parse tree
	 */
	void enterDropStatement(MDLParser.DropStatementContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#dropStatement}.
	 * @param ctx the parse tree
	 */
	void exitDropStatement(MDLParser.DropStatementContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#dropModuleStatement}.
	 * @param ctx the parse tree
	 */
	void enterDropModuleStatement(MDLParser.DropModuleStatementContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#dropModuleStatement}.
	 * @param ctx the parse tree
	 */
	void exitDropModuleStatement(MDLParser.DropModuleStatementContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#createEnumerationStatement}.
	 * @param ctx the parse tree
	 */
	void enterCreateEnumerationStatement(MDLParser.CreateEnumerationStatementContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#createEnumerationStatement}.
	 * @param ctx the parse tree
	 */
	void exitCreateEnumerationStatement(MDLParser.CreateEnumerationStatementContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#enumValueList}.
	 * @param ctx the parse tree
	 */
	void enterEnumValueList(MDLParser.EnumValueListContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#enumValueList}.
	 * @param ctx the parse tree
	 */
	void exitEnumValueList(MDLParser.EnumValueListContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#enumValue}.
	 * @param ctx the parse tree
	 */
	void enterEnumValue(MDLParser.EnumValueContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#enumValue}.
	 * @param ctx the parse tree
	 */
	void exitEnumValue(MDLParser.EnumValueContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#alterEnumerationStatement}.
	 * @param ctx the parse tree
	 */
	void enterAlterEnumerationStatement(MDLParser.AlterEnumerationStatementContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#alterEnumerationStatement}.
	 * @param ctx the parse tree
	 */
	void exitAlterEnumerationStatement(MDLParser.AlterEnumerationStatementContext ctx);
	/**
	 * Enter a parse tree produced by the {@code AddEnumValue}
	 * labeled alternative in {@link MDLParser#alterEnumOperation}.
	 * @param ctx the parse tree
	 */
	void enterAddEnumValue(MDLParser.AddEnumValueContext ctx);
	/**
	 * Exit a parse tree produced by the {@code AddEnumValue}
	 * labeled alternative in {@link MDLParser#alterEnumOperation}.
	 * @param ctx the parse tree
	 */
	void exitAddEnumValue(MDLParser.AddEnumValueContext ctx);
	/**
	 * Enter a parse tree produced by the {@code DropEnumValue}
	 * labeled alternative in {@link MDLParser#alterEnumOperation}.
	 * @param ctx the parse tree
	 */
	void enterDropEnumValue(MDLParser.DropEnumValueContext ctx);
	/**
	 * Exit a parse tree produced by the {@code DropEnumValue}
	 * labeled alternative in {@link MDLParser#alterEnumOperation}.
	 * @param ctx the parse tree
	 */
	void exitDropEnumValue(MDLParser.DropEnumValueContext ctx);
	/**
	 * Enter a parse tree produced by the {@code RenameEnumValue}
	 * labeled alternative in {@link MDLParser#alterEnumOperation}.
	 * @param ctx the parse tree
	 */
	void enterRenameEnumValue(MDLParser.RenameEnumValueContext ctx);
	/**
	 * Exit a parse tree produced by the {@code RenameEnumValue}
	 * labeled alternative in {@link MDLParser#alterEnumOperation}.
	 * @param ctx the parse tree
	 */
	void exitRenameEnumValue(MDLParser.RenameEnumValueContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#dropEnumerationStatement}.
	 * @param ctx the parse tree
	 */
	void enterDropEnumerationStatement(MDLParser.DropEnumerationStatementContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#dropEnumerationStatement}.
	 * @param ctx the parse tree
	 */
	void exitDropEnumerationStatement(MDLParser.DropEnumerationStatementContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#createEntityStatement}.
	 * @param ctx the parse tree
	 */
	void enterCreateEntityStatement(MDLParser.CreateEntityStatementContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#createEntityStatement}.
	 * @param ctx the parse tree
	 */
	void exitCreateEntityStatement(MDLParser.CreateEntityStatementContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#createOrModify}.
	 * @param ctx the parse tree
	 */
	void enterCreateOrModify(MDLParser.CreateOrModifyContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#createOrModify}.
	 * @param ctx the parse tree
	 */
	void exitCreateOrModify(MDLParser.CreateOrModifyContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#entityType}.
	 * @param ctx the parse tree
	 */
	void enterEntityType(MDLParser.EntityTypeContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#entityType}.
	 * @param ctx the parse tree
	 */
	void exitEntityType(MDLParser.EntityTypeContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#attributeList}.
	 * @param ctx the parse tree
	 */
	void enterAttributeList(MDLParser.AttributeListContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#attributeList}.
	 * @param ctx the parse tree
	 */
	void exitAttributeList(MDLParser.AttributeListContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#attribute}.
	 * @param ctx the parse tree
	 */
	void enterAttribute(MDLParser.AttributeContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#attribute}.
	 * @param ctx the parse tree
	 */
	void exitAttribute(MDLParser.AttributeContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#attributeName}.
	 * @param ctx the parse tree
	 */
	void enterAttributeName(MDLParser.AttributeNameContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#attributeName}.
	 * @param ctx the parse tree
	 */
	void exitAttributeName(MDLParser.AttributeNameContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#renamedFromAnnotation}.
	 * @param ctx the parse tree
	 */
	void enterRenamedFromAnnotation(MDLParser.RenamedFromAnnotationContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#renamedFromAnnotation}.
	 * @param ctx the parse tree
	 */
	void exitRenamedFromAnnotation(MDLParser.RenamedFromAnnotationContext ctx);
	/**
	 * Enter a parse tree produced by the {@code StringType}
	 * labeled alternative in {@link MDLParser#dataType}.
	 * @param ctx the parse tree
	 */
	void enterStringType(MDLParser.StringTypeContext ctx);
	/**
	 * Exit a parse tree produced by the {@code StringType}
	 * labeled alternative in {@link MDLParser#dataType}.
	 * @param ctx the parse tree
	 */
	void exitStringType(MDLParser.StringTypeContext ctx);
	/**
	 * Enter a parse tree produced by the {@code IntegerType}
	 * labeled alternative in {@link MDLParser#dataType}.
	 * @param ctx the parse tree
	 */
	void enterIntegerType(MDLParser.IntegerTypeContext ctx);
	/**
	 * Exit a parse tree produced by the {@code IntegerType}
	 * labeled alternative in {@link MDLParser#dataType}.
	 * @param ctx the parse tree
	 */
	void exitIntegerType(MDLParser.IntegerTypeContext ctx);
	/**
	 * Enter a parse tree produced by the {@code LongType}
	 * labeled alternative in {@link MDLParser#dataType}.
	 * @param ctx the parse tree
	 */
	void enterLongType(MDLParser.LongTypeContext ctx);
	/**
	 * Exit a parse tree produced by the {@code LongType}
	 * labeled alternative in {@link MDLParser#dataType}.
	 * @param ctx the parse tree
	 */
	void exitLongType(MDLParser.LongTypeContext ctx);
	/**
	 * Enter a parse tree produced by the {@code DecimalType}
	 * labeled alternative in {@link MDLParser#dataType}.
	 * @param ctx the parse tree
	 */
	void enterDecimalType(MDLParser.DecimalTypeContext ctx);
	/**
	 * Exit a parse tree produced by the {@code DecimalType}
	 * labeled alternative in {@link MDLParser#dataType}.
	 * @param ctx the parse tree
	 */
	void exitDecimalType(MDLParser.DecimalTypeContext ctx);
	/**
	 * Enter a parse tree produced by the {@code BooleanType}
	 * labeled alternative in {@link MDLParser#dataType}.
	 * @param ctx the parse tree
	 */
	void enterBooleanType(MDLParser.BooleanTypeContext ctx);
	/**
	 * Exit a parse tree produced by the {@code BooleanType}
	 * labeled alternative in {@link MDLParser#dataType}.
	 * @param ctx the parse tree
	 */
	void exitBooleanType(MDLParser.BooleanTypeContext ctx);
	/**
	 * Enter a parse tree produced by the {@code DateTimeType}
	 * labeled alternative in {@link MDLParser#dataType}.
	 * @param ctx the parse tree
	 */
	void enterDateTimeType(MDLParser.DateTimeTypeContext ctx);
	/**
	 * Exit a parse tree produced by the {@code DateTimeType}
	 * labeled alternative in {@link MDLParser#dataType}.
	 * @param ctx the parse tree
	 */
	void exitDateTimeType(MDLParser.DateTimeTypeContext ctx);
	/**
	 * Enter a parse tree produced by the {@code DateType}
	 * labeled alternative in {@link MDLParser#dataType}.
	 * @param ctx the parse tree
	 */
	void enterDateType(MDLParser.DateTypeContext ctx);
	/**
	 * Exit a parse tree produced by the {@code DateType}
	 * labeled alternative in {@link MDLParser#dataType}.
	 * @param ctx the parse tree
	 */
	void exitDateType(MDLParser.DateTypeContext ctx);
	/**
	 * Enter a parse tree produced by the {@code AutoNumberType}
	 * labeled alternative in {@link MDLParser#dataType}.
	 * @param ctx the parse tree
	 */
	void enterAutoNumberType(MDLParser.AutoNumberTypeContext ctx);
	/**
	 * Exit a parse tree produced by the {@code AutoNumberType}
	 * labeled alternative in {@link MDLParser#dataType}.
	 * @param ctx the parse tree
	 */
	void exitAutoNumberType(MDLParser.AutoNumberTypeContext ctx);
	/**
	 * Enter a parse tree produced by the {@code BinaryType}
	 * labeled alternative in {@link MDLParser#dataType}.
	 * @param ctx the parse tree
	 */
	void enterBinaryType(MDLParser.BinaryTypeContext ctx);
	/**
	 * Exit a parse tree produced by the {@code BinaryType}
	 * labeled alternative in {@link MDLParser#dataType}.
	 * @param ctx the parse tree
	 */
	void exitBinaryType(MDLParser.BinaryTypeContext ctx);
	/**
	 * Enter a parse tree produced by the {@code EnumerationType}
	 * labeled alternative in {@link MDLParser#dataType}.
	 * @param ctx the parse tree
	 */
	void enterEnumerationType(MDLParser.EnumerationTypeContext ctx);
	/**
	 * Exit a parse tree produced by the {@code EnumerationType}
	 * labeled alternative in {@link MDLParser#dataType}.
	 * @param ctx the parse tree
	 */
	void exitEnumerationType(MDLParser.EnumerationTypeContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#attributeConstraints}.
	 * @param ctx the parse tree
	 */
	void enterAttributeConstraints(MDLParser.AttributeConstraintsContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#attributeConstraints}.
	 * @param ctx the parse tree
	 */
	void exitAttributeConstraints(MDLParser.AttributeConstraintsContext ctx);
	/**
	 * Enter a parse tree produced by the {@code NotNullConstraint}
	 * labeled alternative in {@link MDLParser#attributeConstraint}.
	 * @param ctx the parse tree
	 */
	void enterNotNullConstraint(MDLParser.NotNullConstraintContext ctx);
	/**
	 * Exit a parse tree produced by the {@code NotNullConstraint}
	 * labeled alternative in {@link MDLParser#attributeConstraint}.
	 * @param ctx the parse tree
	 */
	void exitNotNullConstraint(MDLParser.NotNullConstraintContext ctx);
	/**
	 * Enter a parse tree produced by the {@code UniqueConstraint}
	 * labeled alternative in {@link MDLParser#attributeConstraint}.
	 * @param ctx the parse tree
	 */
	void enterUniqueConstraint(MDLParser.UniqueConstraintContext ctx);
	/**
	 * Exit a parse tree produced by the {@code UniqueConstraint}
	 * labeled alternative in {@link MDLParser#attributeConstraint}.
	 * @param ctx the parse tree
	 */
	void exitUniqueConstraint(MDLParser.UniqueConstraintContext ctx);
	/**
	 * Enter a parse tree produced by the {@code DefaultConstraint}
	 * labeled alternative in {@link MDLParser#attributeConstraint}.
	 * @param ctx the parse tree
	 */
	void enterDefaultConstraint(MDLParser.DefaultConstraintContext ctx);
	/**
	 * Exit a parse tree produced by the {@code DefaultConstraint}
	 * labeled alternative in {@link MDLParser#attributeConstraint}.
	 * @param ctx the parse tree
	 */
	void exitDefaultConstraint(MDLParser.DefaultConstraintContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#defaultValue}.
	 * @param ctx the parse tree
	 */
	void enterDefaultValue(MDLParser.DefaultValueContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#defaultValue}.
	 * @param ctx the parse tree
	 */
	void exitDefaultValue(MDLParser.DefaultValueContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#indexClause}.
	 * @param ctx the parse tree
	 */
	void enterIndexClause(MDLParser.IndexClauseContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#indexClause}.
	 * @param ctx the parse tree
	 */
	void exitIndexClause(MDLParser.IndexClauseContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#indexColumnList}.
	 * @param ctx the parse tree
	 */
	void enterIndexColumnList(MDLParser.IndexColumnListContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#indexColumnList}.
	 * @param ctx the parse tree
	 */
	void exitIndexColumnList(MDLParser.IndexColumnListContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#indexColumn}.
	 * @param ctx the parse tree
	 */
	void enterIndexColumn(MDLParser.IndexColumnContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#indexColumn}.
	 * @param ctx the parse tree
	 */
	void exitIndexColumn(MDLParser.IndexColumnContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#dropEntityStatement}.
	 * @param ctx the parse tree
	 */
	void enterDropEntityStatement(MDLParser.DropEntityStatementContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#dropEntityStatement}.
	 * @param ctx the parse tree
	 */
	void exitDropEntityStatement(MDLParser.DropEntityStatementContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#createViewEntityStatement}.
	 * @param ctx the parse tree
	 */
	void enterCreateViewEntityStatement(MDLParser.CreateViewEntityStatementContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#createViewEntityStatement}.
	 * @param ctx the parse tree
	 */
	void exitCreateViewEntityStatement(MDLParser.CreateViewEntityStatementContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#viewAttributeList}.
	 * @param ctx the parse tree
	 */
	void enterViewAttributeList(MDLParser.ViewAttributeListContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#viewAttributeList}.
	 * @param ctx the parse tree
	 */
	void exitViewAttributeList(MDLParser.ViewAttributeListContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#viewAttribute}.
	 * @param ctx the parse tree
	 */
	void enterViewAttribute(MDLParser.ViewAttributeContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#viewAttribute}.
	 * @param ctx the parse tree
	 */
	void exitViewAttribute(MDLParser.ViewAttributeContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#oqlQuery}.
	 * @param ctx the parse tree
	 */
	void enterOqlQuery(MDLParser.OqlQueryContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#oqlQuery}.
	 * @param ctx the parse tree
	 */
	void exitOqlQuery(MDLParser.OqlQueryContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#selectClause}.
	 * @param ctx the parse tree
	 */
	void enterSelectClause(MDLParser.SelectClauseContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#selectClause}.
	 * @param ctx the parse tree
	 */
	void exitSelectClause(MDLParser.SelectClauseContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#selectItems}.
	 * @param ctx the parse tree
	 */
	void enterSelectItems(MDLParser.SelectItemsContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#selectItems}.
	 * @param ctx the parse tree
	 */
	void exitSelectItems(MDLParser.SelectItemsContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#selectItem}.
	 * @param ctx the parse tree
	 */
	void enterSelectItem(MDLParser.SelectItemContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#selectItem}.
	 * @param ctx the parse tree
	 */
	void exitSelectItem(MDLParser.SelectItemContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#fromClause}.
	 * @param ctx the parse tree
	 */
	void enterFromClause(MDLParser.FromClauseContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#fromClause}.
	 * @param ctx the parse tree
	 */
	void exitFromClause(MDLParser.FromClauseContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#fromItem}.
	 * @param ctx the parse tree
	 */
	void enterFromItem(MDLParser.FromItemContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#fromItem}.
	 * @param ctx the parse tree
	 */
	void exitFromItem(MDLParser.FromItemContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#joinClause}.
	 * @param ctx the parse tree
	 */
	void enterJoinClause(MDLParser.JoinClauseContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#joinClause}.
	 * @param ctx the parse tree
	 */
	void exitJoinClause(MDLParser.JoinClauseContext ctx);
	/**
	 * Enter a parse tree produced by the {@code SimpleJoinTarget}
	 * labeled alternative in {@link MDLParser#joinTarget}.
	 * @param ctx the parse tree
	 */
	void enterSimpleJoinTarget(MDLParser.SimpleJoinTargetContext ctx);
	/**
	 * Exit a parse tree produced by the {@code SimpleJoinTarget}
	 * labeled alternative in {@link MDLParser#joinTarget}.
	 * @param ctx the parse tree
	 */
	void exitSimpleJoinTarget(MDLParser.SimpleJoinTargetContext ctx);
	/**
	 * Enter a parse tree produced by the {@code AssociationJoinTarget}
	 * labeled alternative in {@link MDLParser#joinTarget}.
	 * @param ctx the parse tree
	 */
	void enterAssociationJoinTarget(MDLParser.AssociationJoinTargetContext ctx);
	/**
	 * Exit a parse tree produced by the {@code AssociationJoinTarget}
	 * labeled alternative in {@link MDLParser#joinTarget}.
	 * @param ctx the parse tree
	 */
	void exitAssociationJoinTarget(MDLParser.AssociationJoinTargetContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#whereClause}.
	 * @param ctx the parse tree
	 */
	void enterWhereClause(MDLParser.WhereClauseContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#whereClause}.
	 * @param ctx the parse tree
	 */
	void exitWhereClause(MDLParser.WhereClauseContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#groupByClause}.
	 * @param ctx the parse tree
	 */
	void enterGroupByClause(MDLParser.GroupByClauseContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#groupByClause}.
	 * @param ctx the parse tree
	 */
	void exitGroupByClause(MDLParser.GroupByClauseContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#orderByClause}.
	 * @param ctx the parse tree
	 */
	void enterOrderByClause(MDLParser.OrderByClauseContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#orderByClause}.
	 * @param ctx the parse tree
	 */
	void exitOrderByClause(MDLParser.OrderByClauseContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#orderByItem}.
	 * @param ctx the parse tree
	 */
	void enterOrderByItem(MDLParser.OrderByItemContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#orderByItem}.
	 * @param ctx the parse tree
	 */
	void exitOrderByItem(MDLParser.OrderByItemContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#limitClause}.
	 * @param ctx the parse tree
	 */
	void enterLimitClause(MDLParser.LimitClauseContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#limitClause}.
	 * @param ctx the parse tree
	 */
	void exitLimitClause(MDLParser.LimitClauseContext ctx);
	/**
	 * Enter a parse tree produced by the {@code AndExpr}
	 * labeled alternative in {@link MDLParser#expression}.
	 * @param ctx the parse tree
	 */
	void enterAndExpr(MDLParser.AndExprContext ctx);
	/**
	 * Exit a parse tree produced by the {@code AndExpr}
	 * labeled alternative in {@link MDLParser#expression}.
	 * @param ctx the parse tree
	 */
	void exitAndExpr(MDLParser.AndExprContext ctx);
	/**
	 * Enter a parse tree produced by the {@code StringExpr}
	 * labeled alternative in {@link MDLParser#expression}.
	 * @param ctx the parse tree
	 */
	void enterStringExpr(MDLParser.StringExprContext ctx);
	/**
	 * Exit a parse tree produced by the {@code StringExpr}
	 * labeled alternative in {@link MDLParser#expression}.
	 * @param ctx the parse tree
	 */
	void exitStringExpr(MDLParser.StringExprContext ctx);
	/**
	 * Enter a parse tree produced by the {@code IdentExpr}
	 * labeled alternative in {@link MDLParser#expression}.
	 * @param ctx the parse tree
	 */
	void enterIdentExpr(MDLParser.IdentExprContext ctx);
	/**
	 * Exit a parse tree produced by the {@code IdentExpr}
	 * labeled alternative in {@link MDLParser#expression}.
	 * @param ctx the parse tree
	 */
	void exitIdentExpr(MDLParser.IdentExprContext ctx);
	/**
	 * Enter a parse tree produced by the {@code TrueExpr}
	 * labeled alternative in {@link MDLParser#expression}.
	 * @param ctx the parse tree
	 */
	void enterTrueExpr(MDLParser.TrueExprContext ctx);
	/**
	 * Exit a parse tree produced by the {@code TrueExpr}
	 * labeled alternative in {@link MDLParser#expression}.
	 * @param ctx the parse tree
	 */
	void exitTrueExpr(MDLParser.TrueExprContext ctx);
	/**
	 * Enter a parse tree produced by the {@code IsNullExpr}
	 * labeled alternative in {@link MDLParser#expression}.
	 * @param ctx the parse tree
	 */
	void enterIsNullExpr(MDLParser.IsNullExprContext ctx);
	/**
	 * Exit a parse tree produced by the {@code IsNullExpr}
	 * labeled alternative in {@link MDLParser#expression}.
	 * @param ctx the parse tree
	 */
	void exitIsNullExpr(MDLParser.IsNullExprContext ctx);
	/**
	 * Enter a parse tree produced by the {@code StarExpr}
	 * labeled alternative in {@link MDLParser#expression}.
	 * @param ctx the parse tree
	 */
	void enterStarExpr(MDLParser.StarExprContext ctx);
	/**
	 * Exit a parse tree produced by the {@code StarExpr}
	 * labeled alternative in {@link MDLParser#expression}.
	 * @param ctx the parse tree
	 */
	void exitStarExpr(MDLParser.StarExprContext ctx);
	/**
	 * Enter a parse tree produced by the {@code FuncExpr}
	 * labeled alternative in {@link MDLParser#expression}.
	 * @param ctx the parse tree
	 */
	void enterFuncExpr(MDLParser.FuncExprContext ctx);
	/**
	 * Exit a parse tree produced by the {@code FuncExpr}
	 * labeled alternative in {@link MDLParser#expression}.
	 * @param ctx the parse tree
	 */
	void exitFuncExpr(MDLParser.FuncExprContext ctx);
	/**
	 * Enter a parse tree produced by the {@code QualifiedExpr}
	 * labeled alternative in {@link MDLParser#expression}.
	 * @param ctx the parse tree
	 */
	void enterQualifiedExpr(MDLParser.QualifiedExprContext ctx);
	/**
	 * Exit a parse tree produced by the {@code QualifiedExpr}
	 * labeled alternative in {@link MDLParser#expression}.
	 * @param ctx the parse tree
	 */
	void exitQualifiedExpr(MDLParser.QualifiedExprContext ctx);
	/**
	 * Enter a parse tree produced by the {@code DecimalExpr}
	 * labeled alternative in {@link MDLParser#expression}.
	 * @param ctx the parse tree
	 */
	void enterDecimalExpr(MDLParser.DecimalExprContext ctx);
	/**
	 * Exit a parse tree produced by the {@code DecimalExpr}
	 * labeled alternative in {@link MDLParser#expression}.
	 * @param ctx the parse tree
	 */
	void exitDecimalExpr(MDLParser.DecimalExprContext ctx);
	/**
	 * Enter a parse tree produced by the {@code OrExpr}
	 * labeled alternative in {@link MDLParser#expression}.
	 * @param ctx the parse tree
	 */
	void enterOrExpr(MDLParser.OrExprContext ctx);
	/**
	 * Exit a parse tree produced by the {@code OrExpr}
	 * labeled alternative in {@link MDLParser#expression}.
	 * @param ctx the parse tree
	 */
	void exitOrExpr(MDLParser.OrExprContext ctx);
	/**
	 * Enter a parse tree produced by the {@code FalseExpr}
	 * labeled alternative in {@link MDLParser#expression}.
	 * @param ctx the parse tree
	 */
	void enterFalseExpr(MDLParser.FalseExprContext ctx);
	/**
	 * Exit a parse tree produced by the {@code FalseExpr}
	 * labeled alternative in {@link MDLParser#expression}.
	 * @param ctx the parse tree
	 */
	void exitFalseExpr(MDLParser.FalseExprContext ctx);
	/**
	 * Enter a parse tree produced by the {@code InSubqueryExpr}
	 * labeled alternative in {@link MDLParser#expression}.
	 * @param ctx the parse tree
	 */
	void enterInSubqueryExpr(MDLParser.InSubqueryExprContext ctx);
	/**
	 * Exit a parse tree produced by the {@code InSubqueryExpr}
	 * labeled alternative in {@link MDLParser#expression}.
	 * @param ctx the parse tree
	 */
	void exitInSubqueryExpr(MDLParser.InSubqueryExprContext ctx);
	/**
	 * Enter a parse tree produced by the {@code MulDivExpr}
	 * labeled alternative in {@link MDLParser#expression}.
	 * @param ctx the parse tree
	 */
	void enterMulDivExpr(MDLParser.MulDivExprContext ctx);
	/**
	 * Exit a parse tree produced by the {@code MulDivExpr}
	 * labeled alternative in {@link MDLParser#expression}.
	 * @param ctx the parse tree
	 */
	void exitMulDivExpr(MDLParser.MulDivExprContext ctx);
	/**
	 * Enter a parse tree produced by the {@code DivisionExpr}
	 * labeled alternative in {@link MDLParser#expression}.
	 * @param ctx the parse tree
	 */
	void enterDivisionExpr(MDLParser.DivisionExprContext ctx);
	/**
	 * Exit a parse tree produced by the {@code DivisionExpr}
	 * labeled alternative in {@link MDLParser#expression}.
	 * @param ctx the parse tree
	 */
	void exitDivisionExpr(MDLParser.DivisionExprContext ctx);
	/**
	 * Enter a parse tree produced by the {@code CompareExpr}
	 * labeled alternative in {@link MDLParser#expression}.
	 * @param ctx the parse tree
	 */
	void enterCompareExpr(MDLParser.CompareExprContext ctx);
	/**
	 * Exit a parse tree produced by the {@code CompareExpr}
	 * labeled alternative in {@link MDLParser#expression}.
	 * @param ctx the parse tree
	 */
	void exitCompareExpr(MDLParser.CompareExprContext ctx);
	/**
	 * Enter a parse tree produced by the {@code FieldAccessExpr}
	 * labeled alternative in {@link MDLParser#expression}.
	 * @param ctx the parse tree
	 */
	void enterFieldAccessExpr(MDLParser.FieldAccessExprContext ctx);
	/**
	 * Exit a parse tree produced by the {@code FieldAccessExpr}
	 * labeled alternative in {@link MDLParser#expression}.
	 * @param ctx the parse tree
	 */
	void exitFieldAccessExpr(MDLParser.FieldAccessExprContext ctx);
	/**
	 * Enter a parse tree produced by the {@code NotExpr}
	 * labeled alternative in {@link MDLParser#expression}.
	 * @param ctx the parse tree
	 */
	void enterNotExpr(MDLParser.NotExprContext ctx);
	/**
	 * Exit a parse tree produced by the {@code NotExpr}
	 * labeled alternative in {@link MDLParser#expression}.
	 * @param ctx the parse tree
	 */
	void exitNotExpr(MDLParser.NotExprContext ctx);
	/**
	 * Enter a parse tree produced by the {@code IntExpr}
	 * labeled alternative in {@link MDLParser#expression}.
	 * @param ctx the parse tree
	 */
	void enterIntExpr(MDLParser.IntExprContext ctx);
	/**
	 * Exit a parse tree produced by the {@code IntExpr}
	 * labeled alternative in {@link MDLParser#expression}.
	 * @param ctx the parse tree
	 */
	void exitIntExpr(MDLParser.IntExprContext ctx);
	/**
	 * Enter a parse tree produced by the {@code SysVarExpr}
	 * labeled alternative in {@link MDLParser#expression}.
	 * @param ctx the parse tree
	 */
	void enterSysVarExpr(MDLParser.SysVarExprContext ctx);
	/**
	 * Exit a parse tree produced by the {@code SysVarExpr}
	 * labeled alternative in {@link MDLParser#expression}.
	 * @param ctx the parse tree
	 */
	void exitSysVarExpr(MDLParser.SysVarExprContext ctx);
	/**
	 * Enter a parse tree produced by the {@code InExpr}
	 * labeled alternative in {@link MDLParser#expression}.
	 * @param ctx the parse tree
	 */
	void enterInExpr(MDLParser.InExprContext ctx);
	/**
	 * Exit a parse tree produced by the {@code InExpr}
	 * labeled alternative in {@link MDLParser#expression}.
	 * @param ctx the parse tree
	 */
	void exitInExpr(MDLParser.InExprContext ctx);
	/**
	 * Enter a parse tree produced by the {@code ParenExpr}
	 * labeled alternative in {@link MDLParser#expression}.
	 * @param ctx the parse tree
	 */
	void enterParenExpr(MDLParser.ParenExprContext ctx);
	/**
	 * Exit a parse tree produced by the {@code ParenExpr}
	 * labeled alternative in {@link MDLParser#expression}.
	 * @param ctx the parse tree
	 */
	void exitParenExpr(MDLParser.ParenExprContext ctx);
	/**
	 * Enter a parse tree produced by the {@code CaseExpr}
	 * labeled alternative in {@link MDLParser#expression}.
	 * @param ctx the parse tree
	 */
	void enterCaseExpr(MDLParser.CaseExprContext ctx);
	/**
	 * Exit a parse tree produced by the {@code CaseExpr}
	 * labeled alternative in {@link MDLParser#expression}.
	 * @param ctx the parse tree
	 */
	void exitCaseExpr(MDLParser.CaseExprContext ctx);
	/**
	 * Enter a parse tree produced by the {@code AddSubExpr}
	 * labeled alternative in {@link MDLParser#expression}.
	 * @param ctx the parse tree
	 */
	void enterAddSubExpr(MDLParser.AddSubExprContext ctx);
	/**
	 * Exit a parse tree produced by the {@code AddSubExpr}
	 * labeled alternative in {@link MDLParser#expression}.
	 * @param ctx the parse tree
	 */
	void exitAddSubExpr(MDLParser.AddSubExprContext ctx);
	/**
	 * Enter a parse tree produced by the {@code SubqueryExpr}
	 * labeled alternative in {@link MDLParser#expression}.
	 * @param ctx the parse tree
	 */
	void enterSubqueryExpr(MDLParser.SubqueryExprContext ctx);
	/**
	 * Exit a parse tree produced by the {@code SubqueryExpr}
	 * labeled alternative in {@link MDLParser#expression}.
	 * @param ctx the parse tree
	 */
	void exitSubqueryExpr(MDLParser.SubqueryExprContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#comparisonOp}.
	 * @param ctx the parse tree
	 */
	void enterComparisonOp(MDLParser.ComparisonOpContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#comparisonOp}.
	 * @param ctx the parse tree
	 */
	void exitComparisonOp(MDLParser.ComparisonOpContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#functionCall}.
	 * @param ctx the parse tree
	 */
	void enterFunctionCall(MDLParser.FunctionCallContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#functionCall}.
	 * @param ctx the parse tree
	 */
	void exitFunctionCall(MDLParser.FunctionCallContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#aggregateFunction}.
	 * @param ctx the parse tree
	 */
	void enterAggregateFunction(MDLParser.AggregateFunctionContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#aggregateFunction}.
	 * @param ctx the parse tree
	 */
	void exitAggregateFunction(MDLParser.AggregateFunctionContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#caseExpression}.
	 * @param ctx the parse tree
	 */
	void enterCaseExpression(MDLParser.CaseExpressionContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#caseExpression}.
	 * @param ctx the parse tree
	 */
	void exitCaseExpression(MDLParser.CaseExpressionContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#expressionList}.
	 * @param ctx the parse tree
	 */
	void enterExpressionList(MDLParser.ExpressionListContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#expressionList}.
	 * @param ctx the parse tree
	 */
	void exitExpressionList(MDLParser.ExpressionListContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#systemVariable}.
	 * @param ctx the parse tree
	 */
	void enterSystemVariable(MDLParser.SystemVariableContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#systemVariable}.
	 * @param ctx the parse tree
	 */
	void exitSystemVariable(MDLParser.SystemVariableContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#createAssociationStatement}.
	 * @param ctx the parse tree
	 */
	void enterCreateAssociationStatement(MDLParser.CreateAssociationStatementContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#createAssociationStatement}.
	 * @param ctx the parse tree
	 */
	void exitCreateAssociationStatement(MDLParser.CreateAssociationStatementContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#associationType}.
	 * @param ctx the parse tree
	 */
	void enterAssociationType(MDLParser.AssociationTypeContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#associationType}.
	 * @param ctx the parse tree
	 */
	void exitAssociationType(MDLParser.AssociationTypeContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#ownerType}.
	 * @param ctx the parse tree
	 */
	void enterOwnerType(MDLParser.OwnerTypeContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#ownerType}.
	 * @param ctx the parse tree
	 */
	void exitOwnerType(MDLParser.OwnerTypeContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#deleteBehavior}.
	 * @param ctx the parse tree
	 */
	void enterDeleteBehavior(MDLParser.DeleteBehaviorContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#deleteBehavior}.
	 * @param ctx the parse tree
	 */
	void exitDeleteBehavior(MDLParser.DeleteBehaviorContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#dropAssociationStatement}.
	 * @param ctx the parse tree
	 */
	void enterDropAssociationStatement(MDLParser.DropAssociationStatementContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#dropAssociationStatement}.
	 * @param ctx the parse tree
	 */
	void exitDropAssociationStatement(MDLParser.DropAssociationStatementContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#queryStatement}.
	 * @param ctx the parse tree
	 */
	void enterQueryStatement(MDLParser.QueryStatementContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#queryStatement}.
	 * @param ctx the parse tree
	 */
	void exitQueryStatement(MDLParser.QueryStatementContext ctx);
	/**
	 * Enter a parse tree produced by the {@code ShowModules}
	 * labeled alternative in {@link MDLParser#showStatement}.
	 * @param ctx the parse tree
	 */
	void enterShowModules(MDLParser.ShowModulesContext ctx);
	/**
	 * Exit a parse tree produced by the {@code ShowModules}
	 * labeled alternative in {@link MDLParser#showStatement}.
	 * @param ctx the parse tree
	 */
	void exitShowModules(MDLParser.ShowModulesContext ctx);
	/**
	 * Enter a parse tree produced by the {@code ShowEnumerations}
	 * labeled alternative in {@link MDLParser#showStatement}.
	 * @param ctx the parse tree
	 */
	void enterShowEnumerations(MDLParser.ShowEnumerationsContext ctx);
	/**
	 * Exit a parse tree produced by the {@code ShowEnumerations}
	 * labeled alternative in {@link MDLParser#showStatement}.
	 * @param ctx the parse tree
	 */
	void exitShowEnumerations(MDLParser.ShowEnumerationsContext ctx);
	/**
	 * Enter a parse tree produced by the {@code ShowEntities}
	 * labeled alternative in {@link MDLParser#showStatement}.
	 * @param ctx the parse tree
	 */
	void enterShowEntities(MDLParser.ShowEntitiesContext ctx);
	/**
	 * Exit a parse tree produced by the {@code ShowEntities}
	 * labeled alternative in {@link MDLParser#showStatement}.
	 * @param ctx the parse tree
	 */
	void exitShowEntities(MDLParser.ShowEntitiesContext ctx);
	/**
	 * Enter a parse tree produced by the {@code ShowEntity}
	 * labeled alternative in {@link MDLParser#showStatement}.
	 * @param ctx the parse tree
	 */
	void enterShowEntity(MDLParser.ShowEntityContext ctx);
	/**
	 * Exit a parse tree produced by the {@code ShowEntity}
	 * labeled alternative in {@link MDLParser#showStatement}.
	 * @param ctx the parse tree
	 */
	void exitShowEntity(MDLParser.ShowEntityContext ctx);
	/**
	 * Enter a parse tree produced by the {@code ShowAssociations}
	 * labeled alternative in {@link MDLParser#showStatement}.
	 * @param ctx the parse tree
	 */
	void enterShowAssociations(MDLParser.ShowAssociationsContext ctx);
	/**
	 * Exit a parse tree produced by the {@code ShowAssociations}
	 * labeled alternative in {@link MDLParser#showStatement}.
	 * @param ctx the parse tree
	 */
	void exitShowAssociations(MDLParser.ShowAssociationsContext ctx);
	/**
	 * Enter a parse tree produced by the {@code ShowAssociation}
	 * labeled alternative in {@link MDLParser#showStatement}.
	 * @param ctx the parse tree
	 */
	void enterShowAssociation(MDLParser.ShowAssociationContext ctx);
	/**
	 * Exit a parse tree produced by the {@code ShowAssociation}
	 * labeled alternative in {@link MDLParser#showStatement}.
	 * @param ctx the parse tree
	 */
	void exitShowAssociation(MDLParser.ShowAssociationContext ctx);
	/**
	 * Enter a parse tree produced by the {@code DescribeEnumeration}
	 * labeled alternative in {@link MDLParser#describeStatement}.
	 * @param ctx the parse tree
	 */
	void enterDescribeEnumeration(MDLParser.DescribeEnumerationContext ctx);
	/**
	 * Exit a parse tree produced by the {@code DescribeEnumeration}
	 * labeled alternative in {@link MDLParser#describeStatement}.
	 * @param ctx the parse tree
	 */
	void exitDescribeEnumeration(MDLParser.DescribeEnumerationContext ctx);
	/**
	 * Enter a parse tree produced by the {@code DescribeEntity}
	 * labeled alternative in {@link MDLParser#describeStatement}.
	 * @param ctx the parse tree
	 */
	void enterDescribeEntity(MDLParser.DescribeEntityContext ctx);
	/**
	 * Exit a parse tree produced by the {@code DescribeEntity}
	 * labeled alternative in {@link MDLParser#describeStatement}.
	 * @param ctx the parse tree
	 */
	void exitDescribeEntity(MDLParser.DescribeEntityContext ctx);
	/**
	 * Enter a parse tree produced by the {@code DescribeAssociation}
	 * labeled alternative in {@link MDLParser#describeStatement}.
	 * @param ctx the parse tree
	 */
	void enterDescribeAssociation(MDLParser.DescribeAssociationContext ctx);
	/**
	 * Exit a parse tree produced by the {@code DescribeAssociation}
	 * labeled alternative in {@link MDLParser#describeStatement}.
	 * @param ctx the parse tree
	 */
	void exitDescribeAssociation(MDLParser.DescribeAssociationContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#repositoryStatement}.
	 * @param ctx the parse tree
	 */
	void enterRepositoryStatement(MDLParser.RepositoryStatementContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#repositoryStatement}.
	 * @param ctx the parse tree
	 */
	void exitRepositoryStatement(MDLParser.RepositoryStatementContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#commitStatement}.
	 * @param ctx the parse tree
	 */
	void enterCommitStatement(MDLParser.CommitStatementContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#commitStatement}.
	 * @param ctx the parse tree
	 */
	void exitCommitStatement(MDLParser.CommitStatementContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#updateStatement}.
	 * @param ctx the parse tree
	 */
	void enterUpdateStatement(MDLParser.UpdateStatementContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#updateStatement}.
	 * @param ctx the parse tree
	 */
	void exitUpdateStatement(MDLParser.UpdateStatementContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#refreshStatement}.
	 * @param ctx the parse tree
	 */
	void enterRefreshStatement(MDLParser.RefreshStatementContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#refreshStatement}.
	 * @param ctx the parse tree
	 */
	void exitRefreshStatement(MDLParser.RefreshStatementContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#sessionStatement}.
	 * @param ctx the parse tree
	 */
	void enterSessionStatement(MDLParser.SessionStatementContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#sessionStatement}.
	 * @param ctx the parse tree
	 */
	void exitSessionStatement(MDLParser.SessionStatementContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#setStatement}.
	 * @param ctx the parse tree
	 */
	void enterSetStatement(MDLParser.SetStatementContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#setStatement}.
	 * @param ctx the parse tree
	 */
	void exitSetStatement(MDLParser.SetStatementContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#helpStatement}.
	 * @param ctx the parse tree
	 */
	void enterHelpStatement(MDLParser.HelpStatementContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#helpStatement}.
	 * @param ctx the parse tree
	 */
	void exitHelpStatement(MDLParser.HelpStatementContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#exitStatement}.
	 * @param ctx the parse tree
	 */
	void enterExitStatement(MDLParser.ExitStatementContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#exitStatement}.
	 * @param ctx the parse tree
	 */
	void exitExitStatement(MDLParser.ExitStatementContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#executeScriptStatement}.
	 * @param ctx the parse tree
	 */
	void enterExecuteScriptStatement(MDLParser.ExecuteScriptStatementContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#executeScriptStatement}.
	 * @param ctx the parse tree
	 */
	void exitExecuteScriptStatement(MDLParser.ExecuteScriptStatementContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#qualifiedName}.
	 * @param ctx the parse tree
	 */
	void enterQualifiedName(MDLParser.QualifiedNameContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#qualifiedName}.
	 * @param ctx the parse tree
	 */
	void exitQualifiedName(MDLParser.QualifiedNameContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#documentation}.
	 * @param ctx the parse tree
	 */
	void enterDocumentation(MDLParser.DocumentationContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#documentation}.
	 * @param ctx the parse tree
	 */
	void exitDocumentation(MDLParser.DocumentationContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#positionAnnotation}.
	 * @param ctx the parse tree
	 */
	void enterPositionAnnotation(MDLParser.PositionAnnotationContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#positionAnnotation}.
	 * @param ctx the parse tree
	 */
	void exitPositionAnnotation(MDLParser.PositionAnnotationContext ctx);
	/**
	 * Enter a parse tree produced by {@link MDLParser#stringLiteral}.
	 * @param ctx the parse tree
	 */
	void enterStringLiteral(MDLParser.StringLiteralContext ctx);
	/**
	 * Exit a parse tree produced by {@link MDLParser#stringLiteral}.
	 * @param ctx the parse tree
	 */
	void exitStringLiteral(MDLParser.StringLiteralContext ctx);
}