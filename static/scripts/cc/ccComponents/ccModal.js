//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
//!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! J.J. !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
//^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
/**
 * Class for storing Row Checkboxes Configurations
  * Public Domain
 * Licensed Copyright Law of the United States of America, Section 105 (https://www.copyright.gov/title17/92chap1.html#105)
 * Per hoc, facies, scietis quod ille miserit me ut facerem universa quae cernitis et factis: Non est mecum!
 * Published by: Dominic Roche of OIT/IUSG/DASM on 08/01/2023
* @class ccModal
 * @extends {HTMLElement}
 */
"use strict";
class ccModal extends HTMLElement {
    //--------------------------------------------------------Fields
   
 
    //--------------------------------------------------------Fields END
 
    /**
     * Creates an instance of ccModal.
     * @param {*} config
     * @memberof ccModal
     */
    constructor(config) {
        super();
        let ccm = this;
        //Relevant Elements
        ccm._titleText=null;
        ccm.origin="";
        ccm._messageText=null;
        ccm._footText=null;
        ccm._closeOnBackgroundClick=true;
        ccm._hasCloseXButton=false;
        ccm._hasOverlay=true;
        ccm._buttons=[];
        //ccm.modalElement=null;
        //ccm.modalElement=document.querySelector("#" + ccm.id);
        ccm._overlayEle=null;
        ccm.homeOnModal=null;
        ccm.modalTitleElement=null;
        ccm.modalBodyPElement=null;
        ccm.addCloseButton=true;
        ccm.cfg=config||null;
        ccm.docsShown=false;
        //InAction vars
        ccm._isShown=false;
        ccm.isOnModal=false;
 
        //JML        
        ccm._headerJML=null;
        ccm._titleJML=null;
        ccm._titleButtonJML=null;
        ccm._bodyJML=null;
        ccm._modalJML=null;
        ccm._dialogJML=null;
        ccm._contentJML=null;
        ccm._footerJML=null;
        ccm.closeButtonJML={n:"button",type:"button",c:"btn btn-secondary",i:"ccModalFooterBtnClose","data-bs-dismiss":"modal",t:"Close",e: ccm.id + ".close();"};
 
        ccm.throttle=null;
        ccm.callButtonId=null;
        ccm.dialogSize="50";
    }
 
    //--------------------------------------------------------Methods
 
    //Web Component Lifecycle
   
    /**
     * Component connected to DOM
     *
     * @memberof ccModal
     */
    async connectedCallback() {
        let ccm=this;
        const { classAdd, getEle, isJson, makeAlert } = await import('../ccUtilities.js');
        ccm.getEle=getEle;
        ccm.isJson=isJson;
        ccm.makeAlert=makeAlert;
        ccm.classAdd=classAdd;
        //make id if missing
        if ((ccm?.id || "") === ""){
            ccm.id="GUID" + ccModal.guid().replaceAll("-", "");
        }
        //make docs if requested
        ccm.checkDocumentation();
        ccm.modalElement=document.querySelector("#" + ccm.id);
        if (!ccm.name) ccm.name = ccm.id;
        if (!ccm.title || ccm.title === "") ccm.title = ccm._messageText||ccm.message||"Message Not Defined";
        ccm.checkClassAdd(ccm.modalElement,"modal");    
        ccm.checkClassAdd(ccm.modalElement,"fade");
        window[ccm.id]=ccm;
        ccm.mdl=ccm;
    }
 
    static get observedAttributes() {
        return ["data-message","data-title","data-foot-text","data-background-click","data-title-close","data-overlay","data-buttons"];
    }
 
    /**
     * This updates class properties based on html tags being altered
     *
     * @param {*} name name of element tag
     * @param {*} oldValue old value for reference of element tag
     * @param {*} newValue new value of element tag
     * @memberof ccModal
     */
    attributeChangedCallback(name, oldValue, newValue) {
        let ccm=this;
        switch(name){
            case "data-message":
                ccm._messageText=newValue;
                break;
                case "data-title":
                    ccm._titleText=newValue;
                    break;
                case "data-foot-text":
                    ccm._footText=newValue;
                    break;
                case "data-background-click":
                    ccm._closeOnBackgroundClick=ccm.asBool(newValue);
                    break;
                case "data-title-close":
                    ccm._hasCloseXButton=newValue;
                    break;
                case "data-overlay":
                    ccm._hasOverlay=ccm.asBool(newValue);
                    break;
                case "data-buttons":
                    ccm._titleText=newValue;
                    break;
                case "data-documentation":
                    if(ccm.dataset.documentation)
                    {
                        ccm.showDocumentation()
                    }
                    else{
                        ccm.hideDocumentation();
                    }
                    break;
                default:
                    //whatimaboutadowitdat            
                    break;
        }
        window[ccm.id]=ccm;
    }
 
    //END Web Component Lifecycle
 
