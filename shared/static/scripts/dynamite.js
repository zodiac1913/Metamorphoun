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
        bang.closeEditButton= document.querySelector("#CloseEditButton");
        bang.quoteToolsButton=document.querySelector("#QuoteToolsButton");
        bang.selectedImageLibrary=undefined;
        bang.config=null;
        bang.traffic=new traffic({server: window.location.origin});
    }

    wireConfigProperties(config){
        let bang=this;
        bang.config=config;
        let inputz=Array.from(document.querySelectorAll(".primaryInput"));
        let quoteFontRandom=document.querySelector("#quoteFontRandom");
        //version
        document.querySelector("#version").innerText=config.version;
        document.querySelector("#published").innerText=config.published;
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
        let mbcMode=document.querySelector("#mbcMode");
        bang.mbcModeWiring(mbcMode);

        quoteToolsButton.addEventListener("click",async (e)=>{ 
            let url = "quoteTools.html";
            window.open(url);
        });
        //wired
		for(let inpt of inputz){
			//set current
			inpt.value=config[inpt.id];
			//set change functions
			inpt.addEventListener("change",async (e)=>{
				await window.comms.formFieldChanged(e);
                if(e.target.id === "quoteFontRandom") {
                    bang.quoteFontRandomWiring(e);
                }
                if(e.target.id === "quoteAppearanceRandom") {
                    bang.quoteAppearanceRandomWiring(e);
                }
                if(e.target.id === "mbcMode") {
                    bang.mbcModeWiring(e.target);
                }

			});
		}
        //Web Actions (allow for web page instead of system tray)
        //Background Change
        document.querySelector("#callLastBackground").addEventListener("click",async (e)=>{
            bang.showBgProgress();
            try {
                let apicallRtn=await bang.traffic.apiCall(bang.traffic.server + "/lastBackgroundApi",{})
                console.log(apicallRtn);
            } finally { bang.hideBgProgress(); }
        });
        document.querySelector("#callNextBackground").addEventListener("click",async (e)=>{
            bang.showBgProgress();
            try {
                let apicallRtn=await bang.traffic.apiCall(bang.traffic.server + "/nextBackgroundApi",{})
                console.log(apicallRtn);
            } finally { bang.hideBgProgress(); }
        });
        //Favorites Menu
        document.querySelector("#FavsBGWith").addEventListener("click",async (e)=>{
            let apicallRtn=await bang.traffic.apiCall(bang.traffic.server + "/saveFavoriteApi",{type:"BG","save":"quoteOnBG"})
            console.log(apicallRtn);
        });
        document.querySelector("#FavsBGWithout").addEventListener("click",async (e)=>{
            let apicallRtn=await bang.traffic.apiCall(bang.traffic.server + "/saveFavoriteApi",{type:"BG","save":"noQuoteOnBG"})
            console.log(apicallRtn);
        });
        document.querySelector("#FavsBGWithquote").addEventListener("click",async (e)=>{
            let apicallRtn=await bang.traffic.apiCall(bang.traffic.server + "/saveFavoriteApi",{type:"Quote","save":"quote"})
            console.log(apicallRtn);
        });
        //End Web Actions


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
        if(bang.editLibraryButton.dataset.dataWired!=="true"){
            bang.editLibraryButton.addEventListener("click", (e) => {bang.popupEditLibrary();})
            bang.editLibraryButton.dataset.dataWired="true";
        }

        if(bang.selectedImageLibrary!==undefined) {
            bang.addLibraryButton.classList.add("d-none");
            bang.editLibraryButton.classList.remove("d-none");
            bang.removeLibraryButton.classList.remove("d-none");
            bang.useLibraryButton.classList.remove("d-none");
            bang.openFolderButton.classList.remove("d-none");
            bang.closeEditButton.classList.remove("d-none");
            bang.closeEditButton.addEventListener("click",async (e)=>{
                // Close the edit mode for the selected library
                bang.selectedImageLibrary.classList.remove("bg-CornflowerBlue", "imageSelected");
                bang.selectedImageLibrary=undefined; // Clear the selection
                // Reset buttons visibility
                bang.wireAndConfigureImageLibraryToolButtons();
            });                
        }else{
            bang.addLibraryButton.classList.remove("d-none");
            bang.editLibraryButton.classList.add("d-none");
            bang.removeLibraryButton.classList.add("d-none");
            bang.useLibraryButton.classList.add("d-none");
            bang.openFolderButton.classList.add("d-none");
            bang.closeEditButton.classList.add("d-none");
            bang.openFolderButton.addEventListener("click",async (e)=>{
                let loc=bang.selectedImageLibrary.querySelector(".locationurl").dataset.url;
                let dataUp={"id": bang.selectedImageLibrary.id,"loc":loc};
                let apicallRtn=await bang.traffic.apiCall(bang.traffic.server + "/openLocationApi",dataUp);
            });
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

        let cardHeader = {c:"card-header",t:"Add New Image Library",ttl:"Add parts of this library",b:[]};
        let cardBody = {i:"CardBodyAddLibrary",c:"card-body d-flex flex-column align-items-center",b:[]};
    
        //FORM FIELDS
        //USE
        let useCheckboxRow = {c:"flex-row d-flex my-2 w-100",b:[]}
        let useCheckbox = {n:"input",i:"UseLibraryCheckbox",type:"checkbox"};
        let useCheckboxName = {n:"label",for:"LibraryNameInput",t:"Use"}
        useCheckboxRow.b.push(useCheckbox,useCheckboxName);
        cardBody.b.push(useCheckboxRow);

            
        //NAME
        let nameInputRow = {c:"flex-row d-flex my-2 w-100",b:[]}
        let labelName = {n:"label",for:"LibraryNameInput",t:"Name"}
        nameInputRow.b.push(labelName);
        let nameInput = {n:"input",i:"LibraryNameInput",type:"text",c:"form-control ms-2"
            ,placeholder:"Enter library name",required:true};
        nameInputRow.b.push(nameInput);
        cardBody.b.push(nameInputRow);
    
        //TITLE
        let titleInputRow = {c:"flex-row d-flex my-2 w-100",b:[]};
        let labelTitle = {n:"label",for:"LibraryTitleInput",t:"Title:"}
        let titleInput = {n:"input",i:"LibraryTitleInput",type:"text",c:"form-control ms-2"
            ,placeholder:"Enter library title",required:true};
        titleInputRow.b.push(labelTitle,titleInput);
        cardBody.b.push(titleInputRow);

        //LOCATION
        let locationInputRow = {c:"flex-row d-flex my-2 w-100",b:[]};
        let labelLocation = {n:"label",for:"LibraryLocationInput",t:"Location:"}
        let locationInput = {n:"input",i:"LibraryLocationInput",type:"text",c:"form-control ms-2"
            ,placeholder:"Enter library location",required:true};
        locationInputRow.b.push(labelLocation,locationInput);

        
        let locationPopUp = {n:"a",href:"/openLocationApi",target:"_blank",b:[
            {n:"img",src:"/pics/folder.png",alt:"Use this to open a folder and copy the folder path to paste in the input.  Sorry for the overzealous web security issues!",height:32,width:32}
        ]};
        locationInputRow.b.push(locationPopUp);
        cardBody.b.push(locationInputRow);



        card.b.push(cardHeader,cardBody);
        dialog.b.push(card);
        let dialogHTM=jsonToHtml(dialog);
        let mdl = document.querySelector("#modal");
        mdl.insertAdjacentHTML("afterbegin",dialogHTM);
        let dialogEle=mdl.querySelector("dialog"); // Get the dialog element after inserting HTML
        mdl.classList.add("openPopup");
        let cardBodyDiv=document.querySelector("#CardBodyAddLibrary");
        //Operation
        let opInputRow={c:"flex-row d-flex my-2 w-100",b:[]};
        let opLabel={n:"label",for:"LibraryOperationInput",t:"Operation:"}
        let opInput={n:"input",i:"LibraryOperationInput",type:"text",value:"Folder", 
            required:true,c: "form-control ms-2"};
        opInputRow.b.push(opLabel);
        opInputRow.b.push(opInput);
        let opInputRowHtml=jsonToHtml(opInputRow);
        cardBodyDiv.insertAdjacentHTML('beforeend',opInputRowHtml);
        
        let inherentInputRow={c:"flex-row d-flex my-2 w-100",b:[]};
        let inherentInput={n:"input",i:"LibraryInherentInput",type:"hidden",value:false };
        inherentInputRow.b.push(inherentInput);
        let inherentInputRowHtml=jsonToHtml(inherentInputRow);
        cardBodyDiv.insertAdjacentHTML('beforeend',inherentInputRowHtml);

        //Allow (Harsh) Distortions
        let allowDistortCheckboxRow = {c:"flex-row d-flex my-2 w-100",b:[]}
        let allowDistortCheckbox = {n:"input",i:"AllowDistortLibraryCheckbox",type:"checkbox"};
        let allowDistortCheckboxName = {n:"label",for:"allowDistortNameInput",t:"Allow Harsh Distortions"
            ,ttl:"This forbids distortions like Dali or Vortex from distorting to a level that can be unacceptable to the user."}
        allowDistortCheckboxRow.b.push(allowDistortCheckbox,allowDistortCheckboxName);
        cardBodyDiv.b.push(allowDistortCheckboxRow);

                
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
        dialogEle.showModal();
    }

    async popupEditLibrary() {
        let bang = this;
        let selected=document.querySelector(".imageSelected");
        let data=selected.info;
        let appDiv = document.querySelector("#app");
        appDiv.classList.add("d-none");
        for(const mdl of document.querySelectorAll(".popupmdl")){ // remove any existing dialog boxes
            mdl.innerText="";
            mdl.classList.add("d-none");
        }

    
        // Check if the dialog already exists
        let existingDialog = document.querySelector("#EditLibraryModal");
        if (existingDialog) {
            existingDialog.showModal();
            return;
        }
    

        let dialog = {n:"dialog",i:"EditLibraryModal",c:"w-50",b:[]};
        let card = {i:"EditLibraryModalCard",c:"card bg-Bisque",b:[]};
        let cardHeader = {c:"card-header",t:"Edit Image Library",ttl:"Edit parts of this library",b:[]};
        if(data.inherent) {
            cardHeader.t = "Edit Image Library (Inherent)";
            cardHeader.ttl = "This is an inherent library and cannot be edited. You may only view its properties.";
        }
        let cardBody = {i:"CardBodyEditLibrary",c:"card-body d-flex flex-column align-items-center",b:[]};
    
        //FORM FIELDS
        //USE
        let useCheckboxRow = {c:"flex-row d-flex my-2 w-100",b:[]}
        let useCheckbox = {n:"input",i:"UseLibraryCheckbox",type:"checkbox"};
        if(data.use) useCheckbox.checked=true;
        let useCheckboxName = {n:"label",for:"LibraryNameInput",t:"Use"}
        useCheckboxRow.b.push(useCheckbox,useCheckboxName);
        cardBody.b.push(useCheckboxRow);

            
        //NAME
        let nameInputRow = {c:"flex-row d-flex my-2 w-100",b:[]}
        let labelName = {n:"label",for:"LibraryNameInput",t:"Name"}
        nameInputRow.b.push(labelName);
        let nameInput = {n:"input",i:"LibraryNameInput",type:"text",c:"form-control ms-2"
            ,placeholder:"Enter library name",required:true,value:data.name};
        if(data.inherent) nameInput.readonly=true;
        nameInputRow.b.push(nameInput);
        cardBody.b.push(nameInputRow);
    
        //TITLE
        let titleInputRow = {c:"flex-row d-flex my-2 w-100",b:[]};
        let labelTitle = {n:"label",for:"LibraryTitleInput",t:"Title:"}
        let titleInput = {n:"input",i:"LibraryTitleInput",type:"text",c:"form-control ms-2"
            ,placeholder:"Enter library title",required:true,value:data.title};
        if(data.inherent) titleInput.readonly=true;
        titleInputRow.b.push(labelTitle,titleInput);
        cardBody.b.push(titleInputRow);

        //LOCATION
        let locationInputRow = {c:"flex-row d-flex my-2 w-100",b:[]};
        let labelLocation = {n:"label",for:"LibraryLocationInput",t:"Location:"}
        let locationInput = {n:"input",i:"LibraryLocationInput",type:"text",c:"form-control ms-2"
            ,placeholder:"Enter library location",required:true,value:data.location};
        if(data.inherent) locationInput.readonly=true;
        locationInputRow.b.push(labelLocation,locationInput);

        
        let locationPopUp = {n:"a",href:"/openLocationApi",target:"_blank",b:[
            {n:"img",src:"/pics/folder.png",alt:"Use this to open a folder and copy the folder path to paste in the input.  Sorry for the overzealous web security issues!",height:32,width:32}
        ]};
        locationInputRow.b.push(locationPopUp);
        cardBody.b.push(locationInputRow);



        card.b.push(cardHeader,cardBody);
        dialog.b.push(card);
        let dialogHTM=jsonToHtml(dialog);
        let mdl = document.querySelector("#modal");
        mdl.insertAdjacentHTML("afterbegin",dialogHTM);
        let dialogEle=mdl.querySelector("dialog"); // Get the dialog element after inserting HTML
        let cardBodyDiv=document.querySelector("#CardBodyEditLibrary");

        //Operation
        let opInputRow={c:"flex-row d-flex my-2 w-100",b:[]};
        let opLabel={n:"label",for:"LibraryOperationInput",t:"Operation:"}
        let opInput={n:"input",i:"LibraryOperationInput",type:"text",value:"Folder", 
            required:true,c: "form-control ms-2",readonly:true };
        opInputRow.b.push(opLabel);
        opInputRow.b.push(opInput);
        let opInputRowHtml=jsonToHtml(opInputRow);
        cardBodyDiv.insertAdjacentHTML('beforeend',opInputRowHtml);
        
        //Inherent
        let inherentInputRow={c:"flex-row d-flex my-2 w-100",b:[]};
        let inherentInput={n:"input",i:"LibraryInherentInput",type:"hidden",value:false };
        inherentInputRow.b.push(inherentInput);
        let inherentInputRowHtml=jsonToHtml(inherentInputRow);
        cardBodyDiv.insertAdjacentHTML('beforeend',inherentInputRowHtml);

        //Allow (Harsh) Distortions
        let allowDistortCheckboxRow = {c:"flex-row d-flex my-2 w-100",b:[]}
        let allowDistortCheckbox = {n:"input",i:"AllowDistortLibraryCheckbox",type:"checkbox"};
        if(data.allowDistort) allowDistortCheckbox.checked=true;            
        let allowDistortCheckboxName = {n:"label",for:"allowDistortNameInput",t:"Allow Harsh Distortions"
            ,ttl:"This forbids distortions like Dali or Vortex from distorting to a level that can be unacceptable to the user."}
        allowDistortCheckboxRow.b.push(allowDistortCheckbox,allowDistortCheckboxName);
        let allowDistortCheckboxRowHtml=jsonToHtml(allowDistortCheckboxRow);
        cardBodyDiv.insertAdjacentHTML('beforeend',allowDistortCheckboxRowHtml);


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
        if(!data.inherent) {
            buttonsDivRow.b.push(saveButton);
        }
        let buttonsDivRowHtml=jsonToHtml(buttonsDivRow);
        cardBodyDiv.insertAdjacentHTML('beforeend',buttonsDivRowHtml);
        cardBodyDiv.insertAdjacentHTML('beforeend',inherentInputRowHtml);
        mdl.classList.remove("d-none");
        let editImageLibraryCloseButton=document.querySelector("#EditImageLibraryCloseButton");
        editImageLibraryCloseButton.addEventListener('click',async()=>{bang.closeEditLibraryForm();})
        if(!data.inherent) {
            let editImageLibrarySaveButton=document.querySelector("#EditImageLibrarySaveButton");
            editImageLibrarySaveButton.addEventListener('click',async()=>{bang.saveEditLibraryForm();})
        }
        dialogEle.showModal();
    }


    async closeAddLibraryForm(){
        let bang=this;
        let existingDialog = document.querySelector("#AddLibraryModal");
        existingDialog.innerText="";
        existingDialog.remove();
        let appDiv = document.querySelector("#app");
        appDiv.classList.remove("d-none");        
    }

    async closeEditLibraryForm(){
        let bang=this;
        let existingDialog = document.querySelector("#EditLibraryModal");
        existingDialog.innerText="";
        existingDialog.remove();
        let appDiv = document.querySelector("#app");
        appDiv.classList.remove("d-none");        
    }

    async saveAddLibraryForm(){
        let bang=this;
        let existingDialog = document.querySelector("#AddLibraryModal");
        let useCheckbox = document.querySelector("#UseLibraryCheckbox").value;
        let nameInput = document.querySelector("#LibraryNameInput").value.replaceAll(" ","");
        let titleInput = document.querySelector("#LibraryTitleInput").value;
        let locationInput = document.querySelector("#LibraryLocationInput").value;
        let opInput = document.querySelector("#LibraryOperationInput").value;
        let inherentInput = document.querySelector("#LibraryInherentInput").value;
        let allowDistortInput = document.querySelector("#AllowDistortLibraryCheckbox").value;
        let jsonUp={"use":useCheckbox, "name": nameInput, "title": titleInput,
            "location": locationInput, "operation": opInput, "inherent": inherentInput,
            "allowDistort": allowDistortInput
        }
        let apicallRtn=await bang.traffic.apiCall(bang.traffic.server + "/addImagesField",jsonUp)
        console.log(apicallRtn);
        bang.closeAddLibraryForm();
    }    

    async saveEditLibraryForm(){
        let bang=this;
        let existingDialog = document.querySelector("#EditLibraryModal");
        let useCheckbox = document.querySelector("#UseLibraryCheckbox").checked;
        let nameInput = document.querySelector("#LibraryNameInput").value.replaceAll(" ",""); // Remove spaces from the name input
        let titleInput = document.querySelector("#LibraryTitleInput").value;
        let locationInput = document.querySelector("#LibraryLocationInput").value;
        let opInput = document.querySelector("#LibraryOperationInput").value;
        let inherentInput = document.querySelector("#LibraryInherentInput").value;
        let allowDistortInput = document.querySelector("#AllowDistortLibraryCheckbox").value;
        let jsonUp={"use":useCheckbox, "name": nameInput, "title": titleInput,
            "location": locationInput, "operation": opInput,"inherent": inherentInput,
            "allowDistort": allowDistortInput
        }
        //Gotta add this to server
        let apicallRtn=await bang.traffic.apiCall(bang.traffic.server + "/editImagesField",jsonUp)
        console.log(apicallRtn);
        bang.closeEditLibraryForm();
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

    async mbcModeWiring(mbcModeEl) {
        let bang = this;
        let envelope = document.querySelector("#mbcValueEnvelope");
        if (mbcModeEl.checked) {
            envelope.classList.remove("d-none");
            await bang.loadMbcSelect();
        } else {
            envelope.classList.add("d-none");
        }
    }

    async loadMbcSelect() {
        let bang = this;
        let sel = document.querySelector("#mbcValueSelect");
        if (sel.options.length > 0) {
            // Already loaded, just sync selection
            sel.selectedIndex = bang.config.mbcValue || 0;
            return;
        }
        try {
            let res = await fetch("/quotes/mbc.json");
            let quotes = await res.json();
            sel.innerHTML = "";
            for (let i = 0; i < quotes.length; i++) {
                let opt = document.createElement("option");
                opt.value = i;
                let display = quotes[i].statement;
                opt.title = display;
                opt.innerText = display.length > 20 ? display.substring(0, 20) + "…" : display;
                sel.appendChild(opt);
            }
            sel.selectedIndex = bang.config.mbcValue || 0;
            sel.addEventListener("change", async () => {
                let json = { parameter: "mbcValue", value: parseInt(sel.value) };
                await bang.traffic.apiCall(bang.traffic.server + "/inputApi", json);
                await bang.traffic.fetchConfig();
            });
        } catch (err) {
            console.error("Failed to load MBC quotes:", err);
        }
    }

    showBgProgress() {
        let el = document.querySelector("#bgProgressEnvelope");
        if (el) {
            el.classList.remove("d-none");
            el.classList.add("d-flex");
            console.log("showBgProgress — visible");
        }
    }

    hideBgProgress() {
        let el = document.querySelector("#bgProgressEnvelope");
        if (el) {
            el.classList.remove("d-flex");
            el.classList.add("d-none");
            console.log("hideBgProgress — hidden");
        }
    }


}

//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
//!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! S.D.G !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
//^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^