package retro

const (
	// Permission bits for writing files
	permBitsFile = 0644

	// Permission bits for writing directories
	permBitsDirectory = 0755
)

// Server-sent events stub
const devStub = `<script type="module">const dev=new EventSource("/~dev");dev.addEventListener("reload",()=>{localStorage.setItem("/~dev",""+Date.now()),window.location.reload()}),dev.addEventListener("error",e=>{try{console.error(JSON.parse(e.data))}catch{}}),window.addEventListener("storage",e=>{e.key==="/~dev"&&window.location.reload()});</script>`