    /**
     * This version of open is to open a pre-created &lt;cc-modal&gt; element.
     * If you wish to create this modal on the fly use openFromConfig method
     *
     * @param {*} message message for the modal
     * @param {*} title title of the modal
     * @param {*} buttons JML data for buttons in the modal
     *                      example: {n:"button",type:"button",c:"btn btn-primary",
     *                                  i:"ccModalDoSomething",
     *                                  t:"Do It!!!",e: ccm.id + ".doSomething();"};
     * @param {*} footText any text you may want in the footer next to the buttons
     * @memberof ccModal
     */
    open(message,title,buttons,footText){
        let ccm=this;
        let messageIsHtml=false;
        //get the place they are (accessibility)
        ccm.homeOnModal=document.activeElement;
        if(ccm.modalElement.hasAttribute("data-message")){ ccm._messageText=ccm.modalElement.getAttribute("data-message");}
        if(ccm.modalElement.hasAttribute("data-title")){ ccm._titleText=ccm.modalElement.getAttribute("data-title");}
        if(ccm.modalElement.hasAttribute("data-foot-text")){ ccm._footText=ccm.modalElement.getAttribute("data-foot-text");}
        if(ccm.modalElement.hasAttribute("data-background-click")){ ccm._closeOnBackgroundClick=ccm.asBool(ccm.modalElement.getAttribute("data-background-click"));}
        if(ccm.modalElement.hasAttribute("data-title-close")){ ccm._hasCloseXButton=ccm.modalElement.getAttribute("data-title-close");}
        if(ccm.modalElement.hasAttribute("data-overlay")){ ccm.asBool(ccm._hasOverlay=ccm.modalElement.getAttribute("data-overlay"));}
        if(ccm.modalElement.hasAttribute("data-buttons")){
            const btns=ccm.modalElement.getAttribute("data-buttons");
            if(isJson(btns)){
                ccm._buttons=[...ccm._buttons,...JSON.parse(btns)];
            }
        }
        if(ccm.modalElement.hasAttribute("data-close-button")){ ccm.addCloseButton=asBool(ccm.modalElement.getAttribute("data-close-button"));}
        if(!message) message=ccm.messageText||"No Message Given! This is an Error.";
        ccm.messageText=message;        
        if(!title) title=ccm._titleText||"No Title Given! This is an Error.";
        ccm.titleText=title;
        if(buttons && ccm._buttons){
            ccm._buttons=[...ccm._buttons,...buttons];
        }else if(buttons && !ccm._buttons){
            ccm._buttons=[...ccm._buttons,...buttons];
        }else if(!buttons && !ccm._buttons){
            ccm._buttons=[];
        }else{
            //leave it
        }
        if(!footText) footText=ccm._footText||"";
        ccm.tabIndex=-1;
        if(!ccm._dialogJML) ccm._dialogJML={i: ccm.id + "ModalDialog",c:"modal-cc-dialog w-" + ccm.dialogSize,b:[]};
        if(!ccm._contentJML) ccm._contentJML={i: ccm.id + "ModalContent",c:"modal-content",b:[]};
        if(!ccm._headerJML) ccm._headerJML={i: ccm.id + "ModalHeader",c:"modal-header",b:[]};
        if(!ccm._titleJML) ccm._titleJML={i: ccm.id + "ModalTitle",c:"h3 m-auto","autofocus": null, t:ccm.titleText};
        ccm._headerJML.b.push(ccm._titleJML);
        if(ccm._hasCloseXButton){
            ccm._headerJML.b.push({i:ccm.id + "ModalTitleXCloseBtn",type:"button",c:"btn-close",
                "data-bs-dismiss":"modal","aria-label": "Close"});
        }
        ccm._contentJML.b.push(ccm._headerJML);
        if(!ccm._bodyJML) ccm._bodyJML={i: ccm.id + "ModalBody",c: "modal-body",b:[]};
        if(ccm.isJson(ccm.messageText)){
            ccm._bodyJML.b.push(ccm.messageText);
        }else if(this.isHtml(ccm.messageText)){
            messageIsHtml=true;
        }else{
            //text??
            ccm._bodyJML.b.push({n: "p",t: ccm.messageText,title: ccm.messageText })
        }
        ccm._contentJML.b.push(ccm._bodyJML);
        if(!ccm._footerJML) ccm._footerJML={i: ccm.id + "ModalFooter",c: "modal-footer",b:[]};
        if(ccm._footText) ccm._footerJML.b.push({n: "p",c:"fs-6 fw-bold",t:ccm._footText});
        if((ccm.buttons.length<1 || !ccm.buttons.find(b=>b.t==="Close")) && ccm.addCloseButton) ccm.buttons.push(ccm._closeButtonJML);
        for(const btn of ccm.buttons){
            ccm._footerJML.b.push(btn);
        }
        ccm._contentJML.b.push(ccm._footerJML);
        ccm._dialogJML.b.push(ccm._contentJML);
        ccm.modalElement.insertAdjacentHTML("beforeend",ccm.jsonToHtml(ccm._dialogJML));
        if(messageIsHtml) document.querySelector("#" + ccm.id + "ModalBody").insertAdjacentHTML("beforeend",ccm.messageText);
        ccm.dataset.src=message||"Unknown";
        ccm.show();
    }
 
    async openFromConfig(cfg,limit=1000){
        let ccm=this;
        //const { classAdd, getEle, isJson, jsonToHtml, makeAlert } = await import('../ccUtilities.js');
        //get the place they are (accessibility)
        ccm.homeOnModal=document.activeElement;
        //const func=ccm.openFromConfigMethod(cfg);
        if(cfg.origin) ccm.origin=cfg.origin;
        let timer;
        await ccm.openFromConfigMethod(cfg);
        // return (...args)=> {
        //     if(!timer){
        //         func.apply(this,args);
        //     }
        //     clearTimeout(timer);
        //     timer=setTimeout(()=>{
        //         timer=undefined;
        //     },limit);
        // };
    }
    /**
     * This version of open (openFromConfig) is to create this modal on the fly (no need for a
     * page to have the &lt;cc-modal&gt; element )
     * If you wish to just open a preexisting &lt;cc-modal&gt; element use the open() method.
     *
     * @param {*} cfg
     * @memberof ccModal
     */
    async openFromConfigMethod(cfg){
        let ccm=this;
        // const { classAdd, jsonToHtml } = await import('../ccUtilities.js');
        ccm.id=cfg.id||"ccModalPopup";
        ccm.dataset.src=cfg["data-src"]||"Unknown";
        ccm.dialogSize=cfg.dialogSize||"50";
        ccm.homeOnModal=document.activeElement;
        ccm.cfg=cfg;
        if(cfg.origin) ccm.origin=cfg.origin;
        if(cfg.callButtonId){
            ccm.callButtonId=cfg.callButtonId;
            ccm.callButton=document.querySelector("#" + ccm.callButtonId);
            ccm.callButton.disabled="disabled";            
        }
        if(cfg.class) ccm.className += cfg.class;
        ccm.titleText=cfg.titleText||"";
        ccm.messageText=cfg.messageText||"ALERT!!!";
        ccm.footText=cfg.footText||"";
        ccm.closeOnBackgroundClick=cfg.closeOnBackgroundClick;
        if(ccm.closeOnBackgroundClick===null||ccm.closeOnBackgroundClick==undefined) ccm.closeOnBackgroundClick=true;
        ccm.hasCloseXButton=cfg.hasCloseXButton||false;
        ccm._hasOverlay=(cfg.hasOverlay===null||cfg.hasOverlay==undefined)?true:cfg.hasOverlay;
        if(cfg.buttons===null){
            ccm.buttons=[];
        }else if(Array.isArray(cfg.buttons)){
            ccm.buttons=cfg.buttons;
        }else{
            ccm.buttons=[ccm.closeButton()];
        }
        ccm.addCloseButton=cfg.addCloseButton===true?true:false;
        ccm.headerJML=cfg.headerJML||null;
        ccm.titleJML=cfg.titleJML||null;
        ccm.titleButtonJML=cfg.titleButtonJML||null;
        ccm.bodyJML=cfg.bodyJML||null;
        ccm.modalJML=cfg.modalJML||{n: "cc-modal",i: ccm.id,c: "modal fade show unobtrusiveWait",b:[]};
        ccm.dialogJML=cfg.dialogJML||null;
        ccm.contentJML=cfg.contentJML||null;
        ccm.footerJML=cfg.footerJML||null;
        ccm.closeButtonJML=cfg.closeButtonJML||{n:"button",type:"button",c:"btn btn-secondary",i:"ccModalFooterBtnClose","data-bs-dismiss":"modal",t:"Close",e: ccm.id + ".closeForConfig();"};
        let messageIsHtml=false;
        ccm.tabIndex=-1;
        if(!ccm._dialogJML) ccm._dialogJML={i: ccm.id + "ModalDialog",c:"modal-cc-dialog w-" + ccm.dialogSize,b:[]};
        if(!ccm._contentJML) ccm._contentJML={i: ccm.id + "ModalContent",c:"modal-content  border-3 border-dark shadow shadow-xl rounded-3",b:[]};
        if(!ccm._headerJML) ccm._headerJML={i: ccm.id + "ModalHeader",c:"modal-header",b:[]};
        if(!ccm._titleJML) ccm._titleJML={i: ccm.id + "ModalTitle",c:"h3 m-auto modal-title", t:ccm.titleText};
        ccm._headerJML.b.push(ccm._titleJML);
        if(ccm._hasCloseXButton){
            ccm._headerJML.b.push({i:ccm.id + "ModalTitleXCloseBtn",type:"button",c:"btn-close",
                "data-bs-dismiss":"modal","aria-label": "Close"});
        }
        ccm._contentJML.b.push(ccm._headerJML);
        if(!ccm._bodyJML) {
            ccm._bodyJML={i: ccm.id + "ModalBody",c: "modal-body",b:[]};
            if(ccm.isJson(ccm.messageText)){
                ccm._bodyJML.b.push(ccm.messageText);
            }else if(ccm.isHtml(ccm.messageText)){
                messageIsHtml=true;
            }else{
                //text??
                ccm._bodyJML.b.push({n: "p",t: ccm.messageText,title: ccm.messageText })
            }
        }else{
            //bodyJML is already set and it should be fully done by dev
        }
        ccm._contentJML.b.push(ccm._bodyJML);
        if(!ccm._footerJML){
            ccm._footerJML={i: ccm.id + "ModalFooter",c: "modal-footer",b:[]};
            if(ccm._footText) ccm._footerJML.b.push({n: "p",c:"fs-6 fw-bold",t:ccm._footText});
            if((ccm.buttons.length<1 || !ccm.buttons.find(b=>b.t==="Close")) && ccm.addCloseButton) ccm.buttons.push(ccm.closeButton());
            for(const btn of ccm.buttons){
                ccm._footerJML.b.push(btn);
            }
        }else{
            //footerJML is already set and it should be fully done by dev
            if(!ccm.modalElement) ccm.modalElement=document.querySelector("#" + ccm.id);
        }
        ccm._contentJML.b.push(ccm._footerJML);
        ccm._dialogJML.b.push(ccm._contentJML);
        let ccModalGen=ccm.modalJML;
        ccModalGen.b.push(ccm._dialogJML);
        let modalEle=document.querySelector("#" + ccm.id);
        if(!modalEle){
            document.querySelector("body").insertAdjacentHTML('afterbegin',ccm.jsonToHtml(ccModalGen));
            ccm.modalElement=document.querySelector("#" + ccm.id);
        }else{
            ccm.modalElement=document.querySelector("#" + ccm.id);
        }
 
        ccm.modalElement.setAttribute("aria-hidden","false");
        ccm.modalElement.dataset.origin="openFromConfig";
        if(ccm.origin) ccm.modalElement.dataset.origin=ccm.origin;
        if(messageIsHtml) document.querySelector("#" + ccm.id + "ModalBody").insertAdjacentHTML('beforeend',ccm.messageText);
        await ccm.show();  
    }
   
