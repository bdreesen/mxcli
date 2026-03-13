// SPDX-License-Identifier: Apache-2.0

// Code generated from MDL.g4 by ANTLR 4.7.2. DO NOT EDIT.

package parser // MDL
import "github.com/antlr/antlr4/runtime/Go/antlr"

// BaseMDLListener is a complete listener for a parse tree produced by MDLParser.
type BaseMDLListener struct{}

var _ MDLListener = &BaseMDLListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseMDLListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseMDLListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseMDLListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseMDLListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterProgram is called when production program is entered.
func (s *BaseMDLListener) EnterProgram(ctx *ProgramContext) {}

// ExitProgram is called when production program is exited.
func (s *BaseMDLListener) ExitProgram(ctx *ProgramContext) {}

// EnterStatement is called when production statement is entered.
func (s *BaseMDLListener) EnterStatement(ctx *StatementContext) {}

// ExitStatement is called when production statement is exited.
func (s *BaseMDLListener) ExitStatement(ctx *StatementContext) {}

// EnterTerminator is called when production terminator is entered.
func (s *BaseMDLListener) EnterTerminator(ctx *TerminatorContext) {}

// ExitTerminator is called when production terminator is exited.
func (s *BaseMDLListener) ExitTerminator(ctx *TerminatorContext) {}

// EnterConnectionStatement is called when production connectionStatement is entered.
func (s *BaseMDLListener) EnterConnectionStatement(ctx *ConnectionStatementContext) {}

// ExitConnectionStatement is called when production connectionStatement is exited.
func (s *BaseMDLListener) ExitConnectionStatement(ctx *ConnectionStatementContext) {}

// EnterConnectLocal is called when production ConnectLocal is entered.
func (s *BaseMDLListener) EnterConnectLocal(ctx *ConnectLocalContext) {}

// ExitConnectLocal is called when production ConnectLocal is exited.
func (s *BaseMDLListener) ExitConnectLocal(ctx *ConnectLocalContext) {}

// EnterConnectFilesystem is called when production ConnectFilesystem is entered.
func (s *BaseMDLListener) EnterConnectFilesystem(ctx *ConnectFilesystemContext) {}

// ExitConnectFilesystem is called when production ConnectFilesystem is exited.
func (s *BaseMDLListener) ExitConnectFilesystem(ctx *ConnectFilesystemContext) {}

// EnterDisconnectStatement is called when production disconnectStatement is entered.
func (s *BaseMDLListener) EnterDisconnectStatement(ctx *DisconnectStatementContext) {}

// ExitDisconnectStatement is called when production disconnectStatement is exited.
func (s *BaseMDLListener) ExitDisconnectStatement(ctx *DisconnectStatementContext) {}

// EnterStatusStatement is called when production statusStatement is entered.
func (s *BaseMDLListener) EnterStatusStatement(ctx *StatusStatementContext) {}

// ExitStatusStatement is called when production statusStatement is exited.
func (s *BaseMDLListener) ExitStatusStatement(ctx *StatusStatementContext) {}

// EnterDdlStatement is called when production ddlStatement is entered.
func (s *BaseMDLListener) EnterDdlStatement(ctx *DdlStatementContext) {}

// ExitDdlStatement is called when production ddlStatement is exited.
func (s *BaseMDLListener) ExitDdlStatement(ctx *DdlStatementContext) {}

// EnterCreateStatement is called when production createStatement is entered.
func (s *BaseMDLListener) EnterCreateStatement(ctx *CreateStatementContext) {}

// ExitCreateStatement is called when production createStatement is exited.
func (s *BaseMDLListener) ExitCreateStatement(ctx *CreateStatementContext) {}

// EnterCreateModuleStatement is called when production createModuleStatement is entered.
func (s *BaseMDLListener) EnterCreateModuleStatement(ctx *CreateModuleStatementContext) {}

// ExitCreateModuleStatement is called when production createModuleStatement is exited.
func (s *BaseMDLListener) ExitCreateModuleStatement(ctx *CreateModuleStatementContext) {}

// EnterAlterStatement is called when production alterStatement is entered.
func (s *BaseMDLListener) EnterAlterStatement(ctx *AlterStatementContext) {}

// ExitAlterStatement is called when production alterStatement is exited.
func (s *BaseMDLListener) ExitAlterStatement(ctx *AlterStatementContext) {}

// EnterDropStatement is called when production dropStatement is entered.
func (s *BaseMDLListener) EnterDropStatement(ctx *DropStatementContext) {}

// ExitDropStatement is called when production dropStatement is exited.
func (s *BaseMDLListener) ExitDropStatement(ctx *DropStatementContext) {}

// EnterDropModuleStatement is called when production dropModuleStatement is entered.
func (s *BaseMDLListener) EnterDropModuleStatement(ctx *DropModuleStatementContext) {}

// ExitDropModuleStatement is called when production dropModuleStatement is exited.
func (s *BaseMDLListener) ExitDropModuleStatement(ctx *DropModuleStatementContext) {}

// EnterCreateEnumerationStatement is called when production createEnumerationStatement is entered.
func (s *BaseMDLListener) EnterCreateEnumerationStatement(ctx *CreateEnumerationStatementContext) {}

