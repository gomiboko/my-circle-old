const http = require('http')
const fs = require('fs')

const server = http.createServer((req, res) => {
  fs.readFile('csrf-page.html', 'utf-8', (err, data) => {
    res.writeHead(200, {'Content-Type': 'text/html'})
    res.write(data)
    res.end()
  })
})

server.listen(3000, 'localhost', () => {
  console.log(`Server running at http://localhost:3000/`)
})
