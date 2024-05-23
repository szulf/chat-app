async function idk(path) {
    fetch(path).then(
        val => {
            val.text().then(s => {document.getElementById("test").innerHTML = s;});
        },
    ).catch(
        err => {
            console.log("Error(1st fetch): ", err);
        }
    )
}

// idk("http://localhost:3000/fetch");