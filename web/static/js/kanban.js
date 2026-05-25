document.addEventListener(
    "DOMContentLoaded",

    () => {
        const cards = document.querySelectorAll(".task-card");

        const columns = document.querySelectorAll(".tasks-container");

        let dragged;

        cards.forEach((card) => {
            card.draggable = true;

            card.addEventListener(
                "dragstart",

                () => {
                    dragged = card;

                    card.classList.add("dragging");
                },
            );

            card.addEventListener(
                "dragend",

                () => {
                    card.classList.remove("dragging");
                },
            );
        });

        columns.forEach((column) => {
            column.addEventListener(
                "dragover",

                (e) => {
                    e.preventDefault();
                },
            );

            column.addEventListener(
                "drop",

                async () => {
                    column.appendChild(dragged);

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
