import { Main, SideBar } from './components'
import { TSwiftProvider } from './context'

export function App() {

  return (
    <TSwiftProvider>
      <div>
        <SideBar />
        <Main />
      </div>
    </TSwiftProvider>
  )
}

export default App