    isJson(str) {
        if(!str) return false;
        try {
            if(typeof str==="object"){
                return true; //not sure this is ok.  Its a javascript object not json
            }else{
                JSON.parse(str);
            }
        } catch (e) {
            //console.info("this is not JSON (" + str +") because:");
            //console.info(e);
            return false;
        }
        return true;
    }
    isHtml(str) {
        const regex = /<([a-z]+)[^>]*>(.*?)<\/\1>/i; // Matches basic HTML tags (case-insensitive)
        return regex.test(str);
      }
    /**
     * This shows the modal.
     * NOTE!!! This is called by open() or openFromConfig().  You should not call this
     * since this modal relies on data in the element and destructs the inner elements
     * on close().
     *
     * @memberof ccModal
     */
    async show(){
        let ccm = this;
        if(!ccm.classList.contains("unobtrusiveWait")){
             ccm.makeBackDrop();
        }
        ccm.modalBody=ccm.modalElement.querySelector(".modal-body");
        ccm.modalElement.parentNode && ccm.modalElement.parentNode.nodeType === Node.ELEMENT_NODE || document.body.prepend(ccm.modalElement);
        ccm.modalElement.style.cssText="display: block; text-align: -webkit-center;";
 
        ccm.modalElement.setAttribute("role", "dialog");
        //Accessibility
        ccm.modalElement.setAttribute("aria-modal", !0);
        ccm.modalElement.setAttribute("aria-hidden","false");
        ccm.modalElement.setAttribute("aria-labelledby",(ccm.id+"ModalTitle"));        
        ccm.modalElement.setAttribute("aria-describedby",(ccm.id+"AriaModalDescription"));        
        ccm.modalElement.scrollTop = 0;
        ccm.modalBody && (ccm.modalBody.scrollTop = 0);
        ccm.modalElement.classList.remove('d-none');
        ccm._isShown=true;
        let desc=`>Beginning of modal dialog window. It begins with a title (in focus) \
                    of &quot;${ccm.titleText}&quot;. Its message/purpose is \
                    &quot;${ccm.messageText}&quot; Escape will cancel \
                    and close the modal. Tab to move about the modal it may \
                    have multiple buttons.`;
        let ariaModalSecription={i:(ccm.id+"AriaModalDescription"),c:"visually-hidden screen-reader-offscreen"
                                    ,t:desc};
        ccm.modalElement.insertAdjacentHTML('beforeend',ccm.jsonToHtml(ariaModalSecription));
        await ccm.wire();
    }
    /**
     * This is called by show to add all events listeners to the modal
     *
     * @memberof ccModal
     */
    async wire(){
        //wire it up
        let ccm=this;
        ccm.setResizeEvent();
        ccm.setEscapeEvent();
        ccm.setBackdropEvent();
        if(ccm._hasCloseXButton){
            if(ccm.dataset.src){
                document.querySelector("#" + ccm.id + "ModalTitleXCloseBtn").addEventListener("click",()=>ccm.closeForConfig());
            }else{
                document.querySelector("#" + ccm.id + "ModalTitleXCloseBtn").addEventListener("click",()=>ccm.close());
            }
        }
        if(ccm.buttons && ccm.buttons.length>0){
            for(const btn of ccm.buttons){
                if(typeof btn==='object'){
                    var btnFunc=new Function(btn.e);
                    let newBtn=document.querySelector("#" + btn.i);
                    if(newBtn) newBtn.addEventListener("click",btnFunc);
                }
            }
        }
        if(ccm._closeOnBackgroundClick){
            let mdlContent=document.querySelector("#" + ccm.id + "ModalContent");
            if(mdlContent!==null){
                    mdlContent.addEventListener("click",(e)=>{
                    ccm.isOnModal=true;
                    e.stopPropagation();
                });
            }
            document.querySelector("#" + ccm.id).addEventListener("click",(e)=>{
                if(!ccm.isOnModal && ccm._hasOverlay && ccm.closeOnBackgroundClick && ccm.cfg){
                    ccm.closeForConfig();
                }else{
                    ccm.close();
                }
                ccm.isOnModal=false;
                e.stopPropagation();
            });
        }
        //accessibility wiring
        ccm.modalTitleElement=  _.querySelector("#" + ccm.id + "ModalTitle");
        if(ccm.modalTitleElement){
            ccm.modalTitleElement.setAttribute("tabindex","0");
            ccm.modalTitleElement.focus();
        }
        ccm.modalBodyPElement= _.querySelector("#" + ccm.id + "ModalBody>p");
        if(ccm.modalBodyPElement){
            document.querySelector("#" + ccm.id + "ModalBody").insertAdjacentHTML("beforeend","<p></p>");
            ccm.modalBodyPElement= _.querySelector("#" + ccm.id + "ModalBod>p");
            if(ccm.modalBodyPElement) ccm.modalBodyPElement.setAttribute("tabindex","0");
        }
        ccm.trapTabbing();
        window[ccm.id]=ccm;
    }
    /**
     * This method makes the backdrop appear
     * NOTE!!! This is called by show().  You should not call this
     * since this modal relies on data in the element and destructs the inner elements
     * on close().
     * @memberof ccModal
     */
    makeBackDrop() {
        let ccm=this;
        if(ccm._hasOverlay) {
            ccm.modalElement.className+='show shadow-lg border border-3 border-dark';
            let body=document.querySelector("body");
            body.classList.add("modal-open");
            body.style="overflow: hidden; padding-right: 14px;";
            let backDrop={i: ccm.id + "ModalBackDrop", c: "modal-backdrop fade show"};
            body.insertAdjacentHTML('beforeend',ccm.jsonToHtml(backDrop));
        }else{
            //no backdrop and reduce modal
            ccm.modalElement.classList.add("show");
            ccm.modalElement.classList.add("border-0");
            ccm.modalElement.style.width="30%";
            ccm.modalElement.style.height="auto";
            ccm.modalElement.style.top="20%";
            ccm.modalElement.style.left="45%";
        }
    }
 
