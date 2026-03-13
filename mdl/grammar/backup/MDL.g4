/**
 * MDL (Mendix Definition Language) Grammar
 * Version: 2.0
 * Scope: Domain Model (Entities, Attributes, Associations, Enumerations, Views)
 *
 * This grammar can be used to generate parsers for multiple languages
 * using ANTLR4: Go, TypeScript, Java, Python, etc.
 */
grammar MDL;

// ============================================================================
// Parser Rules
// ============================================================================

// Entry point - a program is a sequence of statements
// SLASH tokens can appear between statements as separators
program
    : (SLASH* statement SLASH*)* EOF
    ;

statement
    : connectionStatement
    | ddlStatement
    | queryStatement
    | repositoryStatement
    | sessionStatement
    ;

// Statement terminator - semicolon or forward slash
terminator
    : SEMI
    | SLASH
    ;

// ----------------------------------------------------------------------------
// Connection Statements
// ----------------------------------------------------------------------------

connectionStatement
    : connectStatement
    | disconnectStatement
    | statusStatement
    ;

connectStatement
    : CONNECT LOCAL stringLiteral terminator?                   # ConnectLocal
    | CONNECT TO FILESYSTEM stringLiteral terminator?           # ConnectFilesystem
    ;

disconnectStatement
    : DISCONNECT terminator?
    ;

statusStatement
    : (STATUS | SHOW STATUS) terminator?
    ;

// ----------------------------------------------------------------------------
// DDL Statements (Data Definition Language)
// ----------------------------------------------------------------------------

ddlStatement
    : createStatement
    | alterStatement
    | dropStatement
    ;

createStatement
    : createModuleStatement
    | createEnumerationStatement
    | createEntityStatement
    | createViewEntityStatement
    | createAssociationStatement
    ;

// ----------------------------------------------------------------------------
// Module Statements
// ----------------------------------------------------------------------------

createModuleStatement
    : CREATE MODULE IDENTIFIER terminator?
    ;

alterStatement
    : alterEnumerationStatement
    ;

dropStatement
    : dropModuleStatement
    | dropEnumerationStatement
    | dropEntityStatement
    | dropAssociationStatement
    ;

dropModuleStatement
    : DROP MODULE IDENTIFIER terminator?
    ;

// ----------------------------------------------------------------------------
// Enumeration Statements
// ----------------------------------------------------------------------------

createEnumerationStatement
    : documentation?
      CREATE ENUMERATION qualifiedName
      LPAREN enumValueList RPAREN
      (COMMENT stringLiteral)?
      terminator?
    ;

enumValueList
    : enumValue (COMMA enumValue)*
    ;

enumValue
    : documentation? IDENTIFIER stringLiteral
    ;

alterEnumerationStatement
    : ALTER ENUMERATION qualifiedName alterEnumOperation terminator?
    ;

alterEnumOperation
    : ADD VALUE IDENTIFIER stringLiteral      # AddEnumValue
    | DROP VALUE IDENTIFIER                   # DropEnumValue
    | RENAME VALUE IDENTIFIER TO IDENTIFIER   # RenameEnumValue
    ;

dropEnumerationStatement
    : DROP ENUMERATION qualifiedName terminator?
    ;

// ----------------------------------------------------------------------------
// Entity Statements
// ----------------------------------------------------------------------------

createEntityStatement
    : documentation?
      positionAnnotation?
      createOrModify?
      entityType ENTITY qualifiedName
      LPAREN attributeList? RPAREN
      (POSITION LPAREN INTEGER COMMA INTEGER RPAREN)?
      indexClause*
      (COMMENT stringLiteral)?
      terminator?
    ;

createOrModify
    : CREATE OR MODIFY
    | CREATE
    ;

entityType
    : PERSISTENT
    | NON_PERSISTENT
    ;

attributeList
    : attribute (COMMA attribute)*
    ;

attribute
    : documentation?
      renamedFromAnnotation?
      attributeName COLON dataType attributeConstraints? (COMMENT stringLiteral)?
    ;

// attributeName allows IDENTIFIER and certain reserved words that are commonly used as attribute names
attributeName
    : IDENTIFIER
    | STATUS        // Allow Status as attribute name
    | TYPE          // Allow Type as attribute name
    | VALUE         // Allow Value as attribute name
    | DEFAULT       // Allow Default as attribute name
    | INDEX         // Allow Index as attribute name
    ;

