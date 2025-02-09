//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! J.J. !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!

//^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

 

/**

* Replacement for Javascript Class (with prototype functions) utilities.js

* This is the modular version.  utilities.js is the monolithic version and is not recommended for use.

* Public Domain

* Licensed Copyright Law of the United States of America, Section 105 (https://www.copyright.gov/title17/92chap1.html#105)

* Per hoc, facies, scietis quod ille miserit me ut facerem universa quae cernitis et factis: Non est mecum!

* Published by: Dominic Roche of OIT/IUSG/DASM on 03/13/2024

* @class ccReporting

*/

 

"use strict";

import ccModal from './ccComponents/ccModal.js';
import ccDOMUtil from './ccDOMUtil.js';
//Export Values
// for all DOM
//      import { addJsScript, append, classAdd, classRemove, css, elementHide, elementRemove,elementShow, getEle, getEles, getVal, isCheck, isElement, makeAlert, makeAlertJML,  off, on, placeAfter, placeBefore, unbind} from '/Scripts/Views/Shared/ccComponents/ccUtilities.js';
//
// for all code
//      import {asAlphaNumeric, asBool, asElementId, asFieldNotationString, asHtml5DateFormat, asInt, asPascalCase, asPropertyNotation, asString, camelToTitle, checkSetDefault, chunkString, clip, disableScroll, enableScroll, endsWith, errorLog, fillFrom, findAllIndexesOf, formClear, formEnabled, formFill, formToJson, functionCall, getIndicesOf, getType, guid, htmlDecode, htmlEncode, htmlEncodedStrip, htmlToJML, imageFix, isEmpty, isHtml, isJson, isNode, isNotEmpty, isNumeric, jsonToHtml, logError, objectKeysToCamelCase, placeAfter, placeBefore, repeatGuidFunc, reportErrorObject, startsWith, stringifyNumber, toTitle, trim, messageBox, makeICS, setICSFunc, adjustPopupToKey, functionCallAsync, get, openUrlClean, set, setWhereWasI, simpleCall, simpleCallQuiet, takeMeBack, timeOfDay, textBetween, toString, whereWasI, commsResponseStatusHandle, delay, simplePost} from '/Scripts/Views/Shared/ccComponents/ccUtilities.js';
//
//
//      --Export Properties: whereWasI
//      --DOM related: append, classAdd, classRemove, css, elementHide, elementRemove,
//              elementShow, getEle, getEles, getVal, isCheck, isElement, makeAlert, makeAlertJML, off, on,
//              placeAfter, placeBefore, unbind
//      --Code related: asAlphaNumeric, asBool, asElementId, asFieldNotationString,
//                      asHtml5DateFormat, asInt, asPascalCase, asPropertyNotation, asString,
//                      camelToTitle, checkSetDefault, chunkString, disableScroll, enableScroll,
//                      endsWith, errorLog, fillFrom, formClear, formEnabled, formFill,
//                      formToJson, functionCall, getIndicesOf, getType, guid, simpleCall, simpleCallQuiet,
//                      htmlDecode, htmlEncode, htmlEncodedStrip, htmlToJML, imageFix, 
//                      isEmpty, isHtml, isJson, isNode, isNotEmpty, isNumeric, jsonToHtml,
//                      logError, objectKeysToCamelCase, repeatGuidFunc, reportErrorObject,
//                      startsWith, stringifyNumber, textBetween, toTitle, trim, unobtrusiveWait, unobtrusiveWaitOff
//      --External Functions: messageBox
//      --Create Files Functions: makeICS, setICSFunc
//      --Page Functions: adjustPopupToKey, functionCallAsync, get, openUrlClean, set,
//                          setWhereWasI, takeMeBack, timeOfDay, toString, whereWasI
//      --Comms Function: commsResponseStatusHandle,dispatch, delay, simplePost
 
// Not for Export
//      --Support Functions: classAddForMulti, classRemoveForMulti, guidS4, repeatGuidFunc
 
 
//------------------------------------------------------Export Properties
const controller = new AbortController();
export let whereWasI = null;
export let quietComms = false;
export let trackId =0;
export const signal = controller.signal;
//------------------------------------------------------Export Functions
//                  ~~~~~~~~~~~~~~~~~~~~ DOM Functions ~~~~~~~~~~~~~~~~~~~~~                   \\
 
 
/**
* Adds a CSS file to the page
*
* @export
* @param {*} cssId
* @param {*} url
* @param {*} type
* @param {*} isAsync
* @memberof ccUtilities
*/
export function addCSSLink(cssId,url) {
    if (!document.getElementById(cssId)) {
        let head = document.getElementsByTagName("head")[0];
        let link = document.createElement("link");
        link.id = cssId;
        link.rel = "stylesheet";
        link.type = "text/css";
        link.href = url;
        link.media = "all";
        head.appendChild(link);
    }
}
/**
* Adds a Javascript file to the page
*
* @export
* @param {*} jsId
* @param {*} url
* @param {*} type
* @param {*} isAsync
* @memberof ccUtilities
*/
export function addJsScript(jsId,url,type,isAsync) {
    if(!type) type="text/javascript";
    if (!document.getElementById(jsId)) {
        let head = document.getElementsByTagName("head")[0];
        let link = document.createElement("script");
        link.id = jsId;
        link.type = type;
        link.src =location.origin + url;
        if(isAsync)link.async=""
        head.appendChild(link);
    }
}
/**
* Append the given html to the element specified
*                       in this case as the last child element
* @param {*} element Target Element (The element you wish to place the html or element in or around)
* @param {*} htmlOrElement html or element to add
* @memberof ccUtilities
*/
export function append(element, htmlOrElement) {
    if (typeof htmlOrElement === 'string' || htmlOrElement instanceof String) element.insertAdjacentHTML("beforeend", htmlOrElement);
    if (typeof htmlOrElement === 'object' || htmlOrElement instanceof Object) element.insertAdjacentElement("beforeend", htmlOrElement);
}
/**
* Add a class from an element (Like JQuery .addClass())
*
* @param {*} eles DOM element use
* @param {*} className class name(s) to remove from element (You can add multiple classes via normal class value ie. 'classOne classTwo andSomeOtherClass')
* @returns element so you can chain stuff
* @memberof ccUtilities
*/
export function classAdd(eles, className) {
    if (eles === undefined || eles === null) return eles;
    if (NodeList.prototype.isPrototypeOf(eles)) eles=Array.from(eles);
    if(Array.isArray(eles))
    {
        for (let ele of eles) {
            ele=classAddForMulti(ele,className);
        }
    } else {
        if(className==="") return eles.className="";
        if(eles.classList.length>0 && (className||"").split(' ').some(c=>c.startsWith("d-")) && Array.from(eles.classList).some(c=>c.startsWith("d-"))){
            let clBad = Array.from(eles.classList).filter(c=> c.startsWith("d-"));
            for(const cls of clBad) { eles.classList.remove(cls);}
        }
        classAddForMulti(eles,className);
    }
    return eles;
 
}
/**
* Remove a class from an element (Like JQuery .removeClass())
*
* @param {*} ele DOM element use
* @param {*} className class name to remove from element
* @returns element so you can chain stuff
* @memberof ccUtilities
*/
export function classRemove(ele, className) {
    if (ele === undefined || ele === null) return false;
    if (NodeList.prototype.isPrototypeOf(ele)) {
        let eles = Array.from(ele);
        for (const e of eles) {
            classRemove(e,className);
        }
    } else {
        if (ele.classList) {
            classRemoveForMulti(ele,className);
            //ele.classList.remove(className);
        } else {
            if (ele.classlist !== undefined) {
                classRemoveForMulti(eles,className);
                ele.className = ele.className.replace(new RegExp("(^|\\b)" + className.split(" ").join("|") + "(\\b|$)","gi"),"").trim();
            }
        }
    }
    return ele;
}
/**
* css Get the value of a computed style property for the first element in the set of matched elements
* or set one or more CSS properties for every matched element. Just like jQuery css() command
* @param {*} ele The element(s) (not jQuery object)
* @param {*} attr Attribute to change
* @param {*} val value to change the attribute to
* @returns original element for chaining
* @memberof ccUtilities
*/
export function css(ele, attr, val) {
    if (attr === null) {
        console.log("Utilities css failed to alter " + ele.utlAttr("id") + " due to null attr parameter!");
        return ele;
    }
    if (typeof attr === "object") {
        //Set only with object
        Object.entries(attr).forEach(function([key, value]) {
            ele.style[key] = null;
            ele.style[key] = value;
        });
        return ele;
    } else {
        if (val === undefined) {
            return ele.style;
        } else {
            ele.style[attr] = null;
            ele.style[attr] = val;
        }
        return ele;
    }
}
/**
* Hides Element(s)
*
* @param {*} eles
* @returns element(s) given
* @memberof ccUtilities
*/
export function elementHide(eles){
    if (eles === undefined || eles === null) return eles;
    if (NodeList.prototype.isPrototypeOf(eles)) {
        eles = Array.from(eles);
    } else {
        eles = [eles];
    }
    for (const ele of eles) {
        //ele.style.display="none";
        ele.utlClassAdd("d-none");
        if(ele.classList.contains("panel-collapse")){
            ele.utlClassAdd("collapse");
        }else{
            ele.utlClassAdd("d-none");
        }
 
    }
    return eles;
}              
/**
* Removes element from the dom
*
* @param {*} ele
* @returns
* @memberof ccUtilities
*/
export function elementRemove(eles) {
    if (eles === undefined || eles === null) return eles;
    if (NodeList.prototype.isPrototypeOf(eles)) {
        eles = Array.from(eles);
    } else {
        eles = [eles];
    }
    for (const ele of eles) {
        ele.remove();
    }
    return eles;
}
/**
* Shows Element(s)
*
* @param {*} eles element(s) to effect by this.  Note that this removes bootstrap 4 d-none as well as style.display="none"
* @returns element(s) given
* @memberof ccUtilities
*/
export function elementShow(eles){
    if (eles === undefined || eles === null) return eles;
    if (NodeList.prototype.isPrototypeOf(eles)) {
        eles = Array.from(eles);
    } else {
        eles = [eles];
    }
    for (const ele of eles) {
        ele.style.display="";
        ele.classList.remove('d-none');
        ele.classList.remove('collapse');
    }
    return eles;
}
/**
* Get the first element that matches the query (must be a standard querySelectorAll selectors type)
* replacement for jQuery $() but returns only 1 record
*
* @param {*} src source element to search from (defaults to document)
* if src is a query and q is undefined/null this will search for src query in the document
 * @param {*} q standard querySelectorAll selectors query
* @returns the first DOM Element matching the criteria
* @memberof ccUtilities
*/
export function getEle(src, q) {
    if(!q){
        if(!src) return document.body;
        q=src;
        src=document;
    }
    let rtn = getEles(src, q)[0];
    if (rtn === undefined) {
        //console.log("getEle Error:" + q + " not found in document!");
        rtn = document.createElement("div",{id:"ElementNotFound",length: 0});
        //rtn.length=0;
        return rtn;
    }
    return rtn;
}
/**
* Get the elements that matches the query (must be a standard querySelectorAll selectors type)
* replacement for jQuery $()
*
* @param {*} src source element to search from (defaults to document)
* @param {*} q standard querySelectorAll selectors query
* @returns all DOM Elements matching the criteria
* @memberof ccUtilities
*/
export function getEles(src, q) {
    if(!q){
        if(!src) return document.body;
        q=src;
        src=document;
    }
    if (q === undefined) return document;
    q = q.replaceAll(":input", "input, select, checkbox, textarea");
    return src.querySelectorAll(q);
}
/**
* Gets the value of an :inputs including checkbox, select, textArea and standard inputs.  If the element is NOT an :input
* type this will return the innerHtml of the element so you still get a value if you use this on a div for instance
* Likewise if a val is provided this will set the :input to the value (regardless of type) and if its some other element like div
* it will set that html
*
* @param {*} ele DOM element to use
 * @param {*} val do not provide for get but for set this is what to set the element value/html to
* @returns element given for chaining
* @memberof ccUtilities
*/
export function getVal(ele, val) {
    if (val === undefined) {
        switch (ele.nodeName) {
            case "SELECT": {
                if (ele.hasAttribute("multiple")) {
                    return Array.from(ele.selectedOptions).map(o => o.label).join(',');
                } else {
                    if (ele.selectedIndex < 0) return "";
                    return ele[ele.selectedIndex].value;
                }
            }
        case "TEXTAREA":
            return ele.value;
        case "INPUT":
            if ((ele.type || "").toUpperCase() === "CHECKBOX") {
                return ele.checked;
            } else {
                return ele.value;
            }
        default:
            return ele.innerHTML;
        }
    } else {
        switch (ele.nodeName) {
        case "SELECT":
            if(val===null || val===""){
                ele.value="";
            }else{
                [...ele.options].some((option, index) => {
                    if (option.value == val) {
                        ele.selectedIndex = index;
                    }
                });
            }
            break;
        case "DIV":
            ele.innerHTML = val;
            break;
        default:
            if ((ele.type || "").toUpperCase() === "CHECKBOX") {
                if (
                    val === true ||
            (val || "").toLowerCase() === "true" ||
            (val || "").toLowerCase() === "on"
                ) {
                    ele.checked = true;
                } else {
                    ele.checked = false;
                }
            } else {
                ele.value = val;
            }
            break;
        }
        return ele;
    }
}
/**
* This (like jQuery .is() method) attempts to derive intelligence about an item.
* currently it only works with check types :checked,:visible, or assumes a caseless comparison of the nodeName ro the type given
*
* @param {*} ele DOM element to use
 * @param {*} checkType check tyoe (:checked,:visible, or some node type like DIV)
* @returns
* @memberof ccUtilities
*/
export function isCheck(ele, checkType) {
    switch (checkType) {
    case ":checked":
        return ele.checked;
    case ":visible":
        if (NodeList.prototype.isPrototypeOf(ele)) ele = Array.from(ele);
        if (Array.isArray(ele)) {
            return ele.some(e => e.style.display !== "none" && e.className.indexOf("d-none")<0);
        } else {
            return ele.style.display !== "none" && ele.className.indexOf("d-none")<0 && ele.className.indexOf("collapse")<0;
        }
    default:
        return ele.nodeName.toUpperCase() === checkType.toUpperCase();
    }
}
/**
* this allows removing an event on an element(s)
*
* @param {*} eles element or elements (NodeList or Array) to remove event tfrom
* @param {*} event event type to attach to (abort,click,contextmenu,focus,input, etc)
* @param {*} func anonymous function to call
* @returns true if successful
 * @memberof ccUtilities
 */
