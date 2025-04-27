package ggw_search

import (
    //"fmt"
    //"strings"
    //"regexp"
)

var ModifyMatch func(match , replacementWord string) string;

type Search struct{
    TargetWords []string; // Sifting through data, to instantly search and add the path to the searchable path list.
    ReplacementWords map[string]string; // `(Regex)TargetWord` : "to raplace with" ,
    IsNotRegex bool; // TargetWord is not a regex
    FullSearch bool; // use 'ReplacementWords' keys to search in all lines (and replace or not)
    OnlyGetMatchedClauses bool; // to get all matched items without any replacement
    GetMatchedPerFile bool; // to show matched clauses separately in each file, or only return one as an instance
    IsReplaceAll bool; // Should be set to true to complete the search based on 'ReplacementWords' keys. handle no replacement by OnlyGetMatchedClauses
    MixWithOriginal string; // replacement mode: "replace","after","before","continue"(skip, for testing or ...)
    OutPutMode string; // "println","showPaths","noOut"
    IsPrintToFile bool; // print to file (for big outputs), or print to stdout
    ResultFilePath string; // file path to print to file 
    ValidExtensions []string; // set {"all"} to serach in all files
    IsStaticDIR bool; // search a whole directory, or specific paths?
    StaticDIRs []string; // specify your Static DIRs to avoid dynamic recursive search
    DIRs []string; // specify directories to search dynamic recursive search
    DirToExclude []string ; // specify part of directories to Exclude searching
    HasCondition bool; // whether replacement has condition, if yes, IsReplaceAll should be set to false
    HawMany int; // for conditional replace, set it -1 to do for all matches
    StartFromNumber int; // for conditional replace, replace from which matched item?
    RegxToExistBefore string; // for conditional replace
    RegxToExistAfter string; // for conditional replace
    InsideCond bool; // for conditional replace
    RegxToExistInside []string; // for conditional replace
    Regx_NOT_ExistInside []string; // for conditional replace
    Regx_NOT_ExistBefore string; // for conditional replace
    Regx_NOT_ExistAfter string; // for conditional replace
    IsModifyMatch bool; // to run your custom code (the ModifyMatch function) for each matched
    ModifyMatch func(match , replacementWord string) string; // run your custom code for each matched
} 

//Replacement Words
var (

    TargetWords []string =[]string{
        // Sifting through data, to instantly search and add the path to the searchable path list.

        // "ping"  // serach in all specified files that contain the word "ping"
        // "_all_" ,  // serach in all specified files
    } ; 

    ReplacementWords map[string]string= map[string]string{
        // `RegxTargetWord` : "to raplace with" ,
    };
)
//settings
var(
    FullSearch bool = true ; // use 'ReplacementWords' keys to search

    // replacing confings
    // Be careful, there's no Ctrl+Z for replacement!
    
    OnlyGetMatchedClauses bool = true ; // to get all matched items without any replacement
    GetMatchedPerFile bool = true ; // to show matched clauses separately in each file

    IsReplaceAll bool = true ; // Should be set to true to complete the search based on 'ReplacementWords' keys.

    // If IsReplaceAll is false , come here :
    HawMany int = 1 ; // set it -1 to do for all matches
    StartFromNumber int = 2 ; // also contains this index
    HasCondition bool = false ; // set condition if IsReplaceAll is false

    // to mix replacement with original text
    MixWithOriginal = "" ; //: replace after before continue

    IsModifyMatch bool = false ; // to run modifyMatch func for each matched
    
    OutPutMode string = "println"; //: println showPaths noOut
    IsPrintToFile bool = false ;
    ResultFilePath string = "./ggw_search_result.log" ;

    ValidExtensions []string = []string{".txt"} // set {"all"} to open all files

    IsNotRegex bool = false ;
)
//DIRs
var(
    IsStaticDIR bool = false ; // set StaticDIRs to avoid dynamic recursive search

    DIRs []string = []string{"./"}  ; // paths to dynamic recursive search

    DirToExclude []string = []string{
    };  

    StaticDIRs []string = []string{
    }; 
)
//if HasCondition is true
var (
    
    RegxToExistBefore string = `` ; // ignore if == ""
    RegxToExistAfter string = `` ; // ignore if == ""

    InsideCond bool = false ;
    RegxToExistInside []string = []string{``} ; // should set RegxToExistBefore &&  RegxToExistAfter
    Regx_NOT_ExistInside []string = []string{} ;  // should set RegxToExistBefore &&  RegxToExistAfter

    Regx_NOT_ExistBefore string = "" ; // ignore if == ""
    Regx_NOT_ExistAfter string = "" ; // ignore if == ""
)  
