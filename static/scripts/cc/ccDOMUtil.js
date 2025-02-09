//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
//!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! J.J. !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
//^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
/*
* Class the replaces the prototype side of utilities.js adding functions to javascript.
* utilities.js is obsolete.
* Public Domain
* Licensed Copyright Law of the United States of America, Section 105 (https://www.copyright.gov/title17/92chap1.html#105)
* Per hoc, facies, scietis quod ille miserit me ut facerem universa quae cernitis et factis: Non est mecum!
* Published by: Dominic Roche of OIT/IUSG/DASM on 03/13/2024
*/
 
"use strict";
import {asBool, append,asElementId, camelToTitle, classAdd, classRemove, css, elementHide, elementRemove,elementShow, endsWith,
    findAllIndexesOf, getEle,
    getEles, getVal, guid, isCheck, isElement, jsonToHtml, makeAlert, off, on, startsWith, toTitle, unbind} from './ccUtilities.js';
 
export default class ccDOMUtil {
}
 
/*=============================Javascript Extension Functions-STRING======================== */
/* ~~~~~~~~ ES6 to Linq */
// Do all of the items match the query
if (typeof Array.prototype.all != "function"){
    Object.defineProperty(Array.prototype, "all", {
        value: function(fn) {
            if (this === undefined || this === null) return false;
            if (!Array.isArray(this)) return false;
            if (this.length === 0) return false;
            if (fn === undefined) fn = ()=> true;
            return this.every(fn);
        }
    });
}
/* ~~~~~~~~ ES6 to Linq */
// Are there and elements with given query in the array
if (typeof Array.prototype.any != "function"){
    Object.defineProperty(Array.prototype, "any", {
        value: function(fn) {
            if (this === undefined || this === null) return false;
            if (!Array.isArray(this)) return false;
            if (this.length === 0) return false;
            if (fn === undefined) fn =()=>true;
            return this.some(fn);
        }
    });
}
//Orders an Array
if (typeof Array.prototype.OrderBy != "function"){
    Object.defineProperty(Array.prototype, "OrderBy", {
        value: function(prop, dir) {
            if (dir === undefined) dir = "asc";
            if (this.length < 1) return this;
            if (typeof this[0] != "object") {
                if (dir.toLowerCase() === "desc") {
                    return this.sort(((a, b) => (a > b ? 1 : a < b ? -1 : 0)) * -1);
                } else {
                    return this.sort((a, b) => (a > b ? 1 : a < b ? -1 : 0));
                }
            }
            return this.sort(orderBy(prop, dir));
        }
    });
}
// Adds an element or html after the parent given (outside)
if (typeof Object.prototype.utlAfter != "function") {
    Object.defineProperty(Object.prototype, "utlAfter", {
        value: function(html) {
            if (html === undefined) {
                return this;
            } else {
                if (typeof html === 'string' || html instanceof String) this.insertAdjacentHTML("afterend", html);
                if (typeof html === 'object' || html instanceof Object) this.insertAdjacentElement("afterend", html);
                return this;
            }
        }
    });
}
// Appends an element or html to a parent given at the bottom (inside)
if (typeof Object.prototype.utlAppend != "function"){
    Object.defineProperty(Object.prototype, "utlAppend", {
        value: function(html) {
            if (html === undefined) {
                return this;
            } else {
                if (typeof html === 'string' || html instanceof String) this.insertAdjacentHTML("beforeend", html);
                if (typeof html === 'object' || html instanceof Object) this.insertAdjacentElement("beforeend", html);
                return this;
            }
            //{ append(this, html); return this; }
        }
    });
}
// Converts the object to a bool
if (typeof String.prototype.utlAsBool != "function") {
    String.prototype.utlAsBool = function() {
        return asBool(this);
    };
}
// Ensures the element id given starts with #
if (typeof String.prototype.utlAsElementId != "function") {
    String.prototype.utlAsElementId = function() {
        return asElementId(this);
    };
}
// Returns a string in camel notation with lowercase start. If it already starts lowercase a _ is prepended
if (typeof String.prototype.utlAsFieldNotationString != "function") {
    String.prototype.utlAsFieldNotationString = function() {
        return asFieldNotationString(this);
    };
}
if (typeof String.prototype.utlAsInt != "function") {
    // Returns a string safely (toString craps out)
        Object.defineProperty(Object.prototype, "utlAsInt", {
            value: function() {
                return asInt(this);
            },
            enumerable: false,
            configurable: true
        });
    }
