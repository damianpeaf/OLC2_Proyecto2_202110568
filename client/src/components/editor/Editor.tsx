
import MonacoEditor from '@monaco-editor/react';
import { useEffect, useState } from 'react';
import { useTSwift } from '../../hooks';

export const Editor = () => {

    const { currentDocument, saveDocument } = useTSwift();
    const [value, setValue] = useState('')


    useEffect(() => {
        setValue(currentDocument.content)
    }, [currentDocument.id, currentDocument.content])

    const handleOnChange = (value: string | undefined) => {
        saveDocument({
            ...currentDocument,
            content: value || '',
        })
    }

    return (
        <MonacoEditor
            value={value}
            theme="vs-dark"
            onChange={handleOnChange}
            options={{
                fontSize: "20px",
            }}
            language='swift'
        />
    )
}
