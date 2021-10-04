const React = require("react")
const ReactDOMServer = require("react-dom/server")

// Run `node retro.ssr.js` after `retro build` to render `src/App.js` on the
// server. Modify this script to support static-site generation (SSG) or server-
// side rendering (SSR) more generally.
const App = require("./out/.retro/App.js").default

console.log(ReactDOMServer.renderToString(React.createElement(App)))
// <div class="App" data-reactroot=""><h1>Hello, world!</h1></div>
