const sessionUser = requireAuth();
if (!sessionUser || !sessionUser.is_specialist) {
    window.location.href = "./index.html";
}

const currentSpecialistId = sessionUser?.id || null;

let poolTickets = [];
let myTickets = [];
let selectedTicketId = null;
let activeTab = "pool";

let currentChatId = null;
let chatMessages = [];
let chatSocket = null;

const searchInput = document.getElementById("searchInput");
const refreshBtn = document.getElementById("refreshBtn");
const poolTabBtn = document.getElementById("poolTabBtn");
const myTabBtn = document.getElementById("myTabBtn");
const ticketList = document.getElementById("ticketList");
const ticketDetails = document.getElementById("ticketDetails");
const assignBtn = document.getElementById("assignBtn");
const closeBtn = document.getElementById("closeBtn");

const poolCount = document.getElementById("poolCount");
const myCount = document.getElementById("myCount");
const newCount = document.getElementById("newCount");
const closedCount = document.getElementById("closedCount");

const chatStatus = document.getElementById("chatStatus");
const chatMessagesContainer = document.getElementById("chatMessages");
const chatInput = document.getElementById("chatInput");
const sendMessageBtn = document.getElementById("sendMessageBtn");

const logoutBtn = document.getElementById("logoutBtn");
const specialistWelcome = document.getElementById("specialistWelcome");
const specialistBadge = document.getElementById("specialistBadge");

if (specialistWelcome && sessionUser) {
    specialistWelcome.textContent = `Вы вошли как ${sessionUser.name || "специалист"}. Здесь отображаются все обращения и отдельная вкладка с назначенными на вас задачами.`;
}

if (specialistBadge && sessionUser) {
    specialistBadge.textContent = sessionUser.name || "Специалист";
}

logoutBtn?.addEventListener("click", logout);

function getCurrentTickets() {
    return activeTab === "pool" ? poolTickets : myTickets;
}

function getFilteredTickets() {
    const query = searchInput.value.trim().toLowerCase();
    const current = getCurrentTickets();

    if (!query) return current;

    return current.filter((ticket) =>
        (ticket.message || "").toLowerCase().includes(query) ||
        String(ticket.id).includes(query) ||
        (ticket.category || "").toLowerCase().includes(query) ||
        (ticket.userName || "").toLowerCase().includes(query) ||
        (ticket.specialistName || "").toLowerCase().includes(query)
    );
}

function updateStats() {
    poolCount.textContent = poolTickets.length;
    myCount.textContent = myTickets.length;
    newCount.textContent = poolTickets.filter((t) => t.status === "NEW").length;
    closedCount.textContent = poolTickets.filter((t) => t.status === "CLOSED").length;
}

function renderTabs() {
    poolTabBtn.classList.toggle("active", activeTab === "pool");
    myTabBtn.classList.toggle("active", activeTab === "my");
}

function renderTicketList() {
    const filtered = getFilteredTickets();
    ticketList.innerHTML = "";

    filtered.forEach((ticket) => {
        const div = document.createElement("div");
        div.className = "ticket-item" + (ticket.id === selectedTicketId ? " active" : "");
        div.onclick = async () => {
            selectedTicketId = ticket.id;
            renderTicketList();
            renderTicketDetails();
            await openChatForSelectedTicket();
        };

        div.innerHTML = `
            <div class="ticket-item-top">
                <div>
                    <div><strong>Обращение #${ticket.id}</strong></div>
                    <div class="ticket-meta">${escapeHtml(ticket.userName || "Пользователь")}</div>
                </div>
                <span class="badge ${categoryClass(ticket.category)}">${categoryLabel(ticket.category)}</span>
            </div>
            <div class="ticket-message">${escapeHtml(ticket.message)}</div>
            <div class="ticket-footer">
                <span class="badge ${statusClass(ticket.status)}">${statusLabel(ticket.status)}</span>
                <span class="ticket-meta">${ticket.createdAt || ""}</span>
            </div>
        `;

        ticketList.appendChild(div);
    });

    if (!filtered.find((t) => t.id === selectedTicketId)) {
        selectedTicketId = filtered[0]?.id || null;
    }
}

