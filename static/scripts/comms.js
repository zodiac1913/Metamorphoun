//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
//!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! J.J. !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
//^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

//----------------------------------------
//comms.js handles communications and for now DOM manipulation related to that.
export default class comms{
    constructor(cfg){
        let traffic=this;
        traffic.controller=new AbortController();
        traffic.signal=this.controller.signal;
        traffic.server=cfg?.server||"http://127.0.0.1:3000";
        traffic.imagesDiv=document.querySelector("#ImagesDiv");
        traffic.textLibs=document.querySelector("#textLibraries");
    }

    //====================================================================
    //Communications client/server
    //====================================================================

    //General API caller for service
    async apiCall(url, data, format = "json", method = "post") {
        let traffic = this;
        method = method.toLowerCase();

        let headers = new Headers();
        if (method === "post" && format === "json") {
            headers.set("Content-Type", "application/json");
            data = JSON.stringify(data);
        }

        try {
            const response = await fetch(url, {
                method: method.toUpperCase(),
                cache: "no-cache",
                body: data,
                signal: traffic.signal,
                headers: headers,
            });

            if (!response.ok) {
                // Check for specific status codes
                if (response.status === 400) {
                    console.error("Bad Request:", await response.text()); 
                } else if (response.status === 500) {
                    console.error("Internal Server Error");
                } else {
                    console.error(`HTTP error! status: ${response.status}`); 
                }
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            return await response.json(); 

        } catch (error) {
            console.error("API Call Error:", error);
            return { error: error.message }; 
        }
    }

    //Specific Service for loading the config. Also generates the Images options
    // and sets the other inputs to current config
    async fetchConfig(){
        let traffic=this;
        let cfgData=await traffic.apiCall(traffic.server + "/configApi","");
        traffic.config=cfgData;
        // if(traffic.imagesDiv.innerText.length<3){
        //     traffic.makeImagesItems();
        // }else{
            traffic.imagesDiv.innerText="";
            traffic.makeImagesItems();
        // }
        traffic.textLibs.innerText="";
        traffic.makeTextLibraryItems();
        //traffic.getFonts(); // this is a method to fetch fonts
        traffic.updateInputsFromConfig();
        return cfgData;
    }
    async getConfig(){
        let traffic=this;
        let cfgData=await traffic.apiCall(traffic.server + "/configApi","");
        traffic.config=cfgData;
        return cfgData;
    }
    //====================================================================
    //                                    END Communications client/server
    //====================================================================
  
    //====================================================================
    //DOM Manipulators
    //====================================================================
    async updateInputsFromConfig(){
        let traffic=this;
        let dontProcessFields="images,server_address,serverPort,originalCurrentBackgroundName," +
            "sourceCurrentBackgroundName,sourceCurrentBackgroundFolder," +
            "originalCurrentBackgroundFolder,currentBackgroundName,currentBackgroundFolder" +
            ",backgroundChangingBlock,textLibraries,currentQuoteStatement" + 
            ",currentQuoteAuthor,picHistories".split(',');
        for (const [key, value] of Object.entries(traffic.config)) {
            //console.log(`${key}: ${value}`);
            if(!dontProcessFields.includes(key)){
                let input=document.querySelector("#" + key);
                console.log(`${key}: ${value}`);
                if(input.nodeName === "SELECT"){
                    for(let option of input.options){
                        if(option.value === value){
                            option.selected=true;
                        }
                    }
                }else{
                    switch(input.type){
                        case "checkbox":
                            input.checked=value;
                            break;
                        default:
                            input.value=value;
                            break;
                    }
    
                }
            }else{
                console.log("key(" + key + ") not used");
            }
        }
    }

    async makeImagesItems(){
        let traffic=this;
        traffic.imagesDiv.innerHTML="";
        let headerDiv=document.createElement("div");
        headerDiv.id="HeaderRowDiv";
        headerDiv.className="row";
        let headerUseDiv=document.createElement("div");
        headerUseDiv.id="HeaderRowEnabledDiv";
        headerUseDiv.className="col-1";
        headerUseDiv.innerText="Enabled"
        headerUseDiv.title="Determines if this Library will be included in the random selection"
        headerDiv.appendChild(headerUseDiv);

        let headerNameDiv=document.createElement("div");
        headerNameDiv.id="HeaderRowNameDiv";
        headerNameDiv.className="col-3";
        headerNameDiv.innerText="Library"
        headerNameDiv.title="The Name of the library to use"
        headerDiv.appendChild(headerNameDiv);

        let headerLocationDiv=document.createElement("div");
        headerLocationDiv.id="HeaderRowLocationDiv";
        headerLocationDiv.className="col-8";
        headerLocationDiv.innerText="Location"
        headerLocationDiv.title="The URL or Folder Location of the image library"
        headerDiv.appendChild(headerLocationDiv);



        traffic.imagesDiv.appendChild(headerDiv);


        for(let configItem of traffic.config.images){
            //Row
            let rowDiv=document.createElement("div");
            rowDiv.id=configItem.name + "ImageItemRow";
            rowDiv.className="row libraryRow";
            rowDiv.info=configItem;
            //Checkbox div
            let useDiv=document.createElement("div");
            useDiv.id=configItem.name + "ImageItemRowUseDiv";
            useDiv.className="col-1";
            let useDivCheckBox=document.createElement("input");
            useDivCheckBox.id="Use" + configItem.name;
            useDivCheckBox.type="checkbox";
            useDivCheckBox.className="form-check-input imagesInput"
            useDivCheckBox.checked=configItem.use?"checked":"";
            useDivCheckBox.dataset.name=configItem.name;
            useDiv.appendChild(useDivCheckBox);
            rowDiv.appendChild(useDiv);
            //Name Div
            let nameDiv=document.createElement("div");
            nameDiv.id=configItem.name + "ImageItemRowNameDiv";
            nameDiv.className="col-3";

            let nameTitleDiv=document.createElement("div");
            nameTitleDiv.id="Name" + configItem.name;
            nameTitleDiv.className=""
            nameTitleDiv.innerText=configItem.title
            nameDiv.appendChild(nameTitleDiv);
            rowDiv.appendChild(nameDiv);
            traffic.imagesDiv.appendChild(rowDiv);
            //Location
            let locationDiv=document.createElement("div");
            locationDiv.id=configItem.name + "ImageItemRowLocationDiv";
            locationDiv.className="col-8";

            let locationTitleDiv=document.createElement("div");
            locationTitleDiv.id="Location" + configItem.name;
            locationTitleDiv.className=""
            //locationTitleDiv.innerText=configItem.location
            locationDiv.appendChild(locationTitleDiv);

            let anchor=document.createElement("a");
            anchor.id="Location" + configItem.name + "Anchor";
            anchor.classList.add("locationurl")
            //if(configItem.location.startsWith("http")){
                //anchor.href="#"//configItem.location;
                anchor.innerText=configItem.location;
                anchor.dataset.url=configItem.location;
            // }else{
            //     anchor.href="file://" + configItem.location;
            //     anchor.innerText=configItem.location;
            //}
            anchor.target="_blank";
            anchor.title="Click to open the location";
            locationTitleDiv.appendChild(anchor);
            rowDiv.appendChild(locationDiv);
            traffic.imagesDiv.appendChild(rowDiv);
            rowDiv.config=configItem;
            //wire
            let useItem=document.querySelector("#Use" + configItem.name);
            let linkItem=document.querySelector("#Location" + configItem.name + "Anchor");
            useItem.addEventListener("change",async (e)=>{await window.comms.imagesFieldChanged(e);});
            //linkItem.addEventListener("click",async (e)=>{await window.comms.locationClicked(e);});
        }
        // let imageDivRows = traffic.imagesDiv.querySelectorAll(".libraryRow");

        // for (let imageRow of imageDivRows) {
        //   imageRow.addEventListener("mouseover", (e) => {
        //     // Remove highlight from all rows
        //     imageRow.parentElement.querySelectorAll(".libraryRow").forEach(row => row.classList.remove("border", "border-CornflowerBlue", "border-1"));
        //     // Highlight the current row
        //     e.currentTarget.classList.add("border", "border-CornflowerBlue", "border-1");
        //   });
        //   imageRow.addEventListener("click",(e)=>{
        //     imageRow.parentElement.querySelectorAll(".libraryRow").forEach(row => 
        //         row.classList.remove("bg-CornflowerBlue","imageSelected")
        //     );
        //     e.currentTarget.classList.add("bg-CornflowerBlue","imageSelected")
        //     traffic.selectedImageLibrary=e.currentTarget;
        //   })
        // }
        window.tnt.wireImageLibraryTools();
    }

    async imagesFieldChanged(event){
      let traffic=this;
      let input=event.target;
      let id=input.dataset.name;
      let val=input.checked;
      let json={parameter: id, value: val};
      let jsonString=JSON.stringify(json);
      console.log(json);
      console.log(jsonString);
      let apicallRtn=await traffic.apiCall(traffic.server + "/imagesFieldChangeApi",json)
      //return await traffic.fetchConfig(); //flipping out
    }

    async locationClicked(event){
        let traffic=this;
        let anchor=event.target;
        let idVal=anchor.id;
        let locVal=anchor.dataset.url;
        let json={id: idVal, loc: locVal};
        let jsonString=json;
        console.log(json);
        console.log(jsonString);
        let apicallRtn=await traffic.apiCall(traffic.server + "/openLocationApi",jsonString)
        return await traffic.fetchConfig();
  
    }

    async makeTextLibraryItems(){
        let traffic=this;
        traffic.textLibs.innerHTML="";
        let headerDiv=document.createElement("div");
        headerDiv.id="TLHeaderRowDiv";
        headerDiv.className="row";
        let headerUseDiv=document.createElement("div");
        headerUseDiv.id="TLHeaderRowEnabledDiv";
        headerUseDiv.className="col-1";
        headerUseDiv.innerText="Enabled"
        headerUseDiv.title="Determines if this Library will be included in the random selection"
        headerDiv.appendChild(headerUseDiv);

        let headerNameDiv=document.createElement("div");
        headerNameDiv.id="TLHeaderRowNameDiv";
        headerNameDiv.className="col-3";
        headerNameDiv.innerText="Library"
        headerNameDiv.title="The Name of the library to use"
        headerDiv.appendChild(headerNameDiv);

        let headerLocationDiv=document.createElement("div");
        headerLocationDiv.id="TLHeaderRowLocationDiv";
        headerLocationDiv.className="col-8";
        headerLocationDiv.innerText="Citation"
        headerLocationDiv.title="The Citation and link to where it came from"
        headerDiv.appendChild(headerLocationDiv);



        traffic.textLibs.appendChild(headerDiv);


        for(let configItem of traffic.config.textLibraries){
            //Row
            let rowDiv=document.createElement("div");
            rowDiv.id="TL" + configItem.name + "ImageItemRow";
            rowDiv.className="row libraryRow";
            //Checkbox div
            let useDiv=document.createElement("div");
            useDiv.id="TL" + configItem.name + "ImageItemRowUseDiv";
            useDiv.className="col-1";
            let useDivCheckBox=document.createElement("input");
            useDivCheckBox.id="TL" + "Use" + configItem.name;
            useDivCheckBox.type="checkbox";
            useDivCheckBox.className="form-check-input imagesInput"
            useDivCheckBox.checked=configItem.use?"checked":"";
            useDivCheckBox.dataset.name=configItem.name;
            useDiv.appendChild(useDivCheckBox);
            rowDiv.appendChild(useDiv);
            //Name Div
            let nameDiv=document.createElement("div");
            nameDiv.id="TL" + configItem.name + "ImageItemRowNameDiv";
            nameDiv.className="col-3";

            let nameTitleDiv=document.createElement("div");
            nameTitleDiv.id="TLName" + configItem.name;
            nameTitleDiv.className=""
            nameTitleDiv.innerText=configItem.title
            nameDiv.appendChild(nameTitleDiv);
            rowDiv.appendChild(nameDiv);
            traffic.textLibs.appendChild(rowDiv);
            //Location
            let locationDiv=document.createElement("div");
            locationDiv.id="TL" + configItem.name + "ImageItemRowLocationDiv";
            locationDiv.className="col-8";

            let locationTitleDiv=document.createElement("div");
            locationTitleDiv.id="TL" + "Citation" + configItem.name;
            locationTitleDiv.className="bg-Aquamarine rounded-2"
            locationTitleDiv.style.maxWidth="22px"
            //locationTitleDiv.innerText=configItem.location
            locationDiv.appendChild(locationTitleDiv);

            let anchor=document.createElement("a");
            anchor.href=configItem.citation;
            anchor.target="_blank";
            anchor.title=configItem.name + " ---- " + configItem.info;
            let anchorIcon = document.createElement("img");
            anchorIcon.src = "/pics/citation.png";
            anchorIcon.alt = "Citation";
            anchorIcon.style.width = "20px";
            anchorIcon.style.height = "20px";
            anchor.appendChild(anchorIcon);
            locationTitleDiv.appendChild(anchor);
            rowDiv.appendChild(locationDiv);
            traffic.textLibs.appendChild(rowDiv);
            rowDiv.config=configItem;
            //wire
            let useItem=document.querySelector("#TLUse" + configItem.name);
            useItem.addEventListener("change",async (e)=>{await window.comms.textLibraryChanged(e);});
        }
 

    }

    async getFonts(){
        let traffic=this;
        let json={parameter: "getFonts"};
        let jsonString=json;
        let fonts=await traffic.apiCall(traffic.server + "/localFontApi",jsonString)
        let fontsSelect=document.querySelector("#textFontFile");
        fontsSelect.innerHTML="";
        // let rndOption=document.createElement("option");
        // rndOption.value="random";
        // rndOption.innerText="random";
        // fontsSelect.appendChild(rndOption)
        for(let font of fonts){
            let fontFile = font.split("\\").pop();
            let fontName = fontFile.split(".")[0];
            let option=document.createElement("option");
            option.value=font;
            option.innerText=fontFile;
            option.id=fontName+"FontOption";
            option.dataset.fontFile = fontFile;
            option.dataset.fontPath = font;
            if(font === traffic.config.textFontFile) {
                option.selected = true;
            }
            fontsSelect.appendChild(option);
        }
        //return await traffic.fetchConfig();
    }


    async textLibraryChanged(event){
    let traffic=this;
    let input=event.target;
    let id=input.dataset.name;
    let val=input.checked;
    let json={parameter: id, value: val};
    let jsonString=json;
    console.log(json);
    console.log(jsonString);
    let apicallRtn=await traffic.apiCall(traffic.server + "/textFieldChangeApi",jsonString)
    return await traffic.fetchConfig();

    }
    //====================================================================
    //                                              END   DOM Manipulators
    //====================================================================


    //====================================================================
    //                                                    API Items
    //====================================================================

    async formFieldChanged(event){
        let traffic=this;
        let input=event.target;
        let id=input.id;
        let val="";
        switch(input.type.toLowerCase()){
          case "checkbox":
            val=input.checked;
            break;
          default:
            val=input.value;
            break;
        }
        let json={parameter: id, value: val};
        let jsonString=json;
        console.log(json);
        console.log(jsonString);
        await traffic.apiCall(traffic.server + "/inputApi",jsonString)
        await traffic.fetchConfig();
    }
    //====================================================================
    //                                              END   API Items
    //====================================================================



}


//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
//!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! S.D.G !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
//^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^