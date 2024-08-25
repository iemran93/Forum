async function handleSignup(event) {
    event.preventDefault();

    const form = event.target;
    const formData = new FormData(form);
    const formDataObject = Object.fromEntries(formData.entries());

    const response = await fetch('/signup/submit', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(formDataObject)
    });

    const result = await response.json();

    if (response.ok) {
        notifysuccess('Signup successful');
        // Optionally, redirect to another page
        window.location.href = '/';
    } else {
        const errorMessage = result.message || 'Signup failed';
        const errorElement = document.getElementById('error-message');
        errorElement.textContent = errorMessage;
        errorElement.style.display = 'block';
    }
}