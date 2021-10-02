import nodeReadline from "readline"

export default (function readline(): (() => Promise<string>) {
	async function* createReadlineGenerator(): AsyncGenerator<string> {
		const nodeReadlineInterface = nodeReadline.createInterface({ input: process.stdin })
		for await (const line of nodeReadlineInterface) {
			yield line
		}
	}
	const generate = createReadlineGenerator()
	return async () => {
		const result = await generate.next()
		return result.value
	}
})()