renamedFromAnnotation
    : AT_RENAMED_FROM LPAREN stringLiteral RPAREN
    ;

dataType
    : STRING (LPAREN (INTEGER | UNLIMITED) RPAREN)?     # StringType
    | INTEGER_TYPE                                       # IntegerType
    | LONG                                               # LongType
    | DECIMAL (LPAREN INTEGER (COMMA INTEGER)? RPAREN)? # DecimalType
    | BOOLEAN                                            # BooleanType
    | DATETIME                                           # DateTimeType
    | DATE                                               # DateType
    | AUTONUMBER                                         # AutoNumberType
    | BINARY                                             # BinaryType
    | ENUMERATION LPAREN qualifiedName RPAREN            # EnumerationType
    ;

attributeConstraints
    : attributeConstraint+
    ;

attributeConstraint
    : NOT NULL (ERROR stringLiteral)?         # NotNullConstraint
    | UNIQUE (ERROR stringLiteral)?           # UniqueConstraint
    | DEFAULT defaultValue                    # DefaultConstraint
    ;

defaultValue
    : stringLiteral
    | INTEGER
    | DECIMAL_NUMBER
    | TRUE
    | FALSE
    | CURRENT_TIMESTAMP
    | qualifiedName               // For enum default values
    ;

indexClause
    : COMMA? INDEX LPAREN indexColumnList RPAREN
    ;

indexColumnList
    : indexColumn (COMMA indexColumn)*
    ;

indexColumn
    : attributeName (ASC | DESC)?
    ;

dropEntityStatement
    : DROP ENTITY qualifiedName terminator?
    ;

// ----------------------------------------------------------------------------
// View Entity Statements
// ----------------------------------------------------------------------------

createViewEntityStatement
    : documentation?
      positionAnnotation?
      createOrModify?
      VIEW ENTITY qualifiedName
      LPAREN viewAttributeList? RPAREN
      AS oqlQuery
      terminator?
    ;

viewAttributeList
    : viewAttribute (COMMA viewAttribute)*
    ;

viewAttribute
    : attributeName COLON dataType
    ;

// OQL Query (simplified - captures the query as tokens)
oqlQuery
    : selectClause fromClause whereClause? groupByClause? orderByClause? limitClause?
    ;

selectClause
    : SELECT DISTINCT? selectItems
    ;

selectItems
    : selectItem (COMMA selectItem)*
    ;

selectItem
    : expression (AS attributeName)?
    ;

fromClause
    : FROM fromItem joinClause*
    ;

fromItem
    : qualifiedName (AS attributeName)?
    ;

joinClause
    : (INNER | LEFT | RIGHT | OUTER)? JOIN joinTarget (ON expression)?
    ;

joinTarget
    : qualifiedName (AS attributeName)?                    # SimpleJoinTarget
    | attributeName SLASH qualifiedName SLASH qualifiedName (AS attributeName)?  # AssociationJoinTarget
    ;

whereClause
    : WHERE expression
    ;

groupByClause
    : GROUP BY expressionList
    ;

orderByClause
    : ORDER BY orderByItem (COMMA orderByItem)*
    ;

orderByItem
    : expression (ASC | DESC)?
    ;

limitClause
    : (OFFSET INTEGER)? LIMIT INTEGER
    | LIMIT INTEGER (OFFSET INTEGER)?
    ;

// Expression (simplified for OQL)
expression
    : LPAREN expression RPAREN                                      # ParenExpr
    | expression (STAR | SLASH | PERCENT) expression                # MulDivExpr
    | expression (PLUS | MINUS) expression                          # AddSubExpr
    | expression (COLON) expression                                 # DivisionExpr
    | expression comparisonOp expression                            # CompareExpr
    | expression AND expression                                     # AndExpr
    | expression OR expression                                      # OrExpr
    | NOT expression                                                # NotExpr
    | expression IS NOT? NULL                                       # IsNullExpr
    | expression NOT? IN LPAREN expressionList RPAREN               # InExpr
    | expression NOT? IN LPAREN oqlQuery RPAREN                     # InSubqueryExpr
    | functionCall                                                  # FuncExpr
    | caseExpression                                                # CaseExpr
    | LPAREN oqlQuery RPAREN                                        # SubqueryExpr
    | qualifiedName                                                 # QualifiedExpr
    | IDENTIFIER DOT attributeName                                  # FieldAccessExpr
    | IDENTIFIER                                                    # IdentExpr
    | stringLiteral                                                 # StringExpr
    | INTEGER                                                       # IntExpr
    | DECIMAL_NUMBER                                                # DecimalExpr
    | TRUE                                                          # TrueExpr
    | FALSE                                                         # FalseExpr
    | STAR                                                          # StarExpr
    | systemVariable                                                # SysVarExpr
    ;

