```g4

lexer grammar TSwiftLexer;

// Skip tokens

WS: [ \t\r\n]+ -> skip;
COMMENT: '//' .*? ('\n' | EOF) -> skip;
MULTILINE_COMMENT: '/*' .*? '*/' -> skip;

// Stmts delimiter
SEMICOLON: ';';

// Keywords

// Declarations

LET_KW: 'let';
VAR_KW: 'var';
FUNC_KW: 'func';
STRUCT_KW: 'struct';

// Control flowa

IF_KW: 'if';
ELSE_KW: 'else';
SWITCH_KW: 'switch';
CASE_KW: 'case';
DEFAULT_KW: 'default';
FOR_KW: 'for';
WHILE_KW: 'while';
BREAK_KW: 'break';
CONTINUE_KW: 'continue';
RETURN_KW: 'return';
GUARD_KW: 'guard';
INOUT_KW: 'inout';
MUTATING_KW: 'mutating';

// ...
IN_KW: 'in';

// Types INTEGER_TYPE: 'Int'; FLOAT_TYPE: 'Float'; STRING_TYPE: 'String'; BOOL_TYPE: 'Bool';
// CHARACTER_TYPE: 'Character';

// Literals

INTEGER_LITERAL: [0-9]+;
FLOAT_LITERAL: [0-9]+ '.' [0-9]+;
// STRING LITERAL WITH SCAPED SEQUENCES
STRING_LITERAL: '"' (~["\r\n\\] | ESC_SEQ)* '"';
BOOL_LITERAL: 'true' | 'false';
NIL_LITERAL: 'nil';

// Identifiers
ID: [a-zA-Z_][a-zA-Z0-9_]*;

// Arithmetic operators

PLUS: '+';
MINUS: '-';
MULT: '*';
DIV: '/';
MOD: '%';

// Assignment operators

EQUALS: '=';
PLUS_EQUALS: '+=';
MINUS_EQUALS: '-=';

// Comparison operators

EQUALS_EQUALS: '==';
NOT_EQUALS: '!=';
LESS_THAN: '<';
LESS_THAN_OR_EQUAL: '<=';
GREATER_THAN: '>';
GREATER_THAN_OR_EQUAL: '>=';

// Logical operators

AND: '&&';
OR: '||';
NOT: '!';
// split into two tokens?

// Delimiters

LPAREN: '(';
RPAREN: ')';
LBRACE: '{';
RBRACE: '}';
LBRACK: '[';
RBRACK: ']';

// Other

COMMA: ',';
DOT: '.';
COLON: ':';
ARROW: '->';
INTERROGATION: '?';
ANPERSAND: '&';

// Error?

fragment ESC_SEQ: '\\' [btnfr"'\\];

// ------------ Grammar ------------

parser grammar TSwiftLanguage;

options {
	tokenVocab = TSwiftLexer;
	// language = Swift; superClass = SwiftParserBaseListener;
}

// make stmt* a new rule
program: (stmt)* EOF?;

delimiter: SEMICOLON? | EOF;

stmt:
	decl_stmt delimiter
	| assign_stmt delimiter
	| transfer_stmt delimiter
	| if_stmt
	| switch_stmt
	| while_stmt
	| for_stmt
	| guard_stmt
	| func_call delimiter
	| vector_func delimiter
	| func_dcl
	| strct_dcl;

decl_stmt:
	var_type ID COLON type EQUALS expr		# TypeValueDecl
	| var_type ID EQUALS expr				# ValueDecl
	| var_type ID COLON type INTERROGATION	# TypeDecl;

vector_expr:
	LBRACK (expr (COMMA expr)*)? RBRACK # VectorItemList;

vector_item: id_pattern (LBRACK expr RBRACK)+ # VectorItem;

vector_prop: vector_item DOT id_pattern # VectorProp;
vector_func: vector_item DOT func_call # VectorFunc;

repeating:
	(vector_type | matrix_type) LPAREN ID COLON expr COMMA ID COLON expr RPAREN;

var_type: VAR_KW | LET_KW;

type: ID | vector_type | matrix_type;

vector_type: LBRACK ID RBRACK;

matrix_type: aux_matrix_type | LBRACK LBRACK ID RBRACK RBRACK;

aux_matrix_type: LBRACK matrix_type RBRACK;

assign_stmt:
	id_pattern EQUALS expr											# DirectAssign
	| id_pattern op = (PLUS_EQUALS | MINUS_EQUALS) expr				# ArithmeticAssign
	| vector_item op = (PLUS_EQUALS | MINUS_EQUALS | EQUALS) expr	# VectorAssign;

id_pattern: ID (DOT ID)* # IdPattern;

literal:
	INTEGER_LITERAL		# IntLiteral
	| FLOAT_LITERAL		# FloatLiteral
	| STRING_LITERAL	# StringLiteral
	| BOOL_LITERAL		# BoolLiteral
	| NIL_LITERAL		# NilLiteral;

expr:
	LPAREN expr RPAREN									# ParenExp // (a)
	| func_call											# FuncCallExp // a.a.a()
	| id_pattern										# IdExp // a.a.a
	| vector_item										# VectorItemExp // a.a.a[0]
	| vector_prop										# VectorPropExp // a[0].a.a
	| vector_func										# VectorFuncExp // a[0].a.a()
	| literal											# LiteralExp // 1, 1.0, "a", true, nil
	| vector_expr										# VectorExp // [1, 2, 3]
	| repeating											# RepeatingExp // [ Int ] (repeating: 0, count: 3)
	| struct_vector										# StructVectorExp // [ Int ]()	
	| op = (NOT | MINUS) expr							# UnaryExp // !a, -a	
	| left = expr op = (MULT | DIV | MOD) right = expr	# BinaryExp // a * b, a / b, a % b
	| left = expr op = (PLUS | MINUS) right = expr		# BinaryExp // a + b, a - b
	| left = expr op = (
		LESS_THAN
		| LESS_THAN_OR_EQUAL
		| GREATER_THAN
		| GREATER_THAN_OR_EQUAL
	) right = expr													# BinaryExp // a < b, a <= b, a > b, a >= b
	| left = expr op = (EQUALS_EQUALS | NOT_EQUALS) right = expr	# BinaryExp // a == b, a != b
	| left = expr op = AND right = expr								# BinaryExp // a && b
	| left = expr op = OR right = expr								# BinaryExp; // a || b
// StructMethodCallExp, StructPropertyCallExp, FunctionCallExp, vector, matrix;  (++, --)?

if_stmt: if_chain (ELSE_KW if_chain)* else_stmt? # IfStmt;

if_chain: IF_KW expr LBRACE stmt* RBRACE # IfChain;
else_stmt: ELSE_KW LBRACE stmt* RBRACE # ElseStmt;

switch_stmt:
	SWITCH_KW expr LBRACE switch_case* default_case? RBRACE # SwitchStmt;

switch_case: CASE_KW expr COLON stmt* # SwitchCase;

default_case: DEFAULT_KW COLON stmt* # DefaultCase;

while_stmt: WHILE_KW expr LBRACE stmt* RBRACE # WhileStmt;

for_stmt:
	FOR_KW ID IN_KW (expr | range) LBRACE stmt* RBRACE # ForStmt;

range: expr DOT DOT DOT expr # NumericRange;

guard_stmt:
	GUARD_KW expr ELSE_KW LBRACE stmt* RBRACE # GuardStmt;

transfer_stmt:
	RETURN_KW expr?	# ReturnStmt
	| BREAK_KW		# BreakStmt
	| CONTINUE_KW	# ContinueStmt;

func_call: id_pattern LPAREN arg_list? RPAREN # FuncCall;

// external names -> num: value, num2: value2
arg_list: func_arg (COMMA func_arg)* # ArgList;
func_arg: (ID COLON)? (ANPERSAND)? (id_pattern | expr) # FuncArg; // 

func_dcl:
	FUNC_KW ID LPAREN param_list? RPAREN (ARROW type)? LBRACE stmt* RBRACE # FuncDecl;

param_list: func_param (COMMA func_param)* # ParamList;
func_param: ID? ID COLON INOUT_KW? type # FuncParam;

// * Structs

strct_dcl: STRUCT_KW ID LBRACE struct_prop* RBRACE # StructDecl;

struct_prop:
	var_type ID (COLON type)? (EQUALS expr)?	# StructAttr
	| MUTATING_KW? func_dcl						# StructFunc;

struct_vector: LBRACK ID RBRACK LPAREN RPAREN # StructVector;

```