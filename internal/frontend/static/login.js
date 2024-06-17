function toggleForms() {
    const loginBox = document.getElementById('login-box');
    const registerBox = document.getElementById('register-box');

    if (loginBox.style.display === "none") {
        loginBox.style.display = "block";
        registerBox.style.display = "none";
    } else {
        loginBox.style.display = "none";
        registerBox.style.display = "block";
    }
}

function login() {
    const URL = "http://localhost:8080/login";
    let user = {
        username: document.getElementById("login-email").value,
        password: document.getElementById("login-password").value
    }
    
    let request = new Request(URL, {
        body: JSON.stringify(user),
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        }
    });
    fetch(request).then(response => {
        if (!response.ok) {
            throw new Error("Network resposne was not ok")
        } else {
            window.location.href = "http://localhost:8080/tasks/";
        }
    })
    .catch(() => {
        alert("Fehler bei der Anmeldung!")
    });
    
}

function register() {
    const URL = "http://localhost:8080/register";
    let user = {
        username: document.getElementById("register-email").value,
        password: document.getElementById("register-password").value
    }
    
    let request = new Request(URL, {
        body: JSON.stringify(user),
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        }
    });
    fetch(request).then(response => {
        if (!response.ok) {
            throw new Error("Network resposne was not ok");
        } else {
            window.location.href = "http://localhost:8080/tasks/";
        }
    })
    .catch(() => {
        alert("Fehler bei der Anmeldung!");
    });
    
}
  