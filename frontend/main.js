fetch('http://localhost:8080/items')
    .then(response => response.json())
    .then(data => {
        console.log(data);
        generateTableRows(data);
    })
    .catch(error => console.error('Error fetching data:', error));

function generateTableRows(data) {
    const tbody = document.querySelector("#myTable tbody");

    data.forEach(item => {
        const row = createRow(item);
        tbody.appendChild(row);
    });
}

function createRow(item) {
    const row = document.createElement("tr");
    row.innerHTML = `
        <td>${item.statement_description}</td>
        <td>${item.posting_description}</td>
        <td>${item.income.received_total}</td>
        <td>${item.income.services_performed}</td>
        <td>${item.comment}</td>
        <td>
            <button onclick="editItem(${item.id})">Edit</button>
            <button onclick="deleteItem(${item.id})">Delete</button>
        </td>
    `;
    return row;
}

function addItem() {
}

function editItem(id) {
}

function deleteItem(id) {
    fetch('http://localhost:8080/item', {
        method: 'DELETE',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ id })
    })
        .then(response => {
            if (response.ok) {
                window.alert('Deletion succeeded')
            } else {
                window.alert('Deletion failed')
            }
        })
}