// ExitCreateEnumerationStatement is called when production createEnumerationStatement is exited.
func (s *BaseMDLListener) ExitCreateEnumerationStatement(ctx *CreateEnumerationStatementContext) {}

// EnterEnumValueList is called when production enumValueList is entered.
func (s *BaseMDLListener) EnterEnumValueList(ctx *EnumValueListContext) {}

// ExitEnumValueList is called when production enumValueList is exited.
func (s *BaseMDLListener) ExitEnumValueList(ctx *EnumValueListContext) {}

// EnterEnumValue is called when production enumValue is entered.
func (s *BaseMDLListener) EnterEnumValue(ctx *EnumValueContext) {}

// ExitEnumValue is called when production enumValue is exited.
func (s *BaseMDLListener) ExitEnumValue(ctx *EnumValueContext) {}

// EnterAlterEnumerationStatement is called when production alterEnumerationStatement is entered.
func (s *BaseMDLListener) EnterAlterEnumerationStatement(ctx *AlterEnumerationStatementContext) {}

// ExitAlterEnumerationStatement is called when production alterEnumerationStatement is exited.
func (s *BaseMDLListener) ExitAlterEnumerationStatement(ctx *AlterEnumerationStatementContext) {}

// EnterAddEnumValue is called when production AddEnumValue is entered.
func (s *BaseMDLListener) EnterAddEnumValue(ctx *AddEnumValueContext) {}

// ExitAddEnumValue is called when production AddEnumValue is exited.
func (s *BaseMDLListener) ExitAddEnumValue(ctx *AddEnumValueContext) {}

// EnterDropEnumValue is called when production DropEnumValue is entered.
func (s *BaseMDLListener) EnterDropEnumValue(ctx *DropEnumValueContext) {}

// ExitDropEnumValue is called when production DropEnumValue is exited.
func (s *BaseMDLListener) ExitDropEnumValue(ctx *DropEnumValueContext) {}

// EnterRenameEnumValue is called when production RenameEnumValue is entered.
func (s *BaseMDLListener) EnterRenameEnumValue(ctx *RenameEnumValueContext) {}

// ExitRenameEnumValue is called when production RenameEnumValue is exited.
func (s *BaseMDLListener) ExitRenameEnumValue(ctx *RenameEnumValueContext) {}

// EnterDropEnumerationStatement is called when production dropEnumerationStatement is entered.
func (s *BaseMDLListener) EnterDropEnumerationStatement(ctx *DropEnumerationStatementContext) {}

// ExitDropEnumerationStatement is called when production dropEnumerationStatement is exited.
func (s *BaseMDLListener) ExitDropEnumerationStatement(ctx *DropEnumerationStatementContext) {}

// EnterCreateEntityStatement is called when production createEntityStatement is entered.
func (s *BaseMDLListener) EnterCreateEntityStatement(ctx *CreateEntityStatementContext) {}

// ExitCreateEntityStatement is called when production createEntityStatement is exited.
func (s *BaseMDLListener) ExitCreateEntityStatement(ctx *CreateEntityStatementContext) {}

// EnterCreateOrModify is called when production createOrModify is entered.
func (s *BaseMDLListener) EnterCreateOrModify(ctx *CreateOrModifyContext) {}

// ExitCreateOrModify is called when production createOrModify is exited.
func (s *BaseMDLListener) ExitCreateOrModify(ctx *CreateOrModifyContext) {}

// EnterEntityType is called when production entityType is entered.
func (s *BaseMDLListener) EnterEntityType(ctx *EntityTypeContext) {}

// ExitEntityType is called when production entityType is exited.
func (s *BaseMDLListener) ExitEntityType(ctx *EntityTypeContext) {}

// EnterAttributeList is called when production attributeList is entered.
func (s *BaseMDLListener) EnterAttributeList(ctx *AttributeListContext) {}

// ExitAttributeList is called when production attributeList is exited.
func (s *BaseMDLListener) ExitAttributeList(ctx *AttributeListContext) {}

// EnterAttribute is called when production attribute is entered.
func (s *BaseMDLListener) EnterAttribute(ctx *AttributeContext) {}

// ExitAttribute is called when production attribute is exited.
func (s *BaseMDLListener) ExitAttribute(ctx *AttributeContext) {}

// EnterAttributeName is called when production attributeName is entered.
func (s *BaseMDLListener) EnterAttributeName(ctx *AttributeNameContext) {}

// ExitAttributeName is called when production attributeName is exited.
func (s *BaseMDLListener) ExitAttributeName(ctx *AttributeNameContext) {}

// EnterRenamedFromAnnotation is called when production renamedFromAnnotation is entered.
func (s *BaseMDLListener) EnterRenamedFromAnnotation(ctx *RenamedFromAnnotationContext) {}

// ExitRenamedFromAnnotation is called when production renamedFromAnnotation is exited.
func (s *BaseMDLListener) ExitRenamedFromAnnotation(ctx *RenamedFromAnnotationContext) {}

