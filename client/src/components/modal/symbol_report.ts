
export interface SymbolTableI {
    GlobalScope: ScopeI
}

export interface ScopeI {
    Name: string
    Vars: SymbolI[]
    Funcs: SymbolI[]
    Structs: SymbolI[]
    ChildScopes: ScopeI[] | null
}

export interface SymbolI {
    Name: string
    Type: string
    Line: number
    Column: number
}

class Scope implements ScopeI {

    static idCounter = 0
    static maxDepth = 20 // <-- max depth of the rendered children scopes

    Vars: SymbolI[]
    Funcs: SymbolI[]
    Structs: SymbolI[]
    ChildScopes: ScopeI[] | null
    Name: string = ''
    id: number = 0

    constructor({ Vars: vars, Funcs: funcs, Structs: structs, ChildScopes: childScopes, Name }: ScopeI) {
        this.Vars = vars
        this.Funcs = funcs
        this.Structs = structs
        this.ChildScopes = childScopes
        this.Name = Name
        this.id = Scope.idCounter++
    }

    graphviz() {
        return `
                subgraph cluster_${this.id} {
                label = "${this.Name}"

                ${this.nodesDefinition()}
                ${this.graphvizChildren()}
               
            }
        `
    }

    nodesDefinition() {
        return `N${this.id} [ label=< ${this.nodeLabel()} > ];
        `
    }

    nodeLabel(): string {
        return `
        <table border="1" cellborder="1" cellspacing="0" cellpadding="4">
            <tr>
                <td bgcolor="lightgrey" colspan="4"><b>Variables</b></td>
            </tr>
           
            ${this.Vars.length > 0 ?
                ` 
                <th>
                    <td>Identificador</td>
                    <td>Tipo</td>
                    <td>Linea</td>
                    <td>Columna</td>
                </th>\n`+
                this.Vars.map(variable => `
                <tr>
                    <td>${variable.Name}</td>
                    <td>${variable.Type}</td>
                    <td>${variable.Line}</td>
                    <td>${variable.Column}</td>
                </tr>
            `).join('')
                : `<tr><td colspan="4"> - </td></tr>`
            }
            <tr>
                <td bgcolor="lightgrey" colspan="4"><b>Subrutinas</b></td>
            </tr>
            ${this.Funcs.length > 0
                ?
                ` 
                <tr>
                    <td>Identificador</td>
                    <td>Tipo</td>
                    <td>Linea</td>
                    <td>Columna</td>
                </tr>\n`+
                this.Funcs.map(subroutine => `
                <tr>
                    <td>${subroutine.Name}</td>
                    <td>${subroutine.Type}</td>
                    <td>${subroutine.Line}</td>
                    <td>${subroutine.Column}</td>
                </tr>
            `).join('')
                : `<tr><td colspan="4"> - </td></tr>`
            }
            <tr>
                <td bgcolor="lightgrey" colspan="4"><b>Estructuras</b></td>
            </tr>
            ${this.Structs.length > 0
                ?
                ` 
                <tr>
                    <td>Identificador</td>
                    <td>Tipo</td>
                    <td>Linea</td>
                    <td>Columna</td>
                </tr>\n`+
                this.Structs.map(struct => `
                <tr>
                    <td>${struct.Name}</td>
                    <td>${struct.Type}</td>
                    <td>${struct.Line}</td>
                    <td>${struct.Column}</td>
                </tr>
            `).join('')
                : `<tr><td colspan="4"> - </td></tr>`
            }
        </table>
        `
    }

    graphvizChildren(): string {
        let result = ''

        this.ChildScopes?.forEach((childScope, index) => {
            if (index < Scope.maxDepth) {
                const child = new Scope(childScope)
                result += `
                    ${child.graphviz()}
                    N${this.id} -> N${child.id} [ltail=cluster_${this.id} lhead=cluster_${child.id}]
                `
            }
        })

        return result
    }
}

export const graphvizReport = (symbolTable: SymbolTableI) => {

    Scope.idCounter = 0
    const globalScope = new Scope(symbolTable.GlobalScope)

    return `
            digraph G {

            node [shape=none];
            rankdir=TB;
                ${globalScope.graphviz()}
            }
        `
}