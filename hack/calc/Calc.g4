// Calc.g4
grammar Calc;

//split into two sections, tokens, and rules.
//The tokens(Start with UpperCase) are terminal symbols in the grammar, that is, they are made up of nothing but literal characters.
//lex: char inputstream => tokens
//Whereas rules(Start with LowerCase) will combine tokens and/or other rules to establish an AST(Abstract Syntax Tree)
//and for ANTLR, rules got listeners

// Tokens
MUL: '*';
DIV: '/';
ADD: '+';
SUB: '-';
NUMBER: [0-9]+;
WHITESPACE: [ \r\n\t]+ -> skip;

// Rules
start : expression EOF;

expression
   : expression op=('*'|'/') expression # MulDiv
   | expression op=('+'|'-') expression # AddSub
   | NUMBER                             # Number
   ;