// EnterStringType is called when production StringType is entered.
func (s *BaseMDLListener) EnterStringType(ctx *StringTypeContext) {}

// ExitStringType is called when production StringType is exited.
func (s *BaseMDLListener) ExitStringType(ctx *StringTypeContext) {}

// EnterIntegerType is called when production IntegerType is entered.
func (s *BaseMDLListener) EnterIntegerType(ctx *IntegerTypeContext) {}

// ExitIntegerType is called when production IntegerType is exited.
func (s *BaseMDLListener) ExitIntegerType(ctx *IntegerTypeContext) {}

// EnterLongType is called when production LongType is entered.
func (s *BaseMDLListener) EnterLongType(ctx *LongTypeContext) {}

// ExitLongType is called when production LongType is exited.
func (s *BaseMDLListener) ExitLongType(ctx *LongTypeContext) {}

// EnterDecimalType is called when production DecimalType is entered.
func (s *BaseMDLListener) EnterDecimalType(ctx *DecimalTypeContext) {}

// ExitDecimalType is called when production DecimalType is exited.
func (s *BaseMDLListener) ExitDecimalType(ctx *DecimalTypeContext) {}

// EnterBooleanType is called when production BooleanType is entered.
func (s *BaseMDLListener) EnterBooleanType(ctx *BooleanTypeContext) {}

// ExitBooleanType is called when production BooleanType is exited.
func (s *BaseMDLListener) ExitBooleanType(ctx *BooleanTypeContext) {}

// EnterDateTimeType is called when production DateTimeType is entered.
func (s *BaseMDLListener) EnterDateTimeType(ctx *DateTimeTypeContext) {}

// ExitDateTimeType is called when production DateTimeType is exited.
func (s *BaseMDLListener) ExitDateTimeType(ctx *DateTimeTypeContext) {}

// EnterDateType is called when production DateType is entered.
func (s *BaseMDLListener) EnterDateType(ctx *DateTypeContext) {}

// ExitDateType is called when production DateType is exited.
func (s *BaseMDLListener) ExitDateType(ctx *DateTypeContext) {}

// EnterAutoNumberType is called when production AutoNumberType is entered.
func (s *BaseMDLListener) EnterAutoNumberType(ctx *AutoNumberTypeContext) {}

// ExitAutoNumberType is called when production AutoNumberType is exited.
func (s *BaseMDLListener) ExitAutoNumberType(ctx *AutoNumberTypeContext) {}

// EnterBinaryType is called when production BinaryType is entered.
func (s *BaseMDLListener) EnterBinaryType(ctx *BinaryTypeContext) {}

// ExitBinaryType is called when production BinaryType is exited.
func (s *BaseMDLListener) ExitBinaryType(ctx *BinaryTypeContext) {}

// EnterEnumerationType is called when production EnumerationType is entered.
func (s *BaseMDLListener) EnterEnumerationType(ctx *EnumerationTypeContext) {}

// ExitEnumerationType is called when production EnumerationType is exited.
func (s *BaseMDLListener) ExitEnumerationType(ctx *EnumerationTypeContext) {}

// EnterAttributeConstraints is called when production attributeConstraints is entered.
func (s *BaseMDLListener) EnterAttributeConstraints(ctx *AttributeConstraintsContext) {}

// ExitAttributeConstraints is called when production attributeConstraints is exited.
func (s *BaseMDLListener) ExitAttributeConstraints(ctx *AttributeConstraintsContext) {}

// EnterNotNullConstraint is called when production NotNullConstraint is entered.
func (s *BaseMDLListener) EnterNotNullConstraint(ctx *NotNullConstraintContext) {}

// ExitNotNullConstraint is called when production NotNullConstraint is exited.
func (s *BaseMDLListener) ExitNotNullConstraint(ctx *NotNullConstraintContext) {}

// EnterUniqueConstraint is called when production UniqueConstraint is entered.
func (s *BaseMDLListener) EnterUniqueConstraint(ctx *UniqueConstraintContext) {}

// ExitUniqueConstraint is called when production UniqueConstraint is exited.
func (s *BaseMDLListener) ExitUniqueConstraint(ctx *UniqueConstraintContext) {}

// EnterDefaultConstraint is called when production DefaultConstraint is entered.
func (s *BaseMDLListener) EnterDefaultConstraint(ctx *DefaultConstraintContext) {}

// ExitDefaultConstraint is called when production DefaultConstraint is exited.
func (s *BaseMDLListener) ExitDefaultConstraint(ctx *DefaultConstraintContext) {}

// EnterDefaultValue is called when production defaultValue is entered.
func (s *BaseMDLListener) EnterDefaultValue(ctx *DefaultValueContext) {}

// ExitDefaultValue is called when production defaultValue is exited.
func (s *BaseMDLListener) ExitDefaultValue(ctx *DefaultValueContext) {}

// EnterIndexClause is called when production indexClause is entered.
func (s *BaseMDLListener) EnterIndexClause(ctx *IndexClauseContext) {}

// ExitIndexClause is called when production indexClause is exited.
func (s *BaseMDLListener) ExitIndexClause(ctx *IndexClauseContext) {}

