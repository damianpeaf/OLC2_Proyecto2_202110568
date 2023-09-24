import {
    DocumentFile,
    TSwiftError,
    TSwiftState,
    initialDocument
} from '.';
import { SymbolTableI } from '../components/modal';
import { v4 } from 'uuid'


export type TSwiftActionType =
    {
        type: 'new-tab',
    } |
    {
        type: 'open-terminal'
    } |
    {
        type: 'close-terminal'
    } |
    {
        type: 'close-tab',
        payload: {
            id: string
        }
    } |
    {
        type: 'rename-tab',
        payload: {
            id: string,
            name: string
        }
    } | {
        type: 'open-file',
        payload: {
            document: Omit<DocumentFile, 'id'>
        }
    } |
    {
        type: 'save-file',
        payload: {
            document: DocumentFile
        }
    } |
    {
        type: 'set-current-document',
        payload: {
            id: string
        }
    } | {
        type: 'open-rename-modal'
    } | {
        type: 'close-rename-modal'
    } | {
        type: 'set-terminal-content',
        payload: {
            content: string
        }
    } | {
        type: 'open-ast-modal'
    } | {
        type: 'close-ast-modal'
    } | {
        type: 'set-graphviz-content',
        payload: {
            content: string
        }
    } | {
        type: 'reset-graphviz-content',
    } | {
        type: 'set-errors',
        payload: {
            errors: TSwiftError[]
        }
    } | {
        type: 'set-symbol-table',
        payload: {
            content: SymbolTableI | null
        }
    } | {
        type: 'open-symbol-table-modal'
    } | {
        type: 'close-symbol-table-modal'
    }

export const TSwiftReducer = (state: TSwiftState, action: TSwiftActionType): TSwiftState => {
    switch (action.type) {
        case 'new-tab':

            const document: DocumentFile = {
                id: v4(),
                name: 'Untitled',
                content: '',
            }

            return {
                ...state,
                documents: [...state.documents, document],
                currentDocument: document
            }

        case 'open-terminal':
            return {
                ...state,
                isConsoleOpen: true
            }
        case 'close-terminal':
            return {
                ...state,
                isConsoleOpen: false
            }

        case 'close-tab':
            const emptyDoc = initialDocument()

            return {
                ...state,
                documents: state.documents.filter(doc => doc.id !== action.payload.id).length > 0
                    ? state.documents.filter(doc => doc.id !== action.payload.id)
                    : [emptyDoc],
                currentDocument: state.currentDocument?.id === action.payload.id
                    ? state.documents.find(doc => doc.id !== action.payload.id) || emptyDoc
                    : state.currentDocument
            }
        case 'rename-tab':
            return {
                ...state,
                documents: state.documents.map(doc => {
                    if (doc.id === action.payload.id) {
                        return {
                            ...doc,
                            name: action.payload.name
                        }
                    }
                    return doc;
                }),
                currentDocument: state.currentDocument?.id === action.payload.id
                    ? {
                        ...state.currentDocument,
                        name: action.payload.name
                    }
                    : state.currentDocument
            }

        case 'open-file':

            const doc: DocumentFile = {
                ...action.payload.document,
                id: v4()
            }

            return {
                ...state,
                documents: [...state.documents, doc],
                currentDocument: doc
            }

        case 'save-file':
            return {
                ...state,
                documents: state.documents.map(doc => {
                    if (doc.id === action.payload.document.id) {
                        return action.payload.document;
                    }
                    return doc;
                }),
                currentDocument: action.payload.document
            }

        case 'set-current-document':
            return {
                ...state,
                currentDocument: state.documents.find(doc => doc.id === action.payload.id) || state.currentDocument
            }
        case 'open-rename-modal':
            return {
                ...state,
                isRenameModalOpen: true
            }
        case 'close-rename-modal':
            return {
                ...state,
                isRenameModalOpen: false
            }
        case 'set-terminal-content':
            return {
                ...state,
                terminalContent: action.payload.content
            }
        case 'open-ast-modal':
            return {
                ...state,
                isAstModalOpen: true
            }

        case 'close-ast-modal':
            return {
                ...state,
                isAstModalOpen: false
            }
        case 'set-graphviz-content':
            return {
                ...state,
                graphviz: action.payload.content
            }
        case 'reset-graphviz-content':
            return {
                ...state,
                graphviz: null
            }
        case 'set-errors':
            return {
                ...state,
                errors: action.payload.errors
            }

        case 'set-symbol-table':
            return {
                ...state,
                symbolTable: action.payload.content
            }

        case 'open-symbol-table-modal':
            return {
                ...state,
                isSymbolTableModalOpen: true
            }

        case 'close-symbol-table-modal':
            return {
                ...state,
                isSymbolTableModalOpen: false
            }
        default:
            return state;
    }
}