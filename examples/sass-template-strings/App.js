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

// export function App() {
// 	return (
// 		<>
// 			{sass`
// 				.modalContainer {
// 					position: absolute;
// 					top: 0;
// 					right: 0;
// 					bottom: 0;
// 					left: 0;
// 					// background-color: hsla(0, 0, 0, 0.5);
// 				}
// 				.clickMeContainer {
// 					display: flex;
// 					flex-direction: row;
// 					justify-content: center;
// 					align-items: center;
// 					min-height: 100vh;
// 				}
// 			`}
// 		</>
// 	)
// }