// EnterIndexColumnList is called when production indexColumnList is entered.
func (s *BaseMDLListener) EnterIndexColumnList(ctx *IndexColumnListContext) {}

// ExitIndexColumnList is called when production indexColumnList is exited.
func (s *BaseMDLListener) ExitIndexColumnList(ctx *IndexColumnListContext) {}

// EnterIndexColumn is called when production indexColumn is entered.
func (s *BaseMDLListener) EnterIndexColumn(ctx *IndexColumnContext) {}

// ExitIndexColumn is called when production indexColumn is exited.
func (s *BaseMDLListener) ExitIndexColumn(ctx *IndexColumnContext) {}

// EnterDropEntityStatement is called when production dropEntityStatement is entered.
func (s *BaseMDLListener) EnterDropEntityStatement(ctx *DropEntityStatementContext) {}

// ExitDropEntityStatement is called when production dropEntityStatement is exited.
func (s *BaseMDLListener) ExitDropEntityStatement(ctx *DropEntityStatementContext) {}

// EnterCreateViewEntityStatement is called when production createViewEntityStatement is entered.
func (s *BaseMDLListener) EnterCreateViewEntityStatement(ctx *CreateViewEntityStatementContext) {}

// ExitCreateViewEntityStatement is called when production createViewEntityStatement is exited.
func (s *BaseMDLListener) ExitCreateViewEntityStatement(ctx *CreateViewEntityStatementContext) {}

// EnterViewAttributeList is called when production viewAttributeList is entered.
func (s *BaseMDLListener) EnterViewAttributeList(ctx *ViewAttributeListContext) {}

// ExitViewAttributeList is called when production viewAttributeList is exited.
func (s *BaseMDLListener) ExitViewAttributeList(ctx *ViewAttributeListContext) {}

// EnterViewAttribute is called when production viewAttribute is entered.
func (s *BaseMDLListener) EnterViewAttribute(ctx *ViewAttributeContext) {}

// ExitViewAttribute is called when production viewAttribute is exited.
func (s *BaseMDLListener) ExitViewAttribute(ctx *ViewAttributeContext) {}

// EnterOqlQuery is called when production oqlQuery is entered.
func (s *BaseMDLListener) EnterOqlQuery(ctx *OqlQueryContext) {}

// ExitOqlQuery is called when production oqlQuery is exited.
func (s *BaseMDLListener) ExitOqlQuery(ctx *OqlQueryContext) {}

// EnterSelectClause is called when production selectClause is entered.
func (s *BaseMDLListener) EnterSelectClause(ctx *SelectClauseContext) {}

// ExitSelectClause is called when production selectClause is exited.
func (s *BaseMDLListener) ExitSelectClause(ctx *SelectClauseContext) {}

// EnterSelectItems is called when production selectItems is entered.
func (s *BaseMDLListener) EnterSelectItems(ctx *SelectItemsContext) {}

// ExitSelectItems is called when production selectItems is exited.
func (s *BaseMDLListener) ExitSelectItems(ctx *SelectItemsContext) {}

// EnterSelectItem is called when production selectItem is entered.
func (s *BaseMDLListener) EnterSelectItem(ctx *SelectItemContext) {}

// ExitSelectItem is called when production selectItem is exited.
func (s *BaseMDLListener) ExitSelectItem(ctx *SelectItemContext) {}

// EnterFromClause is called when production fromClause is entered.
func (s *BaseMDLListener) EnterFromClause(ctx *FromClauseContext) {}

// ExitFromClause is called when production fromClause is exited.
func (s *BaseMDLListener) ExitFromClause(ctx *FromClauseContext) {}

// EnterFromItem is called when production fromItem is entered.
func (s *BaseMDLListener) EnterFromItem(ctx *FromItemContext) {}

// ExitFromItem is called when production fromItem is exited.
func (s *BaseMDLListener) ExitFromItem(ctx *FromItemContext) {}

// EnterJoinClause is called when production joinClause is entered.
func (s *BaseMDLListener) EnterJoinClause(ctx *JoinClauseContext) {}

// ExitJoinClause is called when production joinClause is exited.
func (s *BaseMDLListener) ExitJoinClause(ctx *JoinClauseContext) {}

// EnterSimpleJoinTarget is called when production SimpleJoinTarget is entered.
func (s *BaseMDLListener) EnterSimpleJoinTarget(ctx *SimpleJoinTargetContext) {}

// ExitSimpleJoinTarget is called when production SimpleJoinTarget is exited.
func (s *BaseMDLListener) ExitSimpleJoinTarget(ctx *SimpleJoinTargetContext) {}

// EnterAssociationJoinTarget is called when production AssociationJoinTarget is entered.
func (s *BaseMDLListener) EnterAssociationJoinTarget(ctx *AssociationJoinTargetContext) {}

// ExitAssociationJoinTarget is called when production AssociationJoinTarget is exited.
func (s *BaseMDLListener) ExitAssociationJoinTarget(ctx *AssociationJoinTargetContext) {}

// EnterWhereClause is called when production whereClause is entered.
func (s *BaseMDLListener) EnterWhereClause(ctx *WhereClauseContext) {}