comparisonOp
    : EQ | NE | LT | LE | GT | GE | LIKE
    ;

functionCall
    : IDENTIFIER LPAREN (expression (COMMA expression)*)? RPAREN
    | aggregateFunction LPAREN (DISTINCT? expression | STAR)? RPAREN
    | DATEPART LPAREN IDENTIFIER COMMA expression RPAREN               // DATEPART(YEAR, date)
    | DATEDIFF LPAREN IDENTIFIER COMMA expression COMMA expression RPAREN  // DATEDIFF(DAY, date1, date2)
    ;

aggregateFunction
    : COUNT | SUM | AVG | MIN | MAX
    ;

caseExpression
    : CASE (WHEN expression THEN expression)+ (ELSE expression)? END
    ;

expressionList
    : expression (COMMA expression)*
    ;

systemVariable
    : LBRACKET PERCENT IDENTIFIER PERCENT RBRACKET
    ;

// ----------------------------------------------------------------------------
// Association Statements
// ----------------------------------------------------------------------------

createAssociationStatement
    : documentation?
      CREATE ASSOCIATION qualifiedName
      FROM qualifiedName
      TO qualifiedName
      TYPE associationType
      (OWNER ownerType)?
      (DELETE_BEHAVIOR deleteBehavior)?
      (COMMENT stringLiteral)?
      terminator?
    ;

associationType
    : REFERENCE | REF           // Reference, Ref
    | REFERENCE_SET | REFSET    // ReferenceSet, RefSet
    ;

ownerType
    : DEFAULT
    | BOTH
    | PARENT
    | CHILD
    ;

deleteBehavior
    : DELETE_BUT_KEEP_REFERENCES
    | CASCADE
    | DELETE_BOTH
    | KEEP_PARENT_DELETE_CHILD
    | KEEP_CHILD_DELETE_PARENT
    | DELETE_IF_NO_REFERENCES
    ;

dropAssociationStatement
    : DROP ASSOCIATION qualifiedName terminator?
    ;

// ----------------------------------------------------------------------------
// Query Statements
// ----------------------------------------------------------------------------

queryStatement
    : showStatement
    | describeStatement
    ;

showStatement
    : SHOW MODULES terminator?                                    # ShowModules
    | SHOW ENUMERATIONS (IN IDENTIFIER)? terminator?              # ShowEnumerations
    | SHOW ENTITIES (IN IDENTIFIER)? terminator?                  # ShowEntities
    | SHOW ENTITY qualifiedName terminator?                       # ShowEntity
    | SHOW ASSOCIATIONS (IN IDENTIFIER)? terminator?              # ShowAssociations
    | SHOW ASSOCIATION qualifiedName terminator?                  # ShowAssociation
    ;

describeStatement
    : DESCRIBE ENUMERATION qualifiedName terminator?              # DescribeEnumeration
    | DESCRIBE ENTITY qualifiedName terminator?                   # DescribeEntity
    | DESCRIBE ASSOCIATION qualifiedName terminator?              # DescribeAssociation
    ;

// ----------------------------------------------------------------------------
// Repository Statements
// ----------------------------------------------------------------------------

repositoryStatement
    : commitStatement
    | updateStatement
    | refreshStatement
    ;

commitStatement
    : COMMIT (MESSAGE stringLiteral)? terminator?
    ;

updateStatement
    : UPDATE terminator?
    ;

refreshStatement
    : REFRESH terminator?
    ;

// ----------------------------------------------------------------------------
// Session Statements
// ----------------------------------------------------------------------------

