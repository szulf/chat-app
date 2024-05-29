async function idk(path) {
    fetch(path).then(
        val => {
            val.text().then(s => document.querySelector("main").innerHTML = s);
        },
    ).catch(
        err => {
            console.log("Error(1st fetch): ", err);
        }
    )
}

let request = window.location.href.substring(window.location.href.lastIndexOf("/"));
// idk("http://localhost:3000" + request);

//Scroll start at bottom
document.getElementById("test").scrollTop = document.getElementById("test").scrollHeight;