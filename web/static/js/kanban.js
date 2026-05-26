const containers = document.querySelectorAll(".tasks-container");

let dragged = null;

document.querySelectorAll(".task-card").forEach((card) => {
    card.draggable = true;

    card.addEventListener("dragstart", (e) => {
        dragged = card;

        e.dataTransfer.effectAllowed = "move";

        card.classList.add("dragging");
    });

    card.addEventListener("dragend", async (e) => {
        console.log("DRAG END");

        card.classList.remove("dragging");

        await saveOrder();

        dragged = null;
    });
});

containers.forEach((container) => {
    container.addEventListener("dragover", (e) => {
        e.preventDefault();

        const after = getCardAfterCursor(container, e.clientY);

        if (!dragged) return;

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
                return {
                    offset,
                    element: card,
                };
            }

            return closest;
        },
        {
            offset: Number.NEGATIVE_INFINITY,
        },
    ).element;
}

async function saveOrder() {
    console.log("SALVANDO ORDEM");

    const projectId = document.querySelector(".container").dataset.projectId;

    const order = [];

    document.querySelectorAll(".tasks-container").forEach((column) => {
        column.querySelectorAll(".task-card").forEach((card, index) => {
            order.push({
                id: card.dataset.id,

                order: index,
            });
        });
    });

    console.log(
        JSON.stringify(
            {
                project_id: projectId,

                cards: order,
            },
            null,
            2,
        ),
    );

    try {
        await fetch("/cards/reorder", {
            method: "POST",

            headers: {
                "Content-Type": "application/json",
            },

            body: JSON.stringify({
                project_id: projectId,

                cards: order,
            }),
        });
    } catch (err) {
        console.error(err);
    }
}