// Turns a string into Pascal Case
if (typeof String.prototype.utlAsPascalCase != "function") {
    String.prototype.utlAsPascalCase = function(str) {
        return asPascalCase(this, str);
    };
}
// Returns a guid
if (typeof Object.prototype.utlGuid != "function") {
    Object.prototype.utlGuid = function(nodash) {
        return guid(nodash);
    };
}
// Returns a string in camel notation
if (typeof String.prototype.utlAsPropertyNotation != "function") {
    String.prototype.utlAsPropertyNotation = function() {
        return asPropertyNotation(this);
    };
}
if (typeof String.prototype.utlAsString != "function") {
    // Returns a string safely (toString craps out)
    Object.defineProperty(Object.prototype, "utlAsString", {
        value: function() {
            return asString(this || "");
        },
        enumerable: false,
        configurable: true
    });
}
// Adds an element or html before the parent element given
if (typeof Object.prototype.utlBefore != "function") {
    Object.defineProperty(Object.prototype, "utlBefore", {
        value: function(html) {
            if (html === undefined) {
                return this;
            } else {
                if (typeof html === 'string' || html instanceof String) this.insertAdjacentHTML("beforebegin", html);
                if (typeof html === 'object' || html instanceof Object) this.insertAdjacentElement("beforebegin", html);
                return this;
            }
        }
    });
}
// Turns somePascalName to Some Pascal Name
if (typeof String.prototype.utlCamelToTitle != "function") {
    String.prototype.utlCamelToTitle = function(str) {
        return camelToTitle(this, str);
    };
}
// Checks that string ends with the specific string
if (typeof String.prototype.utlChunkString != "function") {
    String.prototype.utlChunkString = function(length) {
        return chunkString(this, length);
    };
}
// adds a class (or multiple classes if you give it a space delimited string)
if (typeof Object.prototype.utlClassAdd != "function") {
    Object.defineProperty(Object.prototype, "utlClassAdd", {
        value: function(className) {
            return classAdd(this,className);
        }
    });
}
// determines if an element has the class or not
if (typeof Object.prototype.utlClassHas != "function") {
    Object.defineProperty(Object.prototype, "utlClassHas", {
        value: function(className) {
            if (this.classList && this.classList.length > 0) {
                return Array.from(this.classList).some(c=>c==className);
            } else {
                return false;
            }
        }
    });
}
//Reomove the specifie class from the element
if (typeof Object.prototype.utlClassRemove != "function"){
    Object.defineProperty(Object.prototype, "utlClassRemove", {
        value: function(className) {
            return classRemove(this, className);
        }
    });
}
//switches one class for another
if (typeof Object.prototype.utlClassSwitch != "function") {
    Object.defineProperty(Object.prototype, "utlClassSwitch", {
        value: function(currClass, newClass) {
            this.utlClassRemove(currClass);
            this.utlClassAdd(newClass);
        }
    });
}
// Clips the string at given max characters
if (typeof String.prototype.utlClip != "function") {
    String.prototype.utlClip = function(maxChar) {
        if(!this || this==="") return "";
        if(!maxChar) return this;
        return this.slice(0,maxChar);
    };
}
// Gets the first element matching your DOM query
if (typeof Object.prototype.utlEle != "function"){
    Object.defineProperty(Object.prototype, "utlEle", {
        value: function(q) {
            let dis = this;
            if (!dis) dis = _;
            if (dis.utlIsElement()) {
                return getEle(dis, q);
            } else {
                return getEle(document, q);
            }
        }
    });
}
// determines if an element exists
if (typeof Object.prototype.utlEleExists != "function") {
    Object.defineProperty(Object.prototype, "utlEleExists", {
        value: function(q) {
            let e=getEle(document, q);
            return !!e && e.length!==0 && e.id!=="ElementNotFound";
        }
    });
}
// Gets all elements matching your DOM query
if (typeof Object.prototype.utlEles != "function"){
    Object.defineProperty(Object.prototype, "utlEles", {
        value: function(q) {
            let dis = this;
            if (!dis) dis = _;
            if (dis.utlIsElement()) {
                return getEles(dis, q);
            } else {
                return getEles(document, q);
            }
        }
    });
}
// Checks that string ends with the specific string
if (typeof String.prototype.utlEndsWith != "function") {
    String.prototype.utlEndsWith = function(str) {
        if(!this) return false;
        return endsWith(this, str, true);
    };
}
// Determined if this DOM Element is present
if (typeof Object.prototype.utlExists != "function"){
    Object.defineProperty(Object.prototype, "utlExists", {
        value: function(q) {
            if(!q){
                if (this.utlIsElement()) {
                    return (this.id!=="ElementNotFound");
                } else{
                    return false;
                }
            } else {
                return (document.getElementById(q));
            }
        }
    });
}
//  Encodes the HTML supplied
if (typeof String.prototype.utlHtmlEncode != "function") {
    String.prototype.utlHtmlEncode = function() {
        return htmlEncode(this);
    };
}
// Fixes image urls
if (typeof String.prototype.utlImageFix != "function") {
    String.prototype.utlImageFix = function() {
        return imageFix(this);
    };
}
// Determined if this object is a DOM Element
if (typeof Object.prototype.utlIsElement != "function"){
    Object.defineProperty(Object.prototype, "utlIsElement", {
        value: function() {
            let e=this;
            if(!e) return false;
            if(e===null || e===undefined) return false;
            return isElement(e);
        },
        enumerable: false
    });
}// Returns a true element has no html
if (typeof String.prototype.utlIsElementEmpty != "function") {
    Object.defineProperty(Object.prototype, "utlIsElementEmpty", {
        value: function() {
            return isElementEmpty(this);
        },
        enumerable: false,
        configurable: true
    });
}
if (typeof String.prototype.utlIsEmpty != "function") {
    // Returns a true if string or Array is empty
    Object.defineProperty(Object.prototype, "utlIsEmpty", {
        value: function() {
            return isEmpty(this);
        },
        enumerable: false,
        configurable: true
    });
}
// Returns a true if string or Array is NOT empty
if (typeof Object.prototype.utlFormOn != "function"){
    Object.defineProperty(Object.prototype, "utlFormOn", {
        value: function(status) {
            return !formEnabled(this,status);
        },
        enumerable: false,
        configurable: true
    })
};
// Moves window to DOM element (If a non-DOM element is used it will move to screen 0,0)
if (typeof Object.prototype.utlGoTo != "function"){
    Object.defineProperty(Object.prototype, "utlGoTo", {
        value: function() {
            if (!this.utlIsElement()) {
                window.scrollTo(0, 0);
            } else {
                let os = this.utlOffset();
                let top = ((os.top-100<0)?50:(os.top-100));
                window.scrollTo(0, top);
            }
        }
    });
}
// Gets the offset information for an element
if (typeof Object.prototype.utlOffset != "function"){
    Object.defineProperty(Object.prototype, "utlOffset", {
        value: function() {
            if (!this.utlIsElement()) {
                return { top: 0, left: 0};
            } else {
                let rect = this.getBoundingClientRect();
                return {
                    top: rect.top + document.body.scrollTop,
                    left: rect.left + document.body.scrollLeft
                };
            }
        }
    });
}
//Gets the height of an element
if (typeof Object.prototype.utlHeight != "function"){
    Object.defineProperty(Object.prototype, "utlHeight", {
        value: function() {
            return parseFloat(getComputedStyle(this, null).height.replace("px", ""));
        }
    });
}
// Hides the element
if (typeof Object.prototype.utlHide != "function") {
    Object.defineProperty(Object.prototype, "utlHide", {
        value: function() {
            return elementHide(this);
        }
    });
}
// Use Converts a DOM Element to JHTML(json)
if (typeof Object.prototype.utlHtmlToJML != "function"){
    Object.defineProperty(Object.prototype, "utlHtmlToJML", {
        value: function() {
            return htmlToJMl(this);
        },
        enumerable: false
    });
}
//Checks values about an element and returns true or false.  Like jQuery .is()
if (typeof Object.prototype.utlIs != "function") {
    Object.defineProperty(Object.prototype, "utlIs", {
        value: function(checkType) {
            return isCheck(this, checkType);
        }
    });
}
//Determines if the element is hidden or not
if (typeof Object.prototype.utlIsHidden != "function") {
    Object.defineProperty(Object.prototype, "utlIsHidden", {
        value: function() {
            return (this.offsetParent === null || this.utlClassHas("d-none"));
        }
    });
}
// Determines if string is html
if (typeof String.prototype.utlIsHtml != "function") {
    Object.defineProperty(Object.prototype, "utlIsHtml", {
        value: function() {
            return isHtml(this);
        },
        enumerable: false,
        configurable: true
    });
}
if (typeof String.prototype.utlIsJson != "function") {
    Object.defineProperty(Object.prototype, "utlIsJson", {
        value: function() {
            return isJson(this);
        },
        enumerable: false,
        configurable: true
    });
}
// Determined if this object is a DOM Node
if (typeof Object.prototype.utlIsNode != "function"){
    Object.defineProperty(Object.prototype, "utlIsNode", {
        value: () =>
            typeof Node === "object"
                ? this instanceof Node
                : this &&
          typeof this === "object" &&
          typeof this.nodeType === "number" &&
          typeof this.nodeName === "string"
    });
}
if (typeof String.prototype.utlIsNotEmpty != "function") {
    // Returns a true if string or Array is NOT empty
    Object.defineProperty(Object.prototype, "utlIsNotEmpty", {
        value: function() {
            return !isEmpty(this);
        },
        enumerable: false,
        configurable: true
    });
}
//Determines if an object is a numeric or not
if (typeof Object.prototype.utlIsNumeric != "function"){
    Object.defineProperty(Object.prototype, "utlIsNumeric", {
        value: function() {
            return isNumeric(this);
        },
        enumerable: false,
        configurable: true
    });
}
// Converts a JHTML(json) object to HTML
if (typeof Object.prototype.utlJsonToHtml != "function"){
    Object.defineProperty(Object.prototype, "utlJsonToHtml", {
        value: function(tabs) {
            return jsonToHtml(this, tabs);
        },
        enumerable: false
    });
}
// removes a function on an element based on event given and function given
if (typeof Object.prototype.utlOff != "function") {
    Object.defineProperty(Object.prototype, "utlOff", {
        value: function(event, func) {
            off(this, event, func);
        }
    });
}
// Sets a function on an element based on event given and function given
if (typeof Object.prototype.utlOn != "function"){
    Object.defineProperty(Object.prototype, "utlOn", {
        value: function(event, func) {
            // if(Array.isArray(this)){
            //     this.forEach(e=> utli.on(e,event,func));
            // }else{
            on(this, event, func);
            //}
        }
    });
}
// Adds an element or html at the top of the parent element given(inside)
if (typeof Object.prototype.utlPrepend != "function") {
    Object.defineProperty(Object.prototype, "utlPrepend", {
        value: function(html) {
            if (html === undefined) {
                return this;
            } else {
                if (typeof html === 'string' || html instanceof String) this.insertAdjacentHTML("afterbegin", html);
                if (typeof html === 'object' || html instanceof Object) this.insertAdjacentElement("afterbegin", html);
                return this;
            }
            //{ append(this, html); return this; }
        }
    });
}
// Gets the previous element
if (typeof Object.prototype.utlPreviousElement != "function") {
    Object.defineProperty(Object.prototype, "utlPreviousElement", {
        value: function() {
            if (!this.utlIsElement()) {
                return _.utlEle("body").utlEle(":not([id=''])");
            } else {
                return prevElement(this);
            }
        }
    });
}
//Gets the previous element with an id attribute set
if (typeof Object.prototype.utlPreviousIdedElement != "function") {
    Object.defineProperty(Object.prototype, "utlPreviousIdedElement", {
        value: function() {
            if (!this.utlIsElement()) {
                return _.utlEle("body").utlEle(":not([id=''])");
            } else {
                return prevElementWithId(this);
            }
        }
    });
}
//Removes the given element from the DOM (like jQuery .remove())
if (typeof Object.prototype.utlRemove != "function"){
    Object.defineProperty(Object.prototype, "utlRemove", {
        value: function() {
            return elementRemove(this);
        }
    });
}
// // Replace all acts like C# Replace
// if (typeof String.prototype.utlReplaceAll != "function") {
//  String.prototype.utlReplaceAll = function(search, replacement, ignore) {
//      return replaceAll(this, search, replacement, ignore);
//  };
// }
// Replace the first occurence of a search
if (typeof String.prototype.utlReplaceFirst != "function") {
    String.prototype.utlReplaceFirst = function(search, replacement) {
        return this.replace(search, replacement);
    };
}
// Replace the last occurence of search string with replacement
if (typeof String.prototype.utlReplaceLast != "function") {
    String.prototype.utlReplaceLast = function(search, replacement) {
        const index = this.lastIndexOf(search);
        if(!replacement) replacement="";
        if (index !== -1) {
            return this.slice(0, index) + replacement + this.slice(index + search.length);
        } else {
            return this;
        }
    };
}
// Shows the element
if (typeof Object.prototype.utlShow != "function"){
    Object.defineProperty(Object.prototype, "utlShow", {
        value: function() {
            return elementShow(this);
        }
    });
}
// Checks that string starts with the specific string
if (typeof String.prototype.utlStartsWith != "function") {
    String.prototype.utlStartsWith = function(str) {
        return startsWith(this, str, true);
    };
}
//Alters the text of an element
if (typeof Object.prototype.utlText != "function"){
    Object.defineProperty(Object.prototype, "utlText", {
        value: function(text) {
            if (text === undefined) {
                return this.textContent;
            } else {
                this.textContent = text;
                return this;
            }
        }
    });
}
//Toggles the element (ie if its hidden show it, if its shown hide it.)
if (typeof Object.prototype.utlToggle != "function") {
    Object.defineProperty(Object.prototype, "utlToggle", {
        value: function() {
            if(this.utlIsHidden()){
                elementShow(this);
            }else{
                elementHide(this);
            }
            return this;
        }
    });
}
// Turns text into title case
if (typeof String.prototype.utlToTitle != "function") {
    String.prototype.utlToTitle = function(str) {
        return toTitle(this, str);
    };
}
 
