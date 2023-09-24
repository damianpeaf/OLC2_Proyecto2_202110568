import { Dialog, Transition } from '@headlessui/react'
import { Fragment, useEffect, useState } from 'react'
import { useTSwift } from '../../hooks'


export const RenameModal = () => {

  const { isRenameModalOpen, closeRenameModal, renameDocument, currentDocument } = useTSwift();

  const [fileName, setFileName] = useState('')

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    renameDocument(currentDocument.id, fileName)
    closeRenameModal()
  }

  useEffect(() => {
    if (isRenameModalOpen) {
      setFileName('')
    }
  }, [isRenameModalOpen])


  return (
    <Transition appear show={isRenameModalOpen} as={Fragment}>
      <Dialog as="div" className="relative z-50" onClose={closeRenameModal}>
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

        <div className="fixed inset-0 overflow-y-auto">
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
              <Dialog.Panel className="w-full max-w-md transform overflow-hidden rounded-2xl bg-slate-900 p-6 text-left align-middle shadow-xl transition-all">
                <Dialog.Title
                  as="h3"
                  className="text-lg font-medium leading-6 text-white"
                >
                  Cambiar nombre de archivo
                </Dialog.Title>
                <form onSubmit={handleSubmit}>
                  <div className="mt-2">
                    <input type="text"
                      value={fileName}
                      onChange={(e) => setFileName(e.target.value)}
                      className="w-full px-3 py-2 text-base text-white placeholder-gray-400 bg-slate-800 border border-transparent rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 focus:placeholder-gray-500 focus:bg-slate-900 focus:text-white sm:text-sm"
                    />
                  </div>
                  <div className="mt-4 flex gap-x-3">
                    <button
                      type="submit"
                      className="inline-flex justify-center rounded-md border border-transparent bg-green-900 px-4 py-2 text-sm font-medium text-green-200 hover:bg-green-800 focus:outline-none focus-visible:ring-2 focus-visible:ring-green-500 focus-visible:ring-offset-2"
                    >
                      Cambiar
                    </button>
                    <button
                      type="button"
                      className="inline-flex justify-center rounded-md border border-transparent bg-blue-950 px-4 py-2 text-sm font-medium text-blue-200 hover:bg-blue-900 focus:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 focus-visible:ring-offset-2"
                      onClick={closeRenameModal}
                    >
                      Cancelar
                    </button>
                  </div>
                </form>
              </Dialog.Panel>
            </Transition.Child>
          </div>
        </div>
      </Dialog>
    </Transition>
  )
}
