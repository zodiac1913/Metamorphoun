//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
//!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! J.J. !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
//^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

//----------------------------------------
//dynamite.js handles DOM manipulation.
import {jsonToHtml} from './cc/ccUtilities.js'
import traffic from './comms.js'
export default class dynamite{
    constructor(cfg){
        let bang=this;
        bang.textLibs=document.querySelector("#textLibraries");
        bang.imagesToolsDiv=document.querySelector("#ImagesSection");
        bang.addLibraryButton=document.querySelector("#AddLibraryButton");
        bang.openFolderButton= document.querySelector("#OpenFolderButton");
        bang.editLibraryButton= document.querySelector("#EditLibraryButton");
        bang.removeLibraryButton= document.querySelector("#RemoveLibraryButton");
        bang.useLibraryButton= document.querySelector("#UseLibraryButton");
        bang.selectedImageLibrary=undefined;
        bang.config=null;
        bang.traffic=new traffic();
    }

    wireConfigProperties(config){
        let bang=this;
        bang.config=config;
        let inputz=Array.from(document.querySelectorAll(".primaryInput"));
        let quoteFontRandom=document.querySelector("#quoteFontRandom");
        //onload
        if(quoteFontRandom.checked){
            document.querySelector("#textFontFileEnvelope").classList.add("d-none");
            document.querySelector("#textFontFile").selectedIndex=0;
        }else{
            document.querySelector("#textFontFileEnvelope").classList.remove("d-none");
        }
        let quoteAppearanceRandom=document.querySelector("#quoteAppearanceRandom");
        if(quoteAppearanceRandom.checked){
            document.querySelector("#quoteAppearanceTextColorEnvelope").classList.add("d-none");
            document.querySelector("#quoteAppearanceBackgroundColorEnvelope").classList.add("d-none");
            document.querySelector("#quoteAppearanceOpacityEnvelope").classList.add("d-none");
        }else{
            document.querySelector("#quoteAppearanceTextColorEnvelope").classList.remove("d-none");
            document.querySelector("#quoteAppearanceBackgroundColorEnvelope").classList.remove("d-none");
            document.querySelector("#quoteAppearanceOpacityEnvelope").classList.remove("d-none");
        }
        //wired
		for(let inpt of inputz){
			//set current
			inpt.value=config[inpt.id];
			//set change functions
			inpt.addEventListener("change",async (e)=>{
				await window.comms.formFieldChanged(e);
				if(e.target.id === "quoteFontRandom") bang.quoteFontRandomWiring(e);
                if(e.target.id === "quoteAppearanceRandom") bang.quoteAppearanceRandomWiring(e);

			});
		}
    }






    wireImageLibraryTools(){
        let bang=this;
        let imagesDiv=document.querySelector("#ImagesDiv");
        let imageDivRows = imagesDiv.querySelectorAll(".libraryRow");

        for (let imageRow of imageDivRows) {
          imageRow.addEventListener("mouseover", (e) => {
            // Remove highlight from all rows
            imageRow.parentElement.querySelectorAll(".libraryRow").forEach(row => row.classList.remove("border", "border-CornflowerBlue", "border-1"));
            // Highlight the current row
            e.currentTarget.classList.add("border", "border-CornflowerBlue", "border-1");
          });
          imageRow.addEventListener("click",(e)=>{
            imageRow.parentElement.querySelectorAll(".libraryRow").forEach(row => 
                row.classList.remove("bg-CornflowerBlue","imageSelected")
            );
            e.currentTarget.classList.add("bg-CornflowerBlue","imageSelected")
            if(bang.selectedImageLibrary===e.currentTarget){
                bang.selectedImageLibrary=undefined;

                e.currentTarget.classList.remove("bg-CornflowerBlue","imageSelected")
            }else{
                bang.selectedImageLibrary=e.currentTarget;
                bang.wireAndConfigureImageLibraryToolButtons();
            }
          })
        }
        bang.wireAndConfigureImageLibraryToolButtons();
    }