// ExitWhereClause is called when production whereClause is exited.
func (s *BaseMDLListener) ExitWhereClause(ctx *WhereClauseContext) {}

// EnterGroupByClause is called when production groupByClause is entered.
func (s *BaseMDLListener) EnterGroupByClause(ctx *GroupByClauseContext) {}

// ExitGroupByClause is called when production groupByClause is exited.
func (s *BaseMDLListener) ExitGroupByClause(ctx *GroupByClauseContext) {}

// EnterOrderByClause is called when production orderByClause is entered.
func (s *BaseMDLListener) EnterOrderByClause(ctx *OrderByClauseContext) {}

// ExitOrderByClause is called when production orderByClause is exited.
func (s *BaseMDLListener) ExitOrderByClause(ctx *OrderByClauseContext) {}

// EnterOrderByItem is called when production orderByItem is entered.
func (s *BaseMDLListener) EnterOrderByItem(ctx *OrderByItemContext) {}

// ExitOrderByItem is called when production orderByItem is exited.
func (s *BaseMDLListener) ExitOrderByItem(ctx *OrderByItemContext) {}

// EnterLimitClause is called when production limitClause is entered.
func (s *BaseMDLListener) EnterLimitClause(ctx *LimitClauseContext) {}

// ExitLimitClause is called when production limitClause is exited.
func (s *BaseMDLListener) ExitLimitClause(ctx *LimitClauseContext) {}

// EnterAndExpr is called when production AndExpr is entered.
func (s *BaseMDLListener) EnterAndExpr(ctx *AndExprContext) {}

// ExitAndExpr is called when production AndExpr is exited.
func (s *BaseMDLListener) ExitAndExpr(ctx *AndExprContext) {}

// EnterStringExpr is called when production StringExpr is entered.
func (s *BaseMDLListener) EnterStringExpr(ctx *StringExprContext) {}

// ExitStringExpr is called when production StringExpr is exited.
func (s *BaseMDLListener) ExitStringExpr(ctx *StringExprContext) {}

// EnterIdentExpr is called when production IdentExpr is entered.
func (s *BaseMDLListener) EnterIdentExpr(ctx *IdentExprContext) {}

// ExitIdentExpr is called when production IdentExpr is exited.
func (s *BaseMDLListener) ExitIdentExpr(ctx *IdentExprContext) {}

// EnterTrueExpr is called when production TrueExpr is entered.
func (s *BaseMDLListener) EnterTrueExpr(ctx *TrueExprContext) {}

// ExitTrueExpr is called when production TrueExpr is exited.
func (s *BaseMDLListener) ExitTrueExpr(ctx *TrueExprContext) {}

// EnterIsNullExpr is called when production IsNullExpr is entered.
func (s *BaseMDLListener) EnterIsNullExpr(ctx *IsNullExprContext) {}

// ExitIsNullExpr is called when production IsNullExpr is exited.
func (s *BaseMDLListener) ExitIsNullExpr(ctx *IsNullExprContext) {}

// EnterStarExpr is called when production StarExpr is entered.
func (s *BaseMDLListener) EnterStarExpr(ctx *StarExprContext) {}

// ExitStarExpr is called when production StarExpr is exited.
func (s *BaseMDLListener) ExitStarExpr(ctx *StarExprContext) {}

// EnterFuncExpr is called when production FuncExpr is entered.
func (s *BaseMDLListener) EnterFuncExpr(ctx *FuncExprContext) {}

// ExitFuncExpr is called when production FuncExpr is exited.
func (s *BaseMDLListener) ExitFuncExpr(ctx *FuncExprContext) {}

// EnterQualifiedExpr is called when production QualifiedExpr is entered.
func (s *BaseMDLListener) EnterQualifiedExpr(ctx *QualifiedExprContext) {}

// ExitQualifiedExpr is called when production QualifiedExpr is exited.
func (s *BaseMDLListener) ExitQualifiedExpr(ctx *QualifiedExprContext) {}

// EnterDecimalExpr is called when production DecimalExpr is entered.
func (s *BaseMDLListener) EnterDecimalExpr(ctx *DecimalExprContext) {}

// ExitDecimalExpr is called when production DecimalExpr is exited.
func (s *BaseMDLListener) ExitDecimalExpr(ctx *DecimalExprContext) {}

// EnterOrExpr is called when production OrExpr is entered.
func (s *BaseMDLListener) EnterOrExpr(ctx *OrExprContext) {}

// ExitOrExpr is called when production OrExpr is exited.
func (s *BaseMDLListener) ExitOrExpr(ctx *OrExprContext) {}

// EnterFalseExpr is called when production FalseExpr is entered.
func (s *BaseMDLListener) EnterFalseExpr(ctx *FalseExprContext) {}

// ExitFalseExpr is called when production FalseExpr is exited.
func (s *BaseMDLListener) ExitFalseExpr(ctx *FalseExprContext) {}

// EnterInSubqueryExpr is called when production InSubqueryExpr is entered.
func (s *BaseMDLListener) EnterInSubqueryExpr(ctx *InSubqueryExprContext) {}

