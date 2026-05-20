const API_BASE = "http://localhost:8078/api/v1";

function getWsBase() {
    return API_BASE.replace(/^http/, "ws");
}

function getChatWsUrl(chatId) {
    return `${getWsBase()}/ws/chats/${chatId}`;
}

function categoryLabel(value) {
    const map = {
        ACCOUNT: "Аккаунт",
        ORDER: "Заказ",
        REFUND: "Возврат",
        PAYMENT: "Оплата",
        DELIVERY: "Доставка",
        SUPPORT: "Поддержка",
    };
    return map[value] || value || "—";
}

function statusLabel(value) {
    const map = {
        NEW: "Новый",
        ASSIGNED: "Назначен",
        CLOSED: "Закрыт",
    };
    return map[value] || value || "—";
}

function categoryClass(value) {
    return (value || "support").toLowerCase();
}

function statusClass(value) {
    return (value || "new").toLowerCase();
}

async function apiGet(path) {
    const response = await fetch(`${API_BASE}${path}`);
    if (!response.ok) {
        const text = await response.text();
        throw new Error(text || "GET request failed");
    }
    return response.json();
}

async function apiPost(path, body) {
    const options = {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
    };

    if (body !== undefined) {
        options.body = JSON.stringify(body);
    }

    const response = await fetch(`${API_BASE}${path}`, options);
    if (!response.ok) {
        const text = await response.text();
        throw new Error(text || "POST request failed");
    }
    return response.json();
}

function normalizeTicket(ticket) {
    const normalizedStatus = (ticket.status || "").toUpperCase();
    const normalizedCategory = (ticket.category || "").toUpperCase();

    return {
        id: ticket.id,
        userId: ticket.user_id,
        message: ticket.message,
        category: normalizedCategory,
        status: normalizedStatus,
        specialistId: ticket.specialist_id || null,
        specialistName: ticket.specialist_name || null,
        createdAt: ticket.created_at,
    };
}

function normalizeChatMessage(message) {
    return {
        id: message.id,
        chatId: message.chat_id,
        senderType: message.sender_type,
        senderId: message.sender_id,
        body: message.body,
        createdAt: message.created_at,
    };
}

function escapeHtml(value) {
    return String(value ?? "")
        .replaceAll("&", "&amp;")
        .replaceAll("<", "&lt;")
        .replaceAll(">", "&gt;")
        .replaceAll('"', "&quot;")
        .replaceAll("'", "&#039;");
}

const userNameCache = new Map();

async function getUserNameById(userId) {
    if (!userId) {
        return null;
    }

    if (userNameCache.has(userId)) {
        return userNameCache.get(userId);
    }

    try {
        const data = await apiGet(`/users/${userId}`);
        const user = data.user || data;
        const name = user?.name || null;

        userNameCache.set(userId, name);
        return name;
    } catch (err) {
        console.error(`Не удалось загрузить пользователя ${userId}`, err);
        userNameCache.set(userId, null);
        return null;
    }
}

async function enrichTicketsWithNames(tickets) {
    const ids = [...new Set(
        tickets.flatMap((ticket) => {
            const result = [];
            if (ticket.userId) result.push(ticket.userId);
            if (ticket.specialistId) result.push(ticket.specialistId);
            return result;
        })
    )];

    await Promise.all(ids.map((id) => getUserNameById(id)));

    tickets.forEach((ticket) => {
        ticket.userName = ticket.userId ? (userNameCache.get(ticket.userId) || null) : null;
        if (!ticket.specialistName && ticket.specialistId) {
            ticket.specialistName = userNameCache.get(ticket.specialistId) || null;
        }
    });

    return tickets;
}

function getStoredUser() {
    const raw = localStorage.getItem("supportDeskUser");
    if (!raw) {
        return null;
    }

    try {
        return JSON.parse(raw);
    } catch (err) {
        console.error("session parse error", err);
        return null;
    }
}

function saveSession(user) {
    localStorage.setItem("supportDeskUser", JSON.stringify(user));
}

function clearSession() {
    localStorage.removeItem("supportDeskUser");
}

function requireAuth() {
    const user = getStoredUser();

    if (!user) {
        window.location.href = "./index.html";
        return null;
    }

    return user;
}

function logout() {
    clearSession();
    window.location.href = "./index.html";
}