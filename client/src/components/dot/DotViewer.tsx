import { Graphviz } from "graphviz-react";
import { DotViewerProps } from ".";
import { useState, useEffect } from 'react';



export const DotViewer = ({ dot }: DotViewerProps) => {

    const [windowSize, setWindowSize] = useState<{ width: number; height: number }>({
        width: window.innerWidth,
        height: window.innerHeight,
    });

    useEffect(() => {
        const handleResize = () => {
            setWindowSize({
                width: window.innerWidth,
                height: window.innerHeight,
            });
        };

        window.addEventListener("resize", handleResize);
        window.addEventListener("load", handleResize);
        return () => {
            window.removeEventListener("resize", handleResize);
            window.removeEventListener("load", handleResize);
        };
    }, []);

    return (
        <Graphviz
            dot={dot}
            options={{
                zoom: true,
                useWorker: false,
                engine: "dot",
                height: (windowSize.height * 0.8),
                width: (windowSize.width * 0.8),
                fit: true,
                scale: 1,
            }}
        />
    )
};
