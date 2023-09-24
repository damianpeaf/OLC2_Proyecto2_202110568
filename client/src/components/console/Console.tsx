import { MdOutlineCloseFullscreen } from 'react-icons/md';
import { useTSwift } from '../../hooks';
import React, { useState, useEffect, useRef, MouseEventHandler } from 'react';
import { ErrorTable } from './ErrorTable';

export const Console = () => {

    const { isConsoleOpen, closeTerminal, terminalContent, errors } = useTSwift()

    const [height, setHeight] = useState(300); // Initial height
    const [isResizing, setIsResizing] = useState(false)

    const onMouseDown = (e: React.MouseEvent<HTMLDivElement>) => {
        const startY = e.clientY;
        const startHeight = height;

        const onMouseMove = (e: MouseEvent) => {
            const diffY = startY - e.clientY; // Calculate difference in reverse
            const newHeight = startHeight + diffY // Limit minimum height
            setHeight(newHeight);
            setIsResizing(true)
        };

        const onMouseUp = () => {
            document.removeEventListener('mousemove', onMouseMove);
            document.removeEventListener('mouseup', onMouseUp);
            setIsResizing(false)
        };

        document.addEventListener('mousemove', onMouseMove);
        document.addEventListener('mouseup', onMouseUp);
    };

    return (
        <section
            className={`
            ${isConsoleOpen ? 'scale-100' : 'scale-0'}
            ${isResizing ? 'select-none' : 'select-auto'}
            absolute
            bottom-0
            left-0
            w-full
            bg-background-dark
            transition-all
            overflow-auto
            resize-y
            cursor-n-resize
            `}
            style={{ height: `${height}px` }}

        >
            <article
                className="
                flex
                justify-between
                px-4
                mt-2
                "
                onMouseDown={onMouseDown}
            >
                <h2
                    className="
                    text-gray-300
                    font-bold
                    text-xl
                    "
                >
                    Consola
                </h2>
                <button
                    className="
                        p-2
                        text-gray-300
                        font-bold
                        text-xl
                    "
                    onClick={closeTerminal}
                >
                    <MdOutlineCloseFullscreen />
                </button>
            </article>
            <article
                className="
                    py-2
                    px-4
                    overflow-y-auto
                    mb-2
                    h-4/5
                    console-font
                    text-gray-400
                "
            >


                {
                    errors.length > 0 &&
                    <>
                        <br />
                        <h3
                            className='text-gray-300 font-bold text-xl'
                        >Errores:</h3>
                        <br />
                        <ErrorTable errors={errors} />
                        <br />
                        <br />
                        <h3
                            className='text-gray-300 font-bold text-xl'
                        >Salida:</h3>
                        <br />
                    </>
                }

                {
                    terminalContent.split('\n').map((line, index) => (
                        <React.Fragment key={index}>
                            <pre className='whitespace-pre-wrap text-xl'>

                                {line.trim() === '' ? <br /> : line.replace(/\t/g, '\u00a0\u00a0\u00a0\u00a0')}
                            </pre>
                        </React.Fragment>
                    ))
                }
            </article>
        </section>
    )
}
