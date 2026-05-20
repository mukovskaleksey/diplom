const sessionUser = requireAuth();
if (!sessionUser || sessionUser.is_specialist) {
    window.location.href = "./index.html";
}

const currentUserId = sessionUser?.id || null;

let tickets = [];
let selectedTicketId = null;

let currentChatId = null;
let chatMessages = [];
let chatSocket = null;

const messageInput = document.getElementById("message");
const searchInput = document.getElementById("searchInput");
const createTicketBtn = document.getElementById("createTicketBtn");
const refreshBtn = document.getElementById("refreshBtn");
const ticketList = document.getElementById("ticketList");
const ticketDetails = document.getElementById("ticketDetails");

const chatStatus = document.getElementById("chatStatus");
const chatMessagesContainer = document.getElementById("chatMessages");
const chatInput = document.getElementById("chatInput");
const sendMessageBtn = document.getElementById("sendMessageBtn");

const logoutBtn = document.getElementById("logoutBtn");
const userWelcome = document.getElementById("userWelcome");
const userBadge = document.getElementById("userBadge");

if (userWelcome && sessionUser) {
    userWelcome.textContent = `Вы вошли как ${sessionUser.name || "пользователь"}. Здесь можно создавать обращения, просматривать свои тикеты и вести чат со специалистом.`;
}

if (userBadge && sessionUser) {
    userBadge.textContent = sessionUser.name || "Пользователь";
}

logoutBtn?.addEventListener("click", logout);

function getFilteredTickets() {
    const query = searchInput.value.trim().toLowerCase();
    if (!query) return tickets;

    return tickets.filter((ticket) =>
        (ticket.message || "").toLowerCase().includes(query) ||
        String(ticket.id).includes(query) ||
        (ticket.category || "").toLowerCase().includes(query)
    );
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
                    <div class="ticket-meta">${escapeHtml(ticket.userName || sessionUser.name || "Пользователь")}</div>
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
    const ticket = tickets.find((t) => t.id === selectedTicketId);

    if (!ticket) {
        ticketDetails.innerHTML = `<p class="muted">Выберите тикет из списка.</p>`;
        return;
    }

    ticketDetails.innerHTML = `
        <div style="margin-bottom: 16px; display: flex; gap: 8px; flex-wrap: wrap;">
            <span class="badge ${categoryClass(ticket.category)}">${categoryLabel(ticket.category)}</span>
            <span class="badge ${statusClass(ticket.status)}">${statusLabel(ticket.status)}</span>
        </div>

        <div class="detail-grid">
            <div class="detail-card">
                <span class="detail-label">Пользователь</span>
                <span class="detail-value">${escapeHtml(ticket.userName || sessionUser.name || "Пользователь")}</span>
            </div>
            <div class="detail-card">
                <span class="detail-label">Специалист</span>
                <span class="detail-value">${escapeHtml(ticket.specialistName || "Не назначен")}</span>
            </div>
            <div class="detail-card">
                <span class="detail-label">Категория</span>
                <span class="detail-value">${categoryLabel(ticket.category)}</span>
            </div>
            <div class="detail-card">
                <span class="detail-label">Статус</span>
                <span class="detail-value">${statusLabel(ticket.status)}</span>
            </div>
        </div>

        <div>
            <p class="detail-label">Сообщение</p>
            <div class="detail-message">${escapeHtml(ticket.message)}</div>
        </div>
    `;
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
            authorLabel = sessionUser.name || "Пользователь";
        } else if (message.senderType === "specialist") {
            authorLabel = "Специалист";
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
            sender_type: "user",
            sender_id: currentUserId,
            body,
        },
    }));

    chatInput.value = "";
}

async function loadUserTickets() {
    if (!currentUserId) {
        tickets = [];
        selectedTicketId = null;
        renderTicketList();
        renderTicketDetails();
        currentChatId = null;
        chatMessages = [];
        renderChatMessages();
        closeChatSocket();
        return;
    }

    try {
        const data = await apiGet(`/users/${currentUserId}/tickets`);
        tickets = (data.tickets || []).map(normalizeTicket);
        await enrichTicketsWithNames(tickets);

        selectedTicketId = tickets[0]?.id || null;
        renderTicketList();
        renderTicketDetails();
        await openChatForSelectedTicket();
    } catch (err) {
        console.error(err);
        alert("Не удалось загрузить обращения пользователя");
    }
}

async function createTicket() {
    const message = messageInput.value.trim();

    if (!message) {
        alert("Введите сообщение");
        return;
    }

    try {
        const data = await apiPost("/tickets", {
            user_id: currentUserId,
            message,
        });

        if (!data.ticket) {
            throw new Error("Некорректный ответ сервера");
        }

        const createdTicket = normalizeTicket(data.ticket);
        await enrichTicketsWithNames([createdTicket]);

        tickets.unshift(createdTicket);
        selectedTicketId = createdTicket.id;
        renderTicketList();
        renderTicketDetails();
        await openChatForSelectedTicket();

        messageInput.value = "";
    } catch (err) {
        console.error(err);
        alert("Не удалось создать обращение");
    }
}

createTicketBtn.addEventListener("click", createTicket);
refreshBtn.addEventListener("click", loadUserTickets);
sendMessageBtn.addEventListener("click", sendChatMessage);

searchInput.addEventListener("input", () => {
    renderTicketList();
    renderTicketDetails();
});

loadUserTickets();