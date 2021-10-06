export const useLayoutEffectSSR = typeof window === "undefined" ? React.useEffect : React.useLayoutEffect
