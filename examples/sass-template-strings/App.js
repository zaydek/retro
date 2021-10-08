import sass from "sass-template-strings"

function ComponentA() {
	return (
		<>
			{sass`
				.ComponentA {
					color: red;
				}
			`}
			<h1 className="ComponentA">
				Hello, world!
			</h1>
		</>
	)
}

function ComponentB() {
	return (
		<>
			{sass`
				.ComponentB {
					color: blue;
				}
			`}
			<h1 className="ComponentB">
				Hello, world!
			</h1>
		</>
	)
}

export function App() {
	return (
		<>
			<ComponentA />
			<ComponentB />
		</>
	)
}
