function handleLike(entityType, entityId, likeType) {
    fetch('/like', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            type: likeType,
            id: entityId,
            entityType: entityType
        }),
    })
    .then(response => response.json())
    .then(data => {
        if (data.likeCount !== undefined) {
            const likeCountElement = document.querySelector(`#${entityType}-${entityId} .like-count`);
            if (likeCountElement) {
                likeCountElement.textContent = data.likeCount;
            }
        } else {
            console.error('Error:', data.message);
        }
    })
    .catch((error) => {
        console.error('Error:', error);
    });
}