    /**
     * This closes the modal and wipes out its innards.
     * NOTE!!! You should not call this
     * since this modal relies on data in the element and destructs the inner elements
     * on close() and also destructs the modal itself on closeForConfig().
     * @return {*}
     * @memberof ccModal
     */
    close(){
        let ccm = this;
        if(ccm.modalElement.dataset.origin==="openFromConfig"){
            ccm.closeForConfig();
        }else{
            Array.from(document.querySelectorAll(".modal-backdrop")).map(e=>e.remove());
            if (!(ccm.modalElement.offsetParent === null || ccm.modalElement.classList.contains("d-none")))
            {
                console.log('ccModal no offset parent')
                return;
            }
            let body=document.querySelector("body");
            body.classList.remove("modal-open");
            ccm.unTrapTabbing();
            if(ccm.homeOnModal) ccm.homeOnModal.focus();
            body.style="";
            ccm.modalElement.setAttribute("aria-hidden","true");
            ccm.modalElement.style.display = "none";
            ccm.modalElement.classList.remove('show');
            ccm.modalElement.classList.add('d-none');
            let modalsExisting=document.querySelectorAll(ccm.id);
            if(modalsExisting.length>1){
                modalsExisting.shift();
                modalsExisting.forEach(m=>m.remove());
            }
            if(!ccm._hasOverlay) {
                //no backdrop and reduce modal
                ccm.modalElement.classList.remove("border-0");
                ccm.modalElement.style.width=null;
                ccm.modalElement.style.height=null;
                ccm.modalElement.style.top=null;
                ccm.modalElement.style.left=null;
                ccm.modalElement.style.borderWidth=null;
            }
            ccm._buttons=[];
            ccm._isShown=false;
            ccm._headerJML=null;
            ccm._titleJML=null;
            ccm._titleButtonJML=null;
            ccm._bodyJML=null;
            ccm._modalJML=null;
            ccm._dialogJML=null;
            ccm._contentJML=null;
            ccm._footerJML=null;
            ccm._footText=null;
            ccm.modalElement.innerText="";
        }
    }
    /**
     * This closes the modal (when opened with openForConfig()) and wipes out the modal.
     * NOTE!!! You should not call this because it destructs the modal itself.
     * @return {*}
     * @memberof ccModal
     */
    closeForConfig(){
        let ccm = this;
        ccm.unTrapTabbing();
        if(ccm.homeOnModal) ccm.homeOnModal.focus();
        Array.from(document.querySelectorAll(".modal-backdrop")).map(e=>e.remove());
        if (!(ccm.modalElement.offsetParent === null || ccm.modalElement.classList.contains("d-none")))  return;
        let body=document.querySelector("body");
        body.classList.remove("modal-open");
        ccm.modalElement.setAttribute("aria-hidden","true");
        body.style="";
        ccm.modalElement.remove();
        Array.from(document.querySelectorAll("#" + ccm.id)).map(e=>e.remove());
        ccm._buttons=[];
        ccm._isShown=false;
        ccm._headerJML=null;
        ccm._titleJML=null;
        ccm._titleButtonJML=null;
        ccm._bodyJML=null;
        ccm._modalJML=null;
        ccm._dialogJML=null;
        ccm._contentJML=null;
        ccm._footerJML=null;
        ccm._footText=null;
        ccm.remove();
    }
   
    hasCloseButton(){
        let ccm=this;
        let hasBtn=ccm.buttons.find(b=>b.t==="Close");
        return hasBtn?true:false;
    }
 
    /**
     * This sets the modal backdrop event if clicking on the modal overlay is set to drop the modal
     *
     * @memberof ccModal
     */
    setBackdropEvent(){
        let ccm = this;
        if(ccm._hasOverlay && ccm._closeOnBackgroundClick && ccm._isShown && !ccm.classList.contains("unobtrusiveWait")){
            document.querySelector("#" + ccm.id + "ModalBackDrop").addEventListener('click', ()=>ccm.close());
        }
    }
    /**
    * This sets the event that handles resizing of the window when a modal is open
     *
     * @memberof ccModal
     */
    setResizeEvent() {
        let ccm = this;
        //ccm._isShown ? j.on(window, Di, (()=>this._adjustDialog())) : j.off(window, Di)
        window.addEventListener('resize', ()=>ccm.adjustDialog());
    }
    /**
     * Sets the even to close the modal on the escape key being pressed
     *
     * @memberof ccModal
     */
    setEscapeEvent(){
        let ccm = this;
        if (ccm._isShown) {
            document.addEventListener("keyup",(e)=>ccm.escapeModal(e));
        }        
    }
    /**
     * Trap tabbing while Modal is open
     *
     * @memberof ccModal
     */
    trapTabbing(){
        let ccm = this;
        document.addEventListener("keyup", ccm.trapTab);
    }
   
    /**
     *Remove the tab trapping when the modal shuts down
     *
     * @memberof ccModal
     */
    unTrapTabbing(){
        let ccm = this;
        document.removeEventListener("keyup", ccm.trapTab);
    }
 
    /**
     * Tab trapping event
     *
     * @param {*} e
     * @memberof ccModal
     */
    trapTab=(e)=>{
        let ccm = this;
        if (e.keyCode===9) {
            if(ccm.modalElement.contains(e.target)){
 
            }else{
                e.preventDefault();
                e.stopImmediatePropagation();
                e.stopPropagation();
                ccm.modalTitleElement.focus();
            }
        }        
    }
 
