const containers = document.querySelectorAll(".tasks-container");

let dragged = null;

document.querySelectorAll(".task-card").forEach((card) => {
    card.addEventListener("dragstart", (e) => {
        dragged = card;
        e.dataTransfer.effectAllowed = "move";
        card.classList.add("dragging");
    });

    card.addEventListener("dragend", async () => {
        card.classList.remove("dragging");
        await saveState();
        dragged = null;
    });
});

containers.forEach((container) => {
    container.addEventListener("dragover", (e) => {
        e.preventDefault();

        if (!dragged) return;

        const after = getCardAfterCursor(container, e.clientY);

        if (!after) {
            container.appendChild(dragged);
        } else {
            container.insertBefore(dragged, after);
        }
    });
});

function getCardAfterCursor(container, y) {
    const cards = [...container.querySelectorAll(".task-card:not(.dragging)")];

    return cards.reduce(
        (closest, card) => {
            const box = card.getBoundingClientRect();
            const offset = y - box.top - box.height / 2;

            if (offset < 0 && offset > closest.offset) {
                return { offset, element: card };
            }

            return closest;
        },
        { offset: Number.NEGATIVE_INFINITY },
    ).element;
}

// Salva ordem E status de cada card conforme a coluna onde está
async function saveState() {
    const projectId = document.querySelector(".container").dataset.projectId;

    const cards = [];

    document.querySelectorAll(".tasks-container").forEach((column) => {
        const status = column.dataset.status;

        column.querySelectorAll(".task-card").forEach((card, index) => {
            cards.push({
                id: card.dataset.id,
                order: index,
                status: status,
            });
        });
    });

    try {
        const response = await fetch("/cards/reorder", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ project_id: projectId, cards }),
        });

        if (!response.ok) {
            console.error("Erro ao salvar:", response.status);
        }
    } catch (err) {
        console.error("Erro ao salvar estado:", err);
    }
}
