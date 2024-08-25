document.addEventListener('DOMContentLoaded', function () {
    async function handleAddComment(event) {
        event.preventDefault();

        const form = event.target;
        const formData = new FormData(form);
        const formDataObject = Object.fromEntries(formData.entries());

        const response = await fetch('/comment', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(formDataObject)
        });

        const result = await response.json();

        if (response.ok) {
            alert('Comment added successfully');
            // redirect to the post
            window.location.href = '/post?id=' + formDataObject.post_id;
        } else {
            const errorMessage = result.message || 'failed';
            const errorElement = document.getElementById('error-message');
            errorElement.textContent = errorMessage;
            errorElement.style.display = 'block';
        }
    }

    const commentForm = document.getElementById('addcomment');
    if (commentForm) {
        commentForm.addEventListener('submit', handleAddComment);
    }
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
        notify(errorMessage);
    }
}