// ExitInSubqueryExpr is called when production InSubqueryExpr is exited.
func (s *BaseMDLListener) ExitInSubqueryExpr(ctx *InSubqueryExprContext) {}

// EnterMulDivExpr is called when production MulDivExpr is entered.
func (s *BaseMDLListener) EnterMulDivExpr(ctx *MulDivExprContext) {}

// ExitMulDivExpr is called when production MulDivExpr is exited.
func (s *BaseMDLListener) ExitMulDivExpr(ctx *MulDivExprContext) {}

// EnterDivisionExpr is called when production DivisionExpr is entered.
func (s *BaseMDLListener) EnterDivisionExpr(ctx *DivisionExprContext) {}

// ExitDivisionExpr is called when production DivisionExpr is exited.
func (s *BaseMDLListener) ExitDivisionExpr(ctx *DivisionExprContext) {}

// EnterCompareExpr is called when production CompareExpr is entered.
func (s *BaseMDLListener) EnterCompareExpr(ctx *CompareExprContext) {}

// ExitCompareExpr is called when production CompareExpr is exited.
func (s *BaseMDLListener) ExitCompareExpr(ctx *CompareExprContext) {}

// EnterFieldAccessExpr is called when production FieldAccessExpr is entered.
func (s *BaseMDLListener) EnterFieldAccessExpr(ctx *FieldAccessExprContext) {}

// ExitFieldAccessExpr is called when production FieldAccessExpr is exited.
func (s *BaseMDLListener) ExitFieldAccessExpr(ctx *FieldAccessExprContext) {}

// EnterNotExpr is called when production NotExpr is entered.
func (s *BaseMDLListener) EnterNotExpr(ctx *NotExprContext) {}

// ExitNotExpr is called when production NotExpr is exited.
func (s *BaseMDLListener) ExitNotExpr(ctx *NotExprContext) {}

// EnterIntExpr is called when production IntExpr is entered.
func (s *BaseMDLListener) EnterIntExpr(ctx *IntExprContext) {}

// ExitIntExpr is called when production IntExpr is exited.
func (s *BaseMDLListener) ExitIntExpr(ctx *IntExprContext) {}

// EnterSysVarExpr is called when production SysVarExpr is entered.
func (s *BaseMDLListener) EnterSysVarExpr(ctx *SysVarExprContext) {}

// ExitSysVarExpr is called when production SysVarExpr is exited.
func (s *BaseMDLListener) ExitSysVarExpr(ctx *SysVarExprContext) {}

// EnterInExpr is called when production InExpr is entered.
func (s *BaseMDLListener) EnterInExpr(ctx *InExprContext) {}

// ExitInExpr is called when production InExpr is exited.
func (s *BaseMDLListener) ExitInExpr(ctx *InExprContext) {}

// EnterParenExpr is called when production ParenExpr is entered.
func (s *BaseMDLListener) EnterParenExpr(ctx *ParenExprContext) {}

// ExitParenExpr is called when production ParenExpr is exited.
func (s *BaseMDLListener) ExitParenExpr(ctx *ParenExprContext) {}

// EnterCaseExpr is called when production CaseExpr is entered.
func (s *BaseMDLListener) EnterCaseExpr(ctx *CaseExprContext) {}

// ExitCaseExpr is called when production CaseExpr is exited.
func (s *BaseMDLListener) ExitCaseExpr(ctx *CaseExprContext) {}

// EnterAddSubExpr is called when production AddSubExpr is entered.
func (s *BaseMDLListener) EnterAddSubExpr(ctx *AddSubExprContext) {}

// ExitAddSubExpr is called when production AddSubExpr is exited.
func (s *BaseMDLListener) ExitAddSubExpr(ctx *AddSubExprContext) {}

// EnterSubqueryExpr is called when production SubqueryExpr is entered.
func (s *BaseMDLListener) EnterSubqueryExpr(ctx *SubqueryExprContext) {}

// ExitSubqueryExpr is called when production SubqueryExpr is exited.
func (s *BaseMDLListener) ExitSubqueryExpr(ctx *SubqueryExprContext) {}

// EnterComparisonOp is called when production comparisonOp is entered.
func (s *BaseMDLListener) EnterComparisonOp(ctx *ComparisonOpContext) {}

// ExitComparisonOp is called when production comparisonOp is exited.
func (s *BaseMDLListener) ExitComparisonOp(ctx *ComparisonOpContext) {}

// EnterFunctionCall is called when production functionCall is entered.
func (s *BaseMDLListener) EnterFunctionCall(ctx *FunctionCallContext) {}

// ExitFunctionCall is called when production functionCall is exited.
func (s *BaseMDLListener) ExitFunctionCall(ctx *FunctionCallContext) {}

// EnterAggregateFunction is called when production aggregateFunction is entered.
func (s *BaseMDLListener) EnterAggregateFunction(ctx *AggregateFunctionContext) {}

// ExitAggregateFunction is called when production aggregateFunction is exited.
func (s *BaseMDLListener) ExitAggregateFunction(ctx *AggregateFunctionContext) {}

// EnterCaseExpression is called when production caseExpression is entered.
func (s *BaseMDLListener) EnterCaseExpression(ctx *CaseExpressionContext) {}

