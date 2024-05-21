async function idk(path) {
    fetch(path).then(
        val => {
            val.text().then(s => {document.getElementById("test").innerHTML = s});
        },
    ).catch(
        err => {
            console.log("Error: ", err);
        }
    )
}

let path = "http://localhost:3000/";
idk(path);