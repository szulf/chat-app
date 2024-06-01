async function logFetchData(host, path) {
    fetch(host+path)
        .then(val => {
            val.text().then((s) => (console.log(s)));
        })
        .catch(err => {
            console.log("Error(1st fetch): ", err);
        });
}

const request = window.location.href.substring(window.location.href.lastIndexOf("/"));

// Dont know if thats good
if (request === "/") {
    window.location.href = "/login";
}
logFetchData("http://localhost:3000", request);

const loginSubmit = document.getElementById("loginSubmit");
if (loginSubmit !== null) {
    loginSubmit.addEventListener("click", e => {
        const username = document.getElementsByName("username")[0].value;
        const password = document.getElementsByName("password")[0].value;

        fetch("http://localhost:3000/login", {
            method: "POST",
            body: JSON.stringify({
                "username": username,
                "password": password
            }),
            credentials: "include"
        }).then(val => {
            val.text().then(s => {
                console.log(s);
            });
        }).catch(err => {
            console.log("Error(login): ", err);
        });

        // Gotta check if this will save cookies
        // window.location.href = "/chats";
    });
}

const registerSubmit = document.getElementById("registerSubmit");
if (registerSubmit !== null) {
    registerSubmit.addEventListener("click", e => {
        const username = document.getElementsByName("username")[0].value;
        const password = document.getElementsByName("password")[0].value;
        const passwordConfirm = document.getElementsByName("passwordConfirm")[0].value;

        // Move to a different function later
        // username(4 chars, no spaces)
        let usernameRgx = /^(?!.*[ ])[a-zA-z\d@$!%*#?&]{4,}/;
        if (username.match(usernameRgx) === null) {
            console.log("username must be at least 4 characters long, and cant contain spaces");
            return;
        }
        // password(1 number, 1 lowercase, 1 uppercase, 1 special?, 8 chars)
        let pswrdRgx = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*#?&])[a-zA-Z\d@$!%*#?&]{8,}/;
        if (password.match(pswrdRgx) === null) {
            console.log("password must contain 1 uppercase letter, 1 lowercase letter, 1 digit and 1 special character(@$!%*#?&), and must be at least 8 characters long");
            return;
        }
        // Are passwords the same
        if (password !== passwordConfirm) {
            console.log("passwords not the same");
            return;
        }

        fetch("http://localhost:3000/login", {
            method: "POST",
            body: JSON.stringify({
                "username": username,
                "password": password
            }),
            credentials: "include"
        }).then(val => {
            val.text().then(s => {
                console.log(s);
            });
        }).catch(err => {
            console.log("Error(login): ", err);
        });

        // Gotta check if this will save cookies
        // window.location.href = "/chats";
    });
}

// Move to a different script later

// Scroll start at bottom
const test = document.getElementById("test");
if (test !== null) {
    test.scrollTop = document.getElementById("test").scrollHeight;
}