function renderTicketDetails() {
    const allTickets = [...poolTickets, ...myTickets];
    const ticket = allTickets.find((t) => t.id === selectedTicketId);

    if (!ticket) {
        ticketDetails.innerHTML = `<p class="muted">Выберите обращение из списка.</p>`;
        return;
    }

    ticketDetails.innerHTML = `
        <div style="margin-bottom: 16px; display: flex; gap: 8px; flex-wrap: wrap;">
            <span class="badge ${categoryClass(ticket.category)}">${categoryLabel(ticket.category)}</span>
            <span class="badge ${statusClass(ticket.status)}">${statusLabel(ticket.status)}</span>
        </div>

        <div class="detail-grid">
            <div class="detail-card">
                <span class="detail-label">Обращение</span>
                <span class="detail-value">#${ticket.id}</span>
            </div>
            <div class="detail-card">
                <span class="detail-label">Пользователь</span>
                <span class="detail-value">${escapeHtml(ticket.userName || "Пользователь")}</span>
            </div>
            <div class="detail-card">
                <span class="detail-label">Специалист</span>
                <span class="detail-value">${escapeHtml(ticket.specialistName || "Не назначен")}</span>
            </div>
            <div class="detail-card">
                <span class="detail-label">Создано</span>
                <span class="detail-value">${ticket.createdAt || "—"}</span>
            </div>
        </div>

        <div>
            <p class="detail-label">Сообщение</p>
            <div class="detail-message">${escapeHtml(ticket.message)}</div>
        </div>
    `;

    const isMine = ticket.specialistId === currentSpecialistId;
    const isAssigned = Boolean(ticket.specialistId);

    assignBtn.disabled = isMine || ticket.status === "CLOSED";
    closeBtn.disabled = !isMine || ticket.status === "CLOSED";

    if (!isAssigned) {
        assignBtn.textContent = "Взять в работу";
    } else if (isMine) {
        assignBtn.textContent = "Назначено на вас";
    } else {
        assignBtn.textContent = "Уже назначено";
        assignBtn.disabled = true;
    }
}

function renderChatMessages() {
    chatMessagesContainer.innerHTML = "";

    if (!currentChatId) {
        chatStatus.textContent = "Выберите обращение, чтобы открыть чат.";
        return;
    }

    chatStatus.textContent = "";

    if (!chatMessages.length) {
        chatMessagesContainer.innerHTML = `<div class="chat-placeholder">Сообщений пока нет.</div>`;
        return;
    }

    chatMessages.forEach((message) => {
        const div = document.createElement("div");
        div.className = `chat-message ${message.senderType}`;

        let authorLabel = "Система";
        if (message.senderType === "user") {
            authorLabel = "Пользователь";
        } else if (message.senderType === "specialist") {
            authorLabel = sessionUser.name || "Специалист";
        }

        div.innerHTML = `
            <div>${escapeHtml(message.body)}</div>
            <span class="chat-message-meta">${escapeHtml(authorLabel)} • ${message.createdAt || ""}</span>
        `;

        chatMessagesContainer.appendChild(div);
    });

    chatMessagesContainer.scrollTop = chatMessagesContainer.scrollHeight;
}

async function openChatForSelectedTicket() {
    if (!selectedTicketId) {
        currentChatId = null;
        chatMessages = [];
        renderChatMessages();
        closeChatSocket();
        return;
    }

    try {
        const openResp = await apiPost(`/tickets/${selectedTicketId}/chat/open`);
        currentChatId = openResp.chat.id;

        const messagesResp = await apiGet(`/chats/${currentChatId}/messages`);
        chatMessages = (messagesResp.messages || []).map(normalizeChatMessage);

        renderChatMessages();
        connectChatSocket(currentChatId);
    } catch (err) {
        console.error(err);
        chatStatus.textContent = "Не удалось открыть чат";
    }
}

function closeChatSocket() {
    if (chatSocket) {
        chatSocket.close();
        chatSocket = null;
    }
}

