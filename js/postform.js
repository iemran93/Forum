async function handlePost(event) {
    event.preventDefault();

    const form = event.target;
    const formData = new FormData(form);
    const formDataObject = Object.fromEntries(formData.entries());

    // Convert categories to an array of IDs
    const selectedCategories = Array.from(form.querySelectorAll('input[name="categories"]:checked'))
        .map(checkbox => parseInt(checkbox.value));

    if (selectedCategories.length === 0) {
        const errorElement = document.getElementById('error-message');
        errorElement.textContent = 'Please select at least one category.';
        errorElement.style.display = 'block';
        return;
    }

    formDataObject.categories = selectedCategories;

    const response = await fetch('/postform/submit', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(formDataObject)
    });
    const result = await response.json();

    if (response.ok) {
        alert('Post it successfully');
        // redirect to the post
        window.location.href = '/post?id=' + result.message;
    } else {
        const errorMessage = result.message || 'Signup failed';
        const errorElement = document.getElementById('error-message');
        errorElement.textContent = errorMessage;
        errorElement.style.display = 'block';
    }
}