// Removes extranious periphrial spaces from a string
if (typeof String.prototype.utlTrim != "function") {
    String.prototype.utlTrim = function() {
        return trim(this);
    };
}
//Unbinds all events from the element
if (typeof Object.prototype.utlUnbind != "function"){
    Object.defineProperty(Object.prototype, "utlUnbind", {
        value: function() {
           return unbind(this);
        }
    });
}
//Gets or sets the value of the given element
if (typeof Object.prototype.utlVal != "function"){
    Object.defineProperty(Object.prototype, "utlVal", {
        value: function(val) {
            return getVal(this, val);
        }
    });
}
//Gets the width of an element
if (typeof Object.prototype.utlWidth != "function"){
    Object.defineProperty(Object.prototype, "utlWidth", {
        value: function() {
            return (this===document.body?document.body.clientWidth:parseFloat(getComputedStyle(this, null).width.replace("px", "")));
        }
    });
}
 
/*==============================================Javascript Extension Object Extensions END=============================================== */
 
if (typeof Object.prototype.getOwnPropertyNames != "function") {
        Object.defineProperty(Object.prototype, "getOwnPropertyNames", {
            value: function() {
            return Object.getOwnPropertyNames(this);
        }
    });
}
 
if (typeof Object.prototype.entries != "function") {
        Object.defineProperty(Object.prototype, "entries", {
            value: function() {
            return Object.entries(this);
        }
    });
}
if (typeof Object.prototype.asDictionary != "function") {
        Object.defineProperty(Object.prototype, "asDictionary", {
            value: function() {
            //return Object.entries(this).map(([key, value]) => ({key,value}));
            return Object.entries(this).map(([key, value]) => value);
        }
    });
}
 