function connectChatSocket(chatId) {
    closeChatSocket();

    chatSocket = new WebSocket(getChatWsUrl(chatId));

    chatSocket.onmessage = (event) => {
        try {
            const payload = JSON.parse(event.data);

            if (payload.type === "message_created") {
                const message = normalizeChatMessage(payload.data);

                if (!chatMessages.some((m) => m.id === message.id)) {
                    chatMessages.push(message);
                    renderChatMessages();
                }
            }
        } catch (err) {
            console.error("ws parse error", err);
        }
    };

    chatSocket.onopen = () => {
        chatStatus.textContent = "";
    };
}

function sendChatMessage() {
    const body = chatInput.value.trim();

    if (!currentChatId) {
        alert("Сначала выберите обращение");
        return;
    }

    if (!body) {
        alert("Введите сообщение");
        return;
    }

    if (!chatSocket || chatSocket.readyState !== WebSocket.OPEN) {
        alert("Нет соединения с чатом");
        return;
    }

    chatSocket.send(JSON.stringify({
        type: "send_message",
        data: {
            sender_type: "specialist",
            sender_id: currentSpecialistId,
            body,
        },
    }));

    chatInput.value = "";
}

async function loadPoolTickets() {
    const data = await apiGet("/specialists/tickets");
    poolTickets = (data.tickets || []).map(normalizeTicket);
    await enrichTicketsWithNames(poolTickets);
}

async function loadMyTickets() {
    const data = await apiGet(`/specialists/tickets?specialist_id=${currentSpecialistId}`);
    myTickets = (data.tickets || []).map(normalizeTicket);
    await enrichTicketsWithNames(myTickets);
}

async function loadCurrentTab() {
    try {
        if (activeTab === "pool") {
            await loadPoolTickets();
        } else {
            await loadMyTickets();
        }

        updateStats();
        renderTabs();
        renderTicketList();
        renderTicketDetails();
        await openChatForSelectedTicket();
    } catch (err) {
        console.error(err);
        alert("Не удалось загрузить обращения");
    }
}

async function refreshAllStats() {
    try {
        await loadPoolTickets();
        await loadMyTickets();
        updateStats();
        renderTabs();
        renderTicketList();
        renderTicketDetails();
        await openChatForSelectedTicket();
    } catch (err) {
        console.error(err);
        alert("Не удалось загрузить обращения");
    }
}

async function assignTicket() {
    if (!selectedTicketId) {
        alert("Выберите обращение");
        return;
    }

    try {
        const data = await apiPost(`/tickets/${selectedTicketId}/assign`, {
            specialist_id: currentSpecialistId,
        });

        if (!data.ticket) {
            throw new Error("Некорректный ответ сервера");
        }

        activeTab = "my";
        await refreshAllStats();
        selectedTicketId = data.ticket.id;
        renderTicketList();
        renderTicketDetails();
        await openChatForSelectedTicket();
    } catch (err) {
        console.error(err);
        alert("Не удалось взять обращение в работу");
    }
}

async function closeTicket() {
    if (!selectedTicketId) {
        alert("Выберите обращение");
        return;
    }

    try {
        const data = await apiPost(`/tickets/${selectedTicketId}/close`);

        if (!data.ticket) {
            throw new Error("Некорректный ответ сервера");
        }

        await refreshAllStats();
        selectedTicketId = data.ticket.id;
        renderTicketList();
        renderTicketDetails();
        await openChatForSelectedTicket();
    } catch (err) {
        console.error(err);
        alert("Не удалось закрыть обращение");
    }
}

poolTabBtn.addEventListener("click", async () => {
    activeTab = "pool";
    selectedTicketId = null;
    await loadCurrentTab();
});

myTabBtn.addEventListener("click", async () => {
    activeTab = "my";
    selectedTicketId = null;
    await loadCurrentTab();
});

refreshBtn.addEventListener("click", refreshAllStats);
assignBtn.addEventListener("click", assignTicket);
closeBtn.addEventListener("click", closeTicket);
sendMessageBtn.addEventListener("click", sendChatMessage);

searchInput.addEventListener("input", () => {
    renderTicketList();
    renderTicketDetails();
});

refreshAllStats();