sessionStatement
    : setStatement
    | helpStatement
    | exitStatement
    | executeScriptStatement
    ;

setStatement
    : SET IDENTIFIER EQ (stringLiteral | INTEGER | TRUE | FALSE) terminator?
    ;

helpStatement
    : (HELP | QUESTION) terminator?
    ;

exitStatement
    : (EXIT | QUIT) terminator?
    ;

executeScriptStatement
    : EXECUTE SCRIPT stringLiteral terminator?
    ;

// ----------------------------------------------------------------------------
// Common Rules
// ----------------------------------------------------------------------------

qualifiedName
    : IDENTIFIER (DOT IDENTIFIER)?
    ;

documentation
    : DOC_COMMENT
    ;

positionAnnotation
    : AT_POSITION LPAREN INTEGER COMMA INTEGER RPAREN
    ;

stringLiteral
    : STRING_LITERAL
    ;

// ============================================================================
// Lexer Rules
// ============================================================================

// Keywords - Connection
CONNECT         : C O N N E C T ;
LOCAL           : L O C A L ;
TO              : T O ;
FILESYSTEM      : F I L E S Y S T E M ;
DISCONNECT      : D I S C O N N E C T ;
STATUS          : S T A T U S ;

// Keywords - DDL
CREATE          : C R E A T E ;
ALTER           : A L T E R ;
DROP            : D R O P ;
ENUMERATION     : E N U M E R A T I O N ;
ENTITY          : E N T I T Y ;
PERSISTENT      : P E R S I S T E N T ;
NON_PERSISTENT  : N O N '-' P E R S I S T E N T ;
VIEW            : V I E W ;
ASSOCIATION     : A S S O C I A T I O N ;
FROM            : F R O M ;
TYPE            : T Y P E ;
OWNER           : O W N E R ;
INDEX           : I N D E X ;
OR              : O R ;
MODIFY          : M O D I F Y ;

// Keywords - Enumeration operations
ADD             : A D D ;
VALUE           : V A L U E ;
RENAME          : R E N A M E ;

// Keywords - Data types
STRING          : S T R I N G ;
INTEGER_TYPE    : I N T E G E R ;
LONG            : L O N G ;
DECIMAL         : D E C I M A L ;
BOOLEAN         : B O O L E A N ;
DATETIME        : D A T E T I M E ;
DATE            : D A T E ;
AUTONUMBER      : A U T O N U M B E R ;
BINARY          : B I N A R Y ;
UNLIMITED       : U N L I M I T E D ;

// Keywords - Constraints
NOT             : N O T ;
NULL            : N U L L ;
UNIQUE          : U N I Q U E ;
DEFAULT         : D E F A U L T ;
ERROR           : E R R O R ;
CURRENT_TIMESTAMP : C U R R E N T '_' T I M E S T A M P ;

// Keywords - Association
REFERENCE       : R E F E R E N C E ;
REF             : R E F ;
REFERENCE_SET   : R E F E R E N C E S E T ;
REFSET          : R E F S E T ;
BOTH            : B O T H ;
PARENT          : P A R E N T ;
CHILD           : C H I L D ;
DELETE_BEHAVIOR : D E L E T E '_' B E H A V I O R ;
DELETE_BUT_KEEP_REFERENCES : D E L E T E '_' B U T '_' K E E P '_' R E F E R E N C E S ;
CASCADE         : C A S C A D E ;
DELETE_BOTH     : D E L E T E '_' B O T H ;
KEEP_PARENT_DELETE_CHILD : K E E P '_' P A R E N T '_' D E L E T E '_' C H I L D ;
KEEP_CHILD_DELETE_PARENT : K E E P '_' C H I L D '_' D E L E T E '_' P A R E N T ;
DELETE_IF_NO_REFERENCES : D E L E T E '_' I F '_' N O '_' R E F E R E N C E S ;

// Keywords - Query
SHOW            : S H O W ;
DESCRIBE        : D E S C R I B E ;
MODULES         : M O D U L E S ;
MODULE          : M O D U L E ;
ENUMERATIONS    : E N U M E R A T I O N S ;
ENTITIES        : E N T I T I E S ;
ASSOCIATIONS    : A S S O C I A T I O N S ;
IN              : I N ;

