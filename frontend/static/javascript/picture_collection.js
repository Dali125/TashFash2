document.addEventListener("DOMContentLoaded", async function () {
    await getFolders();


    await constructDropdown();
    const form = document.getElementById("folder");
    form.addEventListener("submit", async function (event) {
        event.preventDefault();
        const formData = new FormData(form);
        const collectionName = formData.get("collection_name");
        //Get selection value from dropdown
        const dropdown = document.getElementById("dropdown");
        const selectedValue = dropdown.options[dropdown.selectedIndex].value;
        formData.append("id", selectedValue);
        formData.append("collection_name", collectionName);

        // Debugging FormData to ensure correct data
        console.log("Form Data:", Array.from(formData.entries()));

        const response = await fetch("/create-folder", {
            method: "POST",
            body: formData,
        });

        if (!response.ok) {
            console.error("Failed to create folder:", response.statusText);
        } else {
            console.log("Folder created successfully");
            document.getElementById("overlay").style.display = "none";
            // Refresh the page to show the new folder
            window.location.reload();
        }
    });
});

async function getFolders() {
    const myFolders = document.getElementById("myFolders");

    try {
        const response = await fetch('/get-folders');

        if (!response.ok) {
            throw new Error("Error Getting Data: " + response.status);
        }

        const data = await response.json();

        if (data === null || data.length === 0) {
            console.log("No Collections Found");
            myFolders.classList.add("justify-center", "items-center", "flex-col");
            myFolders.innerHTML = `
            <div class="flex flex-col gap-4 border-2">
                <p>No Collections Found</p>
                <button class="font-bold h-14 w-full md:w-32 text-white bg-black rounded-md p-2 text-sm md:text-base" onclick="on()">Add A Collection</button>
            </div>
            `;
        } else {
            // Clear container only once before the loop
            myFolders.innerHTML = "";
            data.forEach(element => {
                console.log(element.folder_name);

                const folderUi = document.createElement("div");
                folderUi.classList = "flex items-center justify-center flex-col gap-4 h-full w-full md:w-full border-2 p-4 rounded-md bg-gray-200";

                const folderName = document.createElement("p");
                folderName.classList = "text-2xl font-bold p-2 md:text-3xl"; // Responsive text size

                const button_container = document.createElement("div");
                button_container.classList = "flex flex-col gap-4 w-full md:flex-row md:gap-4"; // Stack buttons on mobile, row on desktop

                const open_button = document.createElement("button");
                open_button.classList = "font-bold h-12 w-full md:w-32 text-white bg-black rounded-md p-2 text-sm md:text-base"; // Full width on mobile

                open_button.innerText = "Open";
                open_button.type = "button";
                open_button.addEventListener("click", async () => {
                    await openFolder(element.collection_id, element.folder_name);
                });

                const edit_button = document.createElement("button");
                edit_button.classList = "font-bold h-12 w-full md:w-32 text-white bg-black rounded-md p-2 text-sm md:text-base"; // Full width on mobile

                edit_button.type = "button";
                edit_button.innerText = "Edit";
                edit_button.addEventListener("click", () => {
                    on();
                });

                button_container.appendChild(open_button);
                button_container.appendChild(edit_button);

                folderUi.appendChild(button_container);
                folderName.innerText = element.folder_name;
                folderUi.appendChild(folderName);

                myFolders.classList.remove("justify-center", "items-center");
                // Append without clearing
                myFolders.appendChild(folderUi);
            });
        }
    } catch (err) {
        console.error("Error fetching folders:", err);
        myFolders.innerHTML = `<p>Error fetching collections. Please try again later.</p>`;
    }
}

async function constructDropdown() {
    const dropdown = document.getElementById("dropdown");

    try {

        const response = await fetch("/get_dropdownValues", {
            method: 'POST'
        });

        if (!response.ok) {
            console.log("Error Detected" + response.statusText)
        }
        const data = await response.json();


        dropdown.innerHTML = '';
        if (data === null || data.length === 0) {
            console.log("No Categories Found");
            const placeholder = document.createElement("option");
            placeholder.text = "No Collections Available";
            placeholder.disabled = true;
            dropdown.appendChild(placeholder);
        } else {
            console.log(data.collections);

            data.collections.forEach(element => {
                const option = document.createElement("option");
                option.value = element.id;
                option.text = element.collection_name;
                dropdown.appendChild(option);
            });

        }
    } catch (error) {
        console.error(error);

    }

}



async function openFolder(id, name) {
    console.log("Opening Folder:", id, name);
    window.location.href = `/folder?collection_id=${encodeURIComponent(id)}&name=${encodeURIComponent(name)}`;
}

function off() {
    document.getElementById("overlay").style.display = "none";
}

function on() {
    document.getElementById("overlay").style.display = "flex";
}
