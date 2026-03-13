// SPDX-License-Identifier: Apache-2.0

// Code generated from MDL.g4 by ANTLR 4.7.2. DO NOT EDIT.

package parser // MDL
import "github.com/antlr/antlr4/runtime/Go/antlr"

// MDLListener is a complete listener for a parse tree produced by MDLParser.
type MDLListener interface {
	antlr.ParseTreeListener

	// EnterProgram is called when entering the program production.
	EnterProgram(c *ProgramContext)

	// EnterStatement is called when entering the statement production.
	EnterStatement(c *StatementContext)

	// EnterTerminator is called when entering the terminator production.
	EnterTerminator(c *TerminatorContext)

	// EnterConnectionStatement is called when entering the connectionStatement production.
	EnterConnectionStatement(c *ConnectionStatementContext)

	// EnterConnectLocal is called when entering the ConnectLocal production.
	EnterConnectLocal(c *ConnectLocalContext)

	// EnterConnectFilesystem is called when entering the ConnectFilesystem production.
	EnterConnectFilesystem(c *ConnectFilesystemContext)

	// EnterDisconnectStatement is called when entering the disconnectStatement production.
	EnterDisconnectStatement(c *DisconnectStatementContext)

	// EnterStatusStatement is called when entering the statusStatement production.
	EnterStatusStatement(c *StatusStatementContext)

	// EnterDdlStatement is called when entering the ddlStatement production.
	EnterDdlStatement(c *DdlStatementContext)

	// EnterCreateStatement is called when entering the createStatement production.
	EnterCreateStatement(c *CreateStatementContext)

	// EnterCreateModuleStatement is called when entering the createModuleStatement production.
	EnterCreateModuleStatement(c *CreateModuleStatementContext)

	// EnterAlterStatement is called when entering the alterStatement production.
	EnterAlterStatement(c *AlterStatementContext)

	// EnterDropStatement is called when entering the dropStatement production.
	EnterDropStatement(c *DropStatementContext)

	// EnterDropModuleStatement is called when entering the dropModuleStatement production.
	EnterDropModuleStatement(c *DropModuleStatementContext)

	// EnterCreateEnumerationStatement is called when entering the createEnumerationStatement production.
	EnterCreateEnumerationStatement(c *CreateEnumerationStatementContext)

	// EnterEnumValueList is called when entering the enumValueList production.
	EnterEnumValueList(c *EnumValueListContext)

	// EnterEnumValue is called when entering the enumValue production.
	EnterEnumValue(c *EnumValueContext)

	// EnterAlterEnumerationStatement is called when entering the alterEnumerationStatement production.
	EnterAlterEnumerationStatement(c *AlterEnumerationStatementContext)

	// EnterAddEnumValue is called when entering the AddEnumValue production.
	EnterAddEnumValue(c *AddEnumValueContext)

	// EnterDropEnumValue is called when entering the DropEnumValue production.
	EnterDropEnumValue(c *DropEnumValueContext)

	// EnterRenameEnumValue is called when entering the RenameEnumValue production.
	EnterRenameEnumValue(c *RenameEnumValueContext)

	// EnterDropEnumerationStatement is called when entering the dropEnumerationStatement production.
	EnterDropEnumerationStatement(c *DropEnumerationStatementContext)

	// EnterCreateEntityStatement is called when entering the createEntityStatement production.
	EnterCreateEntityStatement(c *CreateEntityStatementContext)

	// EnterCreateOrModify is called when entering the createOrModify production.
	EnterCreateOrModify(c *CreateOrModifyContext)

	// EnterEntityType is called when entering the entityType production.
	EnterEntityType(c *EntityTypeContext)

	// EnterAttributeList is called when entering the attributeList production.
	EnterAttributeList(c *AttributeListContext)

	// EnterAttribute is called when entering the attribute production.
	EnterAttribute(c *AttributeContext)

	// EnterAttributeName is called when entering the attributeName production.
	EnterAttributeName(c *AttributeNameContext)

	// EnterRenamedFromAnnotation is called when entering the renamedFromAnnotation production.
	EnterRenamedFromAnnotation(c *RenamedFromAnnotationContext)

	// EnterStringType is called when entering the StringType production.
	EnterStringType(c *StringTypeContext)

	// EnterIntegerType is called when entering the IntegerType production.
	EnterIntegerType(c *IntegerTypeContext)

	// EnterLongType is called when entering the LongType production.
	EnterLongType(c *LongTypeContext)

	// EnterDecimalType is called when entering the DecimalType production.
	EnterDecimalType(c *DecimalTypeContext)

	// EnterBooleanType is called when entering the BooleanType production.
	EnterBooleanType(c *BooleanTypeContext)

	// EnterDateTimeType is called when entering the DateTimeType production.
	EnterDateTimeType(c *DateTimeTypeContext)

	// EnterDateType is called when entering the DateType production.
	EnterDateType(c *DateTypeContext)

	// EnterAutoNumberType is called when entering the AutoNumberType production.
	EnterAutoNumberType(c *AutoNumberTypeContext)

	// EnterBinaryType is called when entering the BinaryType production.
	EnterBinaryType(c *BinaryTypeContext)

	// EnterEnumerationType is called when entering the EnumerationType production.
	EnterEnumerationType(c *EnumerationTypeContext)

	// EnterAttributeConstraints is called when entering the attributeConstraints production.
	EnterAttributeConstraints(c *AttributeConstraintsContext)

	// EnterNotNullConstraint is called when entering the NotNullConstraint production.
	EnterNotNullConstraint(c *NotNullConstraintContext)

	// EnterUniqueConstraint is called when entering the UniqueConstraint production.
	EnterUniqueConstraint(c *UniqueConstraintContext)

	// EnterDefaultConstraint is called when entering the DefaultConstraint production.
	EnterDefaultConstraint(c *DefaultConstraintContext)

	// EnterDefaultValue is called when entering the defaultValue production.
	EnterDefaultValue(c *DefaultValueContext)

	// EnterIndexClause is called when entering the indexClause production.
	EnterIndexClause(c *IndexClauseContext)

	// EnterIndexColumnList is called when entering the indexColumnList production.
	EnterIndexColumnList(c *IndexColumnListContext)

	// EnterIndexColumn is called when entering the indexColumn production.
	EnterIndexColumn(c *IndexColumnContext)

	// EnterDropEntityStatement is called when entering the dropEntityStatement production.
	EnterDropEntityStatement(c *DropEntityStatementContext)

	// EnterCreateViewEntityStatement is called when entering the createViewEntityStatement production.
	EnterCreateViewEntityStatement(c *CreateViewEntityStatementContext)

	// EnterViewAttributeList is called when entering the viewAttributeList production.
	EnterViewAttributeList(c *ViewAttributeListContext)

	// EnterViewAttribute is called when entering the viewAttribute production.
	EnterViewAttribute(c *ViewAttributeContext)

	// EnterOqlQuery is called when entering the oqlQuery production.
	EnterOqlQuery(c *OqlQueryContext)

	// EnterSelectClause is called when entering the selectClause production.
	EnterSelectClause(c *SelectClauseContext)

	// EnterSelectItems is called when entering the selectItems production.
	EnterSelectItems(c *SelectItemsContext)

	// EnterSelectItem is called when entering the selectItem production.
	EnterSelectItem(c *SelectItemContext)

	// EnterFromClause is called when entering the fromClause production.
	EnterFromClause(c *FromClauseContext)

	// EnterFromItem is called when entering the fromItem production.
	EnterFromItem(c *FromItemContext)

	// EnterJoinClause is called when entering the joinClause production.
	EnterJoinClause(c *JoinClauseContext)

	// EnterSimpleJoinTarget is called when entering the SimpleJoinTarget production.
	EnterSimpleJoinTarget(c *SimpleJoinTargetContext)

	// EnterAssociationJoinTarget is called when entering the AssociationJoinTarget production.
	EnterAssociationJoinTarget(c *AssociationJoinTargetContext)

	// EnterWhereClause is called when entering the whereClause production.
	EnterWhereClause(c *WhereClauseContext)

	// EnterGroupByClause is called when entering the groupByClause production.
	EnterGroupByClause(c *GroupByClauseContext)

	// EnterOrderByClause is called when entering the orderByClause production.
	EnterOrderByClause(c *OrderByClauseContext)

	// EnterOrderByItem is called when entering the orderByItem production.
	EnterOrderByItem(c *OrderByItemContext)

	// EnterLimitClause is called when entering the limitClause production.
	EnterLimitClause(c *LimitClauseContext)

	// EnterAndExpr is called when entering the AndExpr production.
	EnterAndExpr(c *AndExprContext)

	// EnterStringExpr is called when entering the StringExpr production.
	EnterStringExpr(c *StringExprContext)

	// EnterIdentExpr is called when entering the IdentExpr production.
	EnterIdentExpr(c *IdentExprContext)

	// EnterTrueExpr is called when entering the TrueExpr production.
	EnterTrueExpr(c *TrueExprContext)

	// EnterIsNullExpr is called when entering the IsNullExpr production.
	EnterIsNullExpr(c *IsNullExprContext)

	// EnterStarExpr is called when entering the StarExpr production.
	EnterStarExpr(c *StarExprContext)

	// EnterFuncExpr is called when entering the FuncExpr production.
	EnterFuncExpr(c *FuncExprContext)

	// EnterQualifiedExpr is called when entering the QualifiedExpr production.
	EnterQualifiedExpr(c *QualifiedExprContext)

	// EnterDecimalExpr is called when entering the DecimalExpr production.
	EnterDecimalExpr(c *DecimalExprContext)

	// EnterOrExpr is called when entering the OrExpr production.
	EnterOrExpr(c *OrExprContext)

	// EnterFalseExpr is called when entering the FalseExpr production.
	EnterFalseExpr(c *FalseExprContext)

	// EnterInSubqueryExpr is called when entering the InSubqueryExpr production.
	EnterInSubqueryExpr(c *InSubqueryExprContext)

	// EnterMulDivExpr is called when entering the MulDivExpr production.
	EnterMulDivExpr(c *MulDivExprContext)

	// EnterDivisionExpr is called when entering the DivisionExpr production.
	EnterDivisionExpr(c *DivisionExprContext)

	// EnterCompareExpr is called when entering the CompareExpr production.
	EnterCompareExpr(c *CompareExprContext)

	// EnterFieldAccessExpr is called when entering the FieldAccessExpr production.
	EnterFieldAccessExpr(c *FieldAccessExprContext)

	// EnterNotExpr is called when entering the NotExpr production.
	EnterNotExpr(c *NotExprContext)

	// EnterIntExpr is called when entering the IntExpr production.
	EnterIntExpr(c *IntExprContext)

	// EnterSysVarExpr is called when entering the SysVarExpr production.
	EnterSysVarExpr(c *SysVarExprContext)

	// EnterInExpr is called when entering the InExpr production.
	EnterInExpr(c *InExprContext)

	// EnterParenExpr is called when entering the ParenExpr production.
	EnterParenExpr(c *ParenExprContext)

	// EnterCaseExpr is called when entering the CaseExpr production.
	EnterCaseExpr(c *CaseExprContext)

	// EnterAddSubExpr is called when entering the AddSubExpr production.
	EnterAddSubExpr(c *AddSubExprContext)

	// EnterSubqueryExpr is called when entering the SubqueryExpr production.
	EnterSubqueryExpr(c *SubqueryExprContext)

	// EnterComparisonOp is called when entering the comparisonOp production.
	EnterComparisonOp(c *ComparisonOpContext)

	// EnterFunctionCall is called when entering the functionCall production.
	EnterFunctionCall(c *FunctionCallContext)

	// EnterAggregateFunction is called when entering the aggregateFunction production.
	EnterAggregateFunction(c *AggregateFunctionContext)

	// EnterCaseExpression is called when entering the caseExpression production.
	EnterCaseExpression(c *CaseExpressionContext)

	// EnterExpressionList is called when entering the expressionList production.
	EnterExpressionList(c *ExpressionListContext)

	// EnterSystemVariable is called when entering the systemVariable production.
	EnterSystemVariable(c *SystemVariableContext)

	// EnterCreateAssociationStatement is called when entering the createAssociationStatement production.
	EnterCreateAssociationStatement(c *CreateAssociationStatementContext)

	// EnterAssociationType is called when entering the associationType production.
	EnterAssociationType(c *AssociationTypeContext)

	// EnterOwnerType is called when entering the ownerType production.
	EnterOwnerType(c *OwnerTypeContext)

	// EnterDeleteBehavior is called when entering the deleteBehavior production.
	EnterDeleteBehavior(c *DeleteBehaviorContext)

	// EnterDropAssociationStatement is called when entering the dropAssociationStatement production.
	EnterDropAssociationStatement(c *DropAssociationStatementContext)

	// EnterQueryStatement is called when entering the queryStatement production.
	EnterQueryStatement(c *QueryStatementContext)

	// EnterShowModules is called when entering the ShowModules production.
	EnterShowModules(c *ShowModulesContext)

	// EnterShowEnumerations is called when entering the ShowEnumerations production.
	EnterShowEnumerations(c *ShowEnumerationsContext)

	// EnterShowEntities is called when entering the ShowEntities production.
	EnterShowEntities(c *ShowEntitiesContext)

	// EnterShowEntity is called when entering the ShowEntity production.
	EnterShowEntity(c *ShowEntityContext)

	// EnterShowAssociations is called when entering the ShowAssociations production.
	EnterShowAssociations(c *ShowAssociationsContext)

	// EnterShowAssociation is called when entering the ShowAssociation production.
	EnterShowAssociation(c *ShowAssociationContext)

	// EnterDescribeEnumeration is called when entering the DescribeEnumeration production.
	EnterDescribeEnumeration(c *DescribeEnumerationContext)

	// EnterDescribeEntity is called when entering the DescribeEntity production.
	EnterDescribeEntity(c *DescribeEntityContext)

	// EnterDescribeAssociation is called when entering the DescribeAssociation production.
	EnterDescribeAssociation(c *DescribeAssociationContext)

	// EnterRepositoryStatement is called when entering the repositoryStatement production.
	EnterRepositoryStatement(c *RepositoryStatementContext)

	// EnterCommitStatement is called when entering the commitStatement production.
	EnterCommitStatement(c *CommitStatementContext)

	// EnterUpdateStatement is called when entering the updateStatement production.
	EnterUpdateStatement(c *UpdateStatementContext)

	// EnterRefreshStatement is called when entering the refreshStatement production.
	EnterRefreshStatement(c *RefreshStatementContext)

	// EnterSessionStatement is called when entering the sessionStatement production.
	EnterSessionStatement(c *SessionStatementContext)

	// EnterSetStatement is called when entering the setStatement production.
	EnterSetStatement(c *SetStatementContext)

	// EnterHelpStatement is called when entering the helpStatement production.
	EnterHelpStatement(c *HelpStatementContext)

	// EnterExitStatement is called when entering the exitStatement production.
	EnterExitStatement(c *ExitStatementContext)

	// EnterExecuteScriptStatement is called when entering the executeScriptStatement production.
	EnterExecuteScriptStatement(c *ExecuteScriptStatementContext)

	// EnterQualifiedName is called when entering the qualifiedName production.
	EnterQualifiedName(c *QualifiedNameContext)

	// EnterDocumentation is called when entering the documentation production.
	EnterDocumentation(c *DocumentationContext)

	// EnterPositionAnnotation is called when entering the positionAnnotation production.
	EnterPositionAnnotation(c *PositionAnnotationContext)

	// EnterStringLiteral is called when entering the stringLiteral production.
	EnterStringLiteral(c *StringLiteralContext)

	// ExitProgram is called when exiting the program production.
	ExitProgram(c *ProgramContext)

	// ExitStatement is called when exiting the statement production.
	ExitStatement(c *StatementContext)

	// ExitTerminator is called when exiting the terminator production.
	ExitTerminator(c *TerminatorContext)

	// ExitConnectionStatement is called when exiting the connectionStatement production.
	ExitConnectionStatement(c *ConnectionStatementContext)

	// ExitConnectLocal is called when exiting the ConnectLocal production.
	ExitConnectLocal(c *ConnectLocalContext)

	// ExitConnectFilesystem is called when exiting the ConnectFilesystem production.
	ExitConnectFilesystem(c *ConnectFilesystemContext)

	// ExitDisconnectStatement is called when exiting the disconnectStatement production.
	ExitDisconnectStatement(c *DisconnectStatementContext)

	// ExitStatusStatement is called when exiting the statusStatement production.
	ExitStatusStatement(c *StatusStatementContext)

	// ExitDdlStatement is called when exiting the ddlStatement production.
	ExitDdlStatement(c *DdlStatementContext)

	// ExitCreateStatement is called when exiting the createStatement production.
	ExitCreateStatement(c *CreateStatementContext)

	// ExitCreateModuleStatement is called when exiting the createModuleStatement production.
	ExitCreateModuleStatement(c *CreateModuleStatementContext)

	// ExitAlterStatement is called when exiting the alterStatement production.
	ExitAlterStatement(c *AlterStatementContext)

	// ExitDropStatement is called when exiting the dropStatement production.
	ExitDropStatement(c *DropStatementContext)

	// ExitDropModuleStatement is called when exiting the dropModuleStatement production.
	ExitDropModuleStatement(c *DropModuleStatementContext)

	// ExitCreateEnumerationStatement is called when exiting the createEnumerationStatement production.
	ExitCreateEnumerationStatement(c *CreateEnumerationStatementContext)

	// ExitEnumValueList is called when exiting the enumValueList production.
	ExitEnumValueList(c *EnumValueListContext)

	// ExitEnumValue is called when exiting the enumValue production.
	ExitEnumValue(c *EnumValueContext)

	// ExitAlterEnumerationStatement is called when exiting the alterEnumerationStatement production.
	ExitAlterEnumerationStatement(c *AlterEnumerationStatementContext)

	// ExitAddEnumValue is called when exiting the AddEnumValue production.
	ExitAddEnumValue(c *AddEnumValueContext)

	// ExitDropEnumValue is called when exiting the DropEnumValue production.
	ExitDropEnumValue(c *DropEnumValueContext)

	// ExitRenameEnumValue is called when exiting the RenameEnumValue production.
	ExitRenameEnumValue(c *RenameEnumValueContext)

	// ExitDropEnumerationStatement is called when exiting the dropEnumerationStatement production.
	ExitDropEnumerationStatement(c *DropEnumerationStatementContext)

	// ExitCreateEntityStatement is called when exiting the createEntityStatement production.
	ExitCreateEntityStatement(c *CreateEntityStatementContext)

	// ExitCreateOrModify is called when exiting the createOrModify production.
	ExitCreateOrModify(c *CreateOrModifyContext)

	// ExitEntityType is called when exiting the entityType production.
	ExitEntityType(c *EntityTypeContext)

	// ExitAttributeList is called when exiting the attributeList production.
	ExitAttributeList(c *AttributeListContext)

	// ExitAttribute is called when exiting the attribute production.
	ExitAttribute(c *AttributeContext)

	// ExitAttributeName is called when exiting the attributeName production.
	ExitAttributeName(c *AttributeNameContext)

	// ExitRenamedFromAnnotation is called when exiting the renamedFromAnnotation production.
	ExitRenamedFromAnnotation(c *RenamedFromAnnotationContext)

	// ExitStringType is called when exiting the StringType production.
	ExitStringType(c *StringTypeContext)

	// ExitIntegerType is called when exiting the IntegerType production.
	ExitIntegerType(c *IntegerTypeContext)

	// ExitLongType is called when exiting the LongType production.
	ExitLongType(c *LongTypeContext)

	// ExitDecimalType is called when exiting the DecimalType production.
	ExitDecimalType(c *DecimalTypeContext)

	// ExitBooleanType is called when exiting the BooleanType production.
	ExitBooleanType(c *BooleanTypeContext)

	// ExitDateTimeType is called when exiting the DateTimeType production.
	ExitDateTimeType(c *DateTimeTypeContext)

	// ExitDateType is called when exiting the DateType production.
	ExitDateType(c *DateTypeContext)

	// ExitAutoNumberType is called when exiting the AutoNumberType production.
	ExitAutoNumberType(c *AutoNumberTypeContext)

	// ExitBinaryType is called when exiting the BinaryType production.
	ExitBinaryType(c *BinaryTypeContext)

	// ExitEnumerationType is called when exiting the EnumerationType production.
	ExitEnumerationType(c *EnumerationTypeContext)

	// ExitAttributeConstraints is called when exiting the attributeConstraints production.
	ExitAttributeConstraints(c *AttributeConstraintsContext)

	// ExitNotNullConstraint is called when exiting the NotNullConstraint production.
	ExitNotNullConstraint(c *NotNullConstraintContext)

	// ExitUniqueConstraint is called when exiting the UniqueConstraint production.
	ExitUniqueConstraint(c *UniqueConstraintContext)

	// ExitDefaultConstraint is called when exiting the DefaultConstraint production.
	ExitDefaultConstraint(c *DefaultConstraintContext)

	// ExitDefaultValue is called when exiting the defaultValue production.
	ExitDefaultValue(c *DefaultValueContext)

	// ExitIndexClause is called when exiting the indexClause production.
	ExitIndexClause(c *IndexClauseContext)

	// ExitIndexColumnList is called when exiting the indexColumnList production.
	ExitIndexColumnList(c *IndexColumnListContext)

	// ExitIndexColumn is called when exiting the indexColumn production.
	ExitIndexColumn(c *IndexColumnContext)

	// ExitDropEntityStatement is called when exiting the dropEntityStatement production.
	ExitDropEntityStatement(c *DropEntityStatementContext)

	// ExitCreateViewEntityStatement is called when exiting the createViewEntityStatement production.
	ExitCreateViewEntityStatement(c *CreateViewEntityStatementContext)

	// ExitViewAttributeList is called when exiting the viewAttributeList production.
	ExitViewAttributeList(c *ViewAttributeListContext)

	// ExitViewAttribute is called when exiting the viewAttribute production.
	ExitViewAttribute(c *ViewAttributeContext)

	// ExitOqlQuery is called when exiting the oqlQuery production.
	ExitOqlQuery(c *OqlQueryContext)

	// ExitSelectClause is called when exiting the selectClause production.
	ExitSelectClause(c *SelectClauseContext)

	// ExitSelectItems is called when exiting the selectItems production.
	ExitSelectItems(c *SelectItemsContext)

	// ExitSelectItem is called when exiting the selectItem production.
	ExitSelectItem(c *SelectItemContext)

	// ExitFromClause is called when exiting the fromClause production.
	ExitFromClause(c *FromClauseContext)

	// ExitFromItem is called when exiting the fromItem production.
	ExitFromItem(c *FromItemContext)

	// ExitJoinClause is called when exiting the joinClause production.
	ExitJoinClause(c *JoinClauseContext)

	// ExitSimpleJoinTarget is called when exiting the SimpleJoinTarget production.
	ExitSimpleJoinTarget(c *SimpleJoinTargetContext)

	// ExitAssociationJoinTarget is called when exiting the AssociationJoinTarget production.
	ExitAssociationJoinTarget(c *AssociationJoinTargetContext)

	// ExitWhereClause is called when exiting the whereClause production.
	ExitWhereClause(c *WhereClauseContext)

	// ExitGroupByClause is called when exiting the groupByClause production.
	ExitGroupByClause(c *GroupByClauseContext)

	// ExitOrderByClause is called when exiting the orderByClause production.
	ExitOrderByClause(c *OrderByClauseContext)

	// ExitOrderByItem is called when exiting the orderByItem production.
	ExitOrderByItem(c *OrderByItemContext)

	// ExitLimitClause is called when exiting the limitClause production.
	ExitLimitClause(c *LimitClauseContext)

	// ExitAndExpr is called when exiting the AndExpr production.
	ExitAndExpr(c *AndExprContext)

	// ExitStringExpr is called when exiting the StringExpr production.
	ExitStringExpr(c *StringExprContext)

	// ExitIdentExpr is called when exiting the IdentExpr production.
	ExitIdentExpr(c *IdentExprContext)

	// ExitTrueExpr is called when exiting the TrueExpr production.
	ExitTrueExpr(c *TrueExprContext)

	// ExitIsNullExpr is called when exiting the IsNullExpr production.
	ExitIsNullExpr(c *IsNullExprContext)

	// ExitStarExpr is called when exiting the StarExpr production.
	ExitStarExpr(c *StarExprContext)

	// ExitFuncExpr is called when exiting the FuncExpr production.
	ExitFuncExpr(c *FuncExprContext)

	// ExitQualifiedExpr is called when exiting the QualifiedExpr production.
	ExitQualifiedExpr(c *QualifiedExprContext)

	// ExitDecimalExpr is called when exiting the DecimalExpr production.
	ExitDecimalExpr(c *DecimalExprContext)

	// ExitOrExpr is called when exiting the OrExpr production.
	ExitOrExpr(c *OrExprContext)

	// ExitFalseExpr is called when exiting the FalseExpr production.
	ExitFalseExpr(c *FalseExprContext)

	// ExitInSubqueryExpr is called when exiting the InSubqueryExpr production.
	ExitInSubqueryExpr(c *InSubqueryExprContext)

	// ExitMulDivExpr is called when exiting the MulDivExpr production.
	ExitMulDivExpr(c *MulDivExprContext)

	// ExitDivisionExpr is called when exiting the DivisionExpr production.
	ExitDivisionExpr(c *DivisionExprContext)

	// ExitCompareExpr is called when exiting the CompareExpr production.
	ExitCompareExpr(c *CompareExprContext)

	// ExitFieldAccessExpr is called when exiting the FieldAccessExpr production.
	ExitFieldAccessExpr(c *FieldAccessExprContext)

	// ExitNotExpr is called when exiting the NotExpr production.
	ExitNotExpr(c *NotExprContext)

	// ExitIntExpr is called when exiting the IntExpr production.
	ExitIntExpr(c *IntExprContext)

	// ExitSysVarExpr is called when exiting the SysVarExpr production.
	ExitSysVarExpr(c *SysVarExprContext)

	// ExitInExpr is called when exiting the InExpr production.
	ExitInExpr(c *InExprContext)

	// ExitParenExpr is called when exiting the ParenExpr production.
	ExitParenExpr(c *ParenExprContext)

	// ExitCaseExpr is called when exiting the CaseExpr production.
	ExitCaseExpr(c *CaseExprContext)

	// ExitAddSubExpr is called when exiting the AddSubExpr production.
	ExitAddSubExpr(c *AddSubExprContext)

	// ExitSubqueryExpr is called when exiting the SubqueryExpr production.
	ExitSubqueryExpr(c *SubqueryExprContext)

	// ExitComparisonOp is called when exiting the comparisonOp production.
	ExitComparisonOp(c *ComparisonOpContext)

	// ExitFunctionCall is called when exiting the functionCall production.
	ExitFunctionCall(c *FunctionCallContext)

	// ExitAggregateFunction is called when exiting the aggregateFunction production.
	ExitAggregateFunction(c *AggregateFunctionContext)

	// ExitCaseExpression is called when exiting the caseExpression production.
	ExitCaseExpression(c *CaseExpressionContext)

	// ExitExpressionList is called when exiting the expressionList production.
	ExitExpressionList(c *ExpressionListContext)

	// ExitSystemVariable is called when exiting the systemVariable production.
	ExitSystemVariable(c *SystemVariableContext)

	// ExitCreateAssociationStatement is called when exiting the createAssociationStatement production.
	ExitCreateAssociationStatement(c *CreateAssociationStatementContext)

	// ExitAssociationType is called when exiting the associationType production.
	ExitAssociationType(c *AssociationTypeContext)

	// ExitOwnerType is called when exiting the ownerType production.
	ExitOwnerType(c *OwnerTypeContext)

	// ExitDeleteBehavior is called when exiting the deleteBehavior production.
	ExitDeleteBehavior(c *DeleteBehaviorContext)

	// ExitDropAssociationStatement is called when exiting the dropAssociationStatement production.
	ExitDropAssociationStatement(c *DropAssociationStatementContext)

	// ExitQueryStatement is called when exiting the queryStatement production.
	ExitQueryStatement(c *QueryStatementContext)

	// ExitShowModules is called when exiting the ShowModules production.
	ExitShowModules(c *ShowModulesContext)

	// ExitShowEnumerations is called when exiting the ShowEnumerations production.
	ExitShowEnumerations(c *ShowEnumerationsContext)

	// ExitShowEntities is called when exiting the ShowEntities production.
	ExitShowEntities(c *ShowEntitiesContext)

	// ExitShowEntity is called when exiting the ShowEntity production.
	ExitShowEntity(c *ShowEntityContext)

	// ExitShowAssociations is called when exiting the ShowAssociations production.
	ExitShowAssociations(c *ShowAssociationsContext)

	// ExitShowAssociation is called when exiting the ShowAssociation production.
	ExitShowAssociation(c *ShowAssociationContext)

	// ExitDescribeEnumeration is called when exiting the DescribeEnumeration production.
	ExitDescribeEnumeration(c *DescribeEnumerationContext)

	// ExitDescribeEntity is called when exiting the DescribeEntity production.
	ExitDescribeEntity(c *DescribeEntityContext)

	// ExitDescribeAssociation is called when exiting the DescribeAssociation production.
	ExitDescribeAssociation(c *DescribeAssociationContext)

	// ExitRepositoryStatement is called when exiting the repositoryStatement production.
	ExitRepositoryStatement(c *RepositoryStatementContext)

	// ExitCommitStatement is called when exiting the commitStatement production.
	ExitCommitStatement(c *CommitStatementContext)

	// ExitUpdateStatement is called when exiting the updateStatement production.
	ExitUpdateStatement(c *UpdateStatementContext)

	// ExitRefreshStatement is called when exiting the refreshStatement production.
	ExitRefreshStatement(c *RefreshStatementContext)

	// ExitSessionStatement is called when exiting the sessionStatement production.
	ExitSessionStatement(c *SessionStatementContext)

	// ExitSetStatement is called when exiting the setStatement production.
	ExitSetStatement(c *SetStatementContext)

	// ExitHelpStatement is called when exiting the helpStatement production.
	ExitHelpStatement(c *HelpStatementContext)

	// ExitExitStatement is called when exiting the exitStatement production.
	ExitExitStatement(c *ExitStatementContext)

	// ExitExecuteScriptStatement is called when exiting the executeScriptStatement production.
	ExitExecuteScriptStatement(c *ExecuteScriptStatementContext)

	// ExitQualifiedName is called when exiting the qualifiedName production.
	ExitQualifiedName(c *QualifiedNameContext)

	// ExitDocumentation is called when exiting the documentation production.
	ExitDocumentation(c *DocumentationContext)

	// ExitPositionAnnotation is called when exiting the positionAnnotation production.
	ExitPositionAnnotation(c *PositionAnnotationContext)

	// ExitStringLiteral is called when exiting the stringLiteral production.
	ExitStringLiteral(c *StringLiteralContext)
}
