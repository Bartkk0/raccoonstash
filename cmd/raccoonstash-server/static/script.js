let file = document.querySelector('#file')
let preview = document.querySelector('#filePreview')

file.addEventListener('input', e => {
    preview.src = URL.createObjectURL(e.target.files[0])
})