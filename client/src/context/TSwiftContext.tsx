import { createContext } from 'react';
import { TSwiftActionType, TSwiftState } from '.';

interface ContextProps extends TSwiftState {
    dispatch: React.Dispatch<TSwiftActionType>,
    saveState: () => void
}

export const TSwiftContext = createContext({} as ContextProps);