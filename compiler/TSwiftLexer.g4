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

// Control flow

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