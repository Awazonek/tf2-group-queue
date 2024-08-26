document.addEventListener("DOMContentLoaded", () => {
    const joinGroupBtn = document.getElementById("join-group-btn");
    const searchGroupBtn = document.getElementById("search-group-btn");
    const resetUserBtn = document.getElementById("reset-user-btn");
    const groupIdInput = document.getElementById("group-id");
    const statusDiv = document.getElementById("status");

    // Function to ping the server every X seconds
    function startPinging() {
        const pingInterval = 5000; // Ping every 5 seconds

        setInterval(async () => {
            try {
                const groupId = groupIdInput.value.trim();
                const response = await fetch(`http://localhost:8080/user-ping/${groupId}`); 

                if (response.ok) {
                    const data = await response.json();
                    if (data.server) {
                        // If a redirect occurs, handle it (e.g., Steam connection URL)
                        window.location.assign(`steam://connect/${data.server}/quickpick_${data.quickpick}`);
                    }
                    statusDiv.innerText = `Pinged at ${new Date().toLocaleTimeString()}: ${JSON.stringify(data)}`;
                }

            } catch (error) {
                statusDiv.innerText = `Error: ${error.message}`;
            }
        }, pingInterval);
    }

    // Function to handle joining the group
    joinGroupBtn.addEventListener("click", async () => {
        const groupId = groupIdInput.value.trim();
        if (groupId) {
            try {
                const response = await fetch(`http://localhost:8080/join-group/${groupId}`, {
                    method: 'POST',
                });

                if (response.ok) {
                    const data = await response.json();
                    statusDiv.innerText = `Joined group ${JSON.stringify(data)}`;
                } else {
                    statusDiv.innerText = `Failed to join group ${response.status}`;
                }
            } catch (error) {
                statusDiv.innerText = `Failed to join group ${error.message}`;
            }
        } else {
            statusDiv.innerText = `Please enter a group ID`;
        }
    });

    // Function to handle searching the group
    searchGroupBtn.addEventListener("click", async () => {
        const groupId = groupIdInput.value.trim();
        if (groupId) {
            try {
                const response = await fetch(`http://localhost:8080/search-group/${groupId}`, {
                    method: 'POST',
                });

                if (response.ok) {
                    const data = await response.json();
                    statusDiv.innerText = `Search started for group ${data.group}`;
                } else {
                    statusDiv.innerText = `Failed to start search. Status: ${response.status}`;
                }
            } catch (error) {
                statusDiv.innerText = `Error searching group: ${error.message}`;
            }
        } else {
            statusDiv.innerText = `Please enter a Group ID.`;
        }
    });


    // Function to handle searching the group
    resetUserBtn.addEventListener("click", async () => {
        const groupId = groupIdInput.value.trim();
        if (groupId) {
            try {
                const response = await fetch(`http://localhost:8080/reset-user/${groupId}`, {
                    method: 'POST',
                });

                if (response.ok) {
                    const data = await response.json();
                    statusDiv.innerText = `Reset user in group ${data.group}`;
                } else {
                    statusDiv.innerText = `Failed to reset user. Status: ${response.status}`;
                }
            } catch (error) {
                statusDiv.innerText = `Error resetting user: ${error.message}`;
            }
        } else {
            statusDiv.innerText = `Please enter a Group ID.`;
        }
    });

    // Start the background ping when the page loads
    startPinging();
});