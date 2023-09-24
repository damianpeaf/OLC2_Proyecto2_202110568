import { Dialog, Transition } from '@headlessui/react'
import { Fragment, useState, useEffect } from 'react';
import { useTSwift } from '../../hooks'
import RViewerJS from 'viewerjs-react'
import 'viewerjs-react/dist/index.css'

export const AstModal = () => {

  const { isAstModalOpen, closeAstModal, graphviz } = useTSwift();

  const [b64Svg, setB64Svg] = useState<string | null>(null)

  useEffect(() => {

    if (graphviz) setB64Svg(btoa(unescape(encodeURIComponent(graphviz))))
    else setB64Svg(null)

  }, [graphviz])

  return (
    <Transition appear show={isAstModalOpen} as={Fragment}>
      <Dialog as="div" className="relative z-50" onClose={closeAstModal}>
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

        <div className="fixed inset-0 overflow-y-auto overflow-x-hidden">
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
              <Dialog.Panel className="w-full max-w-7xl h-full max-h-screen overflow-y-auto overflow-x-hidden transform overflow-hidden rounded-2xl bg-slate-900 p-6 text-left align-middle shadow-xl transition-all">
                <Dialog.Title
                  as="h3"
                  className="text-lg font-medium leading-6 text-white"
                >
                  Reporte CST
                </Dialog.Title>
                {
                  b64Svg ?
                    <div className="w-full h-full flex justify-center items-center min-h-[80vh] bg-white">
                      <style>
                        {`
                        .viewer-backdrop {
                          background-color: white !important;
                        }
                        `}
                      </style>
                      {/* @ts-ignore */}
                      <RViewerJS
                      >
                        <img src={`data:image/svg+xml;base64,${b64Svg}`} alt="AST" className="w-full h-full bg-white fill-white" />
                      </RViewerJS>
                    </div>
                    :
                    <div className="mt-4">
                      <p className="text-lg text-gray-300">
                        No se ha generado el reporte CST
                      </p>
                    </div>
                }
                <div className="mt-4 flex gap-x-3">
                  <button
                    type="button"
                    className="inline-flex justify-center rounded-md border border-transparent bg-blue-950 px-4 py-2 text-sm font-medium text-blue-200 hover:bg-blue-900 focus:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 focus-visible:ring-offset-2"
                    onClick={closeAstModal}
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