    /**
     * This adjusts the modal on window size change
     *
     * @memberof ccModal
     */
    adjustDialog(){
        let ccm = this;
        let body=document.querySelector("body");
        let widthy=Math.abs(window.innerWidth - document.documentElement.clientWidth);
        const isModalOverflowing = ccm.modalElement.scrollHeight > document.documentElement.clientHeight
            , scrollbarWidth = widthy
            , isBodyOverflowing = scrollbarWidth > 0;
            if ((!isBodyOverflowing && isModalOverflowing) || (isBodyOverflowing && !isModalOverflowing)) {
                this.modalElement.style.paddingLeft = `${scrollbarWidth}px`
            }
       
            if ((isBodyOverflowing && !isModalOverflowing) || (!isBodyOverflowing && isModalOverflowing)) {
            this.modalElement.style.paddingRight = `${scrollbarWidth}px`
            }
    }
    /**
     * Closes the modal when the escape key is hit
     *
     * @param {*} e
     * @memberof ccModal
     */
    escapeModal(e){
        let ccm = this;
        if(e.code==='Escape'){
            if(ccm.dataset.src){
                ccm.closeForConfig();
            }else{
                ccm.close();
            }
        }
    }
 
    jsonToHtml(jsonIn, tabs) {
        let ccm = this;
        if (typeof tabs === "undefined") tabs = 2;
        let rtn = "";
        if (Array.isArray(jsonIn)) {
            for(const jObj of jsonIn){
                rtn += ccm.jsonToHtml(jObj, 1);
            }
        } else {
            let parentEle =   "<" + (jsonIn.nodeType || jsonIn.node || jsonIn.n || "div").toLowerCase();
            let children = "";
            let hasChildren = jsonIn.hasOwnProperty("babies")||jsonIn.hasOwnProperty("b");
            let nodeType=(jsonIn.nodeType || jsonIn.node || jsonIn.n || "div").toLowerCase();
            for (let key in jsonIn) {
                if (key !== "nodeType" && key !== "node" && key !== "n") {
                    if (jsonIn.hasOwnProperty(key)) {
                        if (typeof jsonIn[key] === "object" && jsonIn[key] !== null && key != "e" && key != "event") {
                            for(const obj of jsonIn[key]){
                                children += ccm.jsonToHtml(obj,tabs + 1);
                            }
                        } else if(key==="event" || key==="e"){
                            //do nothing.  This is for storing events in JML
                        } else {
                            if (key !== "innerText" && key !== "text" && key !== "innerHTML"  && key !== "t" && jsonIn[key] != undefined && jsonIn[key] != null) {
                                let q=(jsonIn[key].toString().indexOf('"')>-1?"'":'"');
                                let fullKey=key;
                                if(key==="c") fullKey="class";
                                if(key==="i") fullKey="id";
                                if (key === "v") fullKey = "value";
                                if (key === "o") fullKey = "onclick";
                                if(key==="s") fullKey="style";
                                if(key==="ttl") fullKey="title";
                                if(key==="p") fullKey="placeholder";
                                if(key==="acon") fullKey="aria-controls";
                                if(key==="aexp") fullKey="aria-expanded";
                                if(key==="ahdr") fullKey="aria-header";
                                if(key==="ahid") fullKey="aria-hidden";
                                if(key==="alab") fullKey="aria-label";
                                if(key==="albb") fullKey="aria-labelledby";
                                if(key==="alvl") fullKey="aria-level";
                                if(key==="arol") fullKey="aria-role";
                                parentEle += " " + fullKey + ((jsonIn[key] !== null && jsonIn[key] !== undefined)? ('=' + q + jsonIn[key] + q): "");
                            }
                        }
                    }
                }
            }
            let innerTxt = "";
            innerTxt = (jsonIn.t || jsonIn.innerText || jsonIn.text || jsonIn.innerHTML || "");
            if (hasChildren) {
                rtn =  parentEle + " >" + innerTxt + "\r\n" + "\t".repeat(tabs) + children + "\r\n" + "\t".repeat(tabs) +
            "</" + nodeType + ((jsonIn.hasOwnProperty("id")||jsonIn.hasOwnProperty("i")) ? ' data-end="' + (jsonIn.id||jsonIn.i) + '"' : "") + " >";
            } else {
                //console.log(jsonIn.i||JSON.stringify(jsonIn));
                if (innerTxt === "" &&  !(nodeType === "textarea")) {
                    if(nodeType !== "div") {
                        rtn=parentEle + " />\r\n" + "\t".repeat(tabs);
                    }else{
                        rtn=parentEle + ">" + innerTxt + "</" + nodeType + (jsonIn.hasOwnProperty("id") ? ' data-end="' + jsonIn.id + '"' : "") + " >" + "\t".repeat(tabs);
                    }
                } else {
                    rtn = parentEle + " >" + innerTxt + "\r\n" + "\t".repeat(tabs) + "</" + nodeType + " >";
                }
            }
        }
        if(rtn)
        return rtn.replaceAll("\t", "").replaceAll("\r", "").replaceAll("\n", "");
    }
    /**
     * checks if an element has a class and if not
     * adds it
     *
     * @param {*} element element to check and/or add for class
     * @param {*} cls class to add
     * @memberof ccModal
     */
    checkClassAdd(element,cls){
        if(element){
            if(!element.classList.contains(cls)){
                element.classList.add(cls);
            }
        }
    }
    /**
     * checks if an element has a class and if so
     * removes it
     *
     * @param {*} element element to check and/or remove for class
     * @param {*} cls class to remove
     * @memberof ccModal
     */
    checkClassRemove(element,cls){
        if(!element.classList.contains(cls)){
            element.classList.add(cls);
          }
    }
 
    /**
     * This creates the JML for alerts for missing element tags and problems with this web component
     *
     * @param {*} idAdd id of alert
     * @param {*} text text of alert
     * @param {*} clazz class of alert
     * @return {*}
     * @memberof ccModal
     */
    makeAlertJML(idAdd, text, clazz) {
        let ccac = this;
        if (!idAdd) idAdd = "";
        if (!clazz) clazz = "alert alert-danger alert-dismissible fade show";
        var alert={i:ccModal.guid().replaceAll("-", "") + idAdd, title: text, t: text,c: clazz };
        if(clazz.includes("alert-dismissible")){
          alert.b=[];
          let alertButton={n:"button",type:"button",title:"Close Alert",c: "btn-close","data-bs-dismiss": "alert"};
          alert.b.push(alertButton);
        }
        return alert;
      }      
    /**
     * Creates alerts when required attributes are not given
     *
     * @param {*} attribute The attribute that is missing
     * @param {*} details that it is missing and why this is a required attribute
     * @memberof ccModal
     */
    missingAttribute(attribute,details){
        let ccm = this;
        if(!document.querySelector("[id$='" + attribute + "']")){
            ccm.modalElement.insertAdjacentHTML('beforebegin',ccm.makeAlert("no" + ccm.asPascalCase(attribute),details));
        }
        ccm.attributeIssue=true;
      }
    /**
     * checks if web component documentation should be shown
     *
     * @memberof ccModal
     */
    checkDocumentation(){
        let ccm=this;
        if(!ccm.modalElement)ccm.modalElement=document.querySelector("#" + ccm.id);
        if (ccm.dataset.documentation === true ||ccm.dataset.documentation === "true" && !document.querySelector("#ccModalDocs")){
            ccm.showDocumentation();
        }
    }
 
