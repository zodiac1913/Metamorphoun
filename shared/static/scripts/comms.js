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
        traffic.server=cfg?.server||"http://127.0.0.1:3000";
        traffic.imagesDiv=document.querySelector("#ImagesDiv");
        traffic.textLibs=document.querySelector("#textLibraries");
        traffic.picSumSave="";
        window.pic={};
        traffic.updatePicInfo();

         // Perform any additional actions if the history is updated
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
                //not necessarily needed
                window.pic = JSON.stringify(response);
                let currInfoLoading = document.querySelector("#currentInfoLoading");
                if (currInfoLoading) currInfoLoading.remove();
                traffic.currentInfoUpdate();
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
        traffic.picSumSave=cfgData.picHistories[0].saveName.replaceAll("pic0","picSumCache")
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
            ",currentQuoteAuthor,picHistories,picUpdateCalled,version,published".split(',');
//IF THE input value is null (below) it means the dontProcessFields is not used for a non-DOM config.
//Add the fields not to process above            
        for (const [key, value] of Object.entries(traffic.config)) {
            console.log(`${key}: ${value}`);
            if(!dontProcessFields.includes(key)){
                let input=document.querySelector("#" + key);
                if(input===null) console.log("input is null above. Its because the developer added a config field that is not a DOM element");
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


    async currentInfoUpdate(){
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