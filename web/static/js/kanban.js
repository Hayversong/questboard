document.addEventListener(
    "DOMContentLoaded",

    () => {
        const cards = document.querySelectorAll(".task-card");

        const columns = document.querySelectorAll(".tasks-container");

        let dragged = null;

        cards.forEach((card) => {
            card.setAttribute("draggable", "true");

            card.addEventListener(
                "dragstart",

                (e) => {
                    dragged = card;

                    setTimeout(
                        () => {
                            card.classList.add("dragging");
                        },

                        0,
                    );
                },
            );

            card.addEventListener(
                "dragend",

                () => {
                    card.classList.remove("dragging");

                    columns.forEach((col) => {
                        col.classList.remove("drag-hover");
                    });
                },
            );
        });

        columns.forEach((column) => {
            column.addEventListener(
                "dragover",

                (e) => {
                    e.preventDefault();

                    column.classList.add("drag-hover");
                },
            );

            column.addEventListener(
                "dragleave",

                () => {
                    column.classList.remove("drag-hover");
                },
            );

            column.addEventListener(
                "drop",

                async (e) => {
                    e.preventDefault();

                    if (!dragged) return;

                    column.classList.remove("drag-hover");

                    column.appendChild(dragged);

                    dragged.classList.add("drop-success");

                    setTimeout(
                        () => {
                            dragged.classList.remove("drop-success");
                        },

                        350,
                    );

                    const body = new URLSearchParams();

                    body.append(
                        "project_id",

                        document.getElementById("project-id").value,
                    );

                    body.append(
                        "card_id",

                        dragged.dataset.id,
                    );

                    body.append(
                        "status",

                        column.dataset.status,
                    );

                    await fetch(
                        "/cards/status",

                        {
                            method: "POST",

                            body,
                        },
                    );
                },
            );
        });
    },
);
