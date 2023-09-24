
3 address code = 3AC = TAC

The stmts of 3AC are:

TACBlock -> []TACStmt

0. TACStmt (INTERFACE)
    components:
        id [int]
        line [int]
        ToString() string

    1. Label (STRUCT)

    2. Assignment (INTERFACE)
        components:
            Assignee [SimpleValue]

        CompoundAssignment (STRUCT)
            components:
                Left [SimpleValue]
                Right [SimpleValue]
                Operator [ArithmeticOperator]

        SimpleAssignment (STRUCT)
            components:
                Val [Value]

    3. Jump statements (INTERFACE)

        Conditional jump statements (STRUCT)
            components:
                Condition [BooleanExpression]
                Target [Label]

        Unconditional jump statements (STRUCT)
            components:
                Target [Label]

    4. Method Dcl (STRUCT)
        components:
            Name
            Block [TACBlock]

    5. Method Call (STRUCT)
        components:
            Name

    6. Value (INTERFACE)

        SimpleValue (INTERFACE)
            components:
                Cast [CastType]

            Temporal (STRUCT)
            HeapPointer (STRUCT)
            StackPointer (STRUCT)
            Literal (STRUCT)
                components:
                    Val [int/float]

        IndexedValue (INTERFACE)
            components:
                Index [SimpleValue]

            HeapIndexed (STRUCT)
            StackIndexed (STRUCT)

    7. ArithmeticOperator (Enum)
        [+, -, *, /, %]

    8. CastType (Enum)
        [float]

    9. BooleanExpression (STRUCT)
        components:
            Left [SimpleValue]
            Right [SimpleValue]
            Operator [BooleanOperator]

    10. BooleanOperator (Enum)
        [==, !=, <, >, <=, >=]

    11. print (STRUCT)
        components:
            Val [SimpleValue]

    12. Comment (STRUCT)
            components:
                Val [string]

Basic ideas:

- TACBlock is just a list of TACStmts
- All values must provide a prop named "Cast" which is a CastType
- CastType <-> Enum

```go
type TACStmt interface {
    ToString() string
}
```


std lib:
printString
compareString

1. Value system
    1.1. Primitives: int, double, bool, char
2. Variable asgmt