
import MonacoEditor from '@monaco-editor/react';
import { useEffect, useState } from 'react';
import { useTSwift } from '../../hooks';

export const C3D = () => {

    const { c3dContent, setC3DContent } = useTSwift();
    const [value, setValue] = useState('')


    useEffect(() => {
        setValue(c3dContent)
    }, [c3dContent])

    const handleOnChange = (value: string | undefined) => {
        setC3DContent(value || '')
    }

    return (
        <MonacoEditor
            value={value}
            theme="vs-dark"
            onChange={handleOnChange}
            options={{
                fontSize: "20px",
            }}
            language='c'
        />
    )
}
