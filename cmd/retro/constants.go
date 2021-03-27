package retro

const (
	MODE_DIR  = 0755
	MODE_FILE = 0644
)

const (
	WWW_DIR = "www"
	SRC_DIR = "src"
	OUT_DIR = "out"
)

// Server-sent events stub
const devStub = `const dev=new EventSource("/~dev");dev.addEventListener("reload",()=>{localStorage.setItem("/~dev",""+Date.now()),window.location.reload()}),dev.addEventListener("error",e=>{try{console.error(JSON.parse(e.data))}catch{}}),window.addEventListener("storage",e=>{e.key==="/~dev"&&window.location.reload()});`