    wireAndConfigureImageLibraryToolButtons(){
        let bang=this;
        bang.addLibraryButton=document.querySelector("#AddLibraryButton");
        bang.addLibraryButton.addEventListener("click", (e) => {bang.popupAddLibrary();})

        if(bang.selectedImageLibrary!==undefined) {
            bang.addLibraryButton.classList.add("d-none");
            bang.editLibraryButton.classList.remove("d-none");
            bang.removeLibraryButton.classList.remove("d-none");
            bang.useLibraryButton.classList.remove("d-none");
            bang.openFolderButton.classList.remove("d-none");
        }else{
            bang.addLibraryButton.classList.remove("d-none");
            bang.editLibraryButton.classList.add("d-none");
            bang.removeLibraryButton.classList.add("d-none");
            bang.useLibraryButton.classList.add("d-none");
            bang.openFolderButton.classList.add("d-none");
            bang.openFolderButton.addEventListener("click",async (e)=>{
                let loc=bang.selectedImageLibrary.querySelector(".locationurl").dataset.url;
                let dataUp={"id": bang.selectedImageLibrary.id,"loc":loc};
                let apicallRtn=await bang.traffic.apiCall(bang.traffic.server + "/openLocationApi",dataUp);
            })
            // bang.editLibraryButton.addEventListener("click",async (e)=>{
            //     let loc=bang.selectedImageLibrary.querySelector(".locationurl").dataset.url;
            //     let dataUp={"id": bang.selectedImageLibrary.id,"loc":loc};
            //     let apicallRtn=await bang.traffic.apiCall(bang.traffic.server + "/openLocationApi",dataUp);
            // })
            
        }


    }

