## v0.0.43 (April 11, 2021)

- Lowered esbuild artifacts from `esnext` to `es2018` for better backwards compatibility.
  In the future this will likely be further lowered to `es2017`.

## v0.0.42 (April 10, 2021)

- Normalized use of single-quotes `'` to backticks ``` ` ```.

## v0.0.39 to v0.0.41 migration guide

Upgrading to v0.0.41 from v0.0.39 can be done in a few steps.

- Upgrade dependencies to `@zaydek/retro@0.0.41` and `esbuild@0.11.6`
- Update `www/index.html` to reference the new bundle aliases:

```diff
<!DOCTYPE html>
<html lang="en">

<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>Hello, world!</title>
- <link rel="stylesheet" href="/index.css" />
+ <link rel="stylesheet" href="/bundle.css" />
</head>

<body>
  <div id="root"></div>
- <script src="/react.js"></script>
+ <script src="/vendor.js"></script>
- <script src="/index.js"></script>
+ <script src="/bundle.js"></script>
</body>

</html>
```

Note that you don’t need to rename source files; this is the purpose of bundle aliasing.

## v0.0.41 (April 9, 2021)

- Fixed a bug that lead the `build` command to panic.

## v0.0.40 (April 9, 2021)

- Added Sass and MDX templates.

  You can quickly experiment with Sass and MDX templates by using the command `--template=sass` or `--template=mdx`. For example, `npx @zaydek/create-retro-app app-name --template=sass`.

- Removed `import ReactDOM from "react-dom"` from all template `index.js` files.

  This import is not necessary because `React` and `ReactDOM` are automatically bundled via a shim for convenience. This has always been the case, but now this is reflected in the template files. Note that importing `React` or `ReactDOM` should otherwise be idempotent.

- Added a minimal `retro.config.js` to the root directory of every template. This file is recommended but not required.

  This configuration file tells esbuild to target `"es2017"`, which has better backwards compatibility.

  For example:

  ```js
  // https://esbuild.github.io/api/#build-api

  module.exports = {
    target: ["es2017"],
  }
  ```

- Deprecated the TypeScript template. ([#6](https://github.com/zaydek/retro/issues/6))

  The TypeScript template was deprecated because esbuild can parse TypeScript, therefore a TypeScript template is largely not needed because JavaScript is preferred as the base template. You can still use TypeScript of course, but you don’t need to convert every component to be typed in order to do so. You can simply rename a file from `Component.js` to `Component.tsx`. The referring import statement does not need to be changed because ES Modules imports do not use extensions.

- Fixed a bug that caused Retro to panic on an already binded port.

  Retro now cycles ports until an unbinded port is found.

- Upgraded logging to resemble Create React App experience.

  Previously Retro logged serve events and build errors to the terminal. This behavior has been changed to build success messages and build errors. Build errors are still propagated to the browser build success and build errors. a singular build success message.

- Upgraded esbuild to `0.11.6` and added support bundle aliasing. ([#1](https://github.com/zaydek/retro/issues/1))

  `src/index.js` now aliases to `out/bundle.js`, `src/index.css` now aliases to `out/bundle.css`, and `React` and `ReactDOM` now alias to `out/vendor.js`.

## v0.0.39 (April 4, 2021)

- Improve error messaging for `@zaydek/create-retro-app`