export function off(eles, event, func) {
    if (eles === undefined || eles === null) return false;
    if (NodeList.prototype.isPrototypeOf(eles)) {
        eles = Array.from(eles);
    } else {
        eles = [eles];
    }
    for (const ele of eles) {
        ele.removeEventListener(event, func, false);
    }
    return true;
}              
/**
* this allows setting an event on an element(s)
*
* @param {*} eles element or elements (NodeList or Array) to attach event to
* @param {*} event event type to attach to (abort,click,contextmenu,focus,input, etc)
* @param {*} func anonymous function to call
* @returns true if successful
 * @memberof ccUtilities
 */
export function on(eles, event, func) {
    if (eles === undefined || eles === null) return false;
    if (NodeList.prototype.isPrototypeOf(eles)) {
        eles = Array.from(eles);
    } else {
        if(Array.isArray(eles)){
            //Array of Nodelist
        }else{
            eles = [eles];
        }
    }
    for (const ele of eles) {
        ele.addEventListener(event, func, false);
    }
    return true;
}
/**
* prepend the given html to the element specified
*
* @param {*} element Target Element (The element you wish to place the html or element in or around)
*                       in this case in as first child
* @param {*} html html or element to add
* @memberof ccUtilities
*/
export function prepend(element, htmlOrElement) {
    if (typeof htmlOrElement === 'string' || htmlOrElement instanceof String) element.insertAdjacentHTML("afterbegin", htmlOrElement);
    if (typeof htmlOrElement === 'object' || htmlOrElement instanceof Object) element.insertAdjacentElement("afterbegin", htmlOrElement);
}
 
/**
* placeAfter the given html to the element specified
*                       in this case after the element
* @param {*} element Target Element (The element you wish to place the html or element in or around)
* @param {*} htmlOrElement html or element to add
* @memberof ccUtilities
*/
export function placeAfter(element, htmlOrElement) {
    if (typeof htmlOrElement === 'string' || htmlOrElement instanceof String) element.insertAdjacentHTML("afterend", htmlOrElement);
    if (typeof htmlOrElement === 'object' || htmlOrElement instanceof Object) element.insertAdjacentElement("afterend", htmlOrElement);
}
/**
* placeBefore the given html to the element specified
*                       in this case before the element
* @param {*} element Target Element (The element you wish to place the html or element in or around)
* @param {*} htmlOrElement html or element to add
* @memberof ccUtilities
*/
export function placeBefore(element, htmlOrElement) {
    if (typeof htmlOrElement === 'string' || htmlOrElement instanceof String) element.insertAdjacentHTML("beforebegin", htmlOrElement);
    if (typeof htmlOrElement === 'object' || htmlOrElement instanceof Object) element.insertAdjacentElement("beforebegin", htmlOrElement);
}
 
 
 
/**
* Unbinds the events from element(s) given
*
* @param {*} ele
* @memberof ccUtilities
*/
export function unbind(ele) {
    if (NodeList.prototype.isPrototypeOf(ele)) {
        ele = Array.from(ele);
    } else {
        ele = [ele];
    }
    for (const e of ele) {
        let old_ele = e;
        let new_ele = e.cloneNode(true);
        //new_ele.children =old_ele.children;
        old_ele.parentNode.replaceChild(new_ele, old_ele);
    }
}
 
//                  ~~~~~~~~~~~~~~~~~~~~ DOM Functions END~~~~~~~~~~~~~~~~~~                   \\
 
//                  ~~~~~~~~~~~~~~~~~~~~ Code Functions ~~~~~~~~~~~~~~~~~~~~                   \\
 
/**
* Returns if any items meet the fn
*
* @export
* @param {*} fn
* @return {*}
 * @memberof ccUtilities
*/
export function any(fn) {
    if (this === undefined || this === null) return false;
    if (!Array.isArray(this)) return false;
    if (this.length === 0) return false;
    if (fn === undefined) fn =()=>true;
    return this.some(fn);
}
/**
* Removes all special characters from a string
*
* @param {*} str to convert
* @returns clean alpha/numeric string
* @memberof ccUtilities
*/
export function asAlphaNumeric(str) {
    if (!str) return "";
    return str.replace(/[^\w]/gi, '');
}
/**
* Converts the object given to a bool
*
* @param {*} obj object to convert
* @returns true if true (testing for type of true) or false
* @memberof ccUtilities
*/
export function asBool(obj) {
    if (!obj) return false;
    if (obj === true) return true;
    if (Array.isArray(obj)) return false;
    let objString=trim(obj.toString().toLowerCase());
    if (objString === "true") return true;
    if (objString === "1")  return true;
    if (objString === "y") return true;
    return false;
}
/**
* Just makes sure that a # is prepended to a string when using it for querySelectorAll
*
* @param {*} id name of element to assure has a prefix of #
* @returns correct querySelectorAll id
* @memberof ccUtilities
*/
export function asElementId(id) {
    if (!id) return "";
    id=startsWith(id, "#")?id:"#"+id;
    id=id.replaceAll("##", "#");
    return asString(id);
}
/**
* converts string to camel case
 *
* @param {*} str
* @returns given string in fieldNotation (camel notation)
* NOTE:  if you give it a camel notation string (ie.  somePropertyName) it will prefix it with an
 * underscore. It assumes the given string (since already lowercase start) is a improper property name
* and thus uses the improper field name _somePropertyName.
* @memberof ccUtilities
*/
export function asFieldNotationString(str) {
    if (!str) return "";
    if (str.length < 1) return "";
    if(str.indexOf(" ")>-1) str=str.split(" ").map(s=>s.substring(0, 1).toLowerCase() + s.substring(1)).join("");
    str=replaceAll(str," ","");
    if (str.length === 1) return str.toLowerCase();
    if (str.substring(0, 1) === str.substring(0, 1).toLowerCase()) return "_" + str;
    return str.substring(0, 1).toLowerCase() + str.substring(1);
}
/**
* Given a Javascript Date() this translates it to a text of html5 date (for use with input type date)
*
* @param {*} d given date
* @returns string date formatted fo type input
* @memberof ccUtilities
*/
export function asHtml5DateFormat(d){
    let month = '' + (d.getMonth() + 1);
    let day = '' + d.getDate();
    let year = d.getFullYear();
    if (month.length < 2) month = '0' + month;
    if (day.length < 2) day = '0' + day;
    return [year, month, day].join('-');
}
 
 
export function asDateParts(dateIn){
    const date = new Date(dateIn);
    const day = date.toLocaleString('en-US', { day: '2-digit' });
    const month = date.toLocaleString('en-US', { month: 'short' });
    const year = date.toLocaleString('en-US', { year: 'numeric' });
    const hours = date.toLocaleString('en-US', { hour: '2-digit', hour12: false }).padStart(2, '0');
    const minutes = date.toLocaleString('en-US', { minute: '2-digit' }).padStart(2, '0');
    const seconds = date.toLocaleString('en-US', { second: '2-digit' }).padStart(2, '0');
    return { day, month, year, hours, minutes, seconds };
}
 