    /**
     * Creates the documentation card before the modal element
     *
     * @memberof ccModal
     */
    showDocumentation(){
        let ccm=this;
        ccm.modalElement.insertAdjacentHTML('beforebegin', ccm.jsonToHtml(ccm.makeDocumentation()));
        document.querySelector("#docCardHeadButton").addEventListener("click",(e)=>{ ccm.hideDocumentation();});
    }
 
    /**
     * Hide the Documentation (removes the documentation card)
     *
     * @memberof ccModal
     */
    hideDocumentation(){
        let ccm=this;
        document.querySelector("#ccModalDocs").remove();
    }
 
    /**
     * this generates the documentation for this element in the parent of the element
     *
     * @return {*}
     * @memberof ccModal
     */
    makeDocumentation(){
        let ccm=this;
        let docCard={i:"ccModalDocs",c:"card",b:[]};
        let docCardHead={c:"card-header h3",t:"CatsCRUDL &lt;cc-modal&gt; Modal Web Component Documentation",b:[]};
        docCardHead.b.push({n:"button",type:"button",i:"docCardHeadButton",c:"btn-close float-end","aria-label":"Close"});
        docCard.b.push(docCardHead);
        docCard.b.push(ccm.makeAlertJML("DocDisable","To remove this div remove the cc-data-documentation=\"true\" attribute from your cc-modal element!","alert alert-info alert-dismissible fade show"))
        let docCardBody={c:"card-body",b:[]};
        let cardTitle={n:"h5",c:"card-title",t:"Element Attributes"};
        docCardBody.b.push(cardTitle);
        docCardBody.b.push({n:"p",c:"m-0 p-0",t:"&lt;cc-modal  details on tags follow "});
        docCardBody.b.push({n:"p",c:"m-0 ms-5 p-0",t:"&Tab;&Tab; id=\"someIdName\" The base id for this modal. If you do not set this it will be a guid "});
        docCardBody.b.push({n:"p",c:"m-0 ms-5 p-0",t:"&Tab;&Tab;data-title=\"The text for the BS5 Modal Title\"  Label (Title) for this input.  If you don`t set this it will default to ATTENTION!!! "});
        docCardBody.b.push({n:"p",c:"m-0 ms-5 p-0",t:"&Tab;&Tab;data-documentation=\"true\"  Shows the Documentation for the attributes and events of a cc-modal "});
        docCardBody.b.push({n:"p",c:"m-0 ms-5 p-0",t:"&Tab;&Tab;data-message=\"The Modal Message\"  This is the text and purpose of this Modal "});
        docCardBody.b.push({n:"p",c:"m-0 ms-5 p-0",t:"&Tab;&Tab;data-foot-text=\"The Footer Text if you want any\"  This is the text for the footer of the modal. Its not normally used. "});
        docCardBody.b.push({n:"p",c:"m-0 ms-5 p-0",t:"&Tab;&Tab;data-background-click=true  This allows clicking the background(overlay) to close the modal if true. (defaults to true)  "});
        docCardBody.b.push({n:"p",c:"m-0 ms-5 p-0",t:"&Tab;&Tab;data-title-close=false  This determines if there is a close button in the title (default is false (no button)) "});
        docCardBody.b.push({n:"p",c:"m-0 ms-5 p-0",t:"&Tab;&Tab;data-overlay=true  This determines if a overlay is set over the page (default is true) "});
        docCardBody.b.push({n:"p",c:"m-0 ms-5 p-0",t:"&Tab;&Tab;data-close-button  This determines if the default close button is added. Default is true "});
        docCardBody.b.push({n:"p",c:"m-0 ms-5 p-0",t:"&Tab;&Tab;data-buttons  This allow a tag for JML buttons that will show in the footer.  There is no need to make a close button as that will be there by default "});
        let hr={n:"hr",style:"height: 12px;border: 0;box-shadow: inset 0 12px 12px -12px rgba(0, 0, 0, 0.5);"};
        docCardBody.b.push(hr);
        let cardTitle2={n:"h5",c:"card-title",t:"Functions"};
        docCardBody.b.push(cardTitle2);
        docCardBody.b.push({n:"p",c:"m-0 p-0",t:" Class Methods: "});
        docCardBody.b.push({n:"p",c:"m-0 ms-2 p-0",t:"&Tab;The following methods can be called via the class  "});
        docCardBody.b.push({n:"p",c:"m-0 ms-5 p-0",t:"&Tab;&Tab;The class name will be its html id (usually #ccModal "});
        docCardBody.b.push({n:"p",c:"m-0 ms-5 p-0",t:"&Tab;&Tab;for element created or ccModalPopup for ones created with"});
        docCardBody.b.push({n:"p",c:"m-0 ms-5 p-0",t:"&Tab;&Tab;the openForConfig method)"});
        docCardBody.b.push({n:"p",c:"m-0 ms-2 p-0",t:"&Tab;open() -- opens the modal and allows closing based on config"});
        docCardBody.b.push({n:"p",c:"m-0 ms-2 p-0",t:"&Tab;openForConfig() -- This allows for a config with full control and defaults "});
        docCardBody.b.push({n:"p",c:"m-0 ms-5 p-0",t:"&Tab;&Tab;the easiest way to use this is to call the window.ccModalConfigOpen()  "});
        docCardBody.b.push({n:"p",c:"m-0 ms-5 p-0",t:"&Tab;&Tab;or the ccModalMakeOpen() function.  See Functions Below "});
 
        docCardBody.b.push(hr);        
        docCardBody.b.push({n:"span",c:"m-0 p-0 h5",t:" Global Function(s): "});
        docCardBody.b.push({n:"span",c:"m-2 p-0 fw-bold fs-4 d-block",t:"ccModalMakeOpen(message, title,buttons,footer) -- creates and opens a modal  "});
        let ccModalMakeOpenParamsUl={n:"ul",c:"ms-5 p-0 list-group list-group-flush list-unstyled",t:"---Parameters:",b:[]};
        ccModalMakeOpenParamsUl.b.push({n:"li",c:"ms-5 p-0 list-group-flush",t:"1. message: The main message of the modal (in the body).[Default=Sorry! No reason for this alert!]"});
        ccModalMakeOpenParamsUl.b.push({n:"li",c:"ms-5 p-0 list-group-flush",t:"2. title: The title of the modal (in the header).[Default=ALERT!!!]"});
        ccModalMakeOpenParamsUl.b.push({n:"li",c:"ms-5 p-0 list-group-flush",t:"3. buttons: The JSON Markup Language object for the buttons.[Default=[](Empty Array)]"});
        ccModalMakeOpenParamsUl.b.push({n:"li",c:"ms-5 p-0 list-group-flush",t:"&Tab;&Tab; ex. [{n:\"button\",type:\"button\",c:\"btn btn-secondary\",i:\"ccModalFooterBtnClose\",\"data-bs-dismiss\":\"modal\",t:\"Close\",e: ccm.id + \".closeForConfig();\"}]."});
        ccModalMakeOpenParamsUl.b.push({n:"li",c:"ms-5 p-0 list-group-flush",t:"4. footer: The text in the footer that will be to the left of the buttons. "});
        ccModalMakeOpenParamsUl.b.push({n:"p",c:"ms-3 p-0 mt-2 mb-0 fst-italic fw-bold",t:"Example Listener:"});
        ccModalMakeOpenParamsUl.b.push({n:"p",c:"ms-4 p-0 m-0 mb-0 fst-italic",t:"document.querySelector(\"#someElement\").addEventListener(\"click\",()=>{"});
        ccModalMakeOpenParamsUl.b.push({n:"p",c:"ms-4 p-0 m-0 fst-italic",t:"&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;ccModalConfigOpen('Wow!!! A Popup!!',\"CC Modal Rocks!!\","});
        ccModalMakeOpenParamsUl.b.push({n:"p",c:"ms-4 p-0 m-0 fst-italic",t:"&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;[{\"n\": \"button\",\"i\": \"testButton\",\"c\": \"btn btn-primary\",\"t\": \"Test\",\"e\": \"testFunc();\" }],"});
        ccModalMakeOpenParamsUl.b.push({n:"p",c:"ms-4 p-0 m-0 fst-italic",t:"&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;\"OW! My foot err\");"});
        docCardBody.b.push(ccModalMakeOpenParamsUl);
 
        docCardBody.b.push({n:"span",c:"m-2 p-0 fw-bold fs-4 d-block",t:"ccModalConfigOpen(cfg) -- opens the modal and allows closing based on config"});
        let ccModalConfigOpenParamsUl={n:"ul",c:"ms-5 p-0 list-group list-group-flush list-unstyled",t:"---Config Options:",b:[]};
        ccModalConfigOpenParamsUl.b.push({n:"li",c:"ms-5 p-0 list-group-flush",t:"cfg.id: The element id of the modal.[Default=ccModalPopup]"});
        ccModalConfigOpenParamsUl.b.push({n:"li",c:"ms-5 p-0 list-group-flush",t:"cfg.titleText: The title of the modal (in the header).[Default=\"\"]"});
        ccModalConfigOpenParamsUl.b.push({n:"li",c:"ms-5 p-0 list-group-flush",t:"cfg.messageText: The main message of the modal (in the body).[Default=ALERT!!!]"});
        ccModalConfigOpenParamsUl.b.push({n:"li",c:"ms-5 p-0 list-group-flush",t:"cfg.footText: The text in the footer that will be to the left of the buttons..[Default=\"\"]"});
        ccModalConfigOpenParamsUl.b.push({n:"li",c:"ms-5 p-0 list-group-flush",t:"cfg.closeOnBackgroundClick: Determines if a background click closes the modal. [Default=true]"});
        ccModalConfigOpenParamsUl.b.push({n:"li",c:"ms-5 p-0 list-group-flush",t:"cfg.hasCloseXButton: Determines if the title has an X close button.[Default=false]"});
        ccModalConfigOpenParamsUl.b.push({n:"li",c:"ms-5 p-0 list-group-flush",t:"cfg.hasOverlay: Determines if the modal has a overlay/backdrop. [Default=true]"});
        ccModalConfigOpenParamsUl.b.push({n:"li",c:"ms-5 p-0 list-group-flush",t:"cfg.buttons: The JSON Markup Language object for the buttons.[Default=[](Empty Array)]"});
        ccModalConfigOpenParamsUl.b.push({n:"li",c:"ms-5 p-0 list-group-flush",t:"&Tab;&Tab; ex JML. [{n:\"button\",type:\"button\",c:\"btn btn-secondary\",i:\"ccModalFooterBtnClose\",\"data-bs-dismiss\":\"modal\",t:\"Close\",e: ccm.id + \".closeForConfig();\"}]."});
        ccModalConfigOpenParamsUl.b.push({n:"li",c:"ms-5 p-0 list-group-flush",t:"&Tab;&Tab;NOTE: A Close Button will be added unless you set cfg.addCloseButton to false"});
        ccModalConfigOpenParamsUl.b.push({n:"li",c:"ms-5 p-0 list-group-flush",t:"cfg.headerJML: If you wish to create you own custom header. [Default={i: ccm.id + \"ModalHeader\",c:\"modal-header\",b:[]}]"});
        ccModalConfigOpenParamsUl.b.push({n:"li",c:"ms-5 p-0 list-group-flush",t:"cfg.titleJML: If you wish to create you own custom title div. [Default={i: ccm.id + \"ModalTitle\",c:\"h3 m-auto\", t:ccm.titleText};]"});
        ccModalConfigOpenParamsUl.b.push({n:"li",c:"ms-5 p-0 list-group-flush",t:"cfg.titleButtonJML: If you wish to create you own custom title button. [Default={i:ccm.id + \"ModalTitleXCloseBtn\",type:\"button\",c:\"btn-close\",\"data-bs-dismiss\":\"modal\",\"aria-label\": \"Close\"}]"});
        ccModalConfigOpenParamsUl.b.push({n:"li",c:"ms-5 p-0 list-group-flush",t:"cfg.bodyJML: If you wish to create you own custom modal body. [Default={i: ccm.id + \"ModalBody\",c: \"modal-body\",b:[]}]"});
        ccModalConfigOpenParamsUl.b.push({n:"li",c:"ms-5 p-0 list-group-flush",t:"cfg.modalJML: If you wish to create the generated cc-modal. [Default={n: \"cc-modal\",i: ccm.id,b:[]}]"});
        ccModalConfigOpenParamsUl.b.push({n:"li",c:"ms-5 p-0 list-group-flush",t:"cfg.dialogJML: If you wish to create you own custom modal dialog. [Default={i: ccm.id + \"ModalDialog\",c:\"modal-cc-dialog\",b:[]};]"});
        ccModalConfigOpenParamsUl.b.push({n:"li",c:"ms-5 p-0 list-group-flush",t:"cfg.contentJML: If you wish to create you own custom modal content. [Default={i: ccm.id + \"ModalContent\",c:\"modal-content\",b:[]}]"});
        ccModalConfigOpenParamsUl.b.push({n:"li",c:"ms-5 p-0 list-group-flush",t:"cfg.footerJML: If you wish to create you own custom footer. [Default={i: ccm.id + \"ModalFooter\",c: \"modal-footer\",b:[]}]"});
        ccModalConfigOpenParamsUl.b.push({n:"li",c:"ms-5 p-0 list-group-flush",t:"cfg.closeButtonJML: The main message of the modal (in the body).[Default=Sorry! No reason for this alert!]"});
        ccModalConfigOpenParamsUl.b.push({n:"p",c:"ms-3 p-0 mt-2 mb-0 fst-italic fw-bold",t:"Example Listener:"});
        ccModalConfigOpenParamsUl.b.push({n:"p",c:"ms-4 p-0 m-0 mb-0 fst-italic",t:"document.querySelector(\"#someElement\").addEventListener(\"click\",()=>{"});
        ccModalConfigOpenParamsUl.b.push({n:"p",c:"ms-4 p-0 m-0 fst-italic",t:"&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;const cfg={};"});
        ccModalConfigOpenParamsUl.b.push({n:"p",c:"ms-4 p-0 m-0 fst-italic",t:"&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;cfg.titleText=\"ccModal is out of this world!\";"});
        ccModalConfigOpenParamsUl.b.push({n:"p",c:"ms-4 p-0 m-0 fst-italic",t:"&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;cfg.messageText=\"This Was made by a full config\";"});
        ccModalConfigOpenParamsUl.b.push({n:"p",c:"ms-4 p-0 m-0 fst-italic",t:"&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;cfg.footText=\"Limitless Config!\";"});
        ccModalConfigOpenParamsUl.b.push({n:"p",c:"ms-4 p-0 m-0 fst-italic",t:"&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;//I want a scary red header!!!;"});
        ccModalConfigOpenParamsUl.b.push({n:"p",c:"ms-4 p-0 m-0 fst-italic",t:"&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;cfg.headerJML={i: \"FullConfigccModal\" + \"ModalHeader\",c:\"modal-header bg-danger\",b:[]};"});
        ccModalConfigOpenParamsUl.b.push({n:"p",c:"ms-4 p-0 m-0 fst-italic",t:"&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;&Tab;ccModalConfigOpen(cfg);"});
        docCardBody.b.push(ccModalConfigOpenParamsUl);
 
        //docCardBody.b.push({n:"pre",c:"m-0 p-0",t:" Stub Function(s): "});
        docCard.b.push(docCardBody);
        return docCard;
      }
   