    async popupAddLibrary() {
        let bang = this;
        let appDiv = document.querySelector("#app");
        appDiv.classList.add("d-none");

    
        // Check if the dialog already exists
        let existingDialog = document.querySelector("#AddLibraryModal");
        if (existingDialog) {
            existingDialog.showModal();
            return;
        }
    

        let dialog = {n:"dialog",i:"AddLibraryModal",c:"w-50",b:[]}
        let card = {i:"AddLibraryModalCard",c:"card bg-Bisque",b:[]}
        //let cardHeader = 

        // let dialog = document.createElement("dialog");
        // dialog.id = "AddLibraryModal";
        // dialog.className = "w-50";
    
        // let card = document.createElement("div");
        // card.id = "AddLibraryModalCard";
        // card.className = "card bg-Bisque";
        // dialog.appendChild(card);
    
        let cardHeader = document.createElement("div");
        cardHeader.className = "card-header";
        cardHeader.innerHTML = "Add Image Library";
        card.appendChild(cardHeader);
    
        let cardBody = document.createElement("div");
        cardBody.id="CardBodyAddLibrary"
        cardBody.className = "card-body d-flex flex-column align-items-center";
    
        //FORM FIELDS
        let useCheckbox = document.createElement("input");
        useCheckbox.type = "checkbox";
        useCheckbox.id = "UseLibraryCheckbox";
        useCheckbox.className = "d-none";
        cardBody.appendChild(useCheckbox);
            
        //NAME
        let nameInputRow = document.createElement("div");
        nameInputRow.className = "flex-row d-flex my-2 w-100";

        let labelName = document.createElement("label");
        labelName.htmlFor = "LibraryNameInput";
        labelName.innerHTML = "Name:";
        nameInputRow.appendChild(labelName);

        let nameInput = document.createElement("input");
        nameInput.type = "text";
        nameInput.id = "LibraryNameInput";
        nameInput.className = "form-control ms-2";
        nameInput.placeholder = "Enter library name";
        nameInput.required = true;
        nameInputRow.appendChild(nameInput);
        cardBody.appendChild(nameInputRow);
    
        //TITLE
        let titleInputRow = document.createElement("div");
        titleInputRow.className = "flex-row d-flex my-2 w-100";

        let labelTitle = document.createElement("label");
        labelTitle.htmlFor = "LibraryTitleInput";
        labelTitle.innerHTML = "Title:";
        titleInputRow.appendChild(labelTitle);

        let titleInput = document.createElement("input");
        titleInput.type = "text";
        titleInput.id = "LibraryTitleInput";
        titleInput.className = "form-control ms-2";
        titleInput.placeholder = "Enter library title";
        titleInput.required = true;
        titleInputRow.appendChild(titleInput);
        cardBody.appendChild(titleInputRow);

        //LOCATION
        let locationInputRow = document.createElement("div");
        locationInputRow.className = "flex-row d-flex my-2 w-100";

        let labelLocation = document.createElement("label");
        labelLocation.htmlFor = "LibraryLocationInput";
        labelLocation.innerHTML = "Location:";
        locationInputRow.appendChild(labelLocation);

        let locationInput = document.createElement("input");
        locationInput.type = "text";//"file";
        locationInput.id = "LibraryLocationInput";
        locationInput.className = "form-control ms-2";
        locationInput.placeholder = "Enter library location";
        //locationInput.webkitdirectory=true;
        //locationInput.directory=true;
        locationInput.required = true;
        locationInputRow.appendChild(locationInput);
        let locationPopUp = document.createElement("a");
        locationPopUp.href = "http://127.0.0.1:3000/openLocationApi";
        locationPopUp.target = "_blank";
        let imgIcon=document.createElement("img");
        imgIcon.src="/pics/folder.png";
        imgIcon.alt="Use this to open a folder and copy the folder path to paste in the input.  Sorry for the overzealous web security issues!"
        imgIcon.height=32;
        imgIcon.width-32;
        locationPopUp.appendChild(imgIcon);
        locationInputRow.appendChild(locationPopUp);
        cardBody.appendChild(locationInputRow);



        card.appendChild(cardBody);
        dialog.appendChild(card);
        let mdl = document.querySelector("#modal");
        mdl.appendChild(dialog);
        let cardBodyDiv=document.querySelector("#CardBodyAddLibrary");
        //JML
        let opInputRow={c:"flex-row d-flex my-2 w-100",b:[]};
        let opLabel={n:"label",for:"LibraryOperationInput",t:"Operation:"}
        let opInput={n:"input",i:"LibraryOperationInput",type:"text",value:"Folder", 
            required:true,c: "form-control ms-2",readonly:true };
        opInputRow.b.push(opLabel);
        opInputRow.b.push(opInput);
        let opInputRowHtml=jsonToHtml(opInputRow);
        cardBodyDiv.insertAdjacentHTML('beforeend',opInputRowHtml);
        
        let inherentInputRow={c:"flex-row d-flex my-2 w-100",b:[]};
        let inherentInput={n:"input",i:"LibraryInherentInput",type:"hidden",value:false };
        inherentInputRow.b.push(inherentInput);
        let inherentInputRowHtml=jsonToHtml(inherentInputRow);
        cardBodyDiv.insertAdjacentHTML('beforeend',inherentInputRowHtml);

        let buttonsDivRow={c:"d-flex justify-content-between my-2 w-100",b:[]};
        let closeButton={n:"button",type:"button",
            i:"AddImageLibraryCloseButton",
            title:"Close form do NOT save data",t:"Close",
            c:"btn btn-secondary text-warning mx-2"};
        buttonsDivRow.b.push(closeButton);
        let saveButton={n:"button",type:"button",
            i:"AddImageLibrarySaveButton",
            title:"Save Library",t:"Save",
            c:"btn btn-success text-warning mx-2"};
        buttonsDivRow.b.push(saveButton);
        let buttonsDivRowHtml=jsonToHtml(buttonsDivRow);
        cardBodyDiv.insertAdjacentHTML('beforeend',buttonsDivRowHtml);
        cardBodyDiv.insertAdjacentHTML('beforeend',inherentInputRowHtml);
        mdl.classList.remove("d-none");
        let addImageLibraryCloseButton=document.querySelector("#AddImageLibraryCloseButton");
        addImageLibraryCloseButton.addEventListener('click',async()=>{bang.closeAddLibraryForm();})
        let addImageLibrarySaveButton=document.querySelector("#AddImageLibrarySaveButton");
        addImageLibrarySaveButton.addEventListener('click',async()=>{bang.saveAddLibraryForm();})
        dialog.showModal();
    }