/**
* Converts the given object to an integer
*
* @param {*} obj object to convert to int
* @returns integer representation of the object or -1 if its not numeric
* @memberof ccUtilities
*/
export function asInt(obj) {
    if (!isNumeric(obj)) return -1;
    return parseInt(obj);
}
/**
* Converts given string to pascal case
*
* @param {*} str
* @returns
* @memberof ccUtilities
*/
export function asPascalCase(str) {
    if (str === undefined || str === null) return "";
    if (str.length < 1) return "";
    if (str.length === 1) return str.toUpperCase();
    str = str.substring(0, 1).toUpperCase() + str.substring(1).toLowerCase();
    return str.replace(/\w\S*/g, function(txt) {
        return txt.charAt(0).toUpperCase() + txt.substr(1).toLowerCase();
    });
}
/**
* Converts string to Pascal notation
 *
* @param {*} str string given for converting
* @returns the original string in Pascal notation
* @memberof ccUtilities
*/
export function asPropertyNotation(str) {
    if (!str) return "";
    if (str.length < 1) return "";
    if(str.indexOf(" ")>-1) str=str.split(" ").map(s=>s.substring(0, 1).toUpperCase() + s.substring(1)).join("");
    if (str.length === 1) return str.toUpperCase();
    return str.substring(0, 1).toUpperCase() + str.substring(1);
}
 
/**
* Changes an object to a string.  This can be used with numbers as well as to ensure that an empty string occurs for nulls or undefined
*
* @param {*} obj object to convert
* @param {*} size used for radix in number type
* @returns
* @memberof ccUtilities
*/
export function asString(obj, size) {
    if (!obj) return "";
    return (size !== undefined) ? obj.toString(size) : obj.toString();
}
/**
* Breaks a camel/pascal notation string to a title.  So SomeCoolTitle become Some Cool Title
*
* @param {*} str string to convert
* @returns Title string
* @memberof ccUtilities
*/
export function camelToTitle(str) {
    if(!str) return "";
    return str.replace(/([A-Z])/g, " $1").replace(/^./, function(str) {             return str.toUpperCase();}).trim();
}
/**
* Given an object and a default this sets the property (after creating it)
*
* @param {*} obj Object to work on
* @param {*} props Comma Delimited list of Properties to create and set default value
* @param {*} setTo default value
* @returns
* @memberof ccUtilities
*/
export function checkSetDefault(obj,props,setTo){
    if(!obj || !props || !setTo) return null;
    if(!Array.isArray(props)) props = props.split(',');
    for(const p of props){
        if(Array.isArray(setTo) && setTo.length===0){
            if(!obj[p]) obj[p]=new Array();
        }else{
            if(!obj[p]) obj[p]=setTo;
        }
    }
    return obj;
}
/**
* Chunk a string into chunks
*  courtesy of: Vivin Paliath ( https://stackoverflow.com/users/263004/vivin-paliath )
*
* @param {*} str given string
* @param {*} length length of chunk
* @returns
* @memberof ccUtilities
*/
export function chunkString(str, length) {
    return str.match(new RegExp('.{1,' + length + '}', 'g'));
}
 
/**
* Clips a string to the maxLength
*
* @param {*} text
* @param {*} maxLength
* @return {*}
 */
function clip(text, maxLength) {
    if (!text || !text.trim()) return "";
    if (maxLength < 1) return text;
    return text.length > maxLength ? text.slice(0, maxLength) : text;
  }
/**
* used for delaying action events.  This allows for good searching and defaults to 1 sec
*
* @param {*} ms miliseconds to wait
* @returns promise for delay
* @memberof ccUtilities
*/
export async function delay(ms){
    return new Promise(res => setTimeout(res, ms));
}
/**
* Disables Scrolling
* @memberof ccUtilities
*/
export function disableScroll() {
    // Get the current page scroll position
    var scrollTop = window.pageYOffset;// || document.documentElement.scrollTop;
    var scrollLeft = window.pageXOffset;// || document.documentElement.scrollLeft,
 
        // if any scroll is attempted, set this to the previous value
        window.onscroll = function() {
            window.scrollTo(scrollLeft, scrollTop);
        };
}
/**
* Disables Scrolling
* @memberof ccUtilities
*/
export function enableScroll() {
    window.onscroll = function() {};
}
/**
* determines if a string ends with the given search string
*
* @param {*} str string to search in
* @param {*} check string to check for
* @param {*} ignore ignore case sensitivity
* @returns true if the string ends with the check
* @memberof ccUtilities
*/
export function endsWith(str, check, ignore) {
    if (check === undefined || check === null) return false;
    if (str === undefined || str === null) return false;
    if (typeof ignore === "undefined") ignore = true;
    str = "" + str + "";
    check = "" + check + "";
    if (ignore) str.slice(-check.length) == check;
    return str.slice(-check.length) === check;
}
/**
* Fills one object from another
*
* @param {*} base
* @param {*} fill
* @memberof ccUtilities
*/
export function fillFrom(base,fill){
    let fillPropsLwrs = Object.getOwnPropertyNames(fill).map(f=> ((f.charAt(0)===f.charAt(0).toUpperCase())?f.utlAsFieldNotationString():f));
    for(let prop of fillPropsLwrs){
        base[prop]=fill[prop]||fill[prop.utlAsPropertyNotation()];
    }
    return base;
}
 
/**
* Find all indexes
 *
* @export
* @param {*} text
* @param {*} search
* @return {*}
 */
export function findAllIndexesOf(text, search) {
    const indexes = [];
 
    if (text && search) {
      let index = text.indexOf(search);
 
      while (index !== -1) {
        indexes.push(index);
        index = text.indexOf(search, index + 1);
      }
    }
 
    return indexes;
  }
 
/**
* Gets all indices of a search string in a given string
*
* @param {*} searchIn
* @param {*} searchFor
* @param {*} caseSensitive
* @return {*}
 * @memberof ccUtilities
*/
export function getIndicesOf(searchIn, searchFor, caseSensitive) {
    var searchForLen = searchFor.length;
    if (searchForLen == 0) {
        return [];
    }
    var startIndex = 0, index, indices = [];
    if (!caseSensitive) {
        searchIn = searchIn.toLowerCase();
        searchFor = searchFor.toLowerCase();
    }
    while ((index = searchIn.indexOf(searchFor, startIndex)) > -1) {
        indices.push(index);
        startIndex = index + searchForLen;
    }
    return indices;
}
/**
* Gets the type similar to .net style
*
* @param {*} obj object to get type from
* @return {*} type name
* @memberof ccUtilities
*/
export function getType( obj ) {
    let typeOfObject = typeof obj;
    if ( 'object' === typeOfObject ) {
        typeOfObject = Object.prototype.toString.call( obj );
        if ( '[object Object]' === typeOfObject ) {
            if ( obj.constructor.name ) {
                return obj.constructor.name;
            } else if ( '[' === obj.constructor.toString().charAt(0) ) {
                typeOfObject = typeOfObject.substring( 8,typeOfObject.length - 1 );
            } else {
                typeOfObject = obj.constructor.toString().match( /function\s*(\w+)/ );
                if ( typeOfObject ) {
                    return typeOfObject[1];
                } else {
                    return 'Function';
                }
            }
        } else {
            typeOfObject = typeOfObject.substring( 8,typeOfObject.length - 1 );
        }
    }
 
    return typeOfObject.charAt(0).toUpperCase() + typeOfObject.slice(1);
}              
 
