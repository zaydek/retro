# Introducing Retro

What is Retro? Retro is a React toolkit for building client-side rendered (CSR) and server-side generated (SSG) web apps.

Retro is based on the philosophy that building web apps should be fast, simple, and fun.

In practice, Retro is a React toolkit that builds on top of Go, Node, React, and esbuild. You can use Retro to build client-side rendered (CSR) or server-side generated (SSG) web apps. Note that Retro does not support server-side rendering (SSR) as Retro is not currently designed to emit a binary or a runtime; Retro is designed to only emit HTML, CSS, and JS files. The only difference being whether these files are interpreted eagerly (SSG) or lazily (CSR).

Furthermore, Retro ships its own small standard library to make it easier to get up and running. The standard library currently includes a store and router implementation, designed to make Retro source idiomatic.

**Retro is not for everyone or every project.** Retro only solves for client-side rendered (CSR) and server-side generated (SSG) web apps.

## Get Started

To create a Retro web app, simply run the following command:

```sh
npx @zaydek/create-retro-app my-retro-app
```

This will create a new React app at the directory `<dir>` or in the current directory if omitted.

To start your app :

```sh
retro dev
```

To build your app for production:

```sh
retro build
```

## esbuild-style Configuration

Retro leverages esbuild to bundle apps for both development and production. Retro's build process can be configured using esbuild-style configuration. Simply add `retro.config.js` at the root directory of your project. This makes Retro apps extensible and can be customizable on a per-app basis.

Some examples of how configuration can support your development include:

- Adding support for HTML imports, such as `.svg`
- Adding support for build-time CSS tooling, such as Sass
- Changing the JavaScript lowering target

## Automatic TypeScript Transpilation

As Retro is built on top of esbuild, esbuild transpiles JavaScript React, TypeScript, and TypeScript React source code on-demand. Note that type-checking is not performed on your source code and additional tooling is needed to support this use-case. That being said, you can mix-and-match JavaScript and TypeScript source code. This is the preferred method for authoring complex apps. You don't need to choose a JavaScript or TypeScript template to get started and you won't need to refactor to 100% JavaScript or 100% TypeScript once you've started.

## License

Retro is licensed as [MIT open source](/LICENSE).
