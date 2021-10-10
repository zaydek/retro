import sass from "sass-template-strings"

export function App() {
	return (
		<>
			{sass`
				.App h1 {
					&::before {
						content: "# ";
						color: #2979ff;
					}
				}
			`}
			<div className="App">
				<h1>
					Hello, world!
				</h1>
			</div>
		</>
	)
}
