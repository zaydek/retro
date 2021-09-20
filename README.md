# Introducing Retro

What is Retro? Retro is a React toolkit for building client-side rendered (CSR) and server-side generated (SSG) web apps.

Retro is based on the philosophy that building web apps should be fast, simply, and fun.

In practice, Retro is a toolkit that sits on top of React and esbuild. You can use Retro to build client-side rendered (CSR) or server-side generated (SSG) web apps. Note that Retro does not support server-side rendering (SSR) as Retro is not currently designed to emit a binary; Retro is designed to emit only HTML / CSS / and JS files.

Furthermore, Retro includes a store and router implementation in the `@zaydek/retro-std` standard library. You don't have to use the library and React apps built without Retro can also use the library. The purpose of the library is to make Retro web apps more idiomatic.

**Retro is not for everyone or every project.** Retro simply attempts to solve for building client-side rendered (CSR) and sever-side generated (SSG) web apps. It is not a heavy-handed framework; think of Retro as a React toolkit.