    async popupEditLibrary() {
        let bang = this;
        let appDiv = document.querySelector("#app");
        appDiv.classList.add("d-none");

    
        // Check if the dialog already exists
        let existingDialog = document.querySelector("#EditLibraryModal");
        if (existingDialog) {
            existingDialog.showModal();
            return;
        }
    

        let dialog = {n:"dialog",i:"EditLibraryModal",c:"w-50",b:[]}
        let card = {i:"EditLibraryModalCard",c:"card bg-Bisque",b:[]}
        //let cardHeader = 

        // let dialog = document.createElement("dialog");
        // dialog.id = "AddLibraryModal";
        // dialog.className = "w-50";
    
        // let card = document.createElement("div");
        // card.id = "AddLibraryModalCard";
        // card.className = "card bg-Bisque";
        // dialog.appendChild(card);
    
        let cardHeader = document.createElement("div");
        cardHeader.className = "card-header";
        cardHeader.innerHTML = "Edit Image Library";
        card.appendChild(cardHeader);
    
        let cardBody = document.createElement("div");
        cardBody.id="CardBodyEditLibrary"
        cardBody.className = "card-body d-flex flex-column align-items-center";
    
        //FORM FIELDS
        let useCheckbox = document.createElement("input");
        useCheckbox.type = "checkbox";
        useCheckbox.id = "UseLibraryCheckbox";
        useCheckbox.className = "d-none";
        cardBody.appendChild(useCheckbox);
            
        //NAME
        let nameInputRow = document.createElement("div");
        nameInputRow.className = "flex-row d-flex my-2 w-100";

        let labelName = document.createElement("label");
        labelName.htmlFor = "LibraryNameInput";
        labelName.innerHTML = "Name:";
        nameInputRow.appendChild(labelName);

        let nameInput = document.createElement("input");
        nameInput.type = "text";
        nameInput.id = "LibraryNameInput";
        nameInput.className = "form-control ms-2";
        nameInput.placeholder = "Enter library name";
        nameInput.required = true;
        nameInputRow.appendChild(nameInput);
        cardBody.appendChild(nameInputRow);
    
        //TITLE
        let titleInputRow = document.createElement("div");
        titleInputRow.className = "flex-row d-flex my-2 w-100";

        let labelTitle = document.createElement("label");
        labelTitle.htmlFor = "LibraryTitleInput";
        labelTitle.innerHTML = "Title:";
        titleInputRow.appendChild(labelTitle);

        let titleInput = document.createElement("input");
        titleInput.type = "text";
        titleInput.id = "LibraryTitleInput";
        titleInput.className = "form-control ms-2";
        titleInput.placeholder = "Enter library title";
        titleInput.required = true;
        titleInputRow.appendChild(titleInput);
        cardBody.appendChild(titleInputRow);

        //LOCATION
        let locationInputRow = document.createElement("div");
        locationInputRow.className = "flex-row d-flex my-2 w-100";

        let labelLocation = document.createElement("label");
        labelLocation.htmlFor = "LibraryLocationInput";
        labelLocation.innerHTML = "Location:";
        locationInputRow.appendChild(labelLocation);

        let locationInput = document.createElement("input");
        locationInput.type = "text";//"file";
        locationInput.id = "LibraryLocationInput";
        locationInput.className = "form-control ms-2";
        locationInput.placeholder = "Enter library location";
        //locationInput.webkitdirectory=true;
        //locationInput.directory=true;
        locationInput.required = true;
        locationInputRow.appendChild(locationInput);
        let locationPopUp = document.createElement("a");
        locationPopUp.href = "http://127.0.0.1:3000/openLocationApi";
        locationPopUp.target = "_blank";
        let imgIcon=document.createElement("img");
        imgIcon.src="/pics/folder.png";
        imgIcon.alt="Use this to open a folder and copy the folder path to paste in the input.  Sorry for the overzealous web security issues!"
        imgIcon.height=32;
        imgIcon.width-32;
        locationPopUp.appendChild(imgIcon);
        locationInputRow.appendChild(locationPopUp);
        cardBody.appendChild(locationInputRow);



        card.appendChild(cardBody);
        dialog.appendChild(card);
        let mdl = document.querySelector("#modal");
        mdl.appendChild(dialog);
        let cardBodyDiv=document.querySelector("#CardBodyEditLibrary");
        //JML
        let opInputRow={c:"flex-row d-flex my-2 w-100",b:[]};
        let opLabel={n:"label",for:"LibraryOperationInput",t:"Operation:"}
        let opInput={n:"input",i:"LibraryOperationInput",type:"text",value:"Folder", 
            required:true,c: "form-control ms-2",readonly:true };
        opInputRow.b.push(opLabel);
        opInputRow.b.push(opInput);
        let opInputRowHtml=jsonToHtml(opInputRow);
        cardBodyDiv.insertAdjacentHTML('beforeend',opInputRowHtml);
        
        let inherentInputRow={c:"flex-row d-flex my-2 w-100",b:[]};
        let inherentInput={n:"input",i:"LibraryInherentInput",type:"hidden",value:false };
        inherentInputRow.b.push(inherentInput);
        let inherentInputRowHtml=jsonToHtml(inherentInputRow);
        cardBodyDiv.insertAdjacentHTML('beforeend',inherentInputRowHtml);

        let buttonsDivRow={c:"d-flex justify-content-between my-2 w-100",b:[]};
        let closeButton={n:"button",type:"button",
            i:"EditImageLibraryCloseButton",
            title:"Close form do NOT save data",t:"Close",
            c:"btn btn-secondary text-warning mx-2"};
        buttonsDivRow.b.push(closeButton);
        let saveButton={n:"button",type:"button",
            i:"EditImageLibrarySaveButton",
            title:"Save Library",t:"Save",
            c:"btn btn-success text-warning mx-2"};
        buttonsDivRow.b.push(saveButton);
        let buttonsDivRowHtml=jsonToHtml(buttonsDivRow);
        cardBodyDiv.insertAdjacentHTML('beforeend',buttonsDivRowHtml);
        cardBodyDiv.insertAdjacentHTML('beforeend',inherentInputRowHtml);
        mdl.classList.remove("d-none");
        let addImageLibraryCloseButton=document.querySelector("#EditImageLibraryCloseButton");
        addImageLibraryCloseButton.addEventListener('click',async()=>{bang.closeAddLibraryForm();})
        let addImageLibrarySaveButton=document.querySelector("#EditImageLibrarySaveButton");
        addImageLibrarySaveButton.addEventListener('click',async()=>{bang.saveAddLibraryForm();})
        dialog.showModal();
    }


