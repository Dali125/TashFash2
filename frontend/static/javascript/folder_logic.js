document.addEventListener("DOMContentLoaded", async function () {
    const urlParams = new URLSearchParams(window.location.search);
    const collection_id = urlParams.get('collection_id');
    const name = urlParams.get('name');
    console.log("Collection ID:", collection_id);
    console.log("Name:", name);

    await getPhotos(collection_id);
    
});

async function getPhotos(collection_id) {
    const folderContents = document.getElementById('folderContents');
    const formData = new FormData();

    formData.append("collection_id", collection_id);

    try {
        const response = await fetch("/get_photos", {
            method: "POST",
            body: formData,
            cache: "no-cache",
        });
        
        if (!response.ok) {
            console.error("Failed to fetch photos:", response.status, response.statusText);
            return;
        }

        const data = await response.json();
        console.log("Fetched photos:", data);

        const myPhotos = data.photos;
        myPhotos.forEach(element => {
            console.log("Photo link:", element.image_link);
            
            // Create container for image and buttons
            const imageContainer = document.createElement('div');
            imageContainer.classList = 'image-container overflow-hidden rounded-lg shadow-lg relative group';

            // Create image element
            const imageElement = document.createElement('img');
            imageElement.src = element.image_link;
            imageElement.classList = 'w-full h-full object-cover';

            // Create buttons container
            const buttonsContainer = document.createElement('div');
            buttonsContainer.classList = 'image-buttons justify-center space-x-2 absolute inset-x-0 w-full h-12 items-center bottom-0 pl-16';

            // Create view button
            const viewButton = document.createElement('button');
            viewButton.innerText = 'View';
            viewButton.classList = 'bg-green-500 text-white px-3 py-1 rounded-md hover:bg-green-600';
            viewButton.onclick = function() {
                window.open(element.image_link, '_blank');
            };

            // Create delete button
            const deleteButton = document.createElement('button');
            deleteButton.innerText = 'Delete';
            deleteButton.classList = 'bg-red-500 text-white px-3 py-1 rounded-md hover:bg-red-600';
            deleteButton.onclick = async function() {
                console.log("Attempting to delete image with link:", element.image_link, "and ID:", element.collection_id);
                await deleteImage(element.image_link, element.id);
            };

            // Append buttons to container
            buttonsContainer.append(viewButton, deleteButton);
            // Append image and buttons to image container
            imageContainer.append(imageElement, buttonsContainer);
            // Append image container to the gallery
            folderContents.append(imageContainer);
        });
    } catch (error) {
        console.error("Error fetching photos:", error);
    }
}

async function deleteImage(image_link, image_id) {
    const formData = new FormData();
    formData.append("image_link", image_link);
    formData.append("image_id", image_id);
    console.log(image_id)

    try {
        const response = await fetch("/delete-photo", {
            method: "POST",
            body: formData,
        });

        if (!response.ok) {
            console.error("Failed to delete image:", response.status, response.statusText);
            const errorData = await response.json();
            console.error("Error details:", errorData);
        } else {
            console.log("Image deleted successfully");
            // Refresh the page
            window.location.reload();
        }
    } catch (error) {
        console.error("Error deleting image:", error);
    }
}

async function addImage(params) {
    
}
