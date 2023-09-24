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