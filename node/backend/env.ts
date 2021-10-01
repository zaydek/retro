// Lazily wraps `throw new Error(...)` because throws aren't legal expressions
function InternalError<Type>(returnType: Type): Type {
	throw new Error("Internal Error")
	return returnType
}

export const NODE_ENV = process.env["NODE_ENV"] ?? InternalError("")
export const RETRO_CMD = process.env["RETRO_CMD"] ?? InternalError("")
export const RETRO_WWW_DIR = process.env["RETRO_WWW_DIR"] ?? InternalError("")
export const RETRO_SRC_DIR = process.env["RETRO_SRC_DIR"] ?? InternalError("")
export const RETRO_OUT_DIR = process.env["RETRO_OUT_DIR"] ?? InternalError("")
