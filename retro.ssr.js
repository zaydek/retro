const React = require("react")
const ReactDOMServer = require("react-dom/server")

// After `retro build`, run `node retro.ssr.js` to render `src/App.js` as a
// string on the server (your computer).
//
// The name of this script is arbitrary. Modify this script to support static-
// site generation (SSG) or server-side rendering (SSR) more generally.
const App = require("./out/.retro/App.js").default

console.log(ReactDOMServer.renderToString(React.createElement(App)))
// <div class="App" data-reactroot=""><h1>Hello, world!</h1></div>
