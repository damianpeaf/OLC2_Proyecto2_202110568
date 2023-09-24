import { Editor } from '../editor';
import { Tabs } from '../tabs';
import { Console } from '../console';
import { RenameModal, SymbolTableModal } from '../modal';
import { AstModal } from '../modal/AstModal';
import { Toaster } from 'react-hot-toast'

export const Main = () => {
    return (
        <main
            className="
            fixed
            left-16
            w-[calc(100vw-64px)]
            h-full
        "
        >
            <section
                className='
                    h-full
                    w-full
                    flex
                    flex-col
                    gap-y-1
                '
            >
                <Tabs />
                <article
                    className='h-full'
                >
                    <Editor />
                    <Console />
                </article>
            </section>
            <RenameModal />
            <AstModal />
            <SymbolTableModal />
            <Toaster
                position="top-center"
                toastOptions={{
                    duration: 2000
                }}
            />
        </main>
    )
}
