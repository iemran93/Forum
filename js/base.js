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
});