    async closeAddLibraryForm(){
        let bang=this;
        let existingDialog = document.querySelector("#AddLibraryModal");
        existingDialog.innerText="";
        existingDialog.remove();
        let appDiv = document.querySelector("#app");
        appDiv.classList.remove("d-none");        
    }

    async saveAddLibraryForm(){
        let bang=this;
        let existingDialog = document.querySelector("#AddLibraryModal");
        let useCheckbox = document.querySelector("#UseLibraryCheckbox").value;
        let nameInput = document.querySelector("#LibraryNameInput").value;
        let titleInput = document.querySelector("#LibraryTitleInput").value;
        let locationInput = document.querySelector("#LibraryLocationInput").value;
        let opInput = document.querySelector("#LibraryOperationInput").value;
        let inherentInput = document.querySelector("#LibraryInherentInput").value;
        let jsonUp={"use":useCheckbox, "name": nameInput, "title": titleInput,
            "location": locationInput, "operation": opInput}
        let apicallRtn=await bang.traffic.apiCall(bang.traffic.server + "/addImagesField",jsonUp)
        console.log(apicallRtn);
        bang.closeAddLibraryForm();
    }    

    async getFolderUrl(e) {
        let bang = this;
        try {
            e.preventDefault();
            // Show directory picker
            const directoryHandle = await window.showDirectoryPicker();
            
            // Display the selected folder path
            document.getElementById('folderPath').textContent = `Selected folder: ${directoryHandle.name}`;
        } catch (err) {
            console.error('Error selecting folder:', err);
        }
    }





    //******************************************************************************************
    //                                  Config Properties wired functions
    //******************************************************************************************
    
    /**
     * this handles the hiding of fonts when on random or showing when not
     * 
     * @param {any} e 
     * 
     * @memberOf dynamite
     */
    quoteFontRandomWiring(e) {
        let bang = this;
        if (e.target.id === "quoteFontRandom") {
            if (e.target.checked) {
                document.querySelector("#textFontFile").selectedIndex = 0;
                document.querySelector("#textFontFileEnvelope").classList.add("d-none");
            } else {
                document.querySelector("#textFontFileEnvelope").classList.remove("d-none");
            }
        }
    }
    
    /**
     * This handles the hiding of text color, background color and opacity when on random or showing when not
     * 
     * @param {any} e 
     * 
     * @memberOf dynamite
     */
    quoteAppearanceRandomWiring(e) {
        let bang = this;
        if (e.target.id === "quoteAppearanceRandom") {
            if (e.target.checked) {
                document.querySelector("#quoteAppearanceTextColorEnvelope").classList.add("d-none");
                document.querySelector("#quoteAppearanceBackgroundColorEnvelope").classList.add("d-none");
                document.querySelector("#quoteAppearanceOpacityEnvelope").classList.add("d-none");
            } else {
                document.querySelector("#quoteAppearanceTextColorEnvelope").classList.remove("d-none");
                document.querySelector("#quoteAppearanceBackgroundColorEnvelope").classList.remove("d-none");
                document.querySelector("#quoteAppearanceOpacityEnvelope").classList.remove("d-none");
            }
        }
    }


}

//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
//!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! S.D.G !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
//^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^