// ExitCaseExpression is called when production caseExpression is exited.
func (s *BaseMDLListener) ExitCaseExpression(ctx *CaseExpressionContext) {}

// EnterExpressionList is called when production expressionList is entered.
func (s *BaseMDLListener) EnterExpressionList(ctx *ExpressionListContext) {}

// ExitExpressionList is called when production expressionList is exited.
func (s *BaseMDLListener) ExitExpressionList(ctx *ExpressionListContext) {}

// EnterSystemVariable is called when production systemVariable is entered.
func (s *BaseMDLListener) EnterSystemVariable(ctx *SystemVariableContext) {}

// ExitSystemVariable is called when production systemVariable is exited.
func (s *BaseMDLListener) ExitSystemVariable(ctx *SystemVariableContext) {}

// EnterCreateAssociationStatement is called when production createAssociationStatement is entered.
func (s *BaseMDLListener) EnterCreateAssociationStatement(ctx *CreateAssociationStatementContext) {}

// ExitCreateAssociationStatement is called when production createAssociationStatement is exited.
func (s *BaseMDLListener) ExitCreateAssociationStatement(ctx *CreateAssociationStatementContext) {}

// EnterAssociationType is called when production associationType is entered.
func (s *BaseMDLListener) EnterAssociationType(ctx *AssociationTypeContext) {}

// ExitAssociationType is called when production associationType is exited.
func (s *BaseMDLListener) ExitAssociationType(ctx *AssociationTypeContext) {}

// EnterOwnerType is called when production ownerType is entered.
func (s *BaseMDLListener) EnterOwnerType(ctx *OwnerTypeContext) {}

// ExitOwnerType is called when production ownerType is exited.
func (s *BaseMDLListener) ExitOwnerType(ctx *OwnerTypeContext) {}

// EnterDeleteBehavior is called when production deleteBehavior is entered.
func (s *BaseMDLListener) EnterDeleteBehavior(ctx *DeleteBehaviorContext) {}

// ExitDeleteBehavior is called when production deleteBehavior is exited.
func (s *BaseMDLListener) ExitDeleteBehavior(ctx *DeleteBehaviorContext) {}

// EnterDropAssociationStatement is called when production dropAssociationStatement is entered.
func (s *BaseMDLListener) EnterDropAssociationStatement(ctx *DropAssociationStatementContext) {}

// ExitDropAssociationStatement is called when production dropAssociationStatement is exited.
func (s *BaseMDLListener) ExitDropAssociationStatement(ctx *DropAssociationStatementContext) {}

// EnterQueryStatement is called when production queryStatement is entered.
func (s *BaseMDLListener) EnterQueryStatement(ctx *QueryStatementContext) {}

// ExitQueryStatement is called when production queryStatement is exited.
func (s *BaseMDLListener) ExitQueryStatement(ctx *QueryStatementContext) {}

// EnterShowModules is called when production ShowModules is entered.
func (s *BaseMDLListener) EnterShowModules(ctx *ShowModulesContext) {}

// ExitShowModules is called when production ShowModules is exited.
func (s *BaseMDLListener) ExitShowModules(ctx *ShowModulesContext) {}

// EnterShowEnumerations is called when production ShowEnumerations is entered.
func (s *BaseMDLListener) EnterShowEnumerations(ctx *ShowEnumerationsContext) {}

// ExitShowEnumerations is called when production ShowEnumerations is exited.
func (s *BaseMDLListener) ExitShowEnumerations(ctx *ShowEnumerationsContext) {}

// EnterShowEntities is called when production ShowEntities is entered.
func (s *BaseMDLListener) EnterShowEntities(ctx *ShowEntitiesContext) {}

// ExitShowEntities is called when production ShowEntities is exited.
func (s *BaseMDLListener) ExitShowEntities(ctx *ShowEntitiesContext) {}

// EnterShowEntity is called when production ShowEntity is entered.
func (s *BaseMDLListener) EnterShowEntity(ctx *ShowEntityContext) {}

// ExitShowEntity is called when production ShowEntity is exited.
func (s *BaseMDLListener) ExitShowEntity(ctx *ShowEntityContext) {}

// EnterShowAssociations is called when production ShowAssociations is entered.
func (s *BaseMDLListener) EnterShowAssociations(ctx *ShowAssociationsContext) {}

// ExitShowAssociations is called when production ShowAssociations is exited.
func (s *BaseMDLListener) ExitShowAssociations(ctx *ShowAssociationsContext) {}

// EnterShowAssociation is called when production ShowAssociation is entered.
func (s *BaseMDLListener) EnterShowAssociation(ctx *ShowAssociationContext) {}

// ExitShowAssociation is called when production ShowAssociation is exited.
func (s *BaseMDLListener) ExitShowAssociation(ctx *ShowAssociationContext) {}

// EnterDescribeEnumeration is called when production DescribeEnumeration is entered.
func (s *BaseMDLListener) EnterDescribeEnumeration(ctx *DescribeEnumerationContext) {}

