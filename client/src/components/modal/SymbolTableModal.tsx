import { Dialog, Transition } from '@headlessui/react'
import { Fragment, useState, useEffect } from 'react';
import { useTSwift } from '../../hooks'
import { DotViewer } from '../dot';
import { graphvizReport } from './';




export const SymbolTableModal = () => {

  const { isSymbolTableModalOpen, closeSymbolTableModal, symbolTable } = useTSwift();

  const [graphiz, setGraphiz] = useState<string | null>(null)

  useEffect(() => {
    symbolTable && setGraphiz(graphvizReport(symbolTable))
  }, [symbolTable])

  return (
    <Transition appear show={isSymbolTableModalOpen} as={Fragment}>
      <Dialog as="div" className="relative z-50" onClose={closeSymbolTableModal}>
        <Transition.Child
          as={Fragment}
          enter="ease-out duration-300"
          enterFrom="opacity-0"
          enterTo="opacity-100"
          leave="ease-in duration-200"
          leaveFrom="opacity-100"
          leaveTo="opacity-0"
        >
          <div className="fixed inset-0 bg-black bg-opacity-25" />
        </Transition.Child>

        <div className="fixed inset-0  overflow-y-auto overflow-x-hidden">
          <div className="flex min-h-full items-center justify-center p-4 text-center">
            <Transition.Child
              as={Fragment}
              enter="ease-out duration-300"
              enterFrom="opacity-0 scale-95"
              enterTo="opacity-100 scale-100"
              leave="ease-in duration-200"
              leaveFrom="opacity-100 scale-100"
              leaveTo="opacity-0 scale-95"
            >
              <Dialog.Panel className="w-full max-w-7xl h-full max-h-screen  overflow-y-auto overflow-x-hidden transform overflow-hidden rounded-2xl bg-slate-900 p-6 text-left align-middle shadow-xl transition-all">
                <Dialog.Title
                  as="h3"
                  className="text-lg font-medium leading-6 text-white"
                >
                  Reporte Tabla de Simbolos
                </Dialog.Title>
                {
                  graphiz ?
                    <div className="w-full h-full flex justify-center items-center min-h-[80vh]">
                      <DotViewer dot={graphiz} />
                    </div>
                    :
                    <div className="mt-4">
                      <p className="text-lg text-gray-300">
                        No se ha generado el reporte la tabla de simbolos
                      </p>
                    </div>
                }
                <div className="mt-4 flex gap-x-3">
                  <button
                    type="button"
                    className="inline-flex justify-center rounded-md border border-transparent bg-blue-950 px-4 py-2 text-sm font-medium text-blue-200 hover:bg-blue-900 focus:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 focus-visible:ring-offset-2"
                    onClick={closeSymbolTableModal}
                  >
                    Cerrar
                  </button>
                </div>
              </Dialog.Panel>
            </Transition.Child>
          </div>
        </div>
      </Dialog>
    </Transition>
  )
}
