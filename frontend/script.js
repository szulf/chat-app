async function fetchData(host, path, args, fn) {
    fetch(host+path, args)
        .then(val => {
            val.text().then(fn);
        })
        .catch(err => {
            console.log("Error(fetch): ", err);
        });
}

const host = "http://localhost:3000";
const request = window.location.href.substring(window.location.href.lastIndexOf("/"));

// gonna need to change that to check if user is logged in
// with a server request maybe??
// because cookie checking is on the server
// to like a "/loggedIn" endpoint
if (request === "/") {
    window.location.replace("/login");
}

fetchData(host, request, {}, (s) => {
    let ss = s.split("\n");
    let elem = document.getElementById(ss[0]);
    ss.shift();
    if (elem !== null) {
        if (elem.id === "errMsg") {
            elem.classList.remove("d-none");
        }
        elem.innerHTML = ss.join("");
    }

    if (request === "/login" || request === "/register") {
        const form = document.getElementById("Form");
        form.addEventListener("submit", e => {
            e.preventDefault();
            const username = document.getElementsByName("Username")[0].value;
            const password = document.getElementsByName("Password")[0].value;
            let passwordConfirm;
            if (request === "/register") {
                passwordConfirm = document.getElementsByName("PasswordConfirm")[0].value;
            }
            const errMsg = document.getElementById("errMsg");

            if (!errMsg.classList.contains("d-none")) {
                errMsg.classList.add("d-none");
            }

            let args;
            if (request === "/register") {
                args = {
                    method: "POST",
                    body: JSON.stringify({
                        "username": username,
                        "password": password,
                        "passwordConfirm": passwordConfirm,
                    }),
                    credentials: "include"
                }
            } else if (request === "/login") {
                args = {
                    method: "POST",
                    body: JSON.stringify({
                        "username": username,
                        "password": password,
                    }),
                    credentials: "include"
                }
            }

            fetchData(host, request, args, (s) => {
                let ss = s.split("\n");
                let elem = document.getElementById(ss[0]);
                ss.shift();
                if (elem !== null) {
                    if (elem.id === "errMsg") {
                        elem.classList.remove("d-none");
                    }
                    elem.innerHTML = ss.join("");
                }
            })
            window.history.replaceState(null, "", "/chats");
        });
    }
})

// Move to a different script later
// like a front end script
// which will also contain the thing with classes

// Scroll start at bottom
const test = document.getElementById("test");
if (test !== null) {
    test.scrollTop = document.getElementById("test").scrollHeight;
}
