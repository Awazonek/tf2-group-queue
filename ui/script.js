let selectedGroupId = null;
let selectedGroupIP = null;
let selectedGroupPort = null;
const groupMap = new Map();
document.addEventListener("DOMContentLoaded", () => {
    const joinGroupBtn = document.getElementById("join-group-btn");
    const searchGroupBtn = document.getElementById("search-group-btn");
    const resetUserBtn = document.getElementById("reset-user-btn");
    const groupIdInput = document.getElementById("group-id");
    const statusDiv = document.getElementById("status");
    // Start the background ping when the page loads
    startPinging();
    listGroups();


   async function listGroups() {
        // Call the list-groups endpoint
        try{
                const response = await fetch('/list-groups');
                const data = await response.json();
                
                // Access the 'groups' field in the JSON response
                const groups = data.groups;
                
                populateGroupsContainer(groups);
            }
            catch(error) {
                console.error('Error fetching groups:', error);
            }
    }

    function populateGroupsContainer(groups) {
        const container = document.getElementById('groups-container');
        container.innerHTML = ''; // Clear any existing content

        groups.forEach(group => {
            const groupDiv = document.createElement('div');
            groupDiv.className = 'group-item';

            const groupName = document.createElement('h3');
            groupName.textContent = `Name: ${group.id}`;
            groupDiv.appendChild(groupName);

            const groupMap = document.createElement('p');
            groupMap.textContent = `Current Map: ${group.server_info.map}`;
            groupDiv.appendChild(groupMap);

            const groupIP = document.createElement('p');
            groupIP.textContent = `IP: ${group.server_info.ip}`;
            groupDiv.appendChild(groupIP);

            const joinButton = document.createElement('button');
            joinButton.className = 'join-btn';
            joinButton.textContent = 'Join Group';
            joinButton.onclick = () => joinGroup(group.id);
            groupDiv.appendChild(joinButton);

            container.appendChild(groupDiv);
        });
    }

        // Function to join a group
        async function joinGroup(groupID) {
            try {
                const response = await fetch(`/join-group/${groupID}`, {
                    method: 'POST'
                });
                if (response.ok) {
                    const resp = await response.json();
                    displaySelectedGroup(resp.group)
                    
                    if (resp.group.server_info.ip && resp.group.server_info.port) {
                        const quickPick = incrementGroupValue(resp.group.server_info.ip, resp.group.server_info.port)

                        window.location.assign(`steam://connect/${resp.group.server_info.ip}:${resp.group.server_info.port}/quickpick_${quickPick}`);
                        selectedGroupIP = resp.group.server_info.ip
                        selectedGroupPort = resp.group.server_info.port
                    }
                } else {
                    alert('Failed to join the group.');
                }
            } catch (error) {
                console.error('Error joining group:', error);
            }
        }


    // Function to display the selected group
    function displaySelectedGroup(group) {
        const searchGroups = document.getElementById('group-list');
        searchGroups.innerHTML = ''; // Clear any existing content
        const selectedGroupContainer = document.getElementById('selected-group-container');
        const selectedGroupSection = document.getElementById('selected-group');

        selectedGroupContainer.innerHTML = ''; // Clear any existing content

        const groupName = document.createElement('h3');
        groupName.textContent = `Name: ${group.id}`;
        selectedGroupContainer.appendChild(groupName);

        const groupMap = document.createElement('p');
        groupMap.textContent = `Current Map: ${group.server_info.map}`;
        selectedGroupContainer.appendChild(groupMap);
        const groupPlayers = document.createElement('p');
        groupPlayers.textContent = `Current Players: ${group.server_info.player_count}`;
        selectedGroupContainer.appendChild(groupPlayers);

        const groupIP = document.createElement('p');
        groupIP.textContent = `IP: ${group.server_info.ip}`;
        selectedGroupContainer.appendChild(groupIP);

        const mapsList = document.createElement('div');
        mapsList.id = 'maps-list';

        const mapsString = group.server_parameters.maps.join(', ');
        const mapsText = document.createElement('p');
        mapsText.textContent = `Allowed Maps: ${mapsString}`;
        selectedGroupContainer.appendChild(mapsText);

        selectedGroupContainer.appendChild(mapsList);

        selectedGroupId = group.id
        selectedGroupIP = group.server_info.ip
        selectedGroupPort = group.server_info.port
        selectedGroupSection.style.display = 'block';
        if (group.searching) {
            const searching = document.createElement('p');
            searching.textContent = `Searching!`;
            selectedGroupContainer.appendChild(searching);
        } else {
            const searchButton = document.createElement('button');
            searchButton.className = 'search-btn';
            searchButton.textContent = 'Queue';
            searchButton.onclick = () => searchGroup();
            selectedGroupContainer.appendChild(searchButton);
            searchButton.style.display = 'block';
        }
    }

    async function searchGroup() {
        try {
            // todo validate we have one
            const response = await fetch(`/search-group/${selectedGroupId}`, {
                method: 'POST',
            });

            if (response.ok) {
                const data = await response.json();
                displaySelectedGroup(data.group)
            } else {
                statusDiv.innerText = `Failed to start search. Status: ${response.status}`;
            }
        } catch (error) {
            statusDiv.innerText = `Error searching group: ${error.message}`;
        }
    }

    // Function to ping the server every X seconds
    function startPinging() {
        const pingInterval = 2000; // Ping every 5 seconds

        setInterval(async () => {
            try {
                if (!selectedGroupId) {
                    return
                }
                const response = await fetch(`/user-ping/${selectedGroupId}`); 

                if (response.ok) {
                    const data = await response.json();
                    if (data.group.server_info.ip != selectedGroupIP || data.group.server_info.port != selectedGroupPort) {
                        // If a redirect occurs, handle it (e.g., Steam connection URL)
                        const quickPick = incrementGroupValue(selectedGroupIP, selectedGroupPort)

                        window.location.assign(`steam://connect/${data.group.server_info.ip}:${data.group.server_info.port}/quickpick_${quickPick}`);
                        selectedGroupIP = data.group.server_info.ip
                        selectedGroupPort = data.group.server_info.port
                    }
                    displaySelectedGroup(data.group)
                    
                }

            } catch (error) {
                statusDiv.innerText = `Error: ${error.message}`;
            }
        }, pingInterval);
    }

    // Function to add a group to the map
    function incrementGroupValue(ip, port) {
        const key = `${ip}:${port}`;
        const currentValue = groupMap.get(key) || 2; // start at 2 sdon't ask
        const newValue = currentValue + 1;
        groupMap.set(key, newValue);
        return newValue
    }


});