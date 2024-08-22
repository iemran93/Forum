
document.addEventListener('DOMContentLoaded', function () {
    const logoutLink = document.getElementById('logoutLink');
    if (logoutLink) {
        logoutLink.addEventListener('click', async function (event) {
            event.preventDefault();
            try {
                const response = await fetch('/logout', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    }
                });
                const result = await response.json();
                if (response.ok) {
                    alert(result.message);
                    window.location.href = '/';
                } else {
                    alert('Logout failed: ' + result.message);
                }
            } catch (error) {
                console.error('Error:', error);
            }
        });
    }

    document.getElementById('filterLink').addEventListener('click', function (e) {
        e.preventDefault(); // Prevent the default link behavior

        var filterContainer = document.getElementById('filterContainer');
        if (filterContainer.style.display === 'none') {
            filterContainer.style.display = 'block';
        } else {
            filterContainer.style.display = 'none';
        }
    });
});

async function handleLike(entityType, entityId, likeType) {
    const response = await fetch('/like', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            type: likeType,
            id: parseInt(entityId),
            entityType: entityType
        })
    })
    const result = await response.json();
    const [likes, dislikes] = result.message.split(',').map(Number);
    if (response.ok) {
        const likeCount = document.querySelector(`#${entityType}-${entityId} .like-count`);
        likeCount.innerText = `${likes}`
        const dislikeCount = document.querySelector(`#${entityType}-${entityId} .dislike-count`);
        dislikeCount.innerText = `${dislikes}`
    } else {
        const errorMessage = result.message || 'failed';
        alert(errorMessage)
    }
}

// clear button for the filter form
function clearForm() {
    // reset the checkboxes and radio buttons
    document.getElementById("filterForm").reset();
}
