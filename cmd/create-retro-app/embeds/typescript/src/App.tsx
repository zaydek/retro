import * as React from "react"

import "./App.css"

interface HeaderProps {
	children?: React.ReactNode
}

function Header({ children }: HeaderProps): JSX.Element {
	return <h1>{children}</h1>
}

export default function App(): JSX.Element {
	return (
		<div className="App">
			<Header>Hello, world!</Header>
		</div>
	)
}