/*==============================================Javascript Extension Functions-STRING END=============================================== */
 
//Gets or sets the value of the given element
if (typeof Object.prototype.uVal != "function")
        Object.defineProperty(Object.prototype, "uVal", {
        value: function(val) {
            return getVal(this, val);
        }
    });
 
/*===========END================================Javascript Extension Functions-Simpler Function Names===============END================= */
 
/*================================NOT Mine================================ */
//by Leonardo Filipe
Object.defineProperty(Array.prototype, 'sortIt', {
    value: function(sorts) {
        sorts.map(sort => {            
            sort.uniques = Array.from(
                new Set(this.map(obj => obj[sort.key]))
            );
           
            sort.uniques = sort.uniques.sort((a, b) => {
                if (typeof a == 'string') {
                    return sort.inverse ? b.localeCompare(a) : a.localeCompare(b);
                }
                else if (typeof a == 'number') {
                    return sort.inverse ? b - a : a - b;
                }
                else if (typeof a == 'boolean') {
                    let x = sort.inverse ? (a === b) ? 0 : a? -1 : 1 : (a === b) ? 0 : a? 1 : -1;
                    return x;
                }
                return 0;
            });
        });
   
        const weightOfObject = (obj) => {
            let weight = "";
            sorts.map(sort => {
                let zeropad = `${sort.uniques.length}`.length;
                weight += sort.uniques.indexOf(obj[sort.key]).toString().padStart(zeropad, '0');
            });
            //obj.weight = weight; // if you need to see weights
            return weight;
        }
   
        this.sort((a, b) => {
            return weightOfObject(a).localeCompare( weightOfObject(b) );
        });
       
        return this;
    }
    });



//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
//!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! S.D.G !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
//^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^