// Keywords - OQL
SELECT          : S E L E C T ;
AS              : A S ;
WHERE           : W H E R E ;
AND             : A N D ;
GROUP           : G R O U P ;
BY              : B Y ;
ORDER           : O R D E R ;
ASC             : A S C ;
DESC            : D E S C ;
LIMIT           : L I M I T ;
OFFSET          : O F F S E T ;
JOIN            : J O I N ;
INNER           : I N N E R ;
LEFT            : L E F T ;
RIGHT           : R I G H T ;
OUTER           : O U T E R ;
ON              : O N ;
DISTINCT        : D I S T I N C T ;
CASE            : C A S E ;
WHEN            : W H E N ;
THEN            : T H E N ;
ELSE            : E L S E ;
END             : E N D ;
IS              : I S ;
LIKE            : L I K E ;

// Keywords - Aggregate functions
COUNT           : C O U N T ;
SUM             : S U M ;
AVG             : A V G ;
MIN             : M I N ;
MAX             : M A X ;

// Keywords - Date functions
DATEPART        : D A T E P A R T ;
DATEDIFF        : D A T E D I F F ;

// Keywords - Repository
COMMIT          : C O M M I T ;
MESSAGE         : M E S S A G E ;
UPDATE          : U P D A T E ;
REFRESH         : R E F R E S H ;

// Keywords - Session
SET             : S E T ;
HELP            : H E L P ;
EXIT            : E X I T ;
QUIT            : Q U I T ;
EXECUTE         : E X E C U T E ;
SCRIPT          : S C R I P T ;

// Keywords - Misc
COMMENT         : C O M M E N T ;
POSITION        : P O S I T I O N ;
TRUE            : T R U E ;
FALSE           : F A L S E ;

// Annotations
AT_POSITION     : '@' P O S I T I O N ;
AT_RENAMED_FROM : '@' R E N A M E D F R O M ;

// Symbols
LPAREN          : '(' ;
RPAREN          : ')' ;
LBRACE          : '{' ;
RBRACE          : '}' ;
LBRACKET        : '[' ;
RBRACKET        : ']' ;
COMMA           : ',' ;
SEMI            : ';' ;
SLASH           : '/' ;
COLON           : ':' ;
DOT             : '.' ;
EQ              : '=' ;
NE              : '!=' | '<>' ;
LT              : '<' ;
LE              : '<=' ;
GT              : '>' ;
GE              : '>=' ;
PLUS            : '+' ;
MINUS           : '-' ;
STAR            : '*' ;
PERCENT         : '%' ;
QUESTION        : '?' ;

// Literals
INTEGER
    : '-'? DIGIT+
    ;

DECIMAL_NUMBER
    : '-'? DIGIT+ '.' DIGIT+
    ;

STRING_LITERAL
    : '\'' ( ~['\\\r\n] | '\\' . | '\'\'' )* '\''
    ;

// Documentation comment
DOC_COMMENT
    : '/**' .*? '*/'
    ;

// Identifier
IDENTIFIER
    : LETTER (LETTER | DIGIT | '_')*
    ;

// Whitespace and comments
WS
    : [ \t\r\n]+ -> skip
    ;

LINE_COMMENT
    : '--' ~[\r\n]* -> skip
    ;

BLOCK_COMMENT
    : '/*' .*? '*/' -> skip
    ;

// Fragment rules for case-insensitive keywords
fragment A : [aA] ;
fragment B : [bB] ;
fragment C : [cC] ;
fragment D : [dD] ;
fragment E : [eE] ;
fragment F : [fF] ;
fragment G : [gG] ;
fragment H : [hH] ;
fragment I : [iI] ;
fragment J : [jJ] ;
fragment K : [kK] ;
fragment L : [lL] ;
fragment M : [mM] ;
fragment N : [nN] ;
fragment O : [oO] ;
fragment P : [pP] ;
fragment Q : [qQ] ;
fragment R : [rR] ;
fragment S : [sS] ;
fragment T : [tT] ;
fragment U : [uU] ;
fragment V : [vV] ;
fragment W : [wW] ;
fragment X : [xX] ;
fragment Y : [yY] ;
fragment Z : [zZ] ;

fragment DIGIT  : [0-9] ;
fragment LETTER : [a-zA-Z] ;
