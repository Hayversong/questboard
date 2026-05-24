document.addEventListener("DOMContentLoaded", () => {
    const cards = document.querySelectorAll(".task-card");

    const columns = document.querySelectorAll(".tasks-container");

    cards.forEach((card) => {
        card.draggable = true;

        card.addEventListener("dragstart", () => {
            card.classList.add("dragging");
        });

        card.addEventListener("dragend", () => {
            card.classList.remove("dragging");
        });
    });

    columns.forEach((column) => {
        column.addEventListener("dragover", (e) => {
            e.preventDefault();

            const dragging = document.querySelector(".dragging");

            if (dragging) {
                column.appendChild(dragging);
            }
        });
    });
});