// ExitDescribeEnumeration is called when production DescribeEnumeration is exited.
func (s *BaseMDLListener) ExitDescribeEnumeration(ctx *DescribeEnumerationContext) {}

// EnterDescribeEntity is called when production DescribeEntity is entered.
func (s *BaseMDLListener) EnterDescribeEntity(ctx *DescribeEntityContext) {}

// ExitDescribeEntity is called when production DescribeEntity is exited.
func (s *BaseMDLListener) ExitDescribeEntity(ctx *DescribeEntityContext) {}

// EnterDescribeAssociation is called when production DescribeAssociation is entered.
func (s *BaseMDLListener) EnterDescribeAssociation(ctx *DescribeAssociationContext) {}

// ExitDescribeAssociation is called when production DescribeAssociation is exited.
func (s *BaseMDLListener) ExitDescribeAssociation(ctx *DescribeAssociationContext) {}

// EnterRepositoryStatement is called when production repositoryStatement is entered.
func (s *BaseMDLListener) EnterRepositoryStatement(ctx *RepositoryStatementContext) {}

// ExitRepositoryStatement is called when production repositoryStatement is exited.
func (s *BaseMDLListener) ExitRepositoryStatement(ctx *RepositoryStatementContext) {}

// EnterCommitStatement is called when production commitStatement is entered.
func (s *BaseMDLListener) EnterCommitStatement(ctx *CommitStatementContext) {}

// ExitCommitStatement is called when production commitStatement is exited.
func (s *BaseMDLListener) ExitCommitStatement(ctx *CommitStatementContext) {}

// EnterUpdateStatement is called when production updateStatement is entered.
func (s *BaseMDLListener) EnterUpdateStatement(ctx *UpdateStatementContext) {}

// ExitUpdateStatement is called when production updateStatement is exited.
func (s *BaseMDLListener) ExitUpdateStatement(ctx *UpdateStatementContext) {}

// EnterRefreshStatement is called when production refreshStatement is entered.
func (s *BaseMDLListener) EnterRefreshStatement(ctx *RefreshStatementContext) {}

// ExitRefreshStatement is called when production refreshStatement is exited.
func (s *BaseMDLListener) ExitRefreshStatement(ctx *RefreshStatementContext) {}

// EnterSessionStatement is called when production sessionStatement is entered.
func (s *BaseMDLListener) EnterSessionStatement(ctx *SessionStatementContext) {}

// ExitSessionStatement is called when production sessionStatement is exited.
func (s *BaseMDLListener) ExitSessionStatement(ctx *SessionStatementContext) {}

// EnterSetStatement is called when production setStatement is entered.
func (s *BaseMDLListener) EnterSetStatement(ctx *SetStatementContext) {}

// ExitSetStatement is called when production setStatement is exited.
func (s *BaseMDLListener) ExitSetStatement(ctx *SetStatementContext) {}

// EnterHelpStatement is called when production helpStatement is entered.
func (s *BaseMDLListener) EnterHelpStatement(ctx *HelpStatementContext) {}

// ExitHelpStatement is called when production helpStatement is exited.
func (s *BaseMDLListener) ExitHelpStatement(ctx *HelpStatementContext) {}

// EnterExitStatement is called when production exitStatement is entered.
func (s *BaseMDLListener) EnterExitStatement(ctx *ExitStatementContext) {}

// ExitExitStatement is called when production exitStatement is exited.
func (s *BaseMDLListener) ExitExitStatement(ctx *ExitStatementContext) {}

// EnterExecuteScriptStatement is called when production executeScriptStatement is entered.
func (s *BaseMDLListener) EnterExecuteScriptStatement(ctx *ExecuteScriptStatementContext) {}

// ExitExecuteScriptStatement is called when production executeScriptStatement is exited.
func (s *BaseMDLListener) ExitExecuteScriptStatement(ctx *ExecuteScriptStatementContext) {}

// EnterQualifiedName is called when production qualifiedName is entered.
func (s *BaseMDLListener) EnterQualifiedName(ctx *QualifiedNameContext) {}

// ExitQualifiedName is called when production qualifiedName is exited.
func (s *BaseMDLListener) ExitQualifiedName(ctx *QualifiedNameContext) {}

// EnterDocumentation is called when production documentation is entered.
func (s *BaseMDLListener) EnterDocumentation(ctx *DocumentationContext) {}

// ExitDocumentation is called when production documentation is exited.
func (s *BaseMDLListener) ExitDocumentation(ctx *DocumentationContext) {}

// EnterPositionAnnotation is called when production positionAnnotation is entered.
func (s *BaseMDLListener) EnterPositionAnnotation(ctx *PositionAnnotationContext) {}

// ExitPositionAnnotation is called when production positionAnnotation is exited.
func (s *BaseMDLListener) ExitPositionAnnotation(ctx *PositionAnnotationContext) {}

// EnterStringLiteral is called when production stringLiteral is entered.
func (s *BaseMDLListener) EnterStringLiteral(ctx *StringLiteralContext) {}

// ExitStringLiteral is called when production stringLiteral is exited.
func (s *BaseMDLListener) ExitStringLiteral(ctx *StringLiteralContext) {}
