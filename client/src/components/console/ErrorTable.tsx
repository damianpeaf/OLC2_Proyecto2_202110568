import { TSwiftError } from "../../context"

type ErrorTableProps = {
    errors: TSwiftError[]
}

export const ErrorTable = ({ errors }: ErrorTableProps) => {

    return (

        <div className="relative overflow-x-auto shadow-md rounded-lg">
            <table className="w-full text-sm text-left  text-gray-400">
                <thead className="text-xs  uppercase  bg-gray-700 text-gray-400">
                    <tr>
                        <th scope="col" className="px-6 py-3">
                            #
                        </th>
                        <th scope="col" className="px-6 py-3">
                            Tipo de error
                        </th>
                        <th scope="col" className="px-6 py-3">
                            Descripción
                        </th>
                        <th scope="col" className="px-6 py-3">
                            Linea
                        </th>
                        <th scope="col" className="px-6 py-3">
                            Columna
                        </th>
                    </tr>
                </thead>
                <tbody>
                    {
                        errors.map((error, index) => (
                            <tr className=" border-b bg-gray-900 border-gray-700" key={index}>

                                <td className="px-6 py-4">
                                    {index + 1}
                                </td>
                                <td className="px-6 py-4">
                                    {error.Type}
                                </td>
                                <td className="px-6 py-4">
                                    {error.Msg}
                                </td>
                                <td className="px-6 py-4">
                                    {error.Line}
                                </td>
                                <td className="px-6 py-4">
                                    {error.Column}
                                </td>
                            </tr>
                        ))
                    }
                </tbody>
            </table>
        </div>

    )
}