    //--------------------------------------------------------Methods END
 
    //--------------------------------------------------------Properties
 
    /**
     * This is a JML for the modal main close button in the footer
     *
     * @memberof ccModal
     */
    closeButton(){
        let ccm=this;    
        return ccm._closeButtonJML||{n:"button",type:"button",c:"btn btn-secondary",i:"ccModalFooterBtnClose","data-bs-dismiss":"modal",t:"Close",e: ccm.id + ".close();"};;
    }


   /**
     * titleText property is what will be presented in the Title of the Modal
     *
     * @memberof ccModal
     */
    get titleText() { return this._titleText; }
    set titleText(value) {this._titleText = value;}
    /**
     * messageText property is what will be presented in the message of the Modal
     *
     * @memberof ccModal
     */
    get messageText() { return this._messageText; }
    set messageText(value) {this._messageText = value;}
    /**
     * Footer Text property is what will be presented in the footer of the Modal
     *
     * @memberof ccModal
     */
    get footText() { return this._footText; }
    set footText(value) {this._footText = value;}
    /**
     * This stores buttons as html(BS5) buttons as JML to be added to the modal on open.
     * buttons is an Array so add buttons with push()
     *
     * @memberof ccModal
     */
    get buttons() { return this._buttons; }
    set buttons(value) {this._buttons = value;}
    /**
     * This stores buttons as html(BS5) buttons as JML to be added to the modal on open.
     * buttons is an Array so add buttons with push()
     *
     * @memberof ccModal
     */
    get closeOnBackgroundClick() { return this._closeOnBackgroundClick; }
    set closeOnBackgroundClick(value) {this._closeOnBackgroundClick = value;}    
    /**
     * This determines if the X in the header is there for closing the Modal
     *
     * @memberof ccModal
     */
    get hasCloseXButton() { return this._hasCloseXButton; }
    set hasCloseXButton(value) {this._hasCloseXButton = value;}    
    /**
     * This is a ref to the Modal Element
     *
     * @memberof ccModal
     */
    get modalEle() { return this.modalElement; }
    set modalEle(value) {this.modalElement = value;}    
    /**
     * This is a ref to the Modal Element
     *
     * @memberof ccModal
     */
    get overlayElement() { return this._overlayEle; }
    set overlayElement(value) {this._overlayEle = value;}    
   
