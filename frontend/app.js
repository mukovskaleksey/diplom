const tickets = [
    {
        id: 101,
        userId: 1,
        message: "Я забыл пароль и не могу войти в аккаунт.",
        category: "ACCOUNT",
        status: "NEW",
        specialistName: null,
        createdAt: "2026-04-11 14:10",
    },
    {
        id: 102,
        userId: 1,
        message: "С моей карты списали деньги два раза за одну покупку.",
        category: "PAYMENT",
        status: "ASSIGNED",
        specialistName: "Ирина Волкова",
        createdAt: "2026-04-11 14:18",
    },
    {
        id: 103,
        userId: 2,
        message: "Хочу изменить адрес доставки моего заказа.",
        category: "DELIVERY",
        status: "ASSIGNED",
        specialistName: "Алексей Соколов",
        createdAt: "2026-04-11 14:29",
    },
    {
        id: 104,
        userId: 3,
        message: "Где мой возврат денег за отменённый заказ?",
        category: "REFUND",
        status: "CLOSED",
        specialistName: "Мария Иванова",
        createdAt: "2026-04-11 14:42",
    },
];

let selectedTicketId = tickets[0]?.id || null;

const ticketList = document.getElementById("ticketList");
const ticketDetails = document.getElementById("ticketDetails");
const searchInput = document.getElementById("searchInput");
const createTicketBtn = document.getElementById("createTicketBtn");
const userIdInput = document.getElementById("userId");
const messageInput = document.getElementById("message");

const totalCount = document.getElementById("totalCount");
const newCount = document.getElementById("newCount");
const assignedCount = document.getElementById("assignedCount");
const closedCount = document.getElementById("closedCount");

function categoryLabel(value) {
    const map = {
        ACCOUNT: "Аккаунт",
        ORDER: "Заказ",
        REFUND: "Возврат",
        PAYMENT: "Оплата",
        DELIVERY: "Доставка",
        SUPPORT: "Поддержка",
    };
    return map[value] || value;
}

function statusLabel(value) {
    const map = {
        NEW: "Новый",
        ASSIGNED: "Назначен",
        CLOSED: "Закрыт",
    };
    return map[value] || value;
}

function updateStats() {
    totalCount.textContent = tickets.length;
    newCount.textContent = tickets.filter(t => t.status === "NEW").length;
    assignedCount.textContent = tickets.filter(t => t.status === "ASSIGNED").length;
    closedCount.textContent = tickets.filter(t => t.status === "CLOSED").length;
}

function getFilteredTickets() {
    const query = searchInput.value.trim().toLowerCase();
    if (!query) return tickets;

    return tickets.filter(ticket =>
        ticket.message.toLowerCase().includes(query) ||
        ticket.category.toLowerCase().includes(query) ||
        String(ticket.id).includes(query)
    );
}

function renderTicketList() {
    const filtered = getFilteredTickets();

    ticketList.innerHTML = "";

    filtered.forEach(ticket => {
        const div = document.createElement("div");
        div.className = "ticket-item" + (ticket.id === selectedTicketId ? " active" : "");
        div.onclick = () => {
            selectedTicketId = ticket.id;
            renderTicketList();
            renderTicketDetails();
        };

        div.innerHTML = `
      <div class="ticket-item-top">
        <div>
          <div><strong>Тикет #${ticket.id}</strong></div>
          <div class="ticket-meta">Пользователь #${ticket.userId}</div>
        </div>
        <span class="badge ${ticket.category.toLowerCase()}">${categoryLabel(ticket.category)}</span>
      </div>
      <div class="ticket-message">${ticket.message}</div>
      <div class="ticket-footer">
        <span class="badge ${ticket.status.toLowerCase()}">${statusLabel(ticket.status)}</span>
        <span class="ticket-meta">${ticket.createdAt}</span>
      </div>
    `;

        ticketList.appendChild(div);
    });

    if (!filtered.find(t => t.id === selectedTicketId)) {
        selectedTicketId = filtered[0]?.id || null;
    }
}

function renderTicketDetails() {
    const ticket = tickets.find(t => t.id === selectedTicketId);

    if (!ticket) {
        ticketDetails.innerHTML = `<p class="muted">Выберите тикет из списка.</p>`;
        return;
    }

    ticketDetails.innerHTML = `
    <div style="margin-bottom: 16px; display: flex; gap: 8px; flex-wrap: wrap;">
      <span class="badge ${ticket.category.toLowerCase()}">${categoryLabel(ticket.category)}</span>
      <span class="badge ${ticket.status.toLowerCase()}">${statusLabel(ticket.status)}</span>
    </div>

    <div class="detail-grid">
      <div class="detail-card">
        <span class="detail-label">Пользователь</span>
        <span class="detail-value">#${ticket.userId}</span>
      </div>
      <div class="detail-card">
        <span class="detail-label">Специалист</span>
        <span class="detail-value">${ticket.specialistName || "Не назначен"}</span>
      </div>
      <div class="detail-card">
        <span class="detail-label">Создан</span>
        <span class="detail-value">${ticket.createdAt}</span>
      </div>
      <div class="detail-card">
        <span class="detail-label">Статус</span>
        <span class="detail-value">${statusLabel(ticket.status)}</span>
      </div>
    </div>

    <div>
      <p class="detail-label">Сообщение пользователя</p>
      <div class="detail-message">${ticket.message}</div>
    </div>
  `;
}

function createTicket() {
    const userId = Number(userIdInput.value);
    const message = messageInput.value.trim();

    if (!userId || !message) {
        alert("Заполни userId и сообщение");
        return;
    }

    const ticket = {
        id: Date.now(),
        userId,
        message,
        category: "SUPPORT",
        status: "NEW",
        specialistName: null,
        createdAt: new Date().toLocaleString("ru-RU"),
    };

    tickets.unshift(ticket);
    selectedTicketId = ticket.id;
    messageInput.value = "";

    updateStats();
    renderTicketList();
    renderTicketDetails();
}

createTicketBtn.addEventListener("click", createTicket);
searchInput.addEventListener("input", () => {
    renderTicketList();
    renderTicketDetails();
});

updateStats();
renderTicketList();
renderTicketDetails();