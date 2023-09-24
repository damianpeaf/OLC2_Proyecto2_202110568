import { Menu } from '@headlessui/react'
import { AiFillFileText, AiFillPlusCircle, AiFillSave, AiFillFolderOpen, AiOutlineAppstoreAdd } from "react-icons/ai";
import { BiRename } from "react-icons/bi";

import { SideBarItem } from './';
import { useTSwift } from '../../hooks';
import { useRef } from 'react';

export const FileOptions = () => {

    const openInputRef = useRef<HTMLInputElement>(null);
    const openFormRef = useRef<HTMLFormElement>(null);
    const { newDocument, openRenameModal, openDocument, currentDocument } = useTSwift();

    const handleOpenFile = () => {
        openInputRef.current?.click()
    }

    const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.item(0);
        if (file) {
            const formatedFilename = file.name.split('.').slice(0, -1).join('.');

            const reader = new FileReader();
            reader.onload = () => {
                if (typeof reader.result === 'string') {
                    openDocument({
                        name: formatedFilename,
                        content: reader.result
                    })
                }
            }
            reader.readAsText(file)
        }

        openFormRef.current?.reset();
    }

    const handleSaveFile = () => {
        const blob = new Blob([currentDocument.content], { type: 'text/plain' });

        const link = document.createElement('a');
        link.href = window.URL.createObjectURL(blob);
        link.download = `${currentDocument.name}.swift`;
        link.click();

        window.URL.revokeObjectURL(link.href);
    }

    const fileOptions = [
        {
            name: 'Renombrar',
            icon: <BiRename />,
            onclick: openRenameModal
        },
        {
            name: 'Nuevo Archivo',
            icon: <AiOutlineAppstoreAdd />,
            onclick: newDocument
        }, {
            name: 'Abrir archivo',
            icon: <AiFillFolderOpen />,
            onclick: handleOpenFile
        }, {
            name: 'Guardar Archivo',
            icon: <AiFillSave />,
            onclick: handleSaveFile
        },
    ]

    return (
        <div className='relative mx-auto'>
            <Menu>
                <Menu.Button >
                    <SideBarItem icon={<AiFillFileText />} label="Opciones de archivos" />
                </Menu.Button>
                <div className='absolute top-1/2 left-[50px] p-2 w-[175px] rounded-xl'>
                    <Menu.Items>
                        {
                            fileOptions.map((option, index) => (
                                <Menu.Item key={index}>
                                    {({ active }) => (
                                        <div
                                            className="
                                                bg-secondary-dark
                                            "
                                        >
                                            <button
                                                className="
                                                flex
                                                items-center
                                                w-full
                                                p-2
                                                text-white
                                                hover:bg-text-dark-theme-dark
                                                "
                                                onClick={option.onclick}
                                            >
                                                {option.icon}
                                                <span className='ml-2'>{option.name}</span>
                                            </button>
                                        </div>
                                    )}
                                </Menu.Item>
                            ))
                        }
                    </Menu.Items>
                    <form ref={openFormRef}>
                        <input
                            type="file"
                            ref={openInputRef}
                            className='invisible'
                            onChange={handleFileChange}
                            accept='.swift'
                            multiple={false}
                        />
                    </form>
                </div>
            </Menu>
        </div>
    )
}