    //JML Props
    /**
     * This is a JML for the modal header element .modal-header
     *
     * @memberof ccModal
     */
    get headerJML() { return this._headerJML; }
    set headerJML(value) {this._headerJML = value;}    
    /**
     * This is a JML for the title
     *
     * @memberof ccModal
     */
    get titleJML() { return this._titleJML; }
    set titleJML(value) {this._titleJML = value;}    
    /**
     * This is a JML for the Title Button for close
     *
     * @memberof ccModal
     */
    get titleButtonJML() { return this._titleButtonJML; }
    set titleButtonJML(value) {this._titleButtonJML = value;}    
    /**
     * This is a JML for the body element where text about this popup goes .modal.body
     *
     * @memberof ccModal
     */
    get bodyJML() { return this._bodyJML; }
    set bodyJML(value) {this._bodyJML = value;}    
    /**
     * This is a JML for the modal (the bs5 modal outer element .modal)
     *
     * @memberof ccModal
     */
    get modalJML() { return this._modalJML; }
    set modalJML(value) {this._modalJML = value;}    
    /**
     * This is a JML for the dialog (first interior element in the modal .modal-cc-dialog)
     *
     * @memberof ccModal
     */
    get dialogJML() { return this._dialogJML; }
    set dialogJML(value) {this._dialogJML = value;}    
    /**
     * This is a JML for the content (Header, body, and footer) .modal-content
     *
     * @memberof ccModal
     */
    get contentJML() { return this._contentJML; }
    set contentJML(value) {this._contentJML = value;}    
    /**
     * This is a JML for the modal footer .modal-footer
     *
     * @memberof ccModal
     */
    get footerJML() { return this._footerJML; }
    set footerJML(value) {this._footerJML = value;}    
 
    //--------------------------------------------------------Properties END  
 
}
 
document.addEventListener("DOMContentLoaded", ()=> {
    window._ = document;
 
});
 
window.ccModalConfigOpen=async (cfg)=>{
    let mdl=new ccModal();
    if(cfg.id===undefined) cfg.id="ccModalGlobal";
    await mdl.openFromConfig(cfg);
    return mdl;
};
/*
* This is a global function that creates and opens a modal
* Given Message, Title, Button(s), and footer text.  
* You can just give the message.  Everything else defaults.
* title defaults to ALERT!!!
* buttons defaults to [] if you set it to null specifically (no button will appear. weird modal)
*                        if you set it to an array of JML buttons it will do that
*                        if you leave it as undefined, "", etc it will use the default button
* footer defaults to "" if you set it to null specifically (no footer will appear)  
*
*/
window.ccModalMakeOpen=async (message, title,buttons,footer)=>{
    let cfg={n: "cc-modal",i: "ccModal"};
    cfg.id="ccModalGlobal";
    cfg.messageText=message||"Sorry! No reason for this alert!";
    cfg.titleText=title||"ALERT!!!";
    cfg.buttons=buttons;
    cfg.footText=footer||"";
    let mdl=new ccModal();
    await mdl.openFromConfig(cfg);
    return mdl;
}
//inline is best or it can be missed if in the DOMContentLoaded causing Illegal Constructor Error
window.customElements.define("cc-modal", ccModal);
 
// document.addEventListener("DOMContentLoaded", function () {
//     customElements.define("cc-modal", ccModal);
//   });
 
function ccModalDefine(){
    customElements.define("cc-modal", ccModal);
}
 
export default ccModal;
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
//!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! S.D.G !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
//^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^