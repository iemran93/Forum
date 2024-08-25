document.addEventListener('DOMContentLoaded', function () {
    const logoutLink = document.getElementById('logoutLink');
    if (logoutLink) {
        logoutLink.addEventListener('click', async function (event) {
            event.preventDefault();
            const result = await notifyConfirm('Are you sure you want to logout?');
            if (result.isConfirmed) {
                try {
                    const response = await fetch('/logout', {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json'
                        }
                    });
                    const result = await response.json();
                    if (response.ok) {
                        window.location.href = '/';
                    } else {
                        notify('Logout failed: ' + result.message);
                    }
                } catch (error) {
                    console.error('Error:', error);
                }
            }
        });
    }
});

function notify(mesg) {
    Swal.fire(mesg);
}

function notifySuccess(mesg) {
    Swal.fire({
        position: "top-end",
        icon: "success",
        title: mesg,
        showConfirmButton: false,
        timer: 1500
      });
}

function notifyConfirm(mesg) {
    return Swal.fire({
        title: mesg,
        showCancelButton: true,
        confirmButtonText: 'Yes',
        cancelButtonText: 'No'
    });
}