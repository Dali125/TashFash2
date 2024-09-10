document.addEventListener("DOMContentLoaded", async function () {



  
    await getFolders();

    const form = document.getElementById("folder");
    form.addEventListener("submit", async function(event) {
        event.preventDefault();
        const formData = new FormData(form);
        const collectionName = formData.get("collection_name");
        const response = await fetch("/create-folder", {
            method: "POST",
            body: formData,
        });
        if (!response.ok) {
            console.error("Failed to create folder:", response.statusText);
        } else {
            console.log("Folder created successfully");
            document.getElementById("overlay").style.display = "none";
            //Refresh the page
            window.location.reload();
        }
    })

    
});






async function getFolders(){

    const myFolders = document.getElementById("myFolders");

    try {

        const response = await fetch('/get-folders')

        if (!response.ok){
            throw new Error("Error Getting Data" + response.status + "")
        }

        const data = await response.json();

        if(data ===  null){
            console.log("No Collections Found")
            myFolders.classList.add("justify-center", "items-center")
            myFolders.innerHTML = `
            
            <div class"flex flex-col gap-4 border-2">
                <p>No Collections Found</p>
                  <button onclick="on()">Add A Collection</button>
            </div>
            `
        }else{

            data.forEach(element => {

                console.log(element.folder_name)
                const folderUi = document.createElement("div");
                folderUi.classList = "flex items-center justify-center flex-col gap-4 h-full w-full md:w-full  border-2 p-4 rounded-md bg-gray-200";
                folderUi.innerHTML = `
                <p class="text-3xl font-bold p-2">${element.folder_name}</p>
                <div class="flex flex-row gap-4">
                <button onclick="openFolder(${element.collection_id}, ${element.folder_name})">Open</button>
                <button class="font-bold h-10 w-32 text-white bg-black rounded-md" onclick="on()">Edit Collection</button>
                </div>
                `
                myFolders.appendChild(folderUi)
                
            });
        }
       
        
    } catch (err) {

        console.error(err)
        
    }

}


function openFolder(id, name) {
    var myId = id;
    var myName = name;
    window.location.href = `/folder?collection_id=${encodeURIComponent(myId)}&name=${encodeURIComponent(name)}`;
}

function off() {
    document.getElementById("overlay").style.display = "none";
}