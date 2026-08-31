//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
//!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! J.J. !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
//^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

//----------------------------------------
//comms.js handles communications and for now DOM manipulation related to that.
import {jsonToHtml} from './cc/ccUtilities.js'
export default class comms{
    constructor(cfg){
        let traffic=this;
        traffic.controller=new AbortController();
        traffic.signal=this.controller.signal;
        traffic.server=cfg?.server||window.location.origin;
        traffic.imagesDiv=document.querySelector("#ImagesDiv");
        traffic.textLibs=document.querySelector("#textLibraries");
        traffic.picSumSave="";
        window.pic={};
        traffic.updatePicInfo();

         // Perform any additiona333l actions if the history is updated
         // Set up a 1-minute timer to call on currentInfoApi
         setInterval(async () => {
            await traffic.updatePicInfo();
        }, 10000); // 60000ms = 1 minute


    }

    async updatePicInfo() {
        let traffic=this;
        try {
            let response = await traffic.apiCall(traffic.server + "/currentInfoApi", "");
            //console.log("Pic history update check:", response);
            // Handle the response if needed
            //if (response && JSON.stringify(response) !== JSON.stringify(window.pic)) { // Check if the response is different from the current pic history
                console.log("Pic history has been updated.");
                if(response.imageItem.name==="PicSum"){
                    traffic.picSumSave=await response.saveName.replaceAll("pic0","picSumCache");
                }
                if((!Array.isArray(response.perScreenPics) || response.perScreenPics.length < 2)){
                    let cfg = traffic.config;
                    if(!cfg || cfg.differentWallpaperPerScreen === undefined){
                        cfg = await traffic.getConfig();
                    }
                    if(cfg && cfg.differentWallpaperPerScreen && Array.isArray(cfg.darwinPerScreenPicHistories) && cfg.darwinPerScreenPicHistories.length > 1){
                        response.perScreenPics = cfg.darwinPerScreenPicHistories;
                    }
                }
                //not necessarily needed
                window.pic = JSON.stringify(response);
                let currInfoLoading = document.querySelector("#currentInfoLoading");
                if (currInfoLoading) currInfoLoading.remove();
                traffic.currentInfoUpdate();

                // Sync progress bar with backend picUpdateCalled (e.g. systray-triggered changes)
                let bgEl = document.querySelector("#bgProgressEnvelope");
                if (bgEl) {
                    if (response.picUpdateCalled) {
                        bgEl.classList.remove("d-none");
                        bgEl.classList.add("d-flex");
                    } else {
                        bgEl.classList.remove("d-flex");
                        bgEl.classList.add("d-none");
                    }
                }

                // Perform any additional actions if the history is updated
            // }else{
            //     console.log("No change in pic history.");
            //     traffic.currentInfoUpdate();
            // }
        } catch (error) {
            console.error("Error checking pic history update:", error);
        }
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
            console.log("API Call Error:", error);
            return { error: error.message }; 
        }
    }

    //Specific Service for loading the config. Also generates the Images options
    // and sets the other inputs to current config
    async fetchConfig(){
        let traffic=this;
        let cfgData=await traffic.apiCall(traffic.server + "/configApi","");
        if (cfgData.picHistories && cfgData.picHistories.length > 0) {
            traffic.picSumSave=cfgData.picHistories[0].saveName.replaceAll("pic0","picSumCache")
        } else {
            traffic.picSumSave="";
        }
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
            ",backgroundChangingBlock,backgroundChangeAttempt,textLibraries,currentQuoteStatement" + 
            ",currentQuoteAuthor,picHistories,darwinPerScreenPicHistories,picUpdateCalled,version,published,mbcMonth,mbcValue,perScreenSupported".split(',');
//IF THE input value is null (below) it means the dontProcessFields is not used for a non-DOM config.
//Add the fields not to process above            
        for (const [key, value] of Object.entries(traffic.config)) {
            console.log(`${key}: ${value}`);
            if(!dontProcessFields.includes(key)){
                let input=document.querySelector("#" + key);
                if(input===null) {
                    console.log("input is null above. Its because the developer added a config field that is not a DOM element");
                    continue;
                }
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
        headerLocationDiv.className="col-6";
        headerLocationDiv.innerText="Location"
        headerLocationDiv.title="The URL or Folder Location of the image library"
        headerDiv.appendChild(headerLocationDiv);

        let headerApiDiv=document.createElement("div");
        headerApiDiv.id="HeaderRowApiDiv";
        headerApiDiv.className="col-2";
        headerApiDiv.innerText="API"
        headerApiDiv.title="API-based sources that require an API key"
        headerDiv.appendChild(headerApiDiv);



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
            locationDiv.className="col-6 d-flex align-items-center";

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

            //API column: dedicated cell for API-key sources.
            let apiDiv=document.createElement("div");
            apiDiv.id=configItem.name + "ImageItemRowApiDiv";
            apiDiv.className="col-2 d-flex align-items-center";
            if(configItem.requiresKey){
                let apiStatus=document.createElement("span");
                apiStatus.id="APIStatus" + configItem.name;
                apiStatus.className=configItem.hasApiKey ? "badge bg-DarkGreen me-2" : "badge bg-secondary me-2";
                apiStatus.innerText=configItem.hasApiKey ? "Key set" : "No key";

                let apiKeyButton=document.createElement("button");
                apiKeyButton.id="APIKey" + configItem.name;
                apiKeyButton.type="button";
                apiKeyButton.className="btn btn-sm btn-outline-secondary flex-shrink-0";
                apiKeyButton.dataset.name=configItem.name;
                apiKeyButton.title=configItem.hasApiKey ? "Replace API key" : "Enter API key";
                apiKeyButton.setAttribute("aria-label", apiKeyButton.title + " for " + configItem.title);
                apiKeyButton.innerHTML='<i class="bi bi-key" aria-hidden="true"></i>';
                apiKeyButton.addEventListener("click", async (event)=>{await traffic.apiKeyClicked(event);});

                apiDiv.appendChild(apiStatus);
                apiDiv.appendChild(apiKeyButton);
            }else{
                let apiNone=document.createElement("span");
                apiNone.className="text-muted";
                apiNone.innerText="\u2014"; // em dash: not an API source
                apiDiv.appendChild(apiNone);
            }
            rowDiv.appendChild(apiDiv);

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

    async apiKeyClicked(event){
        let traffic=this;
        let button=event.currentTarget;
        let apiKey=window.prompt("Enter the API key for " + button.dataset.name + ":", "");
        if(apiKey===null){
            return;
        }

        let result=await traffic.apiCall(traffic.server + "/imageApiKeyApi", {
            name: button.dataset.name,
            apiKey: apiKey
        });
        if(!result.error){
            let hasKey=!!apiKey.trim();
            button.title=hasKey ? "Replace API key" : "Enter API key";
            button.setAttribute("aria-label", button.title + " for " + button.closest(".libraryRow").config.title);
            let statusBadge=document.querySelector("#APIStatus" + button.dataset.name);
            if(statusBadge){
                statusBadge.className=hasKey ? "badge bg-DarkGreen me-2" : "badge bg-secondary me-2";
                statusBadge.innerText=hasKey ? "Key set" : "No key";
            }
        }
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
        headerUseDiv.className="col";
        headerUseDiv.innerText="Enabled"
        headerUseDiv.title="Determines if this Library will be included in the random selection"
        headerDiv.appendChild(headerUseDiv);

        let headerNameDiv=document.createElement("div");
        headerNameDiv.id="TLHeaderRowNameDiv";
        headerNameDiv.className="col";
        headerNameDiv.innerText="Library"
        headerNameDiv.title="The Name of the library to use"
        headerDiv.appendChild(headerNameDiv);

        let headerLocationDiv=document.createElement("div");
        headerLocationDiv.id="TLHeaderRowLocationDiv";
        headerLocationDiv.className="col";
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
            useDiv.className="col";
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
            nameDiv.className="col";

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
            locationDiv.className="col";

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
            if(configItem.name!=="MBC Values"){
                let useItem=document.querySelector("#TLUse" + configItem.name);
                useItem.addEventListener("change",async (e)=>{await window.comms.textLibraryChanged(e);});
            }
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


    async currentInfoUpdate(){
        let traffic=this;
        const currentPic=typeof(window.pic)==="object"?window.pic:JSON.parse(window.pic); // this is the current picture object from the server
        let picInfoEle=document.querySelector("#infoPic");
        if(!picInfoEle) return;
        picInfoEle.innerHTML = "";
        picInfoEle.insertAdjacentHTML("beforeend", jsonToHtml(traffic.buildPicCardJml(currentPic, "Current")));
        traffic.wirePicCards(picInfoEle);
    }

    // ---- New compact pic-card renderer (used by Current + History) --------

    // Escape text for safe insertion into HTML attributes/content.
    escapeHtml(value){
        return String(value == null ? "" : value)
            .replaceAll("&","&amp;").replaceAll("<","&lt;").replaceAll(">","&gt;")
            .replaceAll('"',"&quot;").replaceAll("'","&#39;");
    }

    // Return a browsable image src for a local path or http URL, or "" if none.
    picImageSrc(path){
        if(!path) return "";
        if(path.toLowerCase().startsWith("http")) return path;
        return "/imageFileApi?path=" + encodeURIComponent(path);
    }

    buildThumbJml(src, path, caption){
        let fileName=path ? path.split(/[\\/]/).pop() : "";
        let thumbBody = src
            ? {
                n:"img",
                src: src,
                c:"picThumb",
                alt: caption,
                "data-full": src,
                "data-caption": caption,
                loading:"lazy"
            }
            : {n:"div", c:"picThumb picThumbEmpty d-flex align-items-center justify-content-center", t:"n/a"};
        let fileLineChildren = [];
        if(fileName){
            fileLineChildren.push({n:"span", t:fileName});
        }
        return {
            c:"text-center",
            b:[
                ...(caption ? [{c:"small fw-bold text-info mb-1", t:caption}] : []),
                thumbBody,
                {c:"small text-truncate", s:"max-width:130px", b:fileLineChildren}
            ]
        };
    }


    getPerScreenPicsForCard(pic){
        let traffic=this;
        let perScreenPics = Array.isArray(pic?.perScreenPics)
            ? pic.perScreenPics.filter(screenPic => screenPic && (screenPic.saveName || screenPic.originName))
            : [];
        if(perScreenPics.length < 2 && traffic.config?.differentWallpaperPerScreen && pic?.picNum === 0){
            let darwinPerScreenPics = Array.isArray(traffic.config?.darwinPerScreenPicHistories)
                ? traffic.config.darwinPerScreenPicHistories.filter(screenPic => screenPic && (screenPic.saveName || screenPic.originName))
                : [];
            if(darwinPerScreenPics.length > 1){
                perScreenPics = darwinPerScreenPics;
            }
        }
        perScreenPics.sort((leftPic, rightPic)=>{
            let leftIndex = Number.isFinite(leftPic?.picNum) ? leftPic.picNum : Number.MAX_SAFE_INTEGER;
            let rightIndex = Number.isFinite(rightPic?.picNum) ? rightPic.picNum : Number.MAX_SAFE_INTEGER;
            return leftIndex - rightIndex;
        });
        return perScreenPics;
    }

    buildPerScreenColumnJml(pics, side){
        let traffic=this;
        let title = side === "original" ? "Original" : "Altered";
        let items = pics.map((screenPic)=>{
            let path = side === "original" ? (screenPic.originName || "") : (screenPic.saveName || "");
            return {
                c:"picVariantColumnItem",
                b:[traffic.buildThumbJml(traffic.picImageSrc(path), path, "")]
            };
        });
        return {
            c:"picVariantColumn",
            b:[
                {c:"small fw-bold text-info mb-2 text-center", t:title},
                ...items
            ]
        };
    }

    buildPerScreenGroupJml(pics, side){
        let traffic=this;
        let title = side === "original" ? "Original(s)" : "Altered";
        let headingClass = side === "original"
            ? "small fw-bold text-info mb-2 text-start"
            : "small fw-bold text-info mb-2 text-end";
        let groupClass = side === "original"
            ? "picStripGroup picStripGroupOriginal"
            : "picStripGroup picStripGroupAltered";
        let rowClass = side === "original"
            ? "picStripRow picStripRowOriginal"
            : "picStripRow picStripRowAltered";
        let items = pics.map((screenPic, idx)=>{
            let path = side === "original" ? (screenPic.originName || "") : (screenPic.saveName || "");
            return {
                c:"picStripItem",
                b:[traffic.buildThumbJml(traffic.picImageSrc(path), path, "")]
            };
        });
        return {
            c:groupClass,
            b:[
                {c:headingClass, t:title},
                {c:rowClass, b:items}
            ]
        };
    }

    buildPerScreenInfoRowsJml(pics){
        return {
            c:"picScreenInfoList mt-3",
            b:pics.map((screenPic, idx)=>{
                let displayIndex = Number.isFinite(screenPic?.picNum) ? Number(screenPic.picNum) + 1 : idx + 1;
                let filePath = screenPic.originName || screenPic.saveName || "";
                let fileName = filePath ? filePath.split(/[\\/]/).pop() : "";
                let filter = screenPic.filter || "original";
                return {
                    c:"picScreenInfoRow",
                    b:[
                        {n:"span", c:"picScreenInfoIndex", t:String(displayIndex)},
                        {n:"span", c:"picScreenInfoLabel", t:"Screen " + displayIndex},
                        {n:"span", c:"picScreenInfoSource", t:screenPic.imageItem?.name || ""},
                        {n:"span", c:"picScreenInfoFile", t:fileName},
                        {n:"span", c:"picScreenInfoFilter", t:filter}
                    ]
                };
            })
        };
    }

    buildPerScreenMetaJml(pics, label, quote, author){
        return {
            c:"picVariantMetaStack text-center px-2",
            b:[
                {c:"h6 text-Aquamarine mb-2", t:label},
                ...(quote ? [{c:"fst-italic small mb-1", t:'"' + quote + '"'}] : []),
                ...(author ? [{c:"small text-info mb-2", t:"- " + author}] : []),
                ...pics.map((screenPic, idx)=>({
                    c:"picVariantMetaLine",
                    b:[
                        {c:"small fw-bold text-warning", t:"Screen " + (idx + 1)},
                        {c:"small text-LightSalmon text-break", t:screenPic.imageItem?.name || ""},
                        {c:"small text-info text-break", t:screenPic.originName || ""}
                    ]
                }))
            ]
        };
    }
    buildChipJml(label, value, extraChildren){
        if(value === undefined || value === null || value === "") return null;
        let babies = [
            {n:"span", c:"fw-bold text-LightSalmon", t:label},
            {n:"span", c:"text-warning fst-italic ms-1", t:String(value)}
        ];
        if(Array.isArray(extraChildren) && extraChildren.length){
            babies.push(...extraChildren);
        }
        return {n:"span", c:"picChip", b:babies};
    }

    buildScreenVariantJml(pic, index){
        let originalPath = pic.originName || "";
        let alteredPath = pic.saveName || "";
        return {
            c:"picVariantRow",
            b:[
                this.buildThumbJml(this.picImageSrc(originalPath), originalPath, "Original"),
                {
                    c:"picVariantMeta text-center px-2",
                    b:[
                        {c:"h6 text-Aquamarine mb-1", t:"Screen " + (index + 1)},
                        {c:"small text-warning text-break", t:(pic.imageItem && pic.imageItem.name) ? pic.imageItem.name : ""},
                        {c:"small text-info text-break", t:originalPath}
                    ]
                },
                this.buildThumbJml(this.picImageSrc(alteredPath), alteredPath, "Altered")
            ]
        };
    }

    buildPicCardJml(pic, label){
        let traffic=this;
        let ii = pic.imageItem || {};
        let fontPath = pic.quoteFont || "";
        let fontName = fontPath ? fontPath.split(/[\\/]/).pop() : "";
        let perScreenPics = traffic.getPerScreenPicsForCard(pic);
        let showPerScreen = perScreenPics.length > 1;
        let quote = pic.quoteStatement || "";
        let author = pic.quoteAuthor || "";

        let visualSection;
        let summarySection = null;
        let perScreenInfoSection = null;
        if(showPerScreen){
            visualSection = {
                c:"picStripLayout",
                b:[
                    traffic.buildPerScreenGroupJml(perScreenPics, "original"),
                    traffic.buildPerScreenGroupJml(perScreenPics, "altered")
                ]
            };
            summarySection = {
                c:"text-center mt-2",
                b:[
                    {c:"h6 text-Aquamarine mb-1", t:label},
                    ...(quote ? [{c:"fst-italic small mb-1", t:'"' + quote + '"'}] : []),
                    ...(author ? [{c:"small text-info", t:"- " + author}] : [])
                ]
            };
            perScreenInfoSection = traffic.buildPerScreenInfoRowsJml(perScreenPics);
        } else {
            let originalPath = pic.originName || "";
            let alteredPath = pic.saveName || "";
            visualSection = {
                c:"d-flex justify-content-between align-items-start picVariantRow picVariantRowSingle",
                b:[
                    traffic.buildThumbJml(traffic.picImageSrc(originalPath), originalPath, "Original"),
                    {
                        c:"text-center flex-grow-1 px-2",
                        b:[
                            {c:"h6 text-Aquamarine mb-1", t:label},
                            ...(quote ? [{c:"fst-italic small", t:'"' + quote + '"'}] : []),
                            ...(author ? [{c:"small text-info", t:"- " + author}] : [])
                        ]
                    },
                    traffic.buildThumbJml(traffic.picImageSrc(alteredPath), alteredPath, "Altered")
                ]
            };
        }

        let chipRow = [
            traffic.buildChipJml("Source:", ii.name),
            traffic.buildChipJml("Info:", ii.title),
            traffic.buildChipJml("Operation:", ii.operation),
            traffic.buildChipJml("Filter:", pic.filter),
            traffic.buildChipJml("Sizing:", pic.sizing),
            traffic.buildChipJml("Font:", fontName)
        ].filter(Boolean);

        if(showPerScreen){
            chipRow.unshift({n:"div", c:"w-100 small text-secondary mb-1", t:"One picture per screen mode"});
        }

        return {
            c:"picCard bg-dark text-light rounded-3 p-2 mb-3 mx-2",
            b:[
                visualSection,
                ...(summarySection ? [summarySection] : []),
                ...(perScreenInfoSection ? [perScreenInfoSection] : []),
                {c:"d-flex flex-wrap mt-2 small", b:chipRow}
            ]
        };
    }

    // Wire up thumbnail click-to-fullsize within a container.
    wirePicCards(container){
        let traffic=this;
        container.querySelectorAll(".picThumb[data-full]").forEach(img=>{
            img.addEventListener("click", ()=>{
                traffic.showFullSizeImage(img.dataset.full, img.dataset.caption);
            });
        });
    }

    // Full-size image popup with a top-right X to close.
    showFullSizeImage(src, caption){
        if(!src) return;
        let existing=document.querySelector("#picFullSizeOverlay");
        if(existing) existing.remove();
        let overlay=document.createElement("div");
        overlay.id="picFullSizeOverlay";
        overlay.className="picFullSizeOverlay";
        overlay.innerHTML=''
            + '<button type="button" class="picFullSizeClose" aria-label="Close">&times;</button>'
            + '<img src="' + this.escapeHtml(src) + '" alt="' + this.escapeHtml(caption||"") + '" class="picFullSizeImg"/>';
        overlay.addEventListener("click",(e)=>{
            if(e.target===overlay || e.target.classList.contains("picFullSizeClose")){
                overlay.remove();
            }
        });
        document.addEventListener("keydown", function esc(e){
            if(e.key==="Escape"){ let o=document.querySelector("#picFullSizeOverlay"); if(o) o.remove(); document.removeEventListener("keydown", esc); }
        });
        document.body.appendChild(overlay);
    }

    // Render the History tab: current + last 9 (10 total), newest first.
    async renderHistory(){
        let traffic=this;
        let container=document.querySelector("#historyList");
        if(!container) return;
        let cfg=await traffic.getConfig();
        let histories=(cfg && cfg.picHistories) ? cfg.picHistories : [];
        let currentPic=typeof(window.pic)==="object"?window.pic:JSON.parse(window.pic || "{}");
        if(currentPic && (currentPic.saveName || currentPic.originName) && histories.length > 0){
            histories[0] = currentPic;
        }
        if(cfg && cfg.differentWallpaperPerScreen && histories.length > 0 && (!Array.isArray(histories[0].perScreenPics) || histories[0].perScreenPics.length < 2)){
            if(Array.isArray(cfg.darwinPerScreenPicHistories) && cfg.darwinPerScreenPicHistories.length > 1){
                histories[0].perScreenPics = cfg.darwinPerScreenPicHistories;
            }
        }
        if(histories.length===0){
            container.innerHTML='<div class="text-warning p-3">No history yet.</div>';
            return;
        }
        let cardsJml=[];
        histories.slice(0,10).forEach((pic, idx)=>{
            let label = idx===0 ? "Current" : "#" + (idx+1);
            cardsJml.push(traffic.buildPicCardJml(pic, label));
        });
        container.innerHTML="";
        container.insertAdjacentHTML("beforeend", jsonToHtml(cardsJml));
        traffic.wirePicCards(container);
    }

    async currentInfoUpdateOLD(){
        let traffic=this;
        const currentPic=typeof(window.pic)==="object"?window.pic:JSON.parse(window.pic); // this is the current picture object from the server
        let randomIIPicked={i:"ImageItemData",c:"bg-AliceBlue fw-bolder text-FireBrick mx-5",b:[]}
        let flexRow1={i:"ImageItemDataFlexRow1",c:"d-flex justify-content-between mb-3",b:[]}
        flexRow1.b.push(
            {i:"op",c:"p-2 fw-bold",t:"Operation: ",b:[{i:"opVal",c:"text-Maroon float-end ms-2 fst-italic",t: currentPic.imageItem.operation}]},
            {i:"ttl",c:"p-2 fw-bold",t:"Info: ",b:[{i:"titleVal",c:"text-Maroon float-end ms-2 fst-italic",t: currentPic.imageItem.title}]},
            {i:"loc",c:"p-2 fw-bold",t:"Location: ",b:[{i:"locVal","data-url":currentPic.imageItem.location,c:"text-Maroon float-end ms-2 fst-italic opencapable",t: currentPic.imageItem.location}]},
        )
        let flexRow2={i:"ImageItemDataFlexRow2",c:"d-flex justify-content-between mb-3",b:[]}
        flexRow2.b.push(
            {i:"inherent",c:"p-2 fw-bold",t:"Is Inherent: ",ttl:"If inherent this can NOT be changed!",b:[{i:"inherentVal",c:"text-Maroon float-end ms-2 fst-italic",t: currentPic.imageItem.inherent.toString()}]},
        )
        if(currentPic.imageItem.name==="PicSum"){
            flexRow2.b.push(
                {i:"name",c:"p-2 fw-bold",t:"Name: ",b:[{i:"nameVal",c:"text-Maroon float-end ms-2 fst-italic",t: currentPic.imageItem.name}]},
                {i:"use",c:"p-2 fw-bold",t:"Use this: ",ttl:"If true this Library is in use",b:[{i:"useVal","data-url": traffic.picSumSave,c:"text-Maroon float-end ms-2 fst-italic opencapable",t: currentPic.imageItem.use}]},
            )
        }else{
            flexRow2.b.push(
                {i:"name",c:"p-2 fw-bold",t:"Name: ",b:[{i:"nameVal",c:"text-Maroon float-end ms-2 fst-italic",t: currentPic.imageItem.name}]},
                {i:"use",c:"p-2 fw-bold",t:"Use this: ",ttl:"If true this Library is in use",b:[{i:"useVal","data-url":currentPic.imageItem.location,c:"text-Maroon float-end ms-2 fst-italic opencapable",t: currentPic.imageItem.use}]},
            )
        }

        randomIIPicked.b.push(flexRow1,flexRow2);
        let flexRow3={i:"ImageItemDataFlexRow3",c:"d-flex flex-row bg-dark",b:[]}
        let picSourceLink={};
        if(currentPic.originName.toLowerCase().startsWith("http")){
            picSourceLink={n:"a",i:"opOriginName",href:currentPic.originName,target:"_blank","title":"Click to see picture",t:currentPic.originName}
        }else{
            picSourceLink={i:"opOriginName","data-url":currentPic.originName,c:"text-Lavender float-end ms-2 fst-italic opencapable",t: currentPic.originName}
        }
        flexRow3.b.push(
            {i:"source",c:"d-flex p-2 fw-bold text-LightSalmon",t:"Picture Source: "
                ,b:[picSourceLink]},
        )
        randomIIPicked.b.push(flexRow3);
        let flexRow4={i:"ImageItemDataFlexRow4",c:"d-flex flex-row mb-3 bg-dark",b:[]}
        let picSavedLink={};
        if(currentPic.imageItem.name==="PicSum"){
            let picsumCache=currentPic.saveName.replaceAll("pic0.png","imgPicSumCache.png");
            // picSavedLink={n:"a",i:"saveNameVal",href:picsumCache,target:"_blank","title":"Click to see picture",t:picsumCache}
            picSavedLink={i:"saveNameVal","data-url":picsumCache,c:"text-Lavender float-end ms-2 fst-italic opencapable",t: picsumCache}

        }else{
            if(currentPic.saveName.toLowerCase().startsWith("http")){
                picSavedLink={n:"a",i:"saveNameVal",href:currentPic.saveName,target:"_blank","title":"Click to see picture",t:currentPic.saveName}
            }else{
                picSavedLink={i:"saveNameVal","data-url":currentPic.saveName,c:"text-Lavender float-end ms-2 fst-italic opencapable",t: currentPic.saveName}
            }
        }
        flexRow4.b.push(
            {i:"saved",c:"d-flex p-2 fw-bold text-LightSalmon",t:"Picture Saved: "
                ,b:[picSavedLink]},
        )
        randomIIPicked.b.push(flexRow4);

        let flexRow5={i:"ImageItemDataFlexRow5",c:"d-flex justify-content-between bg-dark",b:[]}
        flexRow5.b.push(
            {i:"op",c:"p-2 fw-bold",t:"Sizing/Scaling: ",b:[{i:"sizingVal",c:"text-warning float-end ms-2 fst-italic",t: currentPic.sizing}]},
            {i:"loc",c:"p-2 fw-bold",t:"Image Filter: ",b:[{i:"filterVal",c:"text-warning float-end ms-2 fst-italic",t: currentPic.filter}]},
        )
        randomIIPicked.b.push(flexRow5);
        //quoteFont
        // let flexRowQuoteFont={i:"ImageItemDataFlexRowQuoteFont",c:"d-flex justify-content-between bg-CornflowerBlue",b:[]}
        let fontFolder=currentPic.quoteFont.includes("/") 
            ? currentPic.quoteFont.split("/").slice(0, -1).join("/") 
            : currentPic.quoteFont.split("\\").slice(0, -1).join("\\"); // Extract the folder path if available
        // flexRowQuoteFont.b.push(
        //     {i:"fontfile",c:"p-2 fw-bold",t:"Font File: ",b:[{i:"fontFileVal","data-url":fontFolder,c:"text-warning float-end ms-2 fst-italic opencapable",t: currentPic.quoteFont}]},
        // )
        // randomIIPicked.b.push(flexRowQuoteFont);

        let fontRow=
            await traffic.currentInfoUpdateRow("QuoteFont","Font File:"
                    ,currentPic.quoteFont,"bg-CornflowerBlue", fontFolder);
        randomIIPicked.b.push(fontRow);


        let picInfoEle=document.querySelector("#infoPic");
        picInfoEle.innerHTML=""; // Clear previous content
        picInfoEle.insertAdjacentHTML("beforeend",jsonToHtml(randomIIPicked))

        for(let oc of document.querySelectorAll(".opencapable")){
            oc.addEventListener("click",async (e)=>{
                window.comms.locationClicked(e);
            })
        }

    }

    async currentInfoUpdateRow(parameter,title,value,color,dataUrl){
        let traffic=this;
        let flexRowItem={i:"FlexRow" + parameter,c:"d-flex justify-content-between bg-" + color + "",b:[]}
        let innerText={i:parameter+"Val",c:"text-warning float-end ms-2 fst-italic",t: value};     
        if(dataUrl){ 
            innerText["data-url"]=dataUrl; 
            innerText.c+=" opencapable"; // Add the opencapable class if dataUrl is provided
        }
        let text={i:"fontfile",c:"p-2 fw-bold",t:title,b:[innerText]};
        flexRowItem.b.push(text)
        return flexRowItem;
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
        console.log("formFieldChanged -> sending", traffic.server + "/inputApi", json);
        const response = await traffic.apiCall(traffic.server + "/inputApi",jsonString)
        console.log("formFieldChanged response", response);
        await traffic.fetchConfig();
    }
    //====================================================================
    //                                              END   API Items
    //====================================================================



}


//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
//!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! S.D.G !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
//^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^