const React = require("react")
const ReactDOM = require("react-dom")
const ReactDOMServer = require("react-dom/server")

const App = require("./out/.retro/App.js").default

console.log(ReactDOMServer.renderToString(React.createElement(App)))