/**
* Generates a Random GUID
*
* @export guid
* @return {*}
 * @memberof ccUtilities
*/
export function guid(nodash) {
    if(!nodash) nodash=false;
    //return (repeatGuidFunc(2) + "-" + repeatGuidFunc(4, "-") + "-" + guidS4());
    let sep=nodash?"":"-";
    return (repeatGuidFunc(2) + sep + repeatGuidFunc(4, sep) + guidS4() + guidS4());
}
/**
* Decodes encoded HTML back to HTML
*
* @param {*} value encoded HTML to decode
* @returns HTML
* @memberof ccUtilities
*/
export function htmlDecode(value) {
    if (value == undefined) return "";
    if (value == null) return "";
    if (value == "") return "";
    let val = value
        .split("&amp;apos;")
        .join("'")
        .split("&amp;")
        .join("&")
        .split("&lt;")
        .join("<")
        .split("&gt;")
        .join(">")
        .split("&quot;")
        .join('"');
    return val;
}
/**
* This encodes HTML to render it safe.  So < becomes &lt; and & becomes @amp; etc.
*
* @param {*} value value to encode
* @returns encoded HTML
* @memberof ccUtilities
*/
export function htmlEncode(value) {
    return document.createElement("div").utlText(value).utlHtml();
}
/**
* this decodes HTML and then strips it (the decode step it to ensure that either way you get just the text)
*
* @param {*} value value to strip HTML or encoded HTML from
* @returns Text from within HTML (or endcoded HTML)
* @memberof ccUtilities
*/
export function htmlEncodedStrip(value) {
    if (value == undefined) return "";
    if (value == null) return "";
    if (value == "") return "";
    let val = htmlDecode(value);
    return document.createElement("div").utlHtml(val).utlText();
}
/**
* Converts an existing Element to JHtml (Including Children)
*
* @param {*} htmlEle Element to convert
* @returns JHtml of the element
* @memberof ccUtilities
*/
export function htmlToJML(htmlEle){
    let utl = this;
    let rtn={b:[]};
    //first deal with children
    let children=Array.from(htmlEle.children);
    for(let child of children){
        let cldJhtml=utl.htmlToJML(child);
        rtn.b.push(cldJhtml);
    }
    //now main
    for(let attr of htmlEle.attributes){
        rtn[attr.name]=(attr.value||"");
        let txt=[].reduce.call(htmlEle.childNodes, function(a, b) { return a + (b.nodeType === 3 ? b.textContent : ''); }, '').trim();
        if(txt.length>0) rtn.t=txt;
    }
    if(htmlEle.nodeName!=="DIV") rtn.n=htmlEle.nodeName;
    if(!rtn.i&&!rtn.id) rtn.i=htmlEle.i=(htmlEle.id||(htmlEle.nodeName+guid()));
    return rtn;
}
/**
* removes the '~' for current url style image call
*
* @param {*} image
* @returns
* @memberof ccUtilities
*/
export function imageFix(image) {
    return startsWith(image,"~") ? image.replaceAll("~", "") : image;
}
/**
* determines if the given object (jQuery, node, dom element) is an element
*
* @param {*} obj object to check
* @returns true if the object is a DOM Element
* @memberof ccUtilities
*/
export function isElement(obj) {
    if(obj===null || obj===undefined) return false;
    return (typeof HTMLElement === "object" ? obj instanceof HTMLElement :
        obj && typeof obj === "object" && obj !== null && obj.nodeType === 1 && typeof obj.nodeName==="string");
}
/**
* Is object empty
*
* @param {*} obj object to check if its empty
* @returns true if the object is empty (array too)
* @memberof ccUtilities
*/
export function isEmpty(obj) {
    if (obj === undefined || obj === null) return true;
    if (Array.isArray(obj)) return obj.length < 1;
    return obj.toString().length < 1;
}
/**
* Determines if this string is HTML
*
* @param {*} str
* @return {*}
 * @memberof ccUtilities
*/
export function isHtml(str) {
    const regex = /<[^>]+>/;
    return regex.test(str);
}
/**
* Determines if a given string is JSON
*
* @param {*} str string to check
* @returns true if the string is valid JSON
* @memberof ccUtilities
*/
export function isJson(str) {
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
/**
* determines if the given object (jQuery, node, dom element) is an node
*
* @param {*} obj object to check
* @returns true if the object is a node
* @memberof ccUtilities
*/
export function isNode(obj){
    return (
        typeof Node === "object" ? obj instanceof Node :
            obj && typeof obj === "object" && typeof obj.nodeType === "number" && typeof obj.nodeName==="string");
}
/**
* Is object not empty
*
* @param {*} obj object to check if its not empty
 * @returns true if the object or array HAS CONTENT
* @memberof ccUtilities
*/
export function isNotEmpty(obj) {
    if (obj === undefined || obj === null) return false;
    return isEmpty(obj) !== true;
}
/**
* Checks if the given object is a numeric
*
* @param {*} obj object to check
* @returns true if the object is a numeric value
* @memberof ccUtilities
*/
export function isNumeric(obj) {
    return !isNaN(parseFloat(obj)) && isFinite(obj);
}
/**
* This is the key function of jHtml (Json represented HTML).  This fast method changes jHTML into HTML
* so that you can condense html using JSON.
 * Example:
*      <div class="someClass" id="someId" role="someRole"><div>Blah Blah</div></div>
* becomes
*      {"c:"someCLass","i":"someId","role":"someRole","b":[{"t":"Blah Blah"}]}
*
* @param {*} jsonIn the jHtml to convert to HTML
* @param {*} tabs tabs to prepend to the html added (for prettying the HTML)
* @returns HTML
* @memberof ccUtilities
*/
export function jsonToHtml(jsonIn, tabs) {
    if (typeof tabs === "undefined") tabs = 2;
    let rtn = "";
    if (Array.isArray(jsonIn)) {
        for(const jObj of jsonIn){
            rtn += jsonToHtml(jObj, 1);
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
                            children += jsonToHtml(obj,tabs + 1);
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
                            if(key==="areq") fullKey="aria-required";
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
* this compiles an error object to send to logError to send to the server
*
* @param {*} logMessage Error subject (Less verbose)
* @param {*} logException Exception thrown (if available)
* @param {*} className Name of javascript class or file where the error occurred
* @param {*} methodName Name of the method/function the error occurred in
* @param {*} logSeverity On a scale from 0 to 9 the severity
* @param {*} logType Type (like console log) (ie Info, Warning, Error, Security)
* @param {*} logApp This should be set to the overall app of javascript or also the file where it occurred
* @param {*} server Server (in a cluster) this happens on (Fairly irrelevant for javascript)
* @param {*} taskDefinitionIdentifier task definition identifier if known
* @param {*} taskDefinitionTitle task definition title
* @memberof ccUtilities
*/
export function errorLog(
    logMessage,
    logException,
    className,
    methodName,
    logSeverity,
    logType,
    logApp,
    server,
    taskDefinitionIdentifier,
    taskDefinitionTitle
) {
    logMessage=(logMessage) ? logMessage : "ERROR?";
    logException=(logException) ? logException : "No Exception Given";
    className=(className) ? className : controllerName;
    methodName=(methodName) ? methodName : actionName;
    logSeverity=(logSeverity) ? logApp : 1;
    logType=(logType) ? logType : "Error";
    logApp=(logApp) ? logApp : "_globalJavascript";
    server=(server) ? server : "";
    taskDefinitionIdentifier=(taskDefinitionIdentifier) = typeof TaskDefinitionIdentifier !== "undefined" ? taskDefinitionIdentifier : cso.CurrentTask.taskDefinitionIdentifier || -1;
    taskDefinitionTitle=(taskDefinitionTitle) = typeof TaskDefinitionTitle !== "undefined" ? taskDefinitionTitle : cso.CurrentTask.taskDefinitionTitle || "";
    logError({
        LogMessage: logMessage,
        LogException: logException,
        ClassName: className,
        MethodName: methodName,
        LogSeverity: logSeverity,
        LogType: logType,
        LogApp: logApp,
        Server: server,
        TaskDefinitionIdentifier: taskDefinitionIdentifier,
        TaskDefinitionTitle: taskDefinitionTitle,
        AddUserIdentifier: currentUserId || -1
    });
}
/**
*  this allows sending error logs to the server from javascript so that errors can be recorded in browser code
*
* @param {*} error error log object compiled by error log method
* @memberof ccUtilities
*/
export async function logError(error) {
    let msg = [];
    for (const prop in Object.keys(error)) {
        msg.push({ fieldName: prop, fieldValue: error[prop] });
    }
    let messageId = "LOGERRORkey" + guid();
    let messageMethod = "/TaskSec/LogEntry";
    unobtrusiveWait("Standby.  Logging an Error"); //Infrastructure js function that hides the waiting display
    let data = await simplePost(messageMethod,error);
    if (data.hasOwnProperty("errorObject")) {
        reportErrorObject(data);
    } else {
        unobtrusiveWaitOff();
    }
}
/**
* This checks for errors from the server and logs the errors
*
* @param {*} data data received from server
* @param {*} toDb true to send the error log back to the server
* @param {*} call function to call on error
* @memberof ccUtilities
*/
export async function reportErrorObject(data, toDb,call) {
    if (toDb === undefined) toDb = false;
    if ((data.fireRob || false) === true) navBar.rulesOfBehaviorShow();
    let errorObject=(data.errorObject);
    if (startsWith(errorObject,"Message Failed to make it using var") || data.bypassMessage) {
                    console.error(timeOfDay() +": handleCalls UTIL LOAD FAILED/UTIL_LOAD/ " + errorObject);
                    if(call) call & call();
    } else {
                    if(data.hasOwnProperty("isPopUp")){
                                    if (data.isPopUp) {
                                                    window.open(data.popUpUrl);
                                                    messageBox('A Popup should have occurred.  If it did not please turn of your popup blocker(IE) or look in your browser for pop up cache(Chrome)', 'ERROR POPUP');
                                                    call & call();
                                    } else {
                                                    if(data.additionalInfo) data.additionalInfo="";
                                                    if(data.additionalInfo.length>2) {
                                                                    messageBox(errorObject + "<hr/>Additional Info:<hr/>" + data.additionalInfo, "Ajax Action Failure!");
                                                    }else{
                                                                    messageBox(data.additionalInfo, data.buttonType);
                                                    }
                                                    call & call();
                                    }
                    }else{
                                    console.log(timeOfDay() +": error :" + errorObject);
                                    if(errorObject.indexOf("The user aborted a request")<0){
                                                    messageBox(errorObject, "!!!ERROR!!!");
                                    }
                    }
    }
    if (toDb) {
                    await logError({
                                    LogMessage: errorObject,
                                    LogException: errorObject,
                                    ClassName: controllerName || "UNKNOWN",
                                    MethodName: actionName || "UNKNOWN",
                                    LogSeverity: 5,
                                    LogType: "Error",
                                    LogApp: "utilities.js",
                                    Server: "",
                                    TaskDefinitionIdentifier: cso.CurrentTask.TaskDefinitionIdentifier || -1,
                                    TaskDefinitionTitle: cso.CurrentTask.TaskDefinitionTitle || "",
                                    AddUserIdentifier: currentUserId || -1
                    });
    }
}
/**
* Creates an bootstrap 5 alert element
 *
* @param {*} alert alert text
* @param {*} alias identifier for this alert
* @param {*} type type of alert (ie dnager, warning, etc)
* @returns the alert jHtml
* @memberof ccUtilities
*/
export function makeAlert(alert, alias, type) {
    if (alias === undefined) alias = guid();
    if (type === undefined) type = "danger";
    let level = type === "danger" ? "Error!": type === "warning" ? "Warning!" : type === "info"? "Attention!": "Great News!";
    let alrt = { i: alias + "Alert" + type,c: "alert alert-" + type + "",title: "Alert for " + alert,   babies: [
        { n: "button",i: alias + "Alert" + type + "Btn",type: "button",c: "Dismiss","data-dismiss": "alert","aria-label": "Close",babies: [
            {n: "span",i: alias + "Alert" + type + "BtnSpan","aria-hidden": "true",t: "&times;"
            }]
        },
        { n: "strong",i: alias + "Alert" + type + "Text",t: level + ": " + alert }
    ]};
    return alrt;
}
/**
   * Make an alert using JML
   *
   * @param {*} idAdd - Additional ID to add to the alert
   * @param {*} text - Text to display in the alert
   * @param {*} clazz - Class to add to the alert
   * @return {*} Alert JML
   * @memberof ccUtilities
   */
export function makeAlertJML(idAdd, text, clazz) {
    let ccac = this;
    if (!idAdd) idAdd = "";
    if (!clazz) clazz = "alert alert-danger alert-dismissible fade show";
    let alert = {
      i: guid().replaceAll("-", "") + idAdd,
      title: text,
      t: text,
      c: clazz,
    };
    if (clazz.includes("alert-dismissible")) {
      alert.b = [];
      let alertButton = {
        n: "button",
        type: "button",
        title: "Close Alert",
        c: "btn-close",
        "data-bs-dismiss": "alert",
      };
      alert.b.push(alertButton);
    }
    return alert;
  }
/**
* Changes the property name of a given object to camel notation
*
* @param {*} obj object to fix property names on
* @returns same object with corrected property names (C# to ES6) for javascript
* @memberof ccUtilities
*/
export function objectKeysToCamelCase(obj) {
    return Object.keys(obj).reduce(function(newObj, key) {
        let val = obj[key];
        let newVal = val !== null && typeof val === "object" ? objectKeysToCamelCase(val) : val;
        newObj[asFieldNotationString(key)] = newVal;
        return newObj;
    }, {});
}
/**
* determines if a string starts with the given search string
*
* @param {*} str the string to check if it starts with the check string
* @param {*} check the string to check for at the begining of str
* @param {*} ignore ignore case
* @returns
* @memberof ccUtilities
*/
export function startsWith(str, check, ignore) {
    if (check === undefined || check === null) return false;
    if (str === undefined || str === null) return false;
    if (typeof ignore === "undefined") ignore = true;
    str = "" + str + "";
    check = "" + check + "";
    if (ignore) return str.slice(0, check.length) == check;
    return str.slice(0, check.length) === check;
}
/**
* Converts number into a string version (ie 75 becomes Seventy fifth)
*
* @param {*} n number to stringify
* @returns string representation of number
* @memberof ccUtilities
*/
export function stringifyNumber(n) {
    let special = [
        "zeroth",
        "first",
        "second",
        "third",
        "fourth",
        "fifth",
        "sixth",
        "seventh",
        "eighth",
        "ninth",
        "tenth",
        "eleventh",
        "twelfth",
        "thirteenth",
        "fourteenth",
        "fifteenth",
        "sixteenth",
        "seventeenth",
        "eighteenth",
        "nineteenth"
    ];
    let deca = [
        "twent",
        "thirt",
        "fort",
        "fift",
        "sixt",
        "sevent",
        "eight",
        "ninet"
    ];
    if (n < 20) return special[n];
    if (n % 10 === 0) return deca[Math.floor(n / 10) - 2] + "ieth";
    return deca[Math.floor(n / 10) - 2] + "y-" + special[n % 10];
}
 
/**
* Gets the text between two strings
*
* @export
* @param {*} str string to get text from
* @param {*} firstFind first enveloping string
* @param {*} secondFind second enveloping string
* @return {*} string in the middle of the two strings
*/
export function textBetween(str, firstFind, secondFind) {
    let pos1 = str.indexOf(firstFind) + firstFind.length;
    let pos2 = str.indexOf(secondFind, pos1);
   
    if (pos1 < 0 && pos2 < 0) {
        return str;
    }
   
    if (pos1 < 0) {
        pos1 = 0;
    }
   
    if (pos2 < 0) {
        pos2 = str.length;
    }
    return str.substring(pos1,pos2);
}
/**
* Capitalizes the first letter of each word in a string
*
* @export
* @param {*} str
* @return {*}
 * @memberof ccUtilities
*/
export function toTitle(str) {
    if(!str) return "";
    let rtn="";
    let titleArr=str.split(" ");
    for(let word of titleArr){
        rtn+=word.substring(0,1).toUpperCase()+word.substring(1).toLowerCase()+" ";
    }
    return rtn.trim();
}
/**
* Trims whitespace from a string
*
* @param {*} str string to trim
* @returns string after trim
 * @memberof ccUtilities
*/
export function trim(str) {return str.replace(/^\s+|\s+$/g, "");}
//                  ~~~~~~~~~~~~~~~~~~~~ Code Functions END~~~~~~~~~~~~~~~~~                   \\
//                  ~~~~~~~~~~~~~~~~~~~~ Form Functions ~~~~~~~~~~~~~~~~~~~~                   \\
 
/**
* Clear a form of data
*
* @param {*} form form element to clear
* @memberof ccUtilities
*/
export function formClear(form) {
    if(form===undefined) form=getEle("form");
    if(typeof form==="string") form=getEle(form.asElementId());
    let inputs = form.getEles("input, select, checkbox, textarea");
    for(const item of inputs){
        let fldType = item.utlAttr("data-field-type");
        if (fldType == undefined) fieldType = item.utlAttr("type");
        switch (fldType) {
            case "search":
                item.utlVal("");
                break;
            case "CheckBox":
            case "checkbox":
                item.checked=false;
                break;
            case "Display":
                item.innerText="";
                break;
            default:
                item.utlVal("");
                break;
        }
    }
}
/**
* enables or disables a forms input elements
*
* @param {*} form form element
* @param {*} isEnabled if it should be enabled or disabled
* @memberof ccUtilities
*/
export function formEnabled(form, isEnabled) {
    if(form===undefined) form=getEle("form");
    if(typeof form==="string") form=getEle(form.asElementId());
    for (const inpt of form.querySelectorAll("input, select, checkbox, textarea")) {
        inpt.disabled = !isEnabled;
    }
}
/**
* fill the form created by formFromObject with data given in config.theObject
*
* @param {*} form form element
* @param {*} currForm data for form
* @memberof ccUtilities
*/
export function formFill(form, data) {
    if(form===undefined) form=getEle("form");
    if(typeof form==="string") form=getEle(form.asElementId());
    if(!form.id) form.id=guid().replaceAll("-","");
    form.data = data;
    let inputs = form.getEles("input, select, checkbox, textarea");
    for(const item of inputs){
        let prop=item.dataset.msgProperty||item.id;
        let field = { Name: prop, Value: data[prop] };
        if (field !== undefined) {
            switch (item.type) {
                case "checkbox":
                    item.checked=asBool(field.Value);
                    break;
                case "textarea":
                    {
                        item.utlVal(field.value);
                        let valSafe = field.value || "";
                        let valLen = valSafe.length;
                        let lines = valSafe.split(/\r\n|\r|\n/).length;
                        if (lines > 1) {
                            if (lines > 4) {
                                let hgt = 25 * lines;
                                if (lines < 50) item.utlCSS("height", hgt);
                            }
                        } else {
                            if (valLen > 360) {
                                let charLines = valLen / 70;
                                let hgt = 25 * charLines;
                                if (charLines < 50) item.utlCSS("height", hgt);
                            }
                        }
                    }
                    break;
                default:
                    item.utlVal(field.value);
                    break;
                }
            }
    }
    form.formEnabled(true);
    unobtrusiveWaitOff();
    functionCall(form.id + "LoadEnd", data);
}              
/**
* Generates the message for a form to send up
*
* @param {*} form form DOM element
* @param {*} additionalProperties Any properties to send up outside of the form
 *                                 as part of the object(ex. {recordId:1, formType:"SomeType"})
* @returns message for transmitting form data
* @memberof ccUtilities
*/
export function formToJson(form,additionalProperties) {
    if(form===undefined) form=getEle("form");
    if(typeof form==="string") form=getEle(form.asElementId());
    let currForm = form;
    let formId = form.utlAttr("id");
    let formType = form.utlAttr("data-formType") || "UNKNOWN";
    let rid = form.utlAttr("data-recordId") || "-1";
    let rtnJson = {};
    let inputs = form.getEles(
        "input:not([type=button]):not([type=submit]), select, checkbox, textarea"
    );
    for (const item of inputs) {
        let itemId = item.utlAttr("id");
        let itemVal = item.utlVal();
        rtnJson[itemId]=itemVal;
    }
    if(additionalProperties) rtnJson=Object.assign(rtnJson, additionalProperties);
    rtnJson = functionCall(currForm.objectName + "MsgAppend", rtnJson);
    return rtnJson;
}
 
//                  ~~~~~~~~~~~~~~~~~~~~ Form Functions END~~~~~~~~~~~~~~~~~                   \\
 
//                  ~~~~~~~~~~~~~~~~ Create Files Functions END~~~~~~~~~~~~~                   \\
/**
* Given a ICS configuration (JSON) this will compile a ICS Event (Coding changes may be required for other ICS types)
*
* @param {*} cfg config JSON
* @returns the icsText
* @memberof ccUtilities
*/
export function makeICS(cfg){
    var icsFile=null;
    let icsText="BEGIN:VCALENDAR\n";
    icsText+="VERSION:" + (cfg.VERSION?cfg.VERSION:"2.0") + "\n";
    icsText+="CALSCALE:" + (cfg.CALSCALE?cfg.CALSCALE:"GREGORIAN") + "\n";
    icsText+="METHOD:" + (cfg.METHOD?cfg.METHOD:"GREGORIAN") + "\n";
    icsText+="PRODID:" + (cfg.PRODID?cfg.METHOD:"-//OIT//IUSG//DASM//CAL EVENTS//EN") + "\n";
    icsText+="BEGIN:" + (cfg.BEGIN?cfg.BEGIN:"VEVENT") + "\n";
    icsText+="UID:" + (cfg.UID?cfg.UID: apestalmenos.guid()) + "\n";
    if(cfg.DTSTAMP) icsText+="DTSTAMP:" + cfg.DTSTAMP + "\n";
    if(cfg.ORGANIZER) icsText+="ORGANIZER;CN=" + cfg.ORGANIZER +"\n";
    if(cfg.DTSTART) icsText+="DTSTART:" + cfg.DTSTART + "\n";
    if(cfg.DTEND) icsText+="DTEND:" + cfg.DTEND + "\n";
    if(cfg.DESCRIPTION) icsText+="DESCRIPTION:" + cfg.DESCRIPTION + "\n";
    if(cfg.LOCATION) icsText+="LOCATION:" + cfg.LOCATION + "\n";
    icsText+="END:" + (cfg.BEGIN?cfg.BEGIN:"VEVENT") + "\n";
    icsText+="END:VCALENDAR";
    // let data = new File([icsText], { type: "text/calendar" });
    // if (icsFile !== null) {
    //         (window.URL ? window.URL : window.webkitURL).revokeObjectURL(icsFile);
    // }
    // icsFile = (window.URL ? window.URL : window.webkitURL).createObjectURL(data);
    return icsText;
}
/**
* This will set a onClick function on a link,button,etc which uses the given icsConfig to make
* an ICS file and cause a download of it for the client
*
* @param {*} actionElement The element to set the onClick function on
* @param {*} icsConfig the config for the ICS (SEE makeICS)
* @memberof ccUtilities
*/
export function setICSFunc(actionElement,icsConfig){
    getEle(actionElement.utlAsElementId()).addEventListener("click", function(e){
        window.open( "data:text/calendar;charset=utf8," + escape(makeICS(icsConfig)));
    });
}
 
 
 
//                  ~~~~~~~~~~~~~~~~ Create Files Functions END~~~~~~~~~~~~~                   \\
 
//                  ~~~~~~~~~~~~~ Independent List Functions ~~~~~~~~~~~~~~~                   \\
 
 
 
//                  ~~~~~~~~~~~~~ Independent List Functions END ~~~~~~~~~~~                   \\
 
//                  ~~~~~~~~~~~~~~~~~~~~ External Functions ~~~~~~~~~~~~~~~~~~~~               \\
/**
* ccUtilities messageBox function uses modal to perform the ole messageBox global method
*
* @param {*} text Text inside the modal
* @param {*} title title of modal
* @param {*} button1Text Text for first button
* @param {*} button1Script javascript function for first button
* @param {*} button2TextText Text for second button
* @param {*} button2Script javascript function for second button
* @param {*} button3Text Text for third button
* @param {*} button3Script javascript function for third button
* @returns the modal if you need it immediately otherwise you can always get it by getEle("#utlModalContent").modal or $("#utlModalContent")[0].modal
* @memberof ccUtilities
*/
export function messageBox(text, title, button1Text, button1Script, button2Text, button2Script, button3Text, button3Script) {
    let mdlObj = {
        message: text,
        title: title,
        buttons: []
    };
    if(button1Text)mdlObj.buttons.push({text: button1Text, e: button1Script});
    if(button2Text)mdlObj.buttons.push({text: button2Text, e: button2Script});
    if(button1Text)mdlObj.buttons.push({text: button3Text, e: button3Script});
    let mdl=ccModalMakeOpen=(mdlObj.message, mdlObj.title, mdlObj.buttons);
    return mdl;
}
/**
* pops up a small unobtrusive div with info(if given) to alert the user to ajax or promise functions which may take some time to complete
*
* @param {*} message message to give the user (defaults to ...Loading Please Wait...)
* @param {*} id an identifier for this alert (a guid is used if one is not provided)
* @param {*} top the top position of the div when it appears
* @param {*} left the left position of the div when it appears
* @param {*} center center the div
* @memberof ccUtilities
*/
export async function unobtrusiveWait(message, id, caller, top, left, center) {
    let uwDiv = getEle(".unobtrusiveWait");
    if (uwDiv.id!=="ElementNotFound") {
        Array.from(getEles(".unobtrusiveWait")).map(e=> e.remove());
    }
    if(!id) id = "unUbtrusiveWaitModal";
    message = (message) ? message : "...Loading Please Wait...";
    top = (top) ? top : 0;
    left = (left) ? left : 0;
    center = (center) ? center : true;
    let hNum = 6-(parseInt(message.length/30)+1);
    if(hNum<1) hNum=1;
    if(hNum>5) hNum=5;
    let cfg={};
    cfg.id="unUbtrusiveWaitModal_" + caller;
    cfg.bodyJML={i: id +"ModalBody", c:"modal-body bg-light text-center",b:[                                                                                                                                  
        {n:"button",c:"btn btn-info w-100",               t: message}
    ]},
    cfg.class="unobtrusiveWait d-block fade show";
    //cfg.messageText=message||"Please Wait...Loading!!";
    //cfg.titleText="";
    cfg.headerJML={c:"modal-header justify-content-center",s: "background-color:#c6c8ca",
        b:[
            {i:id + "text",c:"h" + hNum + " p-2 mx-auto my-auto progress-bar progress-bar-striped progress-bar-animated",
                role:"progressbar",s:"width: 100%;","aria-valuenow":"100","aria-valuemin":"0","aria-valuemax":"100"
            }
        ]};
    cfg["data-src"]=caller;
    cfg.buttons=null;
    cfg.hasCloseXButton=true;
    cfg.footerJML==null;
    cfg.addCloseButton=false;
    await ccModalConfigOpen(cfg);
}
/**
* Removes all (by default) unobtrusiveWait divs or just the one identified if id is given
*
* @param {*} id id of unobtrusiveWait div to remove if not give all will be by class
* @memberof ccUtilities
*/
export function unobtrusiveWaitOff(id) {
    if (!id) {
        Array.from(getEles(".unobtrusiveWait")).map(e=> e.remove());
    } else {
        getEle(id.utlAsElementId()).remove();
    }
    Array.from(document.querySelectorAll(".modal-backdrop")).map(e=>e.remove());
    let body=document.querySelector("body");
    body.classList.remove("modal-open");
    body.style="";
}
//                  ~~~~~~~~~~~~~~~~~~~~ External Functions END~~~~~~~~~~~~~~~~~               \\
 
//                  ~~~~~~~~~~~~~~~~~~~ Comms Functions ~~~~~~~~~~~~~~~~~~~~                   \\
/**
* Makes a simple ajax request to the server.  this is just like a straight fetch with some error handling involved
*
* @param {*} name name of request
* @param {*} ajaxCall url to call
* @param {*} dataUp data to send up
* @param {*} messageReceiverDataType type of data to send.  defaults to json
* @returns data sent back from server
* @memberof ccUtilities
*/
export async function simpleCall(name,ajaxCall,dataUp,messageReceiverDataType,httpMethod) {
    let spa = this;
    if (typeof sessionActivity === "function") sessionActivity();
    if (messageReceiverDataType === undefined) messageReceiverDataType = "json";
    if (!httpMethod) httpMethod = "GET";
    let res = {};
    let resObj = {};
    try {
      res = httpMethod === "POST"
          ? await fetch(ajaxCall,
              {method: httpMethod,cache: "no-cache",
              headers: {"Content-Type": "application/json; charset=utf-8"},
              body: JSON.stringify(dataUp),
            })
          : await fetch(ajaxCall,
            {method: httpMethod,cache: "no-cache",
            headers: {"Content-Type": "application/json; charset=utf-8"},
              //signal: window.abortCtrl
            });
      resObj = {ok: res.ok,url: res.url,text: res.statusText,sessionTimeout: false,redirected: res.redirected,status: res.status,
        bodyUsed: res.bodyUsed,
      };
    } catch (err) {
      let logMsg = "";
      let msg = "";
      switch (res.status) {
        case 404:
          logMsg ="Web Server error. Status: " +res.status +"(" +res.statusText +") for location: " +
            res.url +" " +(res.redirected ? "(--redirected--)." : "") +
            ". This URL is not on the server.  Please send a Customer Service Request to fix this.";
          msg = logMsg;
        case 500:
          logMsg ="Web Server error. Status: " + res.status + "(" + res.statusText + ") for location: " +
            res.url + " " + (res.redirected ? "(--redirected--)." : "") +
            ". This URL is not on the server.  Please send a Customer Service Request to fix this.";
        default:
          {
            logMsg = "Web Server error. Status: " + res.status + "(" + res.statusText + ") for location: " +
              res.url + " " + (res.redirected ? "(--redirected--)" : "");
            msg = logMsg;
          }
          break;
      }
      if (logMsg.length) console.log(logMsg);
      resObj.json = msg.length ? { errorObject: msg } : { error: logMsg };
    }
    resObj.json = {};
    try {
      if (res.status === 200) {
        resObj.json = await res.json();
      } else {
        let msg = "";
        let goMsg = true;
        switch (res.status) {
          case 404:
            msg =`Web Server error. Status: ${res.status}(${res.statusText}) for location: ${res.url} ` +
            (res.redirected ? "(--redirected--)." : "") + `. This URL is not on the server.  Please send a Customer Service Request to fix this.`;
            break;
          case 500:
            {
              const resHtml = await res.text();
              if (resHtml.indexOf("set JsonRequestBehavior to AllowGet")) {
                //ignore BS
                resObj.json = [];
                goMsg = false;
              } else {
                msg =
                  "Web Server error. Status: " +
                  res.status +
                  "(" +
                  res.statusText +
                  ") for location: " +
                  res.url +
                  " " +
                  (res.redirected ? "(--redirected--)" : "");
              }
            }
            break;
          default:
            {
              msg =
                "Web Server error. Status: " +
                res.status +
                "(" +
                res.statusText +
                ") for location: " +
                res.url +
                " " +
                (res.redirected ? "(--redirected--)" : "");
            }
            break;
        }
        console.log(msg);
        if (goMsg) resObj.json = { errorObject: msg };
      }
    } catch (err) {
      const contentType = res.headers.get("content-type") || "unknown";
      if (contentType.indexOf("application/json;") > -1 || res.url.indexOf("/Account/AccessError") > -1) {
        let text="";
        try{text = (await res.text()) || "GARBAGE";}catch(err){text = "GARBAGE:" + err.message;}
        if(res.url.indexOf("/Account/AccessError") > -1 ){
            ccModalMakeOpen("You do not have access to " + ajaxCall + ".  Please contact your administrator.","Access Error!!");
        }
        if (res.url.indexOf("/Task/ROB") > -1) {
          console.log("User Session Expired?  Calling ROB");
          resetTimer("Terminated");
          throw Error({ text: "Your Session Has ended", sessionTimeout: true });
        }
        console.log(err);
        resObj.json = { errorObject: err };
      }
    }
    if (resObj.json.errorObject)
    {  
        unobtrusiveWaitOff();
        ccModalMakeOpen(resObj.json.errorObject, "Error",
        [{n:"button",type:"button",c:"btn btn-secondary",i:"ccModalFooterBtnClose","data-bs-dismiss":"modal",t:"Close",e: "ccModalGlobal.close();"}]);
    }   
    if (resObj.url && resObj.url.indexOf("Account/AccessError") > -1) {
        resObj.json = { errorObject: "Access Error Occurred on :" + resObj.url };
    }
    return resObj.json;
}
 
// export async function timeOfDay(){
//     d = new Date();
//     datetext = d.toTimeString();
//     // datestring is "20:32:01 GMT+0530 (India Standard Time)"
//     // Split with ' ' and we get: ["20:32:01", "GMT+0530", "(India", "Standard", "Time)"]
//     // Take the first value from array :)
//     datetext = datetext.split(' ')[0];
// }
 
/**
* This is the main function to call to get data from the API.  It is called by the other functions in this class.
*
* @param {*} ajaxCall The API URL to call
* @param {*} dataUp The data to send up to the API
* @param {*} messageReceiverDataType The type of data to return (json, text, etc.)
* @return {*}  {Promise<any>}
* @memberof ccUtilities
*/
export async function simplePost(ajaxCall, dataUp, messageReceiverDataType) {
let ccac = this;
//if (typeof sessionActivity === "function") sessionActivity();
let trackId = (util.guid()).toString();
if (messageReceiverDataType === undefined) messageReceiverDataType = "json";
let httpMethod = "POST";
let res = {};
let resObj = {};
//ccac.unobtrusiveWait("Getting Data!!");
if (httpMethod === "POST") {
    let fcUp = new FormData();
    // for ( let key in dataUpJson ) {
    //     fcUp.append(key, dataUpJson[key]);
    // }
    fcUp.append("data", dataUp);
    // for (let [key, val] of Object.entries(dataUp)) {
    //     fcUp.append(key, val);
    // }
    fcUp.append("trackId", trackId);
    dataUp = fcUp;
}
try {
    res = await fetch(ajaxCall,{
    method: "POST",
    cache: "no-cache",
    body: dataUp,
    signal: signal,
    });
    resObj = {
    ok: res.ok,
    url: res.url,
    text: res.statusText,
    sessionTimeout: false,
    redirected: res.redirected,
    status: res.status,
    bodyUsed: res.bodyUsed,
    };
} catch (err) {
    let logMsg = "";
    let msg = "";
    switch (res.status) {
    case 404:
        logMsg =
        "Web Server error. Status: " +
        res.status +
        "(" +
        res.statusText +
        ") for location: " +
        res.url +
        " " +
        (res.redirected ? "(--redirected--)." : "") +
        ". This URL is not on the server.  Please send a Customer Service Request to fix this.";
        msg = logMsg;
    case 500:
        logMsg =
        "Web Server error. Status: " +
        res.status +
        "(" +
        res.statusText +
        ") for location: " +
        res.url +
        " " +
        (res.redirected ? "(--redirected--)." : "") +
        ". This URL is not on the server.  Please send a Customer Service Request to fix this.";
    default:
        {
        logMsg =
            "Web Server error. Status: " +
            res.status +
            "(" +
            res.statusText +
            ") for location: " +
            res.url +
            " " +
            (res.redirected ? "(--redirected--)" : "");
        msg = logMsg;
        }
        break;
    }
    if (logMsg.length) console.log(logMsg);
    resObj.json = msg.length ? { errorObject: msg } : { error: logMsg };
}
resObj.json = {};
resObj = await commsResponseStatusHandle(res, resObj.json);
if (resObj.url.indexOf("Account/AccessError") > -1)
    resObj.json = { errorObject: "Access Error Occurred on :" + resObj.url };
//ccac.unobtrusiveWaitOff();
return resObj.json;
}
 
/**
* This is a function to handle the response status of the fetch call.
*
* @param {*} res - the response object
* @param {*} resObj - the response object
* @return {*}
 * @memberof ccAutoComplete
 */
export async function commsResponseStatusHandle(res, resObj) {
let ccac = this;
try {
    if (res.status === 200) {
    resObj.json = await res.json();
    } else {
    let msg = "";
    let goMsg = true;
    switch (res.status) {
        case 404:
        msg =
            "Web Server error. Status: " +
            res.status +
            "(" +
            res.statusText +
            ") for location: " +
            res.url +
            " " +
            (res.redirected ? "(--redirected--)." : "") +
            ". This URL is not on the server.  Please send a Customer Service Request to fix this.";
        break;
        case 500:
        {
            const resHtml = await res.text();
            if (resHtml.indexOf("set JsonRequestBehavior to AllowGet")) {
            //ignore BS
            resObj.json = [];
            goMsg = false;
            } else {
            msg =
                "Web Server error. Status: " +
                res.status +
                "(" +
                res.statusText +
                ") for location: " +
                res.url +
                " " +
                (res.redirected ? "(--redirected--)" : "");
            }
        }
        break;
        default:
        {
            msg =
            "Web Server error. Status: " +
            res.status +
            "(" +
            res.statusText +
            ") for location: " +
            res.url +
            " " +
            (res.redirected ? "(--redirected--)" : "");
        }
        break;
    }
    console.log(msg);
    if (goMsg) resObj.json = { errorObject: msg };
    }
} catch (err) {
    const contentType = res.headers.get("content-type") || "unknown";
    if (contentType.indexOf("application/json;") > -1) {
    const text = (await res.text()) || "GARBAGE";
    if (
        res.url.indexOf("/Account/AccessError") > -1 ||
        res.url.indexOf("/Task/ROB") > -1
    ) {
        console.log("User Session Expired?  Calling ROB");
        resetTimer("Terminated");
        throw Error({ text: "Your Session Has ended", sessionTimeout: true });
    }
    console.log(err);
    resObj.json = { errorObject: err };
    }
}
if (!resObj.url) resObj.url = window.location.href;
return resObj;
}
 
/**
* Use to communicate to classes/components calling a function
*
* @param {*} receiver the window level (if possible) class or component to call
* @param {*} caller the window level (if possible) class or component calling the function
* @param {*} job function name
* @param {*} data any data to be sent to that function
* @return {*}
 * @memberof catsCRUDL
*/
  export async function dispatch(receiver, caller, job, data){
    let res = await receiver[job](data);
    return res;
  }
 
/**
* This checks for errors from the server calls made by catsCRUDL and logs the errors receipt
*
* @param {*} data data received from server
* @param {*} toDb true to send the error log back to the server
* @param {*} call function to call on error
* @memberof ccUtilities
*/
export async function receiptCheckGood(data, toDb, call) {
    if (toDb === undefined) toDb = false;
    if ((data.fireRob || false) === true) navBar.rulesOfBehaviorShow();
    if (typeof data === "string" && isJson(data)) data = JSON.parse(data);
    if (data.hasOwnProperty("Data") && data.Data.hasOwnProperty("errorObject"))
      data = data.Data;
    if (data.hasOwnProperty("Data") && data.Data.hasOwnProperty("isValid"))
      data = data.Data;
    if (data.hasOwnProperty("Data") && data.Data.hasOwnProperty("form"))
      data = data.Data;
    let returnValue = true; //true means move on with callers code while false should end here
    if (data.hasOwnProperty("errorObject")) {
      returnValue = false;
      let errorObject = data.errorObject;
      if (
        errorObject.utlStartsWith("Message Failed to make it using var") ||
        data.bypassMessage
      ) {
        console.error(
          timeOfDay() +
            ": handleCalls DHO LOAD FAILED/DHOLISTLOAD/ " +
            errorObject
        );
        if (call) call & call();
      } else {
        if (data.hasOwnProperty("isPopUp") && data.hasOwnProperty("popUpUrl")) {
          if (data.isPopUp) {
            window.open(data.popUpUrl);
            util.messageBox(
              "A Popup should have occurred.  If it did not please turn of your popup blocker(IE) or look in your browser for pop up cache(Chrome)",
              "ERROR POPUP"
            );
            unobtrusiveWaitOff();
            call & call();
          } else {
            if (data.additionalInfo) data.additionalInfo = "";
            if (data.additionalInfo.length > 2) {
              util.messageBox(
                errorObject +
                  "<hr/>Additional Info:<hr/>" +
                  data.additionalInfo,
                "Ajax Action Failure!"
              );
            } else {
              util.messageBox(data.additionalInfo, data.buttonType);
            }
            unobtrusiveWaitOff();
            call & call();
          }
        } else {
          console.log(timeOfDay() + ": error :" + errorObject);
          if (errorObject.indexOf("The user aborted a request") < 0) {
            util.messageBox(errorObject, "!!!ERROR!!!");
          }
        }
      }
      if (toDb) {
        await util.logError({
          LogMessage: errorObject,
          LogException: errorObject,
          ClassName: controllerName || "UNKNOWN",
          MethodName: actionName || "UNKNOWN",
          LogSeverity: 5,
          LogType: "Error",
          LogApp: "utilities.js",
          Server: "",
          TaskDefinitionIdentifier:
            cso.CurrentTask.TaskDefinitionIdentifier || -1,
          TaskDefinitionTitle: cso.CurrentTask.TaskDefinitionTitle || "",
          AddUserIdentifier: currentUserId || -1,
        });
      }
    }
    if (data.hasOwnProperty("um")) {
      await unobtrusiveWait(
        data.um.message,
        "um" + util.guid().utlReplaceAll("-", ""),
        0,
        200,
        true,
        data.um.timeOut || 30000
      );
    }
    if (data.hasOwnProperty("redirect")) {
      let redirect = data.redirect;
      if (typeof redirect === "string") {
        let surl = redirect.utlStartsWith("http") ? redirect : host + redirect;
        setTimeout((window.location.href = surl), 3000);
      } else {
        let url = redirect.url.utlStartsWith("http")
          ? redirect.url
          : host + redirect.url;
        setTimeout((window.location.href = url), redirect.delay || 3000);
      }
    }
    if (data.hasOwnProperty("isValid")) {
      let frmObj = null;
      if (data.hasOwnProperty("form")) {
        frmObj =
          (typeof data.form === "string" ? JSON.parse(data.form) : data.form) ||
          null;
      }
      if (frmObj) {
        let frm = new utlForm(frmObj.ObjectName, frmObj);
        if (data.form) delete data.form;
        var validateData = data;
        let validationObject = new validateObject(frm, validateData);
        if (!validationObject.isValid) {
          validationObject.fieldNotesMissing = 1;
          validationObject.validateFailed(validationObject);
          _.utlEle("#" + frm.alias + "CardFooter").utlShow();
        }
      } else {
        let isValid =
          data.isValid === true || data.isValid === false
            ? data.isValid
            : false;
        if (!isValid) {
          util.messageBox(
            data.additionalInfo || "Something went wrong!!",
            data.buttonType || "Error Occured!"
          );
        }
      }
    }
    return returnValue;
  }
 
//                  ~~~~~~~~~~~~~~~~~~~ Comms Functions END~~~~~~~~~~~~~~~~~                   \\
 
//                  ~~~~~~~~~~~~~~~~~ Page Function Calls ~~~~~~~~~~~~~~~~~~                   \\
/**
* This allows placing function stubs within other functions that are optional within them
*
* @param {*} func function to call (if it exists)
* @param {*} data data to pass (if given)
* @returns function or just data
* @memberof ccUtilities
*/
export function functionCall(func, data) {
    if (data == undefined) {
        if (typeof window[func] === "function") {
            return window[func]();
        }
    } else {
        if (typeof window[func] === "function") {
            return window[func](data);
        }
    }
    return data;
}
/**
* This allows placing function stubs within other functions that are optional within them
*
* @param {*} func function to call (if it exists)
* @param {*} data data to pass (if given)
* @returns function or just data
* @memberof ccUtilities
*/
export async function functionCallAsync(func, data) {
    if (data == undefined) {
        if (typeof window[func] === "function") {
            return await window[func]();
        }
    } else {
        if (typeof window[func] === "function") {
            return await window[func](data);
        }
    }
    return data;
}
/**
* Records what DOM element the user last used
*
* @param {*} ele element to restore location later on
* @memberof ccUtilities
*/
export function setWhereWasI(ele){
    whereWasI=ele;
}
/**
* using the property whereWasI this uses the utlGoTo function to return the user to where they were on the page
 * before a popup (such as a form) happened
*
* @memberof ccUtilities
*/
export function takeMeBack() {
    if (whereWasI !== undefined) whereWasI.utlGoTo();
    whereWasI = undefined;
}
 
/**
* Adjusts popups based on quadrants
*
* @export
* @param {*} activator button that activates popup
* @param {*} popup popup div
* @param {*} xAdj x adjustment
* @param {*} yAdj y adjustment
* @memberof ccUtilities
*/
export function adjustPopupToKey(activator,popup,xAdj,yAdj){
    if(!xAdj) xAdj=15;
    if(!yAdj) yAdj=15;
    let winHeight = window.innerHeight || document.documentElement.clientHeight || document.body.clientHeight;
    let winHeightThird=winHeight*.333333;
    let winWidth = window.innerWidth || document.documentElement.clientWidth || document.body.clientWidth;
    let winWidthThird = winWidth*.333333;
    let actTop=activator.utlOffset().top;
    let actLeft=activator.utlOffset().left;
    let actWidth=activator.offsetWidth;
    let actHeight=activator.offsetHeight;
    let pupTop=popup.utlOffset().top;
    let pupLeft=popup.utlOffset().left;
    let pupWidth=popup.offsetWidth;
    let pupHeight=popup.offsetHeight;
    let pupBottom=pupTop+pupHeight;
    let pupRight=pupLeft+pupWidth;
    let actBottom=actTop+actHeight;
    let actRight=actLeft+actWidth;
    let actXCenter=actTop+(actHeight/2);
    let actYCenter=actLeft+(actWidth/2);
    let actQuadrent="";
    let actYSection="";
    if(actTop>=0 && actTop <= (winHeightThird)){
        actQuadrent="T";
    }else if(actTop>=(winHeightThird+1) && actTop <= (winHeightThird*2)){
        actQuadrent="C";
    }else{
        actQuadrent="B";
    }
    if(actLeft>=0 && actLeft <= (winWidthThird)){
        actQuadrent+="L";
    }else if(actLeft>=(winWidthThird+1) && actLeft <= (winWidthThird*2)){
        actQuadrent+="C";
    }else{
        actQuadrent+="R";
    }
    switch(actQuadrent){
        case "TL":
            popup.style.left=actRight+xAdj + "px";
            popup.style.top=actBottom+yAdj + "px";
            break;
        case "TC":
            popup.style.left=actXCenter+xAdj + "px";
            popup.style.top=actBottom+yAdj + "px";
            break;
        case "TR":
            popup.style.left=actLeft-(pupWidth+xAdj) + "px";
            popup.style.top=actBottom+(yAdj+pupHeight) + "px";
            break;
        case "BL":
            popup.style.left=actRight+xAdj + "px";
            popup.style.top=actBottom-(yAdj+pupHeight) + "px";
            break;
        case "BC":
            popup.style.left=actXCenter+xAdj + "px";
            popup.style.top=actBottom-(yAdj+pupHeight) + "px";
            break;
        case "BR":
            popup.style.left=actLeft-(xAdj + pupWidth) + "px";
            popup.style.top=actBottom-(yAdj+pupHeight) + "px";
            break;
        case "CL":
            popup.style.left=actLeft+(pupWidth+xAdj) + "px";
            popup.style.top=actBottom+yAdj + "px";
            break;
        case "CR":
            popup.style.left=actLeft-(pupWidth+xAdj) + "px";
            popup.style.top=actBottom+yAdj + "px";
            break;
        default:
            popup.style.left=actLeft-(pupWidth+xAdj) + "px";
            popup.style.top=actBottom-yAdj + "px";
            //do nothing
            break
    }
 
}
/**
* Opens a new window with given url
 *
* @param {*} url url for popup
* @memberof ccUtilities
*/
export function openUrlClean(url) {
    let newTab = window.open();
    newTab.opener = null;
    newTab.location = url;
}
/**
* get a property by name
 *
* @param {*} theProperty property you wish the value from
* @returns value of field specified
* @memberof ccUtilities
*/
export function get(theProperty) {
    return this[theProperty];
}
/**
* sets the value of a property by name
*
* @param {*} theProperty property you wish to set the value of
* @param {*} value value to set it to
* @returns property value given
* @memberof ccUtilities
*/
export function set(theProperty,value) {
    return this[theProperty]=value;
}
/**
* Time of day for logging
*
* @returns
* @memberof ccUtilities
*/
export function timeOfDay() {
    let d = new Date();
    return d.getUTCHours() + ":" + d.getUTCMinutes() + ":" + d.getUTCSeconds();
}
/**
* Returns details of this object
 *
* @returns details of this object
* @memberof ccUtilities
*/
export function toString() {
    let rtn="";
    for(const [key, value] of Object.entries(this)){
        rtn += key + ": " + value + "<br/>";
    }
    return rtn;
}
//                  ~~~~~~~~~~~~~~~~~ Page Function Calls END~~~~~~~~~~~~~~~                   \\
 
//------------------------------------------------------Support Functions
/**
* Assists classAdd when multiple classes are added (ie. 'classOne classTwo andSomeOtherClass'))
*
* @param {*} ele DOM element use
* @param {*} className  class names to remove from element
 * @returns element so you can chain stuff
* @memberof ccUtilities
*/
function classAddForMulti(ele,className){
    for(const cN of className.split(' ')){
        ele.classList.add(cN);
    }
    return ele;
}
/**
* Assists classAdd when multiple classes are added (ie. 'classOne classTwo andSomeOtherClass'))
*
* @param {*} ele DOM element use
* @param {*} className  class names to remove from element
 * @returns element so you can chain stuff
* @memberof ccUtilities
*/
function classRemoveForMulti(ele,className){
    for(const cN of className.split(' ')){
        ele.className = ele.className.replace(new RegExp("(^|\\b)" + cN.split(" ").join("|") + "(\\b|$)","gi"),"").trim();
    }
    return ele;
}
 
/**
* generates a 4 character member of a guid
*
* @static
* @returns 4 character member of a guid
* @memberof ccUtilities
*/
function guidS4() {
    return asString(Math.floor((1 + Math.random()) * 0x10000), 16).substring(1);
}
 
function repeatGuidFunc(times, delim) {
    if (!delim) delim = '';
    let rtn = "";
    for (let i = 0; i < times; i++) {
        rtn += guidS4() + ((i < (times - 1)) ? delim : "");
    }
    return rtn;
}
 
document.addEventListener("DOMContentLoaded", function () {
    if(!window.util)window.util={};
    util.functionCall=functionCall;
});
 
 
 
 
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
//!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! S.D.G !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
//^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^