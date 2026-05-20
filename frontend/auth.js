const loginTabBtn = document.getElementById("loginTabBtn");
const registerTabBtn = document.getElementById("registerTabBtn");

const loginForm = document.getElementById("loginForm");
const registerForm = document.getElementById("registerForm");
const authMessage = document.getElementById("authMessage");

const loginEmailInput = document.getElementById("loginEmail");
const loginPasswordInput = document.getElementById("loginPassword");

const registerNameInput = document.getElementById("registerName");
const registerEmailInput = document.getElementById("registerEmail");
const registerPasswordInput = document.getElementById("registerPassword");

function setAuthMessage(text, isError = false) {
    authMessage.textContent = text;
    authMessage.classList.toggle("error", isError);
    authMessage.classList.toggle("success", !isError && Boolean(text));
}

function showLoginTab() {
    loginTabBtn.classList.add("active");
    registerTabBtn.classList.remove("active");
    loginForm.classList.remove("hidden");
    registerForm.classList.add("hidden");
    setAuthMessage("");
}

function showRegisterTab() {
    registerTabBtn.classList.add("active");
    loginTabBtn.classList.remove("active");
    registerForm.classList.remove("hidden");
    loginForm.classList.add("hidden");
    setAuthMessage("");
}

function getSafeUserName(user) {
    return user?.name || "Пользователь";
}

function redirectByRole(user) {
    if (user?.is_specialist) {
        window.location.href = "./specialist.html";
        return;
    }

    window.location.href = "./user.html";
}

async function login(event) {
    event.preventDefault();

    const email = loginEmailInput.value.trim();
    const password = loginPasswordInput.value.trim();

    if (!email || !password) {
        setAuthMessage("Заполните email и пароль", true);
        return;
    }

    try {
        const data = await apiPost("/auth/login", {
            email,
            password,
        });

        const user = data.user;
        if (!user) {
            throw new Error("Сервер не вернул пользователя");
        }

        saveSession(user);
        setAuthMessage(`Добро пожаловать, ${getSafeUserName(user)}`);

        redirectByRole(user);
    } catch (err) {
        console.error(err);
        setAuthMessage("Не удалось выполнить вход", true);
    }
}

async function register(event) {
    event.preventDefault();

    const name = registerNameInput.value.trim();
    const email = registerEmailInput.value.trim();
    const password = registerPasswordInput.value.trim();

    if (!name || !email || !password) {
        setAuthMessage("Заполните имя, email и пароль", true);
        return;
    }

    try {
        const data = await apiPost("/auth/register", {
            name,
            email,
            password,
        });

        const user = data.user;
        if (!user) {
            throw new Error("Сервер не вернул пользователя");
        }

        saveSession(user);
        setAuthMessage(`Пользователь ${getSafeUserName(user)} зарегистрирован`);

        redirectByRole(user);
    } catch (err) {
        console.error(err);
        setAuthMessage("Не удалось зарегистрироваться", true);
    }
}

loginTabBtn.addEventListener("click", showLoginTab);
registerTabBtn.addEventListener("click", showRegisterTab);
loginForm.addEventListener("submit", login);
registerForm.addEventListener("submit", register);